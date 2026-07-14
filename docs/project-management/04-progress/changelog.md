> 状态：维护中  
> 负责人：vvitem  
> 最后更新：2026-07-14  
> 基线 Commit：`cf25d4e40f34e0c2e835a2cc04b3eeea4d011ad7`  
> 关联 Milestone：`M0-foundation`  
> 关联 Issue/PR：[Issue #7](https://github.com/vvitem/item_all/issues/7) · [PR #11](https://github.com/vvitem/item_all/pull/11)

# 项目管理变更日志

## 2026-07-14 — M0-003 实现评审

### Added

- 新增 `.github/workflows/ci.yml`，建立 `quality`、`generated-and-security`、`windows-build` 三个稳定 Job。
- 新增 `scripts/ci` Go 检查器，覆盖依赖精确版本、pnpm `allowBuilds`、Workflow 安全、远程资源、环境变量读取、范围外能力和 Wails Binding Allowlist。
- 新增 Windows `ItemAll.exe` 启动观察与进程树清理脚本。
- 新增 `docs/development/ci.md`，记录本地复现、故障排查、Action SHA、Gitleaks、Binding 漂移和 Branch Protection 自举。
- 新增完整历史 Gitleaks v8.30.0 扫描与三类生成漂移门禁。

### Changed

- `M0-003` 从 `READY` 更新为 `IN_REVIEW`，关联 [PR #11](https://github.com/vvitem/item_all/pull/11)。
- 项目状态更新为 `DONE=2`、`IN_REVIEW=1`、`NOT_STARTED=47`。
- 所有 Action 使用官方完整 Commit SHA；Workflow 权限固定为 `contents: read`。
- Linux Wails Binding 生成后恢复普通文件模式 `0644`，再执行严格内容与模式 diff。
- Gitleaks 安装路径使用 v8.30.0 `go.mod` 声明的 `github.com/zricethezav/gitleaks/v8`。

### Validation

- Actions Run `29312670714`：三个 Job 全部成功。
- Go format/vet/test/race、Frontend test/typecheck/lint/build 通过。
- 实际 Workflow 自检、依赖政策、仓库边界、Binding/Go Module/pnpm Lockfile 漂移通过。
- 完整历史 Gitleaks 通过。
- Windows Wails production build、非空 EXE 和启动烟测通过。
- 临时 Binding 诊断 Workflow 已删除；最终只保留正式 CI Workflow。
- 受控 RED→GREEN 证明仍待完成，因此 M0-003 不提前标记 `DONE`。

### Business Code

- 未新增 SSH、数据库、SQLite 业务表、AI、Operation Bus、凭据、云服务、遥测或发布功能。

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
