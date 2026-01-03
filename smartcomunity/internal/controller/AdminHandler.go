package controller

import (
	"smartcommunity/internal/model"
	"smartcommunity/internal/service"
	"smartcommunity/pkg/response"
	"strconv"

	"github.com/gin-gonic/gin"
)

type AdminHandler struct {
	Service service.AdminService
}

// --- 角色管理 ---

func (h *AdminHandler) CreateRole(c *gin.Context) {
	var req model.SysRole
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, "参数错误")
		return
	}
	if err := h.Service.CreateRole(&req); err != nil {
		response.Fail(c, "创建失败: "+err.Error())
		return
	}
	response.Success(c, req)
}

func (h *AdminHandler) ListRoles(c *gin.Context) {
	list, err := h.Service.ListRoles()
	if err != nil {
		response.Fail(c, "获取失败")
		return
	}
	response.Success(c, list)
}

func (h *AdminHandler) CreateMenu(c *gin.Context) {
	var req model.SysMenu
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, "参数错误")
		return
	}
	if err := h.Service.CreateMenu(&req); err != nil {
		response.Fail(c, "创建失败")
		return
	}
	response.Success(c, req)
}

func (h *AdminHandler) ListMenus(c *gin.Context) {
	list, err := h.Service.ListMenus()
	if err != nil {
		response.Fail(c, "获取失败")
		return
	}
	response.Success(c, list)
}

// BindRoleMenu 为角色绑定菜单
func (h *AdminHandler) BindRoleMenu(c *gin.Context) {
	var req struct {
		RoleID  int64   `json:"role_id"`
		MenuIDs []int64 `json:"menu_ids"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, "参数错误")
		return
	}
	if err := h.Service.BindRoleMenu(req.RoleID, req.MenuIDs); err != nil {
		response.Fail(c, "绑定失败")
		return
	}
	response.Success(c, nil)
}

// --- 用户管理 ---

func (h *AdminHandler) ListUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))
	keyword := c.Query("keyword")

	list, total, err := h.Service.ListUsers(page, size, keyword)
	if err != nil {
		response.Fail(c, "获取失败")
		return
	}
	response.Success(c, gin.H{"list": list, "total": total})
}

func (h *AdminHandler) FreezeUser(c *gin.Context) {
	var req struct {
		ID     int64 `json:"id"`
		Status int   `json:"status"` // 1正常 0冻结
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, "参数错误")
		return
	}
	if err := h.Service.FreezeUser(req.ID, req.Status); err != nil {
		response.Fail(c, "操作失败")
		return
	}
	response.Success(c, nil)
}

// AssignRole 为用户分配角色
func (h *AdminHandler) AssignRole(c *gin.Context) {
	var req struct {
		UserID   int64  `json:"user_id"`
		RoleCode string `json:"role_code"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, "参数错误")
		return
	}

	if err := h.Service.AssignRole(req.UserID, req.RoleCode); err != nil {
		response.Fail(c, "操作失败: "+err.Error())
		return
	}
	response.Success(c, nil)
}

// UpdateUserBalance 修改用户余额 (充值/扣费)
func (h *AdminHandler) UpdateUserBalance(c *gin.Context) {
	var req struct {
		UserID int64   `json:"user_id"`
		Amount float64 `json:"amount"` // 正数充值，负数扣费
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, "参数错误")
		return
	}

	if err := h.Service.UpdateUserBalance(req.UserID, req.Amount); err != nil {
		response.Fail(c, "操作失败: "+err.Error())
		return
	}
	response.Success(c, nil)
}

// GetDashboardStats 数据大屏统计
func (h *AdminHandler) GetDashboardStats(c *gin.Context) {
	stats, err := h.Service.GetDashboardStats()
	if err != nil {
		response.Fail(c, "获取统计失败")
		return
	}
	response.Success(c, stats)
}
