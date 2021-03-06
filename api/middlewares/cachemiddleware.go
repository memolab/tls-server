package middlewares

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"go.uber.org/zap"

	"tls-server/api/structsz/middcachez"
	"tls-server/api/types"
	"tls-server/utils"

	"github.com/boltdb/bolt"
)

// CacheMiddleware provide url resp cache
type CacheMiddleware struct {
	ctl      types.APICTL
	db       *bolt.DB
	chBucket []byte
}

func NewCacheMiddleware(ctl types.APICTL) *CacheMiddleware {
	db, err := bolt.Open("caching.db", 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		ctl.Log().Error("MiddlewareCache: error init bolt db, values: caching.db,0600,timeout:1", zap.Error(err))
	}

	chBucketName := []byte("ch")
	errtbl := db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(chBucketName)
		return err
	})
	if errtbl != nil {
		ctl.Log().Error("MiddlewareCache: error create bolt bucket for handlers caching", zap.ByteString("values", chBucketName), zap.Error(errtbl))
	}

	return &CacheMiddleware{
		ctl:      ctl,
		db:       db,
		chBucket: chBucketName,
	}
}

// Handler NOT implemented, leaved to CacheHandler
func (cache *CacheMiddleware) Handler() types.MiddlewareHandler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			panic("CacheMiddleware handler not implemented")
		})
	}
}

// CacheHandler MiddlewareCache http handler
// urlKey: ["GET,POST", "/"]
// httpKeys: var, head, token
// {"var", "trm", "var":"id", "head": "Header-Key", "tokenUID": "tkn"}
func (cache *CacheMiddleware) CacheHandler(urlKey []string, httpKeys map[string]string, expires time.Duration) types.MiddlewareHandler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			keys := []string{}

			if r.URL.Path == urlKey[1] && strings.Contains(urlKey[0], r.Method) {
				keys = append(keys, r.Method, urlKey[1])
			} else {
				next.ServeHTTP(rw, r)
				return
			}

			for k, v := range httpKeys {
				_v := ""
				switch k {
				case "var":
					_v = utils.EscapeParam(r.FormValue(v))
				case "head":
					_v = utils.EscapeParam(r.Header.Get(v))
				case "tokenUID":
					_v = r.Context().Value(types.CTXUIDKey{}).(string)
				}

				if _v != "" {
					keys = append(keys, v+_v)
				}
			}

			if len(keys) > 0 {
				key := []byte(strings.Join(keys, ""))
				bval := cache.Get([]byte(key))
				if bval != nil {
					obj := middcachez.GetRootAsCacheHandlersObj(bval, 0)
					if time.Now().Sub(time.Unix(obj.Timed(), 0)) <= expires {
						rw.Header().Set("Content-Type", string(obj.ContentType()))
						rw.WriteHeader(int(obj.Status()))

						if rwLen, err := rw.Write(obj.Body()); err != nil {
							cache.ctl.Log().Error("MiddlewareCache: error ResponseWriter", zap.Error(err))
						} else {
							rw.Header().Set("X-Cache", string(key))
							cache.ctl.Log().Debug("MiddlewareCache: serve from cache", zap.String("fullPath", r.URL.RequestURI()),
								zap.ByteString("key", key), zap.Int("length", rwLen))
							return
						}
					}
				}
				*r = *r.WithContext(context.WithValue(r.Context(), types.CTXCACHEKey{}, key))
			}

			next.ServeHTTP(rw, r)
		})
	}
}

func (cache *CacheMiddleware) RespFlat(rw http.ResponseWriter, r *http.Request, status int, data *[]byte) {
	rw.Header().Set("Content-Type", "arraybuffer")

	rw.WriteHeader(status)
	if _, werr := rw.Write(*data); werr != nil {
		cache.ctl.Log().Error("MiddlewareCache:  error text response writer", zap.Error(werr))
		return
	}

	if k := r.Context().Value(types.CTXCACHEKey{}); k != nil && status == 200 {
		if key, ok := k.([]byte); ok && k != nil {
			cache.writeCacheHandler(key, status, []byte("arraybuffer"), data)
		}
	}
}

// RespJSON responce json content type with cachable
func (cache *CacheMiddleware) RespJSON(rw http.ResponseWriter, r *http.Request, status int, data interface{}) {
	rw.Header().Set("Content-Type", "application/json")

	dataByts, jerr := json.Marshal(data)
	if jerr != nil {
		cache.ctl.Log().Error("MiddlewareCache: error marshal json response", zap.Error(jerr))
		http.Error(rw, "Internal Server Error", 500)
		return
	}

	rw.WriteHeader(status)
	if _, werr := rw.Write(dataByts); werr != nil {
		cache.ctl.Log().Error("MiddlewareCache:  error json response writer", zap.Error(werr))
		return
	} else if k := r.Context().Value(types.CTXCACHEKey{}); k != nil && status == 200 {
		if key, ok := k.([]byte); ok && k != nil {
			cache.writeCacheHandler(key, status, []byte("application/json"), &dataByts)
		}
	}
}

// RespJSONRaw responce json content type with cachable
func (cache *CacheMiddleware) RespJSONRaw(rw http.ResponseWriter, r *http.Request, status int, data *[]byte) {
	rw.Header().Set("Content-Type", "application/json")

	rw.WriteHeader(status)
	if _, werr := rw.Write(*data); werr != nil {
		cache.ctl.Log().Error("MiddlewareCache:  error json response writer", zap.Error(werr))
		return
	} else if k := r.Context().Value(types.CTXCACHEKey{}); k != nil && status == 200 {
		if key, ok := k.([]byte); ok && k != nil {
			cache.writeCacheHandler(key, status, []byte("application/json"), data)
		}
	}
}

// Get return bytes from cache db
func (cache *CacheMiddleware) Get(key []byte) (data []byte) {
	cache.db.View(func(tx *bolt.Tx) error {
		d := tx.Bucket(cache.chBucket).Get(key)
		if d != nil {
			data = make([]byte, len(d))
			copy(data, d)
		}
		return nil
	})
	return
}

func (cache *CacheMiddleware) writeCacheHandler(key []byte, status int, ContentType []byte, data *[]byte) {
	dataDB := middcachez.MakeCacheHandlersObj(status, ContentType, data)
	err := cache.db.Update(func(tx *bolt.Tx) error {
		return tx.Bucket(cache.chBucket).Put(key, *dataDB)
	})

	if err != nil {
		cache.ctl.Log().Error("MiddlewareCache: error writeCacheHandler", zap.Error(err))
	}
}

// LogInfo log all caching data
func (cache *CacheMiddleware) LogInfo() {
	cache.ctl.Log().Debug("CacheMiddleware log DB:")

	cache.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(cache.chBucket))
		c := b.Cursor()
		i := 1
		for k, _ := c.First(); k != nil; k, _ = c.Next() {
			//fmt.Println("row: ", "K", string(k), "V", string(v))
			cache.ctl.Log().Debug("row", zap.ByteString(strconv.Itoa(i), k))
			i = i + 1
		}
		return nil
	})
}

// Drop to drop cache bucket
func (cache *CacheMiddleware) Drop() (err error) {
	err = cache.db.Update(func(tx *bolt.Tx) error {
		err := tx.DeleteBucket(cache.chBucket)
		if err == nil {
			_, err = tx.CreateBucketIfNotExists(cache.chBucket)
			return err
		}
		return err
	})
	return
}

// Close end any pinding tasks
func (cache *CacheMiddleware) Close(wg *sync.WaitGroup) {
	cache.LogInfo()
	cache.db.Close()
}
