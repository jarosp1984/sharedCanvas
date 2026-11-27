# sharedCanvas

## Build, test, and run

The Go module is in `src/` (Go 1.21); run Go commands there. The server serves `../static/index.html` relative to its working directory, so start it from `src/`.

```bash
cd src
go build .
go test ./...
go test -run '^TestSessionLinesAreIsolatedBetweenSessions$' ./...
go vet ./...
go run .  # memory store, port 8080 by default
```

For Redis, from the repository root run `docker compose up -d redis`, then from `src/` run `LINE_STORE=redis REDIS_ADDR=localhost:6379 go run .`. `SC_PORT`, `REDIS_PASSWORD`, and `REDIS_DB` are optional; see `README.md` for defaults. Existing handler tests use an in-memory store and `httptest`, so they do not require Redis.

## Architecture and contracts

- `src/main.go` wires a `LineStore` into a standard-library HTTP mux. `/` generates a UUIDv4 and redirects to `/sessions/{uuid}`; the session page serves `static/index.html`. Sessions are identified by their URL rather than a separate creation record.
- The browser derives the session ID from the URL and talks to `/api/sessions/{uuid}/lines` (GET/POST) and `/api/sessions/{uuid}/lines/clear` (POST). `src/session.go` handles path parsing and UUIDv4 validation; handlers validate input and pass the session ID and request context to the store.
- `LineStore` has memory and Redis implementations. Keep session isolation and behavior consistent in both: line IDs start at 1 per session, `since` returns IDs greater than the cursor, and clearing removes lines without resetting the counter. Redis stores JSON lines and the counter under separate `sharedcanvas:sessions:{uuid}:lines` and `sharedcanvas:sessions:{uuid}:next_id` keys; its `since` implementation currently reads the full list and filters in Go.
- The single-page frontend draws immediately, posts segments, loads saved lines, and polls once per second using `maxRenderedLineID`. When an incremental poll is empty it checks the full list to detect a remote clear. Changes to IDs, clear behavior, or API response shapes must account for this client-side synchronization.
- Line validation in `src/line.go` matches the 800x600 canvas and 1-50 brush range in `static/index.html`. POST requests accept one JSON object with no unknown fields, limited to 1 MiB. Keep UI constraints, handler behavior, tests, and `doc/openapi/openapi.yaml` aligned when changing the API or line model.
