package main

import (
	"net/http"

	"time"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	users := r.PathPrefix("/users").Subrouter()
	users.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(time.Duration(1) * time.Second)
		w.WriteHeader(http.StatusOK)
		return
	}).Methods("GET", "POST")

	users.HandleFunc("/{id}", func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(time.Duration(1) * time.Second)
		w.WriteHeader(http.StatusOK)
		return
	}).Methods("GET")

	users.HandleFunc("/{id}/orders", func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(time.Duration(1) * time.Second)
		w.WriteHeader(http.StatusOK)
		return
	}).Methods("GET")

	r.HandleFunc("/orders", func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(time.Duration(1) * time.Second)
		w.WriteHeader(http.StatusOK)
		return
	}).Methods("GET")

	http.Handle("/", r)

	go http.ListenAndServe(":8080", nil)
	http.ListenAndServe(":8081", nil)
}
