package middlewares

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
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
}

type rateLimiter struct {
	duration time.Duration
	count    int
	clients  map[string]rateLimiterClient
	syncMtx  *sync.Mutex
}
type rateLimiterClient struct {
	time  time.Time
	count int
}

func NewAuthMiddleware(ctl types.APICTL, configHeaderTokenKey string,
	configRateLimiteAPI string, configSecretKey1 string, configSecretKey2 string) *AuthMiddleware {

	scCookie := securecookie.New([]byte(configSecretKey1), []byte(configSecretKey2))
	scCookie.MaxAge(0)

	rateLimiteConf := strings.Split(configRateLimiteAPI, ":")
	i, erri := strconv.Atoi(rateLimiteConf[0])
	t, errs := time.ParseDuration(rateLimiteConf[1])
	rateLimit := &rateLimiter{duration: t, count: i, clients: map[string]rateLimiterClient{}, syncMtx: &sync.Mutex{}}
	if errs != nil || erri != nil {
		ctl.Log().Error("AuthMiddleware: invalid RateLimit config", "erri", erri, "errs", errs)
	}

	return &AuthMiddleware{
		ctl:            ctl,
		headerTokenKey: configHeaderTokenKey,
		rateLimit:      rateLimit,
		scCookie:       scCookie,
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
				if auth.checkRateLimit(uid, auth.rateLimit.duration, auth.rateLimit.count) {
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

func (auth *AuthMiddleware) checkRateLimit(uid string, limitDuration time.Duration, limitCount int) bool {
	re := false
	t := time.Now()

	auth.rateLimit.syncMtx.Lock()
	defer auth.rateLimit.syncMtx.Unlock()

	cu := auth.rateLimit.clients[uid]
	cud := t.Sub(cu.time)

	if cud <= limitDuration {
		if cu.count >= limitCount {
			re = true
		}
	} else {
		cu.count = 0
	}

	cu.count++
	cu.time = t
	auth.rateLimit.clients[uid] = cu

	return re
}

func (auth *AuthMiddleware) RateLimitIPHandler(headerkeysConf string, rateLimiteConfStr string) types.MiddlewareHandler {
	//headerkeys := strings.Split(headerkeysConf, ",")
	rateLimiteConf := strings.Split(rateLimiteConfStr, ":")
	c, erri := strconv.Atoi(rateLimiteConf[0])
	d, errs := time.ParseDuration(rateLimiteConf[1])
	if errs != nil || erri != nil {
		auth.ctl.Log().Error("call rateLimitIPHandler with Invalid config", "erri", erri, "errs", errs)
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			ip := strings.Split(r.RemoteAddr, ":")[0]
			//ip := utils.GetIPAdress(r, headerkeys)
			if ip != "" && auth.checkRateLimit(ip, d, c) {
				auth.ctl.Abort(rw, http.StatusNotAcceptable)
				return
			}

			next.ServeHTTP(rw, r)
		})
	}
}

func (auth *AuthMiddleware) logInfo() {
	auth.ctl.Log().Info("rateLimit Log:", "clients", auth.rateLimit.clients)
}

func (auth *AuthMiddleware) Shutdown() {
	auth.ctl.Log().Info("rateLimit Log:", "clients", auth.rateLimit.clients)
}