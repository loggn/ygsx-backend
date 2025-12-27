package model

import "time"

// 首页轮播图广告图片
type Banner struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	ImageURL  string    `json:"image_url" gorm:"not null"`
	Link      string    `json:"link"`
	Sort      int       `json:"sort" gorm:"default:0"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
