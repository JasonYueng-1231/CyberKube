package serviceapi

import (
    "net/http"
    svc "github.com/JasonYueng-1231/CyberKube/backend/internal/service"
    "github.com/gin-gonic/gin"
)

func RegisterService(r *gin.RouterGroup) {
    g := r.Group("/services")
    g.GET("", func(c *gin.Context) {
        cluster := c.Query("cluster"); ns := c.Query("namespace"); if ns=="" { ns = "default" }
        items, err := svc.ListServices(cluster, ns)
        if err != nil { c.JSON(http.StatusInternalServerError, gin.H{"code":50003, "message": err.Error()}); return }
        c.JSON(http.StatusOK, gin.H{"code":0, "message":"success", "data": gin.H{"items": items}})
    })
    g.GET("/detail", func(c *gin.Context) {
        cluster := c.Query("cluster"); ns := c.Query("namespace"); name := c.Query("name"); if ns=="" { ns = "default" }
        d, err := svc.GetService(cluster, ns, name)
        if err != nil { c.JSON(http.StatusInternalServerError, gin.H{"code":50003, "message": err.Error()}); return }
        c.JSON(http.StatusOK, gin.H{"code":0, "message":"success", "data": d})
    })
    g.POST("/yaml", func(c *gin.Context) {
        var req struct{ Cluster, Namespace, Yaml string }
        if err := c.ShouldBindJSON(&req); err != nil { c.JSON(http.StatusBadRequest, gin.H{"code":40001, "message":"参数错误"}); return }
        if req.Namespace=="" { req.Namespace = "default" }
        if err := svc.CreateServiceYAML(req.Cluster, req.Namespace, req.Yaml); err != nil { c.JSON(http.StatusBadRequest, gin.H{"code":40901, "message": err.Error()}); return }
        c.JSON(http.StatusOK, gin.H{"code":0, "message":"success"})
    })
    g.PUT("/yaml", func(c *gin.Context) {
        var req struct{ Cluster, Namespace, Yaml string }
        if err := c.ShouldBindJSON(&req); err != nil { c.JSON(http.StatusBadRequest, gin.H{"code":40001, "message":"参数错误"}); return }
        if req.Namespace=="" { req.Namespace = "default" }
        if err := svc.UpdateServiceYAML(req.Cluster, req.Namespace, req.Yaml); err != nil { c.JSON(http.StatusBadRequest, gin.H{"code":40901, "message": err.Error()}); return }
        c.JSON(http.StatusOK, gin.H{"code":0, "message":"success"})
    })
    g.DELETE("", func(c *gin.Context) {
        cluster := c.Query("cluster"); ns := c.Query("namespace"); name := c.Query("name"); if ns=="" { ns = "default" }
        if err := svc.DeleteService(cluster, ns, name); err != nil { c.JSON(http.StatusInternalServerError, gin.H{"code":50003, "message": err.Error()}); return }
        c.JSON(http.StatusOK, gin.H{"code":0, "message":"success"})
    })
}

