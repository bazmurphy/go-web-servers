# The Profane

Not only do we validate that Chirps are under 140 characters, but we also have a list of words that are not allowed.

## Assignment

We need to update the `/api/validate_chirp` endpoint to replace all "profane" words with `4` asterisks: `****`.

Assuming the length validation passed, replace any of the following words in the Chirp with the static 4-character string `****`:

- kerfuffle
- sharbert
- fornax

Be sure to match against uppercase versions of the words as well, but not punctuation. "Sharbert!" does _not_ need to be replaced, we'll consider it a different word due to the exclamation point. Finally, instead of the `valid` boolean, your handler should return the cleaned version of the text in a JSON response:

### Example input

```json
{
  "body": "This is a kerfuffle opinion I need to share with the world"
}
```

### Example output

```json
{
  "cleaned_body": "This is a **** opinion I need to share with the world"
}
```

## Tips

_Use an HTTP client like [Thunder Client](https://marketplace.visualstudio.com/items?itemName=rangav.vscode-thunder-client) to test your POST requests._

I'd recommend creating two helper functions:

- `respondWithError(w http.ResponseWriter, code int, msg string)`
- `respondWithJSON(w http.ResponseWriter, code int, payload interface{})`

These helpers are not _required_ but might help [DRY](https://blog.boot.dev/clean-code/dry-code/) up your code when we add more endpoints in the future.

I'd also recommend breaking the bad word replacement into a separate function. You can even write some unit tests for it!

Here are some useful standard library functions:

- [strings.ToLower](https://pkg.go.dev/strings#ToLower)
- [strings.Split](https://pkg.go.dev/strings#Split)
- [strings.Join](https://pkg.go.dev/strings#Join)
