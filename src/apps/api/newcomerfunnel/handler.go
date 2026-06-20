package newcomerfunnelapi

import (
	"yonghemolimis/src/apps/api/response"
	"yonghemolimis/src/usecase/newcomerfunnel"

	"github.com/gin-gonic/gin"
)

func GetNewcomerFunnel(c *gin.Context) {
	data, err := newcomerfunnel.GetNewcomerFunnelAnalysis()
	if err != nil {
		response.Fail(c, "获取新人引导漏斗失败: "+err.Error())
		return
	}
	response.OK(c, data)
}
