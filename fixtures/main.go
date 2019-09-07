package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

func main() {
	mux1 := http.NewServeMux()
	mux1.HandleFunc("/", func(rw http.ResponseWriter, req *http.Request) {
		fmt.Fprint(rw, `<h1>HTTP/2 Fixtures</h1>
		<script>
		const apiURL = "https://localhost:3000";
		fetch(apiURL + "/books.jsonld?preload=/hydra:member/*", {credentials: "include"})
			.then(resp => {
				document.write("<p><code>/books.jsonld</code> loaded...</p>")
				console.log(resp)
				return resp.json()
			})
			.then(json => {
				json["hydra:member"].forEach(rel => {
					fetch(apiURL + rel, {credentials: "include"})
					.then(data => {
						document.write("<p><code>" + rel + "</code> loaded...</p>")
						console.log(data)
					})	
				})
			})
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

		if strings.HasPrefix(req.RequestURI, "/books.jsonld") {
			if myCookieValue == "" {
				http.SetCookie(rw, &http.Cookie{Name: "myCookie", Value: "foo"})
			}

			fmt.Fprint(rw, `{
	"@id": "/books.jsonld",
	"hydra:member": [
		"/books/1.jsonld",
		"/books/2.jsonld"
	],
	"foo": [
		{"bar": [{"a": "b"}, {"c": "d"}], "car": "caz"},
		{"bar": [{"a": "d"}, {"c": "e"}], "car": "caz2"}
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
