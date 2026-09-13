# sharedCanvas

Canvas for shared drawings. This is a test project for exploring agentic development with GitHub Copilot.

Each canvas session has its own UUID in the URL so you can share a link with someone else and collaborate in the same session.

## Run

Memory-backed mode is the default:

```bash
cd src && go run .
```

Redis-backed mode uses the Redis instance from [docker-compose.yml](docker-compose.yml):

```bash
docker compose up -d redis
cd src
LINE_STORE=redis REDIS_ADDR=localhost:6379 go run .
```

The app serves the UI at `http://localhost:8080`. Opening `/` creates a new session and redirects to a shareable URL like `http://localhost:8080/sessions/<uuid>`.

## Environment variables:

- `SC_PORT`: HTTP port, defaults to `8080`
- `LINE_STORE`: `memory` or `redis`, defaults to `memory`
- `REDIS_ADDR`: Redis host and port, defaults to `localhost:6379`
- `REDIS_PASSWORD`: optional Redis password
- `REDIS_DB`: optional Redis database number, defaults to `0`

## REST API

OpenAPI description: [doc/openapi/openapi.yaml](./doc/openapi/openapi.yaml)

Session-scoped endpoints:

- `GET /api/sessions/{uuid}/lines`
- `POST /api/sessions/{uuid}/lines`
- `POST /api/sessions/{uuid}/lines/clear`

## Prompts history

History of major propmpts used in project development (AI backend: Github Copilot): 
[doc/prompts_history.md](./doc/prompts_history.md)