package mis

import (
	"net/http"
	"path/filepath"
	"strings"

	"yonghemolimis/src/apps/api/response"
	"yonghemolimis/src/pkgs/oss"

	"github.com/gin-gonic/gin"
)

const (
	maxImageSize = 5 << 20
	maxFileSize  = 20 << 20
)

var allowedImageTypes = map[string]bool{"image/jpeg": true, "image/png": true, "image/gif": true, "image/webp": true}
var allowedFileExtensions = map[string]bool{".jpg": true, ".jpeg": true, ".png": true, ".gif": true, ".webp": true, ".pdf": true}

// UploadImage 使用 oss-admin 同一套 MinIO 服务上传图片，只接受已登录的 MIS 管理员请求。
func UploadImage(c *gin.Context) { uploadAsset(c, true) }

// UploadFile 上传证书/体检报告等图片或 PDF 文件。
func UploadFile(c *gin.Context) { uploadAsset(c, false) }

func uploadAsset(c *gin.Context, imageOnly bool) {
	svc := oss.GetService()
	if svc == nil {
		response.Error(c, http.StatusServiceUnavailable, "文件存储服务暂不可用")
		return
	}
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		response.Error(c, http.StatusBadRequest, "请选择文件")
		return
	}
	defer file.Close()
	if header.Size <= 0 || header.Size > func() int64 {
		if imageOnly {
			return maxImageSize
		}
		return maxFileSize
	}() {
		response.Error(c, http.StatusBadRequest, "文件大小不符合要求")
		return
	}
	contentType := strings.ToLower(strings.TrimSpace(header.Header.Get("Content-Type")))
	extension := strings.ToLower(filepath.Ext(header.Filename))
	if imageOnly && !allowedImageTypes[contentType] {
		response.Error(c, http.StatusBadRequest, "仅支持 JPG、PNG、GIF、WEBP 图片")
		return
	}
	if !imageOnly && !allowedFileExtensions[extension] {
		response.Error(c, http.StatusBadRequest, "仅支持图片或 PDF 文件")
		return
	}
	var url string
	if imageOnly {
		url, err = svc.UploadImage(c.Request.Context(), file, header.Size, contentType)
	} else {
		url, err = svc.UploadFile(c.Request.Context(), file, header.Size, contentType, filepath.Base(header.Filename))
	}
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "文件上传失败")
		return
	}
	response.OK(c, gin.H{"url": url})
}
