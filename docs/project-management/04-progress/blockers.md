> 状态：维护中  
> 负责人：vvitem  
> 最后更新：2026-07-11  
> 基线 Commit：`b590887fea1e2c43e7831a48932b9be5d44ccfa3`  
> 关联 Milestone：`M0-foundation`  
> 关联 Issue/PR：待创建

# 阻塞项

## 当前阻塞

当前无已确认的外部阻塞。

## 待决策但未阻塞

| ID | 事项 | 解除条件 | 下一步 | 最晚时间 |
|---|---|---|---|---|
| BLK-CAND-001 | 正式名称/Go Module 未确认 | 选择唯一名称和 module path | 完成名称检索并记录 DEC | 第1周 |
| BLK-CAND-002 | Windows Keychain 库未选择 | PoC 验证存取/删除/不可用错误 | 比较原生封装与成熟库 | 第2周 |
| BLK-CAND-003 | SQL Parser 未选择 | MySQL/PG 语法 corpus 通过 | M2 前做技术 Spike | 第10周 |

候选项只有在阻止当前任务验收时才转为正式 `BLOCKED`，并同步 Backlog 状态。
