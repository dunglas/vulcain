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
	m.HandleFunc("/oa/books.json", func(rw http.ResponseWriter, req *http.Request) {
		if _, err := fmt.Fprint(rw, OABooksContent); err != nil {
			panic(err)
		}
	})
	m.HandleFunc("/oa/authors/", func(rw http.ResponseWriter, req *http.Request) {
		if _, err := fmt.Fprint(rw, OAAuthor1Content); err != nil {
			panic(err)
		}
	})
	m.HandleFunc("/oa/books/", func(rw http.ResponseWriter, req *http.Request) {
		if _, err := fmt.Fprint(rw, `{
	"id": `+findID.FindString(req.RequestURI)+`,
	"title": "Book 1",
	"description": "A good book",
	"author": 1
}`); err != nil {
			panic(err)
		}
	})

	m.ServeHTTP(rw, req)
}
