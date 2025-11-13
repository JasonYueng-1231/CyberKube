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
        if auth == "" || !strings.HasPrefix(auth, "Bearer ") {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"code": 40101, "message": "未授权"})
            return
        }
        token := strings.TrimPrefix(auth, "Bearer ")
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

