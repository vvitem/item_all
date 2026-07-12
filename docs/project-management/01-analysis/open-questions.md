> 状态：维护中  
> 负责人：vvitem  
> 最后更新：2026-07-11  
> 基线 Commit：`b590887fea1e2c43e7831a48932b9be5d44ccfa3`  
> 关联 Milestone：`M0-foundation`  
> 关联 Issue/PR：待创建

# 待确认问题

| ID | 问题 | 当前建议 | 影响 | 确认方式 | 截止点 |
|---|---|---|---|---|---|
| OQ-001 | 正式项目名称是否使用 `item_all` | 在 M0 第 1 周确定品牌名 | 包名、应用 ID、签名 | 名称检索与仓库决策 | M0-W1 |
| OQ-002 | Windows 凭据存储具体实现 | 优先 Windows Credential Manager，封装统一 Keychain 接口 | 安全、迁移 | PoC + 威胁模型复核 | M0-W2 |
| OQ-003 | 首发模型 Provider | OpenAI-compatible + Anthropic-compatible，最多两类 | 工期、测试矩阵 | 选 2 个真实端点测试 | M3-W1 |
| OQ-004 | 手工 SQL 写操作是否进入首版 | 默认允许在独立手工编辑器中经高风险确认；AI 永久只读 | 产品边界、审计 | 用户访谈 + ADR | M2-W2 |
| OQ-005 | Evidence 大对象目录与清理周期 | SQLite 存元数据，文件系统存对象，默认 30 天 | 磁盘、隐私 | 真实日志压测 | M3-W2 |
| OQ-006 | Windows 代码签名 | Beta 可暂缓但必须明确提示；稳定版前完成 | 安装信任 | 成本和证书周期评估 | M4-W1 |
| OQ-007 | 是否发布本地 MCP | 仅在 M4 核心门禁提前通过后做 | 范围、外部攻击面 | Stop Line 评审 | M4-W2 |

未确认项不得被文档写成已实现或最终承诺。
