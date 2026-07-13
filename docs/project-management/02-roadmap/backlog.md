> 状态：维护中  
> 负责人：vvitem  
> 最后更新：2026-07-13  
> 基线 Commit：`dbb104e3b457d0ef83ca563ff2ae6106c74b43a1`  
> 关联 Milestone：`M0-foundation → M4-beta`  
> 关联 Issue/PR：[PR #6](https://github.com/vvitem/item_all/pull/6) · [Issue #7](https://github.com/vvitem/item_all/issues/7)

# Backlog

状态只允许：`NOT_STARTED`、`READY`、`IN_PROGRESS`、`BLOCKED`、`IN_REVIEW`、`DONE`、`DEFERRED`、`CANCELLED`。当前共 50 项：`DONE=2`，`READY=1`，`NOT_STARTED=47`。

| ID | Milestone | Epic | 模块 | 任务 | 优先级 | 状态 | 依赖 | 验收标准 | 代码路径 | 测试 | Issue/PR | 阻塞原因 | 下一步 |
|---|---|---|---|---|---|---|---|---|---|---|---|---|---|
| M0-001 | M0-foundation | Project Foundation | docs | 建立项目管理、ADR、进度和交付文档体系 | P0 | DONE | 无 | 所有文件非空、状态基于真实仓库、创建并合并 PR | docs/project-management/** | 链接与完整性检查 | [PR #1](https://github.com/vvitem/item_all/pull/1) / `be37cfb` | 无 | 已完成，持续维护 |
| M0-002 | M0-foundation | Project Foundation | build | 初始化 Go/Wails/React 单仓 | P0 | DONE | M0-001 | Windows 本地可启动空壳，版本信息可见 | `main.go`, `app.go`, `internal/buildinfo`, `frontend/`, `build/` | Go tests/race、Vitest/typecheck/lint/build、Windows build/dev smoke | [Issue #2](https://github.com/vvitem/item_all/issues/2) / [PR #6](https://github.com/vvitem/item_all/pull/6) / `dbb104e` | 无 | 已完成，持续维护 |
| M0-003 | M0-foundation | Quality | ci | 建立 Go/前端 lint、test、Windows build、secret scan | P0 | READY | M0-002 | PR 必须通过全部门禁 | 建议：.github/workflows/ci.yml | 故意失败验证门禁 | [Issue #7](https://github.com/vvitem/item_all/issues/7) | 无 | 编写设计 Spec 与实施计划 |
| M0-004 | M0-foundation | Storage | database | 建立 SQLite WAL、Migration 和 Repository 骨架 | P0 | NOT_STARTED | M0-002 | 首次启动迁移成功；重复运行幂等 | 建议：internal/storage, migrations/ | migration integration test | 待创建 | 无 | 确定 migration 库 |
| M0-005 | M0-foundation | Identity | security | 实现 workspace/actor/device 本地身份模型 | P1 | NOT_STARTED | M0-004 | 默认本地身份可追踪且不需要账号 | 建议：internal/identity | repository tests | 待创建 | 无 | 定义 ID 策略 |
| M0-006 | M0-foundation | Credential | security | 实现 CredentialRef 与 OS Keychain 接口 | P0 | NOT_STARTED | M0-004 | SQLite/日志无明文；Keychain 失败 fail-closed | 建议：internal/security/credential | fake keychain + Windows test | 待创建 | 无 | 完成 Windows PoC |
| M0-007 | M0-foundation | Operation Bus | operation | 实现 OperationRequest、流水线和假适配器 | P0 | NOT_STARTED | M0-004,M0-005 | UI 服务不能直接调用适配器；请求有 TraceID | 建议：internal/operation | pipeline unit tests | 待创建 | 无 | 先定义接口 |
| M0-008 | M0-foundation | Audit | audit | 实现 append-only AuditEvent | P0 | NOT_STARTED | M0-004,M0-007 | requested/completed/failed 事件完整且不可更新 | 建议：internal/audit | repository + tamper tests | 待创建 | 无 | 确定 hash chain 是否首版启用 |
| M0-009 | M0-foundation | Asset | asset | 实现 SSH Asset CRUD 与校验 | P0 | NOT_STARTED | M0-004,M0-006 | 增删改查、标签、环境、凭据引用可用 | 建议：internal/asset, frontend/src/features/assets | service/UI tests | 待创建 | 无 | 定义 DTO |
| M0-010 | M0-foundation | SSH | connection | 完成 SSH 连接测试、指纹确认和错误分类 | P0 | NOT_STARTED | M0-006,M0-007,M0-009 | 首次确认、变化阻断、认证/网络错误区分 | 建议：internal/connection/ssh | real-server integration | 待创建 | 无 | 准备测试 SSH 主机 |
| M1-001 | M1-ssh-diagnosis | Terminal | ssh | 实现 xterm 终端基础会话 | P0 | NOT_STARTED | M0-010 | 输入输出、resize、关闭和错误态可用 | 建议：internal/ssh/session, frontend/.../terminal | Windows E2E | 待创建 | 无 | 先验证 Windows PTY/WebView |
| M1-002 | M1-ssh-diagnosis | Terminal | ssh | 实现断线检测、重连和状态机 | P0 | NOT_STARTED | M1-001 | 断网后状态不假在线；手工重连可恢复 | 建议：internal/ssh/session | network fault test | 待创建 | 无 | 定义 session state |
| M1-003 | M1-ssh-diagnosis | SFTP | sftp | 实现目录浏览、上传和下载 | P0 | NOT_STARTED | M0-010 | 限制并发、进度、取消和错误可见 | 建议：internal/sftp, frontend/.../sftp | integration + E2E | 待创建 | 无 | 定义传输队列 |
| M1-004 | M1-ssh-diagnosis | Import | asset | 实现 SSH Config 导入 | P1 | NOT_STARTED | M0-009 | 支持 Host/Hostname/User/Port/IdentityFile 引用；不复制私钥 | 建议：internal/asset/importer | fixture tests | 待创建 | 无 | 定义冲突策略 |
| M1-005 | M1-ssh-diagnosis | SSH Tools | ssh | 实现 system.info/load/memory/time | P0 | NOT_STARTED | M0-007,M0-010 | SSH-01/02/10 通过，固定模板且可取消 | 建议：internal/ssh/tools | tool + integration tests | 待创建 | 无 | 实现统一 runner |
| M1-006 | M1-ssh-diagnosis | SSH Tools | ssh | 实现 disk.usage/file.top_size | P0 | NOT_STARTED | M1-005 | SSH-03/04 通过，路径和输出受限 | 建议：internal/ssh/tools | path boundary tests | 待创建 | 无 | 定义允许根路径 |
| M1-007 | M1-ssh-diagnosis | SSH Tools | ssh | 实现 service.status/process.list | P0 | NOT_STARTED | M1-005 | SSH-05/06 通过，命令参数脱敏 | 建议：internal/ssh/tools | parser + redaction tests | 待创建 | 无 | 覆盖非 systemd |
| M1-008 | M1-ssh-diagnosis | SSH Tools | ssh | 实现 network.listen/log.tail/log.search | P0 | NOT_STARTED | M1-005 | SSH-07/08/09 通过，日志视为不可信 | 建议：internal/ssh/tools | injection + limit tests | 待创建 | 无 | 定义 log allowlist |
| M1-009 | M1-ssh-diagnosis | Task UI | frontend | 实现只读诊断任务时间线和取消 | P1 | NOT_STARTED | M1-005..008 | 步骤、输出、Evidence 和失败位置可见 | 建议：frontend/src/features/tasks | component + E2E | 待创建 | 无 | 先用固定 Plan |
| M1-010 | M1-ssh-diagnosis | Hardening | quality | 完成 SSH 场景真机回归和诊断包 | P0 | NOT_STARTED | M1-001..009 | 至少8/10场景通过；秘密扫描为0 | 建议：e2e/ssh, internal/diagnostics | Windows/real host report | 待创建 | 无 | 准备回归矩阵 |
| M2-001 | M2-database | Database | database | 实现统一数据库连接配置与池 | P0 | NOT_STARTED | M0-006,M0-007 | MySQL/PG 连接、超时、TLS 错误分类 | 建议：internal/database/connection | container integration | 待创建 | 无 | 定义 Driver 接口 |
| M2-002 | M2-database | Tunnel | connection | 实现数据库 SSH Tunnel | P0 | NOT_STARTED | M0-010,M2-001 | Tunnel 生命周期随连接释放且不泄漏端口 | 建议：internal/connection/tunnel | fault integration | 待创建 | 无 | 定义 lease |
| M2-003 | M2-database | Metadata | database | 实现 schemas/tables/columns/indexes | P0 | NOT_STARTED | M2-001 | DB-02/03/04 元数据跨方言可用 | 建议：internal/database/tools | MySQL+PG tests | 待创建 | 无 | 建立 dialect 层 |
| M2-004 | M2-database | SQL UI | frontend | 实现 SQL 编辑器和结果表格 | P0 | NOT_STARTED | M2-001 | 可编辑、执行只读查询、分页/截断可见 | 建议：frontend/src/features/database | component + E2E | 待创建 | 无 | 选择 Monaco/虚拟表格 |
| M2-005 | M2-database | Readonly | policy | 实现 SQL 分类和只读事务双层约束 | P0 | NOT_STARTED | M2-001,M0-007 | DML/DDL/多语句/危险 CTE 均拒绝 | 建议：internal/database/sqlpolicy | corpus + DB tests | 待创建 | 无 | 选择 parser |
| M2-006 | M2-database | DB Tools | database | 实现 connection_stats/sessions/locks | P0 | NOT_STARTED | M2-003,M2-005 | DB-01/06/07/08 通过 | 建议：internal/database/tools | dialect tests | 待创建 | 无 | 处理权限降级 |
| M2-007 | M2-database | DB Tools | database | 实现 table_stats/query_readonly | P0 | NOT_STARTED | M2-003,M2-005 | DB-09/10 通过，行/字节/时间限制有效 | 建议：internal/database/tools | limit/cancel tests | 待创建 | 无 | 定义结果序列化 |
| M2-008 | M2-database | Explain | database | 实现 db.explain 和计划规范化 | P0 | NOT_STARTED | M2-005 | DB-05 通过；默认不使用 ANALYZE | 建议：internal/database/explain | fixture + integration | 待创建 | 无 | 定义 normalized plan |
| M2-009 | M2-database | History | database | 实现查询历史和敏感值处理 | P1 | NOT_STARTED | M2-004,M2-005 | 只保存 SQL hash/脱敏文本，可清理 | 建议：internal/database/history | redaction tests | 待创建 | 无 | 确定保留周期 |
| M2-010 | M2-database | Hardening | quality | 完成数据库场景回归 | P0 | NOT_STARTED | M2-001..009 | 至少8/10场景通过；写操作0 | 建议：e2e/database | MySQL+PG report | 待创建 | 无 | 建立 Docker fixtures |
| M3-001 | M3-safe-ai | Task Engine | task | 实现 Task/TaskStep 持久状态机 | P0 | NOT_STARTED | M0-004,M0-007 | 状态转换受控、事务一致、重启可恢复 | 建议：internal/task | state/property tests | 待创建 | 无 | 先写 transition table |
| M3-002 | M3-safe-ai | Task Engine | task | 实现超时、取消、暂停和重试 | P0 | NOT_STARTED | M3-001 | 取消传播到适配器；重试幂等 | 建议：internal/task/executor | fault tests | 待创建 | 无 | 定义 retry policy |
| M3-003 | M3-safe-ai | Evidence | evidence | 实现 Evidence Contract 和对象存储 | P0 | NOT_STARTED | M0-008,M3-001 | 每步声明预期证据；大对象不进 SQLite | 建议：internal/evidence | storage/redaction tests | 待创建 | 无 | 确定文件布局 |
| M3-004 | M3-safe-ai | AI Provider | ai | 实现两类 Provider 和模型配置 | P1 | NOT_STARTED | M0-006 | Key 不进入日志；超时/限流错误稳定 | 建议：internal/ai/provider | mock + live smoke | 待创建 | 无 | 确认 Provider |
| M3-005 | M3-safe-ai | AI Plan | ai | 定义最多8步 JSON Schema Plan | P0 | NOT_STARTED | M3-001,M3-004 | 非法工具/资产/参数无法通过验证 | 建议：internal/ai/plan | schema/fuzz tests | 待创建 | 无 | 冻结 v1 schema |
| M3-006 | M3-safe-ai | AI Plan | security | 实现 Plan Scope 锁定和风险重算 | P0 | NOT_STARTED | M3-005,M0-007 | 执行不能扩大计划；风险由本地计算 | 建议：internal/policy | tamper tests | 待创建 | 无 | 定义 scope hash |
| M3-007 | M3-safe-ai | Prompt Security | security | 实现不可信内容隔离与注入测试集 | P0 | NOT_STARTED | M3-004,M3-005 | 日志/SQL结果不能新增工具或指令 | 建议：internal/ai/context | adversarial tests | 待创建 | 无 | 建立攻击语料 |
| M3-008 | M3-safe-ai | Report | ai | 基于 Evidence 生成验证和报告 | P1 | NOT_STARTED | M3-003,M3-004 | 报告引用 Evidence ID，不虚构成功 | 建议：internal/ai/report | golden tests | 待创建 | 无 | 定义 citation format |
| M3-009 | M3-safe-ai | Recovery | task | 实现应用重启恢复和明确终止 | P0 | NOT_STARTED | M3-001..003 | 无 UNKNOWN 任务；在途步骤按规则恢复/失败 | 建议：internal/task/recovery | crash tests | 待创建 | 无 | 定义 lease/heartbeat |
| M3-010 | M3-safe-ai | End-to-End | quality | 完成20场景受约束 AI 回归 | P0 | NOT_STARTED | M1-010,M2-010,M3-001..009 | ≥70%无需改参数；AI写操作0 | 建议：e2e/ai | scenario report | 待创建 | 无 | 锁定测试模型 |
| M4-001 | M4-beta | Runbook | runbook | 实现 Runbook/Version 数据模型 | P0 | NOT_STARTED | M3-005 | 成功 Plan 可保存、版本化、参数化 | 建议：internal/runbook | repository/schema tests | 待创建 | 无 | 定义变量类型 |
| M4-002 | M4-beta | Runbook | runbook | 实现 Runbook 执行仍经过 Operation Bus | P0 | NOT_STARTED | M4-001,M0-007 | 不可绕过 Policy/Audit/Evidence | 建议：internal/runbook/executor | bypass tests | 待创建 | 无 | 复用 Task Engine |
| M4-003 | M4-beta | Templates | product | 内置10–15个诊断模板 | P1 | NOT_STARTED | M3-010,M4-001 | 模板覆盖高频场景且有版本 | 建议：assets/runbooks | golden scenario tests | 待创建 | 无 | 按用户反馈排序 |
| M4-004 | M4-beta | Onboarding | frontend | 实现首次引导和示例诊断 | P0 | NOT_STARTED | M1/M2/M3 | 新用户10分钟内完成首个可信任务 | 建议：frontend/src/features/onboarding | usability E2E | 待创建 | 无 | 招募5名测试用户 |
| M4-005 | M4-beta | Windows | quality | 建立 Windows 核心 E2E | P0 | NOT_STARTED | M1-010,M2-010,M3-010 | 安装、Keychain、SSH、DB、任务恢复通过 | 建议：e2e/windows | GitHub runner+真机 | 待创建 | 无 | 准备签名/无签名路径 |
| M4-006 | M4-beta | Security | security | 完成发布前威胁复核和秘密扫描 | P0 | NOT_STARTED | 全部P0 | 无明文秘密；极高风险均关闭或拒绝发布 | docs/.../security-release-checklist.md | security report | 待创建 | 无 | 安排外部复核可选 |
| M4-007 | M4-beta | Release | build | 建立版本、安装包、校验和和发布流程 | P0 | NOT_STARTED | M0-003,M4-005 | 可重复构建 v0.1.0-beta，产物有 SHA256 | 建议：.github/workflows/release.yml | clean-machine install | 待创建 | 无 | 定义版本注入 |
| M4-008 | M4-beta | Telemetry | privacy | 实现默认关闭的诊断包/可选遥测 | P1 | NOT_STARTED | M3-003 | 用户可预览；默认不上传命令/结果 | 建议：internal/diagnostics | privacy tests | 待创建 | 无 | 定义字段白名单 |
| M4-009 | M4-beta | Community | docs | 完善 CONTRIBUTING、威胁模型和 Good First Issue | P1 | NOT_STARTED | 稳定代码结构 | 陌生贡献者可本地构建测试 | 建议：CONTRIBUTING.md,.github | fresh clone test | 待创建 | 无 | 记录环境依赖 |
| M4-010 | M4-beta | Beta | release | 完成种子用户测试并发布 Beta | P0 | NOT_STARTED | M4-001..009 | 20安装/10连续使用两周/阻断问题关闭 | release notes | beta report | 待创建 | 无 | 建立反馈表 |

## DONE 规则

任务只有在功能/文档完成、正常与失败路径验证、测试通过、边界未绕过、文档同步且有 Commit/PR/报告证据时才能标记 `DONE`。