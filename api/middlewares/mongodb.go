package middlewares

import (
	"time"
	"gopkg.in/mgo.v2"
	"log"
)

var dumpDB *mgo.Session

func InitGlobalDumpDB(dumpDBConfig string){
	mgoConn, err := newMongo(dumpDBConfig)
	if err != nil {
		log.Fatal("Middlewares: initGlobalDumpDB conn", "err", err)
	return
	}

	dumpDB = mgoConn
}

func newMongo(url string) (*mgo.Session, error) {
	sess, err := mgo.DialWithTimeout(url, 3*time.Second)
	if err != nil {
		return nil, err
	}

	sess.SetMode(mgo.Monotonic, true)

	return sess, nil
}
