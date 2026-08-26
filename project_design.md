# 项目设计

系统围绕配电故障恢复计划组织业务状态。

调度员创建故障后，领域拓扑会计算受影响区段；确认计划会生成有序恢复步骤。现场队伍必须按 `pending -> dispatched -> acknowledged -> completed|failed` 的状态机推进步骤，前序步骤未完成时后续步骤不能领取。所有变更写入 JSONL 事件日志，服务启动时由 `internal/service/recovery.go` 重放事件，重建故障、计划、步骤、通知和幂等结果索引。

主要模块：

- `internal/domain` 定义故障、计划、步骤、事件和状态转移约束。
- `internal/service` 实现调度、现场命令、恢复重放、审计和通知触发。
- `internal/store` 提供 JSONL 事件存储、通知队列和快照辅助能力。
- `internal/httpapi` 暴露故障、计划、步骤、通知、审计和健康检查接口。
- `cmd/recovery` 启动 HTTP 服务，`cmd/smoke` 执行端到端业务烟测。

持久化边界是业务可靠性的关键：命令幂等键需要随事件一起落盘，重启后才能安全处理客户端使用原请求标识发起的重试。
