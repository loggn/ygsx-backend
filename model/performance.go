package model

import "time"

type Performance struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"index" json:"user_id"` // 负责人 ID
	Sales     float64   `gorm:"type:decimal(10,2)" json:"sales"`
	Bonus     float64   `gorm:"type:decimal(10,2)" json:"bonus"`
	UserCount int       `json:"user_count"`
	CreatedAt time.Time `json:"created_at"`
}
