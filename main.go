package main

import (
	"os"

	fluentd "github.com/joonix/log"
	log "github.com/sirupsen/logrus"

	_ "net/http/pprof"

	_ "github.com/joho/godotenv/autoload"

	"github.com/dunglas/vulcain/gateway"
)

func init() {
	switch os.Getenv("LOG_FORMAT") {
	case "JSON":
		log.SetFormatter(&log.JSONFormatter{})
		return
	case "FLUENTD":
		log.SetFormatter(fluentd.NewFormatter())
	}
}

func main() {
	g, err := gateway.NewGatewayFromEnv()
	if err != nil {
		log.Panicln(err)
	}

	g.Serve()
}
