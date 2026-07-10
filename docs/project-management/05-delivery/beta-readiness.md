> 状态：维护中  
> 负责人：vvitem  
> 最后更新：2026-07-11  
> 基线 Commit：`b590887fea1e2c43e7831a48932b9be5d44ccfa3`  
> 关联 Milestone：`M4-beta`  
> 关联 Issue/PR：待创建

# Beta 就绪检查

## 用户价值

- [ ] 20 位真实安装用户。
- [ ] 10 位连续使用两周。
- [ ] 首次资产到第一次诊断中位数≤10分钟。
- [ ] 20 场景成功率≥70%。

## 功能

- [ ] SSH/SFTP/MySQL/PostgreSQL 核心路径通过。
- [ ] Operation Bus、Task、Evidence、Audit、Runbook 可用。
- [ ] AI Plan≤8步，AI 写操作=0。

## 可靠性

- [ ] 断网、超时、取消、崩溃恢复无未知状态。
- [ ] 大日志/大结果受限且 UI 不冻结。
- [ ] Migration fresh/upgrade 通过。

## 安全

- [ ] Secret Scan 为0。
- [ ] Prompt Injection 和 SQL 逃逸测试通过。
- [ ] 极高风险均关闭或明确阻断发布。

## 交付

- [ ] Windows 安装包、SHA256、版本信息和 Release Notes 完整。
- [ ] 已知问题、隐私说明、反馈入口和回滚步骤公开。
