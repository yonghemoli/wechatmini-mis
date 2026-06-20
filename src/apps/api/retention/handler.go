package retentionapi

import (
	"strconv"
	"yonghemolimis/src/apps/api/response"
	"yonghemolimis/src/usecase/retention"

	"github.com/gin-gonic/gin"
)

// GetRetention 留存分析
func GetRetention(c *gin.Context) {
	days, _ := strconv.Atoi(c.DefaultQuery("days", "30"))
	if days <= 0 {
		days = 30
	}
	data, err := retention.GetRetentionData(days)
	if err != nil {
		response.Fail(c, "获取留存数据失败: "+err.Error())
		return
	}
	response.OK(c, data)
}
