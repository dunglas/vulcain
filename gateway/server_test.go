package gateway

import (
	"context"
	"crypto/tls"
	"io/ioutil"
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

func createTestingUtils() (*httptest.Server, *Gateway, http.Client) {
	upstream := httptest.NewServer(http.HandlerFunc(api.Fixtures))

	upstreamURL, _ := url.Parse(upstream.URL)
	g := NewGateway(&options{
		Addr:      testAddr,
		MaxPushes: -1,
		Upstream:  upstreamURL,
		CertFile:  "../fixtures/tls/server.crt",
		KeyFile:   "../fixtures/tls/server.key",
	})
	go func() {
		g.Serve()
	}()

	// This is a self-signed certificate
	transport := &http2.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := http.Client{Transport: transport, Timeout: time.Duration(100 * time.Millisecond)}

	return upstream, g, client
}

func TestH2NoPush(t *testing.T) {
	upstream, g, client := createTestingUtils()
	defer upstream.Close()

	// loop until the gateway is ready
	var resp *http.Response
	for resp == nil {
		resp, _ = client.Get(gatewayURL + "/books.jsonld?fields=/hydra:member/*&preload=/hydra:member/*/author")
	}

	b, _ := ioutil.ReadAll(resp.Body)

	assert.Equal(t, []string{"</books/1.jsonld?preload=%2Fauthor>; rel=preload; as=fetch", "</books/2.jsonld?preload=%2Fauthor>; rel=preload; as=fetch"}, resp.Header["Link"])
	assert.Equal(t, `{"hydra:member":["/books/1.jsonld?preload=%2Fauthor","/books/2.jsonld?preload=%2Fauthor"]}`, string(b))
	g.server.Shutdown(context.Background())
}

// Unfortunately, Go's HTTP client doesn't support Pushes yet (https://github.com/golang/go/issues/18594)
// In the meantime, we use Symfony HttpClient
func TestH2Push(t *testing.T) {
	upstream, g, _ := createTestingUtils()
	defer upstream.Close()

	for _, test := range []string{"fields-query", "fields-header", "preload-query", "preload-header", "fields-preload-query", "fields-preload-header"} {
		cmd := exec.Command("../test-push/" + test + ".php")
		cmd.Env = os.Environ()
		cmd.Env = append(cmd.Env, "GATEWAY_URL="+gatewayURL)
		stdoutStderr, err := cmd.CombinedOutput()
		if !assert.NoError(t, err) {
			t.Log(string(stdoutStderr))
		}
	}

	g.server.Shutdown(context.Background())
}

func TestH2PushLimit(t *testing.T) {
	upstream, g, _ := createTestingUtils()
	g.options.MaxPushes = 2
	defer upstream.Close()

	cmd := exec.Command("../test-push/push-limit.php")
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, "GATEWAY_URL="+gatewayURL)
	stdoutStderr, err := cmd.CombinedOutput()
	if !assert.NoError(t, err) {
		t.Log(string(stdoutStderr))
	}

	g.server.Shutdown(context.Background())
}
