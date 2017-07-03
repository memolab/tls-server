package middlewares

import (
	"log"
	"time"

	"gopkg.in/mgo.v2"
)

var dumpDB *mgo.Session

func CreateMongoIndexs(mongo *mgo.Session) []error {
	indxs := []mgo.Index{
		mgo.Index{
			Key:        []string{"-Timed", "Status", "HandlersDuration"},
			Background: true,
			Sparse:     true,
		},
		mgo.Index{
			Key:        []string{"Cached", "Status"},
			Background: true,
			Sparse:     true,
		},
	}

	var errs []error

	for _, indx := range indxs {
		if err := dumpDB.DB("").C("accessLogs").EnsureIndex(indx); err != nil {
			errs = append(errs, err)
		}

	}

	return errs
}

func InitGlobalDumpDB(dumpDBConfig string) {
	mgoConn, err := newMongo(dumpDBConfig)
	if err != nil {
		log.Fatal("Middlewares: initGlobalDumpDB conn ", "err ", err)
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
