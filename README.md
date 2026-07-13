# ItemAll

Local-first SafeOps 桌面工作台，面向个人开发者、独立运维和全栈工程师，聚焦 Linux 故障诊断与 MySQL/PostgreSQL 数据库运维。

> 当前状态：项目处于 `M0-foundation` 工程骨架阶段。实现状态以 [项目管理入口](docs/project-management/README.md) 为准。

## 项目原则

- 本地优先：连接、凭据、审计和证据默认保存在本机。
- 安全优先：AI 只调用结构化、参数化、只读工具，不获得任意 Shell 或数据库写权限。
- 单一执行边界：UI、AI、Runbook 及未来 CLI/MCP 都必须经过 Operation Bus。
- 可验证：任务遵循 `Goal → Plan → Approval → Execute → Verify → Report`，重要结论必须有 Evidence。
- 范围受控：六个月首版仅覆盖 SSH、SFTP、MySQL、PostgreSQL，Windows 正式支持。

## 当前工程骨架

- Go Module：`github.com/vvitem/item_all`
- 桌面框架：Wails v2
- 前端：React、TypeScript、Vite
- 当前 Binding：只读 `GetAppInfo()`
- 当前页面：Loading、Ready、Error、Retry 四种可验证行为
- 当前不包含 SSH、数据库、AI、SQLite 业务表或云服务

## 开发环境基线

- Go 1.26.5
- Wails CLI v2.13.0
- Node.js 24.18.0 LTS
- pnpm 11.12.0
- Windows 10/11 与 WebView2

## 安装开发工具

```powershell
corepack enable
corepack prepare pnpm@11.12.0 --activate
go install github.com/wailsapp/wails/v2/cmd/wails@v2.13.0
wails doctor
```

## 安装依赖

```powershell
go mod download
Set-Location frontend
pnpm install --frozen-lockfile
Set-Location ..
```

`frontend/pnpm-workspace.yaml` 只批准 `esbuild` 的安装构建脚本；其他依赖脚本仍保持拒绝状态。

## 执行验证

```powershell
go test ./...
Set-Location frontend
pnpm typecheck
pnpm lint
pnpm test:run
pnpm build
Set-Location ..
```

## 启动开发模式

```powershell
wails dev
```

## Windows 生产构建

```powershell
$commit = git rev-parse --short=12 HEAD
$buildTime = (Get-Date).ToUniversalTime().ToString("yyyy-MM-ddTHH:mm:ssZ")
wails build -clean -trimpath -ldflags "-X github.com/vvitem/item_all/internal/buildinfo.Version=0.0.0-dev -X github.com/vvitem/item_all/internal/buildinfo.Commit=$commit -X github.com/vvitem/item_all/internal/buildinfo.BuildTime=$buildTime"
```

生成文件位于 `build/bin/ItemAll.exe`。

## 文档入口

- [项目管理总览](docs/project-management/README.md)
- [当前状态](docs/project-management/04-progress/current-status.md)
- [六个月路线图](docs/project-management/02-roadmap/six-month-roadmap.md)
- [后端架构](docs/project-management/03-backend-design/architecture-overview.md)
- [MVP 范围](docs/project-management/01-analysis/mvp-scope.md)
- [M0-002 设计](docs/superpowers/specs/2026-07-13-m0-002-project-scaffold-design.md)
- [M0-002 实现计划](docs/superpowers/plans/2026-07-13-m0-002-project-scaffold.md)
- [原始增强版规划](docs/project-management/00-reference/OpsKat_项目深度分析与差异化产品规划_增强版.md)
