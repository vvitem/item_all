> 状态：维护中  
> 负责人：vvitem  
> 最后更新：2026-07-14  
> 基线 Commit：`620272272991ec47a39458e2dcb03a0ca7648e95`  
> 关联 Milestone：`M0-foundation`  
> 关联 Issue/PR：[Issue #7](https://github.com/vvitem/item_all/issues/7) · [PR #11](https://github.com/vvitem/item_all/pull/11) · [Issue #12](https://github.com/vvitem/item_all/issues/12)

# 当前项目状态

## 基线信息

- 基线分支：`main`
- 当前基线 Commit：`620272272991ec47a39458e2dcb03a0ca7648e95`
- 当前活动项：`M0-004`
- 最后更新时间：2026-07-14
- 当前 Milestone：`M0-foundation`
- 当前周次：第 2 周
- 项目状态：`DONE=3`、`READY=1`、`NOT_STARTED=46`

## 当前正在进行

- `M0-004`：SQLite WAL、Migration 和 Repository 骨架，状态 `READY`。
- 跟踪 Issue：[Issue #12](https://github.com/vvitem/item_all/issues/12)。
- 下一阶段先完成设计 Spec 与 Implementation Plan，批准前不写持久化实现。
- 设计需冻结 SQLite Driver、Migration 方案、数据目录和 Repository/Tx 边界。

## 最近完成

- `M0-003` 已通过 [PR #11](https://github.com/vvitem/item_all/pull/11) squash 合并，Merge Commit `620272272991ec47a39458e2dcb03a0ca7648e95`。
- 正式 CI 包含三个稳定 Job：`quality`、`generated-and-security`、`windows-build`。
- 最终 PR Head Run `29314687975` 三个 Job 全部成功。
- 受控门禁证明：RED Run `29313493020` 精确阻断未格式化 Go，GREEN Run `29313730163` 恢复全绿。
- 代码审查回归证明：RED Run `29314236649` 复现策略缺陷，GREEN Run `29314452432` 完成修复。
- `main` 已配置 Ruleset：PR 合并、分支最新、三个 required checks、禁止删除和 force push。

## 下一步任务

1. 合并本状态同步 PR，验证 `quality`、`generated-and-security`、`windows-build` 作为 required checks 生效。
2. 将 Issue #7 按 completed 关闭。
3. 编写 M0-004 SQLite/Repository 设计 Spec。
4. 设计批准后编写 Implementation Plan，再进入实现。

## 当前阻塞

无。

## 高风险事项

- SQLite Driver 选择会影响 CGO、Windows 构建、二进制体积和跨平台维护成本。
- Migration 必须 fail-closed、可重复执行，并具备损坏数据库与锁冲突测试。
- Repository 不得把 `*sql.DB` 暴露给 UI/Wails Binding。
- required check 名称不得随意修改，否则 Ruleset 会失效。

## 待确认决策

- SQLite Driver：`modernc.org/sqlite` 或 `github.com/mattn/go-sqlite3`。
- Migration：内嵌 SQL 最小 runner 或成熟 migration 库。
- Windows Credential Manager 封装方案。
- 手工 SQL 写操作是否进入首版。
- Windows 稳定版代码签名方案。

## 范围变化

无。M0-004 仅建立本地持久化底座，不包含 CredentialRef、Operation Bus、AuditEvent 业务模型、SSH、数据库运维、AI、云同步或遥测。

## 测试状态

- CI Policy、Go、Frontend、生成漂移、完整历史 Gitleaks 和 Windows Wails build 均已在 PR #11 验证。
- 最终 Head GREEN：Actions Run `29314687975`。
- 本状态同步 PR 必须通过三个 required checks 后才能合并。

## 相关代码路径

- `.github/workflows/ci.yml`
- `scripts/ci/**`
- `docs/development/ci.md`
- 下一阶段建议：`internal/storage`、`internal/storage/migrations`、`internal/repository`

## 相关 Issue / PR / Commit

- [Issue #7](https://github.com/vvitem/item_all/issues/7)：M0-003 CI 门禁，待本状态同步 PR 合并后关闭。
- [PR #11](https://github.com/vvitem/item_all/pull/11)：M0-003 实现，已合并。
- Merge Commit：`620272272991ec47a39458e2dcb03a0ca7648e95`。
- [Issue #12](https://github.com/vvitem/item_all/issues/12)：M0-004，`READY`。

## 本周可演示结果

任意 Pull Request 必须通过 Go/Frontend 质量、生成与安全边界、Windows Wails 构建三个 required checks 后才能合并；下一阶段开始建立本地 SQLite 持久化底座。