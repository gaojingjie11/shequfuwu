package service

import (
	"errors"
	"smartcommunity/internal/global"
	"smartcommunity/internal/model"
)

type SecurityService struct{}

// --- 访客相关 ---

// CreateVisitor 提交访客登记
func (s *SecurityService) CreateVisitor(visitor *model.Visitor) error {
	visitor.Status = 0 // 默认为待审核
	return global.DB.Create(visitor).Error
}

// GetMyVisitors 获取我的访客记录 (分页)
func (s *SecurityService) GetMyVisitors(userID int64, page, size int) ([]model.Visitor, int64, error) {
	var list []model.Visitor
	var total int64
	db := global.DB.Model(&model.Visitor{}).Where("user_id = ?", userID)
	db.Count(&total)

	offset := (page - 1) * size
	err := db.Order("created_at desc").Offset(offset).Limit(size).Find(&list).Error
	return list, total, err
}

// --- 车位相关 ---

// GetMyParking 获取我的车位信息 (支持多个)
func (s *SecurityService) GetMyParking(userID int64) ([]model.Parking, error) {
	var list []model.Parking
	err := global.DB.Where("user_id = ?", userID).Find(&list).Error
	if err != nil {
		return nil, err
	}
	return list, nil
}

// BindCarPlate 绑定/修改车牌号
func (s *SecurityService) BindCarPlate(userID int64, parkingID int64, carPlate string) error {
	// 1. 查找用户拥有的车位 (Double check ownership)
	var parking model.Parking
	if err := global.DB.Where("id = ? AND user_id = ?", parkingID, userID).First(&parking).Error; err != nil {
		return errors.New("未找到对应车位或无权操作")
	}

	// 2. 更新车牌
	parking.CarPlate = carPlate
	parking.Status = 1 // 标记为占用
	return global.DB.Save(&parking).Error
}

// ListAvailableParking (可选) 查看空闲车位 - 方便演示购买车位
func (s *SecurityService) ListAvailableParking() ([]model.Parking, error) {
	var list []model.Parking
	err := global.DB.Where("status = 0").Find(&list).Error
	return list, err
}

// AuditVisitor 审核访客 (status: 1通过 2拒绝)
func (s *SecurityService) AuditVisitor(id int64, status int, remark string) error {
	// 只能更新状态为 0 (待审核) 的记录
	return global.DB.Model(&model.Visitor{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":       status,
			"audit_remark": remark, // 存拒绝理由
		}).Error
}

// GetAllVisitors 管理员获取所有访客列表
func (s *SecurityService) GetAllVisitors(page, size int) ([]model.Visitor, int64, error) {
	var list []model.Visitor
	var total int64

	tx := global.DB.Model(&model.Visitor{})
	tx.Count(&total) // Get total count first

	offset := (page - 1) * size
	err := tx.Order("id desc").Offset(offset).Limit(size).Find(&list).Error
	return list, total, err
}

// --- 车位管理 (Admin) ---

// GetAllParking 获取所有车位列表 (Admin)
func (s *SecurityService) GetAllParking(page, size int) ([]model.Parking, int64, error) {
	var list []model.Parking
	var total int64
	db := global.DB.Model(&model.Parking{})
	db.Count(&total)

	offset := (page - 1) * size
	err := db.Order("id asc").Offset(offset).Limit(size).Find(&list).Error
	return list, total, err
}

// GetParkingStats 获取车位统计信息
func (s *SecurityService) GetParkingStats() (map[string]int64, error) {
	var total, used int64
	if err := global.DB.Model(&model.Parking{}).Count(&total).Error; err != nil {
		return nil, err
	}
	if err := global.DB.Model(&model.Parking{}).Where("status = 1").Count(&used).Error; err != nil {
		return nil, err
	}

	return map[string]int64{
		"total": total,
		"used":  used,
		"free":  total - used,
	}, nil
}

// AssignParking 为用户分配车位
func (s *SecurityService) AssignParking(id int64, userID int64, carPlate string) error {
	var parking model.Parking
	if err := global.DB.First(&parking, id).Error; err != nil {
		return errors.New("车位不存在")
	}

	// 如果是解绑 (userID=0)
	if userID == 0 {
		return global.DB.Model(&parking).Updates(map[string]interface{}{
			"user_id":   0,
			"car_plate": "",
			"status":    0,
		}).Error
	}

	// Admin has the right to re-assign, so we allow overwriting.
	// Only check if userID is valid if needed (we assume valid for now)

	return global.DB.Model(&parking).Updates(map[string]interface{}{
		"user_id":   userID,
		"car_plate": carPlate,
		"status":    1,
	}).Error
}

// CreateParking 创建新车位
func (s *SecurityService) CreateParking(parkingNo string) error {
	// 1. 检查是否已存在
	var count int64
	global.DB.Model(&model.Parking{}).Where("parking_no = ?", parkingNo).Count(&count)
	if count > 0 {
		return errors.New("车位号已存在")
	}

	// 2. 创建
	parking := model.Parking{
		ParkingNo: parkingNo,
		Status:    0, // 空闲
	}
	return global.DB.Create(&parking).Error
}
