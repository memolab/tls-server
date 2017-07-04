package api

import (
	"encoding/json"
	"net/http"
	"time"

	"go.uber.org/zap"

	"tls-server/api/middlewares"
	"tls-server/api/types"

	"gopkg.in/mgo.v2"
)

type (
	// APICtl main common handlers struct
	APICtl struct {
		mongo   *mgo.Session
		regMidd map[string]types.Middleware
		log     *zap.Logger
		auth    *middlewares.AuthMiddleware
		cache   *middlewares.CacheMiddleware
	}

	route struct {
		url         string
		handler     http.HandlerFunc
		middlewares []types.MiddlewareHandler
	}
)

var (
	shutdownAPI func(err error)
)

// InitAPI setup api functions and http handlers
// return muxServe
func InitAPI(config map[string]string) *http.ServeMux {
	var log *zap.Logger

	if config["prod"] == "1" {
		log, _ = zap.NewProduction()
	} else {
		log, _ = zap.NewDevelopment()
	}

	mgoConn, err := newMongo(config["mongoURL"])
	if err != nil {
		log.Fatal("mongoDB conn", zap.Error(err))
	}

	c := &APICtl{
		mongo:   mgoConn,
		regMidd: map[string]types.Middleware{},
		log:     log,
	}

	middlewares.InitGlobalDumpDB(config["dumpDB"])

	mux := http.NewServeMux()

	for _, r := range *initRoutes(c, config) {
		mux.Handle(config["apiPrefix"]+r.url, applyMiddlewares(r.handler, r.middlewares...))
	}

	shutdownAPI = func(err error) {
		if c.mongo != nil {
			c.mongo.Close()
		}

		c.auth.Shutdown()
		c.cache.Shutdown()
		for _, m := range c.regMidd {
			m.Shutdown()
		}

		if err != nil {
			log.Warn("server shutdown error", zap.Error(err))
		}

		log.Warn("SHUTDOWN.")
		log.Sync()
	}

	return mux
}

// ShutdownAPI call on server about to close to free any resources
func ShutdownAPI(err error) {
	shutdownAPI(err)
}

func applyMiddlewares(h http.Handler, handlers ...types.MiddlewareHandler) http.Handler {
	for _, handle := range handlers {
		h = handle(h)
	}
	return h
}

// RespJSON responce json content type
func (c *APICtl) RespJSON(rw http.ResponseWriter, status int, data interface{}) {
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

// Abort abort status responce
func (c *APICtl) Abort(rw http.ResponseWriter, status int) {
	rw.Header().Set("Content-Type", "text/plain")
	rw.WriteHeader(status)
	if _, werr := rw.Write([]byte(http.StatusText(status))); werr != nil {
		c.log.Error("Error json response writer", zap.Error(werr))
	}
}

// Log return logger
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
