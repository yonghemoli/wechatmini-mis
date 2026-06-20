package refresh

import (
	"yonghemolimis/src/apps/api/response"
	"yonghemolimis/src/cron"

	"github.com/gin-gonic/gin"
)

// RefreshDashboard 手动刷新数据总览（快照+画像+分群）
// POST /api/v1/refresh/dashboard
func RefreshDashboard(c *gin.Context) {
	result, err := cron.RunDashboardRefresh()
	if err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.OK(c, result)
}

// RefreshFeatures 手动刷新功能分析（聚合+评分）
// POST /api/v1/refresh/features
func RefreshFeatures(c *gin.Context) {
	result, err := cron.RunFeatureRefresh()
	if err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.OK(c, result)
}

// RefreshAll 手动全量刷新
// POST /api/v1/refresh/all
func RefreshAll(c *gin.Context) {
	result, err := cron.RunAllTasks()
	if err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.OK(c, result)
}
