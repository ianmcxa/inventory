package main

import (
	"html"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func GetRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/location", addLocation).Methods("POST")
	return router
}

func addLocation(w http.ResponseWriter, r *http.Request) {
	var location Location
	// protect against giant blobs being sent
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
}
