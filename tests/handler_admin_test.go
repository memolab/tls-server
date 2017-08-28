package tests

import (
	"net/http"
	"testing"
	"tls-server/api/structsz/accesslogcounts"
)

// go test ./tests -run ^TestAdminOverview$ -v
func TestAdminOverview(t *testing.T) {
	serve := newServing()
	signUserToken("593c4d4d45cf2708b6cb532d")

	w1 := serve("GET", "/admin/overview", ``)
	if w1.Code != http.StatusOK {
		t.Errorf("Get /admin/overview returned %v. Expected %v", w1.Code, http.StatusOK)
	}

	li := accesslogcounts.GetRootAsAccessLogCounts(w1.Body.Bytes(), 0)
	alc := accesslogcounts.AccessLogCount{}

	if li.ListLength() > 0 {
		if !(li.List(&alc, 0)) {
			t.Errorf("fail read accesslogcounts.AccessLogCount listing")

			if !(alc.Count() > 0) {
				t.Errorf("fail in accesslogcounts.AccessLogCount listing")
			}
		}
	}
}
