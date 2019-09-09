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
	n := newPointersTree(true, true)
	n.importPointers(Preload, []string{"/foo", "/bar/baz"})
	n.importPointers(Fields, []string{"/foo", "/baz/bar"})

	u, _ := url.Parse("/test")
	urlRewriter(u, n)

	assert.Equal(t, "/test?fields=%2Ffoo&fields=%2Fbaz%2Fbar&preload=%2Ffoo&preload=%2Fbar%2Fbaz", u.String())
}

func TestTraverseJSONFields(t *testing.T) {
	n := newPointersTree(true, true)
	n.importPointers(Fields, []string{"/notexist", "/bar"})

	result := traverseJSON([]byte(`{"foo": "f", "bar": "b"}`), n, true, urlRewriter)
	assert.Equal(t, `{"bar":"b"}`, string(result))
}

func TestTraverseJSONFieldsRewriteURL(t *testing.T) {
	n := newPointersTree(true, true)
	n.importPointers(Fields, []string{"/foo/*/bar"})

	result := traverseJSON([]byte(`{"foo": ["/a", "/b"]}`), n, true, urlRewriter)
	assert.Equal(t, `{"foo":["/a?fields=%2Fbar","/b?fields=%2Fbar"]}`, string(result))
}
