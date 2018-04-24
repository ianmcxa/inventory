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

}
