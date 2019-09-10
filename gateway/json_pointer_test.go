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
