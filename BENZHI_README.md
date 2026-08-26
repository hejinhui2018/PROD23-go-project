# BENZHI 验证说明

本仓库是 Go 1.23 后端服务，验证入口以标准 Go 命令为主。

常用命令：

```powershell
go test ./...
go vet ./...
go build ./...
```

容器验证使用 `golang:1.23.12`，工作目录挂载到 `/src` 后运行 `sh -c "go test ./... && go vet ./... && go build ./..."`。业务 smoke 需要先启动 `cmd/recovery`，再执行 `go run ./cmd/smoke -base http://localhost:<port>`。
