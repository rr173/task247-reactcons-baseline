# 化学反应网络守恒约束诊断服务 task247-reactcons

面向反应工程研究者的纯后端领域计算服务：导入物种组成、反应式与可选实验约束，构建化学计量矩阵，逐反应检查元素/电荷守恒，求守恒子空间（左零空间），定位破坏守恒的最小反应集（MCS），并发布不可变网络版本。无前端页面，通过 HTTP `/api` 提供 23 个端点。

## 构建与门禁（固定环境）
```bash
export CGO_ENABLED=0 GOTOOLCHAIN=local
GO_BIN=/opt/homebrew/bin/go
"$GO_BIN" build ./...
"$GO_BIN" vet   ./...
"$GO_BIN" test  ./...
"$GO_BIN" run ./cmd/reactcons --smoke-test
```

## 运行
```bash
"$GO_BIN" run ./cmd/reactcons --addr :8080 --db reactcons.db
```

## API 入口（/api）

| 能力 | 端点 |
|---|---|
| 自检 | `GET /api/self-check` |
| 网络 | `POST /api/networks` · `GET /api/networks` · `GET /api/networks/{id}` |
| 物种 | `POST /api/networks/{id}/species` · `GET /api/networks/{id}/species` · `PUT /api/species/{id}` |
| 反应 | `POST /api/networks/{id}/reactions` · `GET /api/networks/{id}/reactions` · `POST /api/reactions/{id}/exempt` |
| 求解 | `POST /api/networks/{id}/solve` · `GET /api/networks/{id}/conservation` · `GET /api/networks/{id}/conserved-pools` · `GET /api/networks/{id}/conflicts` |
| 边界 | `POST /api/networks/{id}/boundaries` · `GET /api/networks/{id}/boundaries` · `DELETE /api/networks/{id}/boundaries/{sid}` |
| 约束 | `POST /api/networks/{id}/constraints` · `GET /api/networks/{id}/constraints` |
| 版本 | `POST /api/networks/{id}/versions` · `GET /api/networks/{id}/versions` · `GET /api/versions/{id}` · `GET /api/versions/{id}/verify` |

## 业务包
`chem`（化学式/方程解析）·`matrix`（计量矩阵/零空间/守恒子空间）·`diagnose`（守恒冲突 + MCS）·`version`（content-hash 不可变快照）·`store`（SQLite 持久化）·`service`（编排）·`httpapi`（HTTP 层）。

## 持久化
SQLite（`modernc.org/sqlite`），表：`networks/species/reactions/conservation_pools/conflict_sets/boundaries/constraints/versions/metrics`。`sealed` 网络与 `published` 版本冻结，禁止修改。
