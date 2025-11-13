package service

import (
    "errors"
    "time"

    "github.com/JasonYueng-1231/CyberKube/backend/internal/database"
    "github.com/JasonYueng-1231/CyberKube/backend/internal/model"
    appjwt "github.com/JasonYueng-1231/CyberKube/backend/pkg/jwt"
    "golang.org/x/crypto/bcrypt"
)

func Register(username, password, nickname string) (*model.User, error) {
    var count int64
    database.DB.Model(&model.User{}).Where("username=?", username).Count(&count)
    if count > 0 {
        return nil, errors.New("用户名已存在")
    }
    hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    u := &model.User{Username: username, Password: string(hash), Nickname: nickname, Role: "admin"}
    if err := database.DB.Create(u).Error; err != nil { return nil, err }
    return u, nil
}

func Login(username, password string) (string, *model.User, error) {
    var u model.User
    if err := database.DB.Where("username=?", username).First(&u).Error; err != nil {
        return "", nil, errors.New("用户不存在")
    }
    if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)); err != nil {
        return "", nil, errors.New("密码错误")
    }
    token, err := appjwt.Sign(u.ID, u.Username, u.Role, 2*time.Hour)
    if err != nil { return "", nil, err }
    return token, &u, nil
}

