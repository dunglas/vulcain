package vulcain

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

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

func TestCanApply(t *testing.T) {
	v := New(Options{})
	assert.False(t, v.CanApply(
		&httptest.ResponseRecorder{},
		&http.Request{URL: &url.URL{}},
		200,
		http.Header{"Content-Type": []string{"text/xml"}},
	))

	assert.False(t, v.CanApply(
		&httptest.ResponseRecorder{},
		&http.Request{URL: &url.URL{}},
		200,
		http.Header{"Content-Type": []string{"application/json"}},
	))

	assert.False(t, v.CanApply(
		&httptest.ResponseRecorder{},
		&http.Request{
			URL:    &url.URL{},
			Header: http.Header{"Preload": []string{`"foo"`}, "Prefer": []string{"selector=css"}},
		},
		200,
		http.Header{"Content-Type": []string{"application/json"}},
	))

	assert.False(t, v.CanApply(
		&httptest.ResponseRecorder{},
		&http.Request{
			URL:    &url.URL{},
			Header: http.Header{"Preload": []string{`"foo"`}, "Prefer": []string{"selector=css"}},
		},
		500,
		http.Header{"Content-Type": []string{"application/json"}},
	))

	assert.True(t, v.CanApply(
		&httptest.ResponseRecorder{},
		&http.Request{
			URL:    &url.URL{},
			Header: http.Header{"Preload": []string{`"foo"`}, "Prefer": []string{"selector=json-pointer"}},
		},
		200,
		http.Header{"Content-Type": []string{"application/json"}},
	))

	assert.True(t, v.CanApply(
		&httptest.ResponseRecorder{},
		&http.Request{
			URL:    &url.URL{},
			Header: http.Header{"Preload": []string{`"foo"`}},
		},
		200,
		http.Header{"Content-Type": []string{"application/ld+json"}},
	))

	assert.True(t, v.CanApply(
		&httptest.ResponseRecorder{},
		&http.Request{
			URL:    &url.URL{},
			Header: http.Header{"Preload": []string{`"foo"`}},
		},
		200,
		http.Header{"Content-Type": []string{"application/ld+json"}},
	))

	assert.True(t, v.CanApply(
		&httptest.ResponseRecorder{},
		&http.Request{
			URL: &url.URL{RawQuery: `preload="foo"`},
		},
		200,
		http.Header{"Content-Type": []string{"application/ld+json"}},
	))
}
