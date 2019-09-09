# Vulcain: the RESTful Way to Fetch All You Need... But Only What You Need
*Protocol and Reference Gateway*

[![GoDoc](https://godoc.org/github.com/dunglas/vulcain?status.svg)](https://godoc.org/github.com/dunglas/vulcain/hub)
[![Build Status](https://travis-ci.com/dunglas/vulcain.svg?branch=master)](https://travis-ci.com/dunglas/vulcain)
[![Coverage Status](https://coveralls.io/repos/github/dunglas/vulcain/badge.svg?branch=master)](https://coveralls.io/github/dunglas/vulcain?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/dunglas/vulcain)](https://goreportcard.com/report/github.com/dunglas/vulcain)

## Example Queries

Preload all relations referenced in the `hydra:member`, then in the author relationship, but only include the title and the author of these relations:

```bash
curl https://localhost:3000/books.jsonld \
    --get \
    --data 'preload=/hydra:member/*/author' \
    --data 'fields=/hydra:member/*/author' \
    --data 'fields=/hydra:member/*/title' \
    --verbose \
    --insecure 
```

Using headers:

```bash
curl https://localhost:3000/books.jsonld \
    --get \
    --header 'Preload: /hydra:member/*/author' \
    --header 'Fields: /hydra:member/*/author' \
    --header 'Fields: /hydra:member/*/title' \
    --verbose \
    --insecure 
```

## Prior Art

* JSON:API
* GraphQL
* Hades
* Prefer-Push I-D
* Transclude I-D

## Credits

Created by [Kévin Dunglas](https://dunglas.fr). Sponsored by [Les-Tilleuls.coop](https://les-tilleuls.coop).

Ideas and code used in Vulcain's reference implementation have been taken from [Hades](https://github.com/gabesullice/hades), an HTTP/2 reverse proxy for JSON:API backend.
