package dashboardapi

import (
	"strconv"
	"time"
	"yonghemolimis/src/apps/api/response"
	"yonghemolimis/src/usecase/dashboard"

	"github.com/gin-gonic/gin"
)

// GetOverview 仪表盘概览
func GetOverview(c *gin.Context) {
	source := c.DefaultQuery("source", "official")
	data, err := dashboard.GetOverview(source)
	if err != nil {
		response.Fail(c, "获取概览失败: "+err.Error())
		return
	}
	response.OK(c, data)
}

// GetDistribution 画像分布
func GetDistribution(c *gin.Context) {
	field := c.DefaultQuery("field", "lifecycle_stage")
	items, err := dashboard.GetDistribution(field)
	if err != nil {
		response.Fail(c, "获取分布失败: "+err.Error())
		return
	}
	response.OK(c, items)
}

// GetDAUTrend DAU 趋势
func GetDAUTrend(c *gin.Context) {
	start := c.Query("start")
	end := c.Query("end")
	// 支持 days 参数快捷计算
	if start == "" || end == "" {
		daysStr := c.DefaultQuery("days", "30")
		days, _ := strconv.Atoi(daysStr)
		if days <= 0 {
			days = 30
		}
		now := time.Now()
		end = now.Format("2006-01-02")
		start = now.AddDate(0, 0, -days).Format("2006-01-02")
	}
	items, err := dashboard.GetDAUTrend(start, end)
	if err != nil {
		response.Fail(c, "获取 DAU 趋势失败: "+err.Error())
		return
	}
	response.OK(c, items)
}

// GetRealmDistribution 境界分布
func GetRealmDistribution(c *gin.Context) {
	items, err := dashboard.GetRealmDistribution()
	if err != nil {
		response.Fail(c, "获取境界分布失败: "+err.Error())
		return
	}
	response.OK(c, items)
}
