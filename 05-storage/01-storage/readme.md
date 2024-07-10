# Storage

Arguably the most important part of your typical web application is the storage of data. It would be pretty useless if each time you logged into your account on YouTube, Twitter or GitHub all of your subscriptions, tweets, or repositories were gone.

Let's talk about how web applications store, or "persist" data to a hard disk.

## Memory vs Disk

When you run a program on your computer, the program is loaded into _memory_. Memory is a lot like a scratch pad. It's fast, but it's not permanent. If the program terminates or restarts, the data in memory is _lost_.

When you're building a web server, any data you store in memory (in your program's variables) is lost when the server is restarted. Any important data needs to be saved to disk via the file system.

## Assignment

Our API needs to support the standard CRUD operations for "chirps". A "chirp" is just a short message that a user can post to the API, like a tweet. For now, we'll just be adding the `POST` and `GET` endpoints to create and read chirps respectively.

### POST /api/chirps

This endpoint should accept a JSON payload with a `body` field.

#### Request body:

```json
{ "body": "Hello, world!" }
```

Delete the `/api/validate_chirp` endpoint that we created before, but port all that logic into this one. Users should not be allowed to create invalid chirps!

If the chirp is valid, you should give it a _unique_ `id` and save it to disk. If all goes well, respond with a `201` status code and the full chirp resource. For now, just use integers for the `id` field, and increment the `id` by 1 for each new chirp.

#### Response body:

```json
{ "id": 1, "body": "Hello, world!" }
```

### GET /api/chirps

This endpoint should return an _array_ of _all_ chirps in the file, ordered by `id` in ascending order. Use a `200` code for success.

#### Response body:

```json
[
  { "id": 1, "body": "First chirp" },
  { "id": 2, "body": "The second chirp that was chirped" }
]
```

## Saving to disk

In a production system, you would almost certainly use a piece of database software like [PostgreSQL](https://www.postgresql.org/) or [MySQL](https://www.mysql.com/) to store your data on disk. We'll be learning SQL soon, but for now, we'll just be using the file system to store our data. It's less efficient, but it will work for our purposes.

Keep your entire "database" in a single file called `database.json` at the root of your project. Make sure to ignore it in Git. Your server should automatically create the file if it doesn't exist upon startup. Here's the structure of the file:

```json
{
  "chirps": {
    "1": { "id": 1, "body": "This is the first chirp ever!" },
    "2": { "id": 2, "body": "Hello, world!" }
  }
}
```

Any time you need to update the database, you should read the entire thing into memory (unmarshal it into a `struct`), update the data, and then write the entire thing back to disk (marshal it back into JSON).

To make sure that multiple requests don't try to write to the database at the same time, you should use a [mutex](https://blog.boot.dev/golang/golang-mutex/) to lock the database while you're using it. Again, I didn't say this will be _efficient_, but it will work!

While not necessary, I recommend encapsulating all of your database logic in an [internal](https://dave.cheney.net/2019/10/06/use-internal-packages-to-reduce-your-public-api-surface) `database` package.

## Tips

**Make sure to delete your `database.json` file every time before you run the tests!!!** The tests assume that they start with a fresh database each time.

Here are _some_ of the types and methods I used to create the database package to get you started:

```go
type DB struct { path string mux *sync.RWMutex }
```

```go
type DBStructure struct { Chirps map[int]Chirp `json:"chirps"` }
```

```go
// NewDB creates a new database connection // and creates the database file if it doesn't exist func NewDB(path string) (*DB, error)
```

```go
// CreateChirp creates a new chirp and saves it to disk func (db *DB) CreateChirp(body string) (Chirp, error)
```

```go
// GetChirps returns all chirps in the database func (db *DB) GetChirps() ([]Chirp, error)
```

```go
// ensureDB creates a new database file if it doesn't exist func (db *DB) ensureDB() error
```

```go
// loadDB reads the database file into memory func (db *DB) loadDB() (DBStructure, error)
```

```go
// writeDB writes the database file to disk func (db *DB) writeDB(dbStructure DBStructure) error
```

Here are some useful standard library functions to know about:

- [os.ReadFile](https://pkg.go.dev/os#ReadFile)
- [os.ErrNotExist](https://pkg.go.dev/os#ErrNotExist)
- [os.WriteFile](https://pkg.go.dev/os#WriteFile)
- [sort.Slice](https://pkg.go.dev/sort#Slice)

## Tests

1. POST /api/chirps
   Request Body: `{ "body": "I had something interesting for breakfast" }` 1. Expecting status code: `201` 2. Expecting JSON at `.id` to be equal to `1` 3. Expecting JSON at `.body` to be equal to `I had something interesting for breakfast`

2. GET /api/chirps

   1. Expecting status code: `200`
   2. Expecting JSON at `.[0].id` to be equal to `1`
   3. Expecting JSON at `.[0].body` to be equal to `I had something interesting for breakfast`

3. POST /api/chirps
   Request Body: `{"body": "What about second breakfast?"}`

   1. Expecting status code: `201`
   2. Expecting JSON at `.id` to be equal to `2`
   3. Expecting JSON at `.body` to be equal to `What about second breakfast?`

4. POST /api/chirps
   Request Body: `{"body": "Supper? Dinner? Do you know about those?"}`

   1. Expecting status code: 201
   2. Expecting JSON at `.id` to be equal to `3`
   3. Expecting JSON at `.body` to be equal to `Supper? Dinner? Do you know about those?`

5. GET /api/chirps
   1. Expecting status code: `200`
   2. Expecting body to contain: `I had something interesting for breakfast`
   3. Expecting body to contain: `What about second breakfast?`
   4. Expecting body to contain: `Supper? Dinner? Do you know about those?`
