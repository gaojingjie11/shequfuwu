package model

import "time"

// HotProduct 热销商品
type HotProduct struct {
	ID        int64     `gorm:"primaryKey" json:"id"`
	ProductID int64     `json:"product_id"`
	CreatedAt time.Time `json:"created_at"`

	// 关联商品
	Product Product `gorm:"foreignKey:ProductID" json:"product"`
}

func (HotProduct) TableName() string {
	return "pms_hot_product"
}
