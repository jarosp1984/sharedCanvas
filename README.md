# sharedCanvas

## Run

Memory-backed mode is the default:

```bash
go run ./src/main.go
```

Redis-backed mode uses the Redis instance from [docker-compose.yml](/home/jarek/projects/sharedCanvas/docker-compose.yml):

```bash
docker compose up -d redis
LINE_STORE=redis REDIS_ADDR=localhost:6379 go run ./src/main.go
```

The app serves the UI at `http://localhost:8080`.

## REST API

Create a line segment with `POST /api/lines`.

Example request:

```json
{
	"x1": 10,
	"y1": 20,
	"x2": 30,
	"y2": 40,
	"color": "#112233",
	"width": 5
}
```

Validation rules:

- `x1`, `x2` must be between `0` and `800`
- `y1`, `y2` must be between `0` and `600`
- `color` must be a 6-digit hex string like `#112233`
- `width` must be between `1` and `50`

Example response:

```json
{
	"status": "ok",
	"line": {
		"x1": 10,
		"y1": 20,
		"x2": 30,
		"y2": 40,
		"color": "#112233",
		"width": 5
	}
}
```

Environment variables:

- `PORT`: HTTP port, defaults to `8080`
- `LINE_STORE`: `memory` or `redis`, defaults to `memory`
- `REDIS_ADDR`: Redis host and port, defaults to `localhost:6379`
- `REDIS_PASSWORD`: optional Redis password
- `REDIS_DB`: optional Redis database number, defaults to `0`