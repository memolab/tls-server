package api

import (
	"encoding/json"
	"net/http"
	"time"

	"go.uber.org/zap"

	"net/http/pprof"

	"tls-server/api/middlewares"
	"tls-server/api/types"

	"gopkg.in/mgo.v2"
)

type APICtl struct {
	mongo   *mgo.Session
	regMidd map[string]types.Middleware
	log     *zap.Logger
}

var (
	shutdownAPI func(err error)
)

// InitAPI setup api functions and http handlers
// return muxServe
func InitAPI(config map[string]string) *http.ServeMux {
	log, _ := newLogger(config["prod"])

	mgoConn, err := newMongo(config["mongoURL"])
	if err != nil {
		log.Error("mongoDB conn", zap.Error(err))
	}

	c := &APICtl{
		mongo:   mgoConn,
		regMidd: map[string]types.Middleware{},
		log:     log,
	}

	middlewares.InitGlobalDumpDB(config["dumpDB"])

	// new FrontMiddleware
	frontMiddleware := middlewares.NewFrontMiddleware(c, config["accessLogsDump"])
	middFront := frontMiddleware.Handler()
	c.regMidd["front"] = frontMiddleware

	// new CacheMiddleware
	cacheMiddleware := middlewares.NewCacheMiddleware(c)
	c.regMidd["cache"] = cacheMiddleware

	// new AuthMiddleware
	authMiddleware := middlewares.NewAuthMiddleware(c, config["headerTokenKey"], config["rateLimiteAPI"],
		config["secretKey1"], config["secretKey2"], config["rateLimiteLogsDump"])
	c.regMidd["auth"] = authMiddleware
	middAuth := authMiddleware.Handler()
	middRateLimitIP := authMiddleware.RateLimitIPHandler(config["headerClientIPs"], config["rateLimiteIP"])

	mux := http.NewServeMux()
	d5min, _ := time.ParseDuration("5m")

	// register handlers
	mux.Handle("/", adapt(http.HandlerFunc(c.indexHanler), middFront))
	mux.Handle("/initdb", adapt(http.HandlerFunc(c.initDBHanler), middAuth, middFront))

	mux.Handle("/user", adapt(http.HandlerFunc(c.userIndexHanler),
		cacheMiddleware.CacheHandler("/user", map[string]string{"tokenUID": ""}, d5min),
		middAuth, middFront))
	mux.Handle("/user2", adapt(http.HandlerFunc(c.user2IndexHanler),
		cacheMiddleware.CacheHandler("/user2", map[string]string{"tokenUID": ""}, d5min),
		middAuth, middFront))

	mux.Handle("/signin", adapt(http.HandlerFunc(c.signInHandler), middRateLimitIP, middFront))
	mux.Handle("/signup", adapt(http.HandlerFunc(c.signUpHandler), middRateLimitIP, middFront))
	//

	// register admin handlers
	isAdmin := middlewares.IsAdmin(c, config["adminID"])
	mux.Handle("/admin", adapt(http.HandlerFunc(c.adminIndexHanler), isAdmin, middAuth, middFront))
	//mux.Handle("/dashboard/", http.StripPrefix("/dashboard", http.FileServer(http.Dir("public/"))))

	//

	if config["pprof"] == "1" {
		mux.HandleFunc("/debug/pprof/", pprof.Index)
		mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
		mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
		mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	}

	shutdownAPI = func(err error) {
		if c.mongo != nil {
			c.mongo.Close()
		}

		for _, m := range c.regMidd {
			m.Shutdown()
		}

		if err != nil {
			c.log.Warn("SHUTDOWN Error", zap.Error(err))
		}
		c.log.Warn("SHUTDOWN.")
	}

	return mux
}

// ShutdownAPI call on server about to close to free any resources
func ShutdownAPI(err error) {
	shutdownAPI(err)
}

func adapt(h http.Handler, handlers ...types.MiddlewareHandler) http.Handler {
	for _, handle := range handlers {
		h = handle(h)
	}
	return h
}

func (c *APICtl) RespJson(rw http.ResponseWriter, status int, data interface{}) {
	rw.Header().Set("Content-Type", "application/json")

	b, jerr := json.Marshal(data)
	if jerr != nil {
		c.log.Error("Error marshal json response", zap.Error(jerr))
		http.Error(rw, "Internal Server Error", 500)
		return
	}

	rw.WriteHeader(status)
	if _, werr := rw.Write(b); werr != nil {
		c.log.Error("Error json response writer", zap.Error(werr))
	}
}

func (c *APICtl) Abort(rw http.ResponseWriter, status int) {
	//c.RespJson(rw, sts, map[string]interface{}{"code": sts, "msg": http.StatusText(sts)})
	rw.Header().Set("Content-Type", "text/plain")
	rw.WriteHeader(status)
	if _, werr := rw.Write([]byte(http.StatusText(status))); werr != nil {
		c.log.Error("Error json response writer", zap.Error(werr))
	}
}

func (c *APICtl) Log() *zap.Logger {
	return c.log
}

func newMongo(url string) (*mgo.Session, error) {
	sess, err := mgo.DialWithTimeout(url, 3*time.Second)
	if err != nil {
		return nil, err
	}

	sess.SetMode(mgo.Monotonic, true)

	return sess, nil
}

func newLogger(prod string) (logger *zap.Logger, err error) {
	if prod == "1" {
		logger, err = zap.NewProduction()
	} else {
		logger, err = zap.NewDevelopment()
	}

	defer logger.Sync()

	return
}

func createMongoIndexs(mongo *mgo.Session) error {
	usrIndx := mgo.Index{
		Key:        []string{"email"},
		Unique:     true,
		Background: true,
		Sparse:     true,
	}
	if err := mongo.DB("").C("users").EnsureIndex(usrIndx); err != nil {
		return err
	}

	return nil
}
