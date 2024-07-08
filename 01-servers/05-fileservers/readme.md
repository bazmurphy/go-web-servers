# Fileservers

A _fileserver_ is a kind of simple web server that serves static files from the host machine. Fileservers are often used to serve static assets for a website, things like:

- HTML
- CSS
- JavaScript
- Images

## Assignment

The Go standard library makes it super easy to build a simple fileserver. Build and run a fileserver that serves a file called `index.html` from its root at `http://localhost:8080`. That file should contain this HTML:

```html
<html>
  <body>
    <h1>Welcome to Chirpy</h1>
  </body>
</html>
```

## Steps

1.  Add the HTML code above to a file called `index.html` in the same root directory as your server
2.  Use the [http.NewServeMux](https://pkg.go.dev/net/http#NewServeMux)'s [.Handle()](https://pkg.go.dev/net/http#ServeMux.Handle) method to add a handler for the root path (`/`).
3.  Use a standard [http.FileServer](https://pkg.go.dev/net/http#FileServer) as the handler
4.  Use [http.Dir](https://pkg.go.dev/net/http#Dir) to convert a filepath (in our case a dot: `.` which indicates the current directory) to a directory for the `http.FileServer`.
5.  Re-build and run your server
6.  Test your server by visiting `http://localhost:8080` in your browser
7.  Run the tests using the CLI

## Tests:

1. GET /

- 1. Expecting status code: 200
- 2. Expecting body to contain: Welcome to Chirpy
