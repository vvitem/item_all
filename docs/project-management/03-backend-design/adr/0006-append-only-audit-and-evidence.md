> 状态：Accepted  
> 负责人：vvitem  
> 最后更新：2026-07-11  
> 基线 Commit：`b590887fea1e2c43e7831a48932b9be5d44ccfa3`  
> 关联 Milestone：`M0-foundation`  
> 关联 Issue/PR：待创建

# ADR 0006：Append-only Audit 与独立 Evidence

- **日期**：2026-07-11
- **状态**：Accepted

## 背景

仅保存聊天或日志无法证明执行决策与结论，完整原始输出又增加泄漏。

## 决策

Audit 保存决策事件和 hash；Evidence 保存脱敏结构化依据及受控对象引用。

## 备选方案

- 维持临时实现并后补边界。
- 使用另一套桌面/执行架构。
- 将能力交给云端服务。

以上方案因安全边界、单人维护或本地优先目标不符合而未选择。

## 影响

增加存储和清理设计，但支持验证、报告和安全追踪。

## 验证

通过对应 Backlog、架构测试和 Beta 门禁验证；若事实推翻决策，创建新 ADR supersede 本记录。
