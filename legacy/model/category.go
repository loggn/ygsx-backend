package model

type Category struct {
	ID   uint   `gorm:"primaryKey" json:"id"`
	Name string `gorm:"size:50;not null" json:"name"`
}
