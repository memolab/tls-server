package api

import (
	"net/http"
	"strings"
)

func (c *APICtl) adminIndexHanler(rw http.ResponseWriter, r *http.Request) {
	if !strings.HasPrefix(r.RemoteAddr, "127.0.0.1") {
		c.Abort(rw, http.StatusForbidden)
		return
	}

	c.RespJson(rw, 200, map[string]interface{}{"msg": "Index Api"})
}
