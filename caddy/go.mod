module github.com/dunglas/vulcain/caddy

go 1.15

replace github.com/dunglas/vulcain => ../

require (
	github.com/caddyserver/caddy/v2 v2.3.0
	github.com/dunglas/vulcain v0.3.0
	go.uber.org/zap v1.16.0
)
