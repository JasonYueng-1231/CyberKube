package workload

import (
    "net/http"

    "github.com/JasonYueng-1231/CyberKube/backend/internal/service"
    "github.com/gin-gonic/gin"
)

func RegisterDeployment(r *gin.RouterGroup) {
    // 列表: GET /api/v1/deployments?cluster=&namespace=
    r.GET("/deployments", func(c *gin.Context) {
        cluster := c.Query("cluster")
        ns := c.Query("namespace")
        if ns == "" { ns = "default" }
        items, err := service.ListDeployments(cluster, ns)
        if err != nil { c.JSON(http.StatusInternalServerError, gin.H{"code":50003, "message": err.Error()}); return }
        c.JSON(http.StatusOK, gin.H{"code":0, "message":"success", "data": gin.H{"items": items}})
    })
    // 伸缩: POST /api/v1/deployments/scale {cluster,namespace,name,replicas}
    r.POST("/deployments/scale", func(c *gin.Context) {
        var req struct{ Cluster, Namespace, Name string; Replicas int32 }
        if err := c.ShouldBindJSON(&req); err != nil { c.JSON(http.StatusBadRequest, gin.H{"code":40001, "message":"参数错误"}); return }
        if req.Namespace == "" { req.Namespace = "default" }
        if err := service.ScaleDeployment(req.Cluster, req.Namespace, req.Name, req.Replicas); err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"code":50003, "message": err.Error()}); return
        }
        c.JSON(http.StatusOK, gin.H{"code":0, "message":"success"})
    })
    // 重启: POST /api/v1/deployments/restart {cluster,namespace,name}
    r.POST("/deployments/restart", func(c *gin.Context) {
        var req struct{ Cluster, Namespace, Name string }
        if err := c.ShouldBindJSON(&req); err != nil { c.JSON(http.StatusBadRequest, gin.H{"code":40001, "message":"参数错误"}); return }
        if req.Namespace == "" { req.Namespace = "default" }
        if err := service.RestartDeployment(req.Cluster, req.Namespace, req.Name); err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"code":50003, "message": err.Error()}); return
        }
        c.JSON(http.StatusOK, gin.H{"code":0, "message":"success"})
    })
}
