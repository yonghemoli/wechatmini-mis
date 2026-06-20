package dungeonyonghemolimis

import (
	"strconv"

	"yonghemolimis/src/apps/api/response"
	uc "yonghemolimis/src/usecase/dungeonanalysis"

	"github.com/gin-gonic/gin"
)

func GetDungeonAnalysis(c *gin.Context) {
	days, _ := strconv.Atoi(c.DefaultQuery("days", "30"))
	data, err := uc.GetDungeonAnalysis(days)
	if err != nil {
		response.Fail(c, "获取副本分析失败: "+err.Error())
		return
	}
	response.OK(c, data)
}

func GetDungeonRealmDist(c *gin.Context) {
	days, _ := strconv.Atoi(c.DefaultQuery("days", "30"))
	data, err := uc.GetDungeonRealmDist(days)
	if err != nil {
		response.Fail(c, "获取副本境界分布失败: "+err.Error())
		return
	}
	response.OK(c, data)
}
