package caddy

import (
	"reflect"
	"testing"

	"github.com/caddyserver/caddy/v2/caddytest"
)

func TestVulcain(t *testing.T) {
	t.Skip("caddytest is currently broken because of https://github.com/caddyserver/caddy/issues/5187")

	tester := caddytest.NewTester(t)
	tester.InitServer(` 
	{
	  order vulcain before request_header

	  skip_install_trust
	  admin localhost:2999
	  http_port     9080
	  https_port    9443
	}

	localhost:9080 {
	  file_server {
	  	root fixtures
	  }
	  header /books* Content-Type application/ld+json
	  vulcain
	}`, "caddyfile")

	resp, _ := tester.AssertGetResponse(`http://localhost:9080/books.jsonld?preload="/hydra:member/*"&fields="/hydra:member/*","/foo/0/bar/*/a"`, 200, `{"hydra:member":["/books-1.jsonld"],"foo":[{"bar":[{"a":"b"},{}]}]}`)

	// Unfortunately, Go's HTTP client doesn't support Pushes yet (https://github.com/golang/go/issues/18594), so we test the fallback
	if !reflect.DeepEqual(resp.Header["Link"], []string{"</books-1.jsonld>; rel=preload; as=fetch"}) {
		t.Errorf("missing link header")
	}
}
