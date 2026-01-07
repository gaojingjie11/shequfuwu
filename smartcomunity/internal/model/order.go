package model

import "time"

// Order 极简订单主表
type Order struct {
	ID          int64     `gorm:"primaryKey" json:"id"`
	OrderNo     string    `json:"order_no"`
	UserID      int64     `json:"user_id"`
	StoreID     int64     `json:"store_id"`
	TotalAmount float64   `json:"total_amount"`
	Status      int       `json:"status"`
	CreatedAt   time.Time `json:"created_at"` // 只保留创建时间

	// 关联明细 (GORM 依然会自动处理这个)
	Items []OrderItem `gorm:"foreignKey:OrderID" json:"items"`

	// 关联门店
	Store Store `gorm:"foreignKey:StoreID" json:"store"`

	// 关联用户
	SysUser SysUser `gorm:"foreignKey:UserID" json:"sys_user"`
}

func (Order) TableName() string {
	return "oms_order"
}

// OrderItem 极简订单明细
type OrderItem struct {
	ID        int64   `gorm:"primaryKey" json:"id"`
	OrderID   int64   `json:"order_id"`
	ProductID int64   `json:"product_id"`
	Price     float64 `json:"price"`
	Quantity  int     `json:"quantity"`

	// 关联商品 (方便前端展示商品名和图片)
	Product Product `gorm:"foreignKey:ProductID" json:"product"`
}

func (OrderItem) TableName() string {
	return "oms_order_item"
}
