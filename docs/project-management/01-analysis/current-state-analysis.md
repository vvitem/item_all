> 状态：已确认  
> 负责人：vvitem  
> 最后更新：2026-07-11  
> 基线 Commit：`b590887fea1e2c43e7831a48932b9be5d44ccfa3`  
> 关联 Milestone：`M0-foundation`  
> 关联 Issue/PR：待创建

# 当前状态分析（As-Is）

## 架构

当前没有应用架构。仓库仅包含 README，不存在桌面进程、IPC、后端模块或存储层。

## 功能

| 领域 | 当前状态 | 证据 |
|---|---|---|
| 资产管理 | 尚未开始 | 无 `internal/asset` 或前端目录 |
| SSH/SFTP | 尚未开始 | 无 Go Module 和协议依赖 |
| 数据库 | 尚未开始 | 无连接器、SQL 编辑器或 Migration |
| AI | 尚未开始 | 无 Provider、Tool Registry 或 Task Engine |
| Operation Bus | 尚未开始 | 无代码 |
| Audit/Evidence | 尚未开始 | 无表、事件或存储接口 |
| Runbook | 尚未开始 | 无模型和页面 |

## 数据模型

不存在数据库文件、Schema 或 Migration。`workspace_id`、`actor_id`、`device_id` 等字段仅是目标设计。

## 执行入口

不存在执行入口。根 README 不包含命令、构建或运行说明。

## 安全边界

目前没有凭据、网络连接或模型调用，因此没有运行时攻击面；同时也没有任何已经验证的安全控制。

## 测试覆盖

- 单元测试：0 个。
- 集成测试：0 个。
- E2E：0 个。
- 安全测试：0 个。

“0 个测试”不是低覆盖率，而是工程尚未初始化。

## 交付能力

不存在 GitHub Actions、构建产物、安装包、版本策略、签名或发布检查表。

## 结论

本项目是绿地项目。正确策略不是迁移或重构，而是先建立执行边界、数据约束和可测试骨架，再逐个交付垂直场景。
