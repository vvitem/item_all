> 状态：维护中  
> 负责人：vvitem  
> 最后更新：2026-07-11  
> 基线 Commit：`b590887fea1e2c43e7831a48932b9be5d44ccfa3`  
> 关联 Milestone：`M4-beta`  
> 关联 Issue/PR：待创建

# 安全发布检查表

- [ ] AI 工具目录不包含任意 Shell/PowerShell/DML/DDL。
- [ ] 所有工具有 Schema、超时、输出、路径/行数限制和错误码。
- [ ] UI、AI、Runbook 无绕过 Operation Bus 的调用。
- [ ] CredentialRef/Keychain 失败路径 fail-closed。
- [ ] 日志、Prompt、Evidence、报告、诊断包 Secret Scan 为0。
- [ ] Redaction 失败阻止持久化和模型回传。
- [ ] 远端日志/文件/SQL 结果按不可信输入处理。
- [ ] SQL Parser 逃逸 corpus 和真实只读权限测试通过。
- [ ] Audit requested/decision/terminal 事件完整。
- [ ] Task 崩溃恢复未知状态为0。
- [ ] 依赖漏洞、许可证和 SBOM 检查完成。
- [ ] Release 产物与校验和来自受控流程。
- [ ] 未关闭极高风险时发布被阻断。
