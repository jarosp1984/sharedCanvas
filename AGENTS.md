# AGENTS.md

## Project Structure
### Directories

- `src/`: Go backend application code, HTTP handlers, domain model, and store implementations.
- `static/`: Frontend static assets served by the Go server.
- `doc/openapi/`: OpenAPI spec for the HTTP API.
- `doc/prompts_history.md`: Historical prompt notes used during development.

### Key Files

- `src/main.go`: Server bootstrap and route registration.
- `src/handler.go`: API handlers for line operations.
- `src/store.go`: Store interface and environment-based store selection.
- `src/store_memory.go`: In-memory line store implementation.
- `src/store_redis.go`: Redis-backed line store implementation.
- `src/handler_test.go`: Backend HTTP handler tests.
- `static/index.html`: Canvas UI and browser-side API integration.
- `docker-compose.yml`: Redis service definition for local development.
