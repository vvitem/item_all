> 状态：维护中  
> 负责人：vvitem  
> 最后更新：2026-07-14  
> 基线 Commit：`cf25d4e40f34e0c2e835a2cc04b3eeea4d011ad7`  
> 关联 Milestone：`M0-foundation`  
> 关联 Issue/PR：[Issue #7](https://github.com/vvitem/item_all/issues/7) · [PR #11](https://github.com/vvitem/item_all/pull/11)

# 周报日志

## 2026-W29（第 2 周）

### 目标

完成 `M0-002` 合并验收，并建立 `M0-003` 正式分层 CI 门禁，使 Go、Frontend、生成文件、安全边界和 Windows Wails 构建在 Pull Request 中可验证。

### 已完成

- [PR #6](https://github.com/vvitem/item_all/pull/6) squash 合并 M0-002，Merge Commit `dbb104e3b457d0ef83ca563ff2ae6106c74b43a1`。
- `M0-002` 更新为 `DONE`，创建 [Issue #7](https://github.com/vvitem/item_all/issues/7) 跟踪 M0-003。
- [PR #9](https://github.com/vvitem/item_all/pull/9) 合并 M0-003 分层 CI 设计。
- [PR #10](https://github.com/vvitem/item_all/pull/10) 合并 M0-003 Implementation Plan。
- 创建 [PR #11](https://github.com/vvitem/item_all/pull/11)，实现单 Workflow 三个稳定 Job：`quality`、`generated-and-security`、`windows-build`。
- 固定 Go 1.26.5、Wails v2.13.0、Node 24.18.0、pnpm 11.12.0、Gitleaks v8.30.0、Ubuntu 24.04 和 Windows Server 2025。
- 所有 Action 固定官方完整 Commit SHA；Workflow 权限仅 `contents: read`，不使用 Secret、`pull_request_target` 或自托管 Runner。
- 新增 `scripts/ci` Go 检查器及依赖、Workflow、安全边界测试。
- 新增完整历史 Gitleaks、Binding/Go Module/pnpm Lockfile 漂移检查。
- 新增 Windows Wails production build、EXE 存在性和启动烟测。
- 定位并修复两类 Runner 差异：前端产物必须先生成供 `go:embed` 使用；Wails Linux 生成 Binding 后需恢复普通文件模式 `0644`。
- 删除所有临时 Binding 诊断 Workflow，最终只保留正式 `.github/workflows/ci.yml`。
- Actions Run `29312670714` 的三个 Job 全部通过。
- 完成受控 RED→GREEN：RED Run `29313493020` 仅在 Go 格式门禁失败；删除 Fixture 后 GREEN Run `29313730163` 三个 Job 全部成功。

### 未完成

- PR #11 尚未完成整分支评审和合并。
- `main` Branch Protection required checks 尚未自举。
- SQLite、Keychain、Operation Bus、SSH、数据库和 AI 均未开始。

### 测试与证据

- Actions Run `29312670714`：`quality`、`generated-and-security`、`windows-build` 全部成功。
- Actions Run `29313298880`：文档和状态同步后的三个 Job 全部成功。
- RED Run `29313493020`：`quality` 在 `Verify Go module and formatting` 失败；`generated-and-security` 与 `windows-build` 成功。
- GREEN Run `29313730163`：删除 `red_gate_probe.go` 后三个 Job 全部成功。
- `quality`：CI Policy tests、Go format/vet/test/race、5 个 Vitest、TypeScript、ESLint 和 Vite build 通过。
- `generated-and-security`：固定工具、实际 Workflow 自检、依赖政策、仓库边界、Binding/Go Module/pnpm Lockfile 漂移和完整历史 Gitleaks 通过。
- `windows-build`：Go tests、Wails production build、非空 `ItemAll.exe` 和 10 秒启动烟测通过。
- 诊断 Artifact 证明 Linux Wails 生成结果内容 Blob 与仓库一致，唯一差异是 `0644 → 0755` 文件模式；正式 Workflow 生成后恢复 `0644` 再执行严格 diff。

### 风险/阻塞

无外部阻塞。当前验收风险只剩合并后正确配置三个 required checks，并用状态同步 PR 验证保护规则。

### 本周下一步

1. 完成 PR #11 整分支审查并合并。
2. 配置 `main` required checks，并通过状态同步 PR 验证。
3. 状态同步完成后将 M0-003 更新为 `DONE`、M0-004 更新为 `READY`。
4. M0-003 关闭前不启动 M0-004 或业务能力。

## 2026-W28（第 1 周）

### 目标

基于空仓库事实建立可执行的项目管理、架构和验收体系。

### 已完成

- 读取仓库元数据、README 和初始 Commit。
- 确认仓库不存在 Go/Wails/React、Migration、CI 和测试。
- 将增强版规划作为不可变 Reference 保存。
- 建立 20 个稳定场景、50 项 Backlog、24 周计划和 6 个 ADR。
- 创建 [PR #1](https://github.com/vvitem/item_all/pull/1)，在本周末仍处于评审阶段。

### 未完成

- W28 周末 PR #1 尚未合并，`M0-001` 当时仍为 `IN_REVIEW`。
- 工程骨架尚未初始化。

### 测试与证据

- 初始基线 Commit：`b590887fea1e2c43e7831a48932b9be5d44ccfa3`。
- 55 个变更文件均非空，项目管理 Markdown 相对链接通过本地完整性检查。
- 参考原文 Git Blob SHA 与本地文件一致：`2cf547356f2810327fce721ce698d108112dc6c1`。

### 风险/阻塞

无外部阻塞；产品名称和 Keychain 方案待确认。

### 下周

完成 `M0-001` 合并和状态同步，启动 `M0-002` 设计。
