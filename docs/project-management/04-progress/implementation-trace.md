> 状态：维护中  
> 负责人：vvitem  
> 最后更新：2026-07-11  
> 基线 Commit：`b590887fea1e2c43e7831a48932b9be5d44ccfa3`  
> 关联 Milestone：`M0-foundation → M4-beta`  
> 关联 Issue/PR：待创建

# 实现追踪

| Requirement/Scenario | Task | Design Doc | Code Path | Test Path | Issue/PR | Status |
|---|---|---|---|---|---|---|
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

代码与测试路径均为建议路径，直到 PR 合并后才能移除“建议”标记。
