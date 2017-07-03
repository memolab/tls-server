package middlewares

import (
	"fmt"
	"net/http"
	"runtime/debug"
	"strings"
	"time"

	"go.uber.org/zap"

	"tls-server/api/types"

	"gopkg.in/mgo.v2/bson"
)

// AResponseWriter custom impl of HttpResponseWriter
type AResponseWriter struct {
	http.ResponseWriter
	status int
	length int
}

// Header impl of httpResponseWriter
func (rw *AResponseWriter) Header() http.Header {
	return rw.ResponseWriter.Header()
}

// Header impl of httpResponseWriter
func (rw *AResponseWriter) Write(b []byte) (i int, e error) {
	i, e = rw.ResponseWriter.Write(b)
	rw.length = i
	return
}

// WriteHeader impl of httpResponseWriter
func (rw *AResponseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}

type FrontMiddleware struct {
	ctl        types.APICTL
	AccessLogs []AccessLog
	accessChan chan AccessLog
	stopLog    chan bool
}

type AccessLog struct {
	ID              bson.ObjectId `bson:"_id,omitempty"`
	RemoteAddr      string        `bson:"RemoteAddr"`
	ReqContentType  string        `bson:"ReqContentType"`
	RespContentType string        `bson:"RespContentType"`
	ReqLength       int           `bson:"ReqLength"`
	RespLength      int           `bson:"RespLength"`
	Status          int           `bson:"Status"`
	Path            string        `bson:"Path"`
	Method          string        `bson:"Method"`
	Cached          string        `bson:"Cached"`
	// get time in Nanosecond
	HandlersDuration time.Duration `bson:"HandlersDuration"`
	Timed            time.Time     `bson:"Timed"`
}

func NewFrontMiddleware(ctl types.APICTL, configAccessLogsDump string) *FrontMiddleware {
	dumpDuration, err := time.ParseDuration(configAccessLogsDump)
	if err != nil {
		ctl.Log().Error("FrontMiddleware: error parsing accessLogsDumpConf duration", zap.Error(err))
	}

	front := &FrontMiddleware{
		ctl:        ctl,
		accessChan: make(chan AccessLog, 200),
		stopLog:    make(chan bool),
	}

	go func() {
		ticker := time.NewTicker(dumpDuration)
		defer ticker.Stop()
		for {
			select {
			case a := <-front.accessChan:
				front.AccessLogs = append(front.AccessLogs, a)
			case <-ticker.C:
				front.dumpLogs()
			case <-front.stopLog:
				return
			}
		}
	}()

	return front
}

func (front *FrontMiddleware) dumpLogs() {
	if !(len(front.AccessLogs) > 0) {
		return
	}
	dbc := dumpDB.Copy()
	defer dbc.Close()

	uc := dbc.DB("").C("accessLogs")
	bulk := uc.Bulk()
	bulk.Unordered()
	for _, v := range front.AccessLogs {
		bulk.Insert(v)
	}

	if _, err := bulk.Run(); err != nil {
		front.ctl.Log().Error("FrontMiddleware: error db dumpLogs", zap.Error(err))
	}

	front.AccessLogs = nil
}

func (front *FrontMiddleware) Handler() types.MiddlewareHandler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(_rw http.ResponseWriter, r *http.Request) {
			sTime := time.Now().UTC()

			rw := &AResponseWriter{ResponseWriter: _rw, status: -1}
			rw.Header().Set("Cache-Control", "no-cache, private")
			rw.Header().Set("expires", "-1")
			rw.Header().Set("Vary", "Accept-Encoding")
			rw.Header().Set("Access-Control-Allow-Origin", "*")
			rw.Header().Set("X-Content-Type-Options", "nosniff")
			rw.Header().Set("x-frame-options", "SAMEORIGIN")
			rw.Header().Set("x-xss-protection", "1; mode=block")

			defer func() {
				if err := recover(); err != nil {
					front.ctl.Log().Error("-- RECOVER PANIC --", zap.Any("PANIC ERR", err))
					fmt.Printf("\x1b[31;1m -- RECOVER PANIC Stack-- %s \x1b[0m\n", err)
					fmt.Printf("\x1b[31;1m -- %s --- \x1b[0m\n", debug.Stack())
					fmt.Println("\x1b[31;1m -- --- --- \x1b[0m")
					http.Error(rw, "Internal Server Error", 500)
				}
			}()

			next.ServeHTTP(rw, r)

			front.accessChan <- AccessLog{RemoteAddr: strings.Split(r.RemoteAddr, ":")[0],
				ReqContentType:   r.Header.Get("Content-Type"),
				RespContentType:  rw.Header().Get("Content-Type"),
				ReqLength:        int(r.ContentLength),
				RespLength:       rw.length,
				Status:           rw.status,
				Path:             r.URL.RequestURI(),
				Method:           r.Method,
				Cached:           rw.Header().Get("X-Cache"),
				HandlersDuration: time.Since(sTime),
				Timed:            sTime}

			front.ctl.Log().Info(r.URL.RequestURI(), zap.Int("status", rw.status),
				zap.Duration("HandlersDuration", time.Since(sTime)))
		})
	}
}

func (front *FrontMiddleware) logInfo() {
	front.ctl.Log().Info("FrontMiddleware Log:", zap.Any("clients", front.AccessLogs))
}

func (front *FrontMiddleware) Shutdown() {
	front.logInfo()
	front.stopLog <- true
}
