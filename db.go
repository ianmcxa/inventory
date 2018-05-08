package main

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

type Location struct {
	Name   string `json:"name"`
	Parent string `json:"parent"`
}

type Item struct {
	Name     string `json:"name"`
	Location string `json:"location"`
}

func GetDB() *sql.DB {
	db, err := sql.Open("sqlite3", "./inventory.db")
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	return db
}

func SetupDB(db *sql.DB) {
	setupStmt := `CREATE TABLE IF NOT EXISTS locations(
	name STRING PRIMARY KEY,
	parent STRING,
	FOREIGN KEY(parent) REFERENCES locations(name)
	);
	CREATE TABLE IF NOT EXISTS items(
	name STRING PRIMARY KEY, 
	location STRING,
	FOREIGN KEY(location) REFERENCES locations(name)
	)`

	_, err := db.Exec(setupStmt)
	if err != nil {
		log.Fatal("%q: %s\n", err, setupStmt)
		os.Exit(1)
	}
}

func AddLocation(db *sql.DB, location Location) {
	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}

	stmt, err := tx.Prepare("INSERT INTO locations VALUES(?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(location.Name, location.Parent)
	if err != nil {
		log.Fatal(err)
	}

	tx.Commit()
}

func RemoveLocation(db *sql.DB, name string) {
	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}

	stmt, err := tx.Prepare("DELETE FROM locations WHERE name = ?")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(name)
	if err != nil {
		log.Fatal(err)
	}

	tx.Commit()
}

func GetLocations(db *sql.DB) []Location {
	locations := make([]Location, 0)

	rows, err := db.Query("SELECT * FROM locations")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var l Location
		err = rows.Scan(&l.Name, &l.Parent)
		if err != nil {
			log.Fatal(err)
		}
		locations = append(locations, l)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	return locations
}

func AddItem(db *sql.DB, item Item) {
	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}

	stmt, err := tx.Prepare("INSERT INTO items VALUES(?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(item.Name, item.Location)
	if err != nil {
		log.Fatal(err)
	}

	tx.Commit()
}

func MoveItem(db *sql.DB, item Item) {
	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}

	stmt, err := tx.Prepare("UPDATE items SET location = ? WHERE name = ?")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(item.Location, item.Name)
	if err != nil {
		log.Fatal(err)
	}

	tx.Commit()
}

func RemoveItem(db *sql.DB, name string) {
	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}

	stmt, err := tx.Prepare("DELETE FROM items WHERE name = ?")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(name)
	if err != nil {
		log.Fatal(err)
	}

	tx.Commit()
}

func GetItems(db *sql.DB) []Item {
	items := make([]Item, 0)

	rows, err := db.Query("SELECT * FROM items")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var i Item
		err = rows.Scan(&i.Name, &i.Location)
		if err != nil {
			log.Fatal(err)
		}
		items = append(items, i)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	return items
}

func GetItemsByLocation(db *sql.DB, location string) []Item {
	items := make([]Item, 0)

	rows, err := db.Query("SELECT * FROM items WHERE location = ?", location)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var i Item
		err = rows.Scan(&i.Name, &i.Location)
		if err != nil {
			log.Fatal(err)
		}
		items = append(items, i)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	return items
}
