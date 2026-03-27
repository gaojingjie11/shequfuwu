package model

import "time"

type Repair struct {
	ID        int64     `gorm:"primaryKey" json:"id"`
	UserID    int64     `json:"user_id"`
	Type      int       `json:"type"`
	Category  string    `json:"category"`
	Content   string    `json:"content"`
	Status    int       `json:"status"`
	Result    string    `json:"result"`
	CreatedAt time.Time `json:"created_at"`
	User      SysUser   `gorm:"foreignKey:UserID" json:"user"`
}

func (Repair) TableName() string {
	return "cms_repair"
}
