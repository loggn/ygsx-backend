package model

import "time"

type Product struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"size:100" json:"name"`
	Picture   string    `gorm:"size:255" json:"picture"`
	Price     float64   `gorm:"type:decimal(10,2)" json:"price"`
	Stock     int       `gorm:"default:0" json:"stock"`
	Status    string    `gorm:"size:10" json:"status"` // on / off
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
