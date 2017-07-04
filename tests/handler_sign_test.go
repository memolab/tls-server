package tests

import (
	"fmt"
	"net/http"
	"testing"
)

func TestSignup(t *testing.T) {
	serve := newServing()

	/*w1 := serve("POST", "/signup", `{"username": "memos", "email": "memo@local.com", "password": "memos"}`)
	if w1.Code != http.StatusCreated {
		t.Errorf("POST /signup [Valid user data] returned %v. Expected %v", w1.Code, http.StatusCreated)
	}*/

	w2 := serve("POST", "/signup", `{"username": "memo_-.", "email": ".com", "password": "memo"}`)
	//fmt.Println(">>>", w2.Body.String())
	if w2.Code != http.StatusNotAcceptable {
		t.Errorf("POST /signup [Invalid user data] returned %v. Expected %v", w2.Code, http.StatusNotAcceptable)
	}

}

func TestSignin(t *testing.T) {
	serve := newServing()

	w1 := serve("POST", "/signin", `{"email": "memo@local.com", "password": "memos"}`)
	fmt.Println("----------------\n", w1.Body.String(), "\n----------------")
	if w1.Code != http.StatusOK {
		t.Errorf("POST /signin [Valid login user data] returned %v. Expected %v", w1.Code, http.StatusOK)
	}

	w2 := serve("POST", "/signin", `{"email": ".com", "password": "memo"}`)
	//fmt.Println(">>>", w2.Body.String())
	if w2.Code != http.StatusBadRequest {
		t.Errorf("POST /signup [Invalid login user data] returned %v. Expected %v", w2.Code, http.StatusBadRequest)
	}

	w3 := serve("POST", "/signin", `{"email": "memo@local.com", "password": "memos_"}`)
	//fmt.Println(">>>", w3.Body.String())
	if w3.Code != http.StatusNotAcceptable {
		t.Errorf("POST /signup [Invalid login user data] returned %v. Expected %v", w3.Code, http.StatusNotAcceptable)
	}

}
