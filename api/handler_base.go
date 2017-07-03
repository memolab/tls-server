package api

import (
	"bytes"
	"net/http"

	"go.uber.org/zap"

	"tls-server/api/middlewares"
	"tls-server/api/types"
)

func (c *APICtl) indexHanler(rw http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		c.Abort(rw, 404)
		return
	}

	c.RespJson(rw, 200, map[string]interface{}{"msg": "Index Api"})
}

func (c *APICtl) userIndexHanler(rw http.ResponseWriter, r *http.Request) {
	uid := r.Context().Value(types.CTXKey("uid")).(string)

	c.regMidd["cache"].(*middlewares.CacheMiddleware).RespJson(rw, r, 200, map[string]interface{}{"msg": uid})
}

func (c *APICtl) user2IndexHanler(rw http.ResponseWriter, r *http.Request) {
	cache := c.regMidd["cache"].(*middlewares.CacheMiddleware)
	data := cache.Get([]byte("/user;593c4d4d45cf2708b6cb532d"))
	cache.RespFlat(rw, r, 200, data)
}

func (c *APICtl) initDBHanler(rw http.ResponseWriter, r *http.Request) {
	dbc := c.mongo.Copy()
	defer dbc.Close()

	uid := r.Context().Value(types.CTXKey("uid")).(string)

	if !bytes.Equal([]byte(uid), []byte("557840937ab117f73048710c")) {
		c.Abort(rw, http.StatusForbidden)
		return
	}

	if err := createMongoIndexs(dbc); err != nil {
		c.log.Error("create mongo indexs", zap.Error(err))
		c.RespJson(rw, 500, map[string]interface{}{"msg": err})
	}

	c.RespJson(rw, 200, map[string]interface{}{"msg": "ok", "uid": uid})
}
