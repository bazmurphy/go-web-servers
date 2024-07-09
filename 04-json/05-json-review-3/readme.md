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

Q: What are the keys in the JSON representation of struct 3?

A: Name, Age

The JSON representation of struct 3 will use the following keys: "Name" "Age"
Note that in this case, the JSON keys will match the struct field names exactly, including the capitalization, because no JSON tags are specified.
