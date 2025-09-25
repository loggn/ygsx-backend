package model

import "time"

type Order struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	UserID     uint      `gorm:"index" json:"user_id"`
	TotalPrice float64   `gorm:"type:decimal(10,2)" json:"total_price"`
	Status     string    `gorm:"size:20" json:"status"` // unpaid / paid / delivering / completed
	Address    string    `gorm:"size:255" json:"address"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
