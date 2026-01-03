package service

import (
	"smartcommunity/internal/global"
	"smartcommunity/internal/model"
)

type MarketingService struct{}

// CreatePromotion 创建/修改活动
func (s *MarketingService) CreatePromotion(p *model.Promotion) error {
	if p.ID > 0 {
		return global.DB.Model(&model.Promotion{}).Where("id = ?", p.ID).Updates(p).Error
	}
	return global.DB.Create(p).Error
}

// ListPromotions 活动列表 (Admin)
func (s *MarketingService) ListPromotions(page, size int) ([]model.Promotion, int64, error) {
	var list []model.Promotion
	var total int64
	db := global.DB.Model(&model.Promotion{})
	db.Count(&total)
	err := db.Order("id desc").Offset((page - 1) * size).Limit(size).Find(&list).Error
	return list, total, err
}

// DeletePromotion 删除活动
func (s *MarketingService) DeletePromotion(id int64) error {
	return global.DB.Delete(&model.Promotion{}, id).Error
}
