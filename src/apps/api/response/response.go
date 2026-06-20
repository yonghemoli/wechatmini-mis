package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// R 统一返回结构
type R struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func OK(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, R{Code: 0, Message: "success", Data: data})
}

func OKMsg(c *gin.Context, msg string) {
	c.JSON(http.StatusOK, R{Code: 0, Message: msg})
}

func Fail(c *gin.Context, msg string) {
	c.JSON(http.StatusOK, R{Code: -1, Message: msg})
}

func FailCode(c *gin.Context, code int, msg string) {
	c.JSON(http.StatusOK, R{Code: code, Message: msg})
}

func Error(c *gin.Context, httpCode int, msg string) {
	c.JSON(httpCode, R{Code: -1, Message: msg})
}

// PageData 分页响应
type PageData struct {
	List  interface{} `json:"list"`
	Total int64       `json:"total"`
	Page  int         `json:"page"`
	Size  int         `json:"size"`
}

func OKPage(c *gin.Context, list interface{}, total int64, page, size int) {
	c.JSON(http.StatusOK, R{
		Code:    0,
		Message: "success",
		Data: PageData{
			List:  list,
			Total: total,
			Page:  page,
			Size:  size,
		},
	})
}
