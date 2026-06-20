package realmbottleneckapi

import (
	"strconv"
	"yonghemolimis/src/apps/api/response"
	"yonghemolimis/src/usecase/realmbottleneck"

	"github.com/gin-gonic/gin"
)

func GetRealmBottleneck(c *gin.Context) {
	days, _ := strconv.Atoi(c.DefaultQuery("inactive_days", "7"))
	if days <= 0 {
		days = 7
	}
	data, err := realmbottleneck.GetRealmBottleneckAnalysis(days)
	if err != nil {
		response.Fail(c, "获取境界卡点分析失败: "+err.Error())
		return
	}
	response.OK(c, data)
}
