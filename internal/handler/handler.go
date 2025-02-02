package main

import (
	"github.com/go-chi/chi"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {
	log.Print("Starting handler")
	r := chi.NewRouter()
	r.HandleFunc("/", requestHandler)
	err := http.ListenAndServe(":8080", r)
	if err != nil {
		log.Fatal(err)
	}
}

func requestHandler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}
	defer r.Body.Close()
	log.Printf("Received request: %s", string(body))
}
