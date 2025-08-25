package vulcain_test

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"net/url"

	"github.com/dunglas/vulcain"
)

func Example() {
	handler := http.NewServeMux()
	handler.Handle("/books.json", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = fmt.Fprintln(w, `{
	"title": "1984",
	"genre": "dystopia",
	"author": "/authors/orwell.json"
}`)
	}))
	handler.Handle("/authors/orwell.json", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = fmt.Fprintln(w, `{
			"name": "George Orwell",
			"birthDate": "1903-06-25"
		}`)
	}))

	backendServer := httptest.NewServer(handler)
	defer backendServer.Close()

	rpURL, err := url.Parse(backendServer.URL)
	if err != nil {
		log.Fatal(err)
	}

	vulcain := vulcain.New()

	rpHandler := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		r := req.WithContext(vulcain.CreateRequestContext(rw, req))
		var wait bool
		defer func() { vulcain.Finish(r, wait) }()

		rp := httputil.NewSingleHostReverseProxy(rpURL)
		rp.ModifyResponse = func(resp *http.Response) error {
			if !vulcain.IsValidRequest(r) || !vulcain.IsValidResponse(r, resp.StatusCode, resp.Header) {
				return nil
			}

			newBody, err := vulcain.Apply(r, rw, resp.Body, resp.Header)
			if newBody == nil {
				return err
			}

			wait = true
			newBodyBuffer := bytes.NewBuffer(newBody)
			resp.Body = io.NopCloser(newBodyBuffer)

			return nil
		}
		rp.ErrorHandler = func(rw http.ResponseWriter, req *http.Request, err error) {
			wait = false
		}

		rp.ServeHTTP(rw, req)
	})

	frontendProxy := httptest.NewServer(rpHandler)
	defer frontendProxy.Close()

	resp, err := http.Get(frontendProxy.URL + `/books.json?preload="/author"&fields="/title","/author"`)
	if err != nil {
		log.Fatal(err)
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	// Go's HTTP client doesn't support HTTP/2 Server Push yet, so a Link rel=preload is added as fallback
	// Browsers and other clients supporting Server Push will receive a push instead
	fmt.Printf("%v\n\n", resp.Header.Values("Link"))
	fmt.Printf("%s", b)

	// Output:
	// [</authors/orwell.json>; rel=preload; as=fetch]
	//
	// {"author":"/authors/orwell.json","title":"1984"}
}
