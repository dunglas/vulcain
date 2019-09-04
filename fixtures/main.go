package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	mux1 := http.NewServeMux()
	mux1.HandleFunc("/", func(rw http.ResponseWriter, req *http.Request) {
		fmt.Fprint(rw, `<h1>HTTP/2 Fixtures</h1>
		<script>
		fetch("https://localhost:3000/books.jsonld", {credentials: "include"})
			.then(console.log)
			.then(
				fetch("https://localhost:3000/books/1.jsonld", {credentials: "include"})
				.then(console.log)
			)
		</script>`)
	})
	s := &http.Server{
		Addr:    ":8081",
		Handler: mux1,
	}
	go func() {
		log.Println("http://localhost:8081 started")
		log.Fatal(s.ListenAndServe())
	}()

	http.HandleFunc("/", func(rw http.ResponseWriter, req *http.Request) {
		rw.Header().Add("Content-Type", "application/ld+json")
		rw.Header().Add("Access-Control-Allow-Origin", "http://localhost:8081")
		rw.Header().Add("Access-Control-Allow-Credentials", "true")

		var myCookieValue string
		if cookie, err := req.Cookie("myCookie"); err == nil {
			myCookieValue = cookie.Value
			rw.Header().Add("Passed-Cookie", myCookieValue)
		}

		if req.RequestURI == "/books.jsonld" {
			if myCookieValue == "" {
				http.SetCookie(rw, &http.Cookie{Name: "myCookie", Value: "foo"})
			}

			fmt.Fprint(rw, `{
	"@id": "/books.jsonld",
	"hydra:member": [
		"/books/1.jsonld",
		"/books/2.jsonld"
	]
}`)

			return
		}

		fmt.Fprint(rw, `{
	"@id": "/books/1.jsonld",
	"title": "Book 1",
	"description": "A good book",
	"author": "/authors/1.jsonld"
}`)
	})

	log.Println("https://localhost:3000 started")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
