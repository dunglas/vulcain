package vulcain

import (
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

// ServerOptions stores the server's options
//
// Deprecated: use the Caddy server module or the standalone library instead
type ServerOptions struct {
	Debug        bool
	Addr         string
	Upstream     *url.URL
	MaxPushes    int
	AcmeHosts    []string
	AcmeCertDir  string
	CertFile     string
	KeyFile      string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	Compress     bool
	OpenAPIFile  string
}

// NewOptionsFromEnv creates a new option instance from environment
// It returns an error if mandatory env env vars are missing
//
// Deprecated: use the Caddy server module or the standalone library instead
func NewOptionsFromEnv() (*ServerOptions, error) {
	var err error

	readTimeout, err := parseDurationFromEnvVar("READ_TIMEOUT", time.Duration(0))
	if err != nil {
		return nil, err
	}

	writeTimeout, err := parseDurationFromEnvVar("WRITE_TIMEOUT", time.Duration(0))
	if err != nil {
		return nil, err
	}

	upstream, err := url.Parse(os.Getenv("UPSTREAM"))
	if err != nil {
		return nil, err
	}

	var maxPushes int
	maxPushesStr := os.Getenv("MAX_PUSHES")
	if maxPushesStr == "" {
		maxPushes = -1
	} else {
		maxPushes, err = strconv.Atoi(maxPushesStr)
		if err != nil {
			return nil, fmt.Errorf(`MAX_PUSHES: invalid value "%s" (%s)`, maxPushesStr, err)
		}
	}

	o := &ServerOptions{
		os.Getenv("DEBUG") == "1",
		os.Getenv("ADDR"),
		upstream,
		maxPushes,
		splitVar(os.Getenv("ACME_HOSTS")),
		os.Getenv("ACME_CERT_DIR"),
		os.Getenv("CERT_FILE"),
		os.Getenv("KEY_FILE"),
		readTimeout,
		writeTimeout,
		os.Getenv("COMPRESS") != "0",
		os.Getenv("OPENAPI_FILE"),
	}

	missingEnv := make([]string, 0, 2)
	if len(o.CertFile) != 0 && len(o.KeyFile) == 0 {
		missingEnv = append(missingEnv, "KEY_FILE")
	}
	if len(o.KeyFile) != 0 && len(o.CertFile) == 0 {
		missingEnv = append(missingEnv, "CERT_FILE")
	}

	if len(missingEnv) > 0 {
		return nil, fmt.Errorf("The following environment variable must be defined: %s", missingEnv)
	}
	return o, nil
}

func splitVar(v string) []string {
	if v == "" {
		return []string{}
	}

	return strings.Split(v, ",")
}

func parseDurationFromEnvVar(k string, d time.Duration) (time.Duration, error) {
	v := os.Getenv(k)
	if v == "" {
		return d, nil
	}

	dur, err := time.ParseDuration(v)
	if err == nil {
		return dur, nil
	}

	return time.Duration(0), fmt.Errorf("%s: %s", k, err)
}
