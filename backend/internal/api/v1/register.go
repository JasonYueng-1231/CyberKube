package v1

import (
    "github.com/JasonYueng-1231/CyberKube/backend/internal/api/v1/auth"
    "github.com/JasonYueng-1231/CyberKube/backend/internal/api/v1/health"
    "github.com/JasonYueng-1231/CyberKube/backend/internal/api/v1/cluster"
    "github.com/JasonYueng-1231/CyberKube/backend/internal/api/v1/namespace"
    "github.com/JasonYueng-1231/CyberKube/backend/internal/api/v1/metrics"
    confapi "github.com/JasonYueng-1231/CyberKube/backend/internal/api/v1/config"
    svcapi "github.com/JasonYueng-1231/CyberKube/backend/internal/api/v1/service"
    netapi "github.com/JasonYueng-1231/CyberKube/backend/internal/api/v1/network"
    "github.com/JasonYueng-1231/CyberKube/backend/internal/api/v1/workload"
    "github.com/JasonYueng-1231/CyberKube/backend/internal/middleware"
    "github.com/gin-gonic/gin"
)

// RegisterHealth 注册健康检查路由
func RegisterHealth(r *gin.RouterGroup) {
    health.Register(r)
}

// RegisterAuth 注册认证路由
func RegisterAuth(r *gin.RouterGroup) { auth.Register(r) }

// RegisterCluster 注册集群管理路由（需鉴权）
func RegisterCluster(r *gin.RouterGroup) {
    g := r.Group("")
    g.Use(middleware.AuthMiddleware())
    cluster.Register(g)
}

func RegisterWorkload(r *gin.RouterGroup) {
    g := r.Group("")
    g.Use(middleware.AuthMiddleware())
    workload.RegisterDeployment(g)
    workload.RegisterPod(g)
    workload.RegisterStream(g)
}

func RegisterNamespace(r *gin.RouterGroup) {
    g := r.Group("")
    g.Use(middleware.AuthMiddleware())
    namespace.Register(g)
}

func RegisterMetrics(r *gin.RouterGroup) {
    g := r.Group("")
    g.Use(middleware.AuthMiddleware())
    metrics.Register(g)
}

func RegisterConfig(r *gin.RouterGroup) {
    g := r.Group("")
    g.Use(middleware.AuthMiddleware())
    confapi.RegisterConfigMap(g)
    confapi.RegisterSecret(g)
}

func RegisterServiceAPI(r *gin.RouterGroup) {
    g := r.Group("")
    g.Use(middleware.AuthMiddleware())
    svcapi.RegisterService(g)
    netapi.RegisterIngress(g)
}
