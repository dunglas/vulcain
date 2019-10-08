package gateway

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnescape(t *testing.T) {
	assert.Equal(t, "~1/0*/", unescape("~01~10~2/"))
}

func TestUrlRewriter(t *testing.T) {
	n := &node{}
	n.importPointers(Preload, []string{"/foo/*", "/bar/baz"})
	n.importPointers(Fields, []string{"/foo/*", "/baz/bar"})

	u, _ := url.Parse("/test")
	urlRewriter(u, n)

	assert.Equal(t, "/test?fields=%2Ffoo%2F%2A&fields=%2Fbaz%2Fbar&preload=%2Ffoo%2F%2A&preload=%2Fbar%2Fbaz", u.String())
}

func urlRewriteRelationHandler(n *node, v string) string {
	u, _ := url.Parse(v)
	urlRewriter(u, n)

	return u.String()
}

func TestTraverseJSONFields(t *testing.T) {
	n := &node{}
	n.importPointers(Fields, []string{"/notexist", "/bar"})

	result := traverseJSON([]byte(`{"foo": "f", "bar": "b"}`), n, true, urlRewriteRelationHandler)
	assert.Equal(t, `{"bar":"b"}`, string(result))
}

func TestTraverseJSONFieldsRewriteURL(t *testing.T) {
	n := &node{}
	n.importPointers(Fields, []string{"/foo/*/bar"})

	result := traverseJSON([]byte(`{"foo": ["/a", "/b"]}`), n, true, urlRewriteRelationHandler)
	assert.Equal(t, `{"foo":["/a?fields=%2Fbar","/b?fields=%2Fbar"]}`, string(result))
}

func TestTraverseJSONPreload(t *testing.T) {
	n := &node{}
	n.importPointers(Preload, []string{"/notexist", "/bar"})

	result := traverseJSON([]byte(`{"foo": "/foo", "bar": "/bar"}`), n, false, urlRewriteRelationHandler)
	assert.Equal(t, `{"foo": "/foo", "bar": "/bar"}`, string(result))
}

func TestTraverseJSONPreloadRewriteURL(t *testing.T) {
	n := &node{}
	n.importPointers(Preload, []string{"/foo/*/rel", "/bar/baz"})

	result := traverseJSON([]byte(`{"foo": ["/a", "/b"], "bar": "/bar"}`), n, false, urlRewriteRelationHandler)
	assert.Equal(t, `{"foo": ["/a?preload=%2Frel", "/b?preload=%2Frel"], "bar": "/bar?preload=%2Fbaz"}`, string(result))
}

func TestTraverseJSONPreloadAndFieldsRewriteURL(t *testing.T) {
	n := &node{}
	n.importPointers(Preload, []string{"/notexist", "/foo/*/rel", "/bar/baz", "/baz"})
	n.importPointers(Fields, []string{"/foo/*", "/bar/baz", "/notexist"})

	result := traverseJSON([]byte(`{"foo": ["/a", "/b"], "bar": "/bar", "baz": "/baz"}`), n, true, urlRewriteRelationHandler)
	assert.Equal(t, `{"bar":"/bar?fields=%2Fbaz\u0026preload=%2Fbaz","foo":["/a?preload=%2Frel","/b?preload=%2Frel"]}`, string(result))
}
