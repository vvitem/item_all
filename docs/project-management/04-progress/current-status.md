> 状态：维护中  
> 负责人：vvitem  
> 最后更新：2026-07-13  
> 基线 Commit：`47ad2e428a7db862195b871abbea42ac4c4e930d`  
> 关联 Milestone：`M0-foundation`  
> 关联 Issue/PR：[Issue #2](https://github.com/vvitem/item_all/issues/2) · [PR #6](https://github.com/vvitem/item_all/pull/6)

# 当前项目状态

## 基线信息

- 基线分支：`main`
- 实现分支：`feat/m0-002-project-scaffold`
- 基线 Commit：`47ad2e428a7db862195b871abbea42ac4c4e930d`
- 最后更新时间：2026-07-13
- 当前 Milestone：`M0-foundation`
- 当前周次：第 2 周

## 本阶段目标

建立可构建的 Go/Wails/React 工程骨架、SQLite WAL、OS Keychain、Operation Bus、append-only Audit，以及 SSH 资产连接测试垂直切片。

## 当前正在进行

- `M0-002`：ItemAll 最小 Go/Wails/React 工程骨架已实现，状态 `IN_REVIEW`。
- 实现 PR：[PR #6](https://github.com/vvitem/item_all/pull/6)。
- 完成证据：Go/前端测试、Windows `wails build`、生产 EXE 启动烟测和 `wails dev` 就绪烟测。
- PR 合并前不得将 `M0-002` 标记为 `DONE`。

## 最近完成

- 产品身份已确认：`ItemAll`、Go Module `github.com/vvitem/item_all`、规范化应用标识 `com.vvitem.itemall`。
- M0-002 设计和实现计划已分别通过 PR #4、PR #5 合并。
- 完成 Go 1.26.5、Wails v2.13.0、React/TypeScript/Vite 工程骨架。
- 完成只读 `GetAppInfo()` Binding、构建元数据和 Loading/Ready/Error/Retry 页面。
- 临时 Windows 验证 Workflow 已从实现分支删除，没有提前落地 M0-003 CI。

## 下一步任务

1. 完成 PR #6 的整体规范、代码质量和范围复核。
2. 将 PR #6 标记为 Ready for Review，等待合并决策。
3. PR #6 合并后通过独立文档 PR 将 `M0-002` 更新为 `DONE`。
4. 合并后将 `M0-003` 更新为 `READY`，开始正式 CI 门禁设计。
5. 在 M0-003 完成前不开始 SSH、数据库或 AI 功能。

## 当前阻塞

无实现阻塞。当前只有合并门：PR #6 尚未合并，因此 `M0-002` 保持 `IN_REVIEW`。

## 高风险事项

- 首个业务代码 PR 必须确认生成 Binding、锁文件和 Wails 平台资源未被手工篡改。
- Windows Keychain/PTTY/WebView 的后续兼容风险尚未在 M0-002 范围内验证。
- Operation Bus 尚未实现，后续运维能力不得绕过其设计边界。
- Codex 生成代码必须持续经过秘密、远程资源和范围外能力扫描。

## 待确认决策

- Windows Credential Manager 封装方案。
- 手工 SQL 写操作是否进入首版。
- Windows 稳定版代码签名方案。

## 范围变化

无。M0-002 只建立桌面工程骨架和只读构建信息，没有引入 SSH、数据库、SQLite 业务表、AI、Operation Bus、云服务或遥测。

## 测试状态

- Go：`go test ./...` 通过，`internal/buildinfo` 竞态检测通过。
- Frontend：4 个 Vitest 用例、TypeScript、ESLint、Vite production build 全部通过。
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

- [Issue #2](https://github.com/vvitem/item_all/issues/2)：M0-002 工程骨架，Open。
- [PR #4](https://github.com/vvitem/item_all/pull/4)：M0-002 设计，已合并。
- [PR #5](https://github.com/vvitem/item_all/pull/5)：M0-002 实现计划，已合并，Merge Commit `47ad2e428a7db862195b871abbea42ac4c4e930d`。
- [PR #6](https://github.com/vvitem/item_all/pull/6)：工程骨架实现，当前 `IN_REVIEW`。
- 自动验证：Actions Run `29242978948`、`29243245959`、`29244060195`。

## 本周可演示结果

Windows 环境可以构建并启动 ItemAll 桌面空壳；页面通过唯一只读 Binding 展示产品名称、版本、Commit、构建时间、Go Runtime 和运行状态，并具备可验证的 Loading、Error 与 Retry 行为。
