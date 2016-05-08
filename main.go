package main

import (
	"log"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Task struct {
	Description string
	Due         time.Time
}

type Category struct {
	Id          bson.ObjectId `bson:"_id,omitempty"`
	Name        string
	Description string
	Tasks       []Task
}

func main() {
	dialInfo := &mgo.DialInfo{
		Addrs:    []string{"localhost"},
		Timeout:  60 * time.Second,
		Database: "mymongodb",
	}

	session, err := mgo.DialWithInfo(dialInfo)
	if err != nil {
		log.Printf("Session failed:", err)
	}

	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("mymongodb").C("categories")

	insertRecord(c)
	retrieveAll(c)
	retrieveOneRecord(c)
	// replace string representation of a bson object id here as one of the arguments enclosed in ""
	// in updateDocument(), removeDoc() and removeDocs() functions.
	// updateDocument(c, "")
	// removeDoc(c, "")
	// removeDocs(c, "")
}

func removeDoc(c *mgo.Collection, _id string) {
	err := c.Remove(bson.M{"_id": _id})
	if err != nil {
		log.Printf("Error while removing document:", err)
		return
	}
	log.Println("Successfully deleted document with id:", _id)
}

func removeDocs(c *mgo.Collection, ids []string) {
	changeInfo, err := c.RemoveAll(ids)
	if err != nil {
		log.Printf("Error while removing documents:", err)
		return
	}
	log.Printf("Successfully removed %d documents!", changeInfo.Removed)
}

func insertRecord(c *mgo.Collection) {
	category := Category{
		bson.NewObjectId(),
		"Category50",
		"This is category 50",
		[]Task{
			Task{"Description for task 50", time.Date(2016, time.May, 5, 9, 0, 0, 0, time.UTC)},
			Task{"Description for task 50", time.Date(2016, time.May, 5, 10, 0, 0, 0, time.UTC)},
		},
	}

	err := c.Insert(category)
	if err != nil {
		log.Printf("Error while inserting a record to the database.")
		return
	}
	count, err := c.Count()
	if err != nil {
		log.Printf("Error while reading the number of records in the database.")
		return
	}
	log.Printf("Success! There are now %d record(s) in the database.", count)
}

func retrieveOneRecord(c *mgo.Collection) {
	result := &Category{}
	err := c.Find(bson.M{"name": "Category50"}).One(result)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Category:%s, Description:%s", result.Name, result.Description)
	tasks := result.Tasks
	for _, v := range tasks {
		log.Printf("Task:%s Due:%v\n", v.Description, v.Due)
	}
	log.Printf("ID:%s", result.Id.String())
}

func updateDocument(c *mgo.Collection, _id string) {
	err := c.Update(bson.M{"_id": _id}, bson.M{"$set": bson.M{
		"description": "Create open source projects",
		"tasks": []Task{
			Task{"Evaluate task 1", time.Date(2016, time.May, 6, 9, 0, 0, 0, time.UTC)},
			Task{"Evaluate task 2", time.Date(2016, time.May, 6, 10, 0, 0, 0, time.UTC)},
		},
	},
	})

	if err != nil {
		log.Printf("Error while updating the document:%v", err)
		return
	}
	log.Println("Successfully updated the document!")
}

func retrieveAll(c *mgo.Collection) {
	iter := c.Find(nil).Sort("-name").Iter()
	result := Category{}

	for iter.Next(&result) {
		log.Printf("Category:%s, Description:%s\n", result.Name, result.Description)
		tasks := result.Tasks
		for _, v := range tasks {
			log.Printf("Task:%s Due:%v\n", v.Description, v.Due)
		}
	}
	if err := iter.Close(); err != nil {
		log.Fatal(err)
	}
}
