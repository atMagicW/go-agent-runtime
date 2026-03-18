# go-agent-runtime

## 环境变量
启动服务前请设置：
```bash
docker compose -f deployments/compose/docker-compose.yaml up -d
export POSTGRES_DSN="postgres://agent:agent@localhost:5432/agent_runtime?sslmode=disable"
export OPENAI_API_KEY="your_openai_api_key"
go run cmd/server/main.go
