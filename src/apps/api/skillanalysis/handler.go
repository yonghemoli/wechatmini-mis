package skillyonghemolimisapi

import (
	"yonghemolimis/src/apps/api/response"
	skillyonghemolimis "yonghemolimis/src/usecase/skillanalysis"

	"github.com/gin-gonic/gin"
)

func GetSkillAnalysis(c *gin.Context) {
	data, err := skillyonghemolimis.GetSkillAnalysis()
	if err != nil {
		response.Fail(c, "获取功法分析失败: "+err.Error())
		return
	}
	response.OK(c, data)
}
