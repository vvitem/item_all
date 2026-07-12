> 状态：草案  
> 负责人：vvitem  
> 最后更新：2026-07-11  
> 基线 Commit：`b590887fea1e2c43e7831a48932b9be5d44ccfa3`  
> 关联 Milestone：`M0-foundation → M4-beta`  
> 关联 Issue/PR：待创建

# 可观测性

## 本地优先原则

默认不上传命令、SQL、日志、结果、资产地址或凭据。可选遥测必须默认关闭、字段白名单、可预览和可撤回。

## 日志

结构化字段：`level`、`timestamp`、`component`、`trace_id`、`task_id`、`operation_id`、`error_code`、`duration_ms`。禁止记录 Secret、完整 DSN、私钥路径内容、原始 SQL 结果。

## 指标

- Operation latency/error/cancel by capability。
- Task terminal state、step retry、recovery count。
- Evidence bytes/truncation/redaction failure。
- SSH/DB connection success and pool health。
- Scenario success rate and time-to-first-diagnosis。

## Trace

统一 `TraceID` 贯穿 Wails 请求、Task、Operation、Adapter、Audit 和 Evidence。首版可采用内部 trace model；OpenTelemetry 接入不应阻塞 Beta。

## 诊断包

用户主动导出，包含版本、平台、配置摘要、脱敏日志、Migration 版本和失败 Trace。生成前执行 Secret Scan，失败则阻止导出。
