package database

import (
    "log"

    "github.com/JasonYueng-1231/CyberKube/backend/internal/model"
    "golang.org/x/crypto/bcrypt"
)

// EnsureDefaultAdmin 若无用户则创建默认管理员
func EnsureDefaultAdmin() error {
    var n int64
    if err := DB.Model(&model.User{}).Count(&n).Error; err != nil { return err }
    if n == 0 {
        hash, _ := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
        u := &model.User{Username: "admin", Password: string(hash), Nickname: "管理员", Role: "admin"}
        if err := DB.Create(u).Error; err != nil { return err }
        log.Println("created default admin: admin/admin123")
    }
    return nil
}

