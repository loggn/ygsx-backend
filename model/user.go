package model

import "time"

type User struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Phone     string    `gorm:"uniqueIndex;size:20" json:"phone"`
	Username  string    `gorm:"size:50" json:"username"`
	Role      string    `gorm:"size:20;default:customer" json:"role"`
	CreatedAt time.Time `json:"created_at"`
}
