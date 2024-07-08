# Routing

The Go standard library has a lot of powerful HTTP features and, as of version 1.22, comes equipped with method-based pattern matching for routing.

Note that there are other powerful routing libraries like [Gorilla Mux](https://github.com/gorilla/mux) and [Chi](https://github.com/go-chi/chi), however, the instructions for this course will assume you are using Go's standard library. Just know that it isn't your only option!

In this lesson, we are going to limit which endpoints are available via which HTTP methods. In our current implementation, we can use any HTTP method to access any endpoint. _This is not ideal._

## Try it!

Run this command to send an empty `POST` request to your running server:

```bash
curl -X POST http://localhost:8080/healthz
```

You should get an `OK` response - but we want this endpoint to only be available via `GET` requests.

## Assignment

Add explicit HTTP methods to our current 2 custom endpoints to only allow for `GET` methods.

- `/healthz`
- `/metrics`

In general, a pattern looks something like this: `[METHOD ][HOST]/[PATH]`

Here are some examples:

```go
mux.HandleFunc("POST /articles", handlerArticlesCreate) mux.HandleFunc("DELETE /articles", handlerArticlesDelete)
```

When a request is made to one of these endpoints with a method other than `GET`, the server should return a `405` (Method Not Allowed) response (this is handled automatically!).
