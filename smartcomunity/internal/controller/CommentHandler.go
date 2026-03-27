package controller

import (
	"log"
	"smartcommunity/internal/service"
	"smartcommunity/pkg/response"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CommentHandler struct {
	Service service.CommentService
}

func (h *CommentHandler) Create(c *gin.Context) {
	userID, ok := c.Get("userID")
	if !ok {
		response.FailWithCode(c, 401, "请先登录")
		return
	}
	uid, ok := userID.(int64)
	if !ok {
		response.FailWithCode(c, 401, "登录状态无效")
		return
	}

	var req struct {
		ProductID int64  `json:"product_id"`
		Content   string `json:"content"`
		Rating    int    `json:"rating"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, "参数错误")
		return
	}

	if err := h.Service.CreateComment(uid, req.ProductID, req.Content, req.Rating); err != nil {
		log.Printf("create comment failed: user_id=%d product_id=%d rating=%d err=%v", uid, req.ProductID, req.Rating, err)
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
