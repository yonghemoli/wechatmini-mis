package guildhealth

import (
	"strconv"

	"yonghemolimis/src/apps/api/response"
	uc "yonghemolimis/src/usecase/guildhealth"

	"github.com/gin-gonic/gin"
)

func GetGuildHealth(c *gin.Context) {
	inactiveDays, _ := strconv.Atoi(c.DefaultQuery("inactive_days", "7"))
	data, err := uc.GetGuildHealthData(inactiveDays)
	if err != nil {
		response.Fail(c, "获取宗门健康数据失败: "+err.Error())
		return
	}
	response.OK(c, data)
}
