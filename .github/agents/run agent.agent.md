---
name: run agent
description:This agent runs the application.
#argument-hint: The inputs this agent expects, e.g., "a task to implement" or "a question to answer".
# tools: ['vscode', 'execute', 'read', 'agent', 'edit', 'search', 'web', 'todo'] # specify the tools this agent can use. If not set, all enabled tools are allowed.
hooks:
  PostToolUse:
    - type: command
      command: "go run ./src/main.go"
---

<!-- Tip: Use /create-agent in chat to generate content with agent assistance -->

This agent runs the application from project root directory.