package database

import (
    "log"

    "github.com/JasonYueng-1231/CyberKube/backend/internal/config"
    "gorm.io/driver/mysql"
    "gorm.io/gorm"
    "gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitMySQL() error {
    if DB != nil { return nil }
    cfg := config.Load()
    gcfg := &gorm.Config{Logger: logger.Default.LogMode(logger.Warn)}
    db, err := gorm.Open(mysql.Open(cfg.MySQLDSN), gcfg)
    if err != nil {
        return err
    }
    DB = db
    log.Println("mysql connected")
    return nil
}

func AutoMigrate(models ...interface{}) error {
    return DB.AutoMigrate(models...)
}

