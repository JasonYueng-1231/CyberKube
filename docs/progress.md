# 开发进度表（持续更新）

更新时间: 2025-11-13

| 阶段 | 事项 | 负责人 | 状态 | 说明 |
|---|---|---|---|---|
| 0 | 仓库初始化与项目骨架 | AI | 已完成 | 已创建 backend/frontend/infra/docs/CI 与 README |
| 1 | 后端基础框架与健康检查 | AI | 已完成 | Gin 路由、中间件、统一响应、/api/v1/health |
| 2 | 前端脚手架与仪表盘首页 | AI | 已完成 | Vite+AntD，按 demo.png 落地首页骨架 |
| 3 | Compose 与 Nginx 配置 | AI | 已完成 | MySQL 8.4、Redis 7.2、Nginx 反代 /api |
| 4 | CI 工作流配置 | AI | 已完成 | 后端/前端构建校验与镜像构建校验 |
| 5 | 推送 GitHub 首版 | AI | 已完成 | 初始化 main 分支并推送 |
| 6 | 认证与集群管理 MVP | AI | 进行中 | JWT 登录、默认 admin、集群创建/列表/删除，前端登录与集群页 |

后续里程碑（MVP）
- 用户与登录（本地用户 + JWT）
- 多集群接入（kubeconfig 加密存储）
- Deployment/Pod/Service/ConfigMap/Secret/Namespace 基础管理
- Pod 日志与 WebShell
- 简易审计日志
