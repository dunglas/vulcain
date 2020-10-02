// Package caddy provides a handler for Caddy Server (https://caddyserver.com/)
// allowing to turn any web API in a one supporting the Vulcain protocol.
package caddy

import (
	"bytes"
	"net/http"
	"strconv"
	"sync"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
	"github.com/caddyserver/caddy/v2/caddyconfig/httpcaddyfile"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
	"github.com/dunglas/vulcain"
	"go.uber.org/zap"
)

func init() {
	caddy.RegisterModule(Vulcain{})
	httpcaddyfile.RegisterHandlerDirective("vulcain", parseCaddyfile)
}

var bufPool = sync.Pool{
	New: func() interface{} {
		return new(bytes.Buffer)
	},
}

type Vulcain struct {
	OpenAPIFile string `json:"openapi_file,omitempty"`
	MaxPushes   int    `json:"max_pushes,omitempty"`

	vulcain *vulcain.Vulcain
	logger  *zap.Logger
}

// CaddyModule returns the Caddy module information.
func (Vulcain) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "http.handlers.vulcain",
		New: func() caddy.Module { return new(Vulcain) },
	}
}

func (v *Vulcain) Provision(ctx caddy.Context) error {
	if v.MaxPushes == 0 {
		v.MaxPushes = -1
	}

	v.logger = ctx.Logger(v)

	v.vulcain = vulcain.New(
		vulcain.WithOpenAPIFile(v.OpenAPIFile),
		vulcain.WithMaxPushes(v.MaxPushes),
		vulcain.WithLogger(ctx.Logger(v)),
	)

	return nil
}

// ServeHTTP applies Vulcain directives.
func (v Vulcain) ServeHTTP(w http.ResponseWriter, r *http.Request, next caddyhttp.Handler) error {
	r = r.WithContext(v.vulcain.CreateRequestContext(w, r))

	var wait bool
	defer func() { v.vulcain.Finish(r, wait) }()

	if !v.vulcain.IsValidRequest(r) {
		return next.ServeHTTP(w, r)
	}

	buf := bufPool.Get().(*bytes.Buffer)
	buf.Reset()
	defer bufPool.Put(buf)

	rec := caddyhttp.NewResponseRecorder(w, buf, func(status int, header http.Header) bool {
		return v.vulcain.IsValidResponse(r, status, header)
	})
	if err := next.ServeHTTP(rec, r); err != nil {
		v.vulcain.Finish(r, false)
		return err
	}
	if !rec.Buffered() {
		return nil
	}

	b, err := v.vulcain.Apply(r, w, rec.Buffer(), rec.Header())
	if err != nil {
		return rec.WriteResponse()
	}

	w.WriteHeader(rec.Status())
	_, err = w.Write(b)
	if err != nil {
		return err
	}

	wait = true

	return nil
}

// UnmarshalCaddyfile sets up the handler from Caddyfile tokens. Syntax:
//
//     vulcain {
//         # path to the OpenAPI file describing the relations (for non-hypermedia APIs)
//	       openapi_file <path>
//         # Maximum number of pushes to do (-1 for unlimited)
//         max_pushes -1
//     }
func (v *Vulcain) UnmarshalCaddyfile(d *caddyfile.Dispenser) error {
	for d.Next() {
		for d.NextBlock(0) {
			switch d.Val() {
			case "openapi_file":
				if !d.NextArg() {
					return d.ArgErr()
				}

				v.OpenAPIFile = d.Val()

			case "max_pushes":
				if !d.NextArg() {
					return d.ArgErr()
				}

				maxPushes, err := strconv.Atoi(d.Val())
				if err != nil {
					return d.Errf("bad max_pushes value '%s': %v", d.Val(), err)
				}

				v.MaxPushes = maxPushes
			}
		}
	}
	return nil
}

// parseCaddyfile unmarshals tokens from h into a new Middleware.
func parseCaddyfile(h httpcaddyfile.Helper) (caddyhttp.MiddlewareHandler, error) {
	var m Vulcain
	err := m.UnmarshalCaddyfile(h.Dispenser)
	return m, err
}

// Interface guards
var (
	_ caddy.Provisioner           = (*Vulcain)(nil)
	_ caddyhttp.MiddlewareHandler = (*Vulcain)(nil)
	_ caddyfile.Unmarshaler       = (*Vulcain)(nil)
)
