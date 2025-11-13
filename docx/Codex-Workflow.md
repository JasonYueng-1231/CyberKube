# Codex 开发工作流规范

**项目**: K8S集群管理平台  
**版本**: v1.0  
**创建日期**: 2025-11-11  
**适用环境**: openEuler + Docker + Codex

---

## 目录

1. [工作流概览](#1-工作流概览)
2. [开发环境准备](#2-开发环境准备)
3. [需求接收与分析](#3-需求接收与分析)
4. [开发流程](#4-开发流程)
5. [代码规范](#5-代码规范)
6. [测试流程](#6-测试流程)
7. [提交规范](#7-提交规范)
8. [问题处理](#8-问题处理)
9. [最佳实践](#9-最佳实践)
10. [检查清单](#10-检查清单)

---

## 1. 工作流概览

### 1.1 完整开发流程

```
┌─────────────────────────────────────────────────────────────┐
│                    Codex 开发工作流                      │
└─────────────────────────────────────────────────────────────┘

第1步: 需求接收
  ↓
  用户描述需求 → Codex理解需求 → 确认需求细节
  ↓
  
第2步: 技术设计
  ↓
  分析技术方案 → 确定实现路径 → 评估工作量
  ↓
  
第3步: 环境准备
  ↓
  启动测试环境 → 检查依赖 → 准备数据
  ↓
  
第4步: 编写代码
  ↓
  后端开发 → 前端开发 → 集成联调
  ↓
  
第5步: 自测验证
  ↓
  单元测试 → 功能测试 → 边界测试
  ↓
  
第6步: 代码审查
  ↓
  代码规范 → 性能检查 → 安全检查
  ↓
  
第7步: 提交交付
  ↓
  整理代码 → 编写文档 → 提供给用户
  ↓
  
第8步: 用户验收
  ↓
  用户测试 → 收集反馈 → 问题修复
  ↓
  
完成 ✓
```

### 1.2 时间分配建议

**小功能（<2小时）**:
- 需求分析: 10%
- 技术设计: 10%
- 编写代码: 50%
- 测试验证: 20%
- 文档整理: 10%

**中等功能（2-8小时）**:
- 需求分析: 15%
- 技术设计: 15%
- 编写代码: 40%
- 测试验证: 20%
- 文档整理: 10%

**大功能（>8小时）**:
- 需求分析: 20%
- 技术设计: 20%
- 编写代码: 35%
- 测试验证: 15%
- 文档整理: 10%

---

## 2. 开发环境准备

### 2.1 环境检查清单

**每次开发前必须检查**:

```bash
#!/bin/bash
# check-env.sh - 环境检查脚本

echo "🔍 开始环境检查..."
echo ""

# 1. 检查Docker服务
echo "1️⃣ 检查Docker服务..."
if ! docker info > /dev/null 2>&1; then
    echo "❌ Docker未运行，请执行: sudo systemctl start docker"
    exit 1
fi
echo "✅ Docker服务正常"

# 2. 检查测试环境容器
echo ""
echo "2️⃣ 检查测试容器..."
if ! docker ps | grep -q "k8s-mgr-mysql-test"; then
    echo "⚠️  MySQL测试容器未运行，是否启动? (y/n)"
    read -r answer
    if [ "$answer" = "y" ]; then
        cd ~/k8s-manager
        ./test-env-start.sh
    fi
else
    echo "✅ MySQL测试容器运行中"
fi

if ! docker ps | grep -q "k8s-mgr-redis-test"; then
    echo "⚠️  Redis测试容器未运行"
else
    echo "✅ Redis测试容器运行中"
fi

# 3. 检查项目目录
echo ""
echo "3️⃣ 检查项目目录..."
if [ ! -d ~/k8s-manager ]; then
    echo "❌ 项目目录不存在"
    exit 1
fi
echo "✅ 项目目录存在"

# 4. 检查Go环境
echo ""
echo "4️⃣ 检查Go环境..."
if ! command -v go &> /dev/null; then
    echo "❌ Go未安装"
    exit 1
fi
echo "✅ Go版本: $(go version)"

# 5. 检查Node环境
echo ""
echo "5️⃣ 检查Node环境..."
if ! command -v node &> /dev/null; then
    echo "⚠️  Node.js未安装（前端开发需要）"
else
    echo "✅ Node版本: $(node --version)"
fi

# 6. 检查数据库连接
echo ""
echo "6️⃣ 检查数据库连接..."
if docker exec k8s-mgr-mysql-test mysql -uroot -ptest123456 -e "SELECT 1" > /dev/null 2>&1; then
    echo "✅ MySQL连接正常"
else
    echo "❌ MySQL连接失败"
    exit 1
fi

# 7. 检查磁盘空间
echo ""
echo "7️⃣ 检查磁盘空间..."
disk_usage=$(df -h / | awk 'NR==2 {print $5}' | sed 's/%//')
if [ "$disk_usage" -gt 80 ]; then
    echo "⚠️  磁盘使用率: ${disk_usage}%（建议清理）"
else
    echo "✅ 磁盘使用率: ${disk_usage}%"
fi

echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "✅ 环境检查完成，可以开始开发！"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
```

**使用方式**:
```bash
# 保存脚本
vim ~/k8s-manager/check-env.sh
# 粘贴上面的内容

# 赋予执行权限
chmod +x ~/k8s-manager/check-env.sh

# 每次开发前执行
cd ~/k8s-manager
./check-env.sh
```

### 2.2 工作区准备

```bash
# 1. 进入项目目录
cd ~/k8s-manager

# 2. 启动测试环境（如果未启动）
./test-env-start.sh

# 3. 创建今日工作分支（可选）
cd backend
git checkout -b feature/$(date +%Y%m%d)-集群管理

# 4. 准备工作日志
echo "## $(date +%Y-%m-%d) 开发日志" >> ~/k8s-manager/dev-log.md
echo "" >> ~/k8s-manager/dev-log.md
```

---

## 3. 需求接收与分析

### 3.1 需求确认模板

**Codex在收到需求后，应该使用以下模板确认**:

```markdown
## 需求确认

### 📋 需求描述
[用户的原始需求]

### 🎯 目标
[明确的开发目标]

### 📦 交付物
- [ ] 后端API接口
- [ ] 前端页面/组件
- [ ] 数据库表结构
- [ ] 接口文档
- [ ] 测试用例

### 🔧 技术方案
[简要的技术实现方案]

### ⏱️ 预计工作量
- 后端开发: X小时
- 前端开发: X小时
- 测试验证: X小时
- 总计: X小时

### ❓ 待确认问题
1. [问题1]
2. [问题2]

### ✅ 确认信息
请确认以上理解是否正确？是否有补充或修改？
```

### 3.2 需求分析示例

**示例1: 简单功能**

```
用户需求: "实现用户登录功能"

Codex确认:
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
📋 需求描述
实现用户登录功能，用户使用用户名和密码登录系统

🎯 目标
1. 用户可以通过用户名和密码登录
2. 登录成功返回JWT token
3. token有效期2小时

📦 交付物
- [x] 后端登录API
- [x] JWT token生成和验证
- [x] 前端登录页面
- [x] 登录状态管理
- [x] API文档

🔧 技术方案
后端:
- 使用bcrypt验证密码
- 使用JWT生成token
- 实现登录API: POST /api/v1/auth/login

前端:
- 使用Ant Design Form组件
- 登录成功后保存token到localStorage
- 自动跳转到首页

⏱️ 预计工作量
- 后端开发: 1小时
- 前端开发: 1.5小时
- 测试验证: 0.5小时
- 总计: 3小时

❓ 待确认问题
1. 是否需要记住密码功能？
2. 是否需要验证码？
3. 登录失败是否有次数限制？

✅ 请确认是否可以开始开发？
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
```

**示例2: 复杂功能**

```
用户需求: "实现Deployment管理功能"

Codex确认:
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
📋 需求描述
实现Kubernetes Deployment的完整管理功能

🎯 目标
1. 查看Deployment列表（支持分页、筛选）
2. 查看Deployment详情
3. 创建新的Deployment
4. 更新Deployment配置
5. 删除Deployment
6. 扩缩容操作
7. 重启Deployment

📦 交付物
- [x] 后端API（7个接口）
- [x] client-go封装
- [x] 前端列表页面
- [x] 前端详情页面
- [x] 前端创建/编辑表单
- [x] API文档
- [x] 测试用例

🔧 技术方案
后端:
1. 封装client-go操作Deployment的方法
2. 实现7个RESTful API
3. 支持YAML和表单两种创建方式
4. 添加操作审计日志

前端:
1. 列表页: Ant Design Table + 搜索/筛选
2. 详情页: Tab页（基本信息/Pods/Events/YAML）
3. 创建页: 表单模式 + YAML编辑器切换
4. 操作确认: Modal二次确认

数据库:
- 无需新建表（K8S数据存储在etcd）
- 审计日志表已存在

⏱️ 预计工作量
- 后端开发: 4小时
- 前端开发: 6小时
- 测试验证: 2小时
- 总计: 12小时

建议分3次完成:
第1次: 列表查询 + 详情查看（4h）
第2次: 创建 + 更新 + 删除（4h）
第3次: 扩缩容 + 重启 + 优化（4h）

❓ 待确认问题
1. 是否支持批量操作？
2. YAML编辑器使用Monaco还是简单的textarea？
3. 是否需要Deployment模板功能？

✅ 请确认方案，我将开始第一阶段开发
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
```

---

## 4. 开发流程

### 4.1 后端开发标准流程

#### Step 1: 阅读技能文档（如果需要）

```bash
# 如果涉及特定技术，先阅读技能文档
# 例如: 操作Excel文件
cat /mnt/skills/public/xlsx/SKILL.md

# 或者: 创建Word文档
cat /mnt/skills/public/docx/SKILL.md
```

#### Step 2: 创建目录结构

```bash
# 示例: 创建集群管理模块
mkdir -p ~/k8s-manager/backend/internal/api/v1/cluster
mkdir -p ~/k8s-manager/backend/internal/service
mkdir -p ~/k8s-manager/backend/internal/k8s
```

#### Step 3: 编写代码（自底向上）

**顺序**: 数据模型 → K8S操作 → Service层 → API层 → 路由

**1. 数据模型** (`internal/model/cluster.go`):
```go
package model

import "time"

type Cluster struct {
    ID          uint      `gorm:"primaryKey" json:"id"`
    Name        string    `gorm:"uniqueIndex;size:100;not null" json:"name"`
    Kubeconfig  string    `gorm:"type:text;not null" json:"-"` // 不返回给前端
    APIServer   string    `gorm:"size:255" json:"api_server"`
    Version     string    `gorm:"size:50" json:"version"`
    Status      int       `gorm:"default:1" json:"status"` // 1-正常 0-异常
    Description string    `gorm:"size:500" json:"description"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}
```

**2. K8S操作封装** (`internal/k8s/client_manager.go`):
```go
package k8s

import (
    "sync"
    "k8s.io/client-go/kubernetes"
    "k8s.io/client-go/tools/clientcmd"
)

type ClusterManager struct {
    clients sync.Map // map[string]*kubernetes.Clientset
}

var Manager = &ClusterManager{}

func (cm *ClusterManager) AddCluster(name, kubeconfig string) error {
    config, err := clientcmd.RESTConfigFromKubeConfig([]byte(kubeconfig))
    if err != nil {
        return err
    }
    
    clientset, err := kubernetes.NewForConfig(config)
    if err != nil {
        return err
    }
    
    cm.clients.Store(name, clientset)
    return nil
}

func (cm *ClusterManager) GetClient(name string) (*kubernetes.Clientset, error) {
    client, ok := cm.clients.Load(name)
    if !ok {
        return nil, errors.New("cluster not found")
    }
    return client.(*kubernetes.Clientset), nil
}
```

**3. Service层** (`internal/service/cluster_service.go`):
```go
package service

import (
    "k8s-manager-backend/internal/model"
    "k8s-manager-backend/internal/k8s"
    "k8s-manager-backend/internal/database"
)

type ClusterService struct{}

func NewClusterService() *ClusterService {
    return &ClusterService{}
}

// 创建集群
func (s *ClusterService) Create(cluster *model.Cluster) error {
    // 1. 测试连接
    if err := k8s.Manager.AddCluster(cluster.Name, cluster.Kubeconfig); err != nil {
        return err
    }
    
    // 2. 加密kubeconfig
    encryptedConfig, err := encrypt.AESEncrypt(cluster.Kubeconfig)
    if err != nil {
        return err
    }
    cluster.Kubeconfig = encryptedConfig
    
    // 3. 保存到数据库
    if err := database.DB.Create(cluster).Error; err != nil {
        return err
    }
    
    return nil
}

// 获取列表
func (s *ClusterService) List(page, pageSize int) ([]model.Cluster, int64, error) {
    var clusters []model.Cluster
    var total int64
    
    offset := (page - 1) * pageSize
    
    if err := database.DB.Model(&model.Cluster{}).Count(&total).Error; err != nil {
        return nil, 0, err
    }
    
    if err := database.DB.Offset(offset).Limit(pageSize).Find(&clusters).Error; err != nil {
        return nil, 0, err
    }
    
    return clusters, total, nil
}
```

**4. API层** (`internal/api/v1/cluster/cluster.go`):
```go
package cluster

import (
    "net/http"
    "strconv"
    "github.com/gin-gonic/gin"
    "k8s-manager-backend/internal/service"
    "k8s-manager-backend/pkg/response"
)

type ClusterHandler struct {
    service *service.ClusterService
}

func NewClusterHandler() *ClusterHandler {
    return &ClusterHandler{
        service: service.NewClusterService(),
    }
}

// 创建集群
// @Summary 创建集群
// @Tags 集群管理
// @Accept json
// @Produce json
// @Param cluster body model.Cluster true "集群信息"
// @Success 200 {object} response.Response
// @Router /api/v1/clusters [post]
func (h *ClusterHandler) Create(c *gin.Context) {
    var cluster model.Cluster
    if err := c.ShouldBindJSON(&cluster); err != nil {
        response.Error(c, 40001, "参数错误: "+err.Error())
        return
    }
    
    if err := h.service.Create(&cluster); err != nil {
        response.Error(c, 50001, "创建失败: "+err.Error())
        return
    }
    
    response.Success(c, cluster)
}

// 获取列表
func (h *ClusterHandler) List(c *gin.Context) {
    page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
    pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
    
    clusters, total, err := h.service.List(page, pageSize)
    if err != nil {
        response.Error(c, 50001, "查询失败: "+err.Error())
        return
    }
    
    response.Success(c, gin.H{
        "total":     total,
        "page":      page,
        "page_size": pageSize,
        "items":     clusters,
    })
}
```

**5. 路由注册** (`internal/api/router.go`):
```go
package api

import (
    "github.com/gin-gonic/gin"
    "k8s-manager-backend/internal/api/v1/cluster"
    "k8s-manager-backend/internal/middleware"
)

func SetupRouter() *gin.Engine {
    r := gin.Default()
    
    // 中间件
    r.Use(middleware.CORS())
    r.Use(middleware.Logger())
    
    // API v1
    v1 := r.Group("/api/v1")
    {
        // 集群管理
        clusterHandler := cluster.NewClusterHandler()
        clusterGroup := v1.Group("/clusters")
        clusterGroup.Use(middleware.Auth()) // 需要认证
        {
            clusterGroup.POST("", clusterHandler.Create)
            clusterGroup.GET("", clusterHandler.List)
            clusterGroup.GET("/:name", clusterHandler.Get)
            clusterGroup.PUT("/:name", clusterHandler.Update)
            clusterGroup.DELETE("/:name", clusterHandler.Delete)
            clusterGroup.POST("/:name/test", clusterHandler.Test)
        }
    }
    
    return r
}
```

#### Step 4: 编译测试

```bash
# 1. 进入后端目录
cd ~/k8s-manager/backend

# 2. 下载依赖（如果有新依赖）
go mod tidy

# 3. 编译检查
go build -o bin/test cmd/server/main.go

# 4. 运行测试
go run cmd/server/main.go

# 5. 在另一个终端测试API
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin123"}'
```

### 4.2 前端开发标准流程

#### Step 1: 创建组件/页面

```bash
# 创建页面目录
mkdir -p ~/k8s-manager/frontend/src/pages/Cluster

# 创建组件文件
touch ~/k8s-manager/frontend/src/pages/Cluster/ClusterList.tsx
touch ~/k8s-manager/frontend/src/pages/Cluster/ClusterList.css
```

#### Step 2: 编写代码（自底向上）

**1. API服务** (`src/services/cluster.ts`):
```typescript
import request from '@/utils/request';

export interface Cluster {
  id: number;
  name: string;
  api_server: string;
  version: string;
  status: number;
  description: string;
  created_at: string;
}

export interface ClusterListResponse {
  total: number;
  page: number;
  page_size: number;
  items: Cluster[];
}

// 获取集群列表
export async function getClusterList(params: {
  page?: number;
  page_size?: number;
}): Promise<ClusterListResponse> {
  return request.get('/clusters', { params });
}

// 创建集群
export async function createCluster(data: {
  name: string;
  kubeconfig: string;
  description?: string;
}): Promise<Cluster> {
  return request.post('/clusters', data);
}

// 删除集群
export async function deleteCluster(name: string): Promise<void> {
  return request.delete(`/clusters/${name}`);
}
```

**2. 页面组件** (`src/pages/Cluster/ClusterList.tsx`):
```typescript
import React, { useState, useEffect } from 'react';
import { Table, Button, Space, Modal, message } from 'antd';
import { PlusOutlined, DeleteOutlined, EyeOutlined } from '@ant-design/icons';
import { getClusterList, deleteCluster, type Cluster } from '@/services/cluster';
import './ClusterList.css';

export const ClusterList: React.FC = () => {
  const [loading, setLoading] = useState(false);
  const [clusters, setClusters] = useState<Cluster[]>([]);
  const [total, setTotal] = useState(0);
  const [page, setPage] = useState(1);
  const [pageSize, setPageSize] = useState(20);

  // 加载数据
  const loadData = async () => {
    setLoading(true);
    try {
      const res = await getClusterList({ page, page_size: pageSize });
      setClusters(res.items);
      setTotal(res.total);
    } catch (error) {
      message.error('加载失败');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    loadData();
  }, [page, pageSize]);

  // 删除集群
  const handleDelete = (record: Cluster) => {
    Modal.confirm({
      title: '确认删除',
      content: `确定要删除集群 "${record.name}" 吗？`,
      onOk: async () => {
        try {
          await deleteCluster(record.name);
          message.success('删除成功');
          loadData();
        } catch (error) {
          message.error('删除失败');
        }
      },
    });
  };

  const columns = [
    {
      title: '集群名称',
      dataIndex: 'name',
      key: 'name',
    },
    {
      title: 'API Server',
      dataIndex: 'api_server',
      key: 'api_server',
    },
    {
      title: '版本',
      dataIndex: 'version',
      key: 'version',
    },
    {
      title: '状态',
      dataIndex: 'status',
      key: 'status',
      render: (status: number) => (
        <span className={`status-badge ${status === 1 ? 'success' : 'error'}`}>
          {status === 1 ? '正常' : '异常'}
        </span>
      ),
    },
    {
      title: '创建时间',
      dataIndex: 'created_at',
      key: 'created_at',
    },
    {
      title: '操作',
      key: 'action',
      render: (_: any, record: Cluster) => (
        <Space>
          <Button 
            type="link" 
            icon={<EyeOutlined />}
            onClick={() => {/* 查看详情 */}}
          >
            详情
          </Button>
          <Button 
            type="link" 
            danger
            icon={<DeleteOutlined />}
            onClick={() => handleDelete(record)}
          >
            删除
          </Button>
        </Space>
      ),
    },
  ];

  return (
    <div className="cluster-list cyber-page">
      <div className="page-header">
        <h2 className="page-title">集群管理</h2>
        <Button 
          type="primary" 
          icon={<PlusOutlined />}
          onClick={() => {/* 打开创建弹窗 */}}
        >
          添加集群
        </Button>
      </div>

      <div className="page-content">
        <Table
          loading={loading}
          columns={columns}
          dataSource={clusters}
          rowKey="id"
          pagination={{
            total,
            current: page,
            pageSize,
            onChange: (p, ps) => {
              setPage(p);
              setPageSize(ps || 20);
            },
          }}
        />
      </div>
    </div>
  );
};
```

**3. 样式文件** (`src/pages/Cluster/ClusterList.css`):
```css
.cluster-list {
  padding: 24px;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
}

.page-title {
  font-size: 24px;
  font-weight: 600;
  color: var(--text-primary);
  margin: 0;
}

.page-content {
  background: var(--bg-card);
  border-radius: 12px;
  padding: 24px;
  border: 1px solid var(--border-glow);
}

.status-badge {
  padding: 4px 12px;
  border-radius: 12px;
  font-size: 12px;
  font-weight: 600;
}

.status-badge.success {
  background: rgba(57, 255, 20, 0.1);
  color: var(--neon-green);
  border: 1px solid var(--neon-green);
}

.status-badge.error {
  background: rgba(255, 0, 110, 0.1);
  color: var(--neon-red);
  border: 1px solid var(--neon-red);
}
```

#### Step 3: 测试运行

```bash
# 1. 进入前端目录
cd ~/k8s-manager/frontend

# 2. 安装依赖（如果有新依赖）
npm install

# 3. 启动开发服务器
npm start

# 4. 在浏览器打开
# http://localhost:3000
```

### 4.3 集成联调流程

```bash
# 终端1: 启动后端
cd ~/k8s-manager/backend
go run cmd/server/main.go

# 终端2: 启动前端
cd ~/k8s-manager/frontend
npm start

# 终端3: 查看测试环境日志（如果需要）
cd ~/k8s-manager
./test-env-logs.sh
```

**联调检查项**:
- [ ] API请求是否成功
- [ ] 数据格式是否正确
- [ ] 错误提示是否友好
- [ ] Loading状态是否正常
- [ ] 界面刷新是否及时

---

## 5. 代码规范

### 5.1 Go代码规范

#### 命名规范

```go
// ✅ 好的命名
type UserService struct {}
func (s *UserService) GetUserByID(id uint) (*User, error) {}
var ErrUserNotFound = errors.New("user not found")

// ❌ 不好的命名
type userservice struct {}
func (s *userservice) get_user(i uint) (*User, error) {}
var err1 = errors.New("not found")
```

#### 注释规范

```go
// ✅ 完整的注释
// UserService 用户服务
// 提供用户相关的业务逻辑处理
type UserService struct {
    repo *UserRepository
}

// GetUserByID 根据ID获取用户信息
// @param id 用户ID
// @return 用户对象和错误信息
func (s *UserService) GetUserByID(id uint) (*User, error) {
    // 参数验证
    if id == 0 {
        return nil, ErrInvalidUserID
    }
    
    // 查询用户
    user, err := s.repo.FindByID(id)
    if err != nil {
        return nil, err
    }
    
    return user, nil
}

// ❌ 缺少注释
type UserService struct {
    repo *UserRepository
}

func (s *UserService) GetUserByID(id uint) (*User, error) {
    if id == 0 {
        return nil, ErrInvalidUserID
    }
    user, err := s.repo.FindByID(id)
    if err != nil {
        return nil, err
    }
    return user, nil
}
```

#### 错误处理规范

```go
// ✅ 完整的错误处理
func (s *ClusterService) Create(cluster *Cluster) error {
    // 1. 参数验证
    if cluster.Name == "" {
        return errors.New("集群名称不能为空")
    }
    
    // 2. 业务逻辑
    if err := s.validateKubeconfig(cluster.Kubeconfig); err != nil {
        return fmt.Errorf("kubeconfig验证失败: %w", err)
    }
    
    // 3. 数据持久化
    if err := s.repo.Create(cluster); err != nil {
        return fmt.Errorf("保存集群失败: %w", err)
    }
    
    return nil
}

// ❌ 不完整的错误处理
func (s *ClusterService) Create(cluster *Cluster) error {
    s.validateKubeconfig(cluster.Kubeconfig)
    s.repo.Create(cluster)
    return nil
}
```

### 5.2 TypeScript代码规范

#### 类型定义

```typescript
// ✅ 完整的类型定义
export interface Cluster {
  id: number;
  name: string;
  api_server: string;
  version: string;
  status: number;
  description?: string;
  created_at: string;
  updated_at: string;
}

export interface ClusterListParams {
  page?: number;
  page_size?: number;
  keyword?: string;
  status?: number;
}

// ❌ 使用any
export interface Cluster {
  id: any;
  name: any;
  [key: string]: any;
}
```

#### 组件规范

```typescript
// ✅ 完整的组件定义
import React, { useState, useEffect } from 'react';

interface ClusterCardProps {
  cluster: Cluster;
  onDelete?: (cluster: Cluster) => void;
  onUpdate?: (cluster: Cluster) => void;
}

/**
 * 集群卡片组件
 * 显示集群的基本信息和操作按钮
 */
export const ClusterCard: React.FC<ClusterCardProps> = ({
  cluster,
  onDelete,
  onUpdate,
}) => {
  const [loading, setLoading] = useState(false);

  const handleDelete = async () => {
    setLoading(true);
    try {
      await onDelete?.(cluster);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="cluster-card">
      {/* 组件内容 */}
    </div>
  );
};

// ❌ 缺少类型和注释
export const ClusterCard = ({ cluster, onDelete }) => {
  return <div>{cluster.name}</div>;
};
```

### 5.3 CSS/样式规范

```css
/* ✅ 好的CSS组织 */
/* ==================== 集群卡片 ==================== */
.cluster-card {
  /* 布局 */
  position: relative;
  padding: 24px;
  
  /* 外观 */
  background: var(--bg-card);
  border: 1px solid var(--border-glow);
  border-radius: 12px;
  
  /* 动画 */
  transition: all 0.3s ease;
}

.cluster-card:hover {
  transform: translateY(-4px);
  box-shadow: var(--shadow-cyan-md);
}

.cluster-card__header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
}

.cluster-card__title {
  font-size: 18px;
  font-weight: 600;
  color: var(--neon-cyan);
}

/* ❌ 混乱的CSS */
.cluster-card {
  padding: 24px;
  background: var(--bg-card);
  transition: all 0.3s ease;
  border: 1px solid var(--border-glow);
  position: relative;
  border-radius: 12px;
}
.cluster-card:hover {transform: translateY(-4px);box-shadow: var(--shadow-cyan-md);}
.cluster-card__header {display: flex;justify-content: space-between;}
```

---

## 6. 测试流程

### 6.1 后端测试清单

#### 单元测试（可选）

```go
// cluster_service_test.go
package service

import (
    "testing"
    "github.com/stretchr/testify/assert"
)

func TestClusterService_Create(t *testing.T) {
    service := NewClusterService()
    
    cluster := &Cluster{
        Name:       "test-cluster",
        Kubeconfig: "...",
    }
    
    err := service.Create(cluster)
    assert.NoError(t, err)
    assert.NotZero(t, cluster.ID)
}
```

#### API测试（必须）

```bash
#!/bin/bash
# test-api.sh - API测试脚本

BASE_URL="http://localhost:8080/api/v1"
TOKEN=""

echo "🧪 开始API测试..."
echo ""

# 1. 测试登录
echo "1️⃣ 测试登录..."
LOGIN_RESP=$(curl -s -X POST ${BASE_URL}/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin123"}')

TOKEN=$(echo $LOGIN_RESP | jq -r '.data.token')

if [ "$TOKEN" = "null" ] || [ -z "$TOKEN" ]; then
    echo "❌ 登录失败"
    exit 1
fi
echo "✅ 登录成功, Token: ${TOKEN:0:20}..."

# 2. 测试集群列表
echo ""
echo "2️⃣ 测试集群列表..."
LIST_RESP=$(curl -s -X GET ${BASE_URL}/clusters \
  -H "Authorization: Bearer ${TOKEN}")

echo $LIST_RESP | jq '.'

if echo $LIST_RESP | jq -e '.code == 0' > /dev/null; then
    echo "✅ 集群列表获取成功"
else
    echo "❌ 集群列表获取失败"
    exit 1
fi

# 3. 测试创建集群
echo ""
echo "3️⃣ 测试创建集群..."
CREATE_RESP=$(curl -s -X POST ${BASE_URL}/clusters \
  -H "Authorization: Bearer ${TOKEN}" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "test-cluster-'$(date +%s)'",
    "kubeconfig": "apiVersion: v1...",
    "description": "测试集群"
  }')

echo $CREATE_RESP | jq '.'

if echo $CREATE_RESP | jq -e '.code == 0' > /dev/null; then
    echo "✅ 创建集群成功"
else
    echo "❌ 创建集群失败"
fi

echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "✅ API测试完成"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
```

### 6.2 前端测试清单

#### 功能测试

- [ ] 页面能正常加载
- [ ] 数据能正确显示
- [ ] 按钮点击有响应
- [ ] 表单能正确提交
- [ ] 错误能正确提示
- [ ] Loading状态正常
- [ ] 分页功能正常

#### 交互测试

- [ ] 悬停效果正常
- [ ] 点击动画流畅
- [ ] 模态框能打开关闭
- [ ] 表单验证正确
- [ ] 删除确认正常

#### 兼容性测试

- [ ] Chrome浏览器正常
- [ ] Firefox浏览器正常
- [ ] 不同分辨率正常
- [ ] 移动端响应式正常（如果需要）

### 6.3 集成测试

```bash
#!/bin/bash
# integration-test.sh - 集成测试脚本

echo "🧪 开始集成测试..."
echo ""

# 1. 检查后端
echo "1️⃣ 检查后端服务..."
if ! curl -s http://localhost:8080/health > /dev/null; then
    echo "❌ 后端服务未启动"
    exit 1
fi
echo "✅ 后端服务正常"

# 2. 检查前端
echo ""
echo "2️⃣ 检查前端服务..."
if ! curl -s http://localhost:3000 > /dev/null; then
    echo "❌ 前端服务未启动"
    exit 1
fi
echo "✅ 前端服务正常"

# 3. 检查数据库
echo ""
echo "3️⃣ 检查数据库连接..."
if ! docker exec k8s-mgr-mysql-test mysql -uroot -ptest123456 -e "SELECT 1" > /dev/null 2>&1; then
    echo "❌ 数据库连接失败"
    exit 1
fi
echo "✅ 数据库连接正常"

# 4. 执行端到端测试
echo ""
echo "4️⃣ 执行端到端测试..."

# 4.1 登录
TOKEN=$(curl -s -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin123"}' \
  | jq -r '.data.token')

if [ "$TOKEN" = "null" ]; then
    echo "❌ 登录失败"
    exit 1
fi
echo "  ✓ 登录成功"

# 4.2 获取集群列表
CLUSTERS=$(curl -s -X GET http://localhost:8080/api/v1/clusters \
  -H "Authorization: Bearer ${TOKEN}" \
  | jq -r '.data.total')

echo "  ✓ 获取集群列表成功 (共 ${CLUSTERS} 个集群)"

# 4.3 创建测试数据
echo "  ✓ 创建测试数据..."

# 4.4 清理测试数据
echo "  ✓ 清理测试数据..."

echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "✅ 集成测试通过"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
```

---

## 7. 提交规范

### 7.1 Git提交规范

#### Commit Message格式

```
<type>(<scope>): <subject>

<body>

<footer>
```

**Type类型**:
- `feat`: 新功能
- `fix`: Bug修复
- `docs`: 文档更新
- `style`: 代码格式（不影响代码运行）
- `refactor`: 重构
- `perf`: 性能优化
- `test`: 测试相关
- `chore`: 构建过程或辅助工具的变动

**Scope范围**:
- `cluster`: 集群管理
- `deployment`: Deployment管理
- `pod`: Pod管理
- `service`: Service管理
- `auth`: 认证授权
- `ui`: 前端界面
- `api`: API接口

**示例**:

```bash
# 好的commit message
git commit -m "feat(cluster): 添加集群列表查询功能

- 实现集群列表API
- 支持分页和筛选
- 添加集群状态显示

Closes #123"

git commit -m "fix(deployment): 修复扩缩容副本数验证问题

副本数应该大于0，之前允许设置为0导致错误

Fixes #456"

git commit -m "docs: 更新API文档

添加集群管理相关接口文档"

# 不好的commit message
git commit -m "update"
git commit -m "fix bug"
git commit -m "add feature"
```

### 7.2 代码提交检查清单

**提交前必须检查**:

```bash
#!/bin/bash
# pre-commit-check.sh - 提交前检查

echo "📋 提交前检查..."
echo ""

# 1. 代码格式检查
echo "1️⃣ 检查代码格式..."
if command -v gofmt &> /dev/null; then
    UNFORMATTED=$(gofmt -l .)
    if [ -n "$UNFORMATTED" ]; then
        echo "❌ 以下文件需要格式化:"
        echo "$UNFORMATTED"
        exit 1
    fi
    echo "✅ 代码格式正确"
fi

# 2. 语法检查
echo ""
echo "2️⃣ 检查语法错误..."
if ! go build ./... > /dev/null 2>&1; then
    echo "❌ 编译失败，请修复错误"
    exit 1
fi
echo "✅ 编译通过"

# 3. 测试检查
echo ""
echo "3️⃣ 运行测试..."
if ! go test ./... > /dev/null 2>&1; then
    echo "❌ 测试失败"
    exit 1
fi
echo "✅ 测试通过"

# 4. 检查TODO
echo ""
echo "4️⃣ 检查TODO..."
TODO_COUNT=$(git diff --cached | grep -c "TODO")
if [ "$TODO_COUNT" -gt 0 ]; then
    echo "⚠️  发现 ${TODO_COUNT} 个TODO，请确认是否需要处理"
fi

echo ""
echo "✅ 所有检查通过，可以提交"
```

---

## 8. 问题处理

### 8.1 常见问题诊断

#### 问题1: 后端编译失败

**症状**:
```
go build: cannot find module providing package xxx
```

**诊断**:
```bash
# 1. 检查依赖
go mod tidy

# 2. 下载依赖
go mod download

# 3. 清理缓存
go clean -modcache
```

#### 问题2: 数据库连接失败

**症状**:
```
Error: dial tcp 127.0.0.1:13306: connect: connection refused
```

**诊断**:
```bash
# 1. 检查容器是否运行
docker ps | grep mysql

# 2. 检查端口
netstat -tlnp | grep 13306

# 3. 重启容器
docker restart k8s-mgr-mysql-test

# 4. 测试连接
docker exec -it k8s-mgr-mysql-test mysql -uroot -ptest123456
```

#### 问题3: K8S API调用失败

**症状**:
```
Error: Unauthorized
Error: the server has asked for the client to provide credentials
```

**诊断**:
```bash
# 1. 检查kubeconfig
kubectl --kubeconfig=/path/to/config get nodes

# 2. 检查权限
kubectl --kubeconfig=/path/to/config auth can-i get pods

# 3. 检查证书
openssl x509 -in /path/to/cert -text -noout
```

#### 问题4: 前端无法访问后端

**症状**:
```
CORS policy: No 'Access-Control-Allow-Origin' header
```

**诊断**:
```bash
# 1. 检查后端CORS配置
# 查看 internal/middleware/cors.go

# 2. 检查前端代理配置
# 查看 vite.config.ts 或 package.json

# 3. 使用curl测试
curl -H "Origin: http://localhost:3000" \
  -H "Access-Control-Request-Method: POST" \
  -X OPTIONS \
  http://localhost:8080/api/v1/clusters
```

### 8.2 调试技巧

#### 后端调试

```go
// 使用log打印
log.Printf("Debug: cluster name = %s", cluster.Name)

// 使用fmt打印
fmt.Printf("Debug: %+v\n", cluster)

// 使用断点调试（VSCode）
// 在 .vscode/launch.json 中配置
{
  "version": "0.2.0",
  "configurations": [
    {
      "name": "Launch",
      "type": "go",
      "request": "launch",
      "mode": "auto",
      "program": "${workspaceFolder}/cmd/server",
      "env": {},
      "args": []
    }
  ]
}
```

#### 前端调试

```typescript
// 使用console.log
console.log('Debug:', cluster);

// 使用debugger断点
function handleDelete(cluster: Cluster) {
  debugger; // 浏览器会在这里暂停
  // ...
}

// 使用React DevTools
// 安装浏览器插件: React Developer Tools
```

#### 网络调试

```bash
# 查看网络请求
curl -v http://localhost:8080/api/v1/clusters

# 使用浏览器开发者工具
# F12 -> Network -> 查看请求详情

# 使用Postman
# 创建Collection -> 保存常用请求
```

---

## 9. 最佳实践

### 9.1 Codex工作最佳实践

#### 实践1: 分步骤开发

```
❌ 不好的方式:
用户: "实现整个Deployment管理功能"
Codex: [一次性写完所有代码]

✅ 好的方式:
用户: "实现Deployment管理功能"
Codex: "我建议分3步完成:
  第1步: 实现列表查询和详情查看
  第2步: 实现创建和更新
  第3步: 实现删除、扩缩容和重启
  
  现在开始第1步，先实现列表查询？"
```

#### 实践2: 先验证后实现

```
❌ 不好的方式:
Codex: [直接写代码]

✅ 好的方式:
Codex: "开始之前，让我先验证一下:
  1. 检查环境: ✓ Docker运行中
  2. 检查依赖: ✓ client-go已安装
  3. 测试连接: ✓ 数据库连接正常
  
  环境正常，开始编写代码..."
```

#### 实践3: 边写边测

```
❌ 不好的方式:
Codex: [写完所有代码] "代码已完成"

✅ 好的方式:
Codex: [写一部分代码]
  "我已经完成了数据模型和Service层，
   让我先测试一下编译是否通过..."
   
   [执行 go build]
   "✓ 编译通过，继续实现API层..."
   
   [写API代码]
   "API层完成，测试一下..."
   
   [执行curl测试]
   "✓ API测试通过，功能完成"
```

#### 实践4: 主动发现问题

```
Codex发现问题:
"⚠️ 注意: 我发现了一个潜在问题:
  
  当前代码在高并发时可能存在竞态条件，
  因为ClusterManager没有加锁保护。
  
  建议修改为:
  [提供改进方案]
  
  是否需要修复？"
```

### 9.2 代码质量最佳实践

#### 实践1: 单一职责

```go
// ✅ 好的设计 - 职责单一
type ClusterService struct {
    repo *ClusterRepository
}

func (s *ClusterService) Create(cluster *Cluster) error {
    // 只负责业务逻辑
}

type ClusterRepository struct {
    db *gorm.DB
}

func (r *ClusterRepository) Create(cluster *Cluster) error {
    // 只负责数据访问
}

// ❌ 不好的设计 - 职责混乱
type ClusterService struct {
    db *gorm.DB
}

func (s *ClusterService) Create(cluster *Cluster) error {
    // 既有业务逻辑又有数据访问
    return s.db.Create(cluster).Error
}
```

#### 实践2: 依赖注入

```go
// ✅ 好的设计 - 依赖注入
type ClusterHandler struct {
    service *ClusterService
}

func NewClusterHandler(service *ClusterService) *ClusterHandler {
    return &ClusterHandler{
        service: service,
    }
}

// ❌ 不好的设计 - 硬编码依赖
type ClusterHandler struct {}

func (h *ClusterHandler) Create(c *gin.Context) {
    service := NewClusterService() // 每次都创建新实例
    // ...
}
```

#### 实践3: 错误处理

```go
// ✅ 好的错误处理
func (s *ClusterService) GetByName(name string) (*Cluster, error) {
    cluster, err := s.repo.FindByName(name)
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, ErrClusterNotFound
        }
        return nil, fmt.Errorf("查询集群失败: %w", err)
    }
    return cluster, nil
}

// ❌ 不好的错误处理
func (s *ClusterService) GetByName(name string) (*Cluster, error) {
    cluster, err := s.repo.FindByName(name)
    if err != nil {
        return nil, err // 丢失了错误上下文
    }
    return cluster, nil
}
```

### 9.3 性能优化最佳实践

#### 实践1: 数据库查询优化

```go
// ✅ 好的查询 - 只查询需要的字段
func (r *ClusterRepository) List() ([]Cluster, error) {
    var clusters []Cluster
    err := r.db.Select("id, name, api_server, version, status, created_at").
        Find(&clusters).Error
    return clusters, err
}

// ❌ 不好的查询 - 查询所有字段
func (r *ClusterRepository) List() ([]Cluster, error) {
    var clusters []Cluster
    err := r.db.Find(&clusters).Error // kubeconfig也会被查询
    return clusters, err
}
```

#### 实践2: 缓存使用

```go
// ✅ 使用缓存
func (s *ClusterService) GetByName(name string) (*Cluster, error) {
    // 1. 先查缓存
    if cached := cache.Get("cluster:" + name); cached != nil {
        return cached.(*Cluster), nil
    }
    
    // 2. 查数据库
    cluster, err := s.repo.FindByName(name)
    if err != nil {
        return nil, err
    }
    
    // 3. 写入缓存
    cache.Set("cluster:"+name, cluster, 5*time.Minute)
    
    return cluster, nil
}
```

#### 实践3: 并发控制

```go
// ✅ 使用sync.Map保证并发安全
type ClusterManager struct {
    clients sync.Map
    mu      sync.RWMutex
}

// ❌ 不安全的并发
type ClusterManager struct {
    clients map[string]*kubernetes.Clientset // 并发读写会panic
}
```

---

## 10. 检查清单

### 10.1 开发完成检查清单

**每个功能开发完成后，Codex应该检查**:

```markdown
## 功能完成检查清单

### 代码质量
- [ ] 代码已格式化（gofmt/prettier）
- [ ] 没有语法错误
- [ ] 没有明显的逻辑错误
- [ ] 变量命名清晰
- [ ] 函数职责单一
- [ ] 有必要的注释

### 功能完整性
- [ ] 后端API已实现
- [ ] 前端页面已实现
- [ ] 前后端已联调
- [ ] 错误处理完善
- [ ] 边界情况已考虑

### 测试验证
- [ ] 编译通过
- [ ] 单元测试通过（如果有）
- [ ] API测试通过
- [ ] 前端功能测试通过
- [ ] 集成测试通过

### 文档完善
- [ ] API文档已更新
- [ ] 代码注释完整
- [ ] 使用说明清晰

### 性能安全
- [ ] 没有SQL注入风险
- [ ] 没有XSS风险
- [ ] 敏感信息已加密
- [ ] 查询已优化
- [ ] 没有内存泄漏

### 用户体验
- [ ] 加载状态显示
- [ ] 错误提示友好
- [ ] 操作反馈及时
- [ ] 界面响应流畅

### 代码提交
- [ ] Commit message规范
- [ ] 代码已推送
- [ ] 文档已同步

### 交付用户
- [ ] 代码已整理
- [ ] 文档已编写
- [ ] 部署说明清晰
```

### 10.2 日常开发检查清单

**每天开始开发前**:
- [ ] 执行环境检查脚本
- [ ] 拉取最新代码
- [ ] 启动测试环境
- [ ] 查看待办事项

**每天结束开发后**:
- [ ] 提交代码
- [ ] 停止测试环境（可选）
- [ ] 记录开发日志
- [ ] 更新待办事项

---

## 附录

### A. 快捷命令

```bash
# ~/.bashrc 或 ~/.zshrc
alias k8s-dev-start="cd ~/k8s-manager && ./test-env-start.sh"
alias k8s-dev-stop="cd ~/k8s-manager && ./test-env-stop.sh"
alias k8s-dev-reset="cd ~/k8s-manager && ./test-env-reset.sh"
alias k8s-dev-logs="cd ~/k8s-manager && ./test-env-logs.sh"
alias k8s-dev-check="cd ~/k8s-manager && ./check-env.sh"

alias k8s-backend="cd ~/k8s-manager/backend && go run cmd/server/main.go"
alias k8s-frontend="cd ~/k8s-manager/frontend && npm start"

alias k8s-test-api="cd ~/k8s-manager && ./test-api.sh"
alias k8s-mysql="docker exec -it k8s-mgr-mysql-test mysql -uroot -ptest123456 k8s_manager_test"
```

### B. 常用代码片段

**Go错误定义**:
```go
// pkg/errors/errors.go
package errors

import "errors"

var (
    ErrInvalidParams    = errors.New("参数错误")
    ErrUnauthorized     = errors.New("未授权")
    ErrForbidden        = errors.New("禁止访问")
    ErrNotFound         = errors.New("资源不存在")
    ErrInternalServer   = errors.New("服务器内部错误")
    ErrDatabaseError    = errors.New("数据库错误")
    ErrK8SAPIError      = errors.New("K8S API调用失败")
)
```

**统一响应**:
```go
// pkg/response/response.go
package response

import (
    "net/http"
    "github.com/gin-gonic/gin"
)

type Response struct {
    Code    int         `json:"code"`
    Message string      `json:"message"`
    Data    interface{} `json:"data,omitempty"`
}

func Success(c *gin.Context, data interface{}) {
    c.JSON(http.StatusOK, Response{
        Code:    0,
        Message: "success",
        Data:    data,
    })
}

func Error(c *gin.Context, code int, message string) {
    c.JSON(http.StatusOK, Response{
        Code:    code,
        Message: message,
    })
}

func ErrorWithData(c *gin.Context, code int, message string, data interface{}) {
    c.JSON(http.StatusOK, Response{
        Code:    code,
        Message: message,
        Data:    data,
    })
}
```

### C. 开发日志模板

```markdown
# 开发日志

## 2025-11-11

### 完成的功能
- [x] 实现集群列表查询API
- [x] 实现集群创建API
- [x] 前端集群列表页面

### 遇到的问题
1. 问题: kubeconfig解析失败
   解决: 使用clientcmd.RESTConfigFromKubeConfig方法

2. 问题: 数据库连接超时
   解决: 增加连接超时时间配置

### 明天计划
- [ ] 实现集群删除功能
- [ ] 实现集群连接测试
- [ ] 优化错误提示

### 工作时长
- 后端开发: 3小时
- 前端开发: 2小时
- 调试测试: 1小时
- 总计: 6小时
```

---

## 总结

这份工作流规范涵盖了：

1. ✅ **完整的开发流程** - 从需求到交付
2. ✅ **环境准备检查** - 自动化脚本
3. ✅ **代码规范要求** - Go/TS/CSS
4. ✅ **测试流程清单** - 单元/集成/E2E
5. ✅ **提交规范模板** - Git Commit
6. ✅ **问题诊断方法** - 常见问题解决
7. ✅ **最佳实践指南** - 代码质量/性能
8. ✅ **完整检查清单** - 功能完成验证

**核心原则**:
- 🎯 **明确目标** - 每次开发前确认需求
- 🔄 **小步快跑** - 分步骤开发，边写边测
- ✅ **质量优先** - 代码质量比速度更重要
- 📝 **文档同步** - 代码和文档同步更新
- 🧪 **充分测试** - 不测试不提交

---

**文档版本**: v1.0  
**最后更新**: 2025-11-11  
**适用对象**: Codex + Jason  
**项目**: K8S集群管理平台
