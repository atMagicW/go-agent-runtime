# go-agent-runtime

## 跑起来的顺序
```
docker compose -f deployments/compose/docker-compose.yaml up -d
export POSTGRES_DSN="postgres://agent:agent@localhost:5432/agent_runtime?sslmode=disable"
go run cmd/server/main.go
```
