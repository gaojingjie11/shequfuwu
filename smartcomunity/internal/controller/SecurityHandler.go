package controller

import (
	"smartcommunity/internal/model"
	"smartcommunity/internal/service"
	"smartcommunity/pkg/response"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type SecurityHandler struct {
	Service service.SecurityService
}

// ... (existing code) -> Removed, actually I'll just remove the block.

// --- 车位管理 (Admin) ---

// GetAllParking 获取所有车位 (Admin)
func (h *SecurityHandler) GetAllParking(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))
	list, total, err := h.Service.GetAllParking(page, size)
	if err != nil {
		response.Fail(c, "获取失败")
		return
	}
	response.Success(c, gin.H{"list": list, "total": total})
}

// GetParkingStats 获取车位统计
func (h *SecurityHandler) GetParkingStats(c *gin.Context) {
	stats, err := h.Service.GetParkingStats()
	if err != nil {
		response.Fail(c, "获取统计失败")
		return
	}
	response.Success(c, stats)
}

// AssignParking 为用户分配车位
func (h *SecurityHandler) AssignParking(c *gin.Context) {
	var req struct {
		ID       int64  `json:"id"`
		UserID   int64  `json:"user_id"`
		CarPlate string `json:"car_plate"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, "参数错误")
		return
	}

	if err := h.Service.AssignParking(req.ID, req.UserID, req.CarPlate); err != nil {
		response.Fail(c, "操作失败: "+err.Error())
		return
	}
	response.Success(c, nil)
}

// --- 访客接口 ---

// CreateVisitor 登记访客
func (h *SecurityHandler) CreateVisitor(c *gin.Context) {
	userID, _ := c.Get("userID")

	var req struct {
		Name      string `json:"visitor_name"`
		Mobile    string `json:"visitor_phone"`
		Reason    string `json:"reason"`
		VisitTime string `json:"visit_time"` // 前端传字符串 "2024-06-20 14:00:00"
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, "参数错误")
		return
	}

	// 时间解析
	vTime, err := time.ParseInLocation("2006-01-02 15:04:05", req.VisitTime, time.Local)
	if err != nil {
		response.Fail(c, "时间格式错误，需为 YYYY-MM-DD HH:mm:ss")
		return
	}

	visitor := model.Visitor{
		UserID:    userID.(int64),
		Name:      req.Name,
		Mobile:    req.Mobile,
		Reason:    req.Reason,
		VisitTime: vTime,
	}

	if err := h.Service.CreateVisitor(&visitor); err != nil {
		response.Fail(c, "提交失败")
		return
	}
	response.Success(c, nil)
}

// ListVisitor 查看记录
func (h *SecurityHandler) ListVisitor(c *gin.Context) {
	userID, _ := c.Get("userID")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))

	list, total, err := h.Service.GetMyVisitors(userID.(int64), page, size)
	if err != nil {
		response.Fail(c, "获取失败")
		return
	}
	response.Success(c, gin.H{"list": list, "total": total})
}

// --- 车位接口 ---

// MyParking 查看我的车位
func (h *SecurityHandler) MyParking(c *gin.Context) {
	userID, _ := c.Get("userID")
	parking, err := h.Service.GetMyParking(userID.(int64))
	if err != nil {
		// 未绑定车位不应视为接口失败，返回 null 即可
		response.Success(c, nil)
		return
	}
	response.Success(c, parking)
}

// BindCar 绑定车牌
func (h *SecurityHandler) BindCar(c *gin.Context) {
	userID, _ := c.Get("userID")
	var req struct {
		ParkingID int64  `json:"parking_id"`
		CarPlate  string `json:"car_plate"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, "参数错误")
		return
	}

	if err := h.Service.BindCarPlate(userID.(int64), req.ParkingID, req.CarPlate); err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.Success(c, nil)
}

// AuditVisitor 审核访客 (Admin)
func (h *SecurityHandler) AuditVisitor(c *gin.Context) {
	var req struct {
		ID     int64  `json:"id"`
		Status int    `json:"status"` // 1:通过 2:拒绝
		Remark string `json:"remark"` // 审核意见
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, "参数错误")
		return
	}

	if err := h.Service.AuditVisitor(req.ID, req.Status, req.Remark); err != nil {
		response.Fail(c, "操作失败")
		return
	}
	response.Success(c, nil)
}

// ListAllVisitor 管理员查看所有
func (h *SecurityHandler) ListAllVisitor(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))

	list, total, err := h.Service.GetAllVisitors(page, size)
	if err != nil {
		response.Fail(c, "获取失败")
		return
	}
	response.Success(c, gin.H{
		"list":  list,
		"total": total,
	})
}

// CreateParking 创建车位 (Admin)
func (h *SecurityHandler) CreateParking(c *gin.Context) {
	var req struct {
		ParkingNo string `json:"parking_no"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, "参数错误")
		return
	}

	if req.ParkingNo == "" {
		response.Fail(c, "车位号不能为空")
		return
	}

	if err := h.Service.CreateParking(req.ParkingNo); err != nil {
		response.Fail(c, "创建失败: "+err.Error())
		return
	}
	response.Success(c, nil)
}
