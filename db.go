package main

import (
	"log"
	"os"
	"time"

	mgo "github.com/globalsign/mgo"
)

var mongo *mgo.Database

type logEntry struct {
	InternalName string
	Pod          string
	Container    string
	Line         string
	Time         time.Time
}

func init() {
	session, err := mgo.Dial(os.Getenv("MONGODB_HOST"))
	if err != nil {
		panic(err)
	}
	mongo = session.DB(os.Getenv("MONGODB_DB"))
}

func sendToDB(entry logEntry) {
	c := mongo.C("app_logs")
	err := c.Insert(entry)
	if err != nil {
		log.Println(err)
	}
}
