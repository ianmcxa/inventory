package main

import (
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func GetRouter(db *sql.DB) *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	//router.HandleFunc("/location", addLocation).Methods("POST")
	router.HandleFunc("/locations", mkGetLocations(db)).Methods("GET")
	router.HandleFunc("/items", mkGetItems(db)).Methods("GET")
	router.HandleFunc("/", serveHtml)
	return router
}

/*func addLocation(w http.ResponseWriter, r *http.Request) {
	var location Location
	// protect against giant blobs being sent
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
}*/

func mkGetLocations(db *sql.DB) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		locations := GetLocations(db)

		json.NewEncoder(w).Encode(locations)
	}
}

func mkGetItems(db *sql.DB) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		items := GetItems(db)

		json.NewEncoder(w).Encode(items)
	}
}

func serveHtml(w http.ResponseWriter, r *http.Request) {
	page, err := ioutil.ReadFile("./index.html")
	if err != nil {
		log.Fatal(err)
	}
	w.Write(page)
}
