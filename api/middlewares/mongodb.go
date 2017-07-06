package middlewares

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var (
	dumpDB  *mgo.Session
	devMode bool
)

func GetAccesslogs(re *[]AccessLog, find bson.M, sort []string) error {
	dbc := dumpDB.Copy()
	defer dbc.Close()

	explain(dbc.DB("").C("accessLogs").Find(find).Sort(sort...))

	return dbc.DB("").C("accessLogs").Find(find).Sort(sort...).All(re)
}

func CreateDumpDBIndexs() []error {
	indxs := map[string][]mgo.Index{
		"accessLogs": []mgo.Index{
			mgo.Index{Key: []string{"-Timed", "Status", "Cached", "Duration"},
				Background: true,
				Sparse:     true,
			},
		},
		"rateLimitLogs": []mgo.Index{
			mgo.Index{Key: []string{"-Timed", "TypeID", "OverCount"},
				Background: true,
				Sparse:     true,
			},
		},
	}

	var errs []error

	for k, indxs := range indxs {
		for _, indx := range indxs {
			if err := dumpDB.DB("").C(k).EnsureIndex(indx); err != nil {
				errs = append(errs, err)
			}
		}

	}

	return errs
}

func explain(qry *mgo.Query) {
	if !devMode {
		return
	}

	exp := bson.M{}
	if err := qry.Explain(exp); err == nil {
		expd, _ := json.MarshalIndent(exp, "", " ")
		fmt.Println(string(expd))
	}
}

func InitGlobalDumpDB(dumpDBConfig string, prod string) {
	mgoConn, err := newMongo(dumpDBConfig)
	if err != nil {
		log.Fatal("Middlewares: initGlobalDumpDB conn ", "err ", err)
		return
	}

	if prod == "0" {
		devMode = true
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
