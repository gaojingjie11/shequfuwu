package controller

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
	"path/filepath"
	"smartcommunity/internal/config"
	"smartcommunity/internal/global"
	"smartcommunity/pkg/response"
	"time"
)

type UploadHandler struct{}

// UploadFile 上传文件到 MinIO
func (h *UploadHandler) UploadFile(c *gin.Context) {
	// 1. 获取上传的文件
	file, err := c.FormFile("file")
	if err != nil {
		response.Fail(c, "请选择文件")
		return
	}

	// 2. 生成唯一文件名 (时间戳 + 随机数 + 后缀)
	ext := filepath.Ext(file.Filename)
	objectName := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
	// 如果想分日期存储，可以拼路径： fmt.Sprintf("%s/%d%s", time.Now().Format("20060102"), time.Now().UnixNano(), ext)

	// 3. 打开文件流
	src, err := file.Open()
	if err != nil {
		response.Fail(c, "文件读取失败")
		return
	}
	defer src.Close()

	// 4. 上传到 MinIO
	ctx := context.Background()
	bucketName := config.Conf.MinIO.Bucket
	contentType := file.Header.Get("Content-Type")

	// 使用 PutObject 上传
	info, err := global.MinioClient.PutObject(ctx, bucketName, objectName, src, file.Size, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		response.Fail(c, "上传云存储失败: "+err.Error())
		return
	}

	// 5. 拼接访问 URL
	// 如果你的 Bucket 是 Public 的，URL 格式通常是: http://endpoint/bucket/objectName
	// 注意：如果 Endpoint 是 docker 内部 IP，这里返回给前端的 URL 需要改成宿主机 IP 或域名
	protocol := "http://"
	if config.Conf.MinIO.UseSSL {
		protocol = "https://"
	}

	// 这里直接拼接出可访问的 URL
	url := fmt.Sprintf("%s%s/%s/%s", protocol, config.Conf.MinIO.Endpoint, bucketName, info.Key)

	response.Success(c, gin.H{
		"url": url,
		"key": info.Key,
	})
}
