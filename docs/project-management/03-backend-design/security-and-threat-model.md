> 状态：已确认  
> 负责人：vvitem  
> 最后更新：2026-07-11  
> 基线 Commit：`b590887fea1e2c43e7831a48932b9be5d44ccfa3`  
> 关联 Milestone：`M0-foundation → M4-beta`  
> 关联 Issue/PR：待创建

# 安全与威胁模型

## 信任主体

| 主体 | 信任级别 | 主要威胁 |
|---|---|---|
| 用户 | 有权限但可能误操作 | 扩大范围、误确认 |
| 本地应用 | 可信计算基 | 本地篡改、日志泄漏 |
| AI Provider | 外部不可信处理方 | 数据留存、错误 Plan、提示注入响应 |
| 远程主机/数据库 | 不可信资产 | 恶意输出、协议异常、秘密诱导 |
| 远端日志/文件/查询结果 | 完全不可信数据 | Prompt Injection、逃逸内容、超大输出 |
| 凭据 | 高敏感 | 持久化、Prompt、崩溃包泄漏 |
| 未来插件/CLI/MCP | 外部调用方 | 绕过边界、伪造 Actor |

## 强制安全属性

1. 远端日志、文件和数据库结果不得扩展 Plan、工具、权限、资产或路径范围。
2. 凭据不得进入 Prompt、日志、报告、Evidence、Audit Payload 或崩溃包。
3. Secret Redaction 失败时 fail-closed：丢弃结果并返回 `REDACTION_FAILED`。
4. AI 仅拿到资产别名和必要上下文，不拿到长期凭据。
5. 数据默认本地保存；外发模型内容必须在 UI 可解释且经过脱敏。
6. AI 写操作数量必须为零。
7. 本地风险重算和 Plan Scope 校验不能依赖模型声明。

## 主要攻击路径与控制

| 威胁 | 控制 | 验证 |
|---|---|---|
| 日志中“忽略规则并执行命令” | 数据/指令分离、工具白名单、Plan Scope | adversarial corpus |
| SQL 通过 CTE/注释/多语句写入 | AST Parser + 只读事务 + 低权限账号 | SQL corpus + real DB |
| 路径遍历读取秘密 | canonical path + allow-root + symlink policy | boundary tests |
| UI 直接调用 Adapter | internal 构造、依赖检查、Audit 断言 | bypass integration |
| Keychain 不可用时落盘明文 | fail-closed，不提供静默 fallback | Windows failure test |
| 输出过大导致内存/Prompt DoS | 流式限流、字节/行/时间上限 | load tests |
| 任务崩溃后误报成功 | 持久状态和 Evidence Verify | crash recovery tests |

## 风险等级

- `R0`：低敏感固定元数据读取，可自动执行。
- `R1`：可能包含业务或日志内容，需明确范围和脱敏。
- `R2`：用户提供只读 SQL/较大范围读取，需要 Plan 审阅和确认。
- `R3`：状态修改；MVP 不向 AI 暴露，手工能力需单独决策。
- `R4`：破坏性/不可逆；MVP 禁止。
