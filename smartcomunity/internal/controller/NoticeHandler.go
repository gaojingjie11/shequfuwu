package controller

import (
	"github.com/gin-gonic/gin"
	"smartcommunity/internal/model"
	"smartcommunity/internal/service"
	"smartcommunity/pkg/response"
	"strconv"
)

type NoticeHandler struct {
	Service service.NoticeService
}

// List 公告列表 (支持分页)
func (h *NoticeHandler) List(c *gin.Context) {
	// 如果传入 page 参数，则走分页逻辑
	if c.Query("page") != "" {
		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))

		list, total, err := h.Service.GetPageList(page, size)
		if err != nil {
			response.Fail(c, "获取失败")
			return
		}
		response.Success(c, gin.H{"list": list, "total": total})
		return
	}

	// 否则走默认 Limit 逻辑 (首页)
	list, err := h.Service.GetList(10) // 默认取最新 10 条
	if err != nil {
		response.Fail(c, "获取失败")
		return
	}
	response.Success(c, list)
}

// Detail 公告详情
func (h *NoticeHandler) Detail(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)

	notice, err := h.Service.GetDetail(id)
	if err != nil {
		response.Fail(c, "公告不存在")
		return
	}
	response.Success(c, notice)
}

// Create 发布公告 (Admin)
func (h *NoticeHandler) Create(c *gin.Context) {
	var req model.Notice
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, "参数错误")
		return
	}
	// 默认发布人
	if req.Publisher == "" {
		req.Publisher = "物业管理员"
	}

	if err := h.Service.Create(&req); err != nil {
		response.Fail(c, "发布失败")
		return
	}
	response.Success(c, req)
}

func (h *NoticeHandler) Delete(c *gin.Context) {
	// 1. 获取路径参数 "id" (字符串类型)
	idStr := c.Param("id")

	// 2. 将字符串转换为 int64
	// base: 10 (十进制), bitSize: 64 (64位整数)
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.Fail(c, "ID格式错误")
		return
	}

	// 3. 调用 Service 删除
	if err := h.Service.Delete(id); err != nil {
		response.Fail(c, "删除失败: "+err.Error())
		return
	}

	response.Success(c, nil)
}

// Read 标记已读
func (h *NoticeHandler) Read(c *gin.Context) {
	userID, _ := c.Get("userID")
	idStr := c.Param("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)

	if err := h.Service.MarkRead(userID.(int64), id); err != nil {
		response.Fail(c, "操作失败")
		return
	}
	response.Success(c, nil)
}
