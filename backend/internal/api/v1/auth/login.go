package auth

import (
    "net/http"

    "github.com/JasonYueng-1231/CyberKube/backend/internal/service"
    "github.com/gin-gonic/gin"
)

type loginReq struct { Username string `json:"username"` ; Password string `json:"password"` }

func Register(r *gin.RouterGroup) {
    g := r.Group("/auth")
    g.POST("/register", func(c *gin.Context) {
        var req struct{ Username, Password, Nickname string }
        if err := c.ShouldBindJSON(&req); err != nil { c.JSON(http.StatusBadRequest, gin.H{"code":40001,"message":"参数错误"}); return }
        u, err := service.Register(req.Username, req.Password, req.Nickname)
        if err != nil { c.JSON(http.StatusBadRequest, gin.H{"code":40901,"message":err.Error()}); return }
        c.JSON(http.StatusOK, gin.H{"code":0, "message":"success", "data": gin.H{"user": u}})
    })
    g.POST("/login", func(c *gin.Context) {
        var req loginReq
        if err := c.ShouldBindJSON(&req); err != nil { c.JSON(http.StatusBadRequest, gin.H{"code":40001,"message":"参数错误"}); return }
        token, u, err := service.Login(req.Username, req.Password)
        if err != nil { c.JSON(http.StatusUnauthorized, gin.H{"code":40101,"message":err.Error()}); return }
        c.JSON(http.StatusOK, gin.H{"code":0, "message":"success", "data": gin.H{"token": token, "user": gin.H{"id":u.ID, "username":u.Username, "nickname":u.Nickname, "role":u.Role}}})
    })
}

