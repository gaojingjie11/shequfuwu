package controller

import (
	"strconv"

	"smartcommunity/internal/model"
	"smartcommunity/internal/service"
	"smartcommunity/pkg/response"

	"github.com/gin-gonic/gin"
)

type AdminHandler struct {
	Service service.AdminService
}

func (h *AdminHandler) CreateRole(c *gin.Context) {
	var req model.SysRole
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, "invalid request parameters")
		return
	}
	if err := h.Service.CreateRole(&req); err != nil {
		response.Fail(c, "create role failed: "+err.Error())
		return
	}
	response.Success(c, req)
}

func (h *AdminHandler) ListRoles(c *gin.Context) {
	list, err := h.Service.ListRoles()
	if err != nil {
		response.Fail(c, "failed to fetch roles")
		return
	}
	response.Success(c, list)
}

func (h *AdminHandler) CreateMenu(c *gin.Context) {
	var req model.SysMenu
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, "invalid request parameters")
		return
	}
	if err := h.Service.CreateMenu(&req); err != nil {
		response.Fail(c, "create menu failed")
		return
	}
	response.Success(c, req)
}

func (h *AdminHandler) ListMenus(c *gin.Context) {
	list, err := h.Service.ListMenus()
	if err != nil {
		response.Fail(c, "failed to fetch menus")
		return
	}
	response.Success(c, list)
}

func (h *AdminHandler) BindRoleMenu(c *gin.Context) {
	var req struct {
		RoleID  int64   `json:"role_id"`
		MenuIDs []int64 `json:"menu_ids"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, "invalid request parameters")
		return
	}
	if err := h.Service.BindRoleMenu(req.RoleID, req.MenuIDs); err != nil {
		response.Fail(c, "bind role menu failed")
		return
	}
	response.Success(c, nil)
}

func (h *AdminHandler) ListUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))
	keyword := c.Query("keyword")

	list, total, err := h.Service.ListUsers(page, size, keyword)
	if err != nil {
		response.Fail(c, "failed to fetch users")
		return
	}
	response.Success(c, gin.H{"list": list, "total": total})
}

func (h *AdminHandler) FreezeUser(c *gin.Context) {
	var req struct {
		ID     int64 `json:"id"`
		Status int   `json:"status"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, "invalid request parameters")
		return
	}
	if err := h.Service.FreezeUser(req.ID, req.Status); err != nil {
		response.Fail(c, "update user status failed")
		return
	}
	response.Success(c, nil)
}

func (h *AdminHandler) AssignRole(c *gin.Context) {
	var req struct {
		UserID   int64  `json:"user_id"`
		RoleCode string `json:"role_code"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, "invalid request parameters")
		return
	}
	if err := h.Service.AssignRole(req.UserID, req.RoleCode); err != nil {
		response.Fail(c, "assign role failed: "+err.Error())
		return
	}
	response.Success(c, nil)
}

func (h *AdminHandler) UpdateUserBalance(c *gin.Context) {
	var req struct {
		UserID int64   `json:"user_id"`
		Amount float64 `json:"amount"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, "invalid request parameters")
		return
	}
	if err := h.Service.UpdateUserBalance(req.UserID, req.Amount); err != nil {
		response.Fail(c, "update user balance failed: "+err.Error())
		return
	}
	response.Success(c, nil)
}

func (h *AdminHandler) GetDashboardStats(c *gin.Context) {
	stats, err := h.Service.GetDashboardStats()
	if err != nil {
		response.Fail(c, "failed to fetch dashboard stats")
		return
	}
	response.Success(c, stats)
}

func (h *AdminHandler) GetAIReport(c *gin.Context) {
	userID, _ := c.Get("userID")
	operatorID, _ := userID.(int64)
	refresh := c.DefaultQuery("refresh", "0") == "1"

	report, err := h.Service.GetAIReport(refresh, operatorID)
	if err != nil {
		response.Fail(c, "failed to generate AI report: "+err.Error())
		return
	}
	response.Success(c, report)
}

func (h *AdminHandler) GenerateAIReport(c *gin.Context) {
	userID, _ := c.Get("userID")
	operatorID, _ := userID.(int64)

	report, err := h.Service.GenerateAIReport(operatorID)
	if err != nil {
		response.Fail(c, "failed to generate AI report: "+err.Error())
		return
	}
	response.Success(c, report)
}

func (h *AdminHandler) ListAIReports(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))

	list, total, err := h.Service.ListAIReports(page, size)
	if err != nil {
		response.Fail(c, "failed to fetch AI report list")
		return
	}
	response.Success(c, gin.H{
		"list":  list,
		"total": total,
	})
}

func (h *AdminHandler) GetAIReportDetail(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		response.Fail(c, "invalid report id")
		return
	}

	report, err := h.Service.GetAIReportDetail(id)
	if err != nil {
		response.Fail(c, "failed to fetch AI report detail")
		return
	}
	response.Success(c, report)
}
