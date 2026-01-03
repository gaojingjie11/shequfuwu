package model

import "time"

type Product struct {
	ID            int64     `gorm:"primaryKey" json:"id"`
	CategoryName  string    `json:"category_name"`
	Name          string    `json:"name"`
	Description   string    `json:"description"`
	Price         float64   `json:"price"`
	OriginalPrice float64   `json:"original_price"`
	Stock         int       `json:"stock"`
	ImageURL      string    `gorm:"column:image_url" json:"image_url"`
	IsPromotion   int       `json:"is_promotion"` // 1:是 0:否
	Sales         int       `json:"sales"`
	Status        int       `json:"status"` // 1:上架 0:下架
	CreatedAt     time.Time `json:"created_at"`
	// 关联分类ID (新增)
	CategoryID int64 `json:"category_id"`
}

func (Product) TableName() string {
	return "pms_product"
}

// ProductCategory 商品分类
type ProductCategory struct {
	ID   int64  `gorm:"primaryKey" json:"id"`
	Name string `json:"name"`
	Icon string `json:"icon"`
	Sort int    `json:"sort"`
}

func (ProductCategory) TableName() string {
	return "pms_product_category"
}

// Promotion 促销活动
type Promotion struct {
	ID        int64     `gorm:"primaryKey" json:"id"`
	Title     string    `json:"title"`
	Type      int       `json:"type"` // 1:秒杀 2:满减
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
	Status    int       `json:"status"` // 1进行中 0已结束
	ProductID int64     `json:"product_id"`
}

func (Promotion) TableName() string {
	return "pms_promotion"
}
