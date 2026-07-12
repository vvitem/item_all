> 状态：已确认  
> 负责人：vvitem  
> 最后更新：2026-07-11  
> 基线 Commit：`b590887fea1e2c43e7831a48932b9be5d44ccfa3`  
> 关联 Milestone：`M0-foundation → M4-beta`  
> 关联 Issue/PR：待创建

# MVP 范围

## 六个月必须完成

### SSH 与主机

- SSH 资产、新增/编辑/删除和连接测试。
- SSH Config 导入。
- 主机指纹首次确认与变化阻断。
- 终端、多标签基础能力、断线状态和重连。
- 基础 SFTP 浏览、上传、下载。

### 数据库

- MySQL、PostgreSQL 与 SSH Tunnel。
- Schema/Table/Column/Index 浏览。
- SQL 编辑、结果表格、历史。
- `EXPLAIN` 与慢 SQL 辅助。
- 只读事务、超时、最大行数和脱敏。

### 安全执行

- Operation Bus 唯一执行入口。
- 持久 Task Engine。
- 结构化只读诊断工具。
- 最多 8 步 AI Plan。
- Evidence、Audit、稳定错误码、超时和取消。
- 成功任务保存为 Runbook。

### 交付

- Windows 核心 E2E。
- 首次引导和 10 分钟首次诊断目标。
- `v0.1.0-beta` 发布。

## 有余力再做

- 只读 `opsctl`。
- 本地 MCP。
- SQLite 浏览。
- macOS/Linux 安装体验优化。
- 超出内置 10–15 个模板的更多诊断模板。

## 六个月明确不做

- RDP/VNC、Kubernetes、Redis、Kafka、MongoDB、对象存储。
- 团队控制平面、账号、RBAC、审批后台、云同步。
- 插件市场、WASM 运行时、计费、移动端、多 Agent。
- 自动执行修复。
- AI 任意 Shell/PowerShell。
- AI 数据库 DML/DDL。

## 范围准入规则

新增功能必须同时满足：

1. 直接提高 SSH-01…SSH-10 或 DB-01…DB-10 的成功率。
2. 不引入新协议或云服务。
3. 有明确验收和测试。
4. 不削弱 Operation Bus、只读边界或脱敏。
5. 不触发 [scope-stop-lines.md](../02-roadmap/scope-stop-lines.md) 的停止条件。
