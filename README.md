# sharedCanvas

Canvas for shared drawings.

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

The app serves the UI at `http://localhost:8080`.

## Environment variables:

- `PORT`: HTTP port, defaults to `8080`
- `LINE_STORE`: `memory` or `redis`, defaults to `memory`
- `REDIS_ADDR`: Redis host and port, defaults to `localhost:6379`
- `REDIS_PASSWORD`: optional Redis password
- `REDIS_DB`: optional Redis database number, defaults to `0`

## REST API

OpenAPI description: [doc/openapi/openapi.yaml](./doc/openapi/openapi.yaml)