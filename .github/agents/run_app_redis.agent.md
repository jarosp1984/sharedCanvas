---
name: run app redis agent
description: This agent runs the application using Redis-backed operation store.
tools: ['execute']
---

This agent runs the application from src directory using this commands:
docker compose up -d redis
cd src
LINE_STORE=redis REDIS_ADDR=localhost:6379 go run .