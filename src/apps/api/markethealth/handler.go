package markethealthapi

import (
	"strconv"
	"yonghemolimis/src/apps/api/response"
	"yonghemolimis/src/usecase/markethealth"

	"github.com/gin-gonic/gin"
)

func GetMarketHealth(c *gin.Context) {
	days, _ := strconv.Atoi(c.DefaultQuery("days", "30"))
	if days <= 0 {
		days = 30
	}
	data, err := markethealth.GetMarketHealth(days)
	if err != nil {
		response.Fail(c, "获取交易市场分析失败: "+err.Error())
		return
	}
	response.OK(c, data)
}
