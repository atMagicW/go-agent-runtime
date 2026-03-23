# 运行模式

## 1. Memory 模式

```yaml
storage:
  mode: "memory"
```

特点：

* 无依赖
* 重启数据丢失
* 适合测试

---

## 2. File 模式（推荐本地开发）

```yaml
storage:
  mode: "file"
  data_dir: "./data"
```

特点：

* session 持久化到文件
* RAG 使用内存
* 无需数据库

---

## 3. Postgres 模式（生产）

```yaml
storage:
  mode: "postgres"
```

特点：

* 完整持久化
* 支持 pgvector
* 支持大规模数据

---

## 模式对比

| 模式       | Session | RAG | 依赖 |
| -------- | ------- | --- | -- |
| memory   | ❌       | ❌   | 无  |
| file     | ✅       | ❌   | 无  |
| postgres | ✅       | ✅   | DB |

---

## 推荐使用方式

### 本地开发

```yaml
mode: file
embedding: mock
```

### 测试

```yaml
mode: memory
```

### 部署

```yaml
mode: postgres
```
