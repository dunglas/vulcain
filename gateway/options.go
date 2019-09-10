package gateway

import (
	"fmt"
	"net/url"
	"os"
	"strings"
	"time"
)

// Options stores the gateway's options
type Options struct {
	Debug               bool
	Addr                string
	Upstream            *url.URL
	AcmeHosts           []string
	AcmeCertDir         string
	CertFile            string
	KeyFile             string
	ReadTimeout         time.Duration
	WriteTimeout        time.Duration
	Compress            bool
	UseForwardedHeaders bool
}

// NewOptionsFromEnv creates a new option instance from environment
// It returns an error if mandatory env env vars are missing
func NewOptionsFromEnv() (*Options, error) {
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

	options := &Options{
		os.Getenv("DEBUG") == "1",
		os.Getenv("ADDR"),
		upstream,
		splitVar(os.Getenv("ACME_HOSTS")),
		os.Getenv("ACME_CERT_DIR"),
		os.Getenv("CERT_FILE"),
		os.Getenv("KEY_FILE"),
		readTimeout,
		writeTimeout,
		os.Getenv("COMPRESS") != "0",
		os.Getenv("USE_FORWARDED_HEADERS") == "1",
	}

	missingEnv := make([]string, 0, 2)
	if len(options.CertFile) != 0 && len(options.KeyFile) == 0 {
		missingEnv = append(missingEnv, "KEY_FILE")
	}
	if len(options.KeyFile) != 0 && len(options.CertFile) == 0 {
		missingEnv = append(missingEnv, "CERT_FILE")
	}

	if len(missingEnv) > 0 {
		return nil, fmt.Errorf("The following environment variable must be defined: %s", missingEnv)
	}
	return options, nil
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
