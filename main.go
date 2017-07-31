package main

import (
	"log"
	"net/http"
	"time"

	"github.com/boltdb/bolt"
	"github.com/gorilla/mux"
)

func main() {
	// Open Bolt DB on file system
	db, err := bolt.Open("db/messages.db", 0600, &bolt.Options{Timeout: 3 * time.Second})
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	store, err := InitializeStore(db)
	if err != nil {
		log.Fatal(err)
	}

	api := APIService{db: store}

	r := mux.NewRouter()
	a := r.PathPrefix("/api").Subrouter()
	a.HandleFunc("/messages", api.ListMessages).Methods("GET")
	a.HandleFunc("/messages/{id}", api.GetMessage).Methods("GET")
	a.HandleFunc("/messages", api.CreateMessage).Methods("POST")
	a.HandleFunc("/messages/{id}", api.DeleteMessage).Methods("DELETE")

	r.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir("./client/dist/"))))

	srv := &http.Server{
		Handler:      r,
		Addr:         "0.0.0.0:8080",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Println("Audition is running...")
	log.Fatal(srv.ListenAndServe())
}
