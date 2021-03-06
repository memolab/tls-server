package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	"tls-server/api/middlewares"
	"tls-server/api/structsz/accessLogs"
	"tls-server/api/structsz/accesslogcounts"
	"tls-server/utils"

	"golang.org/x/crypto/bcrypt"

	"gopkg.in/mgo.v2/bson"
)

func (c *APICtl) adminIndexHandler(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		/*if !strings.HasPrefix(r.RemoteAddr, "127.0.0.1") {
			c.Abort(rw, http.StatusForbidden)
			return
		}*/
		c.RespJSON(rw, 200, map[string]interface{}{"msg": "AdminAPI Index"})

	case "POST":
		params := struct {
			Email    string `json:"email" valid:"req,email"`
			Password string `json:"password" valid:"req,alphaNumu,min=5,max=60"`
		}{}

		if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
			c.Abort(rw, http.StatusBadRequest)
			return
		}

		if errs := utils.ValidateStruct(params); len(errs) > 0 {
			fmt.Println(errs)
			c.RespJSON(rw, http.StatusNotAcceptable, map[string]interface{}{"errs": errs})
			return
		}

		dbc := c.mongo.Clone()
		defer dbc.Close()

		user := User{}
		if err := dbc.DB("").C("users").Find(bson.M{"email": params.Email}).One(&user); err != nil {
			c.RespJSON(rw, http.StatusNotAcceptable, map[string]string{"msg": "Bad Credentials."})
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(user.HPassword), []byte(params.Password)); err != nil {
			c.RespJSON(rw, http.StatusNotAcceptable, map[string]string{"msg": "Bad Credentials."})
			return
		}

		id, _ := user.ID.MarshalText()
		if token, err := c.auth.NewSecretToken(id); err == nil {
			c.RespJSON(rw, http.StatusOK, struct {
				User  UserLoged `json:"user"`
				Token string    `json:"token"`
			}{
				User:  UserLoged{ID: user.ID, Username: user.Username, Avatar: user.Avatar},
				Token: token,
			})
			return
		}

		c.Abort(rw, http.StatusBadRequest)
	default:
		c.Abort(rw, http.StatusMethodNotAllowed)
	}
}

func (c *APICtl) adminOverviewHandler(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		re := []middlewares.AccessLogCount{}
		sel := bson.M{}
		ords := map[string]string{
			"t":  "-Timed",
			"t2": "Timed",
			"c":  "-Count",
			"c2": "Count",
		}

		if duration, err := time.ParseDuration(r.FormValue("in")); err == nil {
			sel["Timed"] = bson.M{"$gt": time.Now().UTC().Add(duration)}
		} else {
			sel["Timed"] = bson.M{"$gt": time.Now().UTC().Add(-6 * time.Hour)}
		}

		ord := []string{}
		if vord, ok := ords[r.FormValue("ord")]; ok {
			ord = append(ord, vord)
		} else {
			ord = append(ord, ords["t"])
		}

		if err := middlewares.GetOverview(&re, sel, ord); err != nil {
			c.Abort(rw, 404)
			return
		}
		c.RespFlat(rw, 200, accesslogcounts.MakeAccessLogCounts(&re))

	default:
		c.Abort(rw, http.StatusMethodNotAllowed)
	}
}

func (c *APICtl) adminAccesslogsHandler(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		re := []middlewares.AccessLog{}
		sel := bson.M{}
		ords := map[string]string{
			"t":  "-Timed",
			"t2": "Timed",
			"d":  "-Duration",
			"d2": "Duration",
		}

		if duration, err := time.ParseDuration(r.FormValue("in")); err == nil {
			sel["Timed"] = bson.M{"$gt": time.Now().UTC().Add(duration)}
		} else {
			sel["Timed"] = bson.M{"$gt": time.Now().UTC().Add(-6 * time.Hour)}
		}

		/*if st, err := strconv.Atoi(r.FormValue("status")); err == nil {
			sel["Status"] = bson.M{"$eq": st}
		}

		if ca, err := strconv.Atoi(r.FormValue("cached")); err == nil && ca > 0 {
			if ca == 1 {
				sel["cached"] = bson.M{"$eq": "null"}
			} else if ca == 2 { //cached
				sel["cached"] = bson.M{"$ne": "null"}
			}
		}*/

		ord := []string{}
		if vord, ok := ords[r.FormValue("ord")]; ok {
			ord = append(ord, vord)
		} else {
			ord = append(ord, ords["t"])
		}

		if err := middlewares.GetAccesslogs(&re, sel, ord); err != nil {
			c.Abort(rw, 404)
			return
		}
		c.RespFlat(rw, 200, accessLogs.MakeAccessLogs(&re))

	default:
		c.Abort(rw, http.StatusMethodNotAllowed)
	}

}
