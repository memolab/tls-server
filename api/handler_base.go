package api

import (
	"bytes"
	"net/http"

	"go.uber.org/zap"

	"tls-server/api/types"
)

func (c *APICtl) indexHanler(rw http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		c.Abort(rw, 404)
		return
	}

	switch r.Method {
	case "GET":
		c.RespJSONRaw(rw, 200, []byte(`{"msg": "API Index"}`))

	default:
		c.Abort(rw, http.StatusMethodNotAllowed)
	}
}

func (c *APICtl) userIndexHanler(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		uid := r.Context().Value(types.CTXUIDKey{}).(string)
		msgb := []byte(`{"msg": "`)
		msgb = append(msgb, uid...)
		msgb = append(msgb, `"}`...)
		c.cache.RespJSONRaw(rw, r, 200, msgb)

	default:
		c.Abort(rw, http.StatusMethodNotAllowed)
	}
}

func (c *APICtl) user2IndexHanler(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		data := c.cache.Get([]byte("/user;593c4d4d45cf2708b6cb532d"))
		c.RespFlat(rw, 200, data)

	default:
		c.Abort(rw, http.StatusMethodNotAllowed)
	}
}

func (c *APICtl) initDBHanler(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		dbc := c.mongo.Copy()
		defer dbc.Close()

		uid := r.Context().Value(types.CTXUIDKey{}).(string)

		if !bytes.Equal([]byte(uid), []byte("557840937ab117f73048710c")) {
			c.Abort(rw, http.StatusForbidden)
			return
		}

		if err := createMongoIndexs(dbc); err != nil {
			c.log.Error("create mongo indexs", zap.Error(err))
			c.RespJSON(rw, 500, map[string]interface{}{"msg": err})
		}

		c.RespJSON(rw, 200, map[string]interface{}{"msg": "ok", "uid": uid})

	default:
		c.Abort(rw, http.StatusMethodNotAllowed)
	}
}
