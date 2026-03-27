package model

import "time"

type GreenPointRecord struct {
	ID        int64     `gorm:"primaryKey" json:"id"`
	UserID    int64     `gorm:"index;not null" json:"user_id"`
	Action    string    `gorm:"size:64;not null" json:"action"`
	Points    int       `gorm:"not null" json:"points"`
	CreatedAt time.Time `json:"created_at"`
}

func (GreenPointRecord) TableName() string {
	return "green_point_record"
}
