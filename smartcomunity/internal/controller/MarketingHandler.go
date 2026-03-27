package controller

import (
	"smartcommunity/internal/model"
	"smartcommunity/internal/service"
	"smartcommunity/pkg/response"
	"strconv"

	"github.com/gin-gonic/gin"
)

type MarketingHandler struct {
	Service service.MarketingService
}

func (h *MarketingHandler) Create(c *gin.Context) {
	var req model.Promotion
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, "参数错误")
		return
	}
	if err := h.Service.CreatePromotion(&req); err != nil {
		response.Fail(c, "操作失败")
		return
	}
	response.Success(c, req)
}

func (h *MarketingHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))
	list, total, err := h.Service.ListPromotions(page, size)
	if err != nil {
		response.Fail(c, "获取失败")
		return
	}
	response.Success(c, gin.H{"list": list, "total": total})
}

func (h *MarketingHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	if err := h.Service.DeletePromotion(id); err != nil {
		response.Fail(c, "删除失败")
		return
	}
	response.Success(c, nil)
}
