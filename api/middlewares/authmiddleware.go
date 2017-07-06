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

// MiddlewareAuth a middleware token authorization
// using securecookie.SecureCookie for encode/decode
type AuthMiddleware struct {
	ctl            types.APICTL
	headerTokenKey string
	rateLimit      *rateLimiter
	scCookie       *securecookie.SecureCookie
	dumpDuration   time.Duration
	stopLog        chan struct{}
}

type rateLimiter struct {
	duration time.Duration
	count    int
	clients  map[string]rateLimiterClient
	syncMtx  *sync.Mutex
}
type rateLimiterClient struct {
	time      time.Time
	count     int
	overCount int
	TID       string
}
type clientLog struct {
	ID        bson.ObjectId `bson:"_id,omitempty"`
	UID       string        `bson:"UID"`
	TypeID    string        `bson:"TypeID"`
	Timed     time.Time     `bson:"Timed"`
	Count     int           `bson:"Count"`
	OverCount int           `bson:"OverCount"`
}

func NewAuthMiddleware(ctl types.APICTL, configHeaderTokenKey string,
	configRateLimiteAPI string, configSecretKey1 string, configSecretKey2 string, configRateLimiteLogsDump string) *AuthMiddleware {

	scCookie := securecookie.New([]byte(configSecretKey1), []byte(configSecretKey2))
	scCookie.MaxAge(0)

	rateLimiteConf := strings.Split(configRateLimiteAPI, ":")
	c, erri := strconv.Atoi(rateLimiteConf[0])
	d, errs := time.ParseDuration(rateLimiteConf[1])
	rateLimit := &rateLimiter{duration: d, count: c, clients: map[string]rateLimiterClient{}, syncMtx: &sync.Mutex{}}
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
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				auth.dumpLogs()
			case <-auth.stopLog:
				auth.dumpLogs()
				return
			}
		}
	}()

	return auth
}

func (auth *AuthMiddleware) dumpLogs() {
	auth.rateLimit.syncMtx.Lock()
	defer auth.rateLimit.syncMtx.Unlock()

	clientsLog := []clientLog{}
	for k, v := range auth.rateLimit.clients {
		if time.Now().Sub(v.time) > auth.dumpDuration {
			clientsLog = append(clientsLog, clientLog{
				UID: k, TypeID: v.TID, Timed: v.time, Count: v.count, OverCount: v.overCount,
			})
			delete(auth.rateLimit.clients, k)
		}
	}

	if !(len(clientsLog) > 0) {
		return
	}

	dbc := dumpDB.Copy()
	defer dbc.Close()

	uc := dbc.DB("").C("rateLimitLogs")
	bulk := uc.Bulk()
	bulk.Unordered()
	for _, v := range clientsLog {
		bulk.Insert(v)
	}

	if _, err := bulk.Run(); err != nil {
		auth.ctl.Log().Error("AuthMiddleware: error db dumpLogs", zap.Error(err))
	}
}

// Handler AuthMiddleware http handler
func (auth *AuthMiddleware) Handler() types.MiddlewareHandler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			uid := ""

			if err := auth.ParseSecretToken(r, &uid); err != nil {
				auth.ctl.Abort(rw, http.StatusForbidden)
				return
			}

			if uid != "" && bson.IsObjectIdHex(uid) {
				if auth.checkRateLimit(uid, auth.rateLimit.duration, auth.rateLimit.count, "tkn") {
					auth.ctl.Abort(rw, http.StatusTooManyRequests)
					return
				}
				next.ServeHTTP(rw, r.WithContext(context.WithValue(r.Context(), types.CTXKey("uid"), uid)))
				return
			}

			auth.ctl.Abort(rw, http.StatusForbidden)
		})
	}
}

// ParseSecretToken parse header token value
func (auth *AuthMiddleware) ParseSecretToken(r *http.Request, tokenVal *string) error {
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
func (auth *AuthMiddleware) NewSecretToken(tokenVal string) (encoded string, err error) {
	encoded, err = auth.scCookie.Encode("i", tokenVal)
	return
}

func (auth *AuthMiddleware) checkRateLimit(uid string, limitDuration time.Duration, limitCount int, tid string) bool {
	re := false
	t := time.Now()

	auth.rateLimit.syncMtx.Lock()
	defer auth.rateLimit.syncMtx.Unlock()

	cu := auth.rateLimit.clients[uid]

	if t.Sub(cu.time) <= limitDuration {
		if cu.count >= limitCount {
			re = true
			cu.overCount++
		} else {
			cu.count++
		}
	} else {
		cu.count = 1
	}

	cu.time = t
	cu.TID = tid
	auth.rateLimit.clients[uid] = cu

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
func (auth *AuthMiddleware) Close() {
	auth.stopLog <- struct{}{}
}
