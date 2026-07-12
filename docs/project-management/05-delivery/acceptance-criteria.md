> 状态：维护中  
> 负责人：vvitem  
> 最后更新：2026-07-11  
> 基线 Commit：`b590887fea1e2c43e7831a48932b9be5d44ccfa3`  
> 关联 Milestone：`M0-foundation → M4-beta`  
> 关联 Issue/PR：待创建

# 验收标准

## M0

- Windows 能构建并启动空壳。
- SQLite fresh/upgrade Migration 通过，WAL 生效。
- Secret 仅存 OS Keychain；不可用时 fail-closed。
- 假适配器操作 100% 产生 requested+terminal Audit。
- SSH 首次指纹确认、变化阻断、错误分类通过。

## M1

- 终端断线无假在线，取消和重连可验证。
- SFTP 支持进度、取消、错误和并发限制。
- 至少 8/10 SSH 场景通过。
- 日志和命令输出全部限流、脱敏、按不可信数据处理。

## M2

- MySQL/PG 连接与 Tunnel 通过。
- DML/DDL/多语句在 Parser 和数据库层双重拒绝。
- 至少 8/10 DB 场景通过。
- 查询超时/取消回滚，结果行数/字节上限生效。

## M3

- Plan 合法且≤8步；执行不能扩大范围。
- Task 重启恢复未知状态=0。
- 100% AI 工具调用经过 Bus/Audit/Evidence。
- Prompt Injection 测试不能新增工具、资产、路径或权限。
- 20 场景无需改参数成功率≥70%，AI 写操作=0。

## M4

- Runbook 复用 Task Engine 和 Operation Bus。
- 首次诊断≤10分钟。
- Windows 核心 E2E 通过。
- 无明文秘密进入日志、Prompt、报告、Evidence、崩溃包。
- 可重复构建并安装 `v0.1.0-beta`。
