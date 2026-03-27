package model

import "time"

type ChatMessage struct {
	ID        int64     `gorm:"primaryKey" json:"id"`
	UserID    int64     `gorm:"column:user_id;index;not null" json:"user_id"`
	Role      string    `gorm:"column:role;type:varchar(16);not null" json:"role"`
	Content   string    `gorm:"column:content;type:longtext;not null" json:"content"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
}

func (ChatMessage) TableName() string {
	return "cms_chat_messages"
}
