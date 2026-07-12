> 状态：维护中  
> 负责人：vvitem  
> 最后更新：2026-07-13  
> 基线 Commit：`be37cfbb2bdca5a05c3303e9ea208992a8bf1721`  
> 关联 Milestone：`M0-foundation`  
> 关联 Issue/PR：[PR #1](https://github.com/vvitem/item_all/pull/1) · [Issue #2](https://github.com/vvitem/item_all/issues/2)

# 当前项目状态

## 基线信息

- 当前分支：`main`
- 当前 Commit：`be37cfbb2bdca5a05c3303e9ea208992a8bf1721`
- 最后更新时间：2026-07-13
- 当前 Milestone：`M0-foundation`
- 当前周次：第 1 周完成，进入第 2 周设计准备

## 本阶段目标

建立可构建的 Go/Wails/React 工程骨架、SQLite WAL、OS Keychain、Operation Bus、append-only Audit，以及 SSH 资产连接测试垂直切片。

## 当前正在进行

- `M0-002`：初始化 Go/Wails/React 单仓，状态 `READY`。
- 已创建 [Issue #2](https://github.com/vvitem/item_all/issues/2)，实现前需确认正式产品名称、Go Module Path 和 Wails Application ID。
- 下一执行动作是完成 `M0-002` 设计 Spec；在 Spec 获得确认前不写工程代码。

## 最近完成

- `M0-001`：项目分析、Roadmap、后端设计、进度追踪和交付门禁文档已通过 [PR #1](https://github.com/vvitem/item_all/pull/1) 合并。
- 合并 Commit：`be37cfbb2bdca5a05c3303e9ea208992a8bf1721`。
- 仓库已具备 20 个稳定场景、50 项 Backlog、24 周计划、6 个 ADR 和统一交付门禁。

## 下一步任务

1. 确认正式产品名称、Go Module Path 和 Wails Application ID。
2. 为 `M0-002` 编写并评审最小工程骨架设计 Spec。
3. 根据已确认 Spec 编写实现计划。
4. 实施 `M0-002`：Windows 可启动的最小 Go/Wails/React 空壳。
5. 完成后再启动 `M0-003` CI 门禁，不并行开发 SSH、AI 或数据库。

## 当前阻塞

没有外部阻塞。`M0-002` 存在一个设计决策门：正式产品名称和 Go Module Path 尚未确认；在该决策完成前不得创建最终 `go.mod`。

## 高风险事项

- 先写协议/UI 再补 Operation Bus 会造成执行边界返工。
- 未先确定 Module Path 和 Application ID 会导致初始化后全仓重命名。
- Windows Keychain/PTTY/WebView 兼容风险尚未验证。
- Codex 生成代码可能在日志和 DTO 中扩散秘密。

## 待确认决策

- 正式项目名称、Go Module Path 和 Wails Application ID。
- Windows Credential Manager 封装方案。
- 手工 SQL 写操作是否进入首版。

## 范围变化

无。当前严格遵循 SSH/SFTP/MySQL/PostgreSQL、只读 AI、六个月不做控制平面的范围。

## 测试状态

- Build：不存在。
- Unit/Integration/E2E：0。
- CI：不存在。
- 文档基线：PR #1 已合并；M0-001 有合并 Commit 作为完成证据。

## 相关代码路径

当前无业务代码。计划路径见 [module-boundaries.md](../03-backend-design/module-boundaries.md)。

## 相关 Issue / PR / Commit

- Initial Commit：`b590887fea1e2c43e7831a48932b9be5d44ccfa3`。
- [PR #1](https://github.com/vvitem/item_all/pull/1)：项目管理与设计基线，已合并。
- Merge Commit：`be37cfbb2bdca5a05c3303e9ea208992a8bf1721`。
- [Issue #2](https://github.com/vvitem/item_all/issues/2)：`M0-002` 工程骨架，状态 Open/READY。

## 本周可演示结果

从 `current-status.md` 可定位唯一当前状态；从 Issue #2 可查看下一项 P0 的边界、验收标准、测试和明确不做事项。