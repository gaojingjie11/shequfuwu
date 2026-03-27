package service

import (
	"errors"
	"fmt"
	"log"
	"smartcommunity/internal/global"
	"smartcommunity/internal/model"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"
)

type AdminService struct{}

func (s *AdminService) CreateRole(role *model.SysRole) error {
	return global.DB.Create(role).Error
}

func (s *AdminService) UpdateRole(role *model.SysRole) error {
	return global.DB.Model(&model.SysRole{}).Where("id = ?", role.ID).Updates(role).Error
}

func (s *AdminService) DeleteRole(id int64) error {
	return global.DB.Delete(&model.SysRole{}, id).Error
}

func (s *AdminService) ListRoles() ([]model.SysRole, error) {
	var roles []model.SysRole
	err := global.DB.Find(&roles).Error
	return roles, err
}

func (s *AdminService) CreateMenu(menu *model.SysMenu) error {
	return global.DB.Create(menu).Error
}

func (s *AdminService) ListMenus() ([]model.SysMenu, error) {
	var menus []model.SysMenu
	err := global.DB.Order("sort asc").Find(&menus).Error
	return menus, err
}

func (s *AdminService) BindRoleMenu(roleID int64, menuIDs []int64) error {
	tx := global.DB.Begin()
	if err := tx.Where("role_id = ?", roleID).Delete(&model.SysRoleMenu{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	roleMenus := make([]model.SysRoleMenu, 0, len(menuIDs))
	for _, menuID := range menuIDs {
		roleMenus = append(roleMenus, model.SysRoleMenu{
			RoleID: roleID,
			MenuID: menuID,
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

func (s *AdminService) ListUsers(page, size int, keyword string) ([]model.SysUser, int64, error) {
	var users []model.SysUser
	var total int64
	db := global.DB.Model(&model.SysUser{})

	if keyword != "" {
		like := "%" + keyword + "%"
		db = db.Where("username LIKE ? OR mobile LIKE ? OR real_name LIKE ?", like, like, like)
	}

	db.Count(&total)
	err := db.Offset((page - 1) * size).Limit(size).Find(&users).Error
	return users, total, err
}

func (s *AdminService) FreezeUser(id int64, status int) error {
	return global.DB.Model(&model.SysUser{}).Where("id = ?", id).Update("status", status).Error
}

func (s *AdminService) AssignRole(userID int64, roleCode string) error {
	return global.DB.Model(&model.SysUser{}).Where("id = ?", userID).Update("role", roleCode).Error
}

func (s *AdminService) UpdateUserBalance(userID int64, amount float64) error {
	if amount == 0 {
		return nil
	}

	return global.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&model.SysUser{}).Where("id = ?", userID).
			Update("balance", gorm.Expr("balance + ?", amount)).Error; err != nil {
			return err
		}

		remark := "Admin adjusted balance"
		if amount > 0 {
			remark = "System balance recharge"
		} else if amount < 0 {
			remark = "System balance deduction"
		}

		transaction := model.SysTransaction{
			UserID:    userID,
			Type:      TransactionTypeTopUp,
			Amount:    amount,
			Remark:    remark,
			CreatedAt: time.Now(),
		}
		return tx.Create(&transaction).Error
	})
}

type DashboardStats struct {
	TotalUsers      int64                    `json:"totalUsers"`
	TodayOrders     int64                    `json:"todayOrders"`
	ParkingRate     string                   `json:"parkingRate"`
	MonthIncome     float64                  `json:"monthIncome"`
	RepairStats     []map[string]interface{} `json:"repairStats"`
	VisitorLogs     []model.Visitor          `json:"visitorLogs"`
	IncomeTrend     []float64                `json:"incomeTrend"`
	IncomeDates     []string                 `json:"incomeDates"`
	YearTotalAmount float64                  `json:"yearTotalAmount"`
	PatrolCount     int64                    `json:"patrolCount"`
	CostStructure   []float64                `json:"costStructure"`
}

func (s *AdminService) GetDashboardStats() (*DashboardStats, error) {
	stats := &DashboardStats{}

	global.DB.Model(&model.SysUser{}).Count(&stats.TotalUsers)

	todayStart := time.Now().Format("2006-01-02") + " 00:00:00"
	global.DB.Model(&model.Order{}).Where("created_at >= ?", todayStart).Count(&stats.TodayOrders)

	monthStart := time.Now().Format("2006-01") + "-01 00:00:00"
	var mallIncome float64
	global.DB.Model(&model.Order{}).
		Where("created_at >= ? AND status in (1,2,3)", monthStart).
		Select("COALESCE(sum(total_amount), 0)").
		Scan(&mallIncome)
	stats.MonthIncome = mallIncome

	var totalParking, occupiedParking int64
	global.DB.Model(&model.Parking{}).Count(&totalParking)
	global.DB.Model(&model.Parking{}).Where("status = ?", 1).Count(&occupiedParking)
	if totalParking > 0 {
		rate := float64(occupiedParking) / float64(totalParking) * 100
		stats.ParkingRate = strconv.FormatFloat(rate, 'f', 0, 64) + "%"
	} else {
		stats.ParkingRate = "0%"
	}
	parkingIncome := float64(occupiedParking) * 30

	var repairs []struct {
		Category string `json:"name"`
		Count    int    `json:"value"`
	}
	global.DB.Model(&model.Repair{}).Select("category, count(*) as count").Group("category").Scan(&repairs)
	for _, repair := range repairs {
		stats.RepairStats = append(stats.RepairStats, map[string]interface{}{
			"name":  repair.Category,
			"value": repair.Count,
		})
	}

	global.DB.Model(&model.Visitor{}).Order("created_at desc").Limit(10).Find(&stats.VisitorLogs)

	for i := 6; i >= 0; i-- {
		day := time.Now().AddDate(0, 0, -i).Format("2006-01-02")
		stats.IncomeDates = append(stats.IncomeDates, day[5:])

		var dayIncome float64
		start := day + " 00:00:00"
		end := day + " 23:59:59"
		global.DB.Model(&model.Order{}).
			Where("created_at BETWEEN ? AND ? AND status in (1,2,3)", start, end).
			Select("COALESCE(sum(total_amount), 0)").
			Scan(&dayIncome)
		stats.IncomeTrend = append(stats.IncomeTrend, dayIncome)
	}

	yearStart := time.Now().Format("2006") + "-01-01 00:00:00"
	var yearMallIncome, yearPropertyIncome float64
	global.DB.Model(&model.Order{}).
		Where("created_at >= ? AND status in (1,2,3)", yearStart).
		Select("COALESCE(sum(total_amount), 0)").
		Scan(&yearMallIncome)
	global.DB.Model(&model.PropertyFee{}).
		Where("status = 1").
		Select("COALESCE(sum(amount), 0)").
		Scan(&yearPropertyIncome)
	stats.YearTotalAmount = yearMallIncome + yearPropertyIncome

	var visitorCount int64
	global.DB.Model(&model.Visitor{}).Count(&visitorCount)
	stats.PatrolCount = visitorCount*2 + 500
	stats.CostStructure = []float64{yearPropertyIncome, parkingIncome * 12, yearMallIncome}

	return stats, nil
}

func (s *AdminService) GenerateAIReport(operatorID int64) (*model.AIReport, error) {
	sevenDaysAgo := time.Now().AddDate(0, 0, -6).Format("2006-01-02") + " 00:00:00"

	report := &model.AIReport{
		GeneratedBy: operatorID,
	}

	global.DB.Model(&model.Repair{}).Where("created_at >= ?", sevenDaysAgo).Count(&report.RepairNewCount)
	global.DB.Model(&model.Repair{}).Where("status <> ?", 2).Count(&report.RepairPendingCount)
	global.DB.Model(&model.Visitor{}).Where("created_at >= ?", sevenDaysAgo).Count(&report.VisitorNewCount)
	global.DB.Model(&model.PropertyFee{}).Where("status = 1 AND pay_time >= ?", sevenDaysAgo).Count(&report.PropertyPaidCount)
	global.DB.Model(&model.PropertyFee{}).
		Where("status = 1 AND pay_time >= ?", sevenDaysAgo).
		Select("COALESCE(sum(amount), 0)").
		Scan(&report.PropertyPaidAmount)

	prompt := fmt.Sprintf(
		"你是一个高级社区物业经理。以下是本社区近7天的数据：报修新增%d条，未处理%d条；访客新增%d条；物业费缴费%d笔，收缴%.2f元。请用 Markdown 生成一份专业的数据分析报告，包含：1）核心数据概览；2）管理风险；3）可执行建议。语言简洁，条理清晰。",
		report.RepairNewCount,
		report.RepairPendingCount,
		report.VisitorNewCount,
		report.PropertyPaidCount,
		report.PropertyPaidAmount,
	)

	reportText, err := (&AIService{}).GenerateCommunityReport(prompt)
	if err != nil {
		log.Printf("generate AI report failed: %v", err)
		reportText = buildFallbackCommunityReport(report)
	}
	report.Report = normalizeAIReportMarkdown(reportText)
	if report.Report == "" {
		report.Report = normalizeAIReportMarkdown(buildFallbackCommunityReport(report))
	}
	report.ReportSummary = buildAIReportSummaryV2(report.Report)

	if err := global.DB.Create(report).Error; err != nil {
		return nil, err
	}
	return report, nil
}

func (s *AdminService) GenerateDailyAIReportIfNeeded(location *time.Location) error {
	if location == nil {
		location = time.Local
	}

	now := time.Now().In(location)
	dayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, location)
	dayEnd := dayStart.Add(24 * time.Hour)

	var count int64
	if err := global.DB.Model(&model.AIReport{}).
		Where("generated_by = ? AND created_at >= ? AND created_at < ?", int64(0), dayStart, dayEnd).
		Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return nil
	}

	_, err := s.GenerateAIReport(0)
	return err
}

func (s *AdminService) GetLatestAIReport() (*model.AIReport, error) {
	var report model.AIReport
	if err := global.DB.Order("id desc").First(&report).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &report, nil
}

func (s *AdminService) GetAIReport(refresh bool, operatorID int64) (*model.AIReport, error) {
	if refresh {
		return s.GenerateAIReport(operatorID)
	}

	report, err := s.GetLatestAIReport()
	if err != nil {
		return nil, err
	}
	if report != nil {
		return report, nil
	}

	return s.GenerateAIReport(operatorID)
}

func (s *AdminService) ListAIReports(page, size int) ([]model.AIReport, int64, error) {
	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 10
	}

	var list []model.AIReport
	var total int64
	db := global.DB.Model(&model.AIReport{})
	db.Count(&total)

	offset := (page - 1) * size
	err := db.Select(
		"id",
		"repair_new_count",
		"repair_pending_count",
		"visitor_new_count",
		"property_paid_count",
		"property_paid_amount",
		"report_summary",
		"created_at",
	).Order("id desc").Offset(offset).Limit(size).Find(&list).Error
	if err == nil {
		for i := range list {
			list[i].ReportSummary = normalizeAIReportSummaryText(list[i].ReportSummary)
		}
	}
	return list, total, err
}

func (s *AdminService) GetAIReportDetail(id int64) (*model.AIReport, error) {
	var report model.AIReport
	if err := global.DB.First(&report, id).Error; err != nil {
		return nil, err
	}
	return &report, nil
}

func buildAIReportSummary(report string) string {
	lines := strings.Split(report, "\n")
	for _, line := range lines {
		text := strings.TrimSpace(line)
		text = strings.TrimLeft(text, "#-*0123456789. ")
		text = strings.TrimSpace(text)
		if text == "" {
			continue
		}
		runes := []rune(text)
		if len(runes) > 42 {
			return string(runes[:42]) + "..."
		}
		return text
	}
	return "AI 社区数据分析报告"
}
