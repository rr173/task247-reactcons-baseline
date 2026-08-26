基于 Go 实现的化学反应网络守恒约束诊断后端服务 Web 项目，一款后端服务，完成化学计量矩阵构建、元素/电荷守恒冲突定位与守恒子空间投影，并发布不可变网络版本。

# BENZHI 评测说明

## 项目类型
化学反应网络守恒约束诊断服务（纯后端，无前端页面）。研究者导入物种组成、反应式与可选实验约束，服务构建化学计量矩阵、检查元素/电荷守恒、求守恒子空间、定位最小冲突反应集，并发布可校验的不可变网络版本。

## 标准命令（固定环境）
```bash
export CGO_ENABLED=0 GOTOOLCHAIN=local
GO_BIN=/opt/homebrew/bin/go

"$GO_BIN" build ./...
"$GO_BIN" vet   ./...
"$GO_BIN" test  ./...
"$GO_BIN" run ./cmd/reactcons --smoke-test
```

## --smoke-test 契约
`--smoke-test` 不启动长驻服务：真实创建反应网络、导入物种与反应、求解守恒冲突、发布版本，随后**关闭并重新打开同一数据库**验证持久化与重启恢复，最后以 0 退出码结束。这是唯一判据。

## Docker 构建与双架构验证
使用项目提供的 `build_benzhi_docker.sh` 构建评测镜像；Dockerfile 不声明端口，服务监听地址由运行参数 `--addr` 指定。

```bash
./build_benzhi_docker.sh task247-reactcons:amd64 linux/amd64
docker run --rm task247-reactcons:amd64 --smoke-test

./build_benzhi_docker.sh task247-reactcons:arm64 linux/arm64
docker run --rm task247-reactcons:arm64 --smoke-test
```

## 运行（长驻 HTTP 服务）
```bash
"$GO_BIN" run ./cmd/reactcons --addr :8080 --db reactcons.db
docker run --rm -P task247-reactcons:amd64 --addr :8080 --db ./app.db
```

## API 前缀
所有端点以 `/api` 开头，共 23 个（网络/物种/反应/守恒/冲突/边界/约束/版本/自检）。详见 `README.md`。

## 持久化
SQLite（modernc.org/sqlite，CGO 无关，离线可构建）。库由 `--db` 指定，默认 `reactcons.db`。
