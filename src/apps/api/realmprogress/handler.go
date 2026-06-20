package realmprogressapi

import (
	"strconv"
	"yonghemolimis/src/apps/api/response"
	"yonghemolimis/src/usecase/realmprogress"

	"github.com/gin-gonic/gin"
)

// GetRealmProgress 境界进阶分析
func GetRealmProgress(c *gin.Context) {
	days, _ := strconv.Atoi(c.DefaultQuery("inactive_days", "7"))
	if days <= 0 {
		days = 7
	}
	data, err := realmprogress.GetRealmProgressAnalysis(days)
	if err != nil {
		response.Fail(c, "获取境界进阶分析失败: "+err.Error())
		return
	}
	response.OK(c, data)
}

// GetRealmPayCorrelation 境界-付费关联
func GetRealmPayCorrelation(c *gin.Context) {
	data, err := realmprogress.GetRealmPayCorrelation()
	if err != nil {
		response.Fail(c, "获取境界付费关联失败: "+err.Error())
		return
	}
	response.OK(c, data)
}
