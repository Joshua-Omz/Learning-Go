# Learning Go 🐹

A personal repository documenting my journey of learning the Go programming language. Each project here represents a concept or milestone I tackled along the way.

---

## Projects

### 1. Persistent Scribes

> **Concepts practiced:** HTTP servers, JSON encoding/decoding, file I/O, middleware, structs, maps, and making outbound HTTP requests.

A simple REST API for managing personal notes. Notes are stored in a local `data.json` file so they survive server restarts. The server also has the ability to "enhance" any note by attaching a random inspirational quote fetched from an external API.

**Location:** [`persistent-scribes/`](./persistent-scribes/)

#### What I Learned Building This

- Setting up an HTTP server using only Go's standard library (`net/http`)
- Defining custom structs and using JSON struct tags
- Reading and writing JSON to disk with `os.ReadFile` / `os.WriteFile`
- Using an in-memory `map` as a simple database
- Writing middleware (API key guard) using function wrapping
- Making outbound HTTP GET requests and decoding the JSON response
- Composing multiple data types into a single response struct

#### Running the Server

```bash
cd persistent-scribes
go run main.go
```

The server starts on **port 8080**.

#### Authentication

All endpoints require an `X-API-Key` header. Any non-empty value is accepted.

```
X-API-Key: my-secret-key
```

#### API Endpoints

| Method   | Endpoint                   | Description                              |
|----------|----------------------------|------------------------------------------|
| `POST`   | `/notes/create`            | Create a new note                        |
| `GET`    | `/notes/read?id=<id>`      | Read a note by ID                        |
| `DELETE` | `/notes/delete?id=<id>`    | Delete a note by ID                      |
| `GET`    | `/notes/enhance?id=<id>`   | Read a note enriched with a random quote |

##### Create a Note

```bash
curl -X POST http://localhost:8080/notes/create \
  -H "X-API-Key: my-secret-key" \
  -H "Content-Type: application/json" \
  -d '{"id": "1", "content": "Learning Go is awesome"}'
```

##### Read a Note

```bash
curl http://localhost:8080/notes/read?id=1 \
  -H "X-API-Key: my-secret-key"
```

##### Delete a Note

```bash
curl -X DELETE http://localhost:8080/notes/delete?id=1 \
  -H "X-API-Key: my-secret-key"
```

##### Enhance a Note

Fetches the note and attaches a random quote from [dummyjson.com](https://dummyjson.com/quotes/random).

```bash
curl http://localhost:8080/notes/enhance?id=1 \
  -H "X-API-Key: my-secret-key"
```

**Example response:**

```json
{
  "id": "1",
  "content": "Learning Go is awesome",
  "qoute": {
    "quote": "The secret of getting ahead is getting started.",
    "author": "Mark Twain"
  }
}
```

---

## Goals

- [x] Build a working HTTP server from scratch using the standard library
- [x] Practice JSON serialization / deserialization
- [x] Implement file-based persistence
- [x] Write HTTP middleware in Go
- [x] Make outbound HTTP requests and handle the response
- [ ] Add a database (e.g., SQLite or PostgreSQL)
- [ ] Write unit tests
- [ ] Explore Go modules and dependency management with third-party packages
- [ ] Learn concurrency with goroutines and channels

---

## Resources

- [The Go Programming Language Tour](https://go.dev/tour/)
- [Go by Example](https://gobyexample.com/)
- [Official Go Documentation](https://pkg.go.dev/)
