package api

import (
    "net/http"
    "time"

    v1 "github.com/JasonYueng-1231/CyberKube/backend/internal/api/v1"
    "github.com/JasonYueng-1231/CyberKube/backend/internal/middleware"
    "github.com/gin-gonic/gin"
)

// SetupRouter 初始化 Gin 引擎与路由
func SetupRouter() *gin.Engine {
    r := gin.New()

    // 基础中间件
    r.Use(gin.Recovery())
    r.Use(middleware.Logger())
    r.Use(middleware.CORS())

    // 健康检查（无前缀）
    r.GET("/healthz", func(c *gin.Context) {
        c.String(http.StatusOK, "ok")
    })

    // API v1
    apiV1 := r.Group("/api/v1")
    {
        v1.RegisterHealth(apiV1)
        v1.RegisterAuth(apiV1)
        v1.RegisterCluster(apiV1)
        v1.RegisterWorkload(apiV1)
        v1.RegisterNamespace(apiV1)
        v1.RegisterMetrics(apiV1)
        v1.RegisterConfig(apiV1)
        v1.RegisterServiceAPI(apiV1)
    }

    // 统一 404
    r.NoRoute(func(c *gin.Context) {
        c.JSON(http.StatusNotFound, gin.H{
            "code":    40401,
            "message": "资源不存在",
            "path":    c.Request.URL.Path,
            "time":    time.Now().Format(time.RFC3339),
        })
    })

    return r
}
