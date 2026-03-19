# go-agent-runtime

一个基于 Go 实现的可插拔 Agent Runtime Demo，支持：

- 多 LLM 路由
- Skill / Tool / MCP 统一 Capability 调度
- 多知识库 RAG 路由
- 意图识别与任务编排
- 会话持久化
- Prompt 模板版本管理
- 模型成本统计
- 审计日志
- 熔断与降级
- 知识库 ingest

## 项目结构

- `cmd/server`：程序入口
- `api/httpapi`：HTTP 接口层
- `internal/app`：应用服务层
- `internal/usecase`：意图识别、规划、编排、路由、治理
- `internal/adapters`：数据库、LLM、RAG、Capability 适配层
- `internal/domain`：核心领域模型
- `migrations`：数据库初始化 SQL
- `configs`：配置文件

## 启动依赖

### 1. 启动 PostgreSQL / pgvector

```bash
make docker-up