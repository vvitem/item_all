> 状态：已确认  
> 负责人：vvitem  
> 最后更新：2026-07-11  
> 基线 Commit：`b590887fea1e2c43e7831a48932b9be5d44ccfa3`  
> 关联 Milestone：`M0-foundation → M4-beta`  
> 关联 Issue/PR：待创建

# 测试策略

## 测试金字塔

1. 纯函数/Schema/Parser/状态机单测。
2. Repository、Operation Pipeline、Credential fake 集成测试。
3. Docker MySQL/PostgreSQL、测试 SSH 主机协议测试。
4. Wails/Windows 核心 E2E。
5. 20 个场景的验收回归和安全对抗测试。

## 必测矩阵

| 领域 | 正常 | 失败 | 取消/超时 | 安全 |
|---|---|---|---|---|
| Operation Bus | typed result | adapter error | context propagation | bypass、risk、redaction |
| Task Engine | dependency execution | step failure | pause/cancel/restart | scope immutability |
| SSH | real connect/tool parse | DNS/auth/fingerprint | disconnect | path/log injection |
| Database | MySQL+PG | TLS/auth/parser | query cancel | DML/DDL escape |
| Credential | store/resolve/delete | keychain unavailable | N/A | no plaintext/leak |
| Evidence | extract/store/export | schema/object error | large output | redaction fail-closed |

## Windows E2E 核心路径

安装/启动 → 创建资产 → Keychain → SSH 指纹/终端 → MySQL/PG 查询 → 固定诊断 Task → 取消 → 重启恢复 → 导出脱敏诊断包。

## 质量门槛

- P0 模块关键分支必须有测试，不用单一覆盖率数字替代风险验证。
- 每个工具至少有 input boundary、output truncation、cancel、redaction 测试。
- 每个修复必须先增加回归测试。
- Beta 前 20 场景成功率≥70%，AI 写操作=0，未知任务状态=0。
