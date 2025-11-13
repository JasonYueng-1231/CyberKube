package model

import "time"

type User struct {
    ID        uint      `gorm:"primaryKey" json:"id"`
    Username  string    `gorm:"uniqueIndex;size:100" json:"username"`
    Password  string    `json:"-"` // bcrypt hash
    Nickname  string    `gorm:"size:100" json:"nickname"`
    Role      string    `gorm:"size:20" json:"role"` // admin/user
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}

