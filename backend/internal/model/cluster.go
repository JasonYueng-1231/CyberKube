package model

import "time"

type Cluster struct {
    ID        uint      `gorm:"primaryKey" json:"id"`
    Name      string    `gorm:"uniqueIndex;size:100" json:"name"`
    Alias     string    `gorm:"size:100" json:"alias"`
    KubeconfigEnc string `json:"-"` // AES 加密内容
    APIServer string    `gorm:"size:255" json:"api_server"`
    Version   string    `gorm:"size:50" json:"version"`
    Status    int       `gorm:"default:1" json:"status"`
    Description string  `gorm:"size:500" json:"description"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}

