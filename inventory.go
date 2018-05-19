package main

import (
	"log"
	"net/http"
)

const port = ":3742"

func main() {
	db := GetDB()
	defer db.Close()
	SetupDB(db)

	router := GetRouter(db)
	log.Printf("Listening on port %s\n", port)
	log.Fatal(http.ListenAndServe(port, router))
}
