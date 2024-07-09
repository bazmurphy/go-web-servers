# JSON

Hopefully, by now you already know what JSON is. If not, you should go back and take the Learn HTTP course [here first](https://boot.dev/learn/learn-http).

What you may be new to is handling and parsing JSON on the server side, rather than sending it as a client.

If you want to take a super deep dive into JSON in Go, then you can [read this post here](https://blog.boot.dev/golang/json-golang/). With that in mind, you don't need to! I'll give you the relevant info below.

## Decode JSON request body

It's _very_ common for `POST` requests to send JSON data in the request body. Here's how you can handle that incoming data:

```json
{
  "name": "John",
  "age": 30
}
```

```go
func handler(w http.ResponseWriter, r *http.Request){
    type parameters struct {
        // these tags indicate how the keys in the JSON should be mapped to the struct fields
        // the struct fields must be exported (start with a capital letter) if you want them parsed
        Name string `json:"name"`
        Age int `json:"age"`
    }

    decoder := json.NewDecoder(r.Body)
    params := parameters{}
    err := decoder.Decode(&params)
    if err != nil {
        // an error will be thrown if the JSON is invalid or has the wrong types
        // any missing fields will simply have their values in the struct set to their zero value
		log.Printf("Error decoding parameters: %s", err)
		w.WriteHeader(500)
		return
    }
    // params is a struct with data populated successfully
    // ...
}
```

## Encode JSON response body

```go
func handler(w http.ResponseWriter, r *http.Request){
    // ...

    type returnVals struct {
        // the key will be the name of struct field unless you give it an explicit JSON tag
        CreatedAt time.Time `json:"created_at"`
        ID int `json:"id"`
    }
    respBody := returnVals{
        CreatedAt: time.Now(),
        ID: 123,
    }
    dat, err := json.Marshal(respBody)
    if err != nil {
      log.Printf("Error marshalling JSON: %s", err)
      w.WriteHeader(500)
      return
    }
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(200)
	  w.Write(dat)
}
```

## Assignment

At Chirpy, we have a silly rule that says all Chirps must be 140 characters long or less.

Add a new endpoint to the Chirpy API that accepts a `POST` request at `/api/validate_chirp`. It should expect a JSON body of this shape:

```json
{
  "body": "This is an opinion I need to share with the world"
}
```

If any errors occur, it should respond with an appropriate HTTP status code and a JSON body of this shape:

```json
{
  "error": "Something went wrong"
}
```

For example, if the Chirp is too long, respond with a `400` code and this body:

```json
{
  "error": "Chirp is too long"
}
```

If the Chirp is valid, respond with a `200` code and this body:

```json
{
  "valid": true
}
```

## Tips

Use an HTTP client like [Thunder Client](https://marketplace.visualstudio.com/items?itemName=rangav.vscode-thunder-client) to test your POST requests.

Use [json.Marshal()](https://pkg.go.dev/encoding/json#Marshal) like the example above to remove whitespace in your encoded data.

```

```
