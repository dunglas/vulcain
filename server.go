package vulcain

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"os"
	"os/signal"

	"github.com/gorilla/handlers"
	"go.uber.org/zap"
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
	var (
		logger *zap.Logger
		err    error
	)
	if options.Debug {
		logger, err = zap.NewDevelopment()
	} else {
		logger, err = zap.NewProduction()
	}

	if err != nil {
		panic(err)
	}

	return &server{
		options: options,
		vulcain: New(WithOpenAPIFile(options.OpenAPIFile), WithMaxPushes(options.MaxPushes), WithLogger(logger)),
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
		if !s.vulcain.CanApply(rw, req, resp.StatusCode, resp.Header) {
			return nil
		}

		newBody, err := s.vulcain.Apply(req, rw, resp.Body, resp.Header)
		if newBody == nil {
			return err
		}

		newBodyBuffer := bytes.NewBuffer(newBody)
		resp.Body = ioutil.NopCloser(newBodyBuffer)

		return nil
	}
	rp.ErrorHandler = func(rw http.ResponseWriter, req *http.Request, err error) {
		// Adapted from the default ErrorHandler
		s.vulcain.logger.Error("http: proxy error", zap.Error(err))
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
			s.vulcain.logger.Error(err.Error())
		}
		s.vulcain.logger.Info("my baby shot me down")
		close(idleConnsClosed)
	}()

	acme := len(s.options.AcmeHosts) > 0
	var err error

	if !acme && s.options.CertFile == "" && s.options.KeyFile == "" {
		s.vulcain.logger.Info("vulcain started", zap.String("protocol", "http"), zap.String("addr", s.options.Addr))
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
					s.vulcain.logger.Fatal(err.Error())
				}
			}()
		}

		s.vulcain.logger.Info("vulcain started", zap.String("protocol", "https"), zap.String("addr", s.options.Addr))
		err = s.server.ListenAndServeTLS(s.options.CertFile, s.options.KeyFile)
	}

	if err != http.ErrServerClosed {
		s.vulcain.logger.Fatal(err.Error())
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
		handlers.RecoveryLogger(zapRecoveryHandlerLogger{s.vulcain.logger}),
		handlers.PrintRecoveryStack(s.options.Debug),
	)(loggingHandler)

	return recoveryHandler
}

type zapRecoveryHandlerLogger struct {
	logger *zap.Logger
}

func (z zapRecoveryHandlerLogger) Println(args ...interface{}) {
	z.logger.Error(fmt.Sprint(args...))
}
