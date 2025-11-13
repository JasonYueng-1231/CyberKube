package health

import (
    "net/http"
    "runtime"
    "time"

    "github.com/gin-gonic/gin"
)

func Register(r *gin.RouterGroup) {
    g := r.Group("/health")
    g.GET("", func(c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{
            "code":    0,
            "message": "success",
            "data": gin.H{
                "status":  "ok",
                "time":    time.Now().Format(time.RFC3339),
                "go":      runtime.Version(),
                "service": "cyberkube-backend",
            },
        })
    })
    g.GET("/ping", func(c *gin.Context) { c.String(http.StatusOK, "pong") })
}

