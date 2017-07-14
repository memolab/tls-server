package tests

import (
	"fmt"
	"log"
	"net/http"
	"testing"
	"time"
	"tls-server/api/middlewares"
	"tls-server/api/structsz/accessLogs"

	"encoding/json"

	mgo "gopkg.in/mgo.v2"
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

func BenchmarkMakeListAccessLogs(b *testing.B) {
	mgoConn, err := mgo.DialWithTimeout("mongodb://localhost/app-go-db-logsDump", 3*time.Second)
	if err != nil {
		log.Fatal("Middlewares: DB conn ", err)
		return
	}
	mgoConn.SetMode(mgo.Monotonic, true)

	list := []middlewares.AccessLog{}
	if err := mgoConn.DB("").C("accessLogs").Find(nil).Sort("-Timed").Limit(100).All(&list); err != nil {
		log.Fatal("Middlewares: conn list", err)
		return
	}

	b.Run("GetListBytes", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			getListBytes(list)
		}
	})

	b.Run("getListJson", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			getListJson(list)
		}
	})
}

func getListBytes(list []middlewares.AccessLog) {
	bts := accessLogs.MakeAccessLogs(list)
	li := accessLogs.GetRootAsAccessLogs(bts, 0)
	for j := 0; j < 5; j++ {
		l := &accessLogs.AccessLog{}
		if li.List(l, j) {
			if string(l.ID()) == "" {
				fmt.Errorf("empty ID in  GetListBytes at %d", j)
			}
		} else {
			fmt.Errorf("error get fom bytes accessLogs.AccessLog at %d", j)
		}
	}
}
func getListJson(list []middlewares.AccessLog) {
	bts, err := json.Marshal(list)
	if err != nil {
		fmt.Errorf("Error json.Marshal", err)
	}

	li := []middlewares.AccessLog{}
	if err := json.Unmarshal(bts, &li); err != nil {
		fmt.Errorf("error get fom bytes json")
	}

	for j := 0; j < 5; j++ {
		if li[j].UID == "" {
			fmt.Errorf("empty ID in middlewares.AccessLog at %d", j)
		}
	}
}