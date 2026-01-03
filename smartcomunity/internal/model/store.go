package model

type Store struct {
	ID            int64  `gorm:"primaryKey" json:"id"`
	Name          string `json:"name"`           // 门店名称
	Address       string `json:"address"`        // 详细地址
	Phone         string `json:"phone"`          // 联系电话
	Region        string `json:"region"`         // 区域 (如 "北京市海淀区")
	BusinessHours string `json:"business_hours"` // 营业时间 "09:00-22:00"
}

func (Store) TableName() string {
	return "pms_store"
}

// StoreProduct 门店商品关联 (库存分配)
type StoreProduct struct {
	ID        int64 `gorm:"primaryKey" json:"id"`
	StoreID   int64 `json:"store_id"`
	ProductID int64 `json:"product_id"`
	Stock     int   `json:"stock"` // 该门店分配的库存
}

func (StoreProduct) TableName() string {
	return "pms_store_product"
}
