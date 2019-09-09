package gateway

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRootNode(t *testing.T) {
	n := newPointersTree(true, true)
	assert.Equal(t, []string{"/"}, n.strings(preloadType, ""))
}

func TestImportPointers(t *testing.T) {
	n := newPointersTree(true, true)
	n.importPointers(preloadType, []string{"/foo", "/bar/foo", "/foo/*", "/bar/foo/*/baz"})
	n.importPointers(fieldsType, []string{"/foo/bat", "/baz", "/baz/*", "/baz"})

	assert.Equal(t, []string{"/foo/*", "/bar/foo/*/baz"}, n.strings(preloadType, ""))
	assert.Equal(t, []string{"/foo/bat", "/baz/*"}, n.strings(fieldsType, ""))
}
