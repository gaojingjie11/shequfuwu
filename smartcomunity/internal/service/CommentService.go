package service

import (
	"smartcommunity/internal/global"
	"smartcommunity/internal/model"
	"time"
)

type CommentService struct{}

// CreateComment 发表评论
func (s *CommentService) CreateComment(userID, productID int64, content string, rating int) error {
	comment := model.ProductComment{
		UserID:    userID,
		ProductID: productID,
		Content:   content,
		Rating:    rating,
		CreatedAt: time.Now(),
	}
	return global.DB.Create(&comment).Error
}

// GetCommentsByProductID 获取商品评论列表 (分页)
func (s *CommentService) GetCommentsByProductID(productID int64, page, size int) ([]model.ProductComment, int64, error) {
	var list []model.ProductComment
	var total int64

	db := global.DB.Model(&model.ProductComment{}).Where("product_id = ?", productID)
	db.Count(&total)

	offset := (page - 1) * size
	// Preload User to show username/avatar
	err := db.Preload("User").Order("created_at desc").Offset(offset).Limit(size).Find(&list).Error
	return list, total, err
}
