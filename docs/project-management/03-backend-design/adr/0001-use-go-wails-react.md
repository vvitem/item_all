> 状态：Accepted  
> 负责人：vvitem  
> 最后更新：2026-07-11  
> 基线 Commit：`b590887fea1e2c43e7831a48932b9be5d44ccfa3`  
> 关联 Milestone：`M0-foundation`  
> 关联 Issue/PR：待创建

# ADR 0001：使用 Go、Wails v2 与 React/TypeScript

- **日期**：2026-07-11
- **状态**：Accepted

## 背景

单人需要同时交付 Windows 桌面、SSH/DB 协议核心和现代 UI。

## 决策

采用 Go 作为核心与 CLI 友好运行时，Wails v2 作为桌面壳，React/TypeScript 作为前端。

## 备选方案

- 维持临时实现并后补边界。
- 使用另一套桌面/执行架构。
- 将能力交给云端服务。

以上方案因安全边界、单人维护或本地优先目标不符合而未选择。

## 影响

减少多语言运行时；可复用 Go 协议生态。代价是 Wails/WebView 平台差异必须通过 Windows E2E 控制。

## 验证

通过对应 Backlog、架构测试和 Beta 门禁验证；若事实推翻决策，创建新 ADR supersede 本记录。
