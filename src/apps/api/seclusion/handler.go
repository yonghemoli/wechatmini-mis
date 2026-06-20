package seclusionapi

import (
	"strconv"
	"yonghemolimis/src/apps/api/response"
	"yonghemolimis/src/usecase/seclusion"

	"github.com/gin-gonic/gin"
)

func GetSeclusion(c *gin.Context) {
	days, _ := strconv.Atoi(c.DefaultQuery("days", "30"))
	if days <= 0 {
		days = 30
	}
	data, err := seclusion.GetSeclusionAnalysis(days)
	if err != nil {
		response.Fail(c, "获取闭关分析失败: "+err.Error())
		return
	}
	response.OK(c, data)
}
