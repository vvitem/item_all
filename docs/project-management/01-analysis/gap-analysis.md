> 状态：维护中  
> 负责人：vvitem  
> 最后更新：2026-07-11  
> 基线 Commit：`b590887fea1e2c43e7831a48932b9be5d44ccfa3`  
> 关联 Milestone：`M0-foundation`  
> 关联 Issue/PR：待创建

# 差距分析

| ID | 领域 | 当前状态 | 目标状态 | 差距 | 优先级 | 关联里程碑 | 建议行动 |
|---|---|---|---|---|---|---|---|
| GAP-001 | 工程骨架 | 只有 README | Go/Wails/React 可构建 | 缺全部工程文件 | P0 | M0 | 初始化单仓并建立 CI |
| GAP-002 | 执行边界 | 无执行入口 | 所有操作经过 Operation Bus | 缺接口、流水线和防绕过测试 | P0 | M0 | 先实现接口和假适配器 |
| GAP-003 | 本地存储 | 无数据库 | SQLite WAL + Migration | 缺 Schema、迁移和备份策略 | P0 | M0 | 建立首版表和 append-only Migration |
| GAP-004 | 凭据 | 无实现 | OS Keychain + 引用模型 | 缺 CredentialRef 和泄漏测试 | P0 | M0 | 只存引用，禁止秘密进 SQLite |
| GAP-005 | SSH | 无实现 | 稳定 SSH/终端/SFTP | 缺连接、指纹、重连和 UI | P0 | M1 | 从连接测试垂直切片开始 |
| GAP-006 | 数据库 | 无实现 | MySQL/PostgreSQL 只读运维 | 缺连接器、Tunnel、SQL UI | P0 | M2 | 先实现连接和元数据，再 Explain |
| GAP-007 | Task Engine | 无实现 | 持久状态机和恢复 | 缺 Task/Step 状态与事务 | P0 | M3 | 先用固定 Plan 验证状态机 |
| GAP-008 | AI | 无实现 | 最多 8 步结构化 Plan | 缺 Provider、Schema、预算和注入防护 | P0 | M3 | 两次模型调用，拒绝自主循环 |
| GAP-009 | Audit/Evidence | 无实现 | 全入口审计与证据 | 缺事件表、哈希和脱敏 | P0 | M0/M3 | 先于真实适配器实现 |
| GAP-010 | Runbook | 无实现 | 成功任务可参数化复用 | 缺版本、变量和执行入口 | P1 | M4 | 基于稳定 Plan Schema 建模 |
| GAP-011 | Windows | 无 CI/E2E | Windows Beta | 缺构建、PTY/WebView/Keychain 验证 | P0 | M4 | 从 M0 开始持续构建，M4 真机硬化 |
| GAP-012 | 开源工程 | 无模板/规范 | 可贡献项目 | 缺 Issue/PR 模板、贡献指南 | P1 | M0/M4 | 本次建立模板，后续补 CONTRIBUTING |

## 处理规则

- 上述差距是当前状态与目标状态的差异，不自动等同“缺陷”。
- 每项差距必须映射到 Backlog，并在代码落地后通过测试和 Commit 关闭。
- 若实际实现需要偏离目标设计，先写 Decision Log；重大变化创建 ADR。
