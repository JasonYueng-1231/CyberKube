package config

import (
    "net/http"
    "github.com/JasonYueng-1231/CyberKube/backend/internal/service"
    "github.com/gin-gonic/gin"
)

func RegisterConfigMap(r *gin.RouterGroup) {
    g := r.Group("/configmaps")
    g.GET("", func(c *gin.Context) {
        cluster := c.Query("cluster"); ns := c.Query("namespace"); if ns=="" { ns = "default" }
        items, err := service.ListConfigMaps(cluster, ns)
        if err != nil { c.JSON(http.StatusInternalServerError, gin.H{"code":50003, "message": err.Error()}); return }
        c.JSON(http.StatusOK, gin.H{"code":0, "message":"success", "data": gin.H{"items": items}})
    })
    g.GET("/detail", func(c *gin.Context) {
        cluster := c.Query("cluster"); ns := c.Query("namespace"); name := c.Query("name"); if ns=="" { ns = "default" }
        d, err := service.GetConfigMap(cluster, ns, name)
        if err != nil { c.JSON(http.StatusInternalServerError, gin.H{"code":50003, "message": err.Error()}); return }
        c.JSON(http.StatusOK, gin.H{"code":0, "message":"success", "data": d})
    })
    g.POST("/yaml", func(c *gin.Context) {
        var req struct{ Cluster, Namespace, Yaml string }
        if err := c.ShouldBindJSON(&req); err != nil { c.JSON(http.StatusBadRequest, gin.H{"code":40001, "message":"参数错误"}); return }
        if req.Namespace=="" { req.Namespace = "default" }
        if err := service.CreateConfigMapYAML(req.Cluster, req.Namespace, req.Yaml); err != nil { c.JSON(http.StatusBadRequest, gin.H{"code":40901, "message": err.Error()}); return }
        c.JSON(http.StatusOK, gin.H{"code":0, "message":"success"})
    })
    g.PUT("/yaml", func(c *gin.Context) {
        var req struct{ Cluster, Namespace, Yaml string }
        if err := c.ShouldBindJSON(&req); err != nil { c.JSON(http.StatusBadRequest, gin.H{"code":40001, "message":"参数错误"}); return }
        if req.Namespace=="" { req.Namespace = "default" }
        if err := service.UpdateConfigMapYAML(req.Cluster, req.Namespace, req.Yaml); err != nil { c.JSON(http.StatusBadRequest, gin.H{"code":40901, "message": err.Error()}); return }
        c.JSON(http.StatusOK, gin.H{"code":0, "message":"success"})
    })
    g.DELETE("", func(c *gin.Context) {
        cluster := c.Query("cluster"); ns := c.Query("namespace"); name := c.Query("name"); if ns=="" { ns = "default" }
        if err := service.DeleteConfigMap(cluster, ns, name); err != nil { c.JSON(http.StatusInternalServerError, gin.H{"code":50003, "message": err.Error()}); return }
        c.JSON(http.StatusOK, gin.H{"code":0, "message":"success"})
    })
}

