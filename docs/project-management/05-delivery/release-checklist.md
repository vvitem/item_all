> 状态：维护中  
> 负责人：vvitem  
> 最后更新：2026-07-11  
> 基线 Commit：`b590887fea1e2c43e7831a48932b9be5d44ccfa3`  
> 关联 Milestone：`M4-beta`  
> 关联 Issue/PR：待创建

# 发布检查表

## 代码与文档

- [ ] 目标 tag 对应干净 Commit。
- [ ] 全部 P0 为 DONE，文档与状态一致。
- [ ] CHANGELOG、README、安装和已知问题更新。

## 构建

- [ ] 锁定 Go/Node/pnpm/Wails 版本。
- [ ] 干净 Windows runner 构建成功。
- [ ] 产物包含版本和 Commit ID。
- [ ] 生成 SHA256SUMS。

## 测试

- [ ] Go/前端/集成/Windows E2E 通过。
- [ ] MySQL/PG 和 SSH 真机 smoke 通过。
- [ ] 安装、升级、卸载和数据保留验证。

## 发布

- [ ] 标记 prerelease。
- [ ] Release Notes 说明范围和明确不支持项。
- [ ] 下载后校验和匹配。
- [ ] 建立回滚/撤回发布步骤。
