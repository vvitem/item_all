> 状态：维护中  
> 负责人：vvitem  
> 最后更新：2026-07-13  
> 基线 Commit：`197b63089c7b0c53b5a4c8cc6257ddbf68bb0e0e`  
> 关联 Milestone：`M0-foundation`  
> 关联 Issue/PR：Issue #2

# 待确认问题

| ID | 问题 | 当前建议 | 影响 | 确认方式 | 截止点 |
|---|---|---|---|---|---|
| OQ-002 | Windows 凭据存储具体实现 | 优先 Windows Credential Manager，封装统一 Keychain 接口 | 安全、迁移 | PoC + 威胁模型复核 | M0-W2 |
| OQ-003 | 首发模型 Provider | 最多支持两类兼容 Provider | 工期、测试矩阵 | 使用两个真实端点测试 | M3-W1 |
| OQ-004 | 手工 SQL 写操作是否进入首版 | 手工编辑器独立确认；AI 永久只读 | 产品边界、审计 | 用户访谈 + ADR | M2-W2 |
| OQ-005 | Evidence 大对象目录与清理周期 | SQLite 存元数据，文件系统存对象，默认 30 天 | 磁盘、隐私 | 真实日志压测 | M3-W2 |
| OQ-006 | Windows 代码签名 | Beta 可暂缓但必须明确提示；稳定版前完成 | 安装信任 | 成本和证书周期评估 | M4-W1 |
| OQ-007 | 是否发布本地 MCP | 仅在 M4 核心门禁提前通过后考虑 | 范围、外部攻击面 | Stop Line 评审 | M4-W2 |

## 已关闭问题

| ID | 结论 | 关闭日期 | 关联决策 |
|---|---|---|---|
| OQ-001 | 产品名称 `ItemAll`；Go Module `github.com/vvitem/item_all`；Application ID `com.vvitem.itemall` | 2026-07-13 | DEC-004 |

未确认项不得被文档写成已实现或最终承诺。
