package v1

import (
    "github.com/JasonYueng-1231/CyberKube/backend/internal/api/v1/health"
    "github.com/gin-gonic/gin"
)

// RegisterHealth 注册健康检查路由
func RegisterHealth(r *gin.RouterGroup) {
    health.Register(r)
}

