package cluster

import (
    "net/http"

    "github.com/JasonYueng-1231/CyberKube/backend/internal/service"
    "github.com/gin-gonic/gin"
)

func Register(r *gin.RouterGroup) {
    g := r.Group("/clusters")
    g.GET("", func(c *gin.Context) {
        list, err := service.ListClusters()
        if err != nil { c.JSON(http.StatusInternalServerError, gin.H{"code":50001, "message":err.Error()}); return }
        c.JSON(http.StatusOK, gin.H{"code":0, "message":"success", "data": gin.H{"items": list}})
    })
    g.POST("", func(c *gin.Context) {
        var req service.ClusterCreateReq
        if err := c.ShouldBindJSON(&req); err != nil { c.JSON(http.StatusBadRequest, gin.H{"code":40001, "message":"参数错误"}); return }
        res, err := service.CreateCluster(req)
        if err != nil { c.JSON(http.StatusBadRequest, gin.H{"code":40901, "message": err.Error()}); return }
        c.JSON(http.StatusOK, gin.H{"code":0, "message":"success", "data": res})
    })
    g.DELETE(":name", func(c *gin.Context) {
        if err := service.DeleteCluster(c.Param("name")); err != nil { c.JSON(http.StatusInternalServerError, gin.H{"code":50001, "message":err.Error()}); return }
        c.JSON(http.StatusOK, gin.H{"code":0, "message":"success"})
    })
    g.GET(":name/health", func(c *gin.Context) {
        cli, err := service.GetClientForCluster(c.Param("name"))
        if err != nil { c.JSON(http.StatusBadRequest, gin.H{"code":40901, "message": err.Error()}); return }
        info, err := service.SimplePing(cli)
        if err != nil { c.JSON(http.StatusInternalServerError, gin.H{"code":50003, "message": err.Error()}); return }
        c.JSON(http.StatusOK, gin.H{"code":0, "message":"success", "data": info})
    })
}
