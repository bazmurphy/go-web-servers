# Stateful Handlers

It's frequently useful to have a way to store and access state in our handlers. For example, we might want to keep track of the number of requests we've received, or we may want to pass around an open connection to a database, or credentials to an API.

## Assignment

The product managers at Chirpy want to know how many requests are being made to serve our homepage - in essence, they want to know how many people are viewing the site!

They have asked for a simple HTTP endpoint they can hit to get the number of requests that have been processed. It will return the count as plain text in the response body.

For now, they just want the number of requests that have been processed since the last time the server was started, we don't need to worry about saving the data between restarts.

## Steps

1.  Create a struct that will hold any stateful, in-memory data we'll need to keep track of. In our case, we just need to keep track of the number of requests we've received.

```go
type apiConfig struct { fileserverHits int }
```

2.  Next, write a new [middleware](https://en.wikipedia.org/wiki/Middleware) method on a `*apiConfig` that increments the `fileserverHits` counter every time it's called. Here's the method signature I used:

```go
func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler { // ... }
```

3.  [Wrap](https://en.wikipedia.org/wiki/Wrapper_function) the `http.FileServer` handler with the middleware method we just wrote. For example:

```go
mux.Handle("/app/*", apiCfg.middlewareMetricsInc(handler))
```

4.  Create a new handler that writes the number of requests that have been counted as plain text in this format to the HTTP response:

Where `x` is the number of requests that have been processed. This handler should be a method on the `*apiConfig` struct so that it can access the `fileserverHits` data.

5.  Register that handler with the serve mux on the `/metrics` path.
6.  Finally, create and register a handler on the `/reset` path that, when hit, will reset your `fileserverHits` back to `0`.

_It should follow the same design as the previous handlers._

Remember, similar to the metrics endpoint, `/reset` will need to be a method on the `*apiConfig` struct so that it can also access the `fileserverHits`

## Tests:

1. GET /reset
   1. Expecting status code: 200
2. GET /app

   1. Headers:`{"Cache-Control": "no-cache"}`
   1. Expecting status code: 200
   1. Expecting body to contain: Welcome to Chirpy

3. GET /metrics
   1. Expecting status code: 200
   1. Expecting body to contain: Hits: 1
4. GET /app

   1. Headers: `{ "Cache-Control": "no-cache" }`
   1. Expecting status code: 200
   1. Expecting body to contain: Welcome to Chirpy

5. GET /app

   1. Headers: `{"Cache-Control": "no-cache"}`
   1. Expecting status code: 200
   1. Expecting body to contain: Welcome to Chirpy

6. GET /app

   1. Headers:`{"Cache-Control": "no-cache"}`
   1. Expecting status code: 200
   1. Expecting body to contain: Welcome to Chirpy

7. GET /metrics
   1. Expecting status code: 200
   1. Expecting body to contain: Hits: 4
