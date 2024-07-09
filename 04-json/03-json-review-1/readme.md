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

Q: What are the keys in the JSON representation of struct 1?

A: name, age, school, name, location

The keys in the JSON representation of struct 1 will be:
"name"
"age"
"school" (which will be an object with its own keys)
"name" (inside the "school" object)
"location" (inside the "school" object)
