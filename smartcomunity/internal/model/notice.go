package model

import "time"

type Notice struct {
	ID        int64     `gorm:"primaryKey" json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Publisher string    `json:"publisher"` // 例如："物业中心"
	ViewCount int       `json:"view_count"`
	CreatedAt time.Time `json:"created_at"`
}

func (Notice) TableName() string {
	return "cms_notice"
}

// NoticeRead 记录用户已读公告
type NoticeRead struct {
	ID       int64     `gorm:"primaryKey" json:"id"`
	UserID   int64     `json:"user_id"`
	NoticeID int64     `json:"notice_id"`
	ReadTime time.Time `json:"read_time"`
}

func (NoticeRead) TableName() string {
	return "cms_notice_read"
}
