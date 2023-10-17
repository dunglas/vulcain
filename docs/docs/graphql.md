# Comparison with GraphQL and Other API Formats

Several API formats and architecture including [GraphQL](https://graphql.org/), [JSON:API](https://api-platform.com/docs/core/serialization/#embedding-relations) and [HAL](https://tools.ietf.org/html/draft-kelly-json-hal) try to solve the under-fetching problem by allowing the client to ask the server to embed related resources in the main HTTP response. Basically, the server will create a big JSON document containing the main requested resource and all its requested relations as nested objects. This solution allows to limit costly round-trip between the client and the server, and is (almost) the only efficient solution when using HTTP/1. However, [this hack introduces unnecessary complexity and hurts HTTP cache mechanisms](https://apisyouwonthate.com/blog/lets-stop-building-apis-around-a-network-hack). Also, it isn't necessary anymore with HTTP/2!

By using HTTP/2 Server Push, Vulcain fixes most problems caused by compound documents and sparse fieldsets based formats such as GraphQL and JSON:API:

* Because each pushed resource is sent in a separate HTTP/2 stream (HTTP/2 multiplexing), related resources can be sent in parallel to the client.
* While embedding resources is a forced push (the client receive the full JSON documents, even if it already has some parts of it), HTTP/2 Server Push allows the client [to cancel the push of resources it already has](cache.md), saving bandwidth and improving performance.
* Consequently, clients and network intermediates (such as [Varnish cache](cache.md)), can store each resource in a specific cache, while resource embedding only allows to have the full big JSON document in cache, [cache invalidation](https://en.wikipedia.org/wiki/Cache_invalidation) is then more efficient with Vulcain, and can be done at the HTTP level.

Specifically with GraphQL, using cache mechanisms provided by the HTTP protocol isn't easy (`POST` requests cannot be cached).

## Using GraphQL as Query Language for Vulcain

As stated in its name, GraphQL is foremost a convenient Query Language for APIs.
Guess what, GraphQL, the query language, is usable as-is with Vulcain servers!

The main idea is to write GraphQL queries client-side, which will be converted in REST requests containing Vulcain headers by a dedicated JavaScript library.

To do so, libraries such as [`apollo-link-rest`](https://www.apollographql.com/docs/link/links/rest/) can be used.
Thanks to `apollo-link-rest` you can write your request in GraphQL, use all [the tools of the frontend ecosystem relying on GraphQL](https://github.com/chentsulin/awesome-graphql#clients), but let the library send REST requests to the Vulcain server to fulfill the GraphQL query.

This approach also fixes [all the problems coming with using GraphQL server-side](https://dunglas.fr/2018/03/symfonylive-paris-slides-rest-vs-graphql-illustrated-examples-with-the-api-platform-framework/)!

Note: a higher-level library dedicated to Vulcain is being written.

## Type System and Introspection

Vulcain focuses on solving the under-fetching and the over-fetching problems. It's out of Vulcain's scope to provide a type system and an introspection mechanism.
However, Vulcain has been designed to play very well with existing formats providing these capabilities.

For hypermedia APIs, we strongly recommend to use [W3C's JSON-LD](https://json-ld.org/spec/latest/json-ld-api-best-practices/) along with [the Hydra Core Vocabulary](http://www.hydra-cg.com/). For less advanced non-hypermedia APIs, we recommend [OpenAPI](https://www.openapis.org/) (formerly known as Swagger).

## Subscription

Vulcain doesn't provide a way to push new versions of resources in real-time, however Vulcain plays very well with [Mercure](https://mercure.rocks) (created by the same author), a modern and RESTful replacement for WebSockets and GraphQL subscriptions.
