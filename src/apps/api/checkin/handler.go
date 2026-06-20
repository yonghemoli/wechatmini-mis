package checkinapi

import (
	"strconv"
	"yonghemolimis/src/apps/api/response"
	"yonghemolimis/src/usecase/checkin"

	"github.com/gin-gonic/gin"
)

func GetCheckIn(c *gin.Context) {
	days, _ := strconv.Atoi(c.DefaultQuery("days", "30"))
	if days <= 0 {
		days = 30
	}
	data, err := checkin.GetCheckInAnalysis(days)
	if err != nil {
		response.Fail(c, "获取签到分析失败: "+err.Error())
		return
	}
	response.OK(c, data)
}
