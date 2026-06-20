package companionapi

import (
	"yonghemolimis/src/apps/api/response"
	"yonghemolimis/src/usecase/companion"

	"github.com/gin-gonic/gin"
)

func GetCompanionAnalysis(c *gin.Context) {
	data, err := companion.GetCompanionAnalysis()
	if err != nil {
		response.Fail(c, "\u83b7\u53d6\u9053\u4fa3\u5206\u6790\u5931\u8d25: "+err.Error())
		return
	}
	response.OK(c, data)
}
