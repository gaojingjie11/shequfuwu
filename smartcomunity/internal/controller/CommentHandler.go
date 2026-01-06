package controller

import (
	"smartcommunity/internal/service"
	"smartcommunity/pkg/response"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CommentHandler struct {
	Service service.CommentService
}

func (h *CommentHandler) Create(c *gin.Context) {
	userID, _ := c.Get("userID")
	var req struct {
		ProductID int64  `json:"product_id"`
		Content   string `json:"content"`
		Rating    int    `json:"rating"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, "参数错误")
		return
	}

	if err := h.Service.CreateComment(userID.(int64), req.ProductID, req.Content, req.Rating); err != nil {
		response.Fail(c, "评论失败")
		return
	}
	response.Success(c, nil)
}

func (h *CommentHandler) List(c *gin.Context) {
	productID, _ := strconv.ParseInt(c.Query("product_id"), 10, 64)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))

	list, total, err := h.Service.GetCommentsByProductID(productID, page, size)
	if err != nil {
		response.Fail(c, "获取失败")
		return
	}
	response.Success(c, gin.H{
		"list":  list,
		"total": total,
	})
}
