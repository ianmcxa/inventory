package main

import (
	"log"
	"net/http"
)

func main() {
	db := GetDB()
	defer db.Close()
	SetupDB(db)

	router := GetRouter()
	log.Fatal(http.ListenAndServe(":3742", router))
}
