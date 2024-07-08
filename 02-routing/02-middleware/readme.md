# Middleware

Middleware is a way to wrap a handler with additional functionality. It is a common pattern in web applications that allows us to write DRY code. For example, we can write a middleware that logs every request to the server. We can then wrap our handler with this middleware and every request will be logged without us having to write the logging code in every handler.

Here are examples of the middleware that we've written so far.

## Keeping track of the number of times a handler has been called

```go
func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler { return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { cfg.fileserverHits++ next.ServeHTTP(w, r) }) }
```

## Logging every request

We haven't written this one yet, but it would look something like this:

```go
func middlewareLog(next http.Handler) http.Handler { return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { log.Printf("%s %s", r.Method, r.URL.Path) next.ServeHTTP(w, r) }) }
```

# Quiz

Q: Middleware...

A: Is a way to inject code into HTTP handlers
