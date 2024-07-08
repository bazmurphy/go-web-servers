# Serving Images

You may be wondering how the fileserver knew to serve the `index.html` file to the root of the server. It's _such_ a common convention on the web to use a file called `index.html` to serve the webpage for a given path, that the Go standard library's [FileServer](https://pkg.go.dev/net/http#FileServer) does it automatically.

When using a standard fileserver, the path to a file on disk is the same as its URL path. An exception is that `index.html` is served from `/` instead of `/index.html`.

## Try it out

Run your chirpy server again, and open `http://localhost:8080/index.html` in a new browser tab. You'll notice that you're redirected to `http://localhost:8080/`.

This works for all directories, not just the root!

For example:

- `/index.html` will be served from `/`
- `/pages/index.html` will be served from `/pages`
- `/pages/about/index.html` will be served from `/pages/about`

Alternatively, try opening a URL that doesn't exist, like `http://localhost:8080/doesntexist.html`. You'll see that the fileserver returns a 404 error.

## Assignment

Let's serve another type of file from our server: an image. Chirpy has a slick logo, and we need to serve it so that our users can load it in their browsers and mobile apps.

Download the Chirpy logo from below and add it to your project directory.

![Chirpy logo](logo.png)

Configure its filepath so that it's accessible from this URL:

```
http://localhost:8080/assets/logo.png
```

## Tests:

1. GET /assets/logo.png
   1. Expecting status code: 200
   2. Expecting "Content-Length" header to contain "35672"
