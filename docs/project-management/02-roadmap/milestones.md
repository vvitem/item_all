> 状态：已确认  
> 负责人：vvitem  
> 最后更新：2026-07-11  
> 基线 Commit：`b590887fea1e2c43e7831a48932b9be5d44ccfa3`  
> 关联 Milestone：`M0-foundation → M4-beta`  
> 关联 Issue/PR：待创建

# Milestone 定义

## M0-foundation（第1–4周）

交付：单仓骨架、CI、SQLite Migration、CredentialRef、Operation Bus、Audit Event、SSH 连接测试。验收重点是“不能绕过边界”，而不是 UI 完整。

## M1-ssh-diagnosis（第5–9周）

交付：终端、SFTP 和 10 个 SSH 结构化工具。必须覆盖网络中断、超时、取消、主机指纹变化和输出截断。

## M2-database（第10–15周）

交付：MySQL/PostgreSQL、Tunnel、Schema/SQL/Explain 和 10 个数据库场景。数据库连接角色、事务和 SQL Parser 共同保证只读。

## M3-safe-ai（第16–20周）

交付：Provider、Plan Schema、Task Engine、Evidence Contract、Report。模型只参与 Plan 和基于证据总结，不执行无限循环。

## M4-beta（第21–24周）

交付：Runbook、首次引导、Windows E2E、诊断包、安装包和 Beta。此阶段只做硬化，不新增协议。

## 完成判定

Milestone 只有在全部 P0 验收完成、无未接受的极高风险、关键 E2E 通过且文档同步后才能标记 `DONE`。
