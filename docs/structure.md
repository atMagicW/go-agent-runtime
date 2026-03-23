# 项目目录结构

```text
go-agent-runtime/
├── cmd/
│   └── server/
│       └── main.go              # 服务入口，负责启动 HTTP Server + Bootstrap
│
├── api/
│   └── httpapi/                # HTTP 接口层（Gin）
│       ├── handler_*.go        # 各类 handler（chat / rag / skill 等）
│       ├── router.go           # 路由注册
│       ├── stream.go           # SSE 流式输出封装
│       └── publisher.go        # 运行时事件 -> SSE
│
├── configs/                    # 配置文件（可热替换）
│   ├── app.yaml
│   ├── models.yaml
│   ├── capabilities.yaml
│   ├── knowledge_bases.yaml
│   ├── fallback.yaml
│   └── pricing.yaml
│
├── internal/
│
│   ├── app/                    # 应用层（组合所有组件）
│   │   ├── bootstrap.go        # 核心：构建整个系统依赖
│   │   ├── agent_service.go    # Agent 主入口
│   │   ├── session_service.go
│   │   ├── rag_service.go
│   │   ├── ingest_service.go
│   │   └── skill_registry.go
│
│   ├── domain/                 # 领域模型（纯结构 + 常量）
│   │   ├── agent/
│   │   ├── capability/
│   │   ├── rag/
│   │   └── model/
│
│   ├── ports/                  # 接口定义（核心抽象层）
│   │   ├── llm.go
│   │   ├── rag_repository.go
│   │   ├── session_repository.go
│   │   ├── capability.go
│   │   └── event_publisher.go
│
│   ├── adapters/               # 适配层（实现 ports）
│   │
│   │   ├── llm/
│   │   │   ├── openai/
│   │   │   └── deepseek/
│   │
│   │   ├── rag/
│   │   │   ├── memory/
│   │   │   ├── pgvector/
│   │   │   └── openai_embedding/
│   │
│   │   ├── persistence/
│   │   │   ├── memory/
│   │   │   ├── file/
│   │   │   └── postgres/
│   │
│   │   ├── capability/
│   │   │   ├── skills/
│   │   │   ├── tools/
│   │   │   └── mcp/
│   │
│   │   └── skillloader/        # 从 skills/ 目录加载 Skill
│
│   ├── usecase/                # 核心业务逻辑（Agent Runtime）
│   │
│   │   ├── intent/             # 意图识别
│   │   │   ├── engine.go
│   │   │   ├── rule_classifier.go
│   │   │   └── llm_classifier.go
│   │
│   │   ├── planner/            # 任务规划
│   │   ├── runtime/            # 执行引擎（orchestrator）
│   │   ├── router/             # model / capability / rag 路由
│   │   └── governance/         # 熔断 / fallback / 限流
│
│   └── pkg/                    # 通用工具
│       ├── config/
│       ├── textsplitter/
│       └── utils/
│
├── skills/                     # Skill 声明（文件化）
│   └── resume-analyzer/
│       ├── manifest.yaml
│       └── SKILL.md
│
├── prompts/                    # Prompt 模板
│
└── docs/                       # 文档
```

---

## 设计特点

* **ports + adapters**：严格分层，可替换实现
* **usecase 独立**：Agent Runtime 不依赖具体实现
* **配置驱动**：模型 / 能力 / fallback 可动态调整
* **多存储模式**：不依赖数据库也可运行
