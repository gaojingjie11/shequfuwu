package model

import "time"

type PropertyFee struct {
	ID          int64      `gorm:"primaryKey" json:"id"`
	UserID      int64      `json:"user_id"`
	Month       string     `gorm:"type:varchar(20)" json:"month"`
	Amount      float64    `gorm:"type:decimal(10,2);not null;default:0.00" json:"amount"`
	UsedPoints  int        `gorm:"column:used_points;not null;default:0" json:"used_points"`
	UsedBalance float64    `gorm:"column:used_balance;type:decimal(10,2);not null;default:0.00" json:"used_balance"`
	Status      int        `json:"status"`
	PayTime     *time.Time `json:"pay_time"`
}

func (PropertyFee) TableName() string {
	return "cms_property_fee"
}
