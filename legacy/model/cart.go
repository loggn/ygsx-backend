package model

import "time"

type Cart struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `json:"user_id"`    // 关联用户ID
	ProductID uint      `json:"product_id"` // 商品ID
	Quantity  int       `json:"quantity"`   // 商品数量
	Checked   bool      `json:"checked"`    // 是否选中(结算)
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// 可选：方便前端展示
	Product Product `gorm:"foreignKey:ProductID" json:"product"`
}
