package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

// JSONLDHandler provides a dummy JSON-LD API
type JSONLDHandler struct {
}

// BooksContent contains the raw JSON of /books.jsonld
const BooksContent = `{
	"@id": "/books.jsonld",
	"hydra:member": [
		"/books/1.jsonld",
		"/books/2.jsonld"
	],
	"foo": [
		{"bar": [{"a": "b"}, {"c": "d"}], "car": "caz"},
		{"bar": [{"a": "d"}, {"c": "e"}], "car": "caz2"}
	]
	}`

// Author1Content contains the raw JSON of /authors/1.jsonld
const Author1Content = `{
	"@id": "/authors/1.jsonld",
	"name": "Kévin"
	}`

func (h *JSONLDHandler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Add("Content-Type", "application/ld+json")
	rw.Header().Add("Access-Control-Allow-Origin", "http://localhost:8081")
	rw.Header().Add("Access-Control-Allow-Credentials", "true")
	rw.Header().Add("Access-Control-Allow-Headers", "Cache-Control")
	rw.Header().Add("Access-Control-Allow-Methods", "GET,HEAD,OPTIONS")
	rw.Header().Add("Cache-Control", "public, max-age=30")

	var myCookieValue string
	if cookie, err := req.Cookie("myCookie"); err == nil {
		myCookieValue = cookie.Value
		rw.Header().Add("Passed-Cookie", myCookieValue)
	}

	if myCookieValue == "" {
		http.SetCookie(rw, &http.Cookie{Name: "myCookie", Value: "foo"})
	}

	m := http.NewServeMux()
	m.HandleFunc("/forwarded", func(rw http.ResponseWriter, req *http.Request) {
		fmt.Fprint(rw, "X-Forwarded-Host: "+req.Header.Get("X-Forwarded-Host")+"\nX-Forwarded-Proto: "+req.Header.Get("X-Forwarded-Proto")+"\nX-Forwarded-For: "+req.Header.Get("X-Forwarded-For"))
	})
	m.HandleFunc("/books.jsonld", func(rw http.ResponseWriter, req *http.Request) {
		fmt.Fprint(rw, BooksContent)
	})
	m.HandleFunc("/authors/", func(rw http.ResponseWriter, req *http.Request) {
		fmt.Fprint(rw, Author1Content)
	})
	m.HandleFunc("/books/", func(rw http.ResponseWriter, req *http.Request) {
		u, _ := url.Parse(req.RequestURI)
		u.RawQuery = ""

		encodedURI, _ := json.Marshal(u.String())
		fmt.Fprint(rw, `{
	"@id": `+string(encodedURI)+`,
	"title": "Book 1",
	"description": "A good book",
	"author": "/authors/1.jsonld",
	"related": "/books/99.jsonld"
	}`)
	})

	m.ServeHTTP(rw, req)
}
