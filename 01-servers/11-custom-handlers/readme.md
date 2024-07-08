# Custom Handlers

In the previous exercise, we used the [http.FileServer](https://pkg.go.dev/net/http#FileServer) function, which simply returns a built-in [http.Handler](https://pkg.go.dev/net/http#Handler).

An `http.Handler` is just an interface:

```go
type Handler interface { ServeHTTP(ResponseWriter, *Request) }
```

Any type with a `ServeHTTP` method that matches the [http.HandlerFunc](https://pkg.go.dev/net/http#HandlerFunc) signature above is an `http.Handler`. Take a second to think about it: it makes a lot of sense! To handle an incoming HTTP request, all a function needs is a way to write a response and the request itself.

## Assignment

Let's add a readiness endpoint to the Chirpy server! Readiness endpoints are commonly used by external systems to check if our server is ready to receive traffic.

The endpoint should be accessible at the `/healthz` path using any HTTP method.

The endpoint should simply return a `200 OK` status code indicating that it has started up successfully and is listening for traffic. The endpoint should return a `Content-Type: text/plain; charset=utf-8` header, and the body will contain a message that simply says "OK" (the text associated with the 200 status code).

_Later this endpoint can be enhanced to return a `503 Service Unavailable` status code if the server is not ready._

### 1\. Add the readiness endpoint

I recommend using the [mux.HandleFunc](https://pkg.go.dev/net/http#ServeMux.HandleFunc) to register your handler. Your handler can just be a function that matches the signature of [http.HandlerFunc](https://pkg.go.dev/net/http#HandlerFunc):

```go
handler func(http.ResponseWriter, *http.Request)
```

Your handler should do the following:

1.  Write the `Content-Type` header
2.  Write the status code using [w.WriteHeader](https://pkg.go.dev/net/http#ResponseWriter.WriteHeader)
3.  Write the body text using [w.Write](https://pkg.go.dev/net/http#ResponseWriter.Write)

### 2\. Update the fileserver path

Now that we've added a new handler, we don't want potential conflicts with the fileserver handler. Update the fileserver to use the `/app/*` path instead of `/`.

Not only will you need to [http.Handle](https://pkg.go.dev/net/http#Handle) the `/app/*` path, you'll also need to strip the `/app` prefix from the request path before passing it to the fileserver handler. You can do this using the [http.StripPrefix](https://pkg.go.dev/net/http#StripPrefix) function.

## Tests:

1. GET /app
   1. Expecting status code: 200
   2. Expecting body to contain: <html
   3. Expecting body to contain: Welcome to Chirpy
   4. Expecting body to contain: </html>
2. GET /app/assets
   1. Expecting status code: 200
   2. Expecting body to contain: <pre>
   3. Expecting body to contain: <a href="logo.png">logo.png</a>
   4. Expecting body to contain: </pre>
3. GET /healthz
   1. Expecting status code: 200
   2. Expecting "Content-Type" header to contain "text/plain; charset=utf-8"
   3. Expecting body to contain: OK
