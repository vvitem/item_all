# item_all

Local-first SafeOps 桌面工作台，面向个人开发者、独立运维和全栈工程师，聚焦 Linux 故障诊断与 MySQL/PostgreSQL 数据库运维。

> 当前状态：项目处于 `M0-foundation` 立项与工程骨架阶段。仓库尚未包含业务代码，所有实现状态以 [项目管理入口](docs/project-management/README.md) 为准。

## 项目原则

- 本地优先：连接、凭据、审计和证据默认保存在本机。
- 安全优先：AI 只调用结构化、参数化、只读工具，不获得任意 Shell 或数据库写权限。
- 单一执行边界：UI、AI、Runbook 及未来 CLI/MCP 都必须经过 Operation Bus。
- 可验证：任务遵循 `Goal → Plan → Approval → Execute → Verify → Report`，重要结论必须有 Evidence。
- 范围受控：六个月首版仅覆盖 SSH、SFTP、MySQL、PostgreSQL，Windows 正式支持。

## 文档入口

- [项目管理总览](docs/project-management/README.md)
- [当前状态](docs/project-management/04-progress/current-status.md)
- [六个月路线图](docs/project-management/02-roadmap/six-month-roadmap.md)
- [后端架构](docs/project-management/03-backend-design/architecture-overview.md)
- [MVP 范围](docs/project-management/01-analysis/mvp-scope.md)
- [原始增强版规划](docs/project-management/00-reference/OpsKat_项目深度分析与差异化产品规划_增强版.md)

## 当前基线

- 仓库：https://github.com/vvitem/item_all
- 默认分支：`main`
- 基线 Commit：`b590887fea1e2c43e7831a48932b9be5d44ccfa3`
- 分析日期：2026-07-11
- 当前代码状态：仅存在本 README，Go/Wails/React/SQLite/CI/测试均尚未初始化。
