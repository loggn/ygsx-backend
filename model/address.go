package model

import "time"

type Address struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	UserID    uint      `json:"user_id" gorm:"not null;index"`
	Receiver  string    `json:"receiver" gorm:"not null"` // 收件人姓名
	Phone     string    `json:"phone" gorm:"not null"`    // 收件人电话
	Province  string    `json:"province"`
	City      string    `json:"city"`
	District  string    `json:"district"`
	Detail    string    `json:"detail" gorm:"not null"`          // 详细地址
	IsDefault bool      `json:"is_default" gorm:"default:false"` // 是否默认
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
