package middlewares

import (
	"fmt"
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
	accessLog  map[string]int
	accessChan chan string
}

func NewFrontMiddleware(ctl types.APICTL) *FrontMiddleware {
	front := &FrontMiddleware{
		ctl:        ctl,
		accessLog:  map[string]int{},
		accessChan: make(chan string, 200),
	}

	go func() {
		for {
			p := <-front.accessChan
			front.accessLog[p]++
		}
	}()

	return front
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
			front.accessChan <- (r.URL.Path + ":" + strconv.Itoa(rw.status))

			colr := ""
			switch rw.status {
			case 500:
				colr = "\x1b[31;1m"
			case 400: // StatusBadRequest
			case 405: // StatusMethodNotAllowed
			case 406: // StatusNotAcceptable
			case 429: // StatusTooManyRequests
				colr = "\x1b[35;1m"
			case 404: // StatusNotFound
				colr = "\x1b[33;1m"
			case 401: // StatusUnauthorized
			case 403: // StatusForbidden
				colr = "\x1b[36;1m"
			case 200: // ok
				colr = "\x1b[34;1m"
			default:
				colr = "\x1b[34;1m"
			}

			fmt.Printf("%s%s %s %s %d [%v]\x1b[0m\n", colr, r.RemoteAddr, r.Method, r.URL.String(), rw.status, time.Since(sTime))
		})
	}
}

func (front *FrontMiddleware) logInfo() {
	t := 0
	for _, i := range front.accessLog {
		t += i
	}
	front.ctl.Log().Info("accessLog:", "total", t, "urls", front.accessLog)

}

func (front *FrontMiddleware) Shutdown() {
	front.logInfo()
}