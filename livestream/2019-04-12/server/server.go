package main

import (
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("request are coming")
		time.Sleep(time.Second * time.Duration(2+rand.Int()%3))
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("慢慢慢！"))

		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			panic(err)
		}

		err = r.Body.Close()
		if err != nil {
			panic(err)
		}

		w.Write(b)
	})
	http.Handle("/", r)

	log.Fatal(http.ListenAndServe("localhost:8080", r))
}
