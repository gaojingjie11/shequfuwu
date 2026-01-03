package controller

import (
	"smartcommunity/internal/service"
	"smartcommunity/pkg/response"
	"strconv"

	"github.com/gin-gonic/gin"
)

type FavoriteHandler struct {
	Service service.FavoriteService
}

// Add 添加收藏
func (h *FavoriteHandler) Add(c *gin.Context) {
	userID, _ := c.Get("userID")
	var req struct {
		ProductID int64 `json:"product_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, "参数错误")
		return
	}

	if err := h.Service.Add(userID.(int64), req.ProductID); err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.Success(c, nil)
}

// Delete 取消收藏
func (h *FavoriteHandler) Delete(c *gin.Context) {
	userID, _ := c.Get("userID")
	var req struct {
		ProductID int64 `json:"product_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, "参数错误")
		return
	}

	if err := h.Service.Delete(userID.(int64), req.ProductID); err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.Success(c, nil)
}

// List 列表
func (h *FavoriteHandler) List(c *gin.Context) {
	userID, _ := c.Get("userID")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))

	list, total, err := h.Service.GetList(userID.(int64), page, size)
	if err != nil {
		response.Fail(c, "获取失败")
		return
	}
	response.Success(c, gin.H{
		"list":  list,
		"total": total,
	})
}

// Check 检查是否收藏
func (h *FavoriteHandler) Check(c *gin.Context) {
	userID, _ := c.Get("userID")
	productID, _ := strconv.ParseInt(c.Query("product_id"), 10, 64)

	isFav, err := h.Service.Check(userID.(int64), productID)
	if err != nil {
		response.Fail(c, "查询失败")
		return
	}
	response.Success(c, gin.H{"is_favorite": isFav})
}
