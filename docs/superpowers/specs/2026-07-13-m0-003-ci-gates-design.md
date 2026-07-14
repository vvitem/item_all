# M0-003 分层 CI 门禁设计

> 状态：待评审  
> 负责人：vvitem  
> 最后更新：2026-07-13  
> 基线 Commit：`351ca3b3466ff347948dec444716720256e8c2c7`  
> 关联 Milestone：`M0-foundation`  
> 关联 Issue/PR：[Issue #7](https://github.com/vvitem/item_all/issues/7)

## 1. 目标

为 ItemAll 建立正式、可复现、默认安全的 Pull Request CI 门禁，使任何进入 `main` 的改动都经过以下验证：

1. Go 格式、模块完整性、静态检查、单元测试和关键竞态检测。
2. Frontend 冻结依赖安装、Vitest、TypeScript、ESLint 和 production build。
3. Wails 生成 Binding 与仓库提交结果无漂移。
4. Go/pnpm 依赖锁文件、精确版本和安装脚本授权策略未被绕过。
5. 仓库不存在真实凭据、私钥、危险远程资源或 M0-003 范围外能力。
6. Windows 可完成 Wails production build，并生成可启动的非空 `ItemAll.exe`。
7. 门禁能够通过一次可控 RED→GREEN 演练证明失败时确实阻断合并。

CI 是代码进入 `main` 的质量和安全底线，不是发布系统。本任务不创建安装包、签名、Release、部署或遥测流程。

## 2. 已确认方案

采用 **方案 B：分层 CI**。

单个 `.github/workflows/ci.yml` 包含三个并行 Job：

```text
pull_request / push(main) / workflow_dispatch
                    │
                    ├── quality (Ubuntu 24.04)
                    ├── generated-and-security (Ubuntu 24.04)
                    └── windows-build (Windows Server 2025)
```

三个 Job 都是独立必需门禁：

- `quality` 负责快速、可读的代码质量反馈。
- `generated-and-security` 负责供应链、生成文件和安全边界。
- `windows-build` 负责桌面应用的真实目标平台构建事实。

不使用单一 Windows Job，也不在 M0-003 提前抽取 reusable workflow 或 composite action。

## 3. 范围与明确不做

### 3.1 本任务交付

- `.github/workflows/ci.yml`
- `docs/development/ci.md`
- 必要的本地验证脚本；优先使用跨平台 Go/Node 脚本，避免复制大量 Bash/PowerShell 逻辑
- 对 README 的 CI 入口说明
- 一次故意失败和恢复绿色的 Actions 运行证据
- 合并后的项目进度同步

### 3.2 明确不做

- SSH、SFTP、数据库、SQLite 业务表、Operation Bus、AI 或凭据能力
- Release、安装包、签名、校验和或 GitHub Release
- macOS/Linux 桌面 production build
- 自托管 Runner
- 上传可执行文件作为长期发布产物
- CodeQL、SBOM、Artifact Attestation 或 OpenSSF Scorecard；这些在功能面扩大后独立评估
- Dependabot 配置；依赖更新流程另开任务
- 覆盖率百分比门槛；M0 阶段以风险路径测试为准

## 4. 设计原则

### 4.1 快速反馈与真实平台分离

Go/Frontend 常规检查放在 Ubuntu，避免每个小改动都串行等待 Windows Wails 构建。Windows Job 独立验证最终桌面构建，不承担所有静态检查。

### 4.2 默认最小权限

Workflow 顶层固定：

```yaml
permissions:
  contents: read
```

指定 `contents: read` 后，其余未声明权限为 `none`。CI 不需要写仓库、评论 PR、上传安全事件或访问环境 Secret。

### 4.3 不信任 PR 内容

- 使用 `pull_request`，禁止 `pull_request_target`。
- 不读取仓库 Secret。
- 不执行来自 PR 标题、分支名、Issue 文本的动态 Shell 片段。
- `actions/checkout` 设置 `persist-credentials: false`。
- 所有第三方输入通过固定参数传递，不拼接未经处理的 GitHub Context。

### 4.4 Action 依赖不可变

所有 `uses:` 必须固定到来源仓库的 **完整 Commit SHA**，并在行尾注释对应发布版本，例如：

```yaml
uses: actions/checkout@<40-char-sha> # v7.0.0
```

实现时采用已经发布并核验的稳定版本；不得使用 `@main`、`@master` 或仅使用可移动大版本 Tag。当前设计核验的官方稳定线包括：

- `actions/checkout` v7.0.0
- `actions/setup-go` v6.5.0
- `actions/setup-node` v7.0.0
- `actions/cache` v6.1.0（仅在确有必要时使用）

具体完整 SHA 在实现计划中冻结，并从官方 Action 仓库 Release/Tag 解析，不从 Fork 获取。

### 4.5 Runner 标签明确

使用：

- `ubuntu-24.04`
- `windows-2025`

不使用 `ubuntu-latest` 或 `windows-latest`，避免 GitHub 迁移 `latest` 指向时无审查改变构建环境。Runner 镜像仍会按周更新，因此工具链必须由 setup action 和仓库文件显式固定。

## 5. Workflow 触发模型

```yaml
on:
  pull_request:
    branches:
      - main
  push:
    branches:
      - main
  workflow_dispatch:
```

### 5.1 不使用路径过滤

本 Workflow 不配置 `paths` 或 `paths-ignore`。

原因：当分支保护将 Job 设为 required checks 时，被路径过滤跳过的 Workflow 可能长期保持 Pending，反而阻止合并。所有 PR 都运行三个门禁，确保状态名称稳定。

### 5.2 Draft PR

Draft PR 也运行 CI。原因：

- 尽早发现工具链和供应链问题。
- 避免转为 Ready 后才第一次触发耗时 Windows 构建。
- 当前单人仓库规模较小，成本可控。

### 5.3 并发取消

Workflow 顶层配置：

```yaml
concurrency:
  group: ${{ github.workflow }}-${{ github.event.pull_request.number || github.ref }}
  cancel-in-progress: true
```

同一 PR 新提交会取消旧运行；`main` push 使用 ref 作为回退。组名包含 workflow 名，避免未来其他 Workflow 互相取消。

## 6. Job 一：`quality`

### 6.1 运行环境

- Runner：`ubuntu-24.04`
- 超时：15 分钟
- 权限：继承顶层 `contents: read`

### 6.2 工具链

- Go：读取 `go.mod` 的 `toolchain go1.26.5`
- Node：读取 `.nvmrc` 的 `24.18.0`
- pnpm：读取 `frontend/package.json` 的 `packageManager: pnpm@11.12.0`
- Corepack：启用后验证实际 pnpm 版本与声明一致

版本不一致必须立即失败，不允许 Runner 预装版本悄悄替代仓库声明。

### 6.3 Go 检查顺序

1. `go version`
2. `go env GOPATH GOMODCACHE GOCACHE`
3. `go mod download`
4. `go mod verify`
5. `gofmt -l` 检查；发现任何 Go 文件需要格式化即失败
6. `go vet ./...`
7. `go test ./...`
8. `go test -race ./internal/buildinfo`

Ubuntu 24.04 上编译 Wails 根包需要 GTK/WebKit 平台依赖。Quality Job 安装固定的系统包集合，并对 Ubuntu 24.04 使用 Wails `webkit2_41` Build Tag。系统依赖安装命令必须集中在单一步骤中，并在 CI 文档说明其原因。

若后续 Wails/Ubuntu 变化导致根包无法稳定测试，可以将纯 Go 测试和桌面根包编译拆开，但不得静默删除 `go test ./...` 门禁。

### 6.4 Frontend 检查顺序

1. 验证 Node 与 pnpm 精确版本
2. `pnpm install --frozen-lockfile`
3. `pnpm test:run`
4. `pnpm typecheck`
5. `pnpm lint`
6. `pnpm build`

测试失败、TypeScript 错误、ESLint warning 或构建错误均使 Job 失败。

### 6.5 缓存

优先使用 `setup-go` 和 `setup-node` 的官方缓存能力：

- Go cache key 由 `go.sum` 驱动
- pnpm cache key 由 `frontend/pnpm-lock.yaml` 驱动

缓存只用于加速，不作为构建输入事实：

- 缓存未命中必须仍能成功。
- 不缓存 `frontend/dist`、测试结果或生成 Binding。
- PR 不允许把缓存内容提交回仓库。

## 7. Job 二：`generated-and-security`

### 7.1 运行环境

- Runner：`ubuntu-24.04`
- 超时：15 分钟
- 权限：`contents: read`

### 7.2 生成文件漂移

安装固定的 Wails CLI：

```text
github.com/wailsapp/wails/v2/cmd/wails@v2.13.0
```

执行 Binding 重新生成后检查：

- `frontend/wailsjs/**`
- `go.mod`
- `go.sum`
- `frontend/pnpm-lock.yaml`

验证方式使用 `git diff --exit-code -- <allowlisted paths>`。

CI 不自动提交生成结果。发现漂移时输出明确修复命令并失败，开发者必须在本地重新生成后提交。

### 7.3 依赖政策

验证：

- `packageManager` 精确为 `pnpm@11.12.0`
- `engines.node` 与 `.nvmrc` 一致
- `package.json` 依赖不得使用 `latest`、`*`、`^`、`~`、`workspace:*` 或无界 Git URL
- `pnpm install --frozen-lockfile` 不修改锁文件
- `pnpm-workspace.yaml` 的 `allowBuilds` 只允许经过评审的包；M0-003 基线仅允许 `esbuild`
- `go.mod` Module Path、Go 版本、toolchain 和 Wails 版本与工程决策一致
- `go mod tidy` 后无漂移

依赖政策检查优先编写为仓库内可测试脚本，避免在 YAML 中维护复杂正则。

### 7.4 Secret Scan

使用固定版本的 Gitleaks CLI，直接下载官方 Release 二进制并校验 SHA256，或使用固定完整 SHA 的官方 Action。推荐 CLI 方式，减少第三方 Action 的 Token 与运行时表面积。

扫描范围：

```text
gitleaks git --redact --no-banner
```

要求：

- 扫描完整可访问 Git 历史，而非只扫工作树。
- 输出必须脱敏。
- 禁止在 Workflow 中加入真实测试 Secret。
- 故意失败演练使用明显的假测试文件和单独分支，验证后必须删除，不得进入 `main`。

### 7.5 静态安全边界扫描

M0-003 继续保留 M0-002 已验证的边界：

- 禁止 `.env`、私钥和证书私钥文件进入仓库
- 禁止前端使用远程 script、stylesheet、font 或 CDN URL
- 禁止新增 SSH、数据库、AI、云 SDK、遥测或环境变量读取能力
- Wails 生成 API 在 M0-003 完成时仍只允许 `GetAppInfo()`

范围扫描不得只依赖一个关键词正则。设计采用两层：

1. 明确文件/依赖 Allowlist。
2. 高风险关键词启发式扫描，命中后失败并要求人工确认或更新允许规则。

未来真正开始 M0-004/M0-006/M0-007 等模块时，必须在对应设计 PR 中更新规则，不能为了通过 CI 直接全局关闭扫描。

### 7.6 Workflow 自身安全检查

检查 `.github/workflows/*.yml`：

- `permissions` 明确存在且没有 `write-all`
- 禁止 `pull_request_target`
- `uses:` 全部为完整 40 字符 SHA，并带版本注释
- 禁止未经环境变量中转直接把不可信 GitHub Context 插入 `run:`
- 禁止长生命周期 Secret 或云凭据

## 8. Job 三：`windows-build`

### 8.1 运行环境

- Runner：`windows-2025`
- 超时：25 分钟
- 权限：`contents: read`

Windows Server 2025 是 M0-003 的真实桌面构建平台。CI 不使用自托管 Windows 主机。

### 8.2 工具链

- Go 1.26.5
- Node 24.18.0
- pnpm 11.12.0
- Wails CLI v2.13.0

Wails CLI 使用精确版本 `go install`。可以缓存 Go Module、Go Build Cache 和 Wails CLI 二进制，但缓存 Key 必须包含：

- Runner OS
- Go toolchain
- Wails 版本
- `go.sum` Hash

### 8.3 验证顺序

1. Checkout，`persist-credentials: false`
2. 安装固定 Go/Node/pnpm/Wails
3. `pnpm install --frozen-lockfile`
4. `go test ./...`
5. `wails build -clean -trimpath`
6. 验证 `build/bin/ItemAll.exe` 存在且大小大于 0
7. 启动 EXE，确认进程在最小观察窗口内没有立即异常退出
8. 清理进程树

### 8.4 烟测定义

烟测只证明：

- 可执行文件能够被 Windows 加载
- Wails/WebView2 初始化没有立即崩溃
- 应用进程在观察窗口内保持运行

烟测不声称验证视觉布局、交互或截图。人工 UI 视觉检查仍属于 PR 评审或后续 Windows E2E。

### 8.5 产物策略

M0-003 默认不上传 EXE Artifact：

- CI 只验证文件存在与启动。
- 正式分发、保留、校验和和签名属于 M4-007。
- 当构建故障需要诊断时，可临时上传日志，不上传含本地路径或敏感内容的整包。

## 9. Failure Semantics

任何 Job 失败都阻止合并；禁止使用以下方式弱化门禁：

- `continue-on-error: true`
- `|| true`
- 忽略退出码
- 把失败命令包装为 warning
- 只在 `main` push 后验证而不验证 PR
- 仅上传报告但让 Job 绿色

每个失败步骤必须输出：

1. 失败的门禁名称。
2. 实际运行命令。
3. 可复现的本地修复命令。
4. 不包含秘密和未脱敏本地路径的必要上下文。

## 10. Branch Protection 契约

实现并验证后，建议将以下 Job 名称设为 `main` 的 required status checks：

- `quality`
- `generated-and-security`
- `windows-build`

Job ID 和显示名称一旦进入 Branch Protection，不得随意重命名。需要重命名时，必须先同步保护规则，避免所有 PR 被永久阻塞。

仓库当前未配置正式 CI，因此 M0-003 的实现 PR 在 Workflow 合并前无法依赖 `main` 上已存在的 required check。实施计划必须包含自举顺序：

1. 在实现 PR 中让 Workflow 运行并变绿。
2. 完成 RED→GREEN 证明。
3. 合并 CI PR。
4. 再配置 Branch Protection required checks。
5. 通过后续文档或空变更 PR 验证保护规则实际生效。

## 11. RED→GREEN 门禁证明

至少进行一次可控失败演练，推荐使用不污染历史的测试分支：

### 11.1 RED

提交一个容易撤销、不会泄露秘密的受控失败，例如：

- 在专用测试文件中加入格式错误，使 `gofmt` 门禁失败；或
- 临时改变生成 Binding，使漂移检查失败；或
- 临时让一个已有 Vitest 断言失败。

要求：

- 失败原因必须唯一且可解释。
- 记录失败 Run ID 和失败 Job。
- 不使用真实 Token、私钥或类似真实秘密的字符串测试 Gitleaks。

### 11.2 GREEN

撤销受控失败，重新运行全部三个 Job：

- 三个 Job 全部成功。
- 记录绿色 Run ID。
- PR 最终 diff 不含故意失败内容。

## 12. 可维护性与本地复现

新增 `docs/development/ci.md`，至少记录：

- 三个 Job 的职责和本地等价命令
- 固定工具链版本
- Wails Ubuntu 24.04 平台依赖和 `webkit2_41` 原因
- 如何重新生成 Binding
- 如何诊断 pnpm `allowBuilds` 失败
- 如何处理 Gitleaks 误报，禁止直接添加宽泛 Allowlist
- 如何验证 Workflow Action SHA 来源
- 如何执行 Windows build 和进程烟测
- required checks 名称和 Branch Protection 配置顺序

复杂检查应落到仓库内脚本并具有单元测试；YAML 主要负责安装环境和编排命令。

## 13. 预期文件结构

```text
.github/
└── workflows/
    └── ci.yml

docs/
└── development/
    └── ci.md

scripts/
└── ci/
    ├── check-dependencies.*
    ├── check-generated.*
    ├── check-security-boundary.*
    └── smoke-windows.*
```

脚本具体语言在实现计划中确定。选择标准：

- 能在本地和 GitHub Runner 复现
- 可测试
- 不依赖不稳定 Shell 行为
- 不复制同一规则的 Bash/PowerShell 两份实现

优先使用 Go 编写跨平台结构化检查；Windows 进程启动/清理由小型 PowerShell 脚本完成。

## 14. 测试设计

### 14.1 Workflow 静态测试

- YAML 可解析
- 仅允许指定 Trigger
- 顶层权限为 `contents: read`
- 无 `pull_request_target`
- Action 全部固定完整 SHA
- 三个 Job ID 与名称稳定
- 超时和并发取消存在

### 14.2 CI 脚本单元测试

- 依赖范围识别正常与非法样例
- `allowBuilds` 只允许批准包
- 远程资源扫描能识别 script/font/CDN
- 范围扫描对当前基线无误报，并能识别受控违规 Fixture
- 生成路径 Allowlist 不覆盖无关文件

### 14.3 集成测试

- 干净 Ubuntu Runner 运行 `quality`
- 干净 Ubuntu Runner 运行 `generated-and-security`
- 干净 Windows Runner 运行 `windows-build`
- 缓存未命中时全部成功
- RED→GREEN 演练证明失败能够阻断

## 15. 验收标准

M0-003 只有同时满足以下条件才能标记 `DONE`：

- [ ] `.github/workflows/ci.yml` 已合并。
- [ ] `quality`、`generated-and-security`、`windows-build` 在干净 Runner 上成功。
- [ ] Go 1.26.5、Node 24.18.0、pnpm 11.12.0、Wails v2.13.0 均明确固定。
- [ ] Action 引用全部固定到完整 Commit SHA，并带版本注释。
- [ ] Workflow 权限为只读，不使用 Secret，不使用 `pull_request_target`。
- [ ] `go test ./...` 和关键竞态检测通过。
- [ ] Frontend 全部测试、typecheck、lint、build 通过。
- [ ] Binding、Go Module 和 pnpm 锁文件无漂移。
- [ ] Secret、依赖政策、远程资源和范围边界扫描通过。
- [ ] Windows production build 和 EXE 启动烟测通过。
- [ ] 至少一次受控 RED 与恢复 GREEN Run 有证据。
- [ ] `docs/development/ci.md` 能让开发者本地复现门禁。
- [ ] PR 最终没有故意失败 Fixture、临时 Workflow 或范围外业务代码。
- [ ] 实现 PR 合并，并通过独立状态同步将 `M0-003` 更新为 `DONE`。

## 16. 风险与缓解

| 风险 | 后果 | 缓解 |
|---|---|---|
| Ubuntu Wails 原生依赖变化 | `go test ./...` 无法编译根包 | 固定 Ubuntu 24.04、记录系统包、使用 `webkit2_41`、失败时显式调整设计 |
| Windows Wails 冷缓存过慢 | PR 反馈延迟 | 独立 Job、合理缓存、25 分钟超时、并发取消 |
| 缓存污染或陈旧 | 假成功或难复现 | 缓存只加速，Key 包含锁文件与工具版本，未命中必须成功 |
| Action Tag 被移动 | 供应链风险 | 完整 Commit SHA Pinning |
| Secret Scan 误报 | 开发者绕过门禁 | 精确规则、局部注释、禁止宽泛全局 Allowlist |
| 路径过滤导致 required check Pending | PR 永久阻塞 | 不使用 paths/paths-ignore |
| Job 名称被修改 | Branch Protection 失效或阻塞 | 固定三个 Job 名，变更需同步保护规则 |
| 受控失败残留 | `main` 被污染 | 最终 diff、状态扫描与 PR Review 强制确认 |
| 范围扫描阻止未来合法功能 | 开发停滞 | 每个后续模块设计 PR 显式更新 Allowlist，不允许临时关闭 |

## 17. 资料依据

设计依据以下官方资料：

- GitHub Actions Workflow Syntax：Trigger、Permissions、Concurrency 与 Job 定义。
- GitHub Actions Secure Use：最小 Token 权限、避免 `pull_request_target`、完整 Commit SHA Pinning。
- GitHub Runner Images：`ubuntu-24.04` 与 `windows-2025` 可用标签及 `latest` 迁移风险。
- Wails v2.13.0 Installation：Ubuntu 24.04 的 GTK/WebKit 依赖与 `webkit2_41` Build Tag。
- actions/checkout、actions/setup-go、actions/setup-node、actions/cache 官方 Release。
- Gitleaks 官方 Release 与 CLI 文档。

## 18. 后续步骤

设计批准后的唯一下一步是编写 M0-003 详细 Implementation Plan。实现前不得直接创建 `.github/workflows/ci.yml`。
