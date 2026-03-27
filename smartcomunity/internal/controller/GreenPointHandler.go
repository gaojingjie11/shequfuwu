package controller

import (
	"strconv"

	"smartcommunity/internal/service"
	"smartcommunity/pkg/response"

	"github.com/gin-gonic/gin"
)

type GreenPointHandler struct {
	Service service.GreenPointService
}

func (h *GreenPointHandler) UploadGarbage(c *gin.Context) {
	userID, _ := c.Get("userID")
	file, err := c.FormFile("file")
	if err != nil {
		response.Fail(c, "please upload an image file")
		return
	}

	result, err := h.Service.UploadGarbage(userID.(int64), file)
	if err != nil {
		response.Fail(c, "garbage recognition failed: "+err.Error())
		return
	}

	response.Success(c, result)
}

func (h *GreenPointHandler) Leaderboard(c *gin.Context) {
	limit, _ := strconv.ParseInt(c.DefaultQuery("limit", "10"), 10, 64)
	list, err := h.Service.GetLeaderboard(limit)
	if err != nil {
		response.Fail(c, "failed to fetch leaderboard")
		return
	}

	response.Success(c, gin.H{"list": list})
}
