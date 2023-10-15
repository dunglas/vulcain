package main

import (
	"github.com/caddyserver/caddy/v2"
	caddycmd "github.com/caddyserver/caddy/v2/cmd"

	"go.uber.org/automaxprocs/maxprocs"

	// plug in Caddy modules here.
	_ "github.com/caddyserver/caddy/v2/modules/standard"
	_ "github.com/dunglas/vulcain/caddy"
)

//nolint:gochecknoinits
func init() {
	//nolint:errcheck
	maxprocs.Set(maxprocs.Logger(caddy.Log().Sugar().Debugf))
}

func main() {
	caddycmd.Main()
}
