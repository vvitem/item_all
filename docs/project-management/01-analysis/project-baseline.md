> 状态：已确认  
> 负责人：vvitem  
> 最后更新：2026-07-11  
> 基线 Commit：`b590887fea1e2c43e7831a48932b9be5d44ccfa3`  
> 关联 Milestone：`M0-foundation`  
> 关联 Issue/PR：待创建

# 项目基线

## 基线信息

| 项目 | 当前值 |
|---|---|
| 仓库 | `https://github.com/vvitem/item_all` |
| 默认分支 | `main` |
| 基线 Commit | `b590887fea1e2c43e7831a48932b9be5d44ccfa3` |
| Commit 信息 | `Initial commit` |
| 分析日期 | `2026-07-11` |
| 仓库可见性 | Private |
| 当前负责人 | `vvitem` |

## 当前目录结构

```text
item_all/
└── README.md
```

GitHub 元数据报告仓库 `size=0`；基线 Commit 只新增两行 README。代码搜索未发现 `go.mod`、Wails 配置或前端工程标识。

## 当前技术栈

| 技术 | 规划状态 | 代码状态 |
|---|---|---|
| Go | 已确认目标 | 未初始化 |
| Wails v2 | 已确认目标 | 未初始化 |
| React + TypeScript | 已确认目标 | 未初始化 |
| SQLite WAL | 已确认目标 | 未初始化 |
| OS Keychain | 已确认目标 | 未初始化 |
| GitHub Actions | 目标交付能力 | 未配置 |

## 当前可运行方式

无。仓库没有可执行程序、构建脚本或依赖清单。当前只能阅读 README 和项目文档。

## 构建与测试状态

- Build：不可执行。
- Unit Test：不存在。
- Integration Test：不存在。
- E2E：不存在。
- CI：不存在。
- Release：不存在。

## 模块状态

### 已经真实实现

- 根目录项目名称和一句话定位。
- 本项目管理文档体系（当前分支/PR 中）。

### 已部分实现

无。

### 已设计但未实现

- Go + Wails + React 桌面架构。
- Operation Bus、Task Engine、Audit/Evidence、结构化只读工具。
- SSH/SFTP、MySQL/PostgreSQL 适配器。
- Windows Beta 交付链路。

### 尚未开始

全部业务模块、数据模型、Migration、测试、CI 和发布。

### 无法确认

- 最终产品名称是否继续使用 `item_all`。
- Windows 代码签名预算与证书获取时间。
- 首发 AI Provider 最小兼容集合。

## 当前主要技术债

绿地仓库不存在历史代码债，但存在“工程债务前置风险”：若先写 UI/协议代码再建立 Operation Bus、审计和测试边界，后续会形成高成本返工。

## 当前主要风险

1. 范围扩张导致六个月无法形成闭环。
2. AI/远端内容造成 Prompt Injection 或读取越界。
3. 凭据进入 Prompt、日志或 Evidence。
4. Operation Bus 被 UI 或适配器绕过。
5. Windows PTY/WebView/Keychain 兼容性晚暴露。
6. 单人依赖 Codex 生成大量不可审查代码。
