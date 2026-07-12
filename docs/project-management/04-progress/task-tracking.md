> 状态：维护中  
> 负责人：vvitem  
> 最后更新：2026-07-13  
> 基线 Commit：`be37cfbb2bdca5a05c3303e9ea208992a8bf1721`  
> 关联 Milestone：`M0-foundation → M4-beta`  
> 关联 Issue/PR：[PR #1](https://github.com/vvitem/item_all/pull/1) · [Issue #2](https://github.com/vvitem/item_all/issues/2)

# 任务追踪

本文件只展示当前活动项和全部 P0，完整 50 项见 [backlog.md](../02-roadmap/backlog.md)。

| ID | Milestone | 任务 | 优先级 | 状态 | 依赖 | 验收 | Issue/PR | 下一步 |
|---|---|---|---|---|---|---|---|---|
| M0-001 | M0-foundation | 建立项目管理、ADR、进度和交付文档体系 | P0 | DONE | 无 | 所有文件非空、状态基于真实仓库、PR 已合并 | [PR #1](https://github.com/vvitem/item_all/pull/1) / `be37cfb` | 持续维护 |
| M0-002 | M0-foundation | 初始化 Go/Wails/React 单仓 | P0 | READY | M0-001 | Windows 本地可启动空壳，版本信息可见 | [Issue #2](https://github.com/vvitem/item_all/issues/2) | 确认名称和 Module Path，完成设计 Spec |
| M0-003 | M0-foundation | 建立 Go/前端 lint、test、Windows build、secret scan | P0 | NOT_STARTED | M0-002 | PR 必须通过全部门禁 | 待创建 | 定义工具版本 |
| M0-004 | M0-foundation | 建立 SQLite WAL、Migration 和 Repository 骨架 | P0 | NOT_STARTED | M0-002 | 首次启动迁移成功；重复运行幂等 | 待创建 | 确定 migration 库 |
| M0-006 | M0-foundation | 实现 CredentialRef 与 OS Keychain 接口 | P0 | NOT_STARTED | M0-004 | SQLite/日志无明文；Keychain 失败 fail-closed | 待创建 | 完成 Windows PoC |
| M0-007 | M0-foundation | 实现 OperationRequest、流水线和假适配器 | P0 | NOT_STARTED | M0-004,M0-005 | UI 服务不能直接调用适配器；请求有 TraceID | 待创建 | 先定义接口 |
| M0-008 | M0-foundation | 实现 append-only AuditEvent | P0 | NOT_STARTED | M0-004,M0-007 | requested/completed/failed 事件完整且不可更新 | 待创建 | 确定 hash chain 是否首版启用 |
| M0-009 | M0-foundation | 实现 SSH Asset CRUD 与校验 | P0 | NOT_STARTED | M0-004,M0-006 | 增删改查、标签、环境、凭据引用可用 | 待创建 | 定义 DTO |
| M0-010 | M0-foundation | 完成 SSH 连接测试、指纹确认和错误分类 | P0 | NOT_STARTED | M0-006,M0-007,M0-009 | 首次确认、变化阻断、认证/网络错误区分 | 待创建 | 准备测试 SSH 主机 |
| M1-001 | M1-ssh-diagnosis | 实现 xterm 终端基础会话 | P0 | NOT_STARTED | M0-010 | 输入输出、resize、关闭和错误态可用 | 待创建 | 先验证 Windows PTY/WebView |
| M1-002 | M1-ssh-diagnosis | 实现断线检测、重连和状态机 | P0 | NOT_STARTED | M1-001 | 断网后状态不假在线；手工重连可恢复 | 待创建 | 定义 session state |
| M1-003 | M1-ssh-diagnosis | 实现目录浏览、上传和下载 | P0 | NOT_STARTED | M0-010 | 限制并发、进度、取消和错误可见 | 待创建 | 定义传输队列 |
| M1-005 | M1-ssh-diagnosis | 实现 system.info/load/memory/time | P0 | NOT_STARTED | M0-007,M0-010 | SSH-01/02/10 通过，固定模板且可取消 | 待创建 | 实现统一 runner |
| M1-006 | M1-ssh-diagnosis | 实现 disk.usage/file.top_size | P0 | NOT_STARTED | M1-005 | SSH-03/04 通过，路径和输出受限 | 待创建 | 定义允许根路径 |
| M1-007 | M1-ssh-diagnosis | 实现 service.status/process.list | P0 | NOT_STARTED | M1-005 | SSH-05/06 通过，命令参数脱敏 | 待创建 | 覆盖非 systemd |
| M1-008 | M1-ssh-diagnosis | 实现 network.listen/log.tail/log.search | P0 | NOT_STARTED | M1-005 | SSH-07/08/09 通过，日志视为不可信 | 待创建 | 定义 log allowlist |
| M1-010 | M1-ssh-diagnosis | 完成 SSH 场景真机回归和诊断包 | P0 | NOT_STARTED | M1-001..009 | 至少8/10场景通过；秘密扫描为0 | 待创建 | 准备回归矩阵 |
| M2-001 | M2-database | 实现统一数据库连接配置与池 | P0 | NOT_STARTED | M0-006,M0-007 | MySQL/PG 连接、超时、TLS 错误分类 | 待创建 | 定义 Driver 接口 |
| M2-002 | M2-database | 实现数据库 SSH Tunnel | P0 | NOT_STARTED | M0-010,M2-001 | Tunnel 生命周期随连接释放且不泄漏端口 | 待创建 | 定义 lease |
| M2-003 | M2-database | 实现 schemas/tables/columns/indexes | P0 | NOT_STARTED | M2-001 | DB-02/03/04 元数据跨方言可用 | 待创建 | 建立 dialect 层 |
| M2-004 | M2-database | 实现 SQL 编辑器和结果表格 | P0 | NOT_STARTED | M2-001 | 可编辑、执行只读查询、分页/截断可见 | 待创建 | 选择 Monaco/虚拟表格 |
| M2-005 | M2-database | 实现 SQL 分类和只读事务双层约束 | P0 | NOT_STARTED | M2-001,M0-007 | DML/DDL/多语句/危险 CTE 均拒绝 | 待创建 | 选择 parser |
| M2-006 | M2-database | 实现 connection_stats/sessions/locks | P0 | NOT_STARTED | M2-003,M2-005 | DB-01/06/07/08 通过 | 待创建 | 处理权限降级 |
| M2-007 | M2-database | 实现 table_stats/query_readonly | P0 | NOT_STARTED | M2-003,M2-005 | DB-09/10 通过，行/字节/时间限制有效 | 待创建 | 定义结果序列化 |
| M2-008 | M2-database | 实现 db.explain 和计划规范化 | P0 | NOT_STARTED | M2-005 | DB-05 通过；默认不使用 ANALYZE | 待创建 | 定义 normalized plan |
| M2-010 | M2-database | 完成数据库场景回归 | P0 | NOT_STARTED | M2-001..009 | 至少8/10场景通过；写操作0 | 待创建 | 建立 Docker fixtures |
| M3-001 | M3-safe-ai | 实现 Task/TaskStep 持久状态机 | P0 | NOT_STARTED | M0-004,M0-007 | 状态转换受控、事务一致、重启可恢复 | 待创建 | 先写 transition table |
| M3-002 | M3-safe-ai | 实现超时、取消、暂停和重试 | P0 | NOT_STARTED | M3-001 | 取消传播到适配器；重试幂等 | 待创建 | 定义 retry policy |
| M3-003 | M3-safe-ai | 实现 Evidence Contract 和对象存储 | P0 | NOT_STARTED | M0-008,M3-001 | 每步声明预期证据；大对象不进 SQLite | 待创建 | 确定文件布局 |
| M3-005 | M3-safe-ai | 定义最多8步 JSON Schema Plan | P0 | NOT_STARTED | M3-001,M3-004 | 非法工具/资产/参数无法通过验证 | 待创建 | 冻结 v1 schema |
| M3-006 | M3-safe-ai | 实现 Plan Scope 锁定和风险重算 | P0 | NOT_STARTED | M3-005,M0-007 | 执行不能扩大计划；风险由本地计算 | 待创建 | 定义 scope hash |
| M3-007 | M3-safe-ai | 实现不可信内容隔离与注入测试集 | P0 | NOT_STARTED | M3-004,M3-005 | 日志/SQL结果不能新增工具或指令 | 待创建 | 建立攻击语料 |
| M3-009 | M3-safe-ai | 实现应用重启恢复和明确终止 | P0 | NOT_STARTED | M3-001..003 | 无 UNKNOWN 任务；在途步骤按规则恢复/失败 | 待创建 | 定义 lease/heartbeat |
| M3-010 | M3-safe-ai | 完成20场景受约束 AI 回归 | P0 | NOT_STARTED | M1-010,M2-010,M3-001..009 | ≥70%无需改参数；AI写操作0 | 待创建 | 锁定测试模型 |
| M4-001 | M4-beta | 实现 Runbook/Version 数据模型 | P0 | NOT_STARTED | M3-005 | 成功 Plan 可保存、版本化、参数化 | 待创建 | 定义变量类型 |
| M4-002 | M4-beta | 实现 Runbook 执行仍经过 Operation Bus | P0 | NOT_STARTED | M4-001,M0-007 | 不可绕过 Policy/Audit/Evidence | 待创建 | 复用 Task Engine |
| M4-004 | M4-beta | 实现首次引导和示例诊断 | P0 | NOT_STARTED | M1/M2/M3 | 新用户10分钟内完成首个可信任务 | 待创建 | 招募5名测试用户 |
| M4-005 | M4-beta | 建立 Windows 核心 E2E | P0 | NOT_STARTED | M1-010,M2-010,M3-010 | 安装、Keychain、SSH、DB、任务恢复通过 | 待创建 | 准备签名/无签名路径 |
| M4-006 | M4-beta | 完成发布前威胁复核和秘密扫描 | P0 | NOT_STARTED | 全部P0 | 无明文秘密；极高风险均关闭或拒绝发布 | 待创建 | 安排外部复核可选 |
| M4-007 | M4-beta | 建立版本、安装包、校验和和发布流程 | P0 | NOT_STARTED | M0-003,M4-005 | 可重复构建 v0.1.0-beta，产物有 SHA256 | 待创建 | 定义版本注入 |
| M4-010 | M4-beta | 完成种子用户测试并发布 Beta | P0 | NOT_STARTED | M4-001..009 | 20安装/10连续使用两周/阻断问题关闭 | 待创建 | 建立反馈表 |

## 状态摘要

- `DONE`: 1
- `READY`: 1
- `IN_REVIEW`: 0
- `IN_PROGRESS`: 0
- `BLOCKED`: 0
- `NOT_STARTED`: 48