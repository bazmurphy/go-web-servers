# JSON Review

## Struct 1

```go
type parameters struct {
    Name string `json:"name"`
    Age int `json:"age"`
    School struct {
        Name string `json:"name"`
        Location string `json:"location"`
    } `json:"school"`
}
```

## Struct 2

```go
type parameters struct {
    name string `json:"name"`
    Age int `json:"age"`
}
```

## Struct 3

```go
type parameters struct {
    Name string
    Age int
}
```

# Quiz

Q: Which keys will be parsed in the JSON representation of struct 2?

A: age

In the JSON representation of struct 2, only one key will be parsed: "age"
The "name" field will not be parsed because it starts with a lowercase letter, making it an unexported field in Go. Unexported fields are not encoded to or decoded from JSON.
