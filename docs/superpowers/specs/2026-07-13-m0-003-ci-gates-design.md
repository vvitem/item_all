# M0-003 分层 CI 门禁设计

> 状态：待评审  
> 负责人：vvitem  
> 最后更新：2026-07-13  
> 基线 Commit：`351ca3b3466ff347948dec444716720256e8c2c7`  
> 关联 Milestone：`M0-foundation`  
> 关联 Issue/PR：[Issue #7](https://github.com/vvitem/item_all/issues/7)

## 1. 目标

为 ItemAll 建立正式、可复现、默认安全的 Pull Request CI 门禁，使任何进入 `main` 的改动都经过：

1. Go 格式、模块完整性、静态检查、单元测试和关键竞态检测。
2. Frontend 冻结依赖安装、Vitest、TypeScript、ESLint 和 production build。
3. Wails Binding、Go Module 与 pnpm 锁文件漂移检查。
4. 依赖版本、安装脚本授权和 Workflow 供应链检查。
5. Secret、远程资源和 M0-003 范围外能力扫描。
6. Windows Wails production build 与 `ItemAll.exe` 启动烟测。
7. 一次受控 RED→GREEN 演练，证明失败确实阻断合并。

CI 是进入 `main` 的质量与安全底线，不是发布系统。本任务不创建安装包、签名、Release、部署或遥测流程。

## 2. 已确认方案

采用 **方案 B：单 Workflow、三个分层 Job**。

```text
pull_request / push(main) / workflow_dispatch
                    │
                    ├── quality (ubuntu-24.04)
                    ├── generated-and-security (ubuntu-24.04)
                    └── windows-build (windows-2025)
```

三个 Job 都是独立必需门禁：

- `quality`：提供快速、可读的 Go/Frontend 质量反馈。
- `generated-and-security`：验证生成文件、依赖政策、Secret 和安全边界。
- `windows-build`：提供真实目标平台的桌面构建与启动事实。

不使用单一 Windows Job，也不在 M0-003 提前抽取 reusable workflow 或 composite action。

## 3. 范围

### 3.1 必须交付

- `.github/workflows/ci.yml`
- `docs/development/ci.md`
- 必要且可测试的本地 CI 检查脚本
- README 中的 CI 入口
- 一次故意失败和恢复绿色的 Actions 运行证据
- 合并后的项目状态同步

### 3.2 明确不做

- SSH、SFTP、数据库、SQLite 业务表、Operation Bus、AI、凭据或遥测
- Release、安装包、签名、校验和或 GitHub Release
- macOS/Linux 桌面 production build
- 自托管 Runner
- 长期上传 EXE Artifact
- CodeQL、SBOM、Artifact Attestation、OpenSSF Scorecard 或 Dependabot
- 覆盖率百分比门槛；M0 以风险路径测试为准

## 4. 信任与供应链边界

### 4.1 最小权限

Workflow 顶层固定：

```yaml
permissions:
  contents: read
```

其余未声明权限均为 `none`。CI 不写仓库、不评论 PR、不上传安全事件，也不访问环境 Secret。

### 4.2 不信任 PR 内容

- 使用 `pull_request`，禁止 `pull_request_target`。
- 不读取仓库 Secret。
- 不把 PR 标题、分支名、Issue 文本等不可信 Context 直接插入 `run:`。
- `actions/checkout` 设置 `persist-credentials: false`。
- 不使用自托管 Runner。

### 4.3 Action 必须不可变

所有 `uses:` 固定到来源仓库的 **完整 40 字符 Commit SHA**，并在行尾标注版本：

```yaml
uses: actions/checkout@<40-char-sha> # v7.0.0
```

禁止 `@main`、`@master` 或只使用可移动大版本 Tag。当前核验的稳定线：

- `actions/checkout` v7.0.0
- `actions/setup-go` v6.5.0
- `actions/setup-node` v7.0.0
- `actions/cache` v6.1.0，仅在内建缓存不足时使用

完整 SHA 在实现计划中从官方 Release/Tag 冻结，不从 Fork 获取。

### 4.4 Runner 明确固定

使用：

- `ubuntu-24.04`
- `windows-2025`

不使用 `*-latest`，避免 GitHub 迁移标签时无审查改变操作系统。Runner 镜像仍会按周更新，因此 Go、Node、pnpm、Wails 和扫描工具版本必须显式固定。

## 5. Trigger 与并发

```yaml
on:
  pull_request:
    branches:
      - main
  push:
    branches:
      - main
  workflow_dispatch:

concurrency:
  group: ${{ github.workflow }}-${{ github.event.pull_request.number || github.ref }}
  cancel-in-progress: true
```

设计规则：

- Draft PR 也运行 CI，尽早发现问题。
- 不配置 `paths` 或 `paths-ignore`，避免 required check 因 Workflow 被跳过而长期 Pending。
- 同一 PR 新提交取消旧运行。
- 并发组包含 Workflow 名，避免未来不同 Workflow 互相取消。

## 6. Job：`quality`

### 6.1 环境

- Runner：`ubuntu-24.04`
- 超时：15 分钟
- Checkout：`fetch-depth: 1`、`persist-credentials: false`
- 权限：`contents: read`

### 6.2 固定工具链

- Go：`go.mod` 中的 `toolchain go1.26.5`
- Node：`.nvmrc` 中的 `24.18.0`
- pnpm：`frontend/package.json` 中的 `pnpm@11.12.0`
- Corepack 启用后必须验证实际 pnpm 版本

Runner 预装版本不得替代仓库声明。

### 6.3 Ubuntu Wails 依赖

Ubuntu 24.04 编译 Wails 根包需要 GCC、GTK3 和 WebKitGTK 4.1。Job 必须：

- 安装明确的系统包集合。
- 对 Go 编译、vet 和 test 使用 `webkit2_41` Build Tag。
- 在 `docs/development/ci.md` 说明该 Tag 的来源与本地复现方式。

不得因为平台依赖复杂而静默删除 `go test ./...`。

### 6.4 Go 检查顺序

1. `go version`
2. `go mod download`
3. `go mod verify`
4. `gofmt -l`，输出非空即失败
5. `go vet ./...`
6. `go test ./...`
7. `go test -race ./internal/buildinfo`

### 6.5 Frontend 检查顺序

1. 验证 Node/pnpm 精确版本
2. `pnpm install --frozen-lockfile`
3. `pnpm test:run`
4. `pnpm typecheck`
5. `pnpm lint`
6. `pnpm build`

任何测试失败、TypeScript 错误、ESLint warning 或构建错误都使 Job 失败。

### 6.6 缓存

优先使用 setup action 内建缓存：

- Go 缓存由 `go.sum` 驱动。
- pnpm 缓存由 `frontend/pnpm-lock.yaml` 驱动。

缓存只用于加速：未命中必须成功，不缓存 `frontend/dist`、测试结果或生成 Binding。

## 7. Job：`generated-and-security`

### 7.1 环境与完整历史

- Runner：`ubuntu-24.04`
- 超时：15 分钟
- Checkout：`fetch-depth: 0`、`persist-credentials: false`
- 权限：`contents: read`

`fetch-depth: 0` 是完整 Git 历史 Secret Scan 的硬要求，不能退回默认浅克隆。

该 Job 同样安装 Ubuntu Wails GTK/WebKit 依赖，并使用 `webkit2_41` Build Tag，确保 Binding 生成和 Go Module 校验能够稳定分析 Wails 根包。

### 7.2 生成文件漂移

安装固定 Wails CLI：

```text
github.com/wailsapp/wails/v2/cmd/wails@v2.13.0
```

重新生成 Binding，并检查：

- `frontend/wailsjs/**`
- `go.mod`
- `go.sum`
- `frontend/pnpm-lock.yaml`

使用 `git diff --exit-code -- <allowlisted paths>`。CI 不自动提交结果；发现漂移时输出本地修复命令并失败。

### 7.3 依赖政策

验证：

- `packageManager` 精确为 `pnpm@11.12.0`
- `engines.node` 与 `.nvmrc` 一致
- npm 依赖不得使用 `latest`、`*`、`^`、`~`、`workspace:*` 或无界 Git URL
- `pnpm install --frozen-lockfile` 不修改锁文件
- `pnpm-workspace.yaml.allowBuilds` 只允许经过评审的包；M0-003 基线仅允许 `esbuild`
- `go.mod` Module Path、Go 版本、toolchain 和 Wails 版本符合 DEC-004
- `go mod tidy` 后无漂移

复杂规则写成仓库内可测试脚本，不在 YAML 中堆叠难维护正则。

### 7.4 Secret Scan

推荐固定版本 Gitleaks CLI：

1. 从官方 Release 下载对应 Runner 的二进制。
2. 校验官方 SHA256。
3. 运行：

```text
gitleaks git --redact --no-banner
```

要求：

- 扫描完整可访问 Git 历史。
- 输出必须脱敏。
- Workflow 不加入真实测试 Secret。
- RED 演练不使用真实 Token、私钥或高度仿真的秘密字符串。

### 7.5 静态安全边界

M0-003 基线继续禁止：

- `.env`、私钥或证书私钥文件
- 前端远程 script、stylesheet、font 或 CDN URL
- SSH、数据库、AI、云 SDK、遥测或环境变量读取能力
- `GetAppInfo()` 之外的 Wails API

扫描采用两层：

1. 明确文件、依赖和 Binding Allowlist。
2. 高风险关键词启发式扫描，命中后失败并要求人工确认。

后续 M0-004/M0-006/M0-007 等模块必须在各自设计 PR 中显式更新规则，禁止为通过 CI 全局关闭扫描。

### 7.6 Workflow 自身检查

检查 `.github/workflows/*.yml`：

- 顶层 `permissions` 存在且没有写权限
- 不存在 `pull_request_target`
- `uses:` 全部为完整 SHA，并带版本注释
- 不可信 GitHub Context 未直接进入 `run:`
- 不存在长生命周期 Secret 或云凭据
- 三个 Job ID、Runner、超时和并发配置符合本设计

## 8. Job：`windows-build`

### 8.1 环境

- Runner：`windows-2025`
- 超时：25 分钟
- Checkout：`fetch-depth: 1`、`persist-credentials: false`
- 权限：`contents: read`

### 8.2 工具链

- Go 1.26.5
- Node 24.18.0
- pnpm 11.12.0
- Wails CLI v2.13.0

可缓存 Go Module、Go Build Cache 和 Wails CLI，但 Key 必须包含 Runner OS、Go toolchain、Wails 版本和 `go.sum` Hash。

### 8.3 验证顺序

1. 安装固定 Go、Node、pnpm 和 Wails
2. `pnpm install --frozen-lockfile`
3. `go test ./...`
4. `wails build -clean -trimpath`
5. 验证 `build/bin/ItemAll.exe` 存在且非空
6. 启动 EXE，确认观察窗口内未立即异常退出
7. 清理完整进程树

### 8.4 烟测语义

烟测只证明：

- Windows 能加载 EXE。
- Wails/WebView2 没有立即崩溃。
- 应用在观察窗口内保持运行。

它不证明视觉布局或交互正确；截图级 UI 验收属于后续 Windows E2E。

### 8.5 产物策略

M0-003 默认不上传 EXE Artifact。正式分发、保留、校验和和签名属于 M4-007。故障诊断可临时上传脱敏日志，但不上传包含本地路径或敏感信息的整包。

## 9. Failure Semantics

任何 Job 失败都阻止合并。禁止：

- `continue-on-error: true`
- `|| true`
- 忽略退出码
- 把失败包装为 warning
- 只在 `main` push 后验证
- 仅上传报告但保持 Job 绿色

失败步骤必须输出：门禁名称、实际命令、本地复现/修复命令，以及经过脱敏的必要上下文。

## 10. Branch Protection 契约

实现并验证后，将以下稳定名称设为 `main` required status checks：

- `quality`
- `generated-and-security`
- `windows-build`

Job 名称进入保护规则后不得随意改变。需要重命名时必须先同步 Branch Protection。

自举顺序：

1. 在实现 PR 中让新 Workflow 运行并变绿。
2. 完成 RED→GREEN 证明。
3. 合并 CI PR。
4. 配置三个 required checks。
5. 用后续文档或空变更 PR 验证保护规则真实生效。

## 11. RED→GREEN 门禁证明

### 11.1 RED

在专用测试分支提交一个唯一、易撤销且无秘密的受控失败，优先选择：

- Go 格式错误；或
- 临时 Binding 漂移；或
- 临时 Vitest 断言失败。

记录失败 Run ID、失败 Job 和具体步骤。

### 11.2 GREEN

撤销受控失败并重新运行全部三个 Job：

- 三个 Job 全部成功。
- 记录绿色 Run ID。
- 最终 PR diff 不包含故意失败内容。

## 12. 本地复现与文件结构

新增：

```text
.github/workflows/ci.yml
docs/development/ci.md
scripts/ci/
```

`docs/development/ci.md` 至少记录：

- 三个 Job 的职责和本地等价命令
- 固定工具链版本
- Ubuntu 24.04 Wails 依赖和 `webkit2_41`
- Binding 重新生成方式
- pnpm `allowBuilds` 故障处理
- Gitleaks 误报处理，禁止宽泛全局 Allowlist
- Action SHA 来源核验方式
- Windows build 和进程烟测
- required checks 与 Branch Protection 自举顺序

复杂检查优先使用 Go 编写跨平台结构化脚本并配单元测试；Windows 进程启动/清理由小型 PowerShell 脚本完成。YAML 只负责环境安装和编排。

## 13. 测试设计

### 13.1 Workflow 静态测试

- YAML 可解析
- Trigger 仅为本设计允许项
- 权限为 `contents: read`
- 无 `pull_request_target`
- Action 全部固定完整 SHA
- 三个 Job ID、Runner、超时和并发稳定
- `generated-and-security` 明确 `fetch-depth: 0`

### 13.2 CI 脚本单元测试

- 依赖范围正常与非法 Fixture
- `allowBuilds` 仅允许批准包
- 远程资源扫描识别 script/font/CDN
- 范围扫描识别受控违规且对当前基线无误报
- 生成路径 Allowlist 不覆盖无关文件
- Workflow 安全规则识别可移动 Action Tag 和写权限

### 13.3 集成验证

- 干净 Ubuntu Runner 运行 `quality`
- 干净 Ubuntu Runner 运行 `generated-and-security`
- 干净 Windows Runner 运行 `windows-build`
- 缓存未命中时全部成功
- RED→GREEN 演练证明失败能够阻断

## 14. 验收标准

M0-003 只有同时满足以下条件才能标记 `DONE`：

- [ ] `.github/workflows/ci.yml` 已合并。
- [ ] 三个 Job 在干净 Runner 上成功。
- [ ] Go 1.26.5、Node 24.18.0、pnpm 11.12.0、Wails v2.13.0 明确固定。
- [ ] Action 全部固定完整 Commit SHA，并带版本注释。
- [ ] Workflow 只读、无 Secret、无 `pull_request_target`。
- [ ] `go test ./...`、关键竞态和 Frontend 全量检查通过。
- [ ] Binding、Go Module 和 pnpm 锁文件无漂移。
- [ ] Secret、依赖政策、远程资源和范围扫描通过。
- [ ] Windows production build 与 EXE 启动烟测通过。
- [ ] 至少一次受控 RED 与恢复 GREEN 有 Run 证据。
- [ ] 本地复现文档完整。
- [ ] 最终 PR 无失败 Fixture、临时 Workflow 或范围外业务代码。
- [ ] 实现 PR 合并，并通过独立状态同步将 M0-003 更新为 `DONE`。

## 15. 风险与缓解

| 风险 | 后果 | 缓解 |
|---|---|---|
| Ubuntu Wails 原生依赖变化 | 根包无法编译 | 固定 Ubuntu 24.04、明确系统包、使用 `webkit2_41` |
| Windows Wails 冷缓存慢 | PR 反馈延迟 | 独立 Job、合理缓存、25 分钟超时、并发取消 |
| 缓存污染 | 假成功或难复现 | 缓存只加速，Key 包含锁文件和工具版本 |
| Action Tag 被移动 | 供应链风险 | 完整 Commit SHA Pinning |
| Secret Scan 误报 | 开发者绕过门禁 | 精确规则、局部说明、禁止宽泛 Allowlist |
| 路径过滤跳过 required check | PR 长期 Pending | 不使用路径过滤 |
| Job 名称变化 | Branch Protection 失效 | 固定三个 Job 名，变更前同步保护规则 |
| 受控失败残留 | `main` 被污染 | 最终 diff 和评审强制确认 |
| 范围扫描阻止后续合法模块 | 开发停滞 | 后续设计 PR 显式更新 Allowlist，不关闭扫描 |

## 16. 资料依据

设计依据以下官方资料：

- GitHub Actions Workflow Syntax：Trigger、Permissions、Concurrency、Job 与路径过滤行为。
- GitHub Actions Secure Use：最小 Token 权限、避免 `pull_request_target`、完整 Commit SHA Pinning。
- GitHub Runner Images：`ubuntu-24.04`、`windows-2025` 标签及 `latest` 迁移风险。
- Wails v2.13.0 Installation：Ubuntu 24.04 GTK/WebKit 依赖与 `webkit2_41`。
- actions/checkout、actions/setup-go、actions/setup-node、actions/cache 官方 Release。
- Gitleaks 官方 Release 与 CLI 文档。

## 17. 后续步骤

设计批准后的唯一下一步是编写 M0-003 Implementation Plan。实现前不得创建正式 `.github/workflows/ci.yml`。
