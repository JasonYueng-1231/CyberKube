# CyberKube —— K8S 集群管理平台（MVP）

CyberKube 是一套面向多集群的 Kubernetes 管理平台，采用 Go + React 技术栈，提供统一的 API 与赛博朋克风格的 Web 界面（参考 docx/demo.png）。本仓库为 monorepo：包含后端、前端、容器编排与文档。

## 目录结构

```
backend/      # Go 后端（Gin + client-go）
frontend/     # React 前端（Vite + Ant Design）
infra/        # 部署与运维（环境变量示例等）
docs/         # 项目文档（含进度表）
.github/      # CI 工作流
```

## 快速开始（本地 Docker Compose）

1) 复制环境变量示例并按需修改

```
cp infra/example.env .env
```

2) 构建并启动

```
docker compose up -d --build
```

3) 访问

- 前端：`http://localhost:80`
- 后端健康检查：`http://localhost:8080/api/v1/health`

默认 `frontend` Nginx 将 `/api` 反向代理到 `backend:8080`。

## 关键环境变量

- `MYSQL_DSN`：后端数据库连接（默认已在 compose 中配置）
- `REDIS_ADDR`：Redis 地址（默认 `redis:6379`）
- `K8S_KUBECONFIG_AES_KEY`：用于加密 kubeconfig 的 32 字节密钥（开发环境可使用示例值）

更多细节见 `docs/` 与各模块 README。

## CI

GitHub Actions 在每次推送与 PR 时执行：
- 后端：`go build`/`go vet`/`go test`（无测试时快速通过）
- 前端：`npm ci && npm run build`
- 镜像：进行构建校验，但不推送远端仓库

## 许可

本仓库仅用于项目开发演示，未附带开源许可证。如需开源请补充许可证文件。

