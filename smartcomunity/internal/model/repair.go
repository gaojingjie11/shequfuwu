package model

import "time"

type Repair struct {
	ID        int64     `gorm:"primaryKey" json:"id"`
	UserID    int64     `json:"user_id"`
	Type      int       `json:"type"`     // 1:报修 2:投诉
	Category  string    `json:"category"` // 如: 水电, 门窗, 噪音, 服务态度
	Content   string    `json:"content"`  // 详细描述
	Status    int       `json:"status"`   // 0:待处理 1:处理中 2:已完成
	Result    string    `json:"result"`   // 处理结果/回复
	CreatedAt time.Time `json:"created_at"`

	// 关联用户
	User SysUser `gorm:"foreignKey:UserID" json:"user"`
}

func (Repair) TableName() string {
	return "cms_repair"
}
