package api

import (
	"net/http"
	"strings"
	"time"
	"tls-server/api/middlewares"

	"gopkg.in/mgo.v2/bson"
)

func (c *APICtl) adminIndexHanler(rw http.ResponseWriter, r *http.Request) {
	if !strings.HasPrefix(r.RemoteAddr, "127.0.0.1") {
		c.Abort(rw, http.StatusForbidden)
		return
	}

	c.RespJSON(rw, 200, map[string]interface{}{"msg": "Index Api"})
}

func (c *APICtl) adminAccesslogsHanler(rw http.ResponseWriter, r *http.Request) {

	re := []middlewares.AccessLog{}

	if err := middlewares.GetAccesslogs(&re, bson.M{"Timed": bson.M{"$gt": time.Now().UTC().Add(-24 * time.Hour)}},
		[]string{"-Timed"}); err != nil {
		c.Abort(rw, 404)
	}

	c.RespJSON(rw, 200, re)
}
