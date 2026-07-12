> 状态：Accepted  
> 负责人：vvitem  
> 最后更新：2026-07-11  
> 基线 Commit：`b590887fea1e2c43e7831a48932b9be5d44ccfa3`  
> 关联 Milestone：`M0-foundation`  
> 关联 Issue/PR：待创建

# ADR 0002：Operation Bus 作为唯一执行边界

- **日期**：2026-07-11
- **状态**：Accepted

## 背景

多个入口独立调用协议会导致策略、审计和脱敏不一致。

## 决策

所有 UI、AI、Runbook、未来 CLI/MCP 请求转换为 OperationRequest，经同一流水线执行。

## 备选方案

- 维持临时实现并后补边界。
- 使用另一套桌面/执行架构。
- 将能力交给云端服务。

以上方案因安全边界、单人维护或本地优先目标不符合而未选择。

## 影响

增加初始抽象成本，但显著降低绕过和重复安全逻辑风险。

## 验证

通过对应 Backlog、架构测试和 Beta 门禁验证；若事实推翻决策，创建新 ADR supersede 本记录。
