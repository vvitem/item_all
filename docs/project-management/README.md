> 状态：维护中  
> 负责人：vvitem  
> 最后更新：2026-07-11  
> 基线 Commit：`b590887fea1e2c43e7831a48932b9be5d44ccfa3`  
> 关联 Milestone：`M0-foundation`  
> 关联 Issue/PR：待创建

# 项目管理与工程设计入口

本目录将产品规划拆解为可执行、可追踪、可验收的工程体系。事实来源优先级为：当前代码与测试 > 已合并 PR/Commit > ADR/Decision Log > 参考规划 > Backlog/Weekly Plan。

## 核心入口

| 目的 | 文档 |
|---|---|
| 查看当前真实状态 | [04-progress/current-status.md](04-progress/current-status.md) |
| 查看仓库基线 | [01-analysis/project-baseline.md](01-analysis/project-baseline.md) |
| 查看 As-Is / To-Be / Gap | [01-analysis/current-state-analysis.md](01-analysis/current-state-analysis.md)、[01-analysis/target-state.md](01-analysis/target-state.md)、[01-analysis/gap-analysis.md](01-analysis/gap-analysis.md) |
| 查看六个月范围与路线图 | [01-analysis/mvp-scope.md](01-analysis/mvp-scope.md)、[02-roadmap/six-month-roadmap.md](02-roadmap/six-month-roadmap.md) |
| 查看后端设计 | [03-backend-design/architecture-overview.md](03-backend-design/architecture-overview.md) |
| 查看任务 | [02-roadmap/backlog.md](02-roadmap/backlog.md)、[04-progress/task-tracking.md](04-progress/task-tracking.md) |
| 查看交付门禁 | [05-delivery/definition-of-done.md](05-delivery/definition-of-done.md)、[05-delivery/beta-readiness.md](05-delivery/beta-readiness.md) |

## 当前事实

仓库在基线 Commit 上只有根目录 `README.md`。不存在 `go.mod`、`wails.json`、`frontend/`、Migration、测试或 GitHub Actions。因此：

- 已实现产品能力：无。
- 已完成工程决策：本文档体系中的已接受 ADR。
- 当前 Milestone：`M0-foundation`。
- 当前开发状态：文档体系进入评审，工程初始化尚未开始。

## 维护规则

每次完成开发任务后，至少同步更新：

1. [current-status.md](04-progress/current-status.md)
2. [task-tracking.md](04-progress/task-tracking.md)
3. [milestone-status.md](04-progress/milestone-status.md)
4. [implementation-trace.md](04-progress/implementation-trace.md)
5. [weekly-log.md](04-progress/weekly-log.md)
6. [changelog.md](04-progress/changelog.md)

任何 `DONE` 必须有 Commit、PR、测试报告或可复现步骤作为证据。计划路径统一标记为“建议路径”，直到代码实际落地。
