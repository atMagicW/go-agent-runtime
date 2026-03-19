#!/usr/bin/env bash
set -e

export POSTGRES_DSN="${POSTGRES_DSN:-postgres://agent:agent@localhost:5432/agent_runtime?sslmode=disable}"
export OPENAI_API_KEY="${OPENAI_API_KEY:-}"

go run cmd/server/main.go