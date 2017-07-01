package middlewares

import (
	"fmt"
	"gopkg.in/mgo.v2/bson"
	"net/http"
	"runtime/debug"
	"strconv"
	"time"
	"tls-server/api/types"
)

// AResponseWriter custom impl of HttpResponseWriter
type AResponseWriter struct {
	http.ResponseWriter
	status int
}

// Header impl of httpResponseWriter
/*func (rw *AResponseWriter) Header() http.Header {
	return rw.ResponseWriter.Header()
}*/

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
	ID               bson.ObjectId `bson:"_id,omitempty"`
	RemoteAddr       string        `bson:"RemoteAddr"`
	ReqContentType   string        `bson:"ReqContentType"`
	RespContentType  string        `bson:"RespContentType"`
	ReqLength        int64         `bson:"ReqLength"`
	RespLength       int64         `bson:"RespLength"`
	Status           int           `bson:"Status"`
	Path             string        `bson:"Path"`
	Method           string        `bson:"Method"`
	Cached           string        `bson:"Cached"`
	// get time in Nanosecond
	HandlersDuration time.Duration `bson:"HandlersDuration"`
	Timed            time.Time     `bson:"Timed"`
}

func NewFrontMiddleware(ctl types.APICTL, configAccessLogsDump string) *FrontMiddleware {
	dumpDuration, err := time.ParseDuration(configAccessLogsDump)
	if err != nil {
		ctl.Log().Error("FrontMiddleware: error parsing accessLogsDumpConf duration", "err", err)
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
	_, err := bulk.Run()
	if err != nil {
		front.ctl.Log().Error("FrontMiddleware: error db dumpLogs", "err", err)
	}

	front.AccessLogs = nil
}

func (front *FrontMiddleware) Handler() types.MiddlewareHandler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(_rw http.ResponseWriter, r *http.Request) {
			sTime := time.Now()

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
					front.ctl.Log().Error("-- RECOVER PANIC --", "err", err)
					fmt.Printf("\x1b[31;1m -- RECOVER PANIC -- %s \x1b[0m\n", err)
					fmt.Printf("\x1b[31;1m -- %s --- \x1b[0m\n", debug.Stack())
					http.Error(rw, "Internal Server Error", 500)
				}
			}()

			next.ServeHTTP(rw, r)

			respLength, _ := strconv.Atoi(rw.Header().Get("X-Bytes"))
			front.accessChan <- AccessLog{RemoteAddr: r.RemoteAddr,
				ReqContentType:  r.Header.Get("Content-Type"),
				RespContentType: rw.Header().Get("Content-Type"),
				ReqLength:       r.ContentLength,
				RespLength:      int64(respLength),
				Status:          rw.status,
				Path:            r.URL.RequestURI(),
				Method:          r.Method,
				Cached:           rw.Header().Get("X-Cache"),
				HandlersDuration: time.Since(sTime),
				Timed:            sTime}
		})
	}
}

func (front *FrontMiddleware) logInfo() {
	front.ctl.Log().Info("FrontMiddleware Log:", "clients", front.AccessLogs)
}

func (front *FrontMiddleware) Shutdown() {
	front.logInfo()
	front.stopLog <- true
}
