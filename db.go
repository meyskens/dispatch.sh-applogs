package main

import (
	"log"
	"os"
	"time"

	mgo "github.com/globalsign/mgo"
)

var mongo *mgo.Database

type logEntry struct {
	InternalName string `bson:"internalName"`
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

	index := mgo.Index{
		Key:        []string{"internalName", "time"},
		Unique:     false,
		DropDups:   false,
		Background: true, // See notes.
	}

	expire := mgo.Index{
		Key:         []string{"time"},
		Unique:      false,
		DropDups:    false,
		Background:  true, // See notes.
		ExpireAfter: 720 * time.Hour,
	}
	mongo.C("app_logs").EnsureIndex(index)
	mongo.C("app_logs").EnsureIndex(expire)

}

func sendToDB(entry logEntry) {
	c := mongo.C("app_logs")
	err := c.Insert(entry)
	if err != nil {
		log.Println(err)
	}
}
