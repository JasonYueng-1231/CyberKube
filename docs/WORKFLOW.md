# 开发/发布工作流

本仓库采用“主干 + 特性分支”的轻量工作流，并用语义化版本(tag)管理回滚点。

## 分支策略
- `main`：仅存放稳定可发布版本。禁止直接推送，必须通过 PR 合并（建议在仓库设置中开启 branch protection）。
- `develop`：日常集成分支（可选）。较大特性在此集成验证，再发起 PR 到 `main`。
- `feature/*`：功能开发分支，例如 `feature/workloads-ws`。
- `hotfix/*`：线上问题修复分支，例如 `hotfix/fix-login-401`。

## 提交流程
1. 从 `develop`（或 `main`）创建 `feature/*` 分支
2. 开发 + 自测 → 推送到远端
3. 发起 PR → CI 通过 → Code Review → 合并入 `develop`（或直接入 `main`）
4. 准备发布时，从 `main` 打 tag（语义化版本）

## 版本与标签
- 语义化版本：`vMAJOR.MINOR.PATCH`（如 `v0.1.0`）
- 每次可用发布必须打 tag；`main` 上的每个 tag 都是可回滚锚点

## CI 建议
- PR 必须通过：后端 build/test、前端 build、镜像构建校验
- main 受保护：仅允许通过 PR 合并，禁止直接 push；必须勾选“Require status checks to pass before merging”

## 回滚策略（快速）
- “代码回滚”：将 `main` 回退到某个已发布 tag（详见 `docs/rollback.md`）
- “发布回滚”：部署层按版本号回滚（镜像/包/Chart），与代码 tag 对齐

