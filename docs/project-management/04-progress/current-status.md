> 状态：维护中  
> 负责人：vvitem  
> 最后更新：2026-07-14  
> 基线 Commit：`cf25d4e40f34e0c2e835a2cc04b3eeea4d011ad7`  
> 关联 Milestone：`M0-foundation`  
> 关联 Issue/PR：[Issue #7](https://github.com/vvitem/item_all/issues/7) · [PR #11](https://github.com/vvitem/item_all/pull/11)

# 当前项目状态

## 基线信息

- 基线分支：`main`
- 当前基线 Commit：`cf25d4e40f34e0c2e835a2cc04b3eeea4d011ad7`
- 实现分支：`feat/m0-003-ci-gates`
- 最后更新时间：2026-07-14
- 当前 Milestone：`M0-foundation`
- 当前周次：第 2 周
- 项目状态：`DONE=2`、`IN_REVIEW=1`、`NOT_STARTED=47`

## 本阶段目标

建立可构建的 Go/Wails/React 工程骨架、正式 CI 门禁、SQLite WAL、OS Keychain、Operation Bus、append-only Audit，以及 SSH 资产连接测试垂直切片。

## 当前正在进行

- `M0-003`：正式分层 CI 门禁已实现，状态 `IN_REVIEW`。
- 实现 PR：[PR #11](https://github.com/vvitem/item_all/pull/11)。
- 三个稳定 Job：`quality`、`generated-and-security`、`windows-build`。
- GREEN 基线：Actions Run `29312670714` 三个 Job 全部成功。
- RED→GREEN 证明：RED Run `29313493020`、GREEN Run `29313730163`。
- 当前剩余工作：整分支代码审查、最终全量 GREEN、合并与 Branch Protection 自举。
- 在 M0-003 完成前不开始 SSH、数据库或 AI 功能。

## 最近完成

- M0-003 设计通过 [PR #9](https://github.com/vvitem/item_all/pull/9) 合并，Merge Commit `3cc4039fcf88882f54be55029541b2f9b638c1d2`。
- M0-003 实现计划通过 [PR #10](https://github.com/vvitem/item_all/pull/10) 合并，Merge Commit `cf25d4e40f34e0c2e835a2cc04b3eeea4d011ad7`。
- 建立只读权限、固定 Action SHA、固定 Runner 和固定工具链的三 Job Workflow。
- 建立 `scripts/ci` Go 检查器，覆盖依赖政策、Workflow 安全和仓库范围边界。
- 建立完整历史 Gitleaks、Binding/Go Module/pnpm Lockfile 漂移检查。
- 建立 Windows Wails production build、EXE 存在性和启动烟测。
- 修复 Linux Wails Binding 生成造成的文件模式噪声：生成后恢复 `0644`，内容漂移仍严格阻断。
- 删除所有临时诊断 Workflow，正式分支只保留 `.github/workflows/ci.yml`。
- 受控提交未格式化 Go Fixture 后，RED Run 仅由 `quality` 的格式门禁失败；安全和 Windows Job 均成功。
- 删除 Fixture 后，GREEN Run `29313730163` 三个 Job 全部恢复成功。

## 下一步任务

1. 完成整分支规范与质量复核，将 PR #11 转为 Ready for Review。
2. 确认最终 Head 对应的三个 Job 全绿，且无 `red_gate_probe.go` 或临时 Workflow。
3. 合并后配置 `quality`、`generated-and-security`、`windows-build` 为 `main` required checks。
4. 通过独立状态同步 PR 将 `M0-003` 更新为 `DONE`；此前 `M0-004` 保持 `NOT_STARTED`。

## 当前阻塞

无实现阻塞。当前仅剩 M0-003 的最终评审、合并和 Branch Protection 自举。

## 高风险事项

- required check 名称进入 Branch Protection 后不得随意修改。
- Secret Scan 必须保持 `fetch-depth: 0`，否则无法覆盖完整历史。
- Wails Linux 生成文件模式差异必须先规范化，再执行内容与模式漂移检查。
- Operation Bus 尚未实现，后续运维能力不得绕过其设计边界。
- Windows Keychain/PTTY/WebView 的后续兼容风险仍需独立验证。

## 待确认决策

- Windows Credential Manager 封装方案。
- 手工 SQL 写操作是否进入首版。
- Windows 稳定版代码签名方案。

## 范围变化

无。M0-003 只建立工程质量和安全门禁，没有引入 SSH、数据库、SQLite 业务表、AI、Operation Bus、凭据、云服务、遥测或发布能力。

## 测试状态

- CI Policy：依赖、Workflow 与仓库边界单元测试通过。
- `quality`：Go format/vet/test/race、5 个 Vitest、TypeScript、ESLint、Vite build 全部通过。
- `generated-and-security`：固定工具安装、实际 Workflow 自检、三类漂移、完整历史 Gitleaks 全部通过。
- `windows-build`：Go tests、Wails build、非空 `ItemAll.exe` 和启动烟测通过。
- GREEN 基线：Actions Run `29312670714`。
- 文档同步 GREEN：Actions Run `29313298880`。
- RED Run `29313493020`：只有 `quality → Verify Go module and formatting` 失败；另外两个 Job 成功。
- GREEN Run `29313730163`：删除受控 Fixture 后三个 Job 全部成功。

## 相关代码路径

- `.github/workflows/ci.yml`
- `scripts/ci/*.go`、`scripts/ci/*_test.go`
- `scripts/ci/smoke-windows.ps1`
- `docs/development/ci.md`
- `README.md`

## 相关 Issue / PR / Commit

- [Issue #7](https://github.com/vvitem/item_all/issues/7)：M0-003 CI 门禁，Open。
- [PR #9](https://github.com/vvitem/item_all/pull/9)：M0-003 设计，已合并。
- [PR #10](https://github.com/vvitem/item_all/pull/10)：M0-003 实现计划，已合并。
- [PR #11](https://github.com/vvitem/item_all/pull/11)：M0-003 实现，当前 `IN_REVIEW`。
- GREEN 基线：Actions Run `29312670714`。
- RED→GREEN：Actions Run `29313493020` → `29313730163`。

## 本周可演示结果

任意 Pull Request 可并行获得 Go/Frontend 质量反馈、生成与安全边界检查以及 Windows Wails 构建事实；受控失败已证明未格式化 Go 会阻断 `quality`，且恢复后三个 Job 全部转绿。
