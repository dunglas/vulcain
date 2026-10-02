package vulcain

import (
	"testing"

	"github.com/dunglas/httpsfv"
	"github.com/stretchr/testify/assert"
)

func TestRootNode(t *testing.T) {
	n := &node{}
	assert.Empty(t, n.httpList(preload, ""))
}

func TestImportPointers(t *testing.T) {
	n := &node{}
	n.importPointers(preload, httpsfv.List{httpsfv.NewItem("/foo"), httpsfv.NewItem("/bar/foo"), httpsfv.NewItem("/foo/*"), httpsfv.NewItem("/bar/foo/*/baz")})
	n.importPointers(fields, httpsfv.List{httpsfv.NewItem("/foo/bat"), httpsfv.NewItem("/baz"), httpsfv.NewItem("/baz/*"), httpsfv.NewItem("/baz")})

	assert.Equal(t, httpsfv.List{httpsfv.NewItem("/foo/*"), httpsfv.NewItem("/bar/foo/*/baz")}, n.httpList(preload, ""))
	assert.Equal(t, httpsfv.List{httpsfv.NewItem("/foo/bat"), httpsfv.NewItem("/baz"), httpsfv.NewItem("/baz")}, n.httpList(fields, ""))
}

func TestString(t *testing.T) {
	n := &node{}
	n.importPointers(preload, httpsfv.List{httpsfv.NewItem("/foo"), httpsfv.NewItem("/bar/foo"), httpsfv.NewItem("/foo/*"), httpsfv.NewItem("/bar/foo/*/baz")})

	assert.Equal(t, "/", n.String())
	assert.Equal(t, "/foo", n.children[0].String())
	assert.Equal(t, "/bar/foo", n.children[1].children[0].String())
	assert.Equal(t, "/foo/*", n.children[0].children[0].String())
	assert.Equal(t, "/bar/foo/*/baz", n.children[1].children[0].children[0].children[0].String())
}

func TestChildDirectivePrecedence(t *testing.T) {
	for _, exactType := range []_type{preload, fields} {
		t.Run(map[_type]string{preload: "preload exact", fields: "fields exact"}[exactType], func(t *testing.T) {
			wildcardType := fields
			if exactType == fields {
				wildcardType = preload
			}

			exact := httpsfv.NewItem("/links/first/title")
			exact.Params.Add("lang", "en")
			wildcard := httpsfv.NewItem("/links/*/href")
			wildcard.Params.Add("type", "application/json")
			n := &node{}
			n.importPointers(exactType, httpsfv.List{exact})
			n.importPointers(wildcardType, httpsfv.List{wildcard})
			originalPreload := n.httpList(preload, "")
			originalFields := n.httpList(fields, "")

			links := n.child([]byte("links"))
			assert.Same(t, n.children[0], links)
			first := links.child([]byte("first"))
			assert.True(t, first.preload)
			assert.True(t, first.fields)
			assert.True(t, first.hasChildren(preload))
			assert.True(t, first.hasChildren(fields))
			assert.Equal(t, httpsfv.List{httpsfv.Item{Value: "/title", Params: exact.Params}}, first.httpList(exactType, ""))
			assert.Equal(t, httpsfv.List{httpsfv.Item{Value: "/href", Params: wildcard.Params}}, first.httpList(wildcardType, ""))
			if exactType == preload {
				assert.Equal(t, "/links/first", first.String())
			} else {
				assert.Equal(t, "/links/*", first.String())
			}

			other := links.child([]byte("other"))
			assert.Same(t, links.children[1], other)
			assert.Nil(t, first.child([]byte("missing")))
			assert.Equal(t, originalPreload, n.httpList(preload, ""))
			assert.Equal(t, originalFields, n.httpList(fields, ""))
		})
	}
}

func TestChildDoesNotRestoreShadowedDirective(t *testing.T) {
	for _, wildcardType := range []_type{preload, fields} {
		t.Run(map[_type]string{preload: "preload wildcard", fields: "fields wildcard"}[wildcardType], func(t *testing.T) {
			exactType := fields
			if wildcardType == fields {
				exactType = preload
			}

			n := &node{}
			n.importPointers(wildcardType, httpsfv.List{httpsfv.NewItem("/links/*/href")})
			n.importPointers(exactType, httpsfv.List{httpsfv.NewItem("/links/*/href"), httpsfv.NewItem("/links/first/title")})
			href := n.child([]byte("links")).child([]byte("first")).child([]byte("href"))
			assert.Equal(t, wildcardType == preload, href.preload)
			assert.Equal(t, wildcardType == fields, href.fields)
			assert.Empty(t, href.httpList(exactType, "/href"))
			assert.Empty(t, href.httpList(exactType, ""))
			assert.False(t, href.hasChildren(exactType))
			assert.Nil(t, href.child([]byte("missing")))
		})
	}
}

func TestHTTPListIgnoresOtherDirectiveDescendants(t *testing.T) {
	for _, leafType := range []_type{preload, fields} {
		t.Run(map[_type]string{preload: "preload leaf", fields: "fields leaf"}[leafType], func(t *testing.T) {
			otherType := fields
			if leafType == fields {
				otherType = preload
			}

			leaf := httpsfv.NewItem("/links/*/href")
			leaf.Params.Add("lang", "en")
			n := &node{}
			n.importPointers(leafType, httpsfv.List{leaf})
			n.importPointers(otherType, httpsfv.List{httpsfv.NewItem("/links/*/href/next"), httpsfv.NewItem("/links/first/id")})
			first := n.child([]byte("links")).child([]byte("first"))
			assert.Equal(t, httpsfv.List{httpsfv.Item{Value: "/href", Params: leaf.Params}}, first.httpList(leafType, ""))
			href := first.child([]byte("href"))
			assert.Equal(t, httpsfv.List{httpsfv.Item{Value: "/href", Params: leaf.Params}}, href.httpList(leafType, "/href"))
			assert.Equal(t, httpsfv.List{leaf}, n.httpList(leafType, ""))
			assert.False(t, href.hasChildren(otherType))
			assert.Nil(t, href.child([]byte("next")))
		})
	}
}
