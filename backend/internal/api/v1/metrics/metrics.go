package metrics

import (
    "net/http"
    "github.com/JasonYueng-1231/CyberKube/backend/internal/service"
    "github.com/gin-gonic/gin"
)

func Register(r *gin.RouterGroup) {
    g := r.Group("/metrics")
    g.GET("/overview", func(c *gin.Context) {
        cluster := c.Query("cluster")
        ov, err := service.GetOverview(cluster)
        if err != nil { c.JSON(http.StatusInternalServerError, gin.H{"code":50003, "message": err.Error()}); return }
        c.JSON(http.StatusOK, gin.H{"code":0, "message":"success", "data": ov})
    })
}

