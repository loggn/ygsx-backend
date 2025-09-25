package model

import "time"

type StockLog struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	ProductID  uint      `gorm:"index" json:"product_id"`
	ChangeNum  int       `json:"change_num"` // 正数=入库，负数=出库
	OperatorID uint      `json:"operator_id"`
	CreatedAt  time.Time `json:"created_at"`
}
