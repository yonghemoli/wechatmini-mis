package funnelapi

import (
	"strconv"
	"yonghemolimis/src/apps/api/response"
	"yonghemolimis/src/usecase/funnel"

	"github.com/gin-gonic/gin"
)

// GetFunnel 转化漏斗
func GetFunnel(c *gin.Context) {
	days, _ := strconv.Atoi(c.DefaultQuery("days", "0"))
	data, err := funnel.GetFunnelData(days)
	if err != nil {
		response.Fail(c, "获取漏斗数据失败: "+err.Error())
		return
	}
	response.OK(c, data)
}
