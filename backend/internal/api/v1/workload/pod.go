package workload

import (
    "net/http"
    "strconv"

    "github.com/JasonYueng-1231/CyberKube/backend/internal/service"
    "github.com/gin-gonic/gin"
    "sigs.k8s.io/yaml"
)

func RegisterPod(r *gin.RouterGroup) {
    // 列表: GET /api/v1/pods?cluster=&namespace=
    r.GET("/pods", func(c *gin.Context) {
        cluster := c.Query("cluster")
        ns := c.Query("namespace")
        if ns == "" { ns = "default" }
        items, err := service.ListPods(cluster, ns)
        if err != nil { c.JSON(http.StatusInternalServerError, gin.H{"code":50003, "message": err.Error()}); return }
        c.JSON(http.StatusOK, gin.H{"code":0, "message":"success", "data": gin.H{"items": items}})
    })
    // 日志: GET /api/v1/pods/logs?cluster=&namespace=&name=&tail=200
    r.GET("/pods/logs", func(c *gin.Context) {
        cluster := c.Query("cluster")
        ns := c.Query("namespace")
        name := c.Query("name")
        tailStr := c.Query("tail")
        if ns == "" { ns = "default" }
        var tail int64
        if tailStr != "" { if v, err := strconv.ParseInt(tailStr, 10, 64); err == nil { tail = v } }
        logs, err := service.GetPodLogs(cluster, ns, name, tail)
        if err != nil { c.JSON(http.StatusInternalServerError, gin.H{"code":50003, "message": err.Error()}); return }
        c.Data(http.StatusOK, "text/plain; charset=utf-8", []byte(logs))
    })
    // 事件: GET /api/v1/pods/events?cluster=&namespace=&name=
    r.GET("/pods/events", func(c *gin.Context) {
        cluster := c.Query("cluster"); ns := c.Query("namespace"); name := c.Query("name")
        if ns == "" { ns = "default" }
        items, err := service.ListPodEvents(cluster, ns, name)
        if err != nil { c.JSON(http.StatusInternalServerError, gin.H{"code":50003, "message": err.Error()}); return }
        c.JSON(http.StatusOK, gin.H{"code":0, "message":"success", "data": gin.H{"items": items}})
    })
    // 详情: GET /api/v1/pods/detail?cluster=&namespace=&name=
    r.GET("/pods/detail", func(c *gin.Context) {
        cluster := c.Query("cluster"); ns := c.Query("namespace"); name := c.Query("name")
        if ns == "" { ns = "default" }
        pod, err := service.GetPod(cluster, ns, name)
        if err != nil { c.JSON(http.StatusInternalServerError, gin.H{"code":50003, "message": err.Error()}); return }
        evs, _ := service.ListPodEvents(cluster, ns, name)
        yamlStr := ""
        if b, err := yaml.Marshal(pod); err == nil {
            yamlStr = string(b)
        }
        c.JSON(http.StatusOK, gin.H{
            "code":    0,
            "message": "success",
            "data": gin.H{
                "pod":    pod,
                "events": evs,
                "yaml":   yamlStr,
            },
        })
    })
}
