> 状态：草案  
> 负责人：vvitem  
> 最后更新：2026-07-11  
> 基线 Commit：`b590887fea1e2c43e7831a48932b9be5d44ccfa3`  
> 关联 Milestone：`M0-foundation → M3-safe-ai`  
> 关联 Issue/PR：待创建

# IPC 与接口契约

## Wails IPC

首版 UI 只通过 Wails IPC 调用 Go Application Service。Binding 使用窄 DTO：

- `AssetCreateRequest/AssetView`
- `OperationSubmitRequest/OperationView`
- `TaskCreateRequest/TaskView`
- `TaskCancelRequest`
- `RunbookCompileRequest`

Binding 负责类型、长度和必填校验；业务规则由 Service/Operation Bus 执行。

## 事件

后端向前端推送版本化事件：

```json
{
  "version": 1,
  "type": "task.step.updated",
  "trace_id": "...",
  "entity_id": "...",
  "sequence": 12,
  "occurred_at": "2026-07-11T00:00:00Z",
  "payload": {}
}
```

前端必须按 `sequence` 去重和检测缺口；缺口时重新读取权威 Task 状态。

## 未来 CLI/MCP

仅预留统一 Envelope，不在六个月实现网络服务：

```text
RequestEnvelope(version, workspace_id, actor_id, device_id, source, trace_id, payload)
```

若未来实现，本地 IPC/Socket 层必须将请求转换为同一个 `OperationRequest`，不得复制执行逻辑。

## 兼容策略

- DTO 增加字段应保持向后兼容。
- 删除/改义字段需要新版本。
- Plan、Runbook、Evidence Schema 必须含版本并提供迁移器。
