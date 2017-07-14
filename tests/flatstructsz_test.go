package tests

import (
	"log"
	"testing"
	"time"
	"tls-server/api/middlewares"
	"tls-server/api/structsz/accessLogs"

	mgo "gopkg.in/mgo.v2"
)

func TestMakeListAccessLogs(t *testing.T) {
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

	bts := accessLogs.MakeAccessLogs(list)

	li := accessLogs.GetRootAsAccessLogs(bts, 0)
	for i := 0; i < 5; i++ {
		l := &accessLogs.AccessLog{}
		if li.List(l, i) {
			if string(l.ID()) != list[i].ID.Hex() {
				t.Errorf("expect ID  %s at %d", list[i].ID.Hex(), i)
			}
		} else {
			t.Errorf("error get fom bytes accessLogs.AccessLog at %d", i)
		}

	}
}
