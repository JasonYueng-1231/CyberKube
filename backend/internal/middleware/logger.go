package middleware

import (
    "log"
    "time"

    "github.com/gin-gonic/gin"
)

// Logger 简易访问日志中间件
func Logger() gin.HandlerFunc {
    return func(c *gin.Context) {
        start := time.Now()
        path := c.Request.URL.Path
        method := c.Request.Method

        c.Next()

        status := c.Writer.Status()
        latency := time.Since(start)
        log.Printf("%s %s -> %d (%s)", method, path, status, latency)
    }
}

