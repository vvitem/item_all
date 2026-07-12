> 状态：已确认  
> 负责人：vvitem  
> 最后更新：2026-07-11  
> 基线 Commit：`b590887fea1e2c43e7831a48932b9be5d44ccfa3`  
> 关联 Milestone：`M0-foundation`  
> 关联 Issue/PR：待创建

# 错误处理

## 目标

错误必须可分类、可追踪、可向用户解释，同时不泄漏凭据、DSN、原始秘密或内部堆栈。

## Error Envelope

```go
type AppError struct {
    Code       string
    MessageKey string
    Retryable  bool
    TraceID    string
    Details    map[string]any // allowlist only
    Cause      error          // internal only
}
```

## 分层职责

- Adapter：保留协议 cause，返回细粒度内部错误。
- Connection：区分 DNS、超时、TLS、认证、指纹。
- Operation Bus：映射稳定错误码，写失败 Audit，清理 Details。
- Task Engine：决定重试、失败或取消，不通过字符串匹配错误。
- UI：按 `MessageKey` 本地化，展示 TraceID 和可行动建议。

## 禁止模式

- `fmt.Errorf("connect %s:%s with %s", host, port, password)`。
- 捕获所有错误后返回“操作失败”。
- 在失败时静默切换到不安全路径。
- 将未知执行状态标记成功。

## 重试

仅网络瞬态、只读且幂等的错误可重试；认证、Policy、Plan Scope、Parser、Redaction 错误不可自动重试。
