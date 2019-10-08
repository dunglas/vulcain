package gateway

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewOpenAPI(t *testing.T) {
	assert.NotNil(t, newOpenAPI("../fixtures/openapi.yaml"))
	assert.Panics(t, func() {
		newOpenAPI("notexists")
	})
}

func TestGetRoute(t *testing.T) {
	oa := newOpenAPI("../fixtures/openapi.yaml")

	u, _ := url.Parse("/oa/books/123")
	assert.NotNil(t, oa.getRoute(u))

	u, _ = url.Parse("/notexists")
	assert.Nil(t, oa.getRoute(u))
}

func TestGetRelation(t *testing.T) {
	oa := newOpenAPI("../fixtures/openapi.yaml")

	u, _ := url.Parse("/oa/books/123")
	r := oa.getRelation(oa.getRoute(u), "/author", "456")
	assert.Equal(t, "/oa/authors/456", r)

	u, _ = url.Parse("/oa/books.json")
	r = oa.getRelation(oa.getRoute(u), "/member/*", "1936")
	assert.Equal(t, "/oa/books/1936", r)

	u, _ = url.Parse("/oa/books.json")
	r = oa.getRelation(oa.getRoute(u), "/notexists", "1891")
	assert.Equal(t, "", r)
}

func TestGenerateLink(t *testing.T) {
	oa := newOpenAPI("../fixtures/openapi.yaml")
	l := oa.generateLink("notexists", "nestor", "makhno")
	assert.Equal(t, "", l)
}
