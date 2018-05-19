package main

import (
	"database/sql"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func GetRouter(db *sql.DB) *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/locations", mkAddLocation(db)).Methods("POST")
	router.HandleFunc("/locations", mkGetLocations(db)).Methods("GET")
	router.HandleFunc("/locations/{name}", mkUpdateLocation(db)).Methods("PUT")
	router.HandleFunc("/locations/{name}", mkDeleteLocation(db)).Methods("DELETE")
	router.HandleFunc("/items", mkGetItems(db)).Methods("GET")
	router.HandleFunc("/items", mkAddItem(db)).Methods("POST")
	router.HandleFunc("/items/{name}", mkUpdateItem(db)).Methods("PUT")
	router.HandleFunc("/items/{name}", mkDeleteItem(db)).Methods("DELETE")
	router.HandleFunc("/", serveHtml)
	return router
}

func mkAddLocation(db *sql.DB) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var location Location
		// protect against giant blobs being sent
		body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
		if err != nil {
			log.Fatal(err)
		}

		if err := r.Body.Close(); err != nil {
			log.Fatal(err)
		}

		if err := json.Unmarshal(body, &location); err != nil {
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.WriteHeader(422) // unprocessable entity
			if err := json.NewEncoder(w).Encode(err); err != nil {
				log.Fatal(err)
			}
		}

		AddLocation(db, location)
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusCreated)
	}
}

func mkUpdateLocation(db *sql.DB) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		name := vars["name"]
		var location Location
		// protect against giant blobs being sent
		body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
		if err != nil {
			log.Fatal(err)
		}

		if err := r.Body.Close(); err != nil {
			log.Fatal(err)
		}

		if err := json.Unmarshal(body, &location); err != nil {
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.WriteHeader(422) // unprocessable entity
			if err := json.NewEncoder(w).Encode(err); err != nil {
				log.Fatal(err)
			}
		}

		UpdateLocation(db, name, location)
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusAccepted)
	}
}

func mkDeleteLocation(db *sql.DB) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		name := vars["name"]
		RemoveLocation(db, name)
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
	}
}

func mkGetLocations(db *sql.DB) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		locations := GetLocations(db)

		json.NewEncoder(w).Encode(locations)
	}
}

func mkUpdateItem(db *sql.DB) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var item Item
		vars := mux.Vars(r)
		name := vars["name"]
		// protect against giant blobs being sent
		body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
		if err != nil {
			log.Fatal(err)
		}

		if err := r.Body.Close(); err != nil {
			log.Fatal(err)
		}

		if err := json.Unmarshal(body, &item); err != nil {
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.WriteHeader(422) // unprocessable entity
			if err := json.NewEncoder(w).Encode(err); err != nil {
				log.Fatal(err)
			}
			return
		}

		UpdateItem(db, name, item)
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusAccepted)
	}
}

func mkAddItem(db *sql.DB) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var item Item
		// protect against giant blobs being sent
		body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
		if err != nil {
			log.Fatal(err)
		}

		if err := r.Body.Close(); err != nil {
			log.Fatal(err)
		}

		if err := json.Unmarshal(body, &item); err != nil {
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.WriteHeader(422) // unprocessable entity
			if err := json.NewEncoder(w).Encode(err); err != nil {
				log.Fatal(err)
			}
		}

		AddItem(db, item)
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusCreated)
	}
}

func mkDeleteItem(db *sql.DB) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		name := vars["name"]
		RemoveItem(db, name)
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
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
