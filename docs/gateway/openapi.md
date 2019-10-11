# Mapping a Non-Hypermedia API using OpenAPI

While [using URLs as value for links](https://en.wikipedia.org/wiki/HATEOAS) should always be preferred, in most existing REST APIs links between resources are known by the server but are not provided in the HTTP response.

The gateway server can still be used for such APIs. To do so, an [OpenAPI specification](https://www.openapis.org/) (formerly known as Swagger) describing links between resources using [Link objects](http://spec.openapis.org/oas/v3.0.2#link-object) must be provided.

Imagine a web API having the following structure:

`/books/1`

```json
{
    "title": "1984",
    "author": 1
}
```

`/authors/1`

```json
{
    "givenName": "George",
    "familyName": "Orwell"
}
```

The link between books and authors, while not explicitly represented as an URL, can be documented in an OpenAPI v3 file:

```yaml
# openapi.yaml
openapi: 3.0.0
# ...
paths:
  '/books/{id}':
    get:
      # ...
      responses:
        default:
          links:
            author:
              operationId: getAuthor
              parameters:
                id: '$response.body#/author'
  '/authors/{id}':
    get:
      operationId: getAuthor
      responses:
        default:
      # ...
```

Then, use the `OPENAPI_FILE` environment variable to reference the OpenAPI file:

    UPSTREAM='http://your-api' OPENAPI_FILE='openapi.yaml' ADDR=':3000' KEY_FILE='tls/key.pem' CERT_FILE='tls/cert.pem' ./vulcain

In response to this request, both `/books/1` and `/authors/1` will be pushed by the Vulcain Gateway Server:

```http
GET /books/1 HTTP/2
Preload: /author
```

## Creating Links from Collections

The Vulcain Gateway Sever supports the [Extended JSON Pointer syntax](../../spec/vulcain.md#extended-json-pointer) to create links between elements of a collection and a related resource:

`/books`

```json
{
    "elements": [
        1,
        2,
        3
    ]
}
```

`/books/1`

```json
{
    "title": "1984",
    "author": 1
}
```

Use the following `links` object to link every item of the collection:

```yaml
# openapi.yaml
openapi: 3.0.0
# ...
paths:
  '/books/':
    get:
      # ...
      responses:
        default:
          links:
            author:
              operationId: getAuthor
              parameters:
                id: '$response.body#/elements/*'
  '/authors/{id}':
    get:
      operationId: getAuthor
      responses:
        default:
      # ...
```

With this HTTP request the server will push all resources linked from this collection. 

```http
GET /books HTTP/2
Preload: /elements/*
```

## Known Issues

* Only `operationId` can be used, `operationRef` is not supported yet, see [getkin/kin-openapi#130](https://github.com/getkin/kin-openapi/issues/130)
* `paths` ending with extensions aren't matched, see [getkin/kin-openapi#129](https://github.com/getkin/kin-openapi/issues/129)
