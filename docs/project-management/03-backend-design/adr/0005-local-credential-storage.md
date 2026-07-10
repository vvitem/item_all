> 状态：Accepted  
> 负责人：vvitem  
> 最后更新：2026-07-11  
> 基线 Commit：`b590887fea1e2c43e7831a48932b9be5d44ccfa3`  
> 关联 Milestone：`M0-foundation`  
> 关联 Issue/PR：待创建

# ADR 0005：凭据使用 OS Keychain 和引用模型

- **日期**：2026-07-11
- **状态**：Accepted

## 背景

SQLite 加密仍涉及主密钥生命周期，且秘密容易进入 DTO、日志和备份。

## 决策

数据库只保存 CredentialRef，秘密由 OS Keychain 管理并以短 Lease 解析。

## 备选方案

- 维持临时实现并后补边界。
- 使用另一套桌面/执行架构。
- 将能力交给云端服务。

以上方案因安全边界、单人维护或本地优先目标不符合而未选择。

## 影响

需处理各平台差异；Keychain 不可用时 fail-closed。

## 验证

通过对应 Backlog、架构测试和 Beta 门禁验证；若事实推翻决策，创建新 ADR supersede 本记录。
