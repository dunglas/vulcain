package gateway

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRootNode(t *testing.T) {
	n := &node{}
	assert.Equal(t, []string{"/"}, n.strings(Preload, ""))
}

func TestImportPointers(t *testing.T) {
	n := &node{}
	n.importPointers(Preload, []string{"/foo", "/bar/foo", "/foo/*", "/bar/foo/*/baz"})
	n.importPointers(Fields, []string{"/foo/bat", "/baz", "/baz/*", "/baz"})

	assert.Equal(t, []string{"/foo/*", "/bar/foo/*/baz"}, n.strings(Preload, ""))
	assert.Equal(t, []string{"/foo/bat", "/baz/*"}, n.strings(Fields, ""))
}

func TestString(t *testing.T) {
	n := &node{}
	n.importPointers(Preload, []string{"/foo", "/bar/foo", "/foo/*", "/bar/foo/*/baz"})

	assert.Equal(t, "/", n.String())
	assert.Equal(t, "/foo", n.children[0].String())
	assert.Equal(t, "/bar/foo", n.children[1].children[0].String())
	assert.Equal(t, "/foo/*", n.children[0].children[0].String())
	assert.Equal(t, "/bar/foo/*/baz", n.children[1].children[0].children[0].children[0].String())
}
