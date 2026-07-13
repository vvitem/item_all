> 状态：维护中  
> 负责人：vvitem  
> 最后更新：2026-07-13  
> 基线 Commit：`197b63089c7b0c53b5a4c8cc6257ddbf68bb0e0e`  
> 关联 Milestone：`M0-foundation`  
> 关联 Issue/PR：[Issue #2](https://github.com/vvitem/item_all/issues/2)

# 决策日志

## DEC-001：采用 Personal-first Local SafeOps 路线

- **日期**：2026-07-11
- **状态**：已确认
- **背景**：单人、六个月、开源增长，不能复制大而全运维客户端。
- **问题**：产品应以资产数量还是诊断闭环为中心。
- **可选方案**：一站式客户端；Personal-first SafeOps；企业访问平台。
- **最终选择**：Personal-first SafeOps。
- **理由**：能够以 SSH/DB 场景形成可验证差异化，并兼容未来团队化。
- **影响**：删除 RDP/K8s/Redis/Kafka 等首版范围。
- **文档**：product-positioning、mvp-scope、roadmap。
- **新增任务**：M0–M4 全部 Backlog。
- **验证**：20 个场景成功率和首次诊断时间。

## DEC-002：仓库按绿地项目处理

- **日期**：2026-07-11
- **状态**：已确认
- **背景**：基线只存在 README。
- **问题**：是否可沿用规划文档中的“现有实现”判断。
- **可选方案**：假定已有 OpsKat 能力；以 item_all 代码事实为准。
- **最终选择**：全部业务能力标记未开始。
- **理由**：规划文档是目标输入，不是本仓库实现证据。
- **影响**：当前 DONE=0，M0 从工程初始化开始。
- **文档**：project-baseline、current-state、gap、progress。
- **新增任务**：M0-002 至 M0-010。
- **验证**：代码/测试落地后逐项更新。

## DEC-003：本次不修改业务代码

- **日期**：2026-07-11
- **状态**：已确认
- **背景**：本任务边界是文档、设计和追踪体系。
- **最终选择**：仅新增/更新 Markdown 和 GitHub 模板。
- **影响**：所有代码路径均标记“建议路径”。
- **验证**：PR changed files 不包含 Go/TS/配置业务实现。

## DEC-004：采用 ItemAll 工程身份和 Wails 标准根目录结构

- **日期**：2026-07-13
- **状态**：已确认
- **背景**：`M0-002` 需要在创建 `go.mod` 和 Wails 配置前冻结产品名称、Module Path、规范化应用标识与工程结构。
- **问题**：选择哪套产品身份和桌面工程布局，才能兼顾单人开发效率与后续模块扩展。
- **可选方案**：
  1. `ItemAll` + Wails 标准根目录结构。
  2. `ItemAll` + `cmd/desktop`。
  3. `apps/desktop` Monorepo。
- **最终选择**：
  - 产品名称：`ItemAll`。
  - Go Module：`github.com/vvitem/item_all`。
  - 规范化应用标识：`com.vvitem.itemall`。
  - 工程结构：根目录保留 `main.go`、`app.go`、`wails.json`，业务模块后续进入 `internal/*`。
- **选择理由**：与 Wails CLI 默认工作方式一致，降低单人和 Codex 辅助开发的脚手架复杂度；当前没有多应用 Monorepo 的实际需求。
- **影响范围**：`go.mod`、构建产物名称、前端展示名称、后续安装器和平台清单标识。
- **配置说明**：Wails v2 的 `wails.json` 没有通用 `applicationID` 字段，不得写入虚构配置键；后续打包任务将规范化标识映射到各平台 Manifest 或 Bundle 配置。
- **需要修改的文档**：`open-questions.md`、M0-002 Design Spec、后续实现追踪与发布配置。
- **需要新增的任务**：无需新增 Backlog；由 `M0-002` 落地。
- **后续验证方式**：Windows 本地执行 `wails dev` 与 `wails build`，确认应用名称、构建信息和平台配置符合本决策。
