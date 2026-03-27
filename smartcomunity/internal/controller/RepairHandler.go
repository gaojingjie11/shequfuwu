package controller

import (
	"smartcommunity/internal/model"
	"smartcommunity/internal/service"
	"smartcommunity/pkg/response"
	"strconv"

	"github.com/gin-gonic/gin"
)

type RepairHandler struct {
	Service service.RepairService
}

// Create 提交接口
func (h *RepairHandler) Create(c *gin.Context) {
	userID, _ := c.Get("userID")

	var req struct {
		Type     int    `json:"type"`     // 1报修 2投诉
		Category string `json:"category"` // 分类
		Content  string `json:"content"`  // 内容
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, "参数错误")
		return
	}

	repair := model.Repair{
		UserID:   userID.(int64),
		Type:     req.Type,
		Category: req.Category,
		Content:  req.Content,
	}

	if err := h.Service.Create(&repair); err != nil {
		response.Fail(c, "提交失败")
		return
	}
	response.Success(c, nil)
}

// List 我的记录接口
func (h *RepairHandler) List(c *gin.Context) {
	userID, _ := c.Get("userID")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))

	list, total, err := h.Service.GetUserList(userID.(int64), page, size)
	if err != nil {
		response.Fail(c, "获取失败")
		return
	}
	response.Success(c, gin.H{"list": list, "total": total})
}

// Process 处理报修 (Admin)
func (h *RepairHandler) Process(c *gin.Context) {
	var req struct {
		ID       int64  `json:"id"`
		Status   int    `json:"status"`   // 1或2
		Feedback string `json:"feedback"` // 例如 "维修工李四已上门修复"
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, "参数错误")
		return
	}

	if err := h.Service.UpdateStatus(req.ID, req.Status, req.Feedback); err != nil {
		response.Fail(c, "操作失败")
		return
	}
	response.Success(c, nil)
}

// ListAll 管理员列表 (分页)
func (h *RepairHandler) ListAll(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))

	list, total, err := h.Service.GetPageList(page, size)
	if err != nil {
		response.Fail(c, "获取失败")
		return
	}
	response.Success(c, gin.H{"list": list, "total": total})
}
