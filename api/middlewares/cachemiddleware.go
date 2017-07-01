package middlewares

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"
	"tls-server/api/structsz/middcachez"
	"tls-server/api/types"

	"github.com/boltdb/bolt"
	"strconv"
)

// MiddlewareCache provide url resp cache
type CacheMiddleware struct {
	ctl      types.APICTL
	db       *bolt.DB
	chBucket []byte
}

func NewCacheMiddleware(ctl types.APICTL) *CacheMiddleware {
	db, err := bolt.Open("caching.db", 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		ctl.Log().Error("MiddlewareCache: error init bolt db", "values", "caching.db,0600,timeout:1", "err", err)
	}

	return &CacheMiddleware{
		ctl:      ctl,
		db:       db,
		chBucket: []byte("ch"),
	}
}

// Handler NOT implemented, leaved to CacheHandler
func (cache *CacheMiddleware) Handler() types.MiddlewareHandler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		})
	}
}

// CacheHandler MiddlewareCache http handler
func (cache *CacheMiddleware) CacheHandler(urlKey string, httpKeys map[string]string, expires time.Duration) types.MiddlewareHandler {
	errtbl := cache.db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(cache.chBucket)
		return err
	})
	if errtbl != nil {
		cache.ctl.Log().Error("MiddlewareCache: error create bolt bucket for handlers caching", "values", string(cache.chBucket), "err", errtbl)
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			keys := []string{}
			if strings.HasPrefix(r.URL.Path, urlKey) {
				keys = append(keys, r.URL.Path)
			}

			for k, v := range httpKeys {
				_v := ""
				switch k {
				case "var":
					_v = r.FormValue(v)
				case "head":
					_v = r.Header.Get(v)
				case "tokenUID":
					_v = r.Context().Value(types.CTXKey("uid")).(string)
				}

				if _v != "" {
					keys = append(keys, _v)
				}
			}

			if len(keys) > 0 {
				key := []byte(strings.Join(keys, ";"))
				bval := cache.Get([]byte(key))

				if bval != nil {
					obj := middcachez.GetRootAsCacheHandlersObj(bval, 0)
					if time.Now().Sub(time.Unix(obj.Timed(), 0)) <= expires {
						rw.Header().Set("Content-Type", string(obj.ContentType()))
						rw.WriteHeader(int(obj.Status()))

						if wlen, err := rw.Write(obj.Body()); err != nil {
							cache.ctl.Log().Error("MiddlewareCache: error ResponseWriter", "err", err)
						} else {
							rw.Header().Set("X-Cache", string(key))
							rw.Header().Set("X-Bytes", strconv.Itoa(wlen))
							//cache.ctl.Log().Debug("MiddlewareCache: serve from cache", "fullPath", r.URL.RequestURI(), "key", string(key), "length", wlen)
							return
						}
					}
				}
				r = r.WithContext(context.WithValue(r.Context(), types.CTXKey("cachein"), key))
			}

			next.ServeHTTP(rw, r)
		})
	}
}

func (cache *CacheMiddleware) RespFlat(rw http.ResponseWriter, r *http.Request, status int, data []byte) {
	rw.Header().Set("Content-Type", "text/plain")

	rw.WriteHeader(status)
	if wlen, werr := rw.Write(data); werr != nil {
		cache.ctl.Log().Error("MiddlewareCache:  error text response writer", "rwWrite", werr)
		return
	}else{
		rw.Header().Set("X-Bytes", strconv.Itoa(wlen))
	}

	if k := r.Context().Value(types.CTXKey("cachein")); k != nil && status == 200 {
		if key, ok := k.([]byte); ok && k != nil {
			cache.writeCacheHandler(key, status, []byte("text/plain"), data)
		}
	}
}

func (cache *CacheMiddleware) RespJson(rw http.ResponseWriter, r *http.Request, status int, data interface{}) {
	rw.Header().Set("Content-Type", "application/json")

	b, jerr := json.Marshal(data)
	if jerr != nil {
		cache.ctl.Log().Error("MiddlewareCache: error json response", "errMarshal", jerr)
		http.Error(rw, "Internal Server Error", 500)
		return
	}

	rw.WriteHeader(status)
	if wlen, werr := rw.Write(b); werr != nil {
		cache.ctl.Log().Error("MiddlewareCache:  error json response writer", "rwWrite", werr)
		return
	}else{
		rw.Header().Set("X-Bytes", strconv.Itoa(wlen))
	}

	if k := r.Context().Value(types.CTXKey("cachein")); k != nil && status == 200 {
		if key, ok := k.([]byte); ok && k != nil {
			cache.writeCacheHandler(key, status, []byte("application/json"), b)
		}
	}
}

func (cache *CacheMiddleware) Get(key []byte) (data []byte) {
	cache.db.View(func(tx *bolt.Tx) error {
		data = tx.Bucket(cache.chBucket).Get(key)
		return nil
	})
	return
}

func (cache *CacheMiddleware) writeCacheHandler(key []byte, status int, ContentType []byte, data []byte) {
	dataDB := middcachez.MakeCacheHandlersObj(status, ContentType, data)
	err := cache.db.Update(func(tx *bolt.Tx) error {
		return tx.Bucket(cache.chBucket).Put(key, dataDB)
	})

	if err != nil {
		cache.ctl.Log().Error("MiddlewareCache: error writeCacheHandler", "err", err)
	}
}

func (cache *CacheMiddleware) logInfo() {
	cache.ctl.Log().Info("CacheMiddleware log DB:")

	cache.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(cache.chBucket))
		c := b.Cursor()
		for k, _ := c.First(); k != nil; k, _ = c.Next() {
			//fmt.Println("row: ", "K", string(k), "V", string(v))
			cache.ctl.Log().Info("row: ", "Key", string(k))
		}
		return nil
	})
}

func (cache *CacheMiddleware) Shutdown() {
	cache.logInfo()
	cache.db.Close()
}
