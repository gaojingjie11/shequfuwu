package model

import "time"

type SysTransaction struct {
	ID        int64     `gorm:"primaryKey" json:"id"`
	UserID    int64     `json:"user_id"`
	Type      int       `json:"type"`
	Amount    float64   `gorm:"type:decimal(10,2);not null;default:0.00" json:"amount"`
	RelatedID int64     `json:"related_id"`
	Remark    string    `json:"remark"`
	CreatedAt time.Time `json:"created_at"`
}

func (SysTransaction) TableName() string {
	return "sys_transaction"
}
