package craftanalysis

import (
	"yonghemolimis/src/apps/api/response"
	craftyonghemolimis "yonghemolimis/src/usecase/craftanalysis"

	"github.com/gin-gonic/gin"
)

func GetCraftAnalysis(c *gin.Context) {
	data, err := craftyonghemolimis.GetCraftAnalysis()
	if err != nil {
		response.Fail(c, "获取炼丹炼器分析失败: "+err.Error())
		return
	}
	response.OK(c, data)
}
