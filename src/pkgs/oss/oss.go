// Package oss 提供 MIS 受保护的 MinIO/S3 上传能力，配置兼容 oss-admin。
package oss

import (
	"context"
	"fmt"
	"io"
	"path/filepath"
	"strings"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type Config struct {
	Endpoint, AccessKey, SecretKey, Bucket, PublicBaseURL string
	UseSSL                                                bool
}

type Service struct {
	client *minio.Client
	bucket string
	base   string
	ssl    bool
	server string
}

var service *Service

func Init(config *Config) error {
	if config == nil || strings.TrimSpace(config.Endpoint) == "" {
		return nil
	}
	client, err := minio.New(config.Endpoint, &minio.Options{Creds: credentials.NewStaticV4(config.AccessKey, config.SecretKey, ""), Secure: config.UseSSL})
	if err != nil {
		return fmt.Errorf("创建 OSS 客户端: %w", err)
	}
	if config.Bucket == "" {
		return fmt.Errorf("OSS bucket 未配置")
	}
	exists, err := client.BucketExists(context.Background(), config.Bucket)
	if err != nil {
		return fmt.Errorf("检查 OSS bucket: %w", err)
	}
	if !exists {
		if err := client.MakeBucket(context.Background(), config.Bucket, minio.MakeBucketOptions{}); err != nil {
			return fmt.Errorf("创建 OSS bucket: %w", err)
		}
	}
	service = &Service{client: client, bucket: config.Bucket, base: strings.TrimRight(config.PublicBaseURL, "/"), ssl: config.UseSSL, server: config.Endpoint}
	return nil
}

func GetService() *Service { return service }

func (s *Service) UploadImage(ctx context.Context, file io.Reader, size int64, contentType string) (string, error) {
	extension := extensionFromContentType(contentType)
	return s.upload(ctx, file, size, contentType, fmt.Sprintf("images/%d%s", time.Now().UnixNano(), extension))
}

func (s *Service) UploadFile(ctx context.Context, file io.Reader, size int64, contentType, originalName string) (string, error) {
	name := strings.TrimSpace(filepath.Base(originalName))
	if name == "" {
		name = "file"
	}
	return s.upload(ctx, file, size, contentType, fmt.Sprintf("files/%d-%s", time.Now().UnixNano(), name))
}

func (s *Service) upload(ctx context.Context, file io.Reader, size int64, contentType, objectName string) (string, error) {
	if _, err := s.client.PutObject(ctx, s.bucket, objectName, file, size, minio.PutObjectOptions{ContentType: contentType}); err != nil {
		return "", fmt.Errorf("上传 OSS 文件: %w", err)
	}
	if s.base != "" {
		return s.base + "/" + s.bucket + "/" + objectName, nil
	}
	scheme := "http"
	if s.ssl {
		scheme = "https"
	}
	return fmt.Sprintf("%s://%s/%s/%s", scheme, s.server, s.bucket, objectName), nil
}

func extensionFromContentType(contentType string) string {
	switch contentType {
	case "image/jpeg":
		return ".jpg"
	case "image/png":
		return ".png"
	case "image/gif":
		return ".gif"
	case "image/webp":
		return ".webp"
	default:
		return ""
	}
}
