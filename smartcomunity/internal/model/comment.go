package model

import "time"

// ProductComment 商品评论
type ProductComment struct {
	ID        int64     `gorm:"primaryKey" json:"id"`
	UserID    int64     `json:"user_id"`
	ProductID int64     `json:"product_id"`
	Content   string    `json:"content"`
	Rating    int       `json:"rating"` // 1-5分
	CreatedAt time.Time `json:"created_at"`

	// 关联
	User SysUser `gorm:"foreignKey:UserID" json:"user,omitempty"`
}
