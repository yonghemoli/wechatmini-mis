package towerstoryapi

import (
	"strconv"
	"yonghemolimis/src/apps/api/response"
	"yonghemolimis/src/usecase/towerstory"

	"github.com/gin-gonic/gin"
)

func GetTowerStory(c *gin.Context) {
	days, _ := strconv.Atoi(c.DefaultQuery("days", "30"))
	if days <= 0 {
		days = 30
	}
	data, err := towerstory.GetTowerStoryAnalysis(days)
	if err != nil {
		response.Fail(c, "获取通天塔主线分析失败: "+err.Error())
		return
	}
	response.OK(c, data)
}

func GetStageDifficulty(c *gin.Context) {
	chapter, _ := strconv.Atoi(c.DefaultQuery("chapter", "0"))
	data, err := towerstory.GetStageDifficulty(chapter)
	if err != nil {
		response.Fail(c, "获取关卡难度分析失败: "+err.Error())
		return
	}
	response.OK(c, data)
}
