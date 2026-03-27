package model

import "time"

type Order struct {
	ID          int64       `gorm:"primaryKey" json:"id"`
	OrderNo     string      `gorm:"column:order_no;type:varchar(64)" json:"order_no"`
	UserID      int64       `json:"user_id"`
	StoreID     int64       `json:"store_id"`
	TotalAmount float64     `gorm:"type:decimal(10,2);not null;default:0.00" json:"total_amount"`
	UsedPoints  int         `gorm:"column:used_points;not null;default:0" json:"used_points"`
	UsedBalance float64     `gorm:"column:used_balance;type:decimal(10,2);not null;default:0.00" json:"used_balance"`
	Status      int         `json:"status"`
	PaidAt      *time.Time  `gorm:"column:paid_at" json:"paid_at"`
	CreatedAt   time.Time   `json:"created_at"`
	Items       []OrderItem `gorm:"foreignKey:OrderID" json:"items"`
	Store       Store       `gorm:"foreignKey:StoreID" json:"store"`
	SysUser     SysUser     `gorm:"foreignKey:UserID" json:"sys_user"`
}

func (Order) TableName() string {
	return "oms_order"
}

type OrderItem struct {
	ID        int64   `gorm:"primaryKey" json:"id"`
	OrderID   int64   `json:"order_id"`
	ProductID int64   `json:"product_id"`
	Price     float64 `json:"price"`
	Quantity  int     `json:"quantity"`
	Product   Product `gorm:"foreignKey:ProductID" json:"product"`
}

func (OrderItem) TableName() string {
	return "oms_order_item"
}
