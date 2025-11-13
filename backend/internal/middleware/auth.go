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
        if strings.HasPrefix(auth, "Bearer ") {
            token = strings.TrimPrefix(auth, "Bearer ")
        } else {
            // 兼容 WebSocket：通过 query token 传递
            q := c.Query("token")
            if q != "" && strings.HasPrefix(q, "Bearer ") { token = strings.TrimPrefix(q, "Bearer ") } else { token = q }
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
