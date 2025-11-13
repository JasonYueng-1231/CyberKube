# 回滚指南

本文档说明如何将代码/部署快速回滚到某个稳定版本(tag)。

## 1. 基于 tag 的代码回滚（临时修复）

最安全的方式是“回滚提交(revert)”，而不是强制改写历史：

1) 找到目标版本 tag（例如 `v0.1.0`）
```
git fetch --tags
git checkout v0.1.0
```

2) 创建 hotfix 分支并将差异提交为反向改动
```
git checkout -b hotfix/revert-to-v0.1.0
# 将当前主分支上的问题提交逐个 git revert
# 或者直接以 v0.1.0 为基线 cherry-pick 需要的修复提交
git push -u origin hotfix/revert-to-v0.1.0
# 发起 PR → CI 通过 → 合并到 main
```

说明：不建议对 `main` 做 `reset --hard` + `--force` 的强推，这会改写历史，影响协作。

## 2. 直接将 main 指回某 tag（紧急止血）

仅在“无人协作、紧急止血”的场景可用（改写历史）：
```
git checkout main
git reset --hard v0.1.0
git push -f origin main
```

随后应尽快同步团队成员并基于最新 `main` 继续工作。

## 3. 部署层回滚

若使用镜像/包发布，请统一以 tag 为版本号：
- 例：`ghcr.io/org/cyberkube/backend:v0.1.0`、`frontend:v0.1.0`
- 回滚只需将部署清单中的镜像 tag 切回目标版本，或 `helm rollback` 到对应 revision

## 4. 数据库与可逆迁移

- 如涉及 DB 变更，按“升级 + 降级”双向迁移脚本维护
- 发布脚本在回滚时先执行降级，再回滚代码/镜像

