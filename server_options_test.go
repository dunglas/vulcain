package vulcain

import (
	"net/url"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewOptionsFromEnv(t *testing.T) {
	testEnv := map[string]string{
		"UPSTREAM":      "http://example.com",
		"SERVER_PUSH":   "1",
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
		os.Setenv(k, v)
		defer os.Unsetenv(k)
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
	os.Setenv("CERT_FILE", "foo")
	defer os.Unsetenv("CERT_FILE")

	_, err := NewOptionsFromEnv()
	assert.EqualError(t, err, "The following environment variable must be defined: [KEY_FILE]")
}

func TestMissingCertFile(t *testing.T) {
	os.Setenv("KEY_FILE", "foo")
	defer os.Unsetenv("KEY_FILE")

	_, err := NewOptionsFromEnv()
	assert.EqualError(t, err, "The following environment variable must be defined: [CERT_FILE]")
}

func TestInvalidDuration(t *testing.T) {
	vars := [2]string{"READ_TIMEOUT", "WRITE_TIMEOUT"}
	for _, elem := range vars {
		os.Setenv(elem, "1 MN (invalid)")
		defer os.Unsetenv(elem)
		_, err := NewOptionsFromEnv()
		assert.EqualError(t, err, elem+`: time: unknown unit " MN (invalid)" in duration "1 MN (invalid)"`)

		os.Unsetenv(elem)
	}
}

func TestInvalidUpstream(t *testing.T) {
	os.Setenv("UPSTREAM", " http://foo.com")
	defer os.Unsetenv("UPSTREAM")
	_, err := NewOptionsFromEnv()
	assert.EqualError(t, err, `parse " http://foo.com": first path segment in URL cannot contain colon`)
}

func TestInvalidMaxPushes(t *testing.T) {
	os.Setenv("MAX_PUSHES", "invalid")
	defer os.Unsetenv("MAX_PUSHES")
	_, err := NewOptionsFromEnv()
	assert.EqualError(t, err, `MAX_PUSHES: invalid value "invalid" (strconv.Atoi: parsing "invalid": invalid syntax)`)
}
