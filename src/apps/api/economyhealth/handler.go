package economyhealth

import (
	"strconv"

	"yonghemolimis/src/apps/api/response"
	uc "yonghemolimis/src/usecase/economyhealth"

	"github.com/gin-gonic/gin"
)

func GetEconomyHealth(c *gin.Context) {
	days, _ := strconv.Atoi(c.DefaultQuery("days", "30"))
	data, err := uc.GetEconomyHealth(days)
	if err != nil {
		response.Fail(c, "获取经济健康数据失败: "+err.Error())
		return
	}
	response.OK(c, data)
}

func GetCurrencyTrend(c *gin.Context) {
	currency := c.DefaultQuery("currency", "SPIRIT_STONE")
	days, _ := strconv.Atoi(c.DefaultQuery("days", "30"))
	data, err := uc.GetCurrencyTrend(currency, days)
	if err != nil {
		response.Fail(c, "获取币种趋势失败: "+err.Error())
		return
	}
	response.OK(c, data)
}
