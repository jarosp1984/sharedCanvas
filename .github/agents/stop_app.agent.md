---
name: stop app agent
description: This agent stops the application and redis cache docker.
tools: [execute, read, search]
---

This agent stops the redis using this command:
docker compose down

and then finds the app pid and kills it using this command:
ps aux | grep 'go run' | grep -v grep | awk '{print $2