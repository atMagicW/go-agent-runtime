# 系统架构

## 分层结构

```text
API Layer
   ↓
Application Layer
   ↓
Usecase（Agent Runtime）
   ↓
Ports（接口）
   ↓
Adapters（实现）
```

---

## 核心组件

### Agent Runtime

* Intent Engine
* Planner
* Orchestrator
* Router（Model / Capability / RAG）
* Response Composer

---

## 能力模型

```text
Capability
├── Skill（任务级）
├── Tool（函数级）
└── MCP Tool（远程）
```

---

## 关键设计原则

### 1. 解耦

* LLM 可替换
* RAG 可替换
* Storage 可切换

### 2. 可扩展

* 新增 Skill / Tool 无需改核心逻辑
* 新增 LLM provider 无侵入

### 3. 配置驱动

* 模型选择
* fallback
* routing

---

## 为什么这样设计

因为 Agent 系统需要：

* 多模型协作
* 多能力组合
* 动态执行流程

