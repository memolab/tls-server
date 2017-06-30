package api

import (
	"encoding/json"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"
	validator "gopkg.in/go-playground/validator.v9"
	"gopkg.in/mgo.v2/bson"
	"tls-server/api/middlewares"
)

func (c *APICtl) signInHandler(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":

		dbc := c.mongo.Copy()
		defer dbc.Close()

		params := struct {
			Email    string `json:"email" valid:"required,email,min=5,max=60"`
			Password string `json:"password" valid:"required,alphanumunicode2,min=5,max=60"`
		}{}

		if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
			c.Abort(rw, http.StatusBadRequest)
			return
		}

		if errs := c.validate.Struct(params); errs != nil {
			c.Abort(rw, http.StatusBadRequest)
			return
		}

		user := User{}
		if err := dbc.DB("").C("users").Find(bson.M{"email": params.Email}).One(&user); err != nil {
			c.RespJson(rw, http.StatusNotAcceptable, map[string]string{"msg": "Bad Credentials."})
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(user.HPassword), []byte(params.Password)); err != nil {
			c.RespJson(rw, http.StatusNotAcceptable, map[string]string{"msg": "Bad Credentials."})
			return
		}

		auth := c.regMidd["auth"].(*middlewares.AuthMiddleware)

		if token, err := auth.NewSecretToken(user.ID.Hex()); err == nil {
			c.RespJson(rw, http.StatusOK, struct {
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

		errs := c.validate.Struct(user)
		if errs != nil {
			rerrs := map[string]string{}
			for _, v := range errs.(validator.ValidationErrors) {
				rerrs[v.Field()] = v.Tag()
			}
			c.RespJson(rw, http.StatusNotAcceptable, map[string]map[string]string{"errs": rerrs})
			return
		}

		var (
			err     error
			newPass []byte
			token   string
		)

		if newPass, err = bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost); err != nil {
			c.log.Error("bcrypt gen pass", "err", err)
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

		auth := c.regMidd["auth"].(*middlewares.AuthMiddleware)
		if token, err = auth.NewSecretToken(user.ID.Hex()); err != nil {
			c.log.Error("gen NewSecretToken", "err", err)
			c.Abort(rw, http.StatusInternalServerError)
			return
		}

		if err = dbc.DB("").C("users").Insert(user); err != nil {
			c.log.Error("db users Insert", "err", err)
			c.Abort(rw, http.StatusInternalServerError)
			return
		}

		c.RespJson(rw, http.StatusCreated, struct {
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
