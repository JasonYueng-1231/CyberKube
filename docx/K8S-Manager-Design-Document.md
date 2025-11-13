# K8S集群管理平台详细设计方案

**项目名称**: K8S-Manager  
**版本**: v1.0  
**创建日期**: 2025-11-10  
**作者**: Jason  
**开发模式**: 个人开发 + AI辅助编程

---

## 目录

1. [项目概述](#1-项目概述)
2. [技术架构](#2-技术架构)
3. [系统设计](#3-系统设计)
4. [数据库设计](#4-数据库设计)
5. [API接口设计](#5-api接口设计)
6. [前端页面设计](#6-前端页面设计)
7. [开发计划](#7-开发计划)
8. [AI辅助开发环境配置](#8-ai辅助开发环境配置)
9. [部署方案](#9-部署方案)
10. [安全设计](#10-安全设计)
11. [附录](#11-附录)

---

## 1. 项目概述

### 1.1 项目背景

开发一个Web化的Kubernetes集群管理平台，提供友好的图形界面来管理多个K8S集群，简化运维操作，提高工作效率。

### 1.2 项目目标

- 支持多K8S集群统一管理
- 提供工作负载（Deployment、Pod等）的完整生命周期管理
- 实现实时日志查看和Web终端功能
- 提供资源监控和可视化
- 完善的权限控制和审计日志

### 1.3 核心功能

#### 阶段一功能（MVP）
- 用户认证与授权
- 多集群管理
- Deployment管理（CRUD、扩缩容、重启）
- Pod管理（查看、删除、日志、WebShell）
- Service管理
- ConfigMap/Secret管理
- 命名空间管理

#### 阶段二功能（扩展）
- StatefulSet、DaemonSet、Job管理
- Ingress管理
- 节点管理
- 资源监控与图表
- RBAC权限管理
- 审计日志

#### 阶段三功能（高级）
- Helm应用市场
- 事件监控与告警
- 资源拓扑图
- 集群备份与恢复

### 1.4 技术选型

| 类型 | 技术 | 版本 | 说明 |
|------|------|------|------|
| 后端语言 | Go | 1.21+ | 高性能、并发支持好 |
| 后端框架 | Gin | 1.9+ | 轻量级、性能优秀 |
| K8S客户端 | client-go | 0.28+ | 官方K8S Go客户端 |
| ORM | GORM | 1.25+ | 功能强大的Go ORM |
| 数据库 | MySQL | 8.0+ | 存储用户、集群信息 |
| 缓存 | Redis | 7.0+ | 会话、配置缓存 |
| 前端框架 | React | 18+ | 主流前端框架 |
| UI组件库 | Ant Design | 5+ | 企业级UI组件 |
| 状态管理 | React Query | 4+ | 数据获取和缓存 |
| 终端组件 | xterm.js | 5+ | Web终端实现 |
| 图表库 | ECharts | 5+ | 数据可视化 |

---

## 2. 技术架构

### 2.1 整体架构

```
┌─────────────────────────────────────────────────────────┐
│                     Browser (用户浏览器)                   │
└──────────────────────┬──────────────────────────────────┘
                       │ HTTPS
                       │
┌──────────────────────▼──────────────────────────────────┐
│                 Nginx (反向代理/静态文件)                   │
└──────────┬──────────────────────────────────┬───────────┘
           │                                  │
    ┌──────▼──────┐                    ┌──────▼──────────┐
    │  前端静态文件  │                    │   后端API服务    │
    │  (React SPA) │                    │   (Go + Gin)    │
    └──────────────┘                    └──────┬──────────┘
                                               │
                    ┌──────────────────────────┼──────────────────┐
                    │                          │                  │
            ┌───────▼────────┐        ┌────────▼────────┐  ┌─────▼─────┐
            │  MySQL数据库    │        │   Redis缓存      │  │  K8S集群   │
            │ (用户/集群/日志) │        │ (Session/Config) │  │ (多个集群)  │
            └────────────────┘        └─────────────────┘  └───────────┘
```

### 2.2 后端架构

```
后端服务
├── API层 (Gin Router + Handler)
│   ├── 认证中间件 (JWT)
│   ├── 权限中间件 (RBAC)
│   ├── 审计中间件 (日志记录)
│   └── CORS中间件
│
├── Service层 (业务逻辑)
│   ├── 用户服务
│   ├── 集群服务
│   ├── 工作负载服务
│   └── 监控服务
│
├── K8S层 (client-go封装)
│   ├── 集群客户端管理器
│   ├── Deployment操作
│   ├── Pod操作
│   ├── Service操作
│   └── ConfigMap/Secret操作
│
└── 数据层 (GORM + MySQL)
    ├── 用户数据
    ├── 集群数据
    └── 审计日志
```

### 2.3 前端架构

```
前端应用
├── 页面层 (Pages)
│   ├── 登录页
│   ├── 集群管理
│   ├── 工作负载
│   ├── 服务与网络
│   └── 配置管理
│
├── 组件层 (Components)
│   ├── 布局组件 (Layout)
│   ├── 业务组件 (ResourceTable, YamlEditor)
│   └── 通用组件 (Modal, Form)
│
├── 服务层 (Services)
│   ├── API请求封装
│   └── WebSocket连接
│
└── 工具层 (Utils)
    ├── 请求拦截器
    ├── 认证工具
    └── 格式化工具
```

---

## 3. 系统设计

### 3.1 项目目录结构

#### 3.1.1 后端目录结构

```
k8s-manager-backend/
├── cmd/
│   └── server/
│       └── main.go                    # 程序入口
│
├── internal/
│   ├── api/                           # API层
│   │   ├── v1/
│   │   │   ├── auth/                  # 认证API
│   │   │   │   ├── login.go
│   │   │   │   └── register.go
│   │   │   ├── cluster/               # 集群管理API
│   │   │   │   ├── cluster.go
│   │   │   │   └── health.go
│   │   │   ├── workload/              # 工作负载API
│   │   │   │   ├── deployment.go
│   │   │   │   ├── pod.go
│   │   │   │   └── statefulset.go
│   │   │   ├── service/               # 服务API
│   │   │   │   ├── service.go
│   │   │   │   └── ingress.go
│   │   │   ├── config/                # 配置API
│   │   │   │   ├── configmap.go
│   │   │   │   └── secret.go
│   │   │   ├── storage/               # 存储API
│   │   │   │   ├── pv.go
│   │   │   │   └── pvc.go
│   │   │   ├── node/                  # 节点API
│   │   │   │   └── node.go
│   │   │   ├── namespace/             # 命名空间API
│   │   │   │   └── namespace.go
│   │   │   └── user/                  # 用户API
│   │   │       └── user.go
│   │   └── router.go                  # 路由注册
│   │
│   ├── service/                       # 服务层
│   │   ├── auth_service.go
│   │   ├── cluster_service.go
│   │   ├── deployment_service.go
│   │   ├── pod_service.go
│   │   ├── service_service.go
│   │   ├── configmap_service.go
│   │   ├── secret_service.go
│   │   ├── namespace_service.go
│   │   ├── node_service.go
│   │   └── user_service.go
│   │
│   ├── k8s/                           # K8S客户端封装
│   │   ├── client_manager.go         # 多集群客户端管理
│   │   ├── deployment.go             # Deployment操作
│   │   ├── pod.go                    # Pod操作
│   │   ├── service.go                # Service操作
│   │   ├── configmap.go              # ConfigMap操作
│   │   ├── secret.go                 # Secret操作
│   │   ├── namespace.go              # Namespace操作
│   │   ├── node.go                   # Node操作
│   │   ├── logs.go                   # 日志操作
│   │   ├── exec.go                   # Exec操作(WebShell)
│   │   └── metrics.go                # Metrics操作
│   │
│   ├── model/                         # 数据模型
│   │   ├── cluster.go                # 集群模型
│   │   ├── user.go                   # 用户模型
│   │   ├── audit_log.go              # 审计日志模型
│   │   └── permission.go             # 权限模型
│   │
│   ├── middleware/                    # 中间件
│   │   ├── auth.go                   # JWT认证中间件
│   │   ├── rbac.go                   # 权限控制中间件
│   │   ├── audit.go                  # 审计日志中间件
│   │   ├── cors.go                   # CORS中间件
│   │   ├── logger.go                 # 日志中间件
│   │   └── recovery.go               # 异常恢复中间件
│   │
│   ├── config/                        # 配置管理
│   │   └── config.go
│   │
│   └── database/                      # 数据库
│       ├── mysql.go                  # MySQL连接
│       └── redis.go                  # Redis连接
│
├── pkg/                               # 公共库
│   ├── response/                      # 统一响应
│   │   └── response.go
│   ├── jwt/                           # JWT工具
│   │   └── jwt.go
│   ├── encrypt/                       # 加密工具
│   │   └── encrypt.go
│   └── utils/                         # 工具函数
│       ├── string.go
│       └── time.go
│
├── config/                            # 配置文件
│   ├── config.yaml                   # 主配置
│   └── config.example.yaml           # 配置示例
│
├── scripts/                           # 脚本
│   ├── init.sql                      # 数据库初始化
│   └── build.sh                      # 构建脚本
│
├── docs/                              # 文档
│   ├── api.md                        # API文档
│   └── deployment.md                 # 部署文档
│
├── .gitignore
├── go.mod
├── go.sum
├── Dockerfile
└── README.md
```

#### 3.1.2 前端目录结构

```
k8s-manager-frontend/
├── public/
│   ├── index.html
│   └── favicon.ico
│
├── src/
│   ├── components/                    # 公共组件
│   │   ├── Layout/                   # 布局组件
│   │   │   ├── BasicLayout.tsx
│   │   │   ├── Header.tsx
│   │   │   └── Sidebar.tsx
│   │   ├── YamlEditor/               # YAML编辑器
│   │   │   └── index.tsx
│   │   ├── WebShell/                 # Web终端
│   │   │   └── index.tsx
│   │   ├── ResourceTable/            # 资源表格
│   │   │   └── index.tsx
│   │   ├── LogViewer/                # 日志查看器
│   │   │   └── index.tsx
│   │   └── Common/                   # 通用组件
│   │       ├── ConfirmModal.tsx
│   │       └── PageLoading.tsx
│   │
│   ├── pages/                         # 页面
│   │   ├── Login/                    # 登录页
│   │   │   └── index.tsx
│   │   ├── Cluster/                  # 集群管理
│   │   │   ├── ClusterList.tsx
│   │   │   ├── ClusterDetail.tsx
│   │   │   └── AddCluster.tsx
│   │   ├── Workload/                 # 工作负载
│   │   │   ├── Deployment/
│   │   │   │   ├── List.tsx
│   │   │   │   ├── Detail.tsx
│   │   │   │   ├── Create.tsx
│   │   │   │   └── Edit.tsx
│   │   │   ├── Pod/
│   │   │   │   ├── List.tsx
│   │   │   │   ├── Detail.tsx
│   │   │   │   ├── Logs.tsx
│   │   │   │   └── Shell.tsx
│   │   │   ├── StatefulSet/
│   │   │   ├── DaemonSet/
│   │   │   └── Job/
│   │   ├── Service/                  # 服务
│   │   │   ├── ServiceList.tsx
│   │   │   └── IngressList.tsx
│   │   ├── Config/                   # 配置
│   │   │   ├── ConfigMapList.tsx
│   │   │   └── SecretList.tsx
│   │   ├── Storage/                  # 存储
│   │   │   ├── PVList.tsx
│   │   │   └── PVCList.tsx
│   │   ├── Node/                     # 节点
│   │   │   ├── NodeList.tsx
│   │   │   └── NodeDetail.tsx
│   │   ├── Namespace/                # 命名空间
│   │   │   └── NamespaceList.tsx
│   │   ├── User/                     # 用户管理
│   │   │   └── UserList.tsx
│   │   └── AuditLog/                 # 审计日志
│   │       └── LogList.tsx
│   │
│   ├── services/                      # API服务
│   │   ├── auth.ts                   # 认证API
│   │   ├── cluster.ts                # 集群API
│   │   ├── deployment.ts             # Deployment API
│   │   ├── pod.ts                    # Pod API
│   │   ├── service.ts                # Service API
│   │   ├── configmap.ts              # ConfigMap API
│   │   ├── secret.ts                 # Secret API
│   │   ├── namespace.ts              # Namespace API
│   │   ├── node.ts                   # Node API
│   │   └── user.ts                   # User API
│   │
│   ├── store/                         # 状态管理
│   │   ├── index.ts
│   │   ├── userStore.ts              # 用户状态
│   │   └── clusterStore.ts           # 集群状态
│   │
│   ├── utils/                         # 工具函数
│   │   ├── request.ts                # 请求封装
│   │   ├── auth.ts                   # 认证工具
│   │   ├── format.ts                 # 格式化工具
│   │   └── constants.ts              # 常量定义
│   │
│   ├── types/                         # TypeScript类型定义
│   │   ├── cluster.ts
│   │   ├── deployment.ts
│   │   ├── pod.ts
│   │   └── user.ts
│   │
│   ├── styles/                        # 样式
│   │   ├── global.css
│   │   └── variables.css
│   │
│   ├── App.tsx                        # 应用根组件
│   ├── index.tsx                      # 入口文件
│   └── routes.tsx                     # 路由配置
│
├── .env                               # 环境变量
├── .env.development                   # 开发环境变量
├── .env.production                    # 生产环境变量
├── .gitignore
├── package.json
├── tsconfig.json
├── Dockerfile
└── README.md
```

### 3.2 核心模块设计

#### 3.2.1 集群客户端管理器

```go
// internal/k8s/client_manager.go
type ClusterManager struct {
    clients sync.Map // map[clusterName]*kubernetes.Clientset
    configs sync.Map // map[clusterName]*rest.Config
    mu      sync.RWMutex
}

// 单例模式
var Manager = &ClusterManager{}

// 添加集群
func (cm *ClusterManager) AddCluster(name, kubeconfig string) error

// 获取客户端
func (cm *ClusterManager) GetClient(name string) (*kubernetes.Clientset, error)

// 移除集群
func (cm *ClusterManager) RemoveCluster(name string)

// 测试连接
func (cm *ClusterManager) TestConnection(name string) error

// 获取集群信息
func (cm *ClusterManager) GetClusterInfo(name string) (*ClusterInfo, error)
```

#### 3.2.2 认证中间件

```go
// internal/middleware/auth.go
func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // 1. 从Header获取token
        token := c.GetHeader("Authorization")
        
        // 2. 验证token
        claims, err := jwt.ValidateToken(token)
        if err != nil {
            response.Error(c, 401, "未授权")
            c.Abort()
            return
        }
        
        // 3. 将用户信息存入Context
        c.Set("userID", claims.UserID)
        c.Set("username", claims.Username)
        
        c.Next()
    }
}
```

#### 3.2.3 审计日志中间件

```go
// internal/middleware/audit.go
func AuditMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // 记录请求开始时间
        startTime := time.Now()
        
        // 处理请求
        c.Next()
        
        // 记录审计日志
        log := &model.AuditLog{
            UserID:       c.GetInt64("userID"),
            Username:     c.GetString("username"),
            ClusterName:  c.Param("cluster"),
            Namespace:    c.Param("namespace"),
            ResourceType: extractResourceType(c.Request.URL.Path),
            ResourceName: c.Param("name"),
            Action:       c.Request.Method,
            Status:       getStatus(c.Writer.Status()),
            IP:           c.ClientIP(),
            UserAgent:    c.Request.UserAgent(),
            Duration:     time.Since(startTime).Milliseconds(),
        }
        
        // 异步保存
        go saveAuditLog(log)
    }
}
```

---

## 4. 数据库设计

### 4.1 数据库选择

- **主数据库**: MySQL 8.0+
- **缓存**: Redis 7.0+

### 4.2 表结构设计

#### 4.2.1 集群表 (clusters)

```sql
CREATE TABLE `clusters` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `name` varchar(100) NOT NULL COMMENT '集群名称(唯一)',
  `alias` varchar(100) DEFAULT NULL COMMENT '集群别名',
  `kubeconfig` text NOT NULL COMMENT 'kubeconfig内容(AES加密存储)',
  `api_server` varchar(255) DEFAULT NULL COMMENT 'API Server地址',
  `version` varchar(50) DEFAULT NULL COMMENT 'K8S版本',
  `node_count` int(11) DEFAULT 0 COMMENT '节点数量',
  `pod_count` int(11) DEFAULT 0 COMMENT 'Pod数量',
  `status` tinyint(4) DEFAULT 1 COMMENT '状态: 1-正常 0-异常',
  `last_check_time` timestamp NULL DEFAULT NULL COMMENT '最后检查时间',
  `description` varchar(500) DEFAULT NULL COMMENT '描述',
  `created_by` bigint(20) unsigned DEFAULT NULL COMMENT '创建人ID',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` timestamp NULL DEFAULT NULL COMMENT '删除时间(软删除)',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_name` (`name`),
  KEY `idx_status` (`status`),
  KEY `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='K8S集群表';
```

#### 4.2.2 用户表 (users)

```sql
CREATE TABLE `users` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `username` varchar(50) NOT NULL COMMENT '用户名(唯一)',
  `password` varchar(255) NOT NULL COMMENT '密码(bcrypt加密)',
  `nickname` varchar(100) DEFAULT NULL COMMENT '昵称',
  `email` varchar(100) DEFAULT NULL COMMENT '邮箱',
  `phone` varchar(20) DEFAULT NULL COMMENT '手机号',
  `avatar` varchar(500) DEFAULT NULL COMMENT '头像URL',
  `role` varchar(20) DEFAULT 'user' COMMENT '角色: admin-管理员, user-普通用户',
  `status` tinyint(4) DEFAULT 1 COMMENT '状态: 1-启用 0-禁用',
  `last_login_time` timestamp NULL DEFAULT NULL COMMENT '最后登录时间',
  `last_login_ip` varchar(50) DEFAULT NULL COMMENT '最后登录IP',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` timestamp NULL DEFAULT NULL COMMENT '删除时间(软删除)',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_username` (`username`),
  UNIQUE KEY `uk_email` (`email`),
  KEY `idx_status` (`status`),
  KEY `idx_role` (`role`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户表';

-- 初始管理员账号
INSERT INTO `users` (`username`, `password`, `nickname`, `role`, `status`) 
VALUES ('admin', '$2a$10$...', '系统管理员', 'admin', 1);
```

#### 4.2.3 审计日志表 (audit_logs)

```sql
CREATE TABLE `audit_logs` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `user_id` bigint(20) unsigned NOT NULL COMMENT '用户ID',
  `username` varchar(50) NOT NULL COMMENT '用户名',
  `cluster_name` varchar(100) DEFAULT NULL COMMENT '集群名称',
  `namespace` varchar(100) DEFAULT NULL COMMENT '命名空间',
  `resource_type` varchar(50) NOT NULL COMMENT '资源类型: deployment/pod/service等',
  `resource_name` varchar(255) DEFAULT NULL COMMENT '资源名称',
  `action` varchar(50) NOT NULL COMMENT '操作: GET/POST/PUT/DELETE/PATCH',
  `method` varchar(20) DEFAULT NULL COMMENT 'HTTP方法',
  `path` varchar(500) DEFAULT NULL COMMENT '请求路径',
  `status` varchar(20) DEFAULT NULL COMMENT '结果: success/failed',
  `status_code` int(11) DEFAULT NULL COMMENT 'HTTP状态码',
  `error_msg` text COMMENT '错误信息',
  `request_body` text COMMENT '请求体(可选)',
  `response_body` text COMMENT '响应体(可选)',
  `ip` varchar(50) DEFAULT NULL COMMENT '客户端IP',
  `user_agent` varchar(500) DEFAULT NULL COMMENT 'User-Agent',
  `duration` bigint(20) DEFAULT NULL COMMENT '请求耗时(毫秒)',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_cluster` (`cluster_name`),
  KEY `idx_resource` (`resource_type`, `resource_name`),
  KEY `idx_created_at` (`created_at`),
  KEY `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='审计日志表';

-- 按月分表策略(可选)
-- audit_logs_202411, audit_logs_202412, ...
```

#### 4.2.4 用户集群权限表 (user_cluster_permissions)

```sql
CREATE TABLE `user_cluster_permissions` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `user_id` bigint(20) unsigned NOT NULL COMMENT '用户ID',
  `cluster_id` bigint(20) unsigned NOT NULL COMMENT '集群ID',
  `namespaces` json DEFAULT NULL COMMENT '可访问的命名空间列表,空表示全部',
  `permissions` json DEFAULT NULL COMMENT '权限列表: ["read","write","delete"]',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_user_cluster` (`user_id`, `cluster_id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_cluster_id` (`cluster_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户集群权限表';
```

#### 4.2.5 配置表 (system_config)

```sql
CREATE TABLE `system_config` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `config_key` varchar(100) NOT NULL COMMENT '配置键(唯一)',
  `config_value` text NOT NULL COMMENT '配置值',
  `config_type` varchar(20) DEFAULT 'string' COMMENT '配置类型: string/int/bool/json',
  `description` varchar(500) DEFAULT NULL COMMENT '配置说明',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_config_key` (`config_key`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='系统配置表';
```

### 4.3 索引设计

#### 索引优化原则
1. 经常作为查询条件的字段建立索引
2. 经常需要排序的字段建立索引
3. 关联查询的外键字段建立索引
4. 控制单表索引数量（5个以内）
5. 避免在低基数字段建索引

#### 关键索引
- `clusters`: `uk_name`, `idx_status`
- `users`: `uk_username`, `uk_email`, `idx_status`
- `audit_logs`: `idx_user_id`, `idx_cluster`, `idx_created_at`（复合索引）

---

## 5. API接口设计

### 5.1 API规范

#### 5.1.1 RESTful设计原则

- 使用HTTP动词: GET(查询)、POST(创建)、PUT(更新)、DELETE(删除)
- URL使用名词复数形式
- 版本控制: `/api/v1/`
- 统一响应格式

#### 5.1.2 统一响应格式

**成功响应**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    // 数据内容
  }
}
```

**失败响应**:
```json
{
  "code": 40001,
  "message": "集群不存在",
  "data": null
}
```

**分页响应**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "total": 100,
    "page": 1,
    "page_size": 20,
    "items": [...]
  }
}
```

#### 5.1.3 错误码定义

| 错误码 | 说明 |
|--------|------|
| 0 | 成功 |
| 40001 | 请求参数错误 |
| 40101 | 未授权 |
| 40301 | 禁止访问 |
| 40401 | 资源不存在 |
| 40901 | 资源冲突 |
| 50001 | 服务器内部错误 |
| 50002 | 数据库错误 |
| 50003 | K8S API错误 |

### 5.2 接口列表

#### 5.2.1 认证接口

**1. 用户登录**
```
POST /api/v1/auth/login
```

请求体:
```json
{
  "username": "admin",
  "password": "123456"
}
```

响应:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": 1,
      "username": "admin",
      "nickname": "管理员",
      "role": "admin"
    }
  }
}
```

**2. 用户注册**
```
POST /api/v1/auth/register
```

**3. 刷新Token**
```
POST /api/v1/auth/refresh
```

**4. 登出**
```
POST /api/v1/auth/logout
```

#### 5.2.2 集群管理接口

**1. 获取集群列表**
```
GET /api/v1/clusters
```

Query参数:
- `page`: 页码（默认1）
- `page_size`: 每页数量（默认20）
- `status`: 状态筛选（可选）
- `keyword`: 名称关键字（可选）

响应:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "total": 5,
    "page": 1,
    "page_size": 20,
    "items": [
      {
        "id": 1,
        "name": "prod-cluster-1",
        "alias": "生产集群1",
        "api_server": "https://10.0.0.1:6443",
        "version": "v1.28.2",
        "node_count": 10,
        "pod_count": 150,
        "status": 1,
        "last_check_time": "2024-11-10T10:00:00Z",
        "created_at": "2024-01-01T00:00:00Z"
      }
    ]
  }
}
```

**2. 获取集群详情**
```
GET /api/v1/clusters/:name
```

响应:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": 1,
    "name": "prod-cluster-1",
    "alias": "生产集群1",
    "api_server": "https://10.0.0.1:6443",
    "version": "v1.28.2",
    "node_count": 10,
    "pod_count": 150,
    "namespace_count": 20,
    "status": 1,
    "components": {
      "etcd": "healthy",
      "scheduler": "healthy",
      "controller-manager": "healthy"
    },
    "resource_usage": {
      "cpu": {
        "total": "40 cores",
        "used": "25 cores",
        "usage_rate": 62.5
      },
      "memory": {
        "total": "160Gi",
        "used": "100Gi",
        "usage_rate": 62.5
      }
    }
  }
}
```

**3. 添加集群**
```
POST /api/v1/clusters
```

请求体:
```json
{
  "name": "test-cluster",
  "alias": "测试集群",
  "kubeconfig": "apiVersion: v1\nkind: Config\n...",
  "description": "测试环境集群"
}
```

**4. 更新集群**
```
PUT /api/v1/clusters/:name
```

**5. 删除集群**
```
DELETE /api/v1/clusters/:name
```

**6. 测试集群连接**
```
POST /api/v1/clusters/:name/test
```

响应:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "status": "success",
    "api_server": "https://10.0.0.1:6443",
    "version": "v1.28.2",
    "latency": 15
  }
}
```

#### 5.2.3 Deployment接口

**1. 获取Deployment列表**
```
GET /api/v1/clusters/:cluster/namespaces/:namespace/deployments
```

Query参数:
- `page`: 页码
- `page_size`: 每页数量
- `name`: 名称筛选（支持模糊匹配）
- `label_selector`: Label选择器（如: app=nginx）

响应:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "total": 10,
    "items": [
      {
        "name": "nginx-deployment",
        "namespace": "default",
        "replicas": 3,
        "available_replicas": 3,
        "ready_replicas": 3,
        "updated_replicas": 3,
        "images": ["nginx:1.19"],
        "labels": {
          "app": "nginx"
        },
        "created_at": "2024-01-01T00:00:00Z",
        "strategy": "RollingUpdate"
      }
    ]
  }
}
```

**2. 获取Deployment详情**
```
GET /api/v1/clusters/:cluster/namespaces/:namespace/deployments/:name
```

**3. 创建Deployment**
```
POST /api/v1/clusters/:cluster/namespaces/:namespace/deployments
```

请求体（支持两种模式）:

模式1: 表单模式
```json
{
  "name": "nginx-deployment",
  "replicas": 3,
  "image": "nginx:1.19",
  "labels": {
    "app": "nginx"
  },
  "env": [
    {
      "name": "ENV",
      "value": "prod"
    }
  ],
  "resources": {
    "limits": {
      "cpu": "1000m",
      "memory": "1Gi"
    },
    "requests": {
      "cpu": "500m",
      "memory": "512Mi"
    }
  }
}
```

模式2: YAML模式
```json
{
  "yaml": "apiVersion: apps/v1\nkind: Deployment\n..."
}
```

**4. 更新Deployment**
```
PUT /api/v1/clusters/:cluster/namespaces/:namespace/deployments/:name
```

**5. 删除Deployment**
```
DELETE /api/v1/clusters/:cluster/namespaces/:namespace/deployments/:name
```

**6. 扩缩容**
```
POST /api/v1/clusters/:cluster/namespaces/:namespace/deployments/:name/scale
```

请求体:
```json
{
  "replicas": 5
}
```

**7. 重启Deployment**
```
POST /api/v1/clusters/:cluster/namespaces/:namespace/deployments/:name/restart
```

**8. 获取Deployment事件**
```
GET /api/v1/clusters/:cluster/namespaces/:namespace/deployments/:name/events
```

**9. 获取Deployment YAML**
```
GET /api/v1/clusters/:cluster/namespaces/:namespace/deployments/:name/yaml
```

#### 5.2.4 Pod接口

**1. 获取Pod列表**
```
GET /api/v1/clusters/:cluster/namespaces/:namespace/pods
```

Query参数:
- `label_selector`: Label选择器
- `field_selector`: 字段选择器
- `phase`: Pod阶段（Running/Pending/Failed等）

**2. 获取Pod详情**
```
GET /api/v1/clusters/:cluster/namespaces/:namespace/pods/:name
```

**3. 删除Pod**
```
DELETE /api/v1/clusters/:cluster/namespaces/:namespace/pods/:name
```

**4. 获取Pod日志**
```
GET /api/v1/clusters/:cluster/namespaces/:namespace/pods/:name/logs
```

Query参数:
- `container`: 容器名称（必填）
- `tail_lines`: 尾部行数（默认100）
- `since_seconds`: 最近N秒的日志
- `timestamps`: 是否显示时间戳（true/false）
- `previous`: 是否获取上一个容器的日志

**5. 实时日志流（WebSocket）**
```
WebSocket: /api/v1/clusters/:cluster/namespaces/:namespace/pods/:name/logs/stream
```

Query参数同上

**6. Web终端（WebSocket）**
```
WebSocket: /api/v1/clusters/:cluster/namespaces/:namespace/pods/:name/shell
```

Query参数:
- `container`: 容器名称（必填）

**7. 获取Pod事件**
```
GET /api/v1/clusters/:cluster/namespaces/:namespace/pods/:name/events
```

**8. 获取Pod Metrics**
```
GET /api/v1/clusters/:cluster/namespaces/:namespace/pods/:name/metrics
```

响应:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "containers": [
      {
        "name": "nginx",
        "cpu": "10m",
        "memory": "50Mi"
      }
    ]
  }
}
```

#### 5.2.5 Service接口

**1. 获取Service列表**
```
GET /api/v1/clusters/:cluster/namespaces/:namespace/services
```

**2. 获取Service详情**
```
GET /api/v1/clusters/:cluster/namespaces/:namespace/services/:name
```

**3. 创建Service**
```
POST /api/v1/clusters/:cluster/namespaces/:namespace/services
```

**4. 更新Service**
```
PUT /api/v1/clusters/:cluster/namespaces/:namespace/services/:name
```

**5. 删除Service**
```
DELETE /api/v1/clusters/:cluster/namespaces/:namespace/services/:name
```

#### 5.2.6 ConfigMap接口

**1. 获取ConfigMap列表**
```
GET /api/v1/clusters/:cluster/namespaces/:namespace/configmaps
```

**2. 获取ConfigMap详情**
```
GET /api/v1/clusters/:cluster/namespaces/:namespace/configmaps/:name
```

**3. 创建ConfigMap**
```
POST /api/v1/clusters/:cluster/namespaces/:namespace/configmaps
```

请求体:
```json
{
  "name": "app-config",
  "data": {
    "database.host": "mysql.default.svc",
    "database.port": "3306"
  }
}
```

**4. 更新ConfigMap**
```
PUT /api/v1/clusters/:cluster/namespaces/:namespace/configmaps/:name
```

**5. 删除ConfigMap**
```
DELETE /api/v1/clusters/:cluster/namespaces/:namespace/configmaps/:name
```

#### 5.2.7 Secret接口

接口设计类似ConfigMap，但响应时需要Base64解码

#### 5.2.8 Namespace接口

**1. 获取Namespace列表**
```
GET /api/v1/clusters/:cluster/namespaces
```

**2. 创建Namespace**
```
POST /api/v1/clusters/:cluster/namespaces
```

**3. 删除Namespace**
```
DELETE /api/v1/clusters/:cluster/namespaces/:name
```

#### 5.2.9 Node接口

**1. 获取Node列表**
```
GET /api/v1/clusters/:cluster/nodes
```

**2. 获取Node详情**
```
GET /api/v1/clusters/:cluster/nodes/:name
```

**3. Cordon节点（标记不可调度）**
```
POST /api/v1/clusters/:cluster/nodes/:name/cordon
```

**4. Uncordon节点**
```
POST /api/v1/clusters/:cluster/nodes/:name/uncordon
```

**5. Drain节点（驱逐Pod）**
```
POST /api/v1/clusters/:cluster/nodes/:name/drain
```

**6. 获取Node Metrics**
```
GET /api/v1/clusters/:cluster/nodes/:name/metrics
```

#### 5.2.10 用户管理接口

**1. 获取用户列表**
```
GET /api/v1/users
```

**2. 创建用户**
```
POST /api/v1/users
```

**3. 更新用户**
```
PUT /api/v1/users/:id
```

**4. 删除用户**
```
DELETE /api/v1/users/:id
```

**5. 修改密码**
```
POST /api/v1/users/:id/password
```

#### 5.2.11 审计日志接口

**1. 获取审计日志列表**
```
GET /api/v1/audit-logs
```

Query参数:
- `user_id`: 用户ID
- `cluster_name`: 集群名称
- `resource_type`: 资源类型
- `action`: 操作类型
- `status`: 状态（success/failed）
- `start_time`: 开始时间
- `end_time`: 结束时间
- `page`: 页码
- `page_size`: 每页数量

---

## 6. 前端页面设计

### 6.1 整体布局

```
┌─────────────────────────────────────────────────────┐
│  Header (顶部导航栏)                                   │
│  [Logo] [集群选择] [命名空间选择]  [用户] [设置] [退出]  │
├──────────┬──────────────────────────────────────────┤
│          │                                          │
│  Sidebar │         Content (主内容区)                │
│  (侧边栏)  │                                          │
│          │                                          │
│  - 集群   │                                          │
│  - 工作负载│                                          │
│  - 服务   │                                          │
│  - 配置   │                                          │
│  - 存储   │                                          │
│  - 节点   │                                          │
│  - 用户   │                                          │
│  - 日志   │                                          │
│          │                                          │
└──────────┴──────────────────────────────────────────┘
```

### 6.2 主要页面

#### 6.2.1 登录页面

**路由**: `/login`

**功能**:
- 用户名密码登录
- 记住密码（可选）
- 登录后跳转到首页

**UI组件**:
- Card容器
- Form表单
- Input输入框
- Button按钮

#### 6.2.2 集群管理页面

**路由**: `/clusters`

**功能**:
- 集群列表展示（表格）
- 添加集群（弹窗）
- 编辑集群
- 删除集群（二次确认）
- 测试连接
- 查看集群详情

**表格列**:
- 集群名称
- API Server地址
- 版本
- 节点数
- Pod数
- 状态（健康/异常）
- 最后检查时间
- 操作

#### 6.2.3 Deployment列表页面

**路由**: `/clusters/:cluster/workloads/deployments`

**功能**:
- Deployment列表展示
- 搜索框（按名称搜索）
- 命名空间筛选
- 创建Deployment
- 查看详情
- 编辑
- 删除（二次确认）
- 扩缩容（弹窗输入副本数）
- 重启

**表格列**:
- 名称
- 命名空间
- 镜像
- 副本数（期望/当前/可用）
- 状态
- 创建时间
- 操作

#### 6.2.4 Deployment详情页面

**路由**: `/clusters/:cluster/workloads/deployments/:namespace/:name`

**功能**:
使用Tab标签页展示不同信息

**Tab1: 基本信息**
- Deployment基本信息
- Selector
- Labels
- Annotations
- 镜像列表
- 环境变量
- 资源限制
- 更新策略

**Tab2: Pod列表**
- 该Deployment管理的Pod列表
- 可点击进入Pod详情

**Tab3: 事件**
- 相关事件列表
- 时间、类型、原因、消息

**Tab4: YAML**
- YAML编辑器
- 支持编辑和应用更新

#### 6.2.5 Pod列表页面

**路由**: `/clusters/:cluster/workloads/pods`

**功能**:
- Pod列表展示
- Label筛选
- 状态筛选
- 查看详情
- 查看日志
- 进入终端
- 删除

**表格列**:
- 名称
- 命名空间
- 状态
- Ready状态
- 重启次数
- 节点
- IP
- 创建时间
- 操作

#### 6.2.6 Pod详情页面

**路由**: `/clusters/:cluster/workloads/pods/:namespace/:name`

**Tab1: 基本信息**
- Pod详细信息
- 容器列表
- 卷信息
- 节点信息

**Tab2: 容器**
- 容器列表
- 每个容器的详细信息
- 资源使用情况

**Tab3: 事件**
- Pod相关事件

**Tab4: YAML**
- Pod YAML

#### 6.2.7 Pod日志页面

**路由**: `/clusters/:cluster/workloads/pods/:namespace/:name/logs`

**功能**:
- 容器选择下拉框
- 实时日志流（WebSocket）
- 历史日志加载
- 自动滚动开关
- 时间戳显示开关
- 行数选择（100/500/1000/全部）
- 日志搜索（客户端搜索）
- 日志下载
- 清空显示

**UI设计**:
- 顶部工具栏（容器选择、配置项、操作按钮）
- 日志显示区域（黑色背景、等宽字体）

#### 6.2.8 Web终端页面

**路由**: `/clusters/:cluster/workloads/pods/:namespace/:name/shell`

**功能**:
- 容器选择
- 使用xterm.js实现Web终端
- 支持复制粘贴
- 支持全屏

#### 6.2.9 Service列表页面

**路由**: `/clusters/:cluster/services`

**功能**:
- Service列表展示
- 创建、编辑、删除
- 查看Endpoints

#### 6.2.10 ConfigMap列表页面

**路由**: `/clusters/:cluster/configs/configmaps`

**功能**:
- ConfigMap列表
- 创建、编辑、删除
- Key-Value编辑器

#### 6.2.11 Secret列表页面

**路由**: `/clusters/:cluster/configs/secrets`

**功能**:
- Secret列表
- 敏感信息脱敏显示（***）
- 点击查看/隐藏
- 创建、编辑、删除

#### 6.2.12 Node列表页面

**路由**: `/clusters/:cluster/nodes`

**功能**:
- Node列表
- 资源使用率图表
- Cordon/Uncordon
- 查看详情
- Pod分布

**表格列**:
- 名称
- 状态
- 角色
- CPU使用率（进度条）
- 内存使用率（进度条）
- Pod数量
- 版本
- 操作

#### 6.2.13 Namespace管理页面

**路由**: `/clusters/:cluster/namespaces`

**功能**:
- Namespace列表
- 创建、删除
- 资源配额显示

#### 6.2.14 用户管理页面

**路由**: `/users`

**功能**:
- 用户列表
- 添加用户
- 编辑用户
- 删除用户
- 重置密码
- 权限分配

#### 6.2.15 审计日志页面

**路由**: `/audit-logs`

**功能**:
- 日志列表
- 高级筛选
  - 用户
  - 集群
  - 资源类型
  - 操作类型
  - 时间范围
- 日志详情（弹窗）
- 导出

### 6.3 公共组件

#### 6.3.1 YamlEditor组件

**功能**:
- 使用Monaco Editor或Ace Editor
- YAML语法高亮
- 语法验证
- 格式化
- 全屏编辑

#### 6.3.2 ResourceTable组件

**功能**:
- 通用资源表格组件
- 支持分页
- 支持排序
- 支持筛选
- 批量操作

#### 6.3.3 WebShell组件

**功能**:
- 基于xterm.js
- WebSocket连接
- 支持复制粘贴
- 全屏模式

#### 6.3.4 LogViewer组件

**功能**:
- 实时日志流
- 虚拟滚动（大量日志性能优化）
- 搜索高亮
- 时间戳显示

---

## 7. 开发计划

### 7.1 开发周期

**总计**: 12-16周（3-4个月）

**工作模式**: 
- 工作日晚上: 2-3小时
- 周末: 8-10小时
- 每周预计: 20-25小时

### 7.2 详细计划

#### Week 1-2: 基础框架 (Foundation)

**目标**: 搭建项目基础框架，完成认证系统

**后端任务**:
- [x] 初始化Go项目
- [x] 配置管理模块
- [x] 数据库连接（MySQL + Redis）
- [x] JWT认证中间件
- [x] 统一响应格式
- [x] 用户登录/注册API
- [x] CORS中间件
- [x] 日志中间件

**前端任务**:
- [x] 初始化React项目
- [x] 安装依赖（antd, axios, react-router-dom等）
- [x] 配置请求拦截器
- [x] 实现登录页面
- [x] 实现基础布局（Header + Sidebar）
- [x] 路由配置

**验收标准**:
- ✅ 可以正常登录
- ✅ Token自动携带
- ✅ 基础布局显示正常

**AI提示词示例**:
```
帮我用Go + Gin实现一个完整的JWT认证系统，包括：
1. JWT token生成和验证
2. 登录API
3. 认证中间件
4. 用户信息从token解析并存入Context
```

---

#### Week 3: 集群管理功能

**目标**: 完成集群的增删改查和连接管理

**后端任务**:
- [ ] 实现ClusterManager（多集群客户端管理器）
- [ ] 集群CRUD API
- [ ] 集群连接测试API
- [ ] 集群信息获取（版本、节点数等）
- [ ] kubeconfig加密存储

**前端任务**:
- [ ] 集群列表页面
- [ ] 添加集群弹窗
- [ ] 集群详情页面
- [ ] 集群切换下拉框（全局状态）
- [ ] 命名空间切换下拉框

**验收标准**:
- ✅ 可以添加真实K8S集群
- ✅ 连接测试成功
- ✅ 可以切换集群
- ✅ 显示集群基本信息

**AI提示词示例**:
```
帮我用client-go实现一个K8S集群客户端管理器，要求：
1. 支持同时管理多个集群
2. 从kubeconfig字符串创建clientset
3. 线程安全（使用sync.Map）
4. 提供连接测试方法
5. 缓存客户端避免重复创建
```

---

#### Week 4-5: Deployment管理

**目标**: 完成Deployment的完整管理功能

**后端任务**:
- [ ] Deployment CRUD操作封装
- [ ] Deployment列表API（支持分页、筛选）
- [ ] Deployment详情API
- [ ] Deployment创建API（表单+YAML两种模式）
- [ ] Deployment更新API
- [ ] Deployment删除API
- [ ] Deployment扩缩容API
- [ ] Deployment重启API
- [ ] Deployment事件API
- [ ] Deployment YAML获取API

**前端任务**:
- [ ] Deployment列表页面
- [ ] Deployment详情页面（Tab页）
  - [ ] 基本信息Tab
  - [ ] Pod列表Tab
  - [ ] 事件Tab
  - [ ] YAML Tab
- [ ] Deployment创建弹窗
  - [ ] 表单模式
  - [ ] YAML模式切换
- [ ] Deployment编辑功能
- [ ] 扩缩容弹窗
- [ ] 二次确认删除

**验收标准**:
- ✅ 可以查看Deployment列表
- ✅ 可以创建新Deployment
- ✅ 可以更新Deployment
- ✅ 可以扩缩容
- ✅ 可以重启
- ✅ 可以删除
- ✅ YAML编辑器正常工作

**AI提示词示例**:
```
帮我用client-go实现Deployment的扩缩容功能，要求：
1. 接收参数: namespace, name, replicas
2. 获取当前Deployment
3. 更新spec.replicas字段
4. 应用更新
5. 返回更新后的Deployment对象
6. 完整的错误处理
```

---

#### Week 6-7: Pod管理与日志

**目标**: 完成Pod查看、日志查询、WebShell

**后端任务**:
- [ ] Pod列表API（支持Label筛选）
- [ ] Pod详情API
- [ ] Pod删除API
- [ ] Pod日志API（历史日志）
- [ ] Pod日志流API（WebSocket）
- [ ] Pod事件API
- [ ] Pod Metrics API
- [ ] WebShell API（WebSocket + Exec）

**前端任务**:
- [ ] Pod列表页面
- [ ] Pod详情页面
- [ ] 日志查看页面
  - [ ] 容器选择
  - [ ] 实时日志（WebSocket）
  - [ ] 历史日志加载
  - [ ] 日志搜索
  - [ ] 日志下载
- [ ] WebShell组件
  - [ ] 集成xterm.js
  - [ ] WebSocket连接
  - [ ] 全屏模式

**验收标准**:
- ✅ 可以查看Pod列表
- ✅ 可以查看Pod详情
- ✅ 可以查看实时日志
- ✅ 可以进入容器终端
- ✅ 终端支持交互

**AI提示词示例**:
```
帮我用Go + Gin + WebSocket实现Pod实时日志流功能，要求：
1. 升级HTTP到WebSocket
2. 使用client-go的GetLogs获取日志流
3. 将日志实时推送到WebSocket客户端
4. 处理客户端断开连接
5. 支持容器选择参数
```

---

#### Week 8: Service与ConfigMap/Secret

**目标**: 完成Service、ConfigMap、Secret管理

**后端任务**:
- [ ] Service CRUD API
- [ ] ConfigMap CRUD API
- [ ] Secret CRUD API（Base64编解码）

**前端任务**:
- [ ] Service列表和详情页面
- [ ] ConfigMap列表和编辑页面
- [ ] Secret列表和编辑页面（敏感信息脱敏）

**验收标准**:
- ✅ 可以管理Service
- ✅ 可以管理ConfigMap
- ✅ 可以管理Secret
- ✅ Secret显示时脱敏

---

#### Week 9: Namespace与Node管理

**目标**: 完成Namespace和Node管理

**后端任务**:
- [ ] Namespace CRUD API
- [ ] Node列表API
- [ ] Node详情API
- [ ] Node Cordon/Uncordon API
- [ ] Node Drain API
- [ ] Node Label管理API
- [ ] Node Taint管理API
- [ ] Node Metrics API

**前端任务**:
- [ ] Namespace管理页面
- [ ] Node列表页面
- [ ] Node详情页面
- [ ] 资源使用率图表（ECharts）

**验收标准**:
- ✅ 可以管理Namespace
- ✅ 可以查看Node列表
- ✅ 可以Cordon/Uncordon节点
- ✅ 显示资源使用率

---

#### Week 10: 审计日志与权限

**目标**: 完成审计日志和基础权限控制

**后端任务**:
- [ ] 审计日志中间件（记录所有操作）
- [ ] 审计日志查询API
- [ ] RBAC中间件（基础权限控制）
- [ ] 用户管理API

**前端任务**:
- [ ] 审计日志查询页面
- [ ] 用户管理页面
- [ ] 权限提示

**验收标准**:
- ✅ 所有操作被记录
- ✅ 可以查询审计日志
- ✅ 管理员可以管理用户

---

#### Week 11: 其他工作负载

**目标**: 完成StatefulSet、DaemonSet、Job等

**后端任务**:
- [ ] StatefulSet CRUD API
- [ ] DaemonSet CRUD API
- [ ] Job CRUD API
- [ ] CronJob CRUD API

**前端任务**:
- [ ] StatefulSet页面
- [ ] DaemonSet页面
- [ ] Job页面
- [ ] CronJob页面

---

#### Week 12: 优化与测试

**目标**: 完善功能、修复Bug、性能优化

**任务**:
- [ ] 代码重构
- [ ] 性能优化
- [ ] UI/UX优化
- [ ] 全面测试
- [ ] Bug修复
- [ ] 编写文档

---

### 7.3 开发建议

#### 每日开发流程

**上午时段（或工作日晚上）**:
1. 设计API接口（10-15分钟）
2. AI生成后端代码框架（15分钟）
3. 实现业务逻辑（30-60分钟）
4. Postman测试API（15分钟）

**下午时段（或周末）**:
1. 设计页面结构（10分钟）
2. AI生成前端组件框架（15分钟）
3. 实现页面逻辑（30-60分钟）
4. 前后端联调（15-30分钟）

#### AI辅助技巧

**代码生成**:
```
请帮我生成一个完整的[功能]实现，包括：
1. [具体需求1]
2. [具体需求2]
3. [具体需求3]
要求：完整的错误处理、代码注释
```

**代码审查**:
```
请帮我审查这段代码，检查：
1. 是否有潜在的bug
2. 是否有性能问题
3. 是否符合最佳实践
[粘贴代码]
```

**调试帮助**:
```
我遇到了这个错误：
[错误信息]

我的代码：
[代码片段]

我的环境：
[环境信息]

请帮我分析问题并给出解决方案
```

---

## 8. AI辅助开发环境配置

### 8.1 环境架构

本项目采用AI辅助开发模式，推荐使用Claude Code在openEuler环境上进行开发。为了保证系统环境的干净和隔离性，所有测试依赖（MySQL、Redis等）都运行在Docker容器中。

#### 8.1.1 架构图

```
openEuler系统
├── Claude Code (AI开发工具)
│   └── 工作目录: ~/k8s-manager/
│       ├── backend/              # 后端代码（本地开发）
│       ├── frontend/             # 前端代码（本地开发）
│       └── scripts/              # 脚本文件
│
├── Docker容器环境 (完全隔离的测试环境)
│   ├── MySQL测试容器  (端口: 13306)
│   ├── Redis测试容器  (端口: 16379)
│   └── phpMyAdmin容器 (端口: 18080)
│
└── 系统环境 (保持干净)
    ├── 仅安装: Docker、Go、Node.js
    └── 不安装: MySQL、Redis等服务
```

#### 8.1.2 为什么使用Docker容器？

**优势**：
- ✅ **环境隔离** - 测试环境与系统环境完全隔离
- ✅ **快速重置** - 一键删除所有测试数据，重新初始化
- ✅ **避免冲突** - 使用非标准端口，不与系统服务冲突
- ✅ **可复现** - 团队成员环境完全一致
- ✅ **接近生产** - 测试环境与生产容器化部署一致

**测试环境 vs 生产环境对比**：
```
测试环境（Docker容器）         生产环境
├── MySQL: :13306       →     MySQL: :3306
├── Redis: :16379       →     Redis: :6379
└── DB: k8s_manager_test →    DB: k8s_manager
```

### 8.2 环境准备

#### 8.2.1 安装Docker

```bash
# 1. 安装Docker
sudo yum install -y docker

# 2. 启动Docker服务
sudo systemctl start docker
sudo systemctl enable docker

# 3. 安装Docker Compose
sudo curl -L "https://github.com/docker/compose/releases/download/v2.23.0/docker-compose-$(uname -s)-$(uname -m)" \
  -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose

# 4. 验证安装
docker --version
docker-compose --version

# 5. 将当前用户加入docker组（避免每次sudo）
sudo usermod -aG docker $USER
# 注意：需要重新登录使生效
```

#### 8.2.2 安装开发工具

```bash
# 安装Go（如果未安装）
sudo yum install -y golang

# 配置Go环境
export GOPROXY=https://goproxy.cn,direct
export GO111MODULE=on
echo 'export GOPROXY=https://goproxy.cn,direct' >> ~/.bashrc
echo 'export GO111MODULE=on' >> ~/.bashrc

# 安装Node.js和npm（如果需要前端开发）
sudo yum install -y nodejs npm
```

#### 8.2.3 创建项目目录

```bash
# 创建项目根目录
mkdir -p ~/k8s-manager
cd ~/k8s-manager

# 创建子目录
mkdir -p backend frontend scripts
```

### 8.3 Docker测试环境配置

#### 8.3.1 创建Docker Compose配置文件

在 `~/k8s-manager/docker-compose.test.yml`:

```yaml
version: '3.8'

services:
  # MySQL测试数据库
  mysql-test:
    image: mysql:8.0
    container_name: k8s-mgr-mysql-test
    environment:
      MYSQL_ROOT_PASSWORD: test123456
      MYSQL_DATABASE: k8s_manager_test
    ports:
      - "13306:3306"  # 使用13306避免与系统MySQL冲突
    volumes:
      - mysql-test-data:/var/lib/mysql
      - ./scripts/init.sql:/docker-entrypoint-initdb.d/init.sql
    command: --default-authentication-plugin=mysql_native_password
    networks:
      - test-network
    restart: unless-stopped
    healthcheck:
      test: ["CMD", "mysqladmin", "ping", "-h", "localhost"]
      interval: 10s
      timeout: 5s
      retries: 5

  # Redis测试缓存
  redis-test:
    image: redis:7-alpine
    container_name: k8s-mgr-redis-test
    ports:
      - "16379:6379"  # 使用16379避免与系统Redis冲突
    networks:
      - test-network
    restart: unless-stopped
    healthcheck:
      test: ["CMD", "redis-cli", "ping"]
      interval: 10s
      timeout: 5s
      retries: 5

  # phpMyAdmin (可选，方便查看和管理数据库)
  phpmyadmin-test:
    image: phpmyadmin:latest
    container_name: k8s-mgr-phpmyadmin-test
    environment:
      PMA_HOST: mysql-test
      PMA_PORT: 3306
      PMA_USER: root
      PMA_PASSWORD: test123456
    ports:
      - "18080:80"
    depends_on:
      - mysql-test
    networks:
      - test-network
    restart: unless-stopped

  # Redis Commander (可选，方便查看Redis数据)
  redis-commander-test:
    image: rediscommander/redis-commander:latest
    container_name: k8s-mgr-redis-commander-test
    environment:
      REDIS_HOSTS: local:redis-test:6379
    ports:
      - "18081:8081"
    depends_on:
      - redis-test
    networks:
      - test-network
    restart: unless-stopped

volumes:
  mysql-test-data:
    name: k8s-mgr-mysql-test-data

networks:
  test-network:
    name: k8s-mgr-test-network
    driver: bridge
```

#### 8.3.2 创建数据库初始化脚本

在 `~/k8s-manager/scripts/init.sql`:

```sql
-- 创建测试数据库
CREATE DATABASE IF NOT EXISTS k8s_manager_test 
DEFAULT CHARACTER SET utf8mb4 
COLLATE utf8mb4_unicode_ci;

USE k8s_manager_test;

-- 用户表
CREATE TABLE IF NOT EXISTS `users` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `username` varchar(50) NOT NULL COMMENT '用户名',
  `password` varchar(255) NOT NULL COMMENT '密码(bcrypt)',
  `nickname` varchar(100) DEFAULT NULL COMMENT '昵称',
  `email` varchar(100) DEFAULT NULL COMMENT '邮箱',
  `role` varchar(20) DEFAULT 'user' COMMENT '角色',
  `status` tinyint(4) DEFAULT 1 COMMENT '状态',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_username` (`username`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户表';

-- 插入测试管理员账号
-- 用户名: admin, 密码: admin123
INSERT INTO `users` (`username`, `password`, `nickname`, `role`) 
VALUES ('admin', '$2a$10$N.zmdr9k7uOCQb376NoUnuTJ8iAt6Z5EHsM8lE9lBOsl7iAt6Z5EH', '测试管理员', 'admin');

-- 集群表
CREATE TABLE IF NOT EXISTS `clusters` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(100) NOT NULL,
  `kubeconfig` text NOT NULL,
  `api_server` varchar(255) DEFAULT NULL,
  `version` varchar(50) DEFAULT NULL,
  `status` tinyint(4) DEFAULT 1,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_name` (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='集群表';

-- 审计日志表
CREATE TABLE IF NOT EXISTS `audit_logs` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `user_id` bigint(20) unsigned NOT NULL,
  `username` varchar(50) NOT NULL,
  `cluster_name` varchar(100) DEFAULT NULL,
  `namespace` varchar(100) DEFAULT NULL,
  `resource_type` varchar(50) NOT NULL,
  `resource_name` varchar(255) DEFAULT NULL,
  `action` varchar(50) NOT NULL,
  `status` varchar(20) DEFAULT NULL,
  `ip` varchar(50) DEFAULT NULL,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='审计日志表';
```

### 8.4 便捷管理脚本

#### 8.4.1 启动测试环境脚本

在 `~/k8s-manager/test-env-start.sh`:

```bash
#!/bin/bash

echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "🚀 启动K8S Manager测试环境"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

cd ~/k8s-manager

# 检查Docker是否运行
if ! docker info > /dev/null 2>&1; then
    echo "❌ Docker未运行，请先启动Docker服务"
    echo "   sudo systemctl start docker"
    exit 1
fi

# 启动Docker容器
echo "📦 启动Docker容器..."
docker-compose -f docker-compose.test.yml up -d

# 等待MySQL就绪
echo "⏳ 等待MySQL启动..."
timeout=30
counter=0
until docker exec k8s-mgr-mysql-test mysqladmin ping -h localhost --silent > /dev/null 2>&1; do
    counter=$((counter+1))
    if [ $counter -gt $timeout ]; then
        echo "❌ MySQL启动超时"
        exit 1
    fi
    echo "   等待中... ($counter/$timeout)"
    sleep 1
done

# 检查容器状态
echo ""
echo "📊 容器状态："
docker-compose -f docker-compose.test.yml ps

echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "✅ 测试环境启动成功！"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""
echo "📦 服务访问地址："
echo "   MySQL:          localhost:13306"
echo "   Redis:          localhost:16379"
echo "   phpMyAdmin:     http://localhost:18080"
echo "   Redis Commander: http://localhost:18081"
echo ""
echo "🔑 MySQL连接信息："
echo "   用户名: root"
echo "   密码:   test123456"
echo "   数据库: k8s_manager_test"
echo ""
echo "💡 测试账号："
echo "   用户名: admin"
echo "   密码:   admin123"
echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
```

#### 8.4.2 停止测试环境脚本

在 `~/k8s-manager/test-env-stop.sh`:

```bash
#!/bin/bash

echo "🛑 停止测试环境..."

cd ~/k8s-manager

# 停止容器
docker-compose -f docker-compose.test.yml down

echo "✅ 测试环境已停止"
echo "💡 提示: 数据已保留，下次启动会恢复"
echo "        如需删除所有数据，请使用: ./test-env-reset.sh"
```

#### 8.4.3 重置测试环境脚本

在 `~/k8s-manager/test-env-reset.sh`:

```bash
#!/bin/bash

echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "⚠️  警告: 即将删除所有测试数据！"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
read -p "确认重置测试环境？(yes/no): " confirm

if [ "$confirm" != "yes" ]; then
    echo "❌ 取消操作"
    exit 0
fi

echo "🔄 重置测试环境..."

cd ~/k8s-manager

# 停止并删除容器和数据卷
docker-compose -f docker-compose.test.yml down -v

echo "⏳ 重新启动环境..."

# 重新启动
docker-compose -f docker-compose.test.yml up -d

# 等待MySQL就绪
sleep 10

echo ""
echo "✅ 测试环境已重置为初始状态！"
echo "   所有数据已删除并重新初始化"
```

#### 8.4.4 查看日志脚本

在 `~/k8s-manager/test-env-logs.sh`:

```bash
#!/bin/bash

echo "📋 查看测试环境日志 (Ctrl+C退出)"
echo ""

cd ~/k8s-manager

# 查看所有容器日志
docker-compose -f docker-compose.test.yml logs -f
```

#### 8.4.5 进入数据库脚本

在 `~/k8s-manager/test-env-mysql.sh`:

```bash
#!/bin/bash

echo "🗄️  连接到MySQL测试数据库..."
echo ""

docker exec -it k8s-mgr-mysql-test mysql -uroot -ptest123456 k8s_manager_test
```

#### 8.4.6 赋予执行权限

```bash
chmod +x ~/k8s-manager/test-env-*.sh
```

### 8.5 后端测试配置

#### 8.5.1 创建测试配置文件

在 `~/k8s-manager/backend/config/config.test.yaml`:

```yaml
server:
  port: 8080
  mode: debug  # 开发模式

database:
  host: localhost
  port: 13306              # Docker映射的端口
  user: root
  password: test123456
  dbname: k8s_manager_test
  max_idle_conns: 10
  max_open_conns: 100

redis:
  host: localhost
  port: 16379              # Docker映射的端口
  password: ""
  db: 0

jwt:
  secret: test-jwt-secret-key-for-testing
  expire: 7200  # 2小时

log:
  level: debug
  file: logs/test.log
```

### 8.6 AI辅助开发工作流

#### 8.6.1 完整开发流程

```bash
# 第1步: 启动测试环境
cd ~/k8s-manager
./test-env-start.sh

# 第2步: 使用Claude Code开发
# 在Claude Code中编写代码

# 第3步: 运行后端测试
cd ~/k8s-manager/backend
go run cmd/server/main.go

# 第4步: 在另一个终端测试API
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin123"}'

# 第5步: 查看数据库数据（如果需要）
./test-env-mysql.sh

# 第6步: 查看日志（如果需要）
./test-env-logs.sh

# 第7步: 开发完成，停止测试环境
./test-env-stop.sh

# 第8步: 如果需要重新测试，重置环境
./test-env-reset.sh
```

#### 8.6.2 Claude Code开发示例

**示例1: 开发用户认证模块**

```
你对Claude Code说：
"请在 ~/k8s-manager/backend 目录下创建用户认证模块，
包括登录、注册、JWT token生成和验证，
使用测试数据库 localhost:13306"

Claude Code会：
1. 创建项目结构
2. 实现用户认证逻辑
3. 配置数据库连接
4. 提供测试代码

你执行测试：
cd ~/k8s-manager/backend
go run cmd/server/main.go

curl测试API是否正常工作
```

**示例2: 开发集群管理模块**

```
你对Claude Code说：
"实现K8S集群管理功能，包括添加集群、列表查询、连接测试"

Claude Code会：
1. 使用client-go实现K8S客户端管理
2. 实现CRUD API
3. 数据持久化到测试数据库
4. 提供测试用的kubeconfig示例

你执行测试：
使用真实的kubeconfig测试集群连接功能
```

### 8.7 常见问题与解决方案

#### Q1: Docker容器启动失败
```bash
# 检查Docker服务状态
sudo systemctl status docker

# 启动Docker服务
sudo systemctl start docker

# 查看容器日志
docker-compose -f docker-compose.test.yml logs
```

#### Q2: 端口已被占用
```bash
# 查看端口占用
sudo netstat -tlnp | grep 13306
sudo netstat -tlnp | grep 16379

# 修改docker-compose.test.yml中的端口映射
# 例如: "13306:3306" 改为 "23306:3306"
```

#### Q3: 数据库连接失败
```bash
# 检查容器是否运行
docker ps | grep k8s-mgr-mysql-test

# 测试数据库连接
docker exec -it k8s-mgr-mysql-test mysql -uroot -ptest123456 -e "SELECT 1"

# 检查后端配置文件中的端口是否正确
cat backend/config/config.test.yaml
```

#### Q4: 容器内存或磁盘不足
```bash
# 查看Docker磁盘使用
docker system df

# 清理未使用的镜像和容器
docker system prune -a

# 清理数据卷
docker volume prune
```

### 8.8 性能优化建议

#### 8.8.1 Docker性能优化

```bash
# 限制MySQL容器内存使用
# 在docker-compose.test.yml中添加:
services:
  mysql-test:
    deploy:
      resources:
        limits:
          memory: 1G
        reservations:
          memory: 512M
```

#### 8.8.2 开发效率优化

```bash
# 使用Air实现Go代码热重载
go install github.com/cosmtrek/air@latest

# 创建.air.toml配置文件
air init

# 使用air运行
air
```

## 9. 部署方案

### 9.1 开发环境（已在第8章详细说明）

开发环境使用Docker容器隔离，详见第8章"AI辅助开发环境配置"。

### 9.2 生产环境部署

### 8.2 生产部署

#### 8.2.1 Docker构建

**后端Dockerfile**:
```dockerfile
# 构建阶段
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o server cmd/server/main.go

# 运行阶段
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/server .
COPY --from=builder /app/config/config.yaml ./config/
EXPOSE 8080
CMD ["./server"]
```

**前端Dockerfile**:
```dockerfile
# 构建阶段
FROM node:18-alpine AS builder
WORKDIR /app
COPY package*.json ./
RUN npm ci
COPY . .
RUN npm run build

# Nginx运行阶段
FROM nginx:alpine
COPY --from=builder /app/build /usr/share/nginx/html
COPY nginx.conf /etc/nginx/conf.d/default.conf
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
```

**前端Nginx配置**:
```nginx
server {
    listen 80;
    server_name _;
    
    location / {
        root /usr/share/nginx/html;
        index index.html;
        try_files $uri $uri/ /index.html;
    }
    
    location /api {
        proxy_pass http://backend:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
    
    location /ws {
        proxy_pass http://backend:8080;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";
    }
}
```

#### 8.2.2 Docker Compose部署

**docker-compose.yml**:
```yaml
version: '3.8'

services:
  mysql:
    image: mysql:8.0
    container_name: k8s-manager-mysql
    environment:
      MYSQL_ROOT_PASSWORD: ${MYSQL_ROOT_PASSWORD}
      MYSQL_DATABASE: k8s_manager
    volumes:
      - mysql-data:/var/lib/mysql
      - ./scripts/init.sql:/docker-entrypoint-initdb.d/init.sql
    ports:
      - "3306:3306"
    networks:
      - k8s-manager-net
    restart: unless-stopped

  redis:
    image: redis:7-alpine
    container_name: k8s-manager-redis
    ports:
      - "6379:6379"
    networks:
      - k8s-manager-net
    restart: unless-stopped

  backend:
    build:
      context: ./k8s-manager-backend
      dockerfile: Dockerfile
    container_name: k8s-manager-backend
    environment:
      DB_HOST: mysql
      DB_PORT: 3306
      DB_USER: root
      DB_PASSWORD: ${MYSQL_ROOT_PASSWORD}
      DB_NAME: k8s_manager
      REDIS_HOST: redis
      REDIS_PORT: 6379
    ports:
      - "8080:8080"
    depends_on:
      - mysql
      - redis
    networks:
      - k8s-manager-net
    restart: unless-stopped

  frontend:
    build:
      context: ./k8s-manager-frontend
      dockerfile: Dockerfile
    container_name: k8s-manager-frontend
    ports:
      - "80:80"
    depends_on:
      - backend
    networks:
      - k8s-manager-net
    restart: unless-stopped

volumes:
  mysql-data:

networks:
  k8s-manager-net:
    driver: bridge
```

**.env文件**:
```bash
MYSQL_ROOT_PASSWORD=your_secure_password_here
```

**启动命令**:
```bash
docker-compose up -d
```

#### 8.2.3 Kubernetes部署

可以将整个平台部署到K8S中：

```yaml
# deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: k8s-manager-backend
spec:
  replicas: 2
  selector:
    matchLabels:
      app: k8s-manager-backend
  template:
    metadata:
      labels:
        app: k8s-manager-backend
    spec:
      containers:
      - name: backend
        image: your-registry/k8s-manager-backend:latest
        ports:
        - containerPort: 8080
        env:
        - name: DB_HOST
          value: mysql-service
        - name: REDIS_HOST
          value: redis-service
---
apiVersion: v1
kind: Service
metadata:
  name: k8s-manager-backend
spec:
  selector:
    app: k8s-manager-backend
  ports:
  - port: 8080
    targetPort: 8080
```

### 8.3 配置管理

**config.yaml** (示例):
```yaml
server:
  port: 8080
  mode: release  # debug/release

database:
  host: localhost
  port: 3306
  user: root
  password: password
  dbname: k8s_manager
  max_idle_conns: 10
  max_open_conns: 100

redis:
  host: localhost
  port: 6379
  password: ""
  db: 0

jwt:
  secret: your-jwt-secret-key
  expire: 7200  # 2小时(秒)

log:
  level: info  # debug/info/warn/error
  file: logs/app.log
```

---

## 10. 安全设计

### 10.1 认证安全

#### JWT Token管理
- Token过期时间: 2小时
- Refresh Token: 7天
- Token存储: localStorage（前端）
- Token传递: Authorization Header

#### 密码安全
- 使用bcrypt加密存储
- 密码强度要求: 最少8位，包含大小写字母和数字
- 支持修改密码功能

### 10.2 授权安全

#### RBAC权限模型
- 角色: admin（管理员）、user（普通用户）
- 管理员: 所有权限
- 普通用户: 根据集群权限配置

#### 集群级别权限
- 每个用户可配置对不同集群的权限
- 权限类型: read（查看）、write（修改）、delete（删除）
- 命名空间级别权限控制

### 10.3 数据安全

#### 敏感信息加密
- kubeconfig: AES-256加密存储
- Secret: Base64编码（K8S原生）
- 数据库密码: 配置文件独立管理

#### 审计日志
- 记录所有操作
- 包含用户、时间、操作类型、资源、结果
- 日志不可删除（仅查询）

### 10.4 网络安全

#### HTTPS/WSS
- 生产环境强制HTTPS
- WebSocket使用WSS

#### CORS配置
- 限制允许的Origin
- 开发环境: 允许localhost
- 生产环境: 配置实际域名

#### 请求限流
- 防止API滥用
- 使用Redis实现限流
- 限制: 每个用户每分钟最多100次请求

### 10.5 K8S访问安全

#### 最小权限原则
- 平台使用的ServiceAccount只授予必要权限
- 避免使用cluster-admin

#### 多集群隔离
- 每个集群独立的kubeconfig
- 集群间完全隔离
- 客户端连接失败自动重试

---

## 11. 附录

### 11.1 技术文档链接

- [Go官方文档](https://golang.org/doc/)
- [Gin框架文档](https://gin-gonic.com/docs/)
- [client-go文档](https://github.com/kubernetes/client-go)
- [GORM文档](https://gorm.io/docs/)
- [React官方文档](https://react.dev/)
- [Ant Design文档](https://ant.design/)
- [xterm.js文档](https://xtermjs.org/)
- [ECharts文档](https://echarts.apache.org/)

### 11.2 常用命令

#### Go开发
```bash
# 初始化项目
go mod init k8s-manager-backend

# 安装依赖
go get github.com/gin-gonic/gin
go get k8s.io/client-go@latest
go get gorm.io/gorm
go get gorm.io/driver/mysql

# 运行
go run cmd/server/main.go

# 构建
go build -o bin/server cmd/server/main.go

# 测试
go test ./...
```

#### React开发
```bash
# 创建项目
npx create-react-app k8s-manager-frontend --template typescript

# 安装依赖
npm install antd axios react-router-dom
npm install @ant-design/icons
npm install xterm xterm-addon-fit
npm install echarts echarts-for-react

# 启动开发服务器
npm start

# 构建生产版本
npm run build
```

#### Docker
```bash
# 构建镜像
docker build -t k8s-manager-backend .
docker build -t k8s-manager-frontend .

# 运行容器
docker run -d -p 8080:8080 k8s-manager-backend
docker run -d -p 80:80 k8s-manager-frontend

# Docker Compose
docker-compose up -d
docker-compose down
docker-compose logs -f backend
```

### 11.3 AI提示词模板

#### 后端代码生成
```
我正在开发一个K8S管理平台，使用Go + Gin + client-go。
请帮我实现[功能名称]，要求：
1. [需求1]
2. [需求2]
3. [需求3]

请提供完整的代码，包括：
- Service层方法
- API Handler
- 路由注册
- 错误处理
- 代码注释
```

#### 前端组件生成
```
我正在开发一个K8S管理平台前端，使用React + TypeScript + Ant Design。
请帮我实现[组件名称]，要求：
1. [需求1]
2. [需求2]
3. [需求3]

请提供完整的代码，包括：
- 组件代码
- TypeScript类型定义
- API调用
- 错误处理
```

#### 调试帮助
```
我在实现[功能]时遇到了问题：
错误信息: [错误信息]
我的代码: [代码片段]
环境信息: [Go版本/K8S版本等]

请帮我：
1. 分析错误原因
2. 给出解决方案
3. 提供修复后的代码
```

### 11.4 Git提交规范

```bash
# 功能开发
git commit -m "feat(模块): 添加xxx功能"

# Bug修复
git commit -m "fix(模块): 修复xxx问题"

# 文档更新
git commit -m "docs: 更新README"

# 代码重构
git commit -m "refactor(模块): 重构xxx代码"

# 样式修改
git commit -m "style: 调整xxx样式"

# 性能优化
git commit -m "perf(模块): 优化xxx性能"
```

### 11.5 常见问题

#### Q1: client-go如何从kubeconfig字符串创建clientset？
```go
import (
    "k8s.io/client-go/kubernetes"
    "k8s.io/client-go/tools/clientcmd"
)

config, err := clientcmd.RESTConfigFromKubeConfig([]byte(kubeconfigStr))
if err != nil {
    return nil, err
}

clientset, err := kubernetes.NewForConfig(config)
if err != nil {
    return nil, err
}
```

#### Q2: 如何实现WebSocket日志流？
```go
import (
    "github.com/gorilla/websocket"
    corev1 "k8s.io/api/core/v1"
)

// 升级HTTP到WebSocket
conn, err := upgrader.Upgrade(w, r, nil)

// 获取日志流
req := client.CoreV1().Pods(namespace).GetLogs(podName, &corev1.PodLogOptions{
    Container: container,
    Follow:    true,
})
stream, err := req.Stream(context.Background())

// 转发日志到WebSocket
buf := make([]byte, 1024)
for {
    n, err := stream.Read(buf)
    if err != nil {
        break
    }
    conn.WriteMessage(websocket.TextMessage, buf[:n])
}
```

#### Q3: 前端如何实现WebSocket连接？
```typescript
const ws = new WebSocket('ws://localhost:8080/api/v1/logs/stream');

ws.onopen = () => {
  console.log('WebSocket连接成功');
};

ws.onmessage = (event) => {
  console.log('收到日志:', event.data);
  // 更新UI
};

ws.onerror = (error) => {
  console.error('WebSocket错误:', error);
};

ws.onclose = () => {
  console.log('WebSocket连接关闭');
};
```

### 11.6 性能优化建议

#### 后端优化
1. 使用连接池（数据库、Redis）
2. 客户端缓存（避免重复创建K8S客户端）
3. 异步处理（审计日志、通知等）
4. 分页查询（避免一次加载大量数据）
5. 使用Redis缓存集群信息

#### 前端优化
1. 列表虚拟滚动（大量数据渲染）
2. 图片懒加载
3. 代码分割（React.lazy + Suspense）
4. 防抖和节流（搜索、滚动等）
5. 使用React Query缓存API响应

---

## 结语

这份设计方案涵盖了K8S集群管理平台的完整设计，从架构设计、数据库设计、API设计到前端页面设计、开发计划和部署方案。

**建议的开发顺序**:
1. 先完成MVP功能（集群管理、Deployment、Pod）
2. 逐步添加其他工作负载和功能
3. 最后完善监控、权限等高级功能

**AI辅助开发要点**:
- 每个功能模块都可以让AI帮助生成代码框架
- 遇到问题及时让AI帮助调试
- 善用AI审查代码质量

**预期成果**:
- 一个完整可用的K8S管理平台
- 支持多集群管理
- 友好的Web界面
- 完善的权限控制和审计

祝开发顺利！如有任何问题，随时咨询AI助手。

---

**文档版本**: v1.0  
**最后更新**: 2025-11-10  
**维护者**: Jason
