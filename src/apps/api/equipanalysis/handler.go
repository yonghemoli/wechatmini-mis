package equipyonghemolimisapi

import (
	"yonghemolimis/src/apps/api/response"
	equipyonghemolimis "yonghemolimis/src/usecase/equipanalysis"

	"github.com/gin-gonic/gin"
)

func GetEquipAnalysis(c *gin.Context) {
	data, err := equipyonghemolimis.GetEquipAnalysis()
	if err != nil {
		response.Fail(c, "获取装备分析失败: "+err.Error())
		return
	}
	response.OK(c, data)
}
