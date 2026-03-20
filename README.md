# go-agent-runtime

一个基于 Go 的可插拔 Agent Runtime Demo，支持多模型路由、Skill/Tool/MCP 统一调度、多知识库 RAG、任务编排与工程化治理

## 核心能力

### 1. Multi-LLM Routing
支持按任务类型、用户指定、降级策略选择不同模型。

### 2. Capability Routing
将本地 Skill、Tool 与远程 MCP Tool 统一为 Capability，通过 Registry 管理并由 Router 分发执行。

### 3. Multi-KB RAG
支持知识库隔离、文本 ingest、切块、向量检索与关键词回退检索。

### 4. Planning & Orchestration
支持意图识别、执行计划生成、步骤依赖与并发执行。

### 5. Governance
支持请求超时、步骤重试、熔断与降级、模型成本统计与审计日志。

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

