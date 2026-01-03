package model

import "time"

type Cart struct {
	ID        int64     `gorm:"primaryKey" json:"id"`
	UserID    int64     `json:"user_id"`
	ProductID int64     `json:"product_id"`
	Quantity  int       `json:"quantity"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// 关联关系：一个购物车条目对应一个商品
	// gorm:"foreignKey:ProductID" 告诉 GORM 用 ProductID 去找 Product 表
	Product Product `gorm:"foreignKey:ProductID" json:"product"`
}

// internal/model/cart.go
type CartItemParam struct {
	CartID   int64 `json:"cart_id"`
	Quantity int   `json:"quantity"`
}

func (Cart) TableName() string {
	return "oms_cart"
}
