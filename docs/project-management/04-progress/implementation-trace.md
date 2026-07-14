> 状态：维护中  
> 负责人：vvitem  
> 最后更新：2026-07-14  
> 基线 Commit：`620272272991ec47a39458e2dcb03a0ca7648e95`  
> 关联 Milestone：`M0-foundation → M4-beta`  
> 关联 Issue/PR：[PR #11](https://github.com/vvitem/item_all/pull/11) · [Issue #12](https://github.com/vvitem/item_all/issues/12)

# 实现追踪

| Requirement/Scenario | Task | Design Doc | Code Path | Test Path | Issue/PR | Status |
|---|---|---|---|---|---|---|
| Project Management Baseline | M0-001 | `docs/project-management/**` | `README.md`, `.github/**`, `docs/project-management/**` | 文档非空、相对链接与 Reference SHA 校验 | [PR #1](https://github.com/vvitem/item_all/pull/1) / `be37cfb` | DONE |
| ItemAll Scaffold | M0-002 | `docs/superpowers/specs/2026-07-13-m0-002-project-scaffold-design.md` | `main.go`, `app.go`, `internal/buildinfo`, `frontend/src`, `build`, `wails.json` | `app_test.go`, `internal/buildinfo/info_test.go`, `frontend/src/App.test.tsx`；Actions `29242978948`、`29243245959`、`29244060195`、`29244941752` | [Issue #2](https://github.com/vvitem/item_all/issues/2) / [PR #6](https://github.com/vvitem/item_all/pull/6) / `dbb104e` | DONE |
| CI Quality Gates | M0-003 | `docs/superpowers/specs/2026-07-13-m0-003-ci-gates-design.md` | `.github/workflows/ci.yml`, `scripts/ci`, `docs/development/ci.md` | Policy tests；Go/Frontend；Binding/lock drift；完整历史 Gitleaks；Windows Wails build/smoke；最终 GREEN `29314687975`；RED→GREEN `29313493020`→`29313730163`；审查回归 `29314236649`→`29314452432` | [Issue #7](https://github.com/vvitem/item_all/issues/7) / [PR #11](https://github.com/vvitem/item_all/pull/11) / `6202722` | DONE |
| Local Storage Foundation | M0-004 | 待编写 SQLite/Repository 设计 Spec | 建议：`internal/storage`, `internal/repository` | Migration 幂等、PRAGMA、并发读写、事务回滚、失败启动 | [Issue #12](https://github.com/vvitem/item_all/issues/12) | READY |
| SSH-01/02/10 | M1-005 | tool-contracts.md | 建议：internal/ssh/tools | 建议：internal/ssh/tools/*_test.go | 待创建 | NOT_STARTED |
| SSH-03/04 | M1-006 | tool-contracts.md | 建议：internal/ssh/tools | 建议：path_boundary_test.go | 待创建 | NOT_STARTED |
| SSH-05/06 | M1-007 | tool-contracts.md | 建议：internal/ssh/tools | 建议：service_process_test.go | 待创建 | NOT_STARTED |
| SSH-07/08/09 | M1-008 | tool-contracts.md | 建议：internal/ssh/tools | 建议：network_log_test.go | 待创建 | NOT_STARTED |
| DB-01/06/07/08 | M2-006 | tool-contracts.md | 建议：internal/database/tools | 建议：sessions_locks_test.go | 待创建 | NOT_STARTED |
| DB-02/03/04 | M2-003 | tool-contracts.md | 建议：internal/database/tools | 建议：metadata_test.go | 待创建 | NOT_STARTED |
| DB-05 | M2-008 | tool-contracts.md | 建议：internal/database/explain | 建议：explain_test.go | 待创建 | NOT_STARTED |
| DB-09/10 | M2-007 | tool-contracts.md | 建议：internal/database/tools | 建议：readonly_limit_test.go | 待创建 | NOT_STARTED |
| Operation Bus | M0-007 | operation-bus.md | 建议：internal/operation | 建议：pipeline_test.go | 待创建 | NOT_STARTED |
| Task Recovery | M3-001/009 | task-engine.md | 建议：internal/task | 建议：recovery_test.go | 待创建 | NOT_STARTED |
| Credential Safety | M0-006 | security-and-threat-model.md | 建议：internal/security | 建议：credential_leak_test.go | 待创建 | NOT_STARTED |
| Windows Beta | M4-005/007 | testing-strategy.md | 建议：.github/workflows,e2e/windows | 建议：e2e/windows | 待创建 | NOT_STARTED |

未合并实现的路径可以在追踪表中记录，但只有合并并具有验证证据后才能更新为 `DONE`。M0-003 已以 Merge Commit `6202722` 结项，M0-004 当前仅进入设计准备。