package shopyonghemolimisapi

import (
	"yonghemolimis/src/apps/api/response"
	shopyonghemolimis "yonghemolimis/src/usecase/shopanalysis"

	"github.com/gin-gonic/gin"
)

func GetShopAnalysis(c *gin.Context) {
	data, err := shopyonghemolimis.GetShopAnalysis()
	if err != nil {
		response.Fail(c, "\u83b7\u53d6\u5546\u5e97\u6d88\u8d39\u5206\u6790\u5931\u8d25: "+err.Error())
		return
	}
	response.OK(c, data)
}
