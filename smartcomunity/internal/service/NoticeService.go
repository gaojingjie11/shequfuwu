package service

import (
	"smartcommunity/internal/global"
	"smartcommunity/internal/model"
	"time"
)

type NoticeService struct{}

// GetList 获取公告列表 (Limit限制条数，通常首页取前5条)
func (s *NoticeService) GetList(limit int) ([]model.Notice, error) {
	var list []model.Notice
	// 按时间倒序排列
	err := global.DB.Order("created_at desc").Limit(limit).Find(&list).Error
	return list, err
}

// GetPageList 分页获取公告
func (s *NoticeService) GetPageList(page, size int) ([]model.Notice, int64, error) {
	var list []model.Notice
	var total int64
	offset := (page - 1) * size

	db := global.DB.Model(&model.Notice{})

	err := db.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = db.Order("created_at desc").Offset(offset).Limit(size).Find(&list).Error
	return list, total, err
}

// GetDetail 获取详情并增加浏览量 (可选功能)
func (s *NoticeService) GetDetail(id int64) (*model.Notice, error) {
	var notice model.Notice
	if err := global.DB.First(&notice, id).Error; err != nil {
		return nil, err
	}

	// 浏览量 +1 (这里为了性能可以不加锁，或者用 Redis 计数)
	// 简单起见直接 SQL 更新
	global.DB.Model(&notice).UpdateColumn("view_count", notice.ViewCount+1)

	return &notice, nil
}

// Create 发布公告
func (s *NoticeService) Create(notice *model.Notice) error {
	return global.DB.Create(notice).Error
}

// Delete 删除公告
func (s *NoticeService) Delete(id int64) error {
	return global.DB.Delete(&model.Notice{}, id).Error
}

// MarkRead 标记已读
func (s *NoticeService) MarkRead(userID, noticeID int64) error {
	// 检查是否已读
	var count int64
	global.DB.Model(&model.NoticeRead{}).Where("user_id = ? AND notice_id = ?", userID, noticeID).Count(&count)
	if count > 0 {
		return nil
	}
	read := model.NoticeRead{
		UserID:   userID,
		NoticeID: noticeID,
		ReadTime: time.Now(),
	}
	return global.DB.Create(&read).Error
}
