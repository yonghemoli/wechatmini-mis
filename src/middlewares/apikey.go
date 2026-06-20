package middlewares

import (
	"net/http"
	"yonghemolimis/src/settings"

	"github.com/gin-gonic/gin"
)

// TrackAPIKeyRequired 校验数据上报接口的 X-API-Key 请求头
// 如果未配置 ANALYTICS_TRACK_API_KEY，则跳过校验（向后兼容）
func TrackAPIKeyRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		expected := settings.Conf.TrackAPIKey
		if expected == "" {
			// 未配置 key，放行
			c.Next()
			return
		}

		key := c.GetHeader("X-API-Key")
		if key == "" || key != expected {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code":    -1,
				"message": "无效的 API Key",
			})
			return
		}

		c.Next()
	}
}
