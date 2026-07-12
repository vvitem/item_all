> 状态：已确认  
> 负责人：vvitem  
> 最后更新：2026-07-11  
> 基线 Commit：`b590887fea1e2c43e7831a48932b9be5d44ccfa3`  
> 关联 Milestone：`M0-foundation`  
> 关联 Issue/PR：待创建

# 实施工作流

## 单任务循环

1. 从 Backlog 选择一个 `READY` 任务并确认依赖。
2. 写 1–2 页设计补充和可执行验收标准。
3. 先建立失败测试：正常、失败、取消、权限/限制。
4. 使用 Codex 生成最小实现，不允许一次跨越多个模块。
5. 人工审查数据流、权限、日志、依赖和错误传播。
6. 真实 Windows/SSH/MySQL/PostgreSQL 环境验证。
7. 更新六个进度文件和实现追踪。
8. PR 只包含一个可演示垂直结果。

## 分支与提交

- 分支：`feat/M0-002-wails-bootstrap`、`fix/M1-002-reconnect-state`。
- Commit：`type(scope): outcome`，例如 `feat(operation): add guarded execution pipeline`。
- PR 必须链接 Backlog ID、场景 ID、设计文档、测试和风险。

## Codex 使用边界

Codex 可用于样板、测试生成、重构建议和文档同步，但不得自行决定：

- 安全策略和风险等级。
- 凭据生命周期。
- 依赖新增。
- Migration 回滚/数据删除。
- 任务标记 `DONE`。

## 合并门禁

- 构建、lint、单元和关联集成测试通过。
- 无 Secret Scan 告警。
- 新执行路径通过 Operation Bus。
- 新工具更新 Tool Contract、威胁模型和测试矩阵。
- 当前状态和追踪文档已同步。
