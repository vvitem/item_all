# ItemAll CI 开发与故障排查

ItemAll 使用 `.github/workflows/ci.yml` 的三个独立门禁：

- `quality`：Go、CI 检查器、Frontend 测试与构建。
- `generated-and-security`：Binding/锁文件漂移、依赖政策、仓库边界和完整历史 Secret Scan。
- `windows-build`：Windows Wails production build 与 `ItemAll.exe` 启动烟测。

三个 Job 在 Pull Request、`main` push 和手工触发时运行。Workflow 不使用路径过滤，Job 名称保持稳定，供后续 Branch Protection 设置 required checks。

## 固定工具链

- Go 1.26.5
- Wails v2.13.0
- Node.js 24.18.0
- pnpm 11.12.0
- Gitleaks v8.30.0
- Ubuntu 24.04
- Windows Server 2025

所有 GitHub Action 均固定到官方 Release/Tag 对应的完整 Commit SHA，并保留版本注释。Workflow 顶层权限只有 `contents: read`，不读取仓库或环境 Secret，不使用 `pull_request_target` 或自托管 Runner。

## 本地质量检查

Windows PowerShell：

```powershell
go test ./scripts/ci -v
go run ./scripts/ci dependencies
go run ./scripts/ci workflow
go run ./scripts/ci boundary

go test ./...
Set-Location frontend
pnpm install --frozen-lockfile
pnpm test:run
pnpm typecheck
pnpm lint
pnpm build
Set-Location ..
```

全量策略入口：

```powershell
go run ./scripts/ci all
```

Ubuntu 24.04 需要安装 Wails 编译依赖：

```bash
sudo apt-get update
sudo apt-get install --yes --no-install-recommends \
  build-essential \
  pkg-config \
  libgtk-3-dev \
  libwebkit2gtk-4.1-dev
```

Ubuntu 上对 Wails 根包执行 `go vet`、`go test` 或 Binding 生成时使用 `webkit2_41` Build Tag：

```bash
go vet -tags=webkit2_41 ./...
go test -tags=webkit2_41 ./...
GOFLAGS='-tags=webkit2_41' wails generate module
```

## 生成文件与锁文件

```powershell
go install github.com/wailsapp/wails/v2/cmd/wails@v2.13.0
wails generate module
go mod tidy
Set-Location frontend
pnpm install --frozen-lockfile
Set-Location ..
```

重新生成后必须检查：

- `frontend/wailsjs`
- `go.mod`
- `go.sum`
- `frontend/pnpm-lock.yaml`

CI 不自动提交生成结果。Wails v2.13.0 在 Linux 上会把生成的 `frontend/wailsjs` 文件模式改为可执行；CI 在生成后将普通文件统一恢复为 `0644`，再执行严格 `git diff`。这只消除跨平台文件模式噪声，不忽略内容漂移。

## pnpm allowBuilds

`frontend/pnpm-workspace.yaml` 当前只批准：

```yaml
allowBuilds:
  esbuild: true
```

新增安装脚本授权必须在独立依赖或功能设计中说明原因，并更新测试。禁止全局允许依赖脚本。

## Secret Scan

Gitleaks 仓库已迁移，但 v8.30.0 的 Go Module 仍声明为 `github.com/zricethezav/gitleaks/v8`，因此安装命令为：

```powershell
go install github.com/zricethezav/gitleaks/v8@v8.30.0
gitleaks git --redact --no-banner .
```

`generated-and-security` 使用 `fetch-depth: 0` 扫描完整可访问 Git 历史，输出保持脱敏。误报只能采用最小、可解释的局部处理；禁止宽泛排除整个目录或规则，也不得把真实 Token、私钥或高度仿真的秘密提交到测试分支。

## Windows 构建与烟测

```powershell
wails build -clean -trimpath
./scripts/ci/smoke-windows.ps1 -Executable build/bin/ItemAll.exe -ObservationSeconds 10
```

烟测验证：

- Windows 能加载 `ItemAll.exe`。
- Wails/WebView2 没有立即崩溃。
- 进程在观察窗口内保持运行。
- 测试结束后清理完整进程树。

烟测不验证视觉布局和交互，不替代人工 UI 检查或后续 Windows E2E。

## Action SHA 核验

当前固定引用：

- `actions/checkout@9c091bb21b7c1c1d1991bb908d89e4e9dddfe3e0` — v7.0.0
- `actions/setup-go@924ae3a1cded613372ab5595356fb5720e22ba16` — v6.5.0
- `actions/setup-node@820762786026740c76f36085b0efc47a31fe5020` — v7.0.0
- `actions/cache@55cc8345863c7cc4c66a329aec7e433d2d1c52a9` — v6.1.0

升级时从官方仓库 Release/Tag 解析完整 SHA，使用独立 PR 更新，不直接使用 `@main`、`@master` 或可移动大版本 Tag。

## required checks 自举

1. 在 CI 实现 PR 中让三个 Job 运行并变绿。
2. 完成一次受控 RED→GREEN 演练。
3. 合并实现 PR。
4. 在 `main` Branch Protection 中要求 `quality`、`generated-and-security`、`windows-build`。
5. 使用后续状态同步 PR 验证 required checks 确实阻止未通过的合并。

Job ID 和显示名称进入保护规则后不得随意改变。需要重命名时必须先同步 Branch Protection。

## 故障原则

失败必须修复根因。禁止：

- `continue-on-error: true`
- `|| true`
- 忽略退出码
- 删除失败检查
- 降低扫描范围或权限约束
- 把失败转换为 warning
- 改用 `ubuntu-latest` 或 `windows-latest` 规避环境问题

每次 Runner 失败都应记录失败 Job、具体步骤、根因和最小修复。

## 当前验证证据

- GREEN Run `29312670714`：`quality`、`generated-and-security`、`windows-build` 全部通过。
- 该 Run 验证了 Go/Frontend 全量检查、真实 Workflow 策略、三类生成漂移、完整历史 Gitleaks、Windows Wails build 和 EXE 启动烟测。
- 受控 RED→GREEN 证明将在实现 PR 合并前追加到本节。
