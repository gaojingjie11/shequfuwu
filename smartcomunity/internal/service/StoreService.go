package service

import (
	"smartcommunity/internal/global"
	"smartcommunity/internal/model"
)

type StoreService struct{}

// GetAll 获取所有门店
func (s *StoreService) GetAll() ([]model.Store, error) {
	var list []model.Store
	// 简单查询所有，也可以按 region 筛选
	err := global.DB.Find(&list).Error
	return list, err
}

// CreateStore 创建门店
func (s *StoreService) CreateStore(store *model.Store) error {
	return global.DB.Create(store).Error
}

// UpdateStore 修改门店
func (s *StoreService) UpdateStore(store *model.Store) error {
	return global.DB.Model(&model.Store{}).Where("id = ?", store.ID).Updates(store).Error
}

// DeleteStore 删除门店
func (s *StoreService) DeleteStore(id int64) error {
	return global.DB.Delete(&model.Store{}, id).Error
}

// BindProduct 为门店分配商品库存
func (s *StoreService) BindProduct(storeID, productID int64, stock int) error {
	var sp model.StoreProduct
	// 检查是否已存在
	if err := global.DB.Where("store_id = ? AND product_id = ?", storeID, productID).First(&sp).Error; err == nil {
		// 已存在则更新库存
		return global.DB.Model(&sp).Update("stock", stock).Error
	}
	// 不存在则创建
	sp = model.StoreProduct{
		StoreID:   storeID,
		ProductID: productID,
		Stock:     stock,
	}
	return global.DB.Create(&sp).Error
}

// GetStoreStock 获取门店某商品库存
func (s *StoreService) GetStoreStock(storeID, productID int64) (int, error) {
	var sp model.StoreProduct
	err := global.DB.Where("store_id = ? AND product_id = ?", storeID, productID).First(&sp).Error
	return sp.Stock, err
}
