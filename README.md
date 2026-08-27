# 城市配电网故障恢复协调服务

服务面向调度员和现场队伍。故障创建后依据馈线拓扑评估受影响区段，确认计划后按顺序执行恢复步骤。步骤状态为 `pending -> dispatched -> acknowledged -> completed|failed`，版本冲突、幂等键和前置约束均由领域服务校验。事件日志采用 JSONL 追加存储，启动时重放；通知队列支持重放。

运行：`go run ./cmd/recovery -config ./config.example.json`

Smoke：`go run ./cmd/smoke -base http://localhost:8080`，预期输出故障 ID、计划 ID、步骤完成以及 `restoration.completed` 通知。

角色：调度员创建/确认计划，现场队伍领取并上报步骤。健康检查：`GET /healthz`，指标：`GET /metrics`。
