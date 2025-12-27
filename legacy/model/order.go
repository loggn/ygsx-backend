package model

import "time"

// 订单主表
type Order struct {
	ID          uint    `json:"id" gorm:"primaryKey"`
	UserID      uint    `json:"user_id" gorm:"index;not null"`
	OrderNumber string  `json:"order_number" gorm:"unique;not null"`
	TotalAmount float64 `json:"total_amount" gorm:"not null"`
	Status      string  `json:"status" gorm:"default:'pending'"` // pending, paid, shipped, finished, canceled

	// 地址快照（防止后期用户修改地址导致历史订单错乱）
	Receiver string `json:"receiver"`
	Phone    string `json:"phone"`
	Province string `json:"province"`
	City     string `json:"city"`
	District string `json:"district"`
	Detail   string `json:"detail"`

	// 一对多关联 —— 一个订单对应多个订单项
	Items []OrderItem `json:"items" gorm:"foreignKey:OrderID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
