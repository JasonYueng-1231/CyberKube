package middleware

import (
    "net/http"
    "strings"

    appjwt "github.com/JasonYueng-1231/CyberKube/backend/pkg/jwt"
    "github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        auth := c.GetHeader("Authorization")
        token := ""
        if auth != "" {
            // 兼容大小写与额外空格: "Bearer <token>"
            parts := strings.Fields(auth)
            if len(parts) == 2 && strings.ToLower(parts[0]) == "bearer" { token = parts[1] }
        }
        if token == "" {
            // 兼容 WebSocket/某些代理剥离 Header 的情况: ?token= 或 ?token=Bearer%20xxx
            q := c.Query("token")
            if strings.HasPrefix(strings.ToLower(q), "bearer ") { token = strings.TrimSpace(q[len("Bearer "):]) } else { token = q }
        }
        if token == "" {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"code": 40101, "message": "未授权"})
            return
        }
        claims, err := appjwt.Parse(token)
        if err != nil {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"code": 40101, "message": "未授权"})
            return
        }
        c.Set("userID", claims.UserID)
        c.Set("username", claims.Username)
        c.Set("role", claims.Role)
        c.Next()
    }
}
