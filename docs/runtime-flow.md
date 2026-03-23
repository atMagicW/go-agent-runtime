# Agent 执行链路

## 整体流程

```text
用户请求
   ↓
HTTP API (/v1/chat)
   ↓
AgentService
   ↓
Intent Engine（意图识别）
   ↓
Planner（任务拆解）
   ↓
Orchestrator（执行引擎）
   ↓
执行多个 Step：
   ├── ModelRouter（LLM）
   ├── RAGRouter（检索）
   ├── CapabilityRouter（Skill / Tool / MCP）
   └── ResponseComposer（最终输出）
   ↓
流式输出（SSE）
```

---

## 1️⃣ 意图识别（Intent Engine）

组合策略：

* RuleClassifier（低成本、快速命中）
* LLMClassifier（复杂语义判断）

输出：

```json
{
  "intent_type": "workflow",
  "requires_rag": true,
  "requires_capability": true,
  "requires_planning": true
}
```

---

## 2️⃣ Planner（任务规划）

根据意图生成执行计划：

```text
Step1: RAG 检索
Step2: Tool / Skill 调用
Step3: LLM 生成
Step4: Response Compose
```

支持：

* 并发执行
* 重试
* 超时控制

---

## 3️⃣ Orchestrator（执行引擎）

负责：

* 按顺序执行 step
* 管理上下文
* 处理失败 / fallback
* 发布事件（SSE）

---

## 4️⃣ Capability Router

统一调度：

```text
Skill      -> 本地任务能力
Tool       -> 原子能力
MCP Tool   -> 远程能力
```

---

## 5️⃣ Model Router

根据：

* task type
* tags
* fallback policy

选择最合适的模型：

```text
intent        -> gpt-4.1-mini
analysis      -> deepseek-reasoner
write         -> deepseek-chat
```

---

## 6️⃣ RAG Router

支持：

* 多知识库
* embedding 检索
* keyword 检索
* rerank

---

## 7️⃣ Response Composer

统一输出：

* 合并 step 结果
* 结构化生成
* 流式输出

---

## 关键设计思想

> Agent = Intent + Plan + Execute + Compose
