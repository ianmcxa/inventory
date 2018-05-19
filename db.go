package main

import (
	"bytes"
	"encoding/gob"
	"log"
	"os"

	bolt "github.com/coreos/bbolt"
)

type Location struct {
	Name   string `json:"name"`
	Parent string `json:"parent"`
}

type Item struct {
	Name     string `json:"name"`
	Location string `json:"location"`
}

func GetDB() *bolt.DB {
	db, err := bolt.Open("inventory.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	return db
}

func SetupDB(db *bolt.DB) {
	db.Batch(func(tx *bolt.Tx) error {
		tx.CreateBucketIfNotExists([]byte("locations"))
		tx.CreateBucketIfNotExists([]byte("items"))
		return nil
	})
}

func AddLocation(db *bolt.DB, location Location) {
	gob.Register(Location{})

	buffer := bytes.Buffer{}
	encoder := gob.NewEncoder(&buffer)
	encoder.Encode(location)

	err := db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("locations"))
		return b.Put([]byte(location.Name), buffer.Bytes())
		return nil
	})

	if err != nil {
		log.Fatal(err)
	}
}

// find location by name and update it
func UpdateLocation(db *bolt.DB, name string, location Location) {
	gob.Register(Location{})

	buffer := bytes.Buffer{}
	encoder := gob.NewEncoder(&buffer)
	encoder.Encode(location)

	err := db.Batch(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("locations"))
		b.Delete([]byte(name))
		return b.Put([]byte(location.Name), buffer.Bytes())
	})

	if err != nil {
		log.Fatal(err)
	}
}

func RemoveLocation(db *bolt.DB, name string) {
	err := db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("locations"))
		return b.Delete([]byte(name))
	})

	if err != nil {
		log.Fatal(err)
	}
}

func GetLocations(db *bolt.DB) []Location {
	locations := make([]Location, 0)
	gob.Register(Location{})

	db.View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		b := tx.Bucket([]byte("locations"))
		// we just return nil if there are no locations stored
		if b == nil {
			return nil
		}

		c := b.Cursor()

		for key, value := c.First(); key != nil; key, value = c.Next() {
			//we must decode the byte array into a location object
			decoder := gob.NewDecoder(bytes.NewReader(value))
			var location Location

			// print error on failure, otherwise add server to array
			err := decoder.Decode(&location)
			if err != nil {
				log.Printf("Unable to decode location %v\n", err)
			} else {
				locations = append(locations, location)
			}
		}

		return nil
	})

	return locations
}

func AddItem(db *bolt.DB, item Item) {
	gob.Register(Item{})

	buffer := bytes.Buffer{}
	encoder := gob.NewEncoder(&buffer)
	encoder.Encode(item)

	err := db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("items"))
		return b.Put([]byte(item.Name), buffer.Bytes())
		return nil
	})

	if err != nil {
		log.Fatal(err)
	}
}

// find item by the name and update its values
func UpdateItem(db *bolt.DB, name string, item Item) {
	gob.Register(Item{})

	buffer := bytes.Buffer{}
	encoder := gob.NewEncoder(&buffer)
	encoder.Encode(item)

	err := db.Batch(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("items"))
		b.Delete([]byte(name))
		return b.Put([]byte(item.Name), buffer.Bytes())
	})

	if err != nil {
		log.Fatal(err)
	}
}

func RemoveItem(db *bolt.DB, name string) {
	err := db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("items"))
		return b.Delete([]byte(name))
	})

	if err != nil {
		log.Fatal(err)
	}
}

func GetItems(db *bolt.DB) []Item {
	items := make([]Item, 0)
	gob.Register(Item{})

	db.View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		b := tx.Bucket([]byte("items"))
		// we just return nil if there are no items stored
		if b == nil {
			return nil
		}

		c := b.Cursor()

		for key, value := c.First(); key != nil; key, value = c.Next() {
			//we must decode the byte array into a item object
			decoder := gob.NewDecoder(bytes.NewReader(value))
			var item Item

			// print error on failure, otherwise add server to array
			err := decoder.Decode(&item)
			if err != nil {
				log.Printf("Unable to decode item %v\n", err)
			} else {
				items = append(items, item)
			}
		}

		return nil
	})

	return items
}
