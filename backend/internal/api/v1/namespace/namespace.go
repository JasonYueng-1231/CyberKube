package namespace

import (
    "net/http"
    "github.com/JasonYueng-1231/CyberKube/backend/internal/service"
    "github.com/gin-gonic/gin"
)

func Register(r *gin.RouterGroup) {
    g := r.Group("/namespaces")
    g.GET("", func(c *gin.Context) {
        cluster := c.Query("cluster")
        items, err := service.ListNamespaces(cluster)
        if err != nil { c.JSON(http.StatusInternalServerError, gin.H{"code":50003, "message": err.Error()}); return }
        c.JSON(http.StatusOK, gin.H{"code":0, "message":"success", "data": gin.H{"items": items}})
    })
}

