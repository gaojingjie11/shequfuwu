package controller

import (
	"smartcommunity/internal/service"
	"smartcommunity/pkg/response"

	"github.com/gin-gonic/gin"
)

type AIHandler struct {
	Service service.AIService
}

func (h *AIHandler) Send(c *gin.Context) {
	var req struct {
		Content string `json:"content"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, "参数错误")
		return
	}

	if req.Content == "" {
		response.Fail(c, "内容不能为空")
		return
	}

	reply, err := h.Service.Chat(req.Content)
	if err != nil {
		response.Fail(c, err.Error())
		return
	}

	response.Success(c, gin.H{
		"reply": reply,
	})
}
