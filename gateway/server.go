package gateway

import (
	"context"
	"net/http"
	"os"
	"os/signal"

	"github.com/gorilla/handlers"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/acme/autocert"
)

// Serve starts the HTTP server
func (g *Gateway) Serve() {
	g.server = &http.Server{
		Addr:         g.Options.Addr,
		Handler:      g.chainHandlers(),
		ReadTimeout:  g.Options.ReadTimeout,
		WriteTimeout: g.Options.WriteTimeout,
	}
	g.server.RegisterOnShutdown(func() {
		// todo
	})

	idleConnsClosed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint

		if err := g.server.Shutdown(context.Background()); err != nil {
			log.Error(err)
		}
		log.Infoln("My Baby Shot Me Down")
		close(idleConnsClosed)
	}()

	acme := len(g.Options.AcmeHosts) > 0
	var err error

	if !acme && g.Options.CertFile == "" && g.Options.KeyFile == "" {
		log.WithFields(log.Fields{"protocol": "http"}).Info("Mercure started")
		err = g.server.ListenAndServe()
	} else {
		// TLS
		if acme {
			certManager := &autocert.Manager{
				Prompt:     autocert.AcceptTOS,
				HostPolicy: autocert.HostWhitelist(g.Options.AcmeHosts...),
			}
			if g.Options.AcmeCertDir != "" {
				certManager.Cache = autocert.DirCache(g.Options.AcmeCertDir)
			}
			g.server.TLSConfig = certManager.TLSConfig()

			// Mandatory for Let's Encrypt http-01 challenge
			go http.ListenAndServe(":http", certManager.HTTPHandler(nil))
		}

		log.WithFields(log.Fields{"protocol": "https"}).Info("Mercure started")
		err = g.server.ListenAndServeTLS(g.Options.CertFile, g.Options.KeyFile)
	}

	if err != http.ErrServerClosed {
		log.Fatal(err)
	}

	<-idleConnsClosed
}

// chainHandlers configures and chains handlers
func (g *Gateway) chainHandlers() http.Handler {
	var useForwardedHeadersHandlers http.Handler
	if g.Options.UseForwardedHeaders {
		useForwardedHeadersHandlers = handlers.ProxyHeaders(g)
	} else {
		useForwardedHeadersHandlers = g
	}

	loggingHandler := handlers.CombinedLoggingHandler(os.Stderr, useForwardedHeadersHandlers)
	recoveryHandler := handlers.RecoveryHandler(
		handlers.RecoveryLogger(log.New()),
		handlers.PrintRecoveryStack(g.Options.Debug),
	)(loggingHandler)

	return recoveryHandler
}
