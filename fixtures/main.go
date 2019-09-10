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

	http.HandleFunc("/", api.Fixtures)

	log.Println("https://localhost:3000 started")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
