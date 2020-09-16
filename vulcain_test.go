package vulcain

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/dunglas/httpsfv"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	g := New(Options{MaxPushes: -1})
	assert.NotNil(t, g)
}

func TestParseRelation(t *testing.T) {
	v := New(Options{OpenAPIFile: openapiFixture, MaxPushes: -1})

	u, _ := url.Parse("/oa/books/123")

	u, _, _ = v.parseRelation("/author", "123", v.getOpenAPIRoute(u, nil, false))
	assert.Equal(t, "/oa/authors/123", u.String())

	u, _, _ = v.parseRelation("/invalid", " http://foo.com", nil)
	assert.Nil(t, u)
}

func TestCanParse(t *testing.T) {
	assert.False(t, canParse(
		http.Header{"Content-Type": []string{"text/xml"}},
		&http.Request{},
		httpsfv.List{httpsfv.NewItem("foo")}, httpsfv.List{},
	))

	assert.False(t, canParse(
		http.Header{"Content-Type": []string{"application/json"}},
		&http.Request{},
		httpsfv.List{},
		httpsfv.List{},
	))

	assert.False(t, canParse(
		http.Header{"Content-Type": []string{"application/json"}},
		&http.Request{Header: http.Header{"Prefer": []string{"selector=css"}}},
		httpsfv.List{httpsfv.NewItem("foo")},
		httpsfv.List{},
	))

	assert.True(t, canParse(
		http.Header{"Content-Type": []string{"application/json"}},
		&http.Request{Header: http.Header{"Prefer": []string{"selector=json-pointer"}}},
		httpsfv.List{httpsfv.NewItem("foo")},
		httpsfv.List{},
	))

	assert.True(t, canParse(
		http.Header{"Content-Type": []string{"application/ld+json"}},
		&http.Request{},
		httpsfv.List{httpsfv.NewItem("foo")},
		httpsfv.List{},
	))

	assert.True(t, canParse(
		http.Header{"Content-Type": []string{"application/ld+json"}},
		&http.Request{},
		httpsfv.List{httpsfv.NewItem("foo")},
		httpsfv.List{},
	))
}
