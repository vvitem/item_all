> 状态：已确认  
> 负责人：vvitem  
> 最后更新：2026-07-11  
> 基线 Commit：`b590887fea1e2c43e7831a48932b9be5d44ccfa3`  
> 关联 Milestone：`M1-ssh-diagnosis / M2-database`  
> 关联 Issue/PR：待创建

# 结构化只读工具契约

## 通用 ToolDefinition

```go
type ToolDefinition struct {
    Name              string
    AssetTypes        []string
    InputSchema       json.RawMessage
    OutputSchema      json.RawMessage
    DefaultTimeout    time.Duration
    MaxOutputBytes    int64
    RiskLevel         RiskLevel
    Cancellable       bool
    EvidenceExtractor string
    ErrorCodes        []string
}
```

所有工具必须经过 Operation Bus；模型只填写 Schema 允许的参数。

## SSH 工具

| 工具 | 关键输入限制 | 默认超时 | 最大输出 | 风险 | Evidence |
|---|---|---:|---:|---|---|
| `system.info` | 无任意命令参数 | 10s | 64KB | R0 | system snapshot |
| `system.load` | sample_seconds 1–10 | 15s | 64KB | R0 | load sample |
| `system.memory` | 固定采集模板 | 10s | 64KB | R0 | memory snapshot |
| `system.time` | 固定模板 | 10s | 64KB | R0 | clock/NTP |
| `disk.usage` | mount filter≤20 | 15s | 256KB | R0 | filesystem usage |
| `file.top_size` | allow-root、depth≤4、limit≤100 | 60s | 1MB | R1 | ranked paths |
| `file.list_recent` | allow-root、window≤24h | 30s | 512KB | R1 | recent file list |
| `service.status` | unit name pattern | 15s | 256KB | R0 | service state |
| `process.list` | limit≤100 | 15s | 256KB | R1 | process snapshot |
| `network.listen` | ports≤20 | 15s | 256KB | R0 | listener list |
| `log.tail` | path allowlist、lines≤1000 | 30s | 1MB | R1 | log excerpt |
| `log.search` | window≤24h、results≤500 | 60s | 2MB | R1 | matched records |

底层实现可使用内置、版本化的命令模板；模板 hash 写入 Audit。模型不能提供命令片段。

## 数据库工具

| 工具 | 关键输入限制 | 默认超时 | 最大输出 | 风险 | Evidence |
|---|---|---:|---:|---|---|
| `db.schemas` | limit≤500 | 10s | 512KB | R0 | schema list |
| `db.tables` | schema required、limit≤1000 | 15s | 1MB | R0 | table list |
| `db.columns` | schema/table required | 10s | 512KB | R0 | column metadata |
| `db.indexes` | schema/table required | 15s | 512KB | R0 | index metadata |
| `db.sessions` | limit≤200、SQL脱敏 | 15s | 1MB | R1 | session snapshot |
| `db.locks` | limit≤200 | 15s | 1MB | R1 | blocking graph |
| `db.connection_stats` | 固定方言查询 | 10s | 256KB | R0 | connection stats |
| `db.table_stats` | schema required、limit≤500 | 30s | 2MB | R1 | table statistics |
| `db.query_readonly` | 单条SELECT/WITH、≤1000行/10MB/30s | 30s | 10MB | R2 | result+SQL hash |
| `db.explain` | 单条SELECT/WITH、默认禁ANALYZE | 30s | 2MB | R1 | normalized plan |

## 明确禁止

```text
ssh.exec(command)
shell.run(script)
db.execute(sql)
file.read(any_path)
```

禁止多语句、DML、DDL、存储过程调用、`COPY ... PROGRAM` 等方言逃逸。`db.query_readonly` 必须同时通过客户端 Parser、只读事务/会话和数据库权限三层约束。
