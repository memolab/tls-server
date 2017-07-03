package tests

import (
	"encoding/json"
	"net/http"
	"strconv"
	"testing"
)

func TestToken(t *testing.T) {
	serve := newServing(t)

	w1 := serve("GET", "/", ``)
	if w1.Code != http.StatusOK {
		t.Errorf("Get / returned %v. Expected %v", w1.Code, http.StatusOK)
	}

	w2 := serve("GET", "/user", ``)
	if w2.Code != http.StatusForbidden {
		t.Errorf("Get /user returned %v. Expected %v", w2.Code, http.StatusForbidden)
	}

	signUserToken("557840937ab117f73048710c")

	w3 := serve("GET", "/user", ``)
	if w3.Code != http.StatusOK {
		t.Errorf("Get /user returned %v. Expected %v", w3.Code, http.StatusOK)
	}

	re := &struct {
		Msg string `json:"msg"`
	}{}
	if err := json.Unmarshal(w3.Body.Bytes(), re); err != nil {
		t.Error("Get /user returned pad resp 'json.Unmarshal' ")
	}

	if re.Msg != "557840937ab117f73048710c" {
		t.Errorf("Get /user returned pad user token MSG=%s", re.Msg)
	}
}

func TestDBIndex(t *testing.T) {
	serve := newServing(t)
	signUserToken("557840937ab117f73048710c")

	w1 := serve("POST", "/initdb", ``)
	if w1.Code != http.StatusOK {
		t.Errorf("POST / returned %v. Expected %v", w1.Code, http.StatusOK)
	}
}

func TestCache(t *testing.T) {
	serve := newServing(t)

	signUserToken("557840937ab117f73048710c")

	for i := 0; i < 5; i++ {
		w := serve("GET", "/user?i="+strconv.Itoa(i), ``)
		if w.Code != http.StatusOK {
			t.Errorf("Get /user returned %v. Expected %v", w.Code, http.StatusOK)
		}

		if i > 0 {
			if w.Header().Get("X-Content") != "cached" {
				t.Errorf("Get /user faild cached header %v. Expected %v", w.Header().Get("X-Content"), "cached")
			}
		}
	}

}
