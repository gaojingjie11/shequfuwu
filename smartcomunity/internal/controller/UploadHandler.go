package controller

import (
	"smartcommunity/internal/service"
	"smartcommunity/pkg/response"

	"github.com/gin-gonic/gin"
)

type UploadHandler struct {
	Service service.StorageService
}

func (h *UploadHandler) UploadFile(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		response.Fail(c, "please select a file")
		return
	}

	url, key, err := h.Service.UploadMultipartFile(file, "common")
	if err != nil {
		response.Fail(c, "upload failed: "+err.Error())
		return
	}

	response.Success(c, gin.H{
		"url": url,
		"key": key,
	})
}
