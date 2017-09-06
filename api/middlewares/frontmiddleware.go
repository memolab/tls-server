package middlewares

import (
	"fmt"
	"net/http"
	"runtime/debug"
	"strings"
	"sync"
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
	ctl         types.APICTL
	AccessLogs  []*AccessLog
	accessChan  chan *AccessLog
	allowHeadrs string
	stopLog     chan struct{}
	Closed      *sync.WaitGroup
}

type AccessLog struct {
	ID              bson.ObjectId `bson:"_id,omitempty"`
	RemoteAddr      string        `bson:"RemoteAddr"`
	UID             string        `bson:"UID"`
	ReqContentType  string        `bson:"ReqContentType"`
	RespContentType string        `bson:"RespContentType"`
	ReqLength       int           `bson:"ReqLength"`
	RespLength      int           `bson:"RespLength"`
	Status          int           `bson:"Status"`
	Path            string        `bson:"Path"`
	Query           string        `bson:"Query"`
	Method          string        `bson:"Method"`
	Cached          string        `bson:"Cached"`
	Duration        time.Duration `bson:"Duration"` // get time duration in Nanosecond
	Timed           time.Time     `bson:"Timed"`
}
type AccessLogCount struct {
	ID    bson.ObjectId `bson:"_id,omitempty"`
	Path  string        `bson:"Path"`
	Count int           `bson:"Count"`
	Timed time.Time     `bson:"Timed"`
}

func NewFrontMiddleware(ctl types.APICTL, configAccessLogsDump string, allowHeadrs string) *FrontMiddleware {
	dumpDuration, err := time.ParseDuration(configAccessLogsDump)
	if err != nil {
		ctl.Log().Error("FrontMiddleware: error parsing accessLogsDumpConf duration", zap.Error(err))
	}

	allowHeadrsArr := strings.Split(strings.ToLower(allowHeadrs), ",")
	for i, s := range allowHeadrsArr {
		allowHeadrsArr[i] = strings.Title(strings.TrimSpace(s))
	}

	front := &FrontMiddleware{
		ctl:         ctl,
		accessChan:  make(chan *AccessLog, 500),
		allowHeadrs: strings.Join(allowHeadrsArr, ","),
		stopLog:     make(chan struct{}),
	}

	go func() {
		ticker := time.NewTicker(dumpDuration)
		defer func() {
			ticker.Stop()
			front.Closed.Done()
		}()
		for {
			select {
			case a := <-front.accessChan:
				front.AccessLogs = append(front.AccessLogs, a)
			case <-ticker.C:
				front.dumpLogs()
			case <-front.stopLog:
				close(front.stopLog)
				front.dumpLogs()
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
	bulk1 := uc.Bulk()
	bulk1.Unordered()
	puc := dbc.DB("").C("accessLogsCounts")
	bulk2 := puc.Bulk()
	bulk2.Unordered()

	pc := map[string]int{}
	for _, v := range front.AccessLogs {
		bulk1.Insert(v)
		pc[v.Path]++
	}

	now := time.Now().UTC()
	y, m, d := now.Date()

	D := time.Date(y, m, d, now.Hour(), 0, 0, 0, time.UTC)
	for p, c := range pc {
		bulk2.Upsert(bson.M{"Timed": D, "Path": p}, bson.M{"$inc": bson.M{"Count": c}})
	}
	if _, err := bulk1.Run(); err != nil {
		front.ctl.Log().Error("FrontMiddleware: error db dumpLogs", zap.Error(err))
	}
	if _, err := bulk2.Run(); err != nil {
		front.ctl.Log().Error("FrontMiddleware: error db dumpLogs", zap.Error(err))
	}

	front.AccessLogs = nil
}

// Handler impl middleware handler
// apply base requstes requirmint, logs and statistics
func (front *FrontMiddleware) Handler() types.MiddlewareHandler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(_rw http.ResponseWriter, r *http.Request) {
			sTime := time.Now().UTC()
			rw := &AResponseWriter{ResponseWriter: _rw, status: -1}

			rw.Header().Set("Cache-Control", "no-cache, private")
			rw.Header().Set("expires", "-1")
			rw.Header().Set("Vary", "Accept-Encoding,Origin")
			rw.Header().Set("Access-Control-Allow-Origin", "*")
			/*rw.Header().Set("X-Content-Type-Options", "nosniff")
			rw.Header().Set("x-frame-options", "SAMEORIGIN")
			rw.Header().Set("x-xss-protection", "1; mode=block")*/

			defer func() {
				if err := recover(); err != nil {
					front.ctl.Log().Error("-- RECOVER PANIC --", zap.Any("PANIC ERR", err))
					fmt.Printf("\x1b[31;1m -- RECOVER PANIC Stack-- %s \x1b[0m\n", err)
					fmt.Printf("\x1b[31;1m -- %s --- \x1b[0m\n", debug.Stack())
					fmt.Println("\x1b[31;1m -- --- --- \x1b[0m")
					http.Error(rw, "Internal Server Error", 500)
				}
			}()

			if r.Method == "OPTIONS" {
				rw.Header().Set("Access-Control-Allow-Headers", front.allowHeadrs)
				rw.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
				rw.Header().Set("Access-Control-Max-Age", "86400")
				front.ctl.Abort(rw, 200)
				return
			}

			next.ServeHTTP(rw, r)

			uid := ""
			if uidx := r.Context().Value(types.CTXUIDKey{}); uidx != nil {
				uid = uidx.(string)
			}
			front.accessChan <- &AccessLog{
				RemoteAddr:      strings.Split(r.RemoteAddr, ":")[0],
				UID:             uid,
				ReqContentType:  r.Header.Get("Content-Type"),
				RespContentType: rw.Header().Get("Content-Type"),
				ReqLength:       int(r.ContentLength),
				RespLength:      rw.length,
				Status:          rw.status,
				Path:            r.URL.Path,
				Query:           r.URL.RawQuery,
				Method:          r.Method,
				Cached:          rw.Header().Get("X-Cache"),
				Duration:        time.Since(sTime),
				Timed:           sTime}

			front.ctl.Log().Debug(r.URL.RequestURI(), zap.Int("status", rw.status),
				zap.Duration("duration", time.Since(sTime)))
		})
	}
}

// LogInfo log all pinding data
func (front *FrontMiddleware) LogInfo() {
	front.ctl.Log().Info("FrontMiddleware Log:", zap.Any("clients", front.AccessLogs))
}

// Close end any pinding tasks
func (front *FrontMiddleware) Close(wg *sync.WaitGroup) {
	wg.Add(1)
	front.Closed = wg
	front.stopLog <- struct{}{}
}
