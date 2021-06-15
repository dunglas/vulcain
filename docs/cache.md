# Cache Considerations

One of the advantages of the approach proposed by Vulcain is that it allows to leverage the full range of capabilities provided by the HTTP protocol, especially [its caching mechanisms](https://developer.mozilla.org/en-US/docs/Web/HTTP/Caching) and [its layered architecture](https://www.ics.uci.edu/~fielding/pubs/dissertation/rest_arch_style.htm#sec_5_1_6).

## Normalizing the API to Maximise Hits

To unleash the full power of Vulcain, API resources should be as granular as possible. Use links as most as possible. For instance, collections should only be collections of links, and should not embed the relations themselves. Doing this allows to have a specific cache for every resources. For instance, when fetching a collection, the client will not re-fetch elements of this collections already in cache.

Example:

* The user arrive on a detail page, `GET /book/1`, stored in cache
* The user click on the links to the collection, `GET /books/`, `/books` contains links to `/books/1` and `/books/2`, only `/books/2` is fetched because `/books/1` is already in cache

## Using Vulcain with Cache Servers such as Varnish

Vulcain plays very well with cache reverse proxies such as [Varnish](https://varnish-cache.org/).
A common setup is to put the Vulcain Gateway Server in at the edge of the network to be able to use Server Push, with a Varnish server behind it, and finally the application server behind Varnish.

The full documents (not filtered) should be stored by Varnish. To do so, Varnish should strip the `Fields` and `Preload` headers before forwarding the request to the backend server.
This way, in most cases, the resources will be fetched directly from the Varnish cache (hit) and will be served almost instantly quickly to the Vulcain Gateway Server. The Vulcain server will then recursively fetch requested relations, and will filter the documents by removing useless fields.

Vulcain also plays very well with cache invalidation mechanisms such as [xkey](https://github.com/varnish/varnish-modules/blob/master/docs/vmod_xkey.rst).

## Preventing to Push Resources Already in Cache (Cache-Digests and CASPer)

When servers push resources that are already in the client cache, [the client will cancel this push](https://tools.ietf.org/html/rfc7540#section-8.2.2).
However, it's even better to not send a push promise at all for resources already in cache. To do so, the new [Cache-Digests for HTTP/2 RFC](https://httpwg.org/http-extensions/cache-digest.html) can be used. Alternatively, [CASPer](https://h2o.examp1e.net/configure/http2_directives.html#http2-casper) (cookie-based cache aware Server Push) can be used.

Note: the Gateway Server [doesn't support Cache Digests nor CASPer yet](https://github.com/dunglas/vulcain/issues/1).
