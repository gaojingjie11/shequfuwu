package model

import "time"

type Product struct {
	ID            int64     `gorm:"primaryKey" json:"id"`
	CategoryName  string    `gorm:"column:category_name;type:varchar(64)" json:"category_name"`
	Name          string    `gorm:"type:varchar(128)" json:"name"`
	Description   string    `gorm:"type:text" json:"description"`
	Price         float64   `gorm:"type:decimal(10,2);not null;default:0.00" json:"price"`
	OriginalPrice float64   `gorm:"type:decimal(10,2);not null;default:0.00" json:"original_price"`
	Stock         int       `json:"stock"`
	ImageURL      string    `gorm:"column:image_url;type:varchar(255)" json:"image_url"`
	IsPromotion   int       `json:"is_promotion"`
	Sales         int       `json:"sales"`
	Status        int       `json:"status"`
	CreatedAt     time.Time `json:"created_at"`
	CategoryID    int64     `json:"category_id"`
}

func (Product) TableName() string {
	return "pms_product"
}

type ProductCategory struct {
	ID   int64  `gorm:"primaryKey" json:"id"`
	Name string `gorm:"type:varchar(64)" json:"name"`
	Icon string `gorm:"type:varchar(255)" json:"icon"`
	Sort int    `json:"sort"`
}

func (ProductCategory) TableName() string {
	return "pms_product_category"
}

type Promotion struct {
	ID        int64     `gorm:"primaryKey" json:"id"`
	Title     string    `gorm:"type:varchar(128)" json:"title"`
	Type      int       `json:"type"`
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
	Status    int       `json:"status"`
	ProductID int64     `json:"product_id"`
}

func (Promotion) TableName() string {
	return "pms_promotion"
}
