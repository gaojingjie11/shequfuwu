package model

import "time"

// PropertyFee 物业费账单
type PropertyFee struct {
	ID      int64      `gorm:"primaryKey" json:"id"`
	UserID  int64      `json:"user_id"`
	Month   string     `json:"month"`    // 例如 "2023-10"
	Amount  float64    `json:"amount"`   // 金额
	Status  int        `json:"status"`   // 0:未缴 1:已缴
	PayTime *time.Time `json:"pay_time"` // 支付时间 (指针允许为null)
	//CreatedAt time.Time  `json:"created_at"`
}

func (PropertyFee) TableName() string {
	return "cms_property_fee"
}
