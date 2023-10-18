FROM caddy:2-alpine

COPY vulcain /usr/bin/caddy
COPY Caddyfile /etc/caddy/Caddyfile
