> 状态：维护中  
> 负责人：vvitem  
> 最后更新：2026-07-11  
> 基线 Commit：`b590887fea1e2c43e7831a48932b9be5d44ccfa3`  
> 关联 Milestone：`M0-foundation`  
> 关联 Issue/PR：本次文档 PR（创建后补充编号）

# 当前项目状态

## 基线信息

- 当前分支：`main`（本次变更将在独立文档分支提交）
- 当前 Commit：`b590887fea1e2c43e7831a48932b9be5d44ccfa3`
- 最后更新时间：2026-07-11
- 当前 Milestone：`M0-foundation`
- 当前周次：第 1 周

## 本阶段目标

建立可构建的 Go/Wails/React 工程骨架、SQLite WAL、OS Keychain、Operation Bus、append-only Audit，以及 SSH 资产连接测试垂直切片。

## 当前正在进行

- `M0-001`：建立项目分析、Roadmap、后端设计、进度追踪和交付门禁文档，状态 `IN_REVIEW`。

## 最近完成

- 仓库初始化：根 README 和项目一句话定位（Commit `b590887fea1e2c43e7831a48932b9be5d44ccfa3`）。
- 完成仓库真实状态扫描：确认尚无业务代码、构建、测试或 CI。

## 下一步任务

1. `M0-002`：初始化 Go/Wails/React 单仓。
2. `M0-003`：建立 CI 门禁。
3. `M0-004`：建立 SQLite WAL 和 Migration。
4. `M0-006`：Windows Keychain PoC。
5. `M0-007`：Operation Bus 接口和假适配器。

## 当前阻塞

无外部阻塞。需要在第 1 周确认正式产品名称和 Go module/application ID。

## 高风险事项

- 先写协议/UI 再补 Operation Bus 会造成执行边界返工。
- Windows Keychain/PTTY/WebView 兼容风险尚未验证。
- Codex 生成代码可能在日志和 DTO 中扩散秘密。

## 待确认决策

- 正式项目名称及 Go module path。
- Windows Credential Manager 封装方案。
- 手工 SQL 写操作是否进入首版。

## 范围变化

无。当前严格遵循 SSH/SFTP/MySQL/PostgreSQL、只读 AI、六个月不做控制平面的范围。

## 测试状态

- Build：不存在。
- Unit/Integration/E2E：0。
- CI：不存在。
- 文档完整性：本次分支生成后执行本地非空和链接检查。

## 相关代码路径

当前无业务代码。计划路径见 [module-boundaries.md](../03-backend-design/module-boundaries.md)。

## 相关 Issue / PR / Commit

- Commit：`b590887fea1e2c43e7831a48932b9be5d44ccfa3` — Initial commit。
- 文档 PR：待创建。

## 本周可演示结果

仓库具备一套基于真实基线的工程文档和任务体系，开发者能从 `current-status.md` 找到唯一当前状态、下一项 P0、对应设计和验收标准。
