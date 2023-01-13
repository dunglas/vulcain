package vulcain

import (
	"context"
	"crypto/tls"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"os/exec"
	"testing"
	"time"

	"github.com/dunglas/vulcain/fixtures/api"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/http2"
)

const testAddr = "127.0.0.1:4343"
const gatewayURL = "https://" + testAddr

func createTestingUtils(openAPIfile string, maxPushes int) (*httptest.Server, *server, http.Client) {
	var handler http.Handler
	if openAPIfile == "" {
		handler = &api.JSONLDHandler{}
	} else {
		handler = &api.OpenAPIHandler{}
	}

	upstream := httptest.NewServer(handler)

	upstreamURL, _ := url.Parse(upstream.URL)
	s := NewServer(&ServerOptions{
		Debug:       true,
		Addr:        testAddr,
		MaxPushes:   maxPushes,
		Upstream:    upstreamURL,
		CertFile:    "./fixtures/tls/server.crt",
		KeyFile:     "./fixtures/tls/server.key",
		OpenAPIFile: openAPIfile,
	})
	go func() {
		s.Serve()
	}()

	// This is a self-signed certificate
	transport := &http2.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := http.Client{Transport: transport, Timeout: time.Duration(100 * time.Millisecond)}

	return upstream, s, client
}

func TestNewServerFromEnv(t *testing.T) {
	s, err := NewServerFromEnv()
	assert.NotNil(t, s)
	assert.Nil(t, err)

	os.Setenv("KEY_FILE", "foo")
	defer os.Unsetenv("KEY_FILE")
	s, err = NewServerFromEnv()
	assert.Nil(t, s)
	assert.Error(t, err)
}

func TestForwardedHeaders(t *testing.T) {
	upstream, s, client := createTestingUtils("", -1)
	defer upstream.Close()

	// loop until the server is ready
	var resp *http.Response
	for resp == nil {
		resp, _ = client.Get(gatewayURL + "/forwarded")
	}

	b, _ := io.ReadAll(resp.Body)

	assert.Equal(t, "X-Forwarded-Host: 127.0.0.1:4343\nX-Forwarded-Proto: https\nX-Forwarded-For: 127.0.0.1", string(b))
	_ = s.server.Shutdown(context.Background())
}

func TestH2NoPush(t *testing.T) {
	upstream, g, client := createTestingUtils("", -1)
	defer upstream.Close()

	// loop until the gateway is ready
	var resp *http.Response
	for resp == nil {
		resp, _ = client.Get(gatewayURL + `/books.jsonld?fields="/hydra:member/*"&preload="/hydra:member/*/author"`)
	}

	b, _ := io.ReadAll(resp.Body)

	assert.Equal(t, []string{"</books/1.jsonld?preload=%22%2Fauthor%22>; rel=preload; as=fetch", "</books/2.jsonld?preload=%22%2Fauthor%22>; rel=preload; as=fetch"}, resp.Header["Link"])
	assert.Equal(t, `{"hydra:member":["/books/1.jsonld?preload=%22%2Fauthor%22","/books/2.jsonld?preload=%22%2Fauthor%22"]}`, string(b))
	_ = g.server.Shutdown(context.Background())
}

func TestMultipleValues(t *testing.T) {
	upstream, g, client := createTestingUtils("", -1)
	defer upstream.Close()

	// loop until the gateway is ready
	var resp *http.Response
	for resp == nil {
		req, _ := http.NewRequest("GET", gatewayURL+"/books/1.jsonld", nil)
		req.Header.Add("Preload", `"/author","/related"`)
		req.Header.Add("Fields", `"/author","/related"`)
		resp, _ = client.Do(req)
	}

	b, _ := io.ReadAll(resp.Body)

	assert.Equal(t, []string{"</authors/1.jsonld>; rel=preload; as=fetch", "</books/99.jsonld>; rel=preload; as=fetch"}, resp.Header["Link"])
	assert.Equal(t, `{"author":"/authors/1.jsonld","related":"/books/99.jsonld"}`, string(b))
	_ = g.server.Shutdown(context.Background())
}

// Unfortunately, Go's HTTP client doesn't support Pushes yet (https://github.com/golang/go/issues/18594)
// In the meantime, we use Symfony HttpClient
func TestH2Push(t *testing.T) {
	upstream, g, _ := createTestingUtils("", -1)
	defer upstream.Close()

	for _, test := range []string{"fields-query", "fields-header", "preload-query", "preload-header", "fields-preload-query", "fields-preload-header"} {
		cmd := exec.Command("./test-push/" + test + ".php")
		cmd.Env = os.Environ()
		cmd.Env = append(cmd.Env, "GATEWAY_URL="+gatewayURL)
		stdoutStderr, err := cmd.CombinedOutput()
		if !assert.NoError(t, err) {
			t.Log(string(stdoutStderr))
		}
	}

	_ = g.server.Shutdown(context.Background())
}

func TestH2PushLimit(t *testing.T) {
	upstream, s, _ := createTestingUtils("", 2)
	defer upstream.Close()

	cmd := exec.Command("./test-push/push-limit.php")
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, "GATEWAY_URL="+gatewayURL)
	stdoutStderr, err := cmd.CombinedOutput()
	if !assert.NoError(t, err) {
		t.Log(string(stdoutStderr))
	}

	_ = s.server.Shutdown(context.Background())
}

func TestH2PushOpenAPI(t *testing.T) {
	upstream, g, _ := createTestingUtils(openapiFixture, -1)
	defer upstream.Close()

	cmd := exec.Command("./test-push/push-openapi.php")
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, "GATEWAY_URL="+gatewayURL)
	stdoutStderr, err := cmd.CombinedOutput()
	if !assert.NoError(t, err) {
		t.Log(string(stdoutStderr))
	}

	_ = g.server.Shutdown(context.Background())
}

func createServers() (*httptest.Server, *httptest.Server) {
	upstream := httptest.NewServer(&api.JSONLDHandler{})

	upstreamURL, _ := url.Parse(upstream.URL)
	s := NewServer(&ServerOptions{Upstream: upstreamURL})
	gateway := httptest.NewServer(s)

	return upstream, gateway
}

func TestNotModified(t *testing.T) {
	upstream, gateway := createServers()
	defer upstream.Close()
	defer gateway.Close()

	resp, _ := http.Get(gateway.URL + "/books.jsonld")
	b, _ := io.ReadAll(resp.Body)

	assert.Equal(t, api.BooksContent, string(b))
}

func TestFieldsQuery(t *testing.T) {
	upstream, gateway := createServers()
	defer upstream.Close()
	defer gateway.Close()

	resp, _ := http.Get(gateway.URL + `/books.jsonld?fields="/@id"`)
	b, _ := io.ReadAll(resp.Body)

	assert.Equal(t, `{"@id":"/books.jsonld"}`, string(b))
}

func TestFieldsHeader(t *testing.T) {
	upstream, gateway := createServers()
	defer upstream.Close()
	defer gateway.Close()

	client := &http.Client{}
	req, _ := http.NewRequest("GET", gateway.URL+"/books.jsonld", nil)
	req.Header.Add("Fields", `"/@id"`)

	resp, _ := client.Do(req)
	b, _ := io.ReadAll(resp.Body)

	assert.Equal(t, "Fields", resp.Header.Get("Vary"))
	assert.Equal(t, `{"@id":"/books.jsonld"}`, string(b))
}

func TestPreloadQuery(t *testing.T) {
	upstream, gateway := createServers()
	defer upstream.Close()
	defer gateway.Close()

	resp, _ := http.Get(gateway.URL + `/books.jsonld?fields="/hydra:member/*"&preload="/hydra:member/*/author"`)
	b, _ := io.ReadAll(resp.Body)

	assert.Equal(t, []string{"</books/1.jsonld?preload=%22%2Fauthor%22>; rel=preload; as=fetch", "</books/2.jsonld?preload=%22%2Fauthor%22>; rel=preload; as=fetch"}, resp.Header["Link"])
	assert.Equal(t, `{"hydra:member":["/books/1.jsonld?preload=%22%2Fauthor%22","/books/2.jsonld?preload=%22%2Fauthor%22"]}`, string(b))
}

func TestPreloadHeader(t *testing.T) {
	upstream, gateway := createServers()
	defer upstream.Close()
	defer gateway.Close()

	client := &http.Client{}
	req, _ := http.NewRequest("GET", gateway.URL+"/books.jsonld", nil)
	req.Header.Add("Fields", `"/hydra:member"`)
	req.Header.Add("Preload", `"/hydra:member/*"`)

	resp, _ := client.Do(req)
	b, _ := io.ReadAll(resp.Body)

	assert.Equal(t, []string{"</books/1.jsonld>; rel=preload; as=fetch", "</books/2.jsonld>; rel=preload; as=fetch"}, resp.Header["Link"])
	assert.Equal(t, []string{"Fields", "Preload"}, resp.Header["Vary"])
	assert.Equal(t, `{"hydra:member":[
		"/books/1.jsonld",
		"/books/2.jsonld"
	]}`, string(b))
}

func TestUpstreamError(t *testing.T) {
	upstreamURL, _ := url.Parse("https://test.invalid")
	g := NewServer(&ServerOptions{Upstream: upstreamURL})
	gateway := httptest.NewServer(g)
	defer gateway.Close()

	client := &http.Client{}
	req, _ := http.NewRequest("GET", gateway.URL+"/error", nil)

	resp, _ := client.Do(req)

	assert.Equal(t, http.StatusBadGateway, resp.StatusCode)
}
