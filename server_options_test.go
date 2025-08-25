package vulcain

import (
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewOptionsFromEnv(t *testing.T) {
	testEnv := map[string]string{
		"UPSTREAM":      "http://example.com",
		"EARLY_HINTS":   "1",
		"MAX_PUSHES":    "-1",
		"ACME_CERT_DIR": "/tmp",
		"ACME_HOSTS":    "example.com,example.org",
		"ADDR":          "127.0.0.1:8080",
		"CERT_FILE":     "foo",
		"COMPRESS":      "0",
		"DEBUG":         "1",
		"KEY_FILE":      "bar",
		"READ_TIMEOUT":  "1m",
		"WRITE_TIMEOUT": "40s",
		"OPENAPI_FILE":  "openapi.yaml",
	}
	for k, v := range testEnv {
		t.Setenv(k, v)
	}

	u, _ := url.Parse("http://example.com")
	opts, err := NewOptionsFromEnv()
	assert.Equal(t, &ServerOptions{
		true,
		"127.0.0.1:8080",
		u,
		true,
		-1,
		[]string{"example.com", "example.org"},
		"/tmp",
		"foo",
		"bar",
		time.Minute,
		40 * time.Second,
		false,
		"openapi.yaml",
	}, opts)
	assert.Nil(t, err)
}

func TestMissingKeyFile(t *testing.T) {
	t.Setenv("CERT_FILE", "foo")

	_, err := NewOptionsFromEnv()
	assert.EqualError(t, err, "the following environment variable must be defined: [KEY_FILE]")
}

func TestMissingCertFile(t *testing.T) {
	t.Setenv("KEY_FILE", "foo")

	_, err := NewOptionsFromEnv()
	assert.EqualError(t, err, "the following environment variable must be defined: [CERT_FILE]")
}

func TestInvalidDuration(t *testing.T) {
	for _, elem := range [2]string{"READ_TIMEOUT", "WRITE_TIMEOUT"} {
		t.Run(elem, func(t *testing.T) {
			t.Setenv(elem, "1 MN (invalid)")

			_, err := NewOptionsFromEnv()
			assert.EqualError(t, err, elem+`: time: unknown unit " MN (invalid)" in duration "1 MN (invalid)"`)
		})
	}
}

func TestInvalidUpstream(t *testing.T) {
	t.Setenv("UPSTREAM", " https://foo.com")

	_, err := NewOptionsFromEnv()
	assert.EqualError(t, err, `parse " https://foo.com": first path segment in URL cannot contain colon`)
}

func TestInvalidMaxPushes(t *testing.T) {
	t.Setenv("MAX_PUSHES", "invalid")

	_, err := NewOptionsFromEnv()
	assert.EqualError(t, err, `MAX_PUSHES: invalid value "invalid" (strconv.Atoi: parsing "invalid": invalid syntax)`)
}
