package main

import (
	"log"
	"net/http"

	"github.com/dunglas/vulcain/fixtures/api"
)

func main() {
	mux1 := http.NewServeMux()
	mux1.Handle("/", http.FileServer(http.Dir("static")))
	s := &http.Server{
		Addr:    ":8081",
		Handler: mux1,
	}
	go func() {
		log.Println("http://localhost:8081 started")
		log.Fatal(s.ListenAndServe())
	}()

	http.HandleFunc("/", api.Fixtures)

	log.Println("https://localhost:3000 started")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
