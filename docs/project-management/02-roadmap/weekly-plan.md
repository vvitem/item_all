> 状态：已确认  
> 负责人：vvitem  
> 最后更新：2026-07-11  
> 基线 Commit：`b590887fea1e2c43e7831a48932b9be5d44ccfa3`  
> 关联 Milestone：`M0-foundation → M4-beta`  
> 关联 Issue/PR：待创建

# 24 周计划

每周投入比例：50% 当前垂直功能、20% 测试和真机、15% 安全审查、10% 文档/社区、5% 技术债。

## 第 1 周

- **本周主要可演示结果**：仓库能以明确品牌和工程决策开始开发
- **Milestone**：`M0-foundation`
- **对应场景**：按任务关联 SSH/DB 场景；基础周为 N/A
- **计划任务**：`M0-001,M0-002`
- **验收标准**：空壳应用设计和初始化方案通过评审
- **测试**：单元测试 + 对应集成/E2E；失败路径必须覆盖。
- **安全检查**：凭据、日志、Operation Bus 和取消传播复核。
- **文档更新**：current-status、task-tracking、implementation-trace、weekly-log、changelog。
- **预计风险**：单人估时偏差、Windows 环境差异、真实资产不稳定。
- **明确不做**：不引入业务功能

## 第 2 周

- **本周主要可演示结果**：Windows 启动桌面空壳并安全保存一个测试秘密
- **Milestone**：`M0-foundation`
- **对应场景**：按任务关联 SSH/DB 场景；基础周为 N/A
- **计划任务**：`M0-002,M0-004,M0-006`
- **验收标准**：Keychain PoC、SQLite WAL、Migration 通过
- **测试**：单元测试 + 对应集成/E2E；失败路径必须覆盖。
- **安全检查**：凭据、日志、Operation Bus 和取消传播复核。
- **文档更新**：current-status、task-tracking、implementation-trace、weekly-log、changelog。
- **预计风险**：单人估时偏差、Windows 环境差异、真实资产不稳定。
- **明确不做**：不做 SSH UI

## 第 3 周

- **本周主要可演示结果**：所有假操作经过 Operation Bus 并写审计
- **Milestone**：`M0-foundation`
- **对应场景**：按任务关联 SSH/DB 场景；基础周为 N/A
- **计划任务**：`M0-005,M0-007,M0-008`
- **验收标准**：绕过测试失败、TraceID 完整
- **测试**：单元测试 + 对应集成/E2E；失败路径必须覆盖。
- **安全检查**：凭据、日志、Operation Bus 和取消传播复核。
- **文档更新**：current-status、task-tracking、implementation-trace、weekly-log、changelog。
- **预计风险**：单人估时偏差、Windows 环境差异、真实资产不稳定。
- **明确不做**：不连接真实资产

## 第 4 周

- **本周主要可演示结果**：添加 SSH 资产并完成指纹确认的连接测试
- **Milestone**：`M0-foundation`
- **对应场景**：按任务关联 SSH/DB 场景；基础周为 N/A
- **计划任务**：`M0-009,M0-010,M0-003`
- **验收标准**：真实主机连接、错误分类、CI 通过
- **测试**：单元测试 + 对应集成/E2E；失败路径必须覆盖。
- **安全检查**：凭据、日志、Operation Bus 和取消传播复核。
- **文档更新**：current-status、task-tracking、implementation-trace、weekly-log、changelog。
- **预计风险**：单人估时偏差、Windows 环境差异、真实资产不稳定。
- **明确不做**：不做终端美化

## 第 5 周

- **本周主要可演示结果**：用户可打开稳定终端会话
- **Milestone**：`M1-ssh-diagnosis`
- **对应场景**：按任务关联 SSH/DB 场景；基础周为 N/A
- **计划任务**：`M1-001`
- **验收标准**：Windows 输入输出/resize/关闭通过
- **测试**：单元测试 + 对应集成/E2E；失败路径必须覆盖。
- **安全检查**：凭据、日志、Operation Bus 和取消传播复核。
- **文档更新**：current-status、task-tracking、implementation-trace、weekly-log、changelog。
- **预计风险**：单人估时偏差、Windows 环境差异、真实资产不稳定。
- **明确不做**：不做分屏和主题

## 第 6 周

- **本周主要可演示结果**：断网后状态明确且可重连
- **Milestone**：`M1-ssh-diagnosis`
- **对应场景**：按任务关联 SSH/DB 场景；基础周为 N/A
- **计划任务**：`M1-002`
- **验收标准**：故障注入无假在线
- **测试**：单元测试 + 对应集成/E2E；失败路径必须覆盖。
- **安全检查**：凭据、日志、Operation Bus 和取消传播复核。
- **文档更新**：current-status、task-tracking、implementation-trace、weekly-log、changelog。
- **预计风险**：单人估时偏差、Windows 环境差异、真实资产不稳定。
- **明确不做**：不做自动无限重连

## 第 7 周

- **本周主要可演示结果**：用户可浏览和传输文件
- **Milestone**：`M1-ssh-diagnosis`
- **对应场景**：按任务关联 SSH/DB 场景；基础周为 N/A
- **计划任务**：`M1-003,M1-004`
- **验收标准**：SFTP 进度/取消和 SSH Config 导入通过
- **测试**：单元测试 + 对应集成/E2E；失败路径必须覆盖。
- **安全检查**：凭据、日志、Operation Bus 和取消传播复核。
- **文档更新**：current-status、task-tracking、implementation-trace、weekly-log、changelog。
- **预计风险**：单人估时偏差、Windows 环境差异、真实资产不稳定。
- **明确不做**：不做在线编辑

## 第 8 周

- **本周主要可演示结果**：完成系统和磁盘诊断
- **Milestone**：`M1-ssh-diagnosis`
- **对应场景**：按任务关联 SSH/DB 场景；基础周为 N/A
- **计划任务**：`M1-005,M1-006`
- **验收标准**：SSH-01/02/03/04/10 通过
- **测试**：单元测试 + 对应集成/E2E；失败路径必须覆盖。
- **安全检查**：凭据、日志、Operation Bus 和取消传播复核。
- **文档更新**：current-status、task-tracking、implementation-trace、weekly-log、changelog。
- **预计风险**：单人估时偏差、Windows 环境差异、真实资产不稳定。
- **明确不做**：不新增协议

## 第 9 周

- **本周主要可演示结果**：完成服务、进程、端口和日志诊断
- **Milestone**：`M1-ssh-diagnosis`
- **对应场景**：按任务关联 SSH/DB 场景；基础周为 N/A
- **计划任务**：`M1-007..010`
- **验收标准**：至少8/10 SSH 场景通过
- **测试**：单元测试 + 对应集成/E2E；失败路径必须覆盖。
- **安全检查**：凭据、日志、Operation Bus 和取消传播复核。
- **文档更新**：current-status、task-tracking、implementation-trace、weekly-log、changelog。
- **预计风险**：单人估时偏差、Windows 环境差异、真实资产不稳定。
- **明确不做**：不做 AI 自由 Plan

## 第 10 周

- **本周主要可演示结果**：MySQL/PG 可连接并支持 Tunnel
- **Milestone**：`M2-database`
- **对应场景**：按任务关联 SSH/DB 场景；基础周为 N/A
- **计划任务**：`M2-001,M2-002`
- **验收标准**：容器+真实 Tunnel 通过
- **测试**：单元测试 + 对应集成/E2E；失败路径必须覆盖。
- **安全检查**：凭据、日志、Operation Bus 和取消传播复核。
- **文档更新**：current-status、task-tracking、implementation-trace、weekly-log、changelog。
- **预计风险**：单人估时偏差、Windows 环境差异、真实资产不稳定。
- **明确不做**：不做高级编辑

## 第 11 周

- **本周主要可演示结果**：用户可浏览数据库对象
- **Milestone**：`M2-database`
- **对应场景**：按任务关联 SSH/DB 场景；基础周为 N/A
- **计划任务**：`M2-003`
- **验收标准**：DB-02/03/04 元数据通过
- **测试**：单元测试 + 对应集成/E2E；失败路径必须覆盖。
- **安全检查**：凭据、日志、Operation Bus 和取消传播复核。
- **文档更新**：current-status、task-tracking、implementation-trace、weekly-log、changelog。
- **预计风险**：单人估时偏差、Windows 环境差异、真实资产不稳定。
- **明确不做**：不做 ER 图

## 第 12 周

- **本周主要可演示结果**：用户可安全执行只读 SQL
- **Milestone**：`M2-database`
- **对应场景**：按任务关联 SSH/DB 场景；基础周为 N/A
- **计划任务**：`M2-004,M2-005`
- **验收标准**：所有写语句被双层拒绝
- **测试**：单元测试 + 对应集成/E2E；失败路径必须覆盖。
- **安全检查**：凭据、日志、Operation Bus 和取消传播复核。
- **文档更新**：current-status、task-tracking、implementation-trace、weekly-log、changelog。
- **预计风险**：单人估时偏差、Windows 环境差异、真实资产不稳定。
- **明确不做**：不做数据编辑

## 第 13 周

- **本周主要可演示结果**：用户可查看会话、锁和连接压力
- **Milestone**：`M2-database`
- **对应场景**：按任务关联 SSH/DB 场景；基础周为 N/A
- **计划任务**：`M2-006`
- **验收标准**：DB-01/06/07/08 通过
- **测试**：单元测试 + 对应集成/E2E；失败路径必须覆盖。
- **安全检查**：凭据、日志、Operation Bus 和取消传播复核。
- **文档更新**：current-status、task-tracking、implementation-trace、weekly-log、changelog。
- **预计风险**：单人估时偏差、Windows 环境差异、真实资产不稳定。
- **明确不做**：不提供 kill

## 第 14 周

- **本周主要可演示结果**：用户可分析 Explain 和表统计
- **Milestone**：`M2-database`
- **对应场景**：按任务关联 SSH/DB 场景；基础周为 N/A
- **计划任务**：`M2-007,M2-008`
- **验收标准**：DB-05/09/10 通过
- **测试**：单元测试 + 对应集成/E2E；失败路径必须覆盖。
- **安全检查**：凭据、日志、Operation Bus 和取消传播复核。
- **文档更新**：current-status、task-tracking、implementation-trace、weekly-log、changelog。
- **预计风险**：单人估时偏差、Windows 环境差异、真实资产不稳定。
- **明确不做**：不自动建索引

## 第 15 周

- **本周主要可演示结果**：数据库闭环可回归
- **Milestone**：`M2-database`
- **对应场景**：按任务关联 SSH/DB 场景；基础周为 N/A
- **计划任务**：`M2-009,M2-010`
- **验收标准**：至少8/10 DB 场景通过
- **测试**：单元测试 + 对应集成/E2E；失败路径必须覆盖。
- **安全检查**：凭据、日志、Operation Bus 和取消传播复核。
- **文档更新**：current-status、task-tracking、implementation-trace、weekly-log、changelog。
- **预计风险**：单人估时偏差、Windows 环境差异、真实资产不稳定。
- **明确不做**：不做第三种数据库

## 第 16 周

- **本周主要可演示结果**：固定 Plan 可由 Task Engine 持久执行
- **Milestone**：`M3-safe-ai`
- **对应场景**：按任务关联 SSH/DB 场景；基础周为 N/A
- **计划任务**：`M3-001,M3-002`
- **验收标准**：状态机/取消/重试通过
- **测试**：单元测试 + 对应集成/E2E；失败路径必须覆盖。
- **安全检查**：凭据、日志、Operation Bus 和取消传播复核。
- **文档更新**：current-status、task-tracking、implementation-trace、weekly-log、changelog。
- **预计风险**：单人估时偏差、Windows 环境差异、真实资产不稳定。
- **明确不做**：自由 Plan 不稳定则不启用

## 第 17 周

- **本周主要可演示结果**：每一步产生可验证 Evidence
- **Milestone**：`M3-safe-ai`
- **对应场景**：按任务关联 SSH/DB 场景；基础周为 N/A
- **计划任务**：`M3-003`
- **验收标准**：Evidence Contract 和清理通过
- **测试**：单元测试 + 对应集成/E2E；失败路径必须覆盖。
- **安全检查**：凭据、日志、Operation Bus 和取消传播复核。
- **文档更新**：current-status、task-tracking、implementation-trace、weekly-log、changelog。
- **预计风险**：单人估时偏差、Windows 环境差异、真实资产不稳定。
- **明确不做**：不上传云端

## 第 18 周

- **本周主要可演示结果**：模型只能输出合法8步 Plan
- **Milestone**：`M3-safe-ai`
- **对应场景**：按任务关联 SSH/DB 场景；基础周为 N/A
- **计划任务**：`M3-004,M3-005`
- **验收标准**：Schema/预算/错误处理通过
- **测试**：单元测试 + 对应集成/E2E；失败路径必须覆盖。
- **安全检查**：凭据、日志、Operation Bus 和取消传播复核。
- **文档更新**：current-status、task-tracking、implementation-trace、weekly-log、changelog。
- **预计风险**：单人估时偏差、Windows 环境差异、真实资产不稳定。
- **明确不做**：不做多 Agent

## 第 19 周

- **本周主要可演示结果**：恶意日志不能扩大权限或范围
- **Milestone**：`M3-safe-ai`
- **对应场景**：按任务关联 SSH/DB 场景；基础周为 N/A
- **计划任务**：`M3-006,M3-007`
- **验收标准**：攻击语料全拒绝
- **测试**：单元测试 + 对应集成/E2E；失败路径必须覆盖。
- **安全检查**：凭据、日志、Operation Bus 和取消传播复核。
- **文档更新**：current-status、task-tracking、implementation-trace、weekly-log、changelog。
- **预计风险**：单人估时偏差、Windows 环境差异、真实资产不稳定。
- **明确不做**：不允许任意工具

## 第 20 周

- **本周主要可演示结果**：任务可恢复并生成有证据的报告
- **Milestone**：`M3-safe-ai`
- **对应场景**：按任务关联 SSH/DB 场景；基础周为 N/A
- **计划任务**：`M3-008..010`
- **验收标准**：20场景≥70%，未知状态0
- **测试**：单元测试 + 对应集成/E2E；失败路径必须覆盖。
- **安全检查**：凭据、日志、Operation Bus 和取消传播复核。
- **文档更新**：current-status、task-tracking、implementation-trace、weekly-log、changelog。
- **预计风险**：单人估时偏差、Windows 环境差异、真实资产不稳定。
- **明确不做**：不做 CLI/MCP

## 第 21 周

- **本周主要可演示结果**：成功任务可保存和再次运行
- **Milestone**：`M4-beta`
- **对应场景**：按任务关联 SSH/DB 场景；基础周为 N/A
- **计划任务**：`M4-001,M4-002`
- **验收标准**：Runbook 仍经 Bus/Audit
- **测试**：单元测试 + 对应集成/E2E；失败路径必须覆盖。
- **安全检查**：凭据、日志、Operation Bus 和取消传播复核。
- **文档更新**：current-status、task-tracking、implementation-trace、weekly-log、changelog。
- **预计风险**：单人估时偏差、Windows 环境差异、真实资产不稳定。
- **明确不做**：不做共享市场

## 第 22 周

- **本周主要可演示结果**：新用户10分钟完成首个诊断
- **Milestone**：`M4-beta`
- **对应场景**：按任务关联 SSH/DB 场景；基础周为 N/A
- **计划任务**：`M4-003,M4-004`
- **验收标准**：5名用户中≥3名独立完成
- **测试**：单元测试 + 对应集成/E2E；失败路径必须覆盖。
- **安全检查**：凭据、日志、Operation Bus 和取消传播复核。
- **文档更新**：current-status、task-tracking、implementation-trace、weekly-log、changelog。
- **预计风险**：单人估时偏差、Windows 环境差异、真实资产不稳定。
- **明确不做**：不增加模板数量换质量

## 第 23 周

- **本周主要可演示结果**：Windows 候选版通过安全和E2E门禁
- **Milestone**：`M4-beta`
- **对应场景**：按任务关联 SSH/DB 场景；基础周为 N/A
- **计划任务**：`M4-005,M4-006,M4-008`
- **验收标准**：核心矩阵通过、秘密0
- **测试**：单元测试 + 对应集成/E2E；失败路径必须覆盖。
- **安全检查**：凭据、日志、Operation Bus 和取消传播复核。
- **文档更新**：current-status、task-tracking、implementation-trace、weekly-log、changelog。
- **预计风险**：单人估时偏差、Windows 环境差异、真实资产不稳定。
- **明确不做**：停止其他平台优化

## 第 24 周

- **本周主要可演示结果**：可重复发布 v0.1.0-beta
- **Milestone**：`M4-beta`
- **对应场景**：按任务关联 SSH/DB 场景；基础周为 N/A
- **计划任务**：`M4-007,M4-009,M4-010`
- **验收标准**：安装、校验和、发布说明和反馈入口完成
- **测试**：单元测试 + 对应集成/E2E；失败路径必须覆盖。
- **安全检查**：凭据、日志、Operation Bus 和取消传播复核。
- **文档更新**：current-status、task-tracking、implementation-trace、weekly-log、changelog。
- **预计风险**：单人估时偏差、Windows 环境差异、真实资产不稳定。
- **明确不做**：禁止临时加入协议
