package model

import "time"

// SysTransaction 对应数据库 sys_transaction 表
type SysTransaction struct {
	ID        int64     `gorm:"primaryKey" json:"id"`
	UserID    int64     `json:"user_id"`
	Type      int       `json:"type"`       // 1:商城订单 2:物业费
	Amount    float64   `json:"amount"`     // 交易金额 (保持与 SysUser.Balance 一致使用 float64)
	RelatedID int64     `json:"related_id"` // 关联的 ID (Order.ID 或 PropertyFee.ID)
	Remark    string    `json:"remark"`     // 备注，例如 "支付商城订单"
	CreatedAt time.Time `json:"created_at"` // 流水必须要有时间，否则无法对账
}

func (SysTransaction) TableName() string {
	return "sys_transaction"
}
