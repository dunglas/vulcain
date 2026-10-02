package vulcain_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/dunglas/vulcain"
	"github.com/stretchr/testify/assert"
)

func TestApplyFieldsAndPreloadMatchIndependently(t *testing.T) {
	for _, test := range []struct {
		name    string
		body    string
		fields  string
		preload string
		query   bool
		want    string
		links   []string
	}{
		{
			name:    "array objects",
			body:    `{"items":[{"title":"A","author":"/authors/1"},{"title":"B","author":"/authors/2"}]}`,
			fields:  `"/items/*"`,
			preload: `"/items/0/author"`,
			want:    `{"items":[{"title":"A","author":"/authors/1"},{"title":"B","author":"/authors/2"}]}`,
			links:   []string{"</authors/1>; rel=preload; as=fetch"},
		},
		{
			name:    "nested object fields",
			body:    `{"items":[{"title":"A","author":"/authors/1"},{"title":"B","author":"/authors/2"}]}`,
			fields:  `"/items/*/author"`,
			preload: `"/items/0/author"`,
			want:    `{"items":[{"author":"/authors/1"},{"author":"/authors/2"}]}`,
			links:   []string{"</authors/1>; rel=preload; as=fetch"},
		},
		{
			name:    "wildcard fields and exact preload URLs",
			body:    `{"items":["/books/1","/books/2"]}`,
			fields:  `"/items/*/title"`,
			preload: `"/items/0/author"`,
			query:   true,
			want:    `{"items":["/books/1?fields=%22%2Ftitle%22&preload=%22%2Fauthor%22","/books/2?fields=%22%2Ftitle%22"]}`,
			links:   []string{"</books/1?fields=%22%2Ftitle%22&preload=%22%2Fauthor%22>; rel=preload; as=fetch"},
		},
		{
			name:    "exact fields and wildcard preload URLs",
			body:    `{"items":["/books/1","/books/2"]}`,
			fields:  `"/items/0/title"`,
			preload: `"/items/*/author"`,
			query:   true,
			want:    `{"items":["/books/1?fields=%22%2Ftitle%22&preload=%22%2Fauthor%22"]}`,
			links:   []string{"</books/1?fields=%22%2Ftitle%22&preload=%22%2Fauthor%22>; rel=preload; as=fetch"},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			body, headers := applyDirectives(t, test.body, test.fields, test.preload, test.query)
			assert.Equal(t, test.want, string(body))
			assert.Equal(t, test.links, headers.Values("Link"))
		})
	}
}

func TestApplyFieldsPreservesAcceptedInvalidUnicode(t *testing.T) {
	for _, test := range []struct {
		name   string
		body   string
		fields string
		want   string
	}{
		{
			name:   "unpaired surrogate value",
			body:   `{"title":"\ud800","ignored":true}`,
			fields: `"/title"`,
			want:   `{"title":"\ud800"}`,
		},
		{
			name:   "unpaired surrogate member name",
			body:   `{"\ud800":{"title":"A","ignored":true}}`,
			fields: `"/*/title"`,
			want:   `{"\ud800":{"title":"A"}}`,
		},
		{
			name:   "invalid UTF-8 value",
			body:   "{\"title\":\"\xff\",\"ignored\":true}",
			fields: `"/title"`,
			want:   "{\"title\":\"\xff\"}",
		},
		{
			name:   "invalid UTF-8 member name",
			body:   "{\"\xff\":{\"title\":\"A\",\"ignored\":true}}",
			fields: `"/*/title"`,
			want:   "{\"\xff\":{\"title\":\"A\"}}",
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			body, _ := applyDirectives(t, test.body, test.fields, "", false)
			assert.Equal(t, test.want, string(body))
		})
	}
}

func applyDirectives(t *testing.T, body, fields, preload string, query bool) ([]byte, http.Header) {
	t.Helper()
	v := vulcain.New()
	rw := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "http://example.com/books", nil)
	if query {
		q := req.URL.Query()
		q.Set("fields", fields)
		q.Set("preload", preload)
		req.URL.RawQuery = q.Encode()
	} else {
		req.Header.Set("Fields", fields)
		req.Header.Set("Preload", preload)
	}
	req = req.WithContext(v.CreateRequestContext(rw, req))
	defer v.Finish(req, false)
	headers := http.Header{"Content-Type": []string{"application/json"}}
	if !v.IsValidRequest(req) || !v.IsValidResponse(req, http.StatusOK, headers) {
		t.Fatal("invalid Vulcain request or response")
	}
	result, err := v.Apply(req, rw, strings.NewReader(body), headers)
	if err != nil {
		t.Fatal(err)
	}

	return result, headers
}
