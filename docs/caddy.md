# Vulcain for Caddy

[Caddy](https://caddyserver.com/) is a powerful, enterprise-ready, open source web server with automatic HTTPS written in Go.

The Vulcain module for Caddy 2 allows to turn any existing web API in a one supporting all features of Vulcain in a few minutes.

## Install

Use [xcaddy](https://github.com/caddyserver/xcaddy) to create a build of Caddy containing the Vulcain module.

1. Install xcaddy
2. Run: `xcaddy build --with github.com/dunglas/vulcain/caddy`

## Configuration

Example configuration:

```caddyfile
{
    order vulcain before request_header
	experimental_http3 # optional, enables HTTP/3
}

my-site.com

reverse_proxy my-api:8080 # all other handlers such as the static file server and custom handlers are also supported
vulcain {
    openapi_file my-openapi-description.yaml # optional
    max_pushes 100 # optional
}
```

## Start the Server

Just run `./caddy run`.

## Cache Handler

Vulcain is best used with an HTTP cache server. The Caddy and the Vulcain team maintain together a distributed HTTP cache module supporting most of the RFC. To build Caddy with this module and Vulcain run: `xcaddy build --with github.com/dunglas/vulcain/caddy --with github.com/caddyserver/cache-handler`
