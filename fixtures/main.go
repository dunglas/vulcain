package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/dunglas/vulcain/fixtures/api"
)

func main() {
	mux1 := http.NewServeMux()
	mux1.HandleFunc("/", func(rw http.ResponseWriter, req *http.Request) {
		fmt.Fprint(rw, `<h1>HTTP/2 Fixtures</h1>
		<script>
		const apiURL = "https://localhost:3000";
		fetch(apiURL + "/books.jsonld?preload=/hydra:member/*/author", {credentials: "include", headers: {"Cache-Control": "no-cache, no-store"}})
			.then(booksResp => {
				document.write("<p>Books: <code>/books.jsonld</code> loaded...</p>")
				console.log(booksResp)
				return booksResp.json()
			})
			.then(booksJSON => {
				booksJSON["hydra:member"].forEach(bookPath => {
					fetch(apiURL + bookPath, {credentials: "include"})
					.then(bookResp => {
						document.write("<p>Book: <code>" + bookPath + "</code> loaded...</p>")
						console.log(bookResp)
						return bookResp.json()
					}).then(bookJSON => {
						fetch(apiURL + bookJSON.author, {credentials: "include"})
						.then(authorResp => {
							document.write("<p>Author: <code>" + bookJSON.author + "</code> loaded...</p>")
							console.log(authorResp)
						})
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

	http.HandleFunc("/", api.Fixtures)

	log.Println("https://localhost:3000 started")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
