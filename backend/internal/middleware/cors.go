package middleware

import "github.com/gin-gonic/gin"

// CORS 允许开发阶段的跨域请求
func CORS() gin.HandlerFunc {
    return func(c *gin.Context) {
        h := c.Writer.Header()
        h.Set("Access-Control-Allow-Origin", "*")
        h.Set("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,PATCH,OPTIONS")
        h.Set("Access-Control-Allow-Headers", "Authorization,Content-Type,X-Requested-With")
        h.Set("Access-Control-Expose-Headers", "Content-Disposition")

        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(204)
            return
        }
        c.Next()
    }
}

