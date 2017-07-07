package tests

import (
	"net/http"
	"testing"
)

func TestUserHandler(t *testing.T) {
	serve := newServing()
	signUserToken("593c4d4d45cf2708b6cb532d")

	for i := 0; i < 4; i++ {
		w := serve("GET", "/user", ``)
		if w.Code != http.StatusOK {
			t.Errorf("Get /user returned %v. Expected %v", w.Code, http.StatusOK)
		}
	}
}

// go test ./tests -run ^% -bench ^BenchmarkUserHandler$ -benchtime 80s -v

func BenchmarkUserHandler(b *testing.B) {
	serve := newServing()
	signUserToken("593c4d4d45cf2708b6cb532d")

	for i := 0; i < b.N; i++ {
		serve("GET", "/user?test", `{"e": "a"}`)
	}

}
