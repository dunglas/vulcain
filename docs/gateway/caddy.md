# Vulcain for Caddy

The Vulcain Gateway Server can be put in front of **any existing REST API** to transform it in a Vulcain-enabled API.
It works with hypermedia APIs ([JSON-LD](https://json-ld.org), [JSON:API](https://jsonapi.org/), [HAL](https://tools.ietf.org/html/draft-kelly-json-hal), [Siren](https://github.com/kevinswiber/siren) ...) as well as [with other non-hypermedia APIs by configuring the server using a subset of the OpenAPI specification](openapi.md).

Tip: the easiest way to create a hypermedia API is to use [the API Platform framework](https://api-platform.com) (by the same author than Vulcain).

The Gateway Server is a module for [the Caddy server](https://caddyserver.com): a powerful, enterprise-ready, open source web server with automatic HTTPS written in Go.

The Vulcain module for Caddy allows to turn any existing web API in a one supporting all features of Vulcain in a few minutes.

## Install

### Docker

The easiest way to get started is to use Docker:

```console
docker run -e VULCAIN_UPSTREAM='http://your-api' -p 80:80 -p 443:443 dunglas/vulcain 
```

The configuration file is located at `/etc/caddy/Caddyfile`.

### Binaries

1. Go on [the Caddy server download page](https://caddyserver.com/download)
2. Select the `github.com/dunglas/vulcain/caddy` module
3. Select other modules you're interested in such as [the cache module](https://github.com/caddyserver/cache-handler) or [Mercure](https://mercure.rocks) (optional)
4. Download and enjoy!

Alternatively, you can use [xcaddy](https://github.com/caddyserver/xcaddy) to create a custom build of Caddy containing the Vulcain module: `xcaddy build --with github.com/dunglas/vulcain/caddy`

Pre-built binaries are also available for download [on the releases page](https://github.com/dunglas/vulcain/releases).

## Configuration

Example configuration:

```caddyfile
{
    order vulcain before request_header
}

example.com {
    vulcain {
        openapi_file my-openapi-description.yaml # optional
        max_pushes 100 # optional
        early_hints # optional, usually not necessary
    }
    reverse_proxy my-api:8080 # all other handlers such as the static file server and custom handlers are also supported
}
```

All other [Caddyfile directives](https://caddyserver.com/docs/caddyfile) are also supported.

## Start the Server

Just run `./caddy run`.

## Cache Handler

Vulcain is best used with an HTTP cache server. The Caddy and the Vulcain team maintain together a [distributed HTTP cache module](https://github.com/caddyserver/cache-handler) built on top of [Souin](https://github.com/darkweak/souin) supporting most of the RFC.

## 103 Early Hints

The gateway server can trigger 103 "Early Hints" responses including Preload hints automatically.
However, enabling this feature is usually useless because the gateway server doesn't supports JSON streaming (yet).

Consequently the server will have to wait for the full JSON response to be received from upstream before being able to compute the Link headers to send.

When the full response is available, we can send the final response directly.

For best performance, better send Early Hints responses as soon as possible, directly from the upstream application.

The gateway server will automatically and instantly forward all 103 responses coming from upstream, even if the `early_hints` directive is not set.
