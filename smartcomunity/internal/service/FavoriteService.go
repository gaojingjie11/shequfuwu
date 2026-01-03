package service

import (
	"errors"
	"smartcommunity/internal/global"
	"smartcommunity/internal/model"
)

type FavoriteService struct{}

// Add 添加收藏
func (s *FavoriteService) Add(userID, productID int64) error {
	var count int64
	global.DB.Model(&model.Favorite{}).Where("user_id = ? AND product_id = ?", userID, productID).Count(&count)
	if count > 0 {
		return errors.New("已收藏该商品")
	}

	fav := model.Favorite{
		UserID:    userID,
		ProductID: productID,
	}
	return global.DB.Create(&fav).Error
}

// Delete 取消收藏
func (s *FavoriteService) Delete(userID, productID int64) error {
	result := global.DB.Where("user_id = ? AND product_id = ?", userID, productID).Delete(&model.Favorite{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("未收藏该商品")
	}
	return nil
}

// GetList 获取收藏列表
func (s *FavoriteService) GetList(userID int64, page, size int) ([]model.Favorite, int64, error) {
	var list []model.Favorite
	var total int64

	db := global.DB.Model(&model.Favorite{}).Where("user_id = ?", userID)

	db.Count(&total)

	offset := (page - 1) * size
	err := db.Preload("Product").Order("created_at desc").Offset(offset).Limit(size).Find(&list).Error

	return list, total, err
}

// Check 是否收藏
func (s *FavoriteService) Check(userID, productID int64) (bool, error) {
	var count int64
	err := global.DB.Model(&model.Favorite{}).Where("user_id = ? AND product_id = ?", userID, productID).Count(&count).Error
	return count > 0, err
}
