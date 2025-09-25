package model

import "time"

type User struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	OpenID     string    `gorm:"uniqueIndex;size:64" json:"openid"`
	Phone      string    `gorm:"size:20" json:"phone"`
	Role       string    `gorm:"size:20" json:"role"` // customer / manager / keeper / boss
	InviteCode string    `gorm:"size:50" json:"invite_code"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
