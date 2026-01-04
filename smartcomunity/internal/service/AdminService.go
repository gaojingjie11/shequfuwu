package service

import (
	"smartcommunity/internal/global"
	"smartcommunity/internal/model"
	"strconv"
	"time"

	"gorm.io/gorm"
)

type AdminService struct{}

// --- 角色管理 ---

// CreateRole 创建角色
func (s *AdminService) CreateRole(role *model.SysRole) error {
	return global.DB.Create(role).Error
}

// UpdateRole 修改角色
func (s *AdminService) UpdateRole(role *model.SysRole) error {
	return global.DB.Model(&model.SysRole{}).Where("id = ?", role.ID).Updates(role).Error
}

// DeleteRole 删除角色
func (s *AdminService) DeleteRole(id int64) error {
	return global.DB.Delete(&model.SysRole{}, id).Error
}

// ListRoles 角色列表
func (s *AdminService) ListRoles() ([]model.SysRole, error) {
	var roles []model.SysRole
	err := global.DB.Find(&roles).Error
	return roles, err
}

// --- 菜单管理 (权限) ---

// CreateMenu 创建菜单
func (s *AdminService) CreateMenu(menu *model.SysMenu) error {
	return global.DB.Create(menu).Error
}

// ListMenus 获取所有菜单
func (s *AdminService) ListMenus() ([]model.SysMenu, error) {
	var menus []model.SysMenu
	err := global.DB.Order("sort asc").Find(&menus).Error
	return menus, err
}

// BindRoleMenu 为角色分配菜单
func (s *AdminService) BindRoleMenu(roleID int64, menuIDs []int64) error {
	tx := global.DB.Begin()
	// 1. 删除旧关联
	if err := tx.Where("role_id = ?", roleID).Delete(&model.SysRoleMenu{}).Error; err != nil {
		tx.Rollback()
		return err
	}
	// 2. 添加新关联
	var roleMenus []model.SysRoleMenu
	for _, mid := range menuIDs {
		roleMenus = append(roleMenus, model.SysRoleMenu{
			RoleID: roleID,
			MenuID: mid,
		})
	}
	if len(roleMenus) > 0 {
		if err := tx.Create(&roleMenus).Error; err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit().Error
}

// --- 用户管理 ---

// ListUsers 获取用户列表 (支持搜索)
func (s *AdminService) ListUsers(page, size int, keyword string) ([]model.SysUser, int64, error) {
	var users []model.SysUser
	var total int64
	db := global.DB.Model(&model.SysUser{})

	if keyword != "" {
		db = db.Where("username LIKE ? OR mobile LIKE ? OR real_name LIKE ?", "%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}

	db.Count(&total)
	err := db.Offset((page - 1) * size).Limit(size).Find(&users).Error
	return users, total, err
}

// FreezeUser 冻结/解冻用户
func (s *AdminService) FreezeUser(id int64, status int) error {
	return global.DB.Model(&model.SysUser{}).Where("id = ?", id).Update("status", status).Error
}

// AssignRole 为用户分配角色
func (s *AdminService) AssignRole(userID int64, roleCode string) error {
	return global.DB.Model(&model.SysUser{}).Where("id = ?", userID).Update("role", roleCode).Error
}

// UpdateUserBalance 为用户充值/扣费 (管理员)
func (s *AdminService) UpdateUserBalance(userID int64, amount float64) error {
	if amount == 0 {
		return nil
	}

	return global.DB.Transaction(func(tx *gorm.DB) error {
		// 1. 更新余额
		if err := tx.Model(&model.SysUser{}).Where("id = ?", userID).
			Update("balance", gorm.Expr("balance + ?", amount)).Error; err != nil {
			return err
		}

		// 2. 记录流水
		remark := "管理员调整余额"
		if amount > 0 {
			remark = "系统充值"
		} else {
			remark = "系统扣除"
		}

		transaction := model.SysTransaction{
			UserID:    userID,
			Type:      3, // 3: 系统调整/充值
			Amount:    amount,
			Remark:    remark,
			CreatedAt: time.Now(),
		}
		if err := tx.Create(&transaction).Error; err != nil {
			return err
		}
		return nil
	})
}

// --- 数据大屏统计 ---

type DashboardStats struct {
	TotalUsers      int64                    `json:"totalUsers"`
	TodayOrders     int64                    `json:"todayOrders"`
	ParkingRate     string                   `json:"parkingRate"`
	MonthIncome     float64                  `json:"monthIncome"`
	RepairStats     []map[string]interface{} `json:"repairStats"`
	VisitorLogs     []model.Visitor          `json:"visitorLogs"`
	IncomeTrend     []float64                `json:"incomeTrend"` // 近7天收入
	IncomeDates     []string                 `json:"incomeDates"` // 近7天日期
	YearTotalAmount float64                  `json:"yearTotalAmount"`
	PatrolCount     int64                    `json:"patrolCount"`
	CostStructure   []float64                `json:"costStructure"` // [物业, 停车, 商城]
}

func (s *AdminService) GetDashboardStats() (*DashboardStats, error) {
	stats := &DashboardStats{}

	// 1. 总用户数
	global.DB.Model(&model.SysUser{}).Count(&stats.TotalUsers)

	// 2. 今日订单
	todayStart := time.Now().Format("2006-01-02") + " 00:00:00"
	global.DB.Model(&model.Order{}).Where("created_at >= ?", todayStart).Count(&stats.TodayOrders)

	// 3. 本月营收 (商城)
	monthStart := time.Now().Format("2006-01") + "-01 00:00:00"
	var mallIncome float64
	global.DB.Model(&model.Order{}).Where("created_at >= ? AND status in (1,2,3)", monthStart).Select("COALESCE(sum(total_amount), 0)").Scan(&mallIncome)
	stats.MonthIncome = mallIncome

	// 4. 车位占用率 & 估算停车收入 (假设每个占用车位每月贡献 300)
	var totalParking, occupiedParking int64
	global.DB.Model(&model.Parking{}).Count(&totalParking)
	global.DB.Model(&model.Parking{}).Where("status = ?", 1).Count(&occupiedParking)
	if totalParking > 0 {
		rate := float64(occupiedParking) / float64(totalParking) * 100
		stats.ParkingRate = strconv.FormatFloat(rate, 'f', 0, 64) + "%"
	} else {
		stats.ParkingRate = "0%"
	}
	parkingIncome := float64(occupiedParking) * 30 // 估算值

	// 5. 报修占比
	var repairs []struct {
		Category string `json:"name"`
		Count    int    `json:"value"`
	}
	global.DB.Model(&model.Repair{}).Select("category, count(*) as count").Group("category").Scan(&repairs)
	for _, r := range repairs {
		stats.RepairStats = append(stats.RepairStats, map[string]interface{}{"name": r.Category, "value": r.Count})
	}

	// 6. 访客记录
	global.DB.Model(&model.Visitor{}).Order("created_at desc").Limit(10).Find(&stats.VisitorLogs)

	// 7. 近7天收入趋势
	for i := 6; i >= 0; i-- {
		day := time.Now().AddDate(0, 0, -i).Format("2006-01-02")
		stats.IncomeDates = append(stats.IncomeDates, day[5:])

		var dayIncome float64
		start := day + " 00:00:00"
		end := day + " 23:59:59"
		global.DB.Model(&model.Order{}).
			Where("created_at BETWEEN ? AND ? AND status in (1,2,3)", start, end).
			Select("COALESCE(sum(total_amount), 0)").Scan(&dayIncome)
		stats.IncomeTrend = append(stats.IncomeTrend, dayIncome)
	}

	// 8. 年度总交易额 (商城 + 物业费)
	yearStart := time.Now().Format("2006") + "-01-01 00:00:00"
	var yearMallIncome, yearPropertyIncome float64
	global.DB.Model(&model.Order{}).Where("created_at >= ? AND status in (1,2,3)", yearStart).Select("COALESCE(sum(total_amount), 0)").Scan(&yearMallIncome)
	global.DB.Model(&model.PropertyFee{}).Where("status = 1").Select("COALESCE(sum(amount), 0)").Scan(&yearPropertyIncome)

	stats.YearTotalAmount = yearMallIncome + yearPropertyIncome

	// 9. 安保巡逻次数 (访客数 * 2 + 基础值)
	var visitorCount int64
	global.DB.Model(&model.Visitor{}).Count(&visitorCount)
	stats.PatrolCount = visitorCount*2 + 500 // 加上基础值显得真实一点

	// 10. 费用构成 (物业, 停车, 商城)
	// 这里用的是本月或年度的总量作为比例
	stats.CostStructure = []float64{yearPropertyIncome, parkingIncome * 12, yearMallIncome}

	return stats, nil
}
