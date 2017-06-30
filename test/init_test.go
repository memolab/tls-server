package test

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/gorilla/securecookie"

	"tls-server/api"
)

var (
	setit          = false
	usrToken       string
	mux            *http.ServeMux
	scCookie       *securecookie.SecureCookie
	addr           string
	headerTokenKey string
)

type Serving func(method string, url string, params string) *httptest.ResponseRecorder

func setup() {
	if setit {
		return
	}

	log.SetFlags(log.LUTC | log.Lmicroseconds | log.Llongfile)

	var config map[string]string
	file, _ := os.Open("../config.dev.json")
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		fmt.Println("CONFIG ERROR: ../config.dev.json")
		panic(err)
	}

	scCookie = securecookie.New([]byte(config["secretKey1"]), []byte(config["secretKey2"]))
	scCookie.MaxAge(0)

	mux = api.InitAPI(config)
	addr = config["addr"]
	headerTokenKey = config["headerTokenKey"]

	setit = true
	fmt.Println("init Routs ...")
}

func newServing(t *testing.T) Serving {
	setup()
	return func(method string, url string, params string) *httptest.ResponseRecorder {
		req, err := http.NewRequest(method, ("https://" + addr + url), strings.NewReader(params))
		if err != nil {
			log.Fatal(err)
			t.Errorf("%v", err)
		}
		req.Header.Set("Content-Type", "application/json")
		if usrToken != "" {
			req.Header.Set(headerTokenKey, usrToken)
		}

		req.Body.Close()

		rw := httptest.NewRecorder()
		mux.ServeHTTP(rw, req)

		return rw
	}
}

func signUserToken(uid string) {
	var err error
	var token string

	if token, err = scCookie.Encode("i", uid); err != nil {
		log.Fatal("Error generate token: ", token, err)
	}

	usrToken = token
}
