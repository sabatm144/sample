package dbCon

import (
	"log"
	"sync"

	"gopkg.in/mgo.v2"
)

var (
	mongodb *mgo.Database
	mgoOnce sync.Once
)

//ConnectMongoDB will dial mongo db once
func ConnectMongoDB() {
	mgoOnce.Do(func() {
		mongodb = dialMgoDB()
	})
}

func dialMgoDB() (db *mgo.Database) {
	session, conErr := mgo.Dial("localhost:27017")
	if conErr != nil {
		log.Fatalf("Unable to dial mongo server: %v", conErr)
	}

	session.SetMode(mgo.Monotonic, true)
	db = session.DB("sample")

	return db
}

//CopyMongoDB will return a copy of mgo session
func CopyMongoDB() *mgo.Database {
	session := mongodb.Session.Copy()
	return session.DB("sample")
}

//CloseMongoDB will close mongo db connection
func CloseMongoDB() {
	mgoOnce.Do(func() {
		mongodb.Session.Close()
	})
}
