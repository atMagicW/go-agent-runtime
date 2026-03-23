# Go Agent Runtime

一个面向生产架构设计的 Agent Runtime（Golang 实现），支持：

* 多 LLM 调度（OpenAI / DeepSeek）
* Skill / Tool / MCP 统一能力模型
* RAG（检索增强生成）
* 意图识别 + 任务编排（Planner）
* 多存储模式（memory / file / postgres）
* 流式输出（SSE）
* 成本统计 / 熔断 / 降级

---

## ✨ 特性

### 🧠 Agent Runtime 核心能力

* 意图识别（规则 + LLM 混合）
* Planner（任务拆解）
* Orchestrator（执行编排）
* Capability Router（统一调度 Skill / Tool / MCP）
* Response Composer（统一输出）

### 🤖 多模型调度

* 支持 OpenAI / DeepSeek
* 按任务类型自动选择模型
* 支持 fallback / 熔断

### 🔧 能力体系

* Skill（任务型能力，支持文件化声明）
* Tool（原子能力）
* MCP Tool（远程能力）

### 📚 RAG

* 多知识库
* 向量检索 + 关键词检索
* rerank（可选）

### 💾 存储模式

* memory（无依赖）
* file（本地持久化）
* postgres（生产级）

---

## 🚀 快速启动

### 1. 修改配置

```yaml
# configs/app.yaml
storage:
  mode: "file"   # memory / file / postgres

rag:
  embedding_provider: "mock"
  seed_on_bootstrap: true
```

---

### 2. 启动服务

```bash
go run cmd/server/main.go
```

---

### 3. 测试接口

```bash
curl localhost:8080/v1/chat \
  -H "Content-Type: application/json" \
  -d '{
    "session_id":"demo",
    "user_id":"u1",
    "message":"请分析这个agent架构",
    "stream":false
  }'
```

---

## 📂 文档

* [架构说明](docs/architecture.md)
* [目录结构](docs/structure.md)
* [执行链路](docs/runtime-flow.md)
* [运行方式](docs/deployment.md)

---

## 🧩 项目定位

本项目不是简单 demo，而是：

> 一个具备生产演进能力的 Agent Runtime 架构实现

适用于：

* AI 应用后端
* 企业级 Agent 平台
* 多模型调度系统
* RAG + Tool + Workflow 场景

---

## 📌 后续扩展方向

* 分布式执行（task queue）
* 多 Agent 协作
* 可视化工作流
* 权限系统（Skill / MCP）
* 长期记忆（Memory Layer）

---

## 📝 License

MIT
