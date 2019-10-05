package api

import (
	"fmt"
	"net/http"
	"regexp"
)

var findID = regexp.MustCompile("[0-9]+")

// OpenAPIHandler provides a dummy API documented with OpenAPI
type OpenAPIHandler struct {
}

// OABooksContent contains the raw JSON of /books
const OABooksContent = `{
	"member": [
		1,
		2
	]
}`

// OAAuthor1Content contains the raw JSON of /authors/1
const OAAuthor1Content = `{
	"id": 1,
	"name": "Kévin"
}`

func (h *OpenAPIHandler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Add("Content-Type", "application/json")

	m := http.NewServeMux()
	m.HandleFunc("/books.json", func(rw http.ResponseWriter, req *http.Request) {
		fmt.Fprint(rw, OABooksContent)
	})
	m.HandleFunc("/authors/", func(rw http.ResponseWriter, req *http.Request) {
		fmt.Fprint(rw, OAAuthor1Content)
	})
	m.HandleFunc("/books/", func(rw http.ResponseWriter, req *http.Request) {
		fmt.Fprint(rw, `{
	"id": `+findID.FindString(req.RequestURI)+`,
	"title": "Book 1",
	"description": "A good book",
	"author": 1
}`)
	})

	m.ServeHTTP(rw, req)
}
