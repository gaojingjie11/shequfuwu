package model

type Store struct {
	ID            int64  `gorm:"primaryKey" json:"id"`
	Name          string `gorm:"type:varchar(128)" json:"name"`
	Address       string `gorm:"type:varchar(255)" json:"address"`
	Phone         string `gorm:"type:varchar(32)" json:"phone"`
	Region        string `gorm:"type:varchar(128)" json:"region"`
	BusinessHours string `gorm:"column:business_hours;type:varchar(64)" json:"business_hours"`
}

func (Store) TableName() string {
	return "pms_store"
}

type StoreProduct struct {
	ID        int64 `gorm:"primaryKey" json:"id"`
	StoreID   int64 `json:"store_id"`
	ProductID int64 `json:"product_id"`
	Stock     int   `json:"stock"`
}

func (StoreProduct) TableName() string {
	return "pms_store_product"
}
