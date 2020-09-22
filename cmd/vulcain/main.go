package main

import (
	_ "github.com/joho/godotenv/autoload"

	"github.com/dunglas/vulcain"
)

func main() {
	s, err := vulcain.NewServerFromEnv() //nolint:staticcheck
	if err != nil {
		panic(err)
	}

	s.Serve() //nolint:staticcheck
}
