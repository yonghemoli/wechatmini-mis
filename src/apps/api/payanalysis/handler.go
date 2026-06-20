package payyonghemolimisapi

import (
	"strconv"
	"yonghemolimis/src/apps/api/response"
	payyonghemolimis "yonghemolimis/src/usecase/payanalysis"

	"github.com/gin-gonic/gin"
)

func GetPayAnalysis(c *gin.Context) {
	days, _ := strconv.Atoi(c.DefaultQuery("days", "30"))
	if days <= 0 {
		days = 30
	}
	data, err := payyonghemolimis.GetPayAnalysis(days)
	if err != nil {
		response.Fail(c, "获取付费分析失败: "+err.Error())
		return
	}
	response.OK(c, data)
}
