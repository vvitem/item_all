> 状态：已确认  
> 负责人：vvitem  
> 最后更新：2026-07-11  
> 基线 Commit：`b590887fea1e2c43e7831a48932b9be5d44ccfa3`  
> 关联 Milestone：`M0-foundation`  
> 关联 Issue/PR：待创建

# 产品定位

## 一句话定位

面向个人开发者、独立运维和全栈工程师的 **Local-first SafeOps 桌面工作台**：把 Linux 故障和数据库问题转换为可审阅的诊断计划，通过受限制的只读工具采集证据，并输出可复用 Runbook。

## 目标用户

- 维护 5–100 台 Linux 主机的个人开发者或独立运维。
- 维护 1–20 个 MySQL/PostgreSQL 数据库的全栈工程师。
- 需要在本地或内网环境使用 BYOK 模型的用户。
- 重视审计和可解释性，但尚不需要企业级身份平台的小团队成员。

## 核心任务

1. 快速判断服务异常、资源瓶颈、端口和日志问题。
2. 安全查询数据库元数据、会话、锁、执行计划和只读业务数据。
3. 将一次诊断过程保存为参数化 Runbook。
4. 保留完整操作、证据、错误和用户决策轨迹。

## 差异化

- **Plan-first**：AI 先生成最多 8 步的结构化 Plan，用户可删改范围。
- **Evidence-first**：没有 Evidence Contract 满足，不宣告成功。
- **Safe-by-construction**：模型没有 `ssh.exec(command)`、`db.execute(sql)` 等任意执行工具。
- **Local-first**：凭据、连接、审计和结果默认在本机处理。
- **Future-compatible**：数据模型保留 `workspace_id`、`actor_id`、`device_id`，但六个月不建设控制平面。

## 非目标

- 不做通用服务器面板、云管平台或零信任访问网关。
- 不追求数据库客户端全功能替代 DBeaver/DataGrip。
- 不做 RDP、Kubernetes、Redis、Kafka、MongoDB、对象存储。
- 不允许 AI 自动修复生产环境。
- 不以“支持模型数量”或“资产数量”作为北极星指标。

## 北极星指标

- 用户从添加资产到完成第一次可信诊断的时间 ≤ 10 分钟。
- 预定义 20 个场景中 ≥ 70% 无需手工修改工具参数即可完成。
- 100% AI 工具调用进入 Operation Bus 和 Audit。
- AI 写操作数量为 0。
- 任务结束时未知状态数量为 0。
