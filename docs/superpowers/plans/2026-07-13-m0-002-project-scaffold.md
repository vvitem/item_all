# M0-002 ItemAll 工程骨架 Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 建立一个可在 Windows 本地开发、测试和构建的 ItemAll Go + Wails v2 + React + TypeScript 最小桌面工程，并通过只读 `GetAppInfo()` 展示可信构建信息。

**Architecture:** 使用 Wails v2 标准根目录结构。根目录 `main.go` 只负责应用装配，`app.go` 只暴露只读 DTO Binding，构建元数据封装在 `internal/buildinfo`；前端是单页 React 应用，通过生成的 Wails Binding 获取 AppInfo，并只维护 `LOADING`、`READY`、`ERROR` 三种状态。

**Tech Stack:** Go 1.26.5、Wails v2.13.0、Node.js 24.18.0 LTS、pnpm 11.12.0、React 19.1.0、TypeScript 5.6.3、Vite 7.0.0、Vitest 3.2.4、React Testing Library 16.3.0。

## Global Constraints

- 产品名称固定为 `ItemAll`。
- Go Module 固定为 `github.com/vvitem/item_all`。
- 规范化应用标识固定为 `com.vvitem.itemall`；Wails v2 `wails.json` 不写入不存在的 `applicationID` 字段。
- 首发平台为 Windows；macOS/Linux 仅保持可编译结构，不做专项优化。
- 使用 Wails v2 标准根目录结构，不建立 `apps/desktop` 或提前建立多应用 Monorepo。
- `main.go` 只做 Composition Root；`app.go` 只做 DTO 转换和 Wails Binding。
- 本任务不实现 SSH、SFTP、数据库、SQLite 业务表、Operation Bus、Task Engine、Audit、Evidence、AI、Runbook、账号、云服务或遥测。
- 不读取或展示环境变量，不加入 `.env`，不写入 Token、API Key、密码或示例密钥。
- 不加载 CDN、远程字体或远程脚本。
- 所有 npm 依赖必须以精确版本写入 `package.json`，并提交 `pnpm-lock.yaml`。
- 前端不引入 Router、全局状态库、请求缓存库、UI 组件库或 Tailwind。
- 业务代码必须按 TDD 顺序实施：失败测试 → 最小实现 → 通过测试 → 提交。
- `M0-002` 在实现 PR 合并前只能标记为 `IN_REVIEW`；合并后才能标记为 `DONE`。

---

## File Map

### 根目录

- Create: `go.mod` — 固定 Go Module、Go toolchain 和 Wails v2.13.0。
- Create: `go.sum` — 由 `go mod tidy` 生成并提交。
- Create: `main.go` — Wails 应用装配、前端资源嵌入和窗口参数。
- Create: `app.go` — `App`、`AppInfo` 和 `GetAppInfo()` Binding。
- Create: `app_test.go` — Binding DTO 单元测试。
- Create: `wails.json` — Wails v2 项目配置和 pnpm 命令。
- Create: `.nvmrc` — 固定 Node.js 24.18.0。
- Create: `.npmrc` — 固定精确依赖策略。
- Create: `.gitignore` — 排除构建产物、依赖和编辑器文件。
- Create: `Makefile` — 可选命令别名；底层命令仍必须能在 Windows PowerShell 直接运行。
- Modify: `README.md` — 增加开发、测试和 Windows 构建步骤。

### Go 内部包

- Create: `internal/buildinfo/info.go` — 构建字段、默认值和运行时信息。
- Create: `internal/buildinfo/info_test.go` — 默认值、注入值和平台字段测试。

### Wails 资源

- Create: `build/appicon.png` — 从 Wails v2.13.0 官方模板生成。
- Create: `build/darwin/Info.dev.plist` — 从官方模板生成。
- Create: `build/darwin/Info.plist` — 从官方模板生成。
- Create: `build/windows/icon.ico` — 从官方模板生成。
- Create: `build/windows/info.json` — 从官方模板生成并将产品名改为 ItemAll。
- Create: `build/windows/wails.exe.manifest` — 从官方模板生成。

### 前端配置

- Create: `frontend/package.json` — 精确版本、脚本和 pnpm 版本。
- Create: `frontend/pnpm-lock.yaml` — 由 pnpm 11.12.0 生成。
- Create: `frontend/index.html` — Vite 入口。
- Create: `frontend/tsconfig.json` — 浏览器 TypeScript 配置。
- Create: `frontend/tsconfig.node.json` — Vite 配置文件 TypeScript 配置。
- Create: `frontend/vite.config.ts` — React、Vitest 和 jsdom 配置。
- Create: `frontend/eslint.config.js` — ESLint flat config。

### 前端代码

- Create: `frontend/src/main.tsx` — React 挂载入口。
- Create: `frontend/src/App.tsx` — AppInfo 三态页面。
- Create: `frontend/src/App.test.tsx` — Loading、Ready、Error、Retry 测试。
- Create: `frontend/src/types.ts` — `AppInfo` TypeScript 接口。
- Create: `frontend/src/services/appInfoClient.ts` — 生成 Binding 的薄封装。
- Create: `frontend/src/test/setup.ts` — jest-dom 初始化。
- Create: `frontend/src/styles.css` — 本地静态样式。
- Generate: `frontend/wailsjs/go/main/App.js` — Wails 自动生成，不手工编辑。
- Generate: `frontend/wailsjs/go/main/App.d.ts` — Wails 自动生成，不手工编辑。
- Generate: `frontend/wailsjs/go/models.ts` — Wails 自动生成，不手工编辑。

### 进度文档

- Modify: `docs/project-management/02-roadmap/backlog.md`
- Modify: `docs/project-management/04-progress/current-status.md`
- Modify: `docs/project-management/04-progress/milestone-status.md`
- Modify: `docs/project-management/04-progress/task-tracking.md`
- Modify: `docs/project-management/04-progress/implementation-trace.md`
- Modify: `docs/project-management/04-progress/weekly-log.md`
- Modify: `docs/project-management/04-progress/changelog.md`

---

### Task 1: 固定工具链并建立可编译的 Wails 项目边界

**Files:**
- Create: `go.mod`
- Create: `.nvmrc`
- Create: `.npmrc`
- Create: `.gitignore`
- Create: `wails.json`
- Generate: `build/**`

**Interfaces:**
- Consumes: 已确认的工程身份 `ItemAll`、`github.com/vvitem/item_all`、`com.vvitem.itemall`。
- Produces: 可供后续任务使用的 Go Module、Wails CLI、Node/pnpm 工具链和官方平台资源。

- [ ] **Step 1: 验证本机基础工具版本**

在 Windows PowerShell 运行：

```powershell
go version
node --version
corepack --version
git --version
```

Expected:

```text
go version go1.26.5 windows/amd64
v24.18.0
Corepack 输出一个可执行版本号
Git 输出一个可执行版本号
```

当 Go 或 Node 版本不一致时，先安装指定版本，不继续创建项目文件。

- [ ] **Step 2: 固定 pnpm 和 Wails CLI**

```powershell
corepack enable
corepack prepare pnpm@11.12.0 --activate
go install github.com/wailsapp/wails/v2/cmd/wails@v2.13.0
pnpm --version
wails version
```

Expected:

```text
11.12.0
Wails CLI v2.13.0
```

- [ ] **Step 3: 创建 Go Module 文件**

Create `go.mod`:

```go
module github.com/vvitem/item_all

go 1.26.0

toolchain go1.26.5

require github.com/wailsapp/wails/v2 v2.13.0
```

Run:

```powershell
go list -m
```

Expected:

```text
github.com/vvitem/item_all
```

- [ ] **Step 4: 创建 Node 和依赖精确版本策略**

Create `.nvmrc`:

```text
24.18.0
```

Create `.npmrc`:

```ini
save-exact=true
engine-strict=true
fund=false
audit=false
```

Run:

```powershell
Get-Content .nvmrc
Get-Content .npmrc
```

Expected: 内容与上述文件完全一致。

- [ ] **Step 5: 创建安全的忽略规则**

Create `.gitignore`:

```gitignore
# Go / Wails output
build/bin/
*.exe
*.app/

# Frontend output and dependencies
frontend/node_modules/
frontend/dist/
frontend/.vite/
frontend/coverage/

# Temporary scaffold
.wails-scaffold/

# Logs and local editor state
*.log
.vscode/
.idea/
.DS_Store
Thumbs.db
```

Run:

```powershell
git check-ignore frontend/node_modules/test build/bin/ItemAll.exe
```

Expected: 两个路径均被打印，表示忽略规则生效。

- [ ] **Step 6: 生成官方 Wails 平台资源**

```powershell
Remove-Item -Recurse -Force .wails-scaffold -ErrorAction SilentlyContinue
wails init -n ItemAll -d .wails-scaffold -t react-ts -q
Copy-Item .wails-scaffold\build . -Recurse -Force
Remove-Item -Recurse -Force .wails-scaffold
```

Run:

```powershell
Get-ChildItem build -Recurse | Select-Object FullName
```

Expected: 至少存在 `build/appicon.png`、`build/darwin` 和 `build/windows` 下的模板文件；仓库中不存在 `.wails-scaffold`。

- [ ] **Step 7: 创建 Wails v2 配置**

Create `wails.json`:

```json
{
  "$schema": "https://wails.io/schemas/config.v2.json",
  "name": "ItemAll",
  "outputfilename": "ItemAll",
  "frontend:install": "pnpm install --frozen-lockfile",
  "frontend:build": "pnpm run build",
  "frontend:dev:watcher": "pnpm run dev",
  "frontend:dev:serverUrl": "auto",
  "author": {
    "name": "vvitem",
    "email": ""
  }
}
```

Run:

```powershell
Get-Content wails.json | ConvertFrom-Json | Select-Object name, outputfilename
```

Expected:

```text
name    outputfilename
----    --------------
ItemAll ItemAll
```

- [ ] **Step 8: 检查本任务边界**

```powershell
Get-ChildItem -Recurse -File | Select-String -Pattern "applicationID|API_KEY|TOKEN|PASSWORD|BEGIN .*PRIVATE KEY" -CaseSensitive
```

Expected: 没有命中。`com.vvitem.itemall` 只记录在决策与设计文档中，不伪造 Wails v2 配置字段。

- [ ] **Step 9: 提交工具链和平台资源**

```powershell
git add go.mod .nvmrc .npmrc .gitignore wails.json build
git commit -m "build: initialise Wails toolchain"
```

Expected: 一个只包含工具链、Wails 配置和官方平台资源的提交。

---

### Task 2: 以 TDD 实现不可变构建信息

**Files:**
- Create: `internal/buildinfo/info_test.go`
- Create: `internal/buildinfo/info.go`

**Interfaces:**
- Consumes: Go Module `github.com/vvitem/item_all`。
- Produces: `buildinfo.Info` 和 `buildinfo.Current() Info`，供 `app.go` 使用。

- [ ] **Step 1: 写入失败测试**

Create `internal/buildinfo/info_test.go`:

```go
package buildinfo

import (
	"runtime"
	"testing"
)

func TestCurrentUsesStableDefaults(t *testing.T) {
	originalVersion := Version
	originalCommit := Commit
	originalBuildTime := BuildTime
	t.Cleanup(func() {
		Version = originalVersion
		Commit = originalCommit
		BuildTime = originalBuildTime
	})

	Version = ""
	Commit = ""
	BuildTime = ""

	got := Current()

	if got.Name != "ItemAll" {
		t.Fatalf("Name = %q, want ItemAll", got.Name)
	}
	if got.Version != "dev" {
		t.Fatalf("Version = %q, want dev", got.Version)
	}
	if got.Commit != "unknown" {
		t.Fatalf("Commit = %q, want unknown", got.Commit)
	}
	if got.BuildTime != "unknown" {
		t.Fatalf("BuildTime = %q, want unknown", got.BuildTime)
	}
	if got.GoVersion != runtime.Version() {
		t.Fatalf("GoVersion = %q, want %q", got.GoVersion, runtime.Version())
	}
	if got.OS != runtime.GOOS {
		t.Fatalf("OS = %q, want %q", got.OS, runtime.GOOS)
	}
	if got.Arch != runtime.GOARCH {
		t.Fatalf("Arch = %q, want %q", got.Arch, runtime.GOARCH)
	}
}

func TestCurrentPreservesInjectedValues(t *testing.T) {
	originalVersion := Version
	originalCommit := Commit
	originalBuildTime := BuildTime
	t.Cleanup(func() {
		Version = originalVersion
		Commit = originalCommit
		BuildTime = originalBuildTime
	})

	Version = "0.0.0-dev"
	Commit = "abc123def456"
	BuildTime = "2026-07-13T08:00:00Z"

	got := Current()

	if got.Version != "0.0.0-dev" {
		t.Fatalf("Version = %q", got.Version)
	}
	if got.Commit != "abc123def456" {
		t.Fatalf("Commit = %q", got.Commit)
	}
	if got.BuildTime != "2026-07-13T08:00:00Z" {
		t.Fatalf("BuildTime = %q", got.BuildTime)
	}
}
```

- [ ] **Step 2: 运行测试并确认失败**

```powershell
go test ./internal/buildinfo -run TestCurrent -v
```

Expected: FAIL，错误包含未定义的 `Version`、`Commit`、`BuildTime`、`Current` 或 `Info`。

- [ ] **Step 3: 写入最小实现**

Create `internal/buildinfo/info.go`:

```go
package buildinfo

import "runtime"

const Name = "ItemAll"

var (
	Version   = "dev"
	Commit    = "unknown"
	BuildTime = "unknown"
)

type Info struct {
	Name      string `json:"name"`
	Version   string `json:"version"`
	Commit    string `json:"commit"`
	BuildTime string `json:"buildTime"`
	GoVersion string `json:"goVersion"`
	OS        string `json:"os"`
	Arch      string `json:"arch"`
}

func Current() Info {
	return Info{
		Name:      Name,
		Version:   valueOrDefault(Version, "dev"),
		Commit:    valueOrDefault(Commit, "unknown"),
		BuildTime: valueOrDefault(BuildTime, "unknown"),
		GoVersion: runtime.Version(),
		OS:        runtime.GOOS,
		Arch:      runtime.GOARCH,
	}
}

func valueOrDefault(value string, fallback string) string {
	if value == "" {
		return fallback
	}
	return value
}
```

- [ ] **Step 4: 格式化并运行测试**

```powershell
gofmt -w internal/buildinfo/info.go internal/buildinfo/info_test.go
go test ./internal/buildinfo -run TestCurrent -v
```

Expected: 两个测试均 PASS。

- [ ] **Step 5: 运行竞态检测**

```powershell
go test -race ./internal/buildinfo
```

Expected: PASS，且没有 race report。

- [ ] **Step 6: 提交构建信息包**

```powershell
git add internal/buildinfo
git commit -m "feat: add immutable build information"
```

---

### Task 3: 以 TDD 实现只读 Wails AppInfo Binding

**Files:**
- Create: `app_test.go`
- Create: `app.go`
- Generate: `go.sum`
- Generate: `frontend/wailsjs/**`

**Interfaces:**
- Consumes: `buildinfo.Current() buildinfo.Info`。
- Produces: `NewApp() *App`、`(*App).GetAppInfo() AppInfo` 和前端生成 Binding `GetAppInfo(): Promise<AppInfo>`。

- [ ] **Step 1: 写入失败测试**

Create `app_test.go`:

```go
package main

import (
	"strings"
	"testing"

	"github.com/vvitem/item_all/internal/buildinfo"
)

func TestGetAppInfoReturnsSafeReadOnlyMetadata(t *testing.T) {
	originalVersion := buildinfo.Version
	originalCommit := buildinfo.Commit
	originalBuildTime := buildinfo.BuildTime
	t.Cleanup(func() {
		buildinfo.Version = originalVersion
		buildinfo.Commit = originalCommit
		buildinfo.BuildTime = originalBuildTime
	})

	buildinfo.Version = "0.0.0-dev"
	buildinfo.Commit = "abc123def456"
	buildinfo.BuildTime = "2026-07-13T08:00:00Z"

	got := NewApp().GetAppInfo()

	if got.Name != "ItemAll" {
		t.Fatalf("Name = %q", got.Name)
	}
	if got.Tagline != "Local-first SafeOps" {
		t.Fatalf("Tagline = %q", got.Tagline)
	}
	if got.Version != "0.0.0-dev" {
		t.Fatalf("Version = %q", got.Version)
	}
	if got.Commit != "abc123def456" {
		t.Fatalf("Commit = %q", got.Commit)
	}
	if got.BuildTime != "2026-07-13T08:00:00Z" {
		t.Fatalf("BuildTime = %q", got.BuildTime)
	}
	if got.Status != "ready" {
		t.Fatalf("Status = %q", got.Status)
	}
	if !strings.Contains(got.Runtime, "/") {
		t.Fatalf("Runtime = %q, want Go version and OS/arch", got.Runtime)
	}
}
```

- [ ] **Step 2: 运行测试并确认失败**

```powershell
go test . -run TestGetAppInfoReturnsSafeReadOnlyMetadata -v
```

Expected: FAIL，错误包含未定义的 `NewApp` 或 `AppInfo`。

- [ ] **Step 3: 写入最小 Binding 实现**

Create `app.go`:

```go
package main

import (
	"fmt"

	"github.com/vvitem/item_all/internal/buildinfo"
)

const tagline = "Local-first SafeOps"

const appStatusReady = "ready"

type App struct{}

type AppInfo struct {
	Name      string `json:"name"`
	Tagline   string `json:"tagline"`
	Version   string `json:"version"`
	Commit    string `json:"commit"`
	BuildTime string `json:"buildTime"`
	Runtime   string `json:"runtime"`
	Status    string `json:"status"`
}

func NewApp() *App {
	return &App{}
}

func (a *App) GetAppInfo() AppInfo {
	info := buildinfo.Current()
	return AppInfo{
		Name:      info.Name,
		Tagline:   tagline,
		Version:   info.Version,
		Commit:    info.Commit,
		BuildTime: info.BuildTime,
		Runtime:   fmt.Sprintf("%s %s/%s", info.GoVersion, info.OS, info.Arch),
		Status:    appStatusReady,
	}
}
```

- [ ] **Step 4: 整理依赖并运行全部 Go 测试**

```powershell
gofmt -w app.go app_test.go
go mod tidy
go test ./... -v
```

Expected: `app` 和 `internal/buildinfo` 测试全部 PASS，并生成 `go.sum`。

- [ ] **Step 5: 生成 Wails TypeScript Binding**

```powershell
wails generate module
```

Expected: `frontend/wailsjs/go/main/App.js`、`App.d.ts` 和 `models.ts` 存在；生成文件中包含 `GetAppInfo`。

Run:

```powershell
Get-ChildItem frontend\wailsjs -Recurse -File | Select-String -Pattern "GetAppInfo"
```

Expected: 至少命中 `App.js` 和 `App.d.ts`。

- [ ] **Step 6: 确认 Binding 没有额外能力**

```powershell
Get-Content frontend\wailsjs\go\main\App.d.ts
```

Expected: 只暴露 `GetAppInfo`，不包含文件系统、网络、Shell、数据库或环境变量方法。

- [ ] **Step 7: 提交 Binding 和模块锁定**

```powershell
git add app.go app_test.go go.mod go.sum frontend/wailsjs
git commit -m "feat: expose read-only app information"
```

---

### Task 4: 以 TDD 实现 React 三态页面

**Files:**
- Create: `frontend/package.json`
- Create: `frontend/pnpm-lock.yaml`
- Create: `frontend/index.html`
- Create: `frontend/tsconfig.json`
- Create: `frontend/tsconfig.node.json`
- Create: `frontend/vite.config.ts`
- Create: `frontend/eslint.config.js`
- Create: `frontend/src/types.ts`
- Create: `frontend/src/services/appInfoClient.ts`
- Create: `frontend/src/test/setup.ts`
- Create: `frontend/src/App.test.tsx`
- Create: `frontend/src/App.tsx`
- Create: `frontend/src/main.tsx`
- Create: `frontend/src/styles.css`

**Interfaces:**
- Consumes: 生成函数 `GetAppInfo()` 和 DTO 字段 `name`、`tagline`、`version`、`commit`、`buildTime`、`runtime`、`status`。
- Produces: 可测试的 `<App loadAppInfo={fn} />`，默认调用 Wails Binding；展示 Loading、Ready、Error 和 Retry。

- [ ] **Step 1: 创建精确版本的前端清单**

Create `frontend/package.json`:

```json
{
  "name": "itemall-frontend",
  "private": true,
  "version": "0.0.0",
  "type": "module",
  "packageManager": "pnpm@11.12.0",
  "engines": {
    "node": "24.18.0",
    "pnpm": "11.12.0"
  },
  "scripts": {
    "dev": "vite",
    "build": "tsc -b && vite build",
    "typecheck": "tsc -b --pretty false",
    "lint": "eslint . --max-warnings 0",
    "test": "vitest",
    "test:run": "vitest run",
    "preview": "vite preview"
  },
  "dependencies": {
    "react": "19.1.0",
    "react-dom": "19.1.0"
  },
  "devDependencies": {
    "@eslint/js": "9.30.1",
    "@testing-library/jest-dom": "6.6.3",
    "@testing-library/react": "16.3.0",
    "@types/react": "19.1.0",
    "@types/react-dom": "19.1.0",
    "@vitejs/plugin-react": "5.0.0",
    "eslint": "9.30.1",
    "eslint-plugin-react-hooks": "5.2.0",
    "eslint-plugin-react-refresh": "0.4.20",
    "globals": "16.3.0",
    "jsdom": "26.1.0",
    "typescript": "5.6.3",
    "typescript-eslint": "8.35.1",
    "vite": "7.0.0",
    "vitest": "3.2.4"
  }
}
```

Run:

```powershell
Set-Location frontend
pnpm install
pnpm install --frozen-lockfile
Set-Location ..
```

Expected: 两次安装成功，生成 `frontend/pnpm-lock.yaml`，第二次不修改锁文件。

- [ ] **Step 2: 创建 TypeScript、Vite、Vitest 和 ESLint 配置**

Create `frontend/tsconfig.json`:

```json
{
  "compilerOptions": {
    "target": "ES2022",
    "useDefineForClassFields": true,
    "lib": ["ES2022", "DOM", "DOM.Iterable"],
    "allowJs": false,
    "skipLibCheck": true,
    "esModuleInterop": true,
    "allowSyntheticDefaultImports": true,
    "strict": true,
    "forceConsistentCasingInFileNames": true,
    "module": "ESNext",
    "moduleResolution": "Bundler",
    "resolveJsonModule": true,
    "isolatedModules": true,
    "noEmit": true,
    "jsx": "react-jsx",
    "types": ["vitest/globals", "@testing-library/jest-dom"]
  },
  "include": ["src", "vite.config.ts"]
}
```

Create `frontend/tsconfig.node.json`:

```json
{
  "compilerOptions": {
    "composite": true,
    "skipLibCheck": true,
    "module": "ESNext",
    "moduleResolution": "Bundler",
    "allowImportingTsExtensions": true
  },
  "include": ["vite.config.ts", "eslint.config.js"]
}
```

Create `frontend/vite.config.ts`:

```ts
import { defineConfig } from 'vitest/config'
import react from '@vitejs/plugin-react'

export default defineConfig({
  plugins: [react()],
  test: {
    environment: 'jsdom',
    setupFiles: './src/test/setup.ts',
    globals: true,
    css: true,
  },
})
```

Create `frontend/eslint.config.js`:

```js
import js from '@eslint/js'
import globals from 'globals'
import reactHooks from 'eslint-plugin-react-hooks'
import reactRefresh from 'eslint-plugin-react-refresh'
import tseslint from 'typescript-eslint'

export default tseslint.config(
  { ignores: ['dist', 'coverage', 'wailsjs'] },
  {
    files: ['**/*.{ts,tsx}'],
    extends: [js.configs.recommended, ...tseslint.configs.recommended],
    languageOptions: {
      ecmaVersion: 2022,
      globals: globals.browser,
    },
    plugins: {
      'react-hooks': reactHooks,
      'react-refresh': reactRefresh,
    },
    rules: {
      ...reactHooks.configs.recommended.rules,
      'react-refresh/only-export-components': ['warn', { allowConstantExport: true }],
    },
  },
)
```

Create `frontend/src/test/setup.ts`:

```ts
import '@testing-library/jest-dom/vitest'
```

- [ ] **Step 3: 创建 DTO 和 Binding Client**

Create `frontend/src/types.ts`:

```ts
export interface AppInfo {
  name: string
  tagline: string
  version: string
  commit: string
  buildTime: string
  runtime: string
  status: string
}

export type LoadAppInfo = () => Promise<AppInfo>
```

Create `frontend/src/services/appInfoClient.ts`:

```ts
import { GetAppInfo } from '../../wailsjs/go/main/App'
import type { AppInfo } from '../types'

export async function getAppInfo(): Promise<AppInfo> {
  return (await GetAppInfo()) as AppInfo
}
```

- [ ] **Step 4: 写入失败的 UI 测试**

Create `frontend/src/App.test.tsx`:

```tsx
import { fireEvent, render, screen, waitFor } from '@testing-library/react'
import { describe, expect, it, vi } from 'vitest'
import App from './App'
import type { AppInfo } from './types'

const readyInfo: AppInfo = {
  name: 'ItemAll',
  tagline: 'Local-first SafeOps',
  version: '0.0.0-dev',
  commit: 'abc123def456',
  buildTime: '2026-07-13T08:00:00Z',
  runtime: 'go1.26.5 windows/amd64',
  status: 'ready',
}

describe('App', () => {
  it('shows loading while app information is pending', () => {
    const pending = new Promise<AppInfo>(() => undefined)

    render(<App loadAppInfo={() => pending} />)

    expect(screen.getByRole('status')).toHaveTextContent('正在读取应用信息')
  })

  it('shows trusted app information when loading succeeds', async () => {
    render(<App loadAppInfo={async () => readyInfo} />)

    expect(await screen.findByRole('heading', { name: 'ItemAll' })).toBeInTheDocument()
    expect(screen.getByText('Local-first SafeOps')).toBeInTheDocument()
    expect(screen.getByText('0.0.0-dev')).toBeInTheDocument()
    expect(screen.getByText('abc123def456')).toBeInTheDocument()
    expect(screen.getByText('go1.26.5 windows/amd64')).toBeInTheDocument()
    expect(screen.getByText('运行正常')).toBeInTheDocument()
  })

  it('shows a safe error without exposing the thrown message', async () => {
    render(<App loadAppInfo={async () => Promise.reject(new Error('C:\\Users\\secret\\token.txt'))} />)

    expect(await screen.findByRole('alert')).toHaveTextContent('无法读取应用信息')
    expect(screen.queryByText(/secret|token\.txt/i)).not.toBeInTheDocument()
  })

  it('retries after a failed request', async () => {
    const loadAppInfo = vi
      .fn<() => Promise<AppInfo>>()
      .mockRejectedValueOnce(new Error('first failure'))
      .mockResolvedValueOnce(readyInfo)

    render(<App loadAppInfo={loadAppInfo} />)

    const retryButton = await screen.findByRole('button', { name: '重新读取' })
    fireEvent.click(retryButton)

    await waitFor(() => expect(loadAppInfo).toHaveBeenCalledTimes(2))
    expect(await screen.findByRole('heading', { name: 'ItemAll' })).toBeInTheDocument()
  })
})
```

- [ ] **Step 5: 运行测试并确认失败**

```powershell
Set-Location frontend
pnpm test:run -- App.test.tsx
Set-Location ..
```

Expected: FAIL，因为 `frontend/src/App.tsx` 尚不存在。

- [ ] **Step 6: 写入最小 React 实现**

Create `frontend/src/App.tsx`:

```tsx
import { useEffect, useState } from 'react'
import { getAppInfo } from './services/appInfoClient'
import type { AppInfo, LoadAppInfo } from './types'

type ViewState =
  | { kind: 'loading' }
  | { kind: 'ready'; info: AppInfo }
  | { kind: 'error' }

interface AppProps {
  loadAppInfo?: LoadAppInfo
}

export default function App({ loadAppInfo = getAppInfo }: AppProps) {
  const [attempt, setAttempt] = useState(0)
  const [state, setState] = useState<ViewState>({ kind: 'loading' })

  useEffect(() => {
    let active = true
    setState({ kind: 'loading' })

    loadAppInfo()
      .then((info) => {
        if (active) {
          setState({ kind: 'ready', info })
        }
      })
      .catch(() => {
        if (active) {
          setState({ kind: 'error' })
        }
      })

    return () => {
      active = false
    }
  }, [attempt, loadAppInfo])

  if (state.kind === 'loading') {
    return (
      <main className="shell shell--centered">
        <p role="status" className="status-card">
          正在读取应用信息
        </p>
      </main>
    )
  }

  if (state.kind === 'error') {
    return (
      <main className="shell shell--centered">
        <section role="alert" className="status-card status-card--error">
          <h1>无法读取应用信息</h1>
          <p>应用未能完成本地初始化，请重新读取。</p>
          <button type="button" onClick={() => setAttempt((value) => value + 1)}>
            重新读取
          </button>
        </section>
      </main>
    )
  }

  const { info } = state

  return (
    <main className="shell">
      <header className="hero">
        <div>
          <p className="eyebrow">Local desktop operations workspace</p>
          <h1>{info.name}</h1>
          <p className="tagline">{info.tagline}</p>
        </div>
        <span className="health-badge">运行正常</span>
      </header>

      <section className="metadata" aria-label="应用构建信息">
        <article>
          <span>Version</span>
          <strong>{info.version}</strong>
        </article>
        <article>
          <span>Commit</span>
          <strong>{info.commit}</strong>
        </article>
        <article>
          <span>Build Time</span>
          <strong>{info.buildTime}</strong>
        </article>
        <article>
          <span>Runtime</span>
          <strong>{info.runtime}</strong>
        </article>
      </section>
    </main>
  )
}
```

Create `frontend/src/main.tsx`:

```tsx
import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import App from './App'
import './styles.css'

const root = document.getElementById('root')

if (root === null) {
  throw new Error('ItemAll root element is missing')
}

createRoot(root).render(
  <StrictMode>
    <App />
  </StrictMode>,
)
```

Create `frontend/index.html`:

```html
<!doctype html>
<html lang="zh-CN">
  <head>
    <meta charset="UTF-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0" />
    <meta name="color-scheme" content="dark" />
    <title>ItemAll</title>
  </head>
  <body>
    <div id="root"></div>
    <script type="module" src="/src/main.tsx"></script>
  </body>
</html>
```

- [ ] **Step 7: 写入本地静态样式**

Create `frontend/src/styles.css`:

```css
:root {
  font-family:
    Inter, ui-sans-serif, system-ui, -apple-system, BlinkMacSystemFont, "Segoe UI",
    sans-serif;
  color: #e2e8f0;
  background: #020617;
  font-synthesis: none;
  text-rendering: optimizeLegibility;
}

* {
  box-sizing: border-box;
}

body {
  min-width: 320px;
  min-height: 100vh;
  margin: 0;
  background:
    radial-gradient(circle at top left, rgba(14, 165, 233, 0.18), transparent 38%),
    #020617;
}

button {
  border: 1px solid #38bdf8;
  border-radius: 0.75rem;
  padding: 0.75rem 1rem;
  color: #e0f2fe;
  background: #075985;
  cursor: pointer;
}

button:hover {
  background: #0369a1;
}

.shell {
  width: min(960px, calc(100% - 3rem));
  margin: 0 auto;
  padding: 4rem 0;
}

.shell--centered {
  min-height: 100vh;
  display: grid;
  place-items: center;
  padding: 2rem 0;
}

.hero {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 2rem;
  padding: 2rem;
  border: 1px solid #1e293b;
  border-radius: 1.25rem;
  background: rgba(15, 23, 42, 0.9);
}

.eyebrow {
  margin: 0 0 0.75rem;
  color: #7dd3fc;
  font-size: 0.8rem;
  font-weight: 700;
  letter-spacing: 0.12em;
  text-transform: uppercase;
}

h1 {
  margin: 0;
  font-size: clamp(2.5rem, 7vw, 5rem);
  line-height: 1;
}

.tagline {
  margin: 1rem 0 0;
  color: #94a3b8;
  font-size: 1.1rem;
}

.health-badge {
  border: 1px solid #22c55e;
  border-radius: 999px;
  padding: 0.5rem 0.75rem;
  color: #bbf7d0;
  background: rgba(22, 101, 52, 0.35);
  white-space: nowrap;
}

.metadata {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 1rem;
  margin-top: 1rem;
}

.metadata article,
.status-card {
  padding: 1.25rem;
  border: 1px solid #1e293b;
  border-radius: 1rem;
  background: rgba(15, 23, 42, 0.85);
}

.metadata span {
  display: block;
  margin-bottom: 0.5rem;
  color: #94a3b8;
  font-size: 0.8rem;
  text-transform: uppercase;
}

.metadata strong {
  overflow-wrap: anywhere;
}

.status-card--error {
  max-width: 32rem;
}

.status-card--error h1 {
  font-size: 2rem;
}

@media (max-width: 640px) {
  .hero {
    flex-direction: column;
  }

  .metadata {
    grid-template-columns: 1fr;
  }
}
```

- [ ] **Step 8: 运行前端测试**

```powershell
Set-Location frontend
pnpm test:run -- App.test.tsx
Set-Location ..
```

Expected: 4 个测试全部 PASS。

- [ ] **Step 9: 运行前端静态验证**

```powershell
Set-Location frontend
pnpm typecheck
pnpm lint
pnpm build
Set-Location ..
```

Expected: 三个命令均成功，`frontend/dist/index.html` 存在。

- [ ] **Step 10: 提交前端三态页面**

```powershell
git add frontend
git commit -m "feat: add app information screen"
```

---

### Task 5: 装配桌面应用、构建注入和开发命令

**Files:**
- Create: `main.go`
- Create: `Makefile`
- Modify: `README.md`
- Modify: `build/windows/info.json`

**Interfaces:**
- Consumes: `NewApp()`、`frontend/dist`、Wails v2 配置和 buildinfo ldflags 字段。
- Produces: 可执行的 `wails dev` 与 `wails build`，以及可复现的开发文档。

- [ ] **Step 1: 写入 Wails Composition Root**

Create `main.go`:

```go
package main

import (
	"embed"
	"log"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	app := NewApp()

	err := wails.Run(&options.App{
		Title:     "ItemAll",
		Width:     1024,
		Height:    700,
		MinWidth:  800,
		MinHeight: 560,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 2, G: 6, B: 23, A: 1},
		Bind: []interface{}{
			app,
		},
	})
	if err != nil {
		log.Fatal(err)
	}
}
```

Run:

```powershell
gofmt -w main.go
go mod tidy
go test ./...
```

Expected: PASS；`go test` 同时验证根包能够编译并成功嵌入 `frontend/dist`。

- [ ] **Step 2: 写入可选 Makefile**

Create `Makefile`:

```makefile
VERSION ?= 0.0.0-dev
COMMIT ?= $(shell git rev-parse --short=12 HEAD)
BUILD_TIME ?= $(shell date -u +%Y-%m-%dT%H:%M:%SZ)
LDFLAGS = -X github.com/vvitem/item_all/internal/buildinfo.Version=$(VERSION) -X github.com/vvitem/item_all/internal/buildinfo.Commit=$(COMMIT) -X github.com/vvitem/item_all/internal/buildinfo.BuildTime=$(BUILD_TIME)

.PHONY: test frontend-test frontend-build dev build

test:
	go test ./...
	cd frontend && pnpm typecheck && pnpm lint && pnpm test:run

frontend-test:
	cd frontend && pnpm test:run

frontend-build:
	cd frontend && pnpm build

dev:
	wails dev

build:
	wails build -clean -trimpath -ldflags "$(LDFLAGS)"
```

- [ ] **Step 3: 更新 Windows 产品元数据**

Open `build/windows/info.json` and set the product-facing fields to the following values while preserving the template schema:

```json
{
  "fixed": {
    "file_version": "0.0.0.0",
    "product_version": "0.0.0.0"
  },
  "info": {
    "0000": {
      "CompanyName": "vvitem",
      "FileDescription": "ItemAll Local-first SafeOps",
      "FileVersion": "0.0.0-dev",
      "InternalName": "ItemAll",
      "LegalCopyright": "Copyright © 2026 vvitem",
      "OriginalFilename": "ItemAll.exe",
      "ProductName": "ItemAll",
      "ProductVersion": "0.0.0-dev"
    }
  }
}
```

Run:

```powershell
Get-Content build\windows\info.json | ConvertFrom-Json | Out-Null
```

Expected: JSON 解析成功。

- [ ] **Step 4: 更新 README 开发说明**

Append the following sections to `README.md`, preserving the existing product introduction:

```markdown
## Development baseline

- Go 1.26.5
- Wails CLI v2.13.0
- Node.js 24.18.0 LTS
- pnpm 11.12.0
- Windows 10/11 with WebView2

## Install tools

```powershell
corepack enable
corepack prepare pnpm@11.12.0 --activate
go install github.com/wailsapp/wails/v2/cmd/wails@v2.13.0
wails doctor
```

## Install dependencies

```powershell
go mod download
Set-Location frontend
pnpm install --frozen-lockfile
Set-Location ..
```

## Validate

```powershell
go test ./...
Set-Location frontend
pnpm typecheck
pnpm lint
pnpm test:run
pnpm build
Set-Location ..
```

## Run in development

```powershell
wails dev
```

## Build on Windows

```powershell
$commit = git rev-parse --short=12 HEAD
$buildTime = (Get-Date).ToUniversalTime().ToString("yyyy-MM-ddTHH:mm:ssZ")
wails build -clean -trimpath -ldflags "-X github.com/vvitem/item_all/internal/buildinfo.Version=0.0.0-dev -X github.com/vvitem/item_all/internal/buildinfo.Commit=$commit -X github.com/vvitem/item_all/internal/buildinfo.BuildTime=$buildTime"
```

The binary is written to `build/bin/ItemAll.exe`.
```

- [ ] **Step 5: 运行全量自动验证**

```powershell
go test ./...
Set-Location frontend
pnpm install --frozen-lockfile
pnpm typecheck
pnpm lint
pnpm test:run
pnpm build
Set-Location ..
```

Expected: 所有命令成功；测试总数至少为 Go 4 个、前端 4 个。

- [ ] **Step 6: 执行 Windows 生产构建**

```powershell
$commit = git rev-parse --short=12 HEAD
$buildTime = (Get-Date).ToUniversalTime().ToString("yyyy-MM-ddTHH:mm:ssZ")
wails build -clean -trimpath -ldflags "-X github.com/vvitem/item_all/internal/buildinfo.Version=0.0.0-dev -X github.com/vvitem/item_all/internal/buildinfo.Commit=$commit -X github.com/vvitem/item_all/internal/buildinfo.BuildTime=$buildTime"
Get-Item build\bin\ItemAll.exe
```

Expected: `build/bin/ItemAll.exe` 存在且文件大小大于 0。

- [ ] **Step 7: 执行 Windows 开发启动验证**

```powershell
wails dev
```

Manual acceptance:

1. 窗口标题为 `ItemAll`。
2. 页面显示 `ItemAll` 和 `Local-first SafeOps`。
3. 状态为“运行正常”。
4. Version、Commit、Build Time、Runtime 均可见。
5. 窗口缩小到最小宽高后内容仍可访问。
6. 开发者工具控制台没有未处理 Promise rejection。
7. 应用不发起网络请求。

关闭应用后，记录 Windows 版本、Go 版本、Wails 版本和验证日期到实现 PR 描述。

- [ ] **Step 8: 扫描秘密和范围外能力**

```powershell
Get-ChildItem -Recurse -File -Exclude pnpm-lock.yaml,go.sum | Select-String -Pattern "BEGIN .*PRIVATE KEY|api[_-]?key|access[_-]?token|password\s*=|ssh\.exec|db\.execute|shell\.run" -CaseSensitive:$false
```

Expected: 没有真实秘密，也没有范围外工具接口。文档中描述禁止项的文字命中不算代码违规，必须人工确认命中路径。

- [ ] **Step 9: 提交桌面装配和开发文档**

```powershell
git add main.go Makefile README.md build/windows/info.json go.mod go.sum
git commit -m "build: wire ItemAll desktop application"
```

---

### Task 6: 同步进度、创建实现 PR 并在合并后闭环

**Files:**
- Modify: `docs/project-management/02-roadmap/backlog.md`
- Modify: `docs/project-management/04-progress/current-status.md`
- Modify: `docs/project-management/04-progress/milestone-status.md`
- Modify: `docs/project-management/04-progress/task-tracking.md`
- Modify: `docs/project-management/04-progress/implementation-trace.md`
- Modify: `docs/project-management/04-progress/weekly-log.md`
- Modify: `docs/project-management/04-progress/changelog.md`

**Interfaces:**
- Consumes: 全部测试结果、Windows 启动记录、实现 Commit 和 Issue #2。
- Produces: 状态真实的实现 PR；合并后 `M0-002=DONE`、`M0-003=READY`。

- [ ] **Step 1: 在实现分支将 M0-002 标记为 IN_REVIEW**

Apply these exact state rules before opening the implementation PR:

```text
M0-001 = DONE
M0-002 = IN_REVIEW
M0-003 = NOT_STARTED
全项目：DONE=1, IN_REVIEW=1, NOT_STARTED=48
M0-foundation：DONE=1, IN_REVIEW=1, 总任务=10
```

`current-status.md` 的“当前正在进行”必须写明：

```markdown
- `M0-002`：最小 Go/Wails/React 工程骨架已实现，正在通过实现 PR 评审。
- 完成证据：Go/前端测试、Windows `wails dev`、Windows `wails build` 和实现 PR。
```

`implementation-trace.md` 增加：

```markdown
| ItemAll Scaffold | M0-002 | 2026-07-13-m0-002-project-scaffold-design.md | `main.go`, `app.go`, `internal/buildinfo`, `frontend/src` | `app_test.go`, `internal/buildinfo/info_test.go`, `frontend/src/App.test.tsx` | Issue #2 / 实现 PR | IN_REVIEW |
```

- [ ] **Step 2: 提交进度更新**

```powershell
git add docs/project-management
git commit -m "docs: track M0-002 implementation review"
```

- [ ] **Step 3: 最终验证干净工作树**

```powershell
go test ./...
Set-Location frontend
pnpm install --frozen-lockfile
pnpm typecheck
pnpm lint
pnpm test:run
pnpm build
Set-Location ..
git status --short
```

Expected: 所有命令成功，`git status --short` 没有输出。

- [ ] **Step 4: 推送实现分支并创建 Draft PR**

Use branch:

```text
feat/m0-002-project-scaffold
```

Use PR title:

```text
feat: initialise ItemAll desktop scaffold
```

PR body must include:

```markdown
## Scope

- Go 1.26.5 + Wails v2.13.0 root project
- React/TypeScript/Vite frontend
- read-only `GetAppInfo()` Binding
- Loading, Ready, Error and Retry UI states
- exact dependency lock files
- Windows local development and production build validation

## Validation

- `go test ./...`
- `pnpm install --frozen-lockfile`
- `pnpm typecheck`
- `pnpm lint`
- `pnpm test:run`
- `pnpm build`
- `wails dev` on Windows
- `wails build -clean -trimpath` on Windows

## Boundaries

No SSH, database, SQLite business tables, AI, Operation Bus, cloud services, credentials or telemetry.

Closes #2 after merge and post-merge status synchronization.
```

- [ ] **Step 5: 评审实现 PR**

Review requirements:

1. Changed files match this plan.
2. Generated `wailsjs` files are not hand-edited.
3. No dependencies use `latest`, `*`, caret or tilde ranges.
4. No environment variables, secrets or remote assets appear.
5. `GetAppInfo()` is the only Binding.
6. All automatic checks and Windows manual checks have evidence.
7. M0-002 remains `IN_REVIEW` until merge.

- [ ] **Step 6: 合并实现 PR 后创建状态同步分支**

After squash merge, create a docs-only branch from the new `main` and apply:

```text
M0-002: IN_REVIEW -> DONE
M0-003: NOT_STARTED -> READY
全项目：DONE=2, READY=1, NOT_STARTED=47
M0-foundation：DONE=2/10 = 20%
```

Update `implementation-trace.md` with the actual PR number and merge Commit, and change Status to `DONE`.

- [ ] **Step 7: 合并状态同步 PR并关闭 Issue #2**

Use commit and PR title:

```text
docs: complete M0-002 and ready M0-003
```

Close Issue #2 only after the docs-only synchronization PR is merged. The closing comment must include:

```markdown
M0-002 已完成并有以下证据：

- 实现 PR：<actual PR number>
- Merge Commit：<actual merge SHA>
- Go 与前端自动测试通过
- Windows `wails dev` 验证通过
- Windows `wails build` 生成 `build/bin/ItemAll.exe`
- 项目进度文档已同步
```

---

## Final Verification Matrix

| Requirement | Evidence |
|---|---|
| 产品身份固定 | `go.mod`, `wails.json`, AppInfo, README |
| Wails 标准根目录 | `main.go`, `app.go`, `build/`, `frontend/` |
| Windows 开发启动 | `wails dev` 手工记录 |
| Windows 生产构建 | `build/bin/ItemAll.exe` 和构建命令 |
| 构建信息展示 | `internal/buildinfo`, `GetAppInfo`, UI |
| Loading/Ready/Error/Retry | `frontend/src/App.test.tsx` |
| Go 测试 | `go test ./...` |
| Frontend typecheck | `pnpm typecheck` |
| Frontend lint | `pnpm lint` |
| Frontend tests | `pnpm test:run` |
| Frontend build | `pnpm build` |
| 精确依赖 | `frontend/package.json`, `pnpm-lock.yaml`, `go.sum` |
| 无范围外能力 | PR changed-files review and secret/scope scan |
| 进度真实 | M0-002 merge前 `IN_REVIEW`，merge后 `DONE` |

## Self-Review Result

- Spec coverage: 已覆盖工程身份、根目录结构、构建信息、AppInfo、三态 UI、错误安全、精确依赖、Windows 开发与生产构建、范围边界和进度同步。
- Placeholder scan: 没有未定义的实现占位项；需要替换的 PR 号和 Merge SHA 仅存在于合并后证据模板，并明确由实际 GitHub 结果填入。
- Type consistency: Go `AppInfo` 与 TypeScript `AppInfo` 字段保持 `name`、`tagline`、`version`、`commit`、`buildTime`、`runtime`、`status` 一致。
- Scope: 本计划只实现 `M0-002`；`M0-003` CI、SQLite、Keychain、Operation Bus、SSH、数据库和 AI 均未混入。
