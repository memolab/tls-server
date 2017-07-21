package api

import (
	"encoding/json"
	"net/http"
	"time"

	"go.uber.org/zap"

	"tls-server/utils"

	"golang.org/x/crypto/bcrypt"
	"gopkg.in/mgo.v2/bson"
)

func (c *APICtl) signInHandler(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":

		dbc := c.mongo.Copy()
		defer dbc.Close()

		params := struct {
			Email    string `json:"email" valid:"req,email"`
			Password string `json:"password" valid:"req,alphaNumu,min=5,max=60"`
		}{}

		if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
			c.Abort(rw, http.StatusBadRequest)
			return
		}

		if errs := utils.ValidateStruct(params); len(errs) > 0 {
			c.Abort(rw, http.StatusBadRequest)
			return
		}

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
		return
	}

	c.Abort(rw, http.StatusMethodNotAllowed)
}

func (c *APICtl) signUpHandler(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		dbc := c.mongo.Copy()
		defer dbc.Close()

		user := User{}

		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			c.Abort(rw, http.StatusBadRequest)
			return
		}

		if errs := utils.ValidateStruct(user); len(errs) > 0 {
			c.RespJSON(rw, http.StatusNotAcceptable, errs)
			return
		}

		var (
			err     error
			newPass []byte
			token   string
		)

		if newPass, err = bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost); err != nil {
			c.log.Error("bcrypt gen pass", zap.Error(err))
			c.Abort(rw, http.StatusInternalServerError)
			return
		}

		user.HPassword = string(newPass)
		// new user
		dated := time.Now().UTC()
		user.ID = bson.NewObjectId()
		user.IsActive = true
		user.Dated = &dated
		user.Updated = &dated

		id, _ := user.ID.MarshalText()
		if token, err = c.auth.NewSecretToken(id); err != nil {
			c.log.Error("gen NewSecretToken", zap.Error(err))
			c.Abort(rw, http.StatusInternalServerError)
			return
		}

		if err = dbc.DB("").C("users").Insert(user); err != nil {
			c.log.Error("db users Insert", zap.Error(err))
			c.Abort(rw, http.StatusInternalServerError)
			return
		}

		c.RespJSON(rw, http.StatusCreated, struct {
			User  UserLoged `json:"user"`
			Token string    `json:"token"`
		}{
			User:  UserLoged{ID: user.ID, Username: user.Username, Avatar: user.Avatar},
			Token: token,
		})
		return
	}

	c.Abort(rw, http.StatusMethodNotAllowed)
}
