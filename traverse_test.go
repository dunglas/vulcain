package vulcain

import (
	"net/url"
	"testing"

	"github.com/dunglas/httpsfv"
	"github.com/stretchr/testify/assert"
)

func TestUnescape(t *testing.T) {
	assert.Equal(t, "~1/0*/", unescape("~01~10~2/"))
}

func TestUrlRewriter(t *testing.T) {
	n := &node{}
	n.importPointers(Preload, httpsfv.List{httpsfv.NewItem("/foo/*"), httpsfv.NewItem("/bar/baz")})
	n.importPointers(Fields, httpsfv.List{httpsfv.NewItem("/foo/*"), httpsfv.NewItem("/baz/bar")})

	u, _ := url.Parse("/test")
	urlRewriter(u, n)

	assert.Equal(t, "/test?fields=%22%2Ffoo%2F%2A%22%2C+%22%2Fbaz%2Fbar%22&preload=%22%2Ffoo%2F%2A%22%2C+%22%2Fbar%2Fbaz%22", u.String())
}

func urlRewriteRelationHandler(n *node, v string) string {
	u, _ := url.Parse(v)
	urlRewriter(u, n)

	return u.String()
}

func TestTraverseJSONFields(t *testing.T) {
	n := &node{}
	n.importPointers(Fields, httpsfv.List{httpsfv.NewItem("/notexist"), httpsfv.NewItem("/bar")})

	result := traverseJSON([]byte(`{"foo": "f", "bar": "b"}`), n, true, urlRewriteRelationHandler)
	assert.Equal(t, `{"bar":"b"}`, string(result))
}

func TestTraverseJSONFieldsRewriteURL(t *testing.T) {
	n := &node{}
	n.importPointers(Fields, httpsfv.List{httpsfv.NewItem("/foo/*/bar")})

	result := traverseJSON([]byte(`{"foo": ["/a", "/b"]}`), n, true, urlRewriteRelationHandler)
	assert.Equal(t, `{"foo":["/a?fields=%22%2Fbar%22","/b?fields=%22%2Fbar%22"]}`, string(result))
}

func TestTraverseJSONPreload(t *testing.T) {
	n := &node{}
	n.importPointers(Preload, httpsfv.List{httpsfv.NewItem("/notexist"), httpsfv.NewItem("/bar")})

	result := traverseJSON([]byte(`{"foo": "/foo", "bar": "/bar"}`), n, false, urlRewriteRelationHandler)
	assert.Equal(t, `{"foo": "/foo", "bar": "/bar"}`, string(result))
}

func TestTraverseJSONPreloadRewriteURL(t *testing.T) {
	n := &node{}
	n.importPointers(Preload, httpsfv.List{httpsfv.NewItem("/foo/*/rel"), httpsfv.NewItem("/bar/baz")})

	result := traverseJSON([]byte(`{"foo": ["/a", "/b"], "bar": "/bar"}`), n, false, urlRewriteRelationHandler)
	assert.Equal(t, `{"foo": ["/a?preload=%22%2Frel%22", "/b?preload=%22%2Frel%22"], "bar": "/bar?preload=%22%2Fbaz%22"}`, string(result))
}

func TestTraverseJSONPreloadAndFieldsRewriteURL(t *testing.T) {
	n := &node{}
	n.importPointers(Preload, httpsfv.List{httpsfv.NewItem("/notexist"), httpsfv.NewItem("/foo/*/rel"), httpsfv.NewItem("/bar/baz"), httpsfv.NewItem("/baz")})
	n.importPointers(Fields, httpsfv.List{httpsfv.NewItem("/foo/*"), httpsfv.NewItem("/bar/baz"), httpsfv.NewItem("/notexist")})

	result := traverseJSON([]byte(`{"foo": ["/a", "/b"], "bar": "/bar", "baz": "/baz"}`), n, true, urlRewriteRelationHandler)
	assert.Equal(t, `{"foo":["/a?preload=%22%2Frel%22","/b?preload=%22%2Frel%22"],"bar":"/bar?fields=%22%2Fbaz%22\u0026preload=%22%2Fbaz%22"}`, string(result))
}
