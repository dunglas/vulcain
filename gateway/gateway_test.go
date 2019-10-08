package gateway

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"

	"github.com/dunglas/vulcain/fixtures/api"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
)

func TestNewGateway(t *testing.T) {
	g := NewGateway(&Options{})
	assert.NotNil(t, g)
}

func TestNewGatewayFromEnv(t *testing.T) {
	g, err := NewGatewayFromEnv()
	assert.NotNil(t, g)
	assert.Nil(t, err)

	os.Setenv("KEY_FILE", "foo")
	defer os.Unsetenv("KEY_FILE")
	g, err = NewGatewayFromEnv()
	assert.Nil(t, g)
	assert.Error(t, err)
}

func createServers() (*httptest.Server, *httptest.Server) {
	upstream := httptest.NewServer(&api.JSONLDHandler{})

	upstreamURL, _ := url.Parse(upstream.URL)
	g := NewGateway(&Options{Upstream: upstreamURL})
	gateway := httptest.NewServer(g)

	return upstream, gateway
}

func TestNotModified(t *testing.T) {
	upstream, gateway := createServers()
	defer upstream.Close()
	defer gateway.Close()

	resp, _ := http.Get(gateway.URL + "/books.jsonld")
	b, _ := ioutil.ReadAll(resp.Body)

	assert.Equal(t, api.BooksContent, string(b))
}

func TestFieldsQuery(t *testing.T) {
	upstream, gateway := createServers()
	defer upstream.Close()
	defer gateway.Close()

	resp, _ := http.Get(gateway.URL + "/books.jsonld?fields=/@id")
	b, _ := ioutil.ReadAll(resp.Body)

	assert.Equal(t, `{"@id":"/books.jsonld"}`, string(b))
}

func TestFieldsHeader(t *testing.T) {
	upstream, gateway := createServers()
	defer upstream.Close()
	defer gateway.Close()

	client := &http.Client{}
	req, _ := http.NewRequest("GET", gateway.URL+"/books.jsonld", nil)
	req.Header.Add("Fields", `/@id`)

	resp, _ := client.Do(req)
	b, _ := ioutil.ReadAll(resp.Body)

	assert.Equal(t, "Fields", resp.Header.Get("Vary"))
	assert.Equal(t, `{"@id":"/books.jsonld"}`, string(b))
}

func TestPreloadQuery(t *testing.T) {
	upstream, gateway := createServers()
	defer upstream.Close()
	defer gateway.Close()

	resp, _ := http.Get(gateway.URL + "/books.jsonld?fields=/hydra:member/*&preload=/hydra:member/*/author")
	b, _ := ioutil.ReadAll(resp.Body)

	assert.Equal(t, []string{"</books/1.jsonld?preload=%2Fauthor>; rel=preload; as=fetch", "</books/2.jsonld?preload=%2Fauthor>; rel=preload; as=fetch"}, resp.Header["Link"])
	assert.Equal(t, `{"hydra:member":["/books/1.jsonld?preload=%2Fauthor","/books/2.jsonld?preload=%2Fauthor"]}`, string(b))
}

func TestPreloadHeader(t *testing.T) {
	upstream, gateway := createServers()
	defer upstream.Close()
	defer gateway.Close()

	client := &http.Client{}
	req, _ := http.NewRequest("GET", gateway.URL+"/books.jsonld", nil)
	req.Header.Add("Fields", `/hydra:member`)
	req.Header.Add("Preload", `/hydra:member/*`)

	resp, _ := client.Do(req)
	b, _ := ioutil.ReadAll(resp.Body)

	assert.Equal(t, []string{"</books/1.jsonld>; rel=preload; as=fetch", "</books/2.jsonld>; rel=preload; as=fetch"}, resp.Header["Link"])
	assert.Equal(t, "Fields, Preload", resp.Header.Get("Vary"))
	assert.Equal(t, `{"hydra:member":[
		"/books/1.jsonld",
		"/books/2.jsonld"
	]}`, string(b))
}

func TestUpstreamError(t *testing.T) {
	hook := test.NewGlobal()

	upstreamURL, _ := url.Parse("https://notexist")
	g := NewGateway(&Options{Upstream: upstreamURL})
	gateway := httptest.NewServer(g)
	defer gateway.Close()

	client := &http.Client{}
	req, _ := http.NewRequest("GET", gateway.URL+"/error", nil)

	resp, _ := client.Do(req)

	assert.Equal(t, http.StatusBadGateway, resp.StatusCode)
	assert.Equal(t, "http: proxy error: dial tcp: lookup notexist: no such host", hook.LastEntry().Message)
}
