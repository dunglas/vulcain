package gateway

import (
	"testing"

	"github.com/dunglas/httpsfv"
	"github.com/stretchr/testify/assert"
)

func TestRootNode(t *testing.T) {
	n := &node{}
	assert.Empty(t, n.httpList(Preload, ""))
}

func TestImportPointers(t *testing.T) {
	n := &node{}
	n.importPointers(Preload, httpsfv.List{httpsfv.NewItem("/foo"), httpsfv.NewItem("/bar/foo"), httpsfv.NewItem("/foo/*"), httpsfv.NewItem("/bar/foo/*/baz")})
	n.importPointers(Fields, httpsfv.List{httpsfv.NewItem("/foo/bat"), httpsfv.NewItem("/baz"), httpsfv.NewItem("/baz/*"), httpsfv.NewItem("/baz")})

	assert.Equal(t, httpsfv.List{httpsfv.NewItem("/foo/*"), httpsfv.NewItem("/bar/foo/*/baz")}, n.httpList(Preload, ""))
	assert.Equal(t, httpsfv.List{httpsfv.NewItem("/foo/bat"), httpsfv.NewItem("/baz/*")}, n.httpList(Fields, ""))
}

func TestString(t *testing.T) {
	n := &node{}
	n.importPointers(Preload, httpsfv.List{httpsfv.NewItem("/foo"), httpsfv.NewItem("/bar/foo"), httpsfv.NewItem("/foo/*"), httpsfv.NewItem("/bar/foo/*/baz")})

	assert.Equal(t, "/", n.String())
	assert.Equal(t, "/foo", n.children[0].String())
	assert.Equal(t, "/bar/foo", n.children[1].children[0].String())
	assert.Equal(t, "/foo/*", n.children[0].children[0].String())
	assert.Equal(t, "/bar/foo/*/baz", n.children[1].children[0].children[0].children[0].String())
}
