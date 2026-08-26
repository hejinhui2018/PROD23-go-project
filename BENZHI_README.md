# BENZHI 验证说明

本项目是面向城市配电网调度员与现场抢修队伍的故障恢复协调后端服务，负责管理恢复计划、执行现场步骤、持久化恢复事件并发送完成通知。

常用命令：

```powershell
go test ./...
go vet ./...
go build ./...
```

容器验证使用 `golang:1.23.12`，工作目录挂载到 `/src` 后运行 `sh -c "go test ./... && go vet ./... && go build ./..."`。业务 smoke 需要先启动 `cmd/recovery`，再执行 `go run ./cmd/smoke -base http://localhost:<port>`。
