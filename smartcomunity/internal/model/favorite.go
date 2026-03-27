package model

import "time"

// Favorite 商品收藏
type Favorite struct {
	ID        int64     `gorm:"primaryKey" json:"id"`
	UserID    int64     `json:"user_id"`
	ProductID int64     `json:"product_id"`
	CreatedAt time.Time `json:"created_at"`

	// 关联商品
	Product Product `gorm:"foreignKey:ProductID" json:"product"`
}

func (Favorite) TableName() string {
	return "pms_favorite"
}
