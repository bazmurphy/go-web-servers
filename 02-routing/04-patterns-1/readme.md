# Patterns

A pattern is a string that specifies the set of URL paths that should be matched to handle HTTP requests. Go's `ServeMux` router uses these patterns to dispatch requests to the appropriate handler functions based on the URL path of the request. As we saw in the previous lesson, patterns help organize the handling of different routes efficiently.

As previously mentioned, patterns generally look like this: `[METHOD ][HOST]/[PATH]`. Note that all three parts are optional.

## Rules and Definitions

### Fixed URL Paths

A pattern that exactly matches the URL path. For example, if you have a pattern `/about`, it will match the URL path `/about` and no other paths.

### Subtree Paths

If a pattern ends with a slash `/`, it matches all URL paths that have the same prefix. For example, a pattern `/images/` matches `/images/`, `/images/logo.png`, and `/images/css/style.css`. As we saw with our `/app/*` path, this is useful for serving a directory of static files or for structuring your application into sub-sections.

### Longest Match Wins

If more than one pattern matches a request path, the longest match is chosen. This allows more specific handlers to override more general ones. For example, if you have patterns `/` (root) and `/images/`, and the request path is `/images/logo.png`, the `/images/` handler will be used because it's the longest match.

### Host-specific Patterns

We won't be using this but be aware that patterns can also start with a hostname (e.g., `www.example.com/`). This allows you to serve different content based on the Host header of the request. If both host-specific and non-host-specific patterns match, the host-specific pattern takes precedence.

If you're interested, you can read more in the [ServeMux docs](https://pkg.go.dev/net/http#ServeMux).

# Quiz

Q: Patterns require an explicit HTTP method

A: False
