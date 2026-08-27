# BENZHI 验证说明

城市配电网故障恢复协调服务是一套面向配电网调度员和现场抢修队伍的故障恢复计划协同后端服务。

常用命令：

```powershell
go test ./...
go vet ./...
go build ./...
```

容器验证使用 `golang:1.23.12`，工作目录挂载到 `/src` 后运行 `sh -c "go test ./... && go vet ./... && go build ./..."`。业务 smoke 需要先启动 `cmd/recovery`，再执行 `go run ./cmd/smoke -base http://localhost:<port>`。
