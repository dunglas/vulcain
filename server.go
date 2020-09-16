package vulcain

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"os"
	"os/signal"

	"github.com/gorilla/handlers"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/acme/autocert"
)

// NewServerFromEnv creates a server using the configuration set in env vars
//
// Deprecated: use the Caddy server module or the standalone library instead
func NewServerFromEnv() (*server, error) {
	options, err := NewOptionsFromEnv()
	if err != nil {
		return nil, err
	}

	return NewServer(options), nil
}

// NewServer creates a Vulcain server
//
// Deprecated: use the Caddy server module or the standalone library instead
func NewServer(options *ServerOptions) *server {
	return &server{
		options: options,
		vulcain: New(Options{options.OpenAPIFile, options.MaxPushes}),
	}
}

type server struct {
	options *ServerOptions
	server  *http.Server
	vulcain *Vulcain
}

// ServeHTTP starts a reverse proxy and apply Vulcain queries on its response
//
// Deprecated: use the Caddy server module or the standalone library instead
func (s *server) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	rp := httputil.NewSingleHostReverseProxy(s.options.Upstream)
	rp.ModifyResponse = func(resp *http.Response) error {
		newBody, newHeaders, err := s.vulcain.Apply(req, rw, resp.Body, resp.Header)
		if newBody == nil {
			return err
		}

		newBodyBuffer := bytes.NewBuffer(newBody)
		resp.Body = ioutil.NopCloser(newBodyBuffer)
		resp.Header = newHeaders

		return nil
	}
	rp.ErrorHandler = func(rw http.ResponseWriter, req *http.Request, err error) {
		// Adapted from the default ErrorHandler
		log.Errorf("http: proxy error: %v", err)
		rw.WriteHeader(http.StatusBadGateway)
	}

	proto := "https"
	if req.TLS == nil {
		proto = "http"
	}

	req.Header.Set("X-Forwarded-Proto", proto)
	req.Header.Set("X-Forwarded-Host", req.Host)
	req.Header.Del("X-Forwarded-For")
	rp.ServeHTTP(rw, req)
}

// Serve starts the HTTP server
//
// Deprecated: use the Caddy server module or the standalone library instead
func (s *server) Serve() {
	s.server = &http.Server{
		Addr:         s.options.Addr,
		Handler:      s.chainHandlers(),
		ReadTimeout:  s.options.ReadTimeout,
		WriteTimeout: s.options.WriteTimeout,
	}

	idleConnsClosed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint

		if err := s.server.Shutdown(context.Background()); err != nil {
			log.Error(err)
		}
		log.Infoln("My Baby Shot Me Down")
		close(idleConnsClosed)
	}()

	acme := len(s.options.AcmeHosts) > 0
	var err error

	if !acme && s.options.CertFile == "" && s.options.KeyFile == "" {
		log.WithFields(log.Fields{"protocol": "http", "addr": s.options.Addr}).Info("Vulcain started")
		err = s.server.ListenAndServe()
	} else {
		// TLS
		if acme {
			certManager := &autocert.Manager{
				Prompt:     autocert.AcceptTOS,
				HostPolicy: autocert.HostWhitelist(s.options.AcmeHosts...),
			}
			if s.options.AcmeCertDir != "" {
				certManager.Cache = autocert.DirCache(s.options.AcmeCertDir)
			}
			s.server.TLSConfig = certManager.TLSConfig()

			// Mandatory for Let's Encrypt http-01 challenge
			go func() {
				if err := http.ListenAndServe(":http", certManager.HTTPHandler(nil)); err != nil {
					log.Fatal(err)
				}
			}()
		}

		log.WithFields(log.Fields{"protocol": "https", "addr": s.options.Addr}).Info("Vulcain started")
		err = s.server.ListenAndServeTLS(s.options.CertFile, s.options.KeyFile)
	}

	if err != http.ErrServerClosed {
		log.Fatal(err)
	}

	<-idleConnsClosed
}

// chainHandlers configures and chains handlers
func (s *server) chainHandlers() http.Handler {
	var compressHandler http.Handler
	if s.options.Compress {
		compressHandler = handlers.CompressHandler(s)
	} else {
		compressHandler = s
	}

	loggingHandler := handlers.CombinedLoggingHandler(os.Stderr, compressHandler)
	recoveryHandler := handlers.RecoveryHandler(
		handlers.RecoveryLogger(log.New()),
		handlers.PrintRecoveryStack(s.options.Debug),
	)(loggingHandler)

	return recoveryHandler
}
