> 状态：已确认  
> 负责人：vvitem  
> 最后更新：2026-07-11  
> 基线 Commit：`b590887fea1e2c43e7831a48932b9be5d44ccfa3`  
> 关联 Milestone：`M0-foundation`  
> 关联 Issue/PR：待创建

# Definition of Done

任务标记 `DONE` 必须同时满足：

- 实现或文档已落入目标分支。
- 验收标准逐条通过，并记录证据。
- 正常、失败、取消/超时路径经过验证。
- 相关单元/集成/E2E 测试通过。
- 没有绕过 Operation Bus、Audit、Evidence、CredentialRef。
- 新增外部输出经过限流和脱敏。
- 安全影响和依赖变更已审查。
- 文档、Backlog、Current Status、Trace 已同步。
- 存在 Commit/PR/测试报告或可复现步骤。
- 没有未解释的 `TODO`、静默 fallback 或未知状态。

文档任务还必须满足：文件非空、链接有效、不将目标设计写成当前实现。
