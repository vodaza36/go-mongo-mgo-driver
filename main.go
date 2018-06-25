package main

import (
	"fmt"
	"log"
	"os"
	"runtime/trace"
	"sync"
	"time"

	"gopkg.in/mgo.v2/bson"

	"gopkg.in/mgo.v2"
)

// User entity
type User struct {
	ID    string `bson:"id"`
	Name  string `bson:"name"`
	Email string `bson:"email"`
}

func insertUser(session *mgo.Session, user *User, waitGroup *sync.WaitGroup) error {
	newSession := session.Copy()
	defer newSession.Close()
	defer waitGroup.Done()
	return newSession.DB("test").C("user").Insert(user)
}

func findUserByID(session *mgo.Session, userID string) (*User, error) {
	newSession := session.Copy()
	defer newSession.Close()
	var result User
	err := newSession.DB("test").C("user").Find(bson.M{"id": userID}).One(&result)
	return &result, err
}

func main() {
	// start tracing
	f, err := os.Create("trace.out")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	trace.Start(f)
	defer trace.Stop()

	// create a Mongo DB connection
	session, err := mgo.Dial("localhost:27017")
	if err != nil {
		log.Fatalf("Error creating connection to Mongo DB: %v", err)
	}
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)

	// first drop existing databse
	err = session.DB("test").DropDatabase()

	if err != nil {
		log.Fatalf("Error droping database test: %v", err)
	}

	// create index
	index := mgo.Index{
		Key:        []string{"id"},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}

	coll := session.DB("test").C("user")
	coll.EnsureIndex(index)

	start := time.Now()

	var waitGroup sync.WaitGroup
	const numUsers = 100
	waitGroup.Add(numUsers)

	// insert numUsers new user records
	for index := 0; index < numUsers; index++ {
		newUser := User{ID: fmt.Sprintf("xhocht%d", index), Name: "Test", Email: fmt.Sprintf("test%d@hochbichler.at", index)}
		go insertUser(session, &newUser, &waitGroup)
	}

	waitGroup.Wait()

	log.Printf("%d users created in %s", numUsers, time.Since(start))
	// find a user
	foundUser, err := findUserByID(session, "xhocht1")

	if err != nil {
		log.Fatalf("Error finding user: %v", err)
	}

	log.Printf("Found user: %#v", foundUser)
}
