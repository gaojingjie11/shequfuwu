package service

import (
	"context"
	"fmt"
	"mime/multipart"
	"path/filepath"
	"smartcommunity/internal/config"
	"smartcommunity/internal/global"
	"time"

	"github.com/minio/minio-go/v7"
)

type StorageService struct{}

func (s *StorageService) UploadMultipartFile(fileHeader *multipart.FileHeader, dir string) (string, string, error) {
	src, err := fileHeader.Open()
	if err != nil {
		return "", "", err
	}
	defer src.Close()

	ext := filepath.Ext(fileHeader.Filename)
	objectName := fmt.Sprintf("%s/%d%s", dir, time.Now().UnixNano(), ext)

	ctx := context.Background()
	bucketName := config.Conf.MinIO.Bucket
	contentType := fileHeader.Header.Get("Content-Type")

	info, err := global.MinioClient.PutObject(ctx, bucketName, objectName, src, fileHeader.Size, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return "", "", err
	}

	protocol := "http://"
	if config.Conf.MinIO.UseSSL {
		protocol = "https://"
	}
	url := fmt.Sprintf("%s%s/%s/%s", protocol, config.Conf.MinIO.Endpoint, bucketName, info.Key)

	return url, info.Key, nil
}
