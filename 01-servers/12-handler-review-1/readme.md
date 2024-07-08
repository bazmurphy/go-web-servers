# Handler Review

# Handler

An [http.Handler](https://pkg.go.dev/net/http#Handler) is any [defined type](https://go.dev/ref/spec#Type_definitions) that implements the set of methods defined by the `Handler` [interface](https://go.dev/tour/methods/9), specifically the `ServeHTTP` method.

```go
type Handler interface { ServeHTTP(ResponseWriter, *Request) }
```

The [ServeMux](https://pkg.go.dev/net/http#ServeMux) you used in the previous exercise is an `http.Handler`.

You will typically use a `Handler` for more complex use cases, such as when you want to implement a custom router, middleware, or other custom logic.

## HandlerFunc

```go
type HandlerFunc func(ResponseWriter, *Request)
```

You'll typically use a `HandlerFunc` when you want to implement a simple handler. The `HandlerFunc` type is just a function that matches the `ServeHTTP` signature above.

## Why this signature?

The `Request` argument is fairly obvious: it contains all the information about the incoming request, such as the HTTP method, path, headers, and body.

The `ResponseWriter` is less intuitive in my opinion. The response is an _argument_, not a _return type_. Instead of returning a value all at once from the handler function, we _write_ the response to the `ResponseWriter`.

# Quiz

Q: Which underlying type might implement the Handler interface?

A: Any of these
