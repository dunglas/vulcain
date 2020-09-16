package main

import (
	"os"

	fluentd "github.com/joonix/log"
	log "github.com/sirupsen/logrus"

	_ "github.com/joho/godotenv/autoload"

	"github.com/dunglas/vulcain"
)

func init() {
	switch os.Getenv("LOG_FORMAT") {
	case "JSON":
		log.SetFormatter(&log.JSONFormatter{})
		return
	case "FLUENTD":
		log.SetFormatter(fluentd.NewFormatter())
	}

	if os.Getenv("DEBUG") == "1" {
		log.SetLevel(log.DebugLevel)
	}
}

func main() {
	s, err := vulcain.NewServerFromEnv()
	if err != nil {
		log.Panicln(err)
	}

	s.Serve()
}
