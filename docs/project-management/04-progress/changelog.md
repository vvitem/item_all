> 状态：维护中  
> 负责人：vvitem  
> 最后更新：2026-07-14  
> 基线 Commit：`620272272991ec47a39458e2dcb03a0ca7648e95`  
> 关联 Milestone：`M0-foundation`  
> 关联 Issue/PR：[PR #11](https://github.com/vvitem/item_all/pull/11) · [Issue #12](https://github.com/vvitem/item_all/issues/12)

# 项目管理变更日志

## 2026-07-14 — M0-003 完成同步

### Added

- 创建 [Issue #12](https://github.com/vvitem/item_all/issues/12)，跟踪 `M0-004` SQLite WAL、Migration 和 Repository 骨架。
- 为 M0-004 记录 Driver、Migration、数据目录、Repository/Tx 边界与失败路径验收要求。

### Changed

- [PR #11](https://github.com/vvitem/item_all/pull/11) 已 squash 合并，Merge Commit `620272272991ec47a39458e2dcb03a0ca7648e95`。
- `M0-003` 从 `IN_REVIEW` 更新为 `DONE`。
- `M0-004` 从 `NOT_STARTED` 更新为 `READY`。
- 项目状态更新为 `DONE=3`、`READY=1`、`NOT_STARTED=46`。
- `M0-foundation` 完成度更新为 `3/10 = 30%`，全项目完成度更新为 `3/50 = 6%`。
- `main` Ruleset 配置 PR 合并、分支最新、三个 required checks、禁止删除和 force push。

### Validation

- 最终 PR Head Run `29314687975`：`quality`、`generated-and-security`、`windows-build` 全部成功。
- 受控门禁 RED→GREEN：`29313493020` → `29313730163`。
- 代码审查回归 RED→GREEN：`29314236649` → `29314452432`。
- 本状态同步 PR 必须通过三个 required checks 后才可合并。

### Business Code

- 本次状态同步不修改业务代码。
- M0-004 仅进入设计准备，尚未引入 SQLite Driver、Migration 或 Repository 实现。

## 2026-07-13 — M0-002 完成同步

### Added

- 创建 [Issue #7](https://github.com/vvitem/item_all/issues/7)，跟踪 `M0-003` 正式 CI 门禁。
- 为 M0-003 记录 Go、Frontend、Windows build、Binding 漂移、Secret Scan 和故意失败验收要求。

### Changed

- [PR #6](https://github.com/vvitem/item_all/pull/6) 已 squash 合并，Merge Commit `dbb104e3b457d0ef83ca563ff2ae6106c74b43a1`。
- `M0-002` 从 `IN_REVIEW` 更新为 `DONE`。
- `M0-003` 从 `NOT_STARTED` 更新为 `READY`。
- 项目状态更新为 `DONE=2`、`READY=1`、`NOT_STARTED=47`。
- `M0-foundation` 完成度更新为 `2/10 = 20%`，全项目完成度更新为 `2/50 = 4%`。
- 当前基线更新为 M0-002 Merge Commit。

### Validation

- M0-002 以 PR #6、Merge Commit、Go/Frontend 测试、Windows production/dev smoke 和安全边界扫描作为完成证据。
- 状态同步只修改项目管理文档，没有修改业务代码。

### Business Code

- 未修改业务代码。

## 2026-07-13 — M0-002 实现评审

### Added

- 新增 ItemAll Go 1.26.5 + Wails v2.13.0 标准根目录工程。
- 新增只读 `GetAppInfo()` Binding 和 `internal/buildinfo`。
- 新增 React/TypeScript/Vite 页面及 Loading、Ready、Error、Retry 测试。
- 新增同步 Binding 抛错回归测试，验证异常文本和本地路径不会进入用户界面。
- 新增官方 Wails 平台资源、Go/pnpm 锁文件、中文开发说明和 Makefile。
- 新增 pnpm 供应链策略，仅批准 `esbuild` 安装构建脚本。

### Changed

- `M0-002` 从 `READY` 更新为 `IN_REVIEW`，关联 [PR #6](https://github.com/vvitem/item_all/pull/6)。
- 项目状态更新为 `DONE=1`、`IN_REVIEW=1`、`NOT_STARTED=48`。
- 实现路径由建议路径更新为 `main.go`、`app.go`、`internal/buildinfo`、`frontend/src` 和 `build`。
- Binding 调用统一进入 Promise 链，同步抛错与异步 rejection 均进入安全错误状态。

### Validation

- 5 个前端测试、TypeScript、ESLint、Vite build、Go tests、race、Binding 漂移检查通过。
- Windows production build、EXE 启动烟测和 `wails dev` 就绪烟测通过。
- 精确依赖、秘密、范围外能力和远程前端资源扫描通过。
- 同步异常测试完成 RED→GREEN；最终前端验证 Run `29244941752` 通过。
- 临时验证 Workflow 已删除；正式 CI 留在 M0-003。

### Business Code

- 新增最小桌面工程骨架和构建信息展示。
- 未新增 SSH、数据库、SQLite 业务表、AI、Operation Bus、云服务、凭据或遥测。

## 2026-07-13 — M0-001 状态同步

### Added

- 创建 [Issue #2](https://github.com/vvitem/item_all/issues/2)，跟踪 `M0-002` 最小 Go/Wails/React 工程骨架。
- 为 `M0-002` 记录验收标准、测试要求、范围边界和设计前置决策。

### Changed

- [PR #1](https://github.com/vvitem/item_all/pull/1) 已合并，`M0-001` 从 `IN_REVIEW` 更新为 `DONE`。
- 当前基线更新为 `be37cfbb2bdca5a05c3303e9ea208992a8bf1721`。
- `M0-foundation` 完成度更新为 `1/10`，全项目完成度更新为 `1/50`。
- 同步当前状态、Backlog、任务追踪、实现追踪和周报。

### Business Code

- 未修改业务代码。

## 2026-07-11

### Added

- 建立 `docs/project-management` 五大核心区域。
- 保存原始增强版规划 Reference。
- 建立仓库基线、As-Is/To-Be/Gap、MVP 和 20 场景。
- 建立 5 个 Milestone、50 项 Backlog 和 24 周计划。
- 建立 Operation Bus、Task Engine、工具、数据、安全、审计和测试设计。
- 建立 6 个 ADR、当前状态、风险、阻塞和实现追踪。
- 建立交付、Beta 和安全发布门禁。
- 建立 GitHub Issue/PR 模板。

### Changed

- 根 README 增加项目原则、范围、文档入口和真实基线。

### Business Code

- 未修改业务代码。