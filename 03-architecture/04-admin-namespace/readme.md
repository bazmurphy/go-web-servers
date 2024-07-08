# Admin Namespace

One of the advantages of a monolithic architecture is that it's fairly simple to inject data directly into the HTML of a web page.

## Assignment

Let's swap out the `/metrics` endpoint, which just returns plain text, for an `/admin/metrics` (not under the `/api` namespace) endpoint that returns HTML intended to be rendered in the browser. It should be accessible via `GET` requests only.

Use this template:

```html
<html>
  <body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited %d times!</p>
  </body>
</html>
```

Where `%d` is replaced with the number of times the page has been loaded.

Make sure you use the `Content-Type` header to set the response type to `text/html` so that the browser knows how to render it.

Try loading `http://localhost:8080/admin/metrics` in your browser, and in another tab load `http://localhost:8080/app` a few times. Refreshing the admin page should show the updated count.

Run and submit the HTTP tests using the CLI tool.

## Tests:

1. GET /api/reset
   1. Expecting status code: 200
2. GET /admin/metrics
   1. Expecting status code: 200
   1. Expecting body to contain: Welcome, Chirpy Admin
   1. Expecting body to contain: Chirpy has been visited 0 times!
3. GET /app
   1. Expecting status code: 200
   1. Expecting body to contain: Welcome to Chirpy
4. GET /admin/metrics
   1. Expecting status code: 200
   1. Expecting body to contain: Chirpy has been visited 1 times!
5. GET /app
   1. Expecting status code: 200
   1. Expecting body to contain: Welcome to Chirpy
6. GET /app
   1. Expecting status code: 200
   1. Expecting body to contain: Welcome to Chirpy
7. GET /app
   1. Expecting status code: 200
   1. Expecting body to contain: Welcome to Chirpy
8. GET /admin/metrics
   1. Expecting status code: 200
   1. Expecting body to contain: Chirpy has been visited 4 times!
