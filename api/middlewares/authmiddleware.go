package middlewares

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"go.uber.org/zap"

	"tls-server/api/types"

	"github.com/gorilla/securecookie"
	"gopkg.in/mgo.v2/bson"
)

// AuthMiddleware a middleware token authorization
// using securecookie.SecureCookie for encode/decode
type AuthMiddleware struct {
	ctl            types.APICTL
	headerTokenKey string
	rateLimit      *rateLimiter
	scCookie       *securecookie.SecureCookie
	dumpDuration   time.Duration
	stopLog        chan struct{}
	Closed         *sync.WaitGroup
}

type rateLimiter struct {
	duration time.Duration
	count    int
	clients  map[string]*rateLimiterClient
	syncMtx  *sync.Mutex
}
type rateLimiterClient struct {
	time      time.Time
	count     int
	overCount int
	tid       string
}

func NewAuthMiddleware(ctl types.APICTL, configHeaderTokenKey string,
	configRateLimiteAPI string, configSecretKey1 string, configSecretKey2 string, configRateLimiteLogsDump string) *AuthMiddleware {

	scCookie := securecookie.New([]byte(configSecretKey1), []byte(configSecretKey2))
	scCookie.MaxAge(0)
	scCookie.SetSerializer(securecookie.NopEncoder{})

	rateLimiteConf := strings.Split(configRateLimiteAPI, ":")
	c, erri := strconv.Atoi(rateLimiteConf[0])
	d, errs := time.ParseDuration(rateLimiteConf[1])
	rateLimit := &rateLimiter{duration: d, count: c,
		clients: map[string]*rateLimiterClient{},
		syncMtx: &sync.Mutex{}}
	if errs != nil || erri != nil {
		ctl.Log().Error("AuthMiddleware: invalid RateLimit config", zap.Errors("erri,errs", []error{erri, errs}))
	}

	dumpDuration, errd := time.ParseDuration(configRateLimiteLogsDump)
	if errd != nil {
		ctl.Log().Error("AuthMiddleware: invalid configRateLimiteLogsDump config", zap.Error(errd))
	}

	auth := &AuthMiddleware{
		ctl:            ctl,
		headerTokenKey: strings.Title(strings.ToLower(configHeaderTokenKey)),
		rateLimit:      rateLimit,
		scCookie:       scCookie,
		dumpDuration:   dumpDuration,
		stopLog:        make(chan struct{}),
	}

	go func() {
		ticker := time.NewTicker(dumpDuration)
		defer func() {
			ticker.Stop()
			auth.Closed.Done()
		}()
		for {
			select {
			case <-ticker.C:
				auth.dumpLogs(false)
			case <-auth.stopLog:
				close(auth.stopLog)
				auth.dumpLogs(true)
				return
			}
		}
	}()

	return auth
}

func (auth *AuthMiddleware) dumpLogs(force bool) {
	if !(len(auth.rateLimit.clients) > 0) {
		return
	}

	clientsLog := map[string]*rateLimiterClient{}
	auth.rateLimit.syncMtx.Lock()
	for k, v := range auth.rateLimit.clients {
		if force || time.Now().UTC().Sub(v.time) >= auth.dumpDuration {
			clientsLog[k] = v
			delete(auth.rateLimit.clients, k)
		}
	}
	auth.rateLimit.syncMtx.Unlock()

	if !(len(clientsLog) > 0) {
		return
	}

	dbc := dumpDB.Copy()
	defer dbc.Close()

	uc := dbc.DB("").C("rateLimitLogs")
	bulk := uc.Bulk()
	bulk.Unordered()

	for uid, v := range clientsLog {
		y, m, d := v.time.Date()
		D := time.Date(y, m, d, 0, 0, 0, 0, time.UTC)
		bulk.Upsert(bson.M{"UID": uid, "Timed": D},
			bson.M{
				"$set": bson.M{"TypeID": v.tid},
				"$inc": bson.M{"Count": v.count, "OverCount": v.overCount},
			},
		)
		delete(clientsLog, uid)
	}

	if _, err := bulk.Run(); err != nil {
		auth.ctl.Log().Error("AuthMiddleware: error db dumpLogs", zap.Error(err))
	}
}

// Handler AuthMiddleware http handler
func (auth *AuthMiddleware) Handler() types.MiddlewareHandler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			uid := make([]byte, 24)

			if err := auth.ParseSecretToken(r, &uid); err != nil {
				auth.ctl.Abort(rw, http.StatusForbidden)
				return
			}
			uidStr := string(uid)
			if bson.IsObjectIdHex(uidStr) {
				*r = *r.WithContext(context.WithValue(r.Context(), types.CTXUIDKey{}, uidStr))
				if auth.checkRateLimit(uidStr, auth.rateLimit.duration, auth.rateLimit.count, "tkn") {
					auth.ctl.Abort(rw, http.StatusTooManyRequests)
					return
				}
				next.ServeHTTP(rw, r)
				return
			}

			auth.ctl.Abort(rw, http.StatusForbidden)
		})
	}
}

// ParseSecretToken parse header token value
func (auth *AuthMiddleware) ParseSecretToken(r *http.Request, tokenVal *[]byte) error {
	token := r.Header.Get(auth.headerTokenKey)
	if token == "" {
		token = r.URL.Query().Get(auth.headerTokenKey)
	}

	if token == "" {
		return errors.New("empty token")
	}

	return auth.scCookie.Decode("i", token, tokenVal)
}

// NewSecretToken generate encoded token
func (auth *AuthMiddleware) NewSecretToken(tokenVal []byte) (encoded string, err error) {
	encoded, err = auth.scCookie.Encode("i", tokenVal)
	return
}

func (auth *AuthMiddleware) checkRateLimit(uid string, limitDuration time.Duration, limitCount int, tid string) bool {
	re := false
	t := time.Now().UTC()

	auth.rateLimit.syncMtx.Lock()
	defer auth.rateLimit.syncMtx.Unlock()

	cu := auth.rateLimit.clients[uid]
	if cu == nil {
		auth.rateLimit.clients[uid] = &rateLimiterClient{tid: tid, time: t}
		cu = auth.rateLimit.clients[uid]
	}

	if t.Sub(cu.time) <= limitDuration {
		if cu.count >= limitCount {
			re = true
			cu.overCount++
		} else {
			cu.count++
		}
	} else {
		cu.count = 1
		cu.time = t
	}
	return re
}

// RateLimitIPHandler limit access by client ip
func (auth *AuthMiddleware) RateLimitIPHandler(headerkeysConf string, rateLimiteConfStr string) types.MiddlewareHandler {
	//headerkeys := strings.Split(headerkeysConf, ",")
	rateLimiteConf := strings.Split(rateLimiteConfStr, ":")
	c, erri := strconv.Atoi(rateLimiteConf[0])
	d, errs := time.ParseDuration(rateLimiteConf[1])
	if errs != nil || erri != nil {
		auth.ctl.Log().Error("call rateLimitIPHandler with Invalid config", zap.Error(erri), zap.Error(errs))
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			ip := strings.Split(r.RemoteAddr, ":")[0]
			//ip := utils.GetIPAdress(r, headerkeys)
			if ip != "" && auth.checkRateLimit(ip, d, c, "ip") {
				auth.ctl.Abort(rw, http.StatusNotAcceptable)
				return
			}

			next.ServeHTTP(rw, r)
		})
	}
}

// LogInfo log all pinding data
func (auth *AuthMiddleware) LogInfo() {
	auth.ctl.Log().Info("AuthMiddleware Log:", zap.Any("clients", auth.rateLimit.clients))
}

// Close stop dump data
func (auth *AuthMiddleware) Close(wg *sync.WaitGroup) {
	wg.Add(1)
	auth.Closed = wg
	auth.stopLog <- struct{}{}
}
