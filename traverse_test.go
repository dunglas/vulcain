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
	n.importPointers(preload, httpsfv.List{httpsfv.NewItem("/foo/*"), httpsfv.NewItem("/bar/baz")})
	n.importPointers(fields, httpsfv.List{httpsfv.NewItem("/foo/*"), httpsfv.NewItem("/baz/bar")})

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
	n.importPointers(fields, httpsfv.List{httpsfv.NewItem("/notexist"), httpsfv.NewItem("/bar")})

	result := New().traverseJSON([]byte(`{"foo": "f", "bar": "b"}`), n, true, urlRewriteRelationHandler)
	assert.Equal(t, `{"bar":"b"}`, string(result))
}

func TestTraverseJSONFieldsRewriteURL(t *testing.T) {
	n := &node{}
	n.importPointers(fields, httpsfv.List{httpsfv.NewItem("/foo/*/bar")})

	result := New().traverseJSON([]byte(`{"foo": ["/a", "/b"]}`), n, true, urlRewriteRelationHandler)
	assert.Equal(t, `{"foo":["/a?fields=%22%2Fbar%22","/b?fields=%22%2Fbar%22"]}`, string(result))
}

func TestTraverseJSONPreload(t *testing.T) {
	n := &node{}
	n.importPointers(preload, httpsfv.List{httpsfv.NewItem("/notexist"), httpsfv.NewItem("/bar")})

	result := New().traverseJSON([]byte(`{"foo": "/foo", "bar": "/bar"}`), n, false, urlRewriteRelationHandler)
	assert.Equal(t, `{"foo": "/foo", "bar": "/bar"}`, string(result))
}

func TestTraverseJSONPreloadRewriteURL(t *testing.T) {
	n := &node{}
	n.importPointers(preload, httpsfv.List{httpsfv.NewItem("/foo/*/rel"), httpsfv.NewItem("/bar/baz")})

	result := New().traverseJSON([]byte(`{"foo": ["/a", "/b"], "bar": "/bar"}`), n, false, urlRewriteRelationHandler)
	assert.Equal(t, `{"foo": ["/a?preload=%22%2Frel%22", "/b?preload=%22%2Frel%22"], "bar": "/bar?preload=%22%2Fbaz%22"}`, string(result))
}

func TestTraverseJSONPreloadAndFieldsRewriteURL(t *testing.T) {
	n := &node{}
	n.importPointers(preload, httpsfv.List{httpsfv.NewItem("/notexist"), httpsfv.NewItem("/foo/*/rel"), httpsfv.NewItem("/bar/baz"), httpsfv.NewItem("/baz")})
	n.importPointers(fields, httpsfv.List{httpsfv.NewItem("/foo/*"), httpsfv.NewItem("/bar/baz"), httpsfv.NewItem("/notexist")})

	result := New().traverseJSON([]byte(`{"foo": ["/a", "/b"], "bar": "/bar", "baz": "/baz"}`), n, true, urlRewriteRelationHandler)
	assert.Equal(t, `{"foo":["/a?preload=%22%2Frel%22","/b?preload=%22%2Frel%22"],"bar":"/bar?fields=%22%2Fbaz%22&preload=%22%2Fbaz%22"}`, string(result))
}

func TestTraverseJSONFieldsKeyWithDot(t *testing.T) {
	n := &node{}
	n.importPointers(fields, httpsfv.List{httpsfv.NewItem("/a.b"), httpsfv.NewItem("/c~1d")})

	result := New().traverseJSON([]byte(`{"a.b": "x", "a": {"b": "y"}, "c/d": "z"}`), n, true, urlRewriteRelationHandler)
	assert.Equal(t, `{"a.b":"x","c/d":"z"}`, string(result))
}

func TestTraverseJSONPreloadArrayIndex(t *testing.T) {
	n := &node{}
	n.importPointers(preload, httpsfv.List{httpsfv.NewItem("/foo/1/rel")})

	result := New().traverseJSON([]byte(`{"foo": ["/a", "/b"]}`), n, false, urlRewriteRelationHandler)
	assert.Equal(t, `{"foo": ["/a", "/b?preload=%22%2Frel%22"]}`, string(result))
}

func TestTraverseJSONPreloadNumber(t *testing.T) {
	n := &node{}
	n.importPointers(preload, httpsfv.List{httpsfv.NewItem("/foo/rel")})

	result := New().traverseJSON([]byte(`{"foo": 1.5}`), n, false, urlRewriteRelationHandler)
	assert.Equal(t, `{"foo": "1.5?preload=%22%2Frel%22"}`, string(result))
}

func TestTraverseJSONInvalid(t *testing.T) {
	n := &node{}
	n.importPointers(fields, httpsfv.List{httpsfv.NewItem("/foo")})

	result := New().traverseJSON([]byte(`{"foo": "/foo", "bar"`), n, true, urlRewriteRelationHandler)
	assert.Equal(t, `{"foo": "/foo", "bar"`, string(result))
}

func TestTraverseJSONPreloadObjectWildcard(t *testing.T) {
	n := &node{}
	n.importPointers(preload, httpsfv.List{httpsfv.NewItem("/links/*/href/rel")})

	result := New().traverseJSON([]byte(`{"links": {"self": {"href": "/a"}, "author": {"href": "/b"}}}`), n, false, urlRewriteRelationHandler)
	assert.Equal(t, `{"links": {"self": {"href": "/a?preload=%22%2Frel%22"}, "author": {"href": "/b?preload=%22%2Frel%22"}}}`, string(result))
}

func TestTraverseJSONFieldsObjectWildcard(t *testing.T) {
	n := &node{}
	n.importPointers(fields, httpsfv.List{httpsfv.NewItem("/links/*/href")})

	result := New().traverseJSON([]byte(`{"links": {"self": {"href": "/a", "title": "A"}, "author": {"href": "/b"}}, "x": "y"}`), n, true, urlRewriteRelationHandler)
	assert.Equal(t, `{"links":{"self":{"href":"/a"},"author":{"href":"/b"}}}`, string(result))
}

func TestTraverseJSONFieldsExactMatchOverWildcard(t *testing.T) {
	n := &node{}
	n.importPointers(fields, httpsfv.List{httpsfv.NewItem("/m/*/a"), httpsfv.NewItem("/m/b/c"), httpsfv.NewItem("/m/~2/d")})

	result := New().traverseJSON([]byte(`{"m": {"b": {"a": "1", "c": "2"}, "e": {"a": "3", "c": "4"}, "*": {"a": "5", "d": "6"}}}`), n, true, urlRewriteRelationHandler)
	assert.Equal(t, `{"m":{"b":{"c":"2"},"e":{"a":"3"},"*":{"d":"6"}}}`, string(result))
}

func TestTraverseJSONFieldsArrayIndices(t *testing.T) {
	for _, test := range []struct {
		name     string
		pointers []string
		body     string
		want     string
	}{
		{"single", []string{"/foo/2"}, `{"foo":["a","b","c","d"]}`, `{"foo":[null,null,"c"]}`},
		{"multiple", []string{"/foo/3/id", "/foo/1/id"}, `{"foo":[{"id":"a"},{"id":"b"},{"id":"c"},{"id":"d"}]}`, `{"foo":[null,{"id":"b"},null,{"id":"d"}]}`},
		{"first", []string{"/foo/0", "/foo/2"}, `{"foo":["a","b","c","d"]}`, `{"foo":["a",null,"c"]}`},
		{"missing", []string{"/foo/8"}, `{"foo":["a","b","c"]}`, `{"foo":[]}`},
	} {
		t.Run(test.name, func(t *testing.T) {
			n := &node{}
			var pointers httpsfv.List
			for _, pointer := range test.pointers {
				pointers = append(pointers, httpsfv.NewItem(pointer))
			}
			n.importPointers(fields, pointers)

			result := New().traverseJSON([]byte(test.body), n, true, urlRewriteRelationHandler)
			assert.Equal(t, test.want, string(result))
		})
	}
}
