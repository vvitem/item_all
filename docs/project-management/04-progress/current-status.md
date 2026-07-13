> 状态：维护中  
> 负责人：vvitem  
> 最后更新：2026-07-13  
> 基线 Commit：`dbb104e3b457d0ef83ca563ff2ae6106c74b43a1`  
> 关联 Milestone：`M0-foundation`  
> 关联 Issue/PR：[Issue #2](https://github.com/vvitem/item_all/issues/2) · [PR #6](https://github.com/vvitem/item_all/pull/6) · [Issue #7](https://github.com/vvitem/item_all/issues/7)

# 当前项目状态

## 基线信息

- 基线分支：`main`
- 当前基线 Commit：`dbb104e3b457d0ef83ca563ff2ae6106c74b43a1`
- 最后更新时间：2026-07-13
- 当前 Milestone：`M0-foundation`
- 当前周次：第 2 周
- 项目状态：`DONE=2`、`READY=1`、`NOT_STARTED=47`

## 本阶段目标

建立可构建的 Go/Wails/React 工程骨架、正式 CI 门禁、SQLite WAL、OS Keychain、Operation Bus、append-only Audit，以及 SSH 资产连接测试垂直切片。

## 当前正在进行

- `M0-003`：建立 Go、Frontend、Windows build、生成文件漂移和秘密扫描 CI 门禁，状态 `READY`。
- 跟踪入口：[Issue #7](https://github.com/vvitem/item_all/issues/7)。
- 下一执行动作是编写 M0-003 设计与实现计划，再按 TDD 和故意失败验证实施正式 Workflow。
- 在 M0-003 完成前不开始 SSH、数据库或 AI 功能。

## 最近完成

- `M0-002` 已通过 [PR #6](https://github.com/vvitem/item_all/pull/6) squash 合并，状态更新为 `DONE`。
- 合并 Commit：`dbb104e3b457d0ef83ca563ff2ae6106c74b43a1`。
- 完成 Go 1.26.5、Wails v2.13.0、React/TypeScript/Vite 最小桌面工程骨架。
- 完成唯一只读 `GetAppInfo()` Binding、构建元数据及 Loading/Ready/Error/Retry 页面。
- 完成同步 Binding 抛错与异步 rejection 的统一安全错误处理。
- 完成 Windows production build、EXE 启动、Wails development readiness、Go/Frontend 测试和范围扫描。

## 下一步任务

1. 完成 M0-002 合并后的状态同步 PR。
2. 关闭 Issue #2，并保留 PR #6 与 Merge Commit 作为完成证据。
3. 为 Issue #7 编写 M0-003 设计 Spec 与实施计划。
4. 建立永久 CI Workflow，并通过故意失败与恢复绿色验证门禁。
5. M0-003 合并后再评估 M0-004 与 M0-006 的执行顺序。

## 当前阻塞

无外部阻塞。M0-003 的主要设计约束是：CI 必须覆盖 Windows Wails 构建，但不能把发布、代码签名或业务功能混入本任务。

## 高风险事项

- CI 只验证 Linux 会遗漏 Wails/WebView2 与 Windows 构建问题。
- 未验证生成 Binding 和锁文件漂移会导致本地与提交内容不一致。
- Secret Scan 若只检查常见扩展名，可能遗漏文档、日志或生成文件中的秘密。
- Operation Bus 尚未实现，后续运维能力不得绕过其设计边界。
- Windows Keychain/PTTY/WebView 的后续兼容风险仍需独立验证。

## 待确认决策

- Windows Credential Manager 封装方案。
- 手工 SQL 写操作是否进入首版。
- Windows 稳定版代码签名方案。

## 范围变化

无。M0-002 只建立桌面工程骨架和只读构建信息；M0-003 只建立工程质量门禁，均不包含 SSH、数据库、SQLite 业务表、AI、Operation Bus、云服务或遥测。

## 测试状态

- Go：`go test ./...` 通过，`internal/buildinfo` 竞态检测通过。
- Frontend：5 个 Vitest 用例、TypeScript、ESLint、Vite production build 全部通过。
- 错误回归：同步 Binding 抛错与异步 rejection 均进入安全错误页，不展示异常路径或秘密文本。
- Wails Binding：仅生成 `GetAppInfo()`，生成文件和锁文件无漂移。
- Windows Production：`ItemAll.exe` 构建成功，文件大小 11,413,504 bytes，启动烟测通过。
- Windows Development：`wails dev` 完成编译、WebView2 初始化并进入目录监听。
- 安全边界：精确依赖检查、秘密扫描、范围外能力扫描和远程资源扫描通过。

## 相关代码路径

- `main.go`、`app.go`、`internal/buildinfo/`
- `frontend/src/`、`frontend/wailsjs/`
- `go.mod`、`frontend/package.json`、`frontend/pnpm-lock.yaml`
- `build/`、`wails.json`、`Makefile`

## 相关 Issue / PR / Commit

- [Issue #2](https://github.com/vvitem/item_all/issues/2)：M0-002 工程骨架，等待状态同步后关闭。
- [PR #4](https://github.com/vvitem/item_all/pull/4)：M0-002 设计，已合并。
- [PR #5](https://github.com/vvitem/item_all/pull/5)：M0-002 实现计划，已合并。
- [PR #6](https://github.com/vvitem/item_all/pull/6)：M0-002 实现，已 squash 合并。
- M0-002 Merge Commit：`dbb104e3b457d0ef83ca563ff2ae6106c74b43a1`。
- [Issue #7](https://github.com/vvitem/item_all/issues/7)：M0-003 CI 门禁，`READY`。
- 自动验证：Actions Run `29242978948`、`29243245959`、`29244060195`、`29244941752`。

## 本周可演示结果

Windows 环境可以构建并启动 ItemAll 桌面空壳；页面通过唯一只读 Binding 展示产品名称、版本、Commit、构建时间、Go Runtime 和运行状态，并具备可验证的 Loading、Error 与 Retry 行为。