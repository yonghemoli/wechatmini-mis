package alertapi

import (
	"strconv"
	"time"
	"yonghemolimis/src/apps/api/response"
	"yonghemolimis/src/dao/analytics"
	"yonghemolimis/src/dao/db"

	"github.com/gin-gonic/gin"
)

// GetHighChurnUsers 高流失风险用户
func GetHighChurnUsers(c *gin.Context) {
	threshold, _ := strconv.Atoi(c.DefaultQuery("threshold", "60"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))

	list, err := analytics.GetHighChurnProfiles(threshold, limit)
	if err != nil {
		response.Fail(c, "查询失败: "+err.Error())
		return
	}
	response.OK(c, list)
}

// GetStuckUsers 卡关用户
func GetStuckUsers(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))

	list, err := analytics.GetStuckProfiles(limit)
	if err != nil {
		response.Fail(c, "查询失败: "+err.Error())
		return
	}
	response.OK(c, list)
}

// GetProfileStats 画像统计摘要
func GetProfileStats(c *gin.Context) {
	stats, err := analytics.GetProfileStats()
	if err != nil {
		response.Fail(c, "查询统计失败: "+err.Error())
		return
	}
	response.OK(c, stats)
}

// GetResourceAlertUsers 资源告急用户
func GetResourceAlertUsers(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	list, err := analytics.GetResourceAlertProfiles(limit)
	if err != nil {
		response.Fail(c, "查询失败: "+err.Error())
		return
	}
	response.OK(c, list)
}

// GetNewbieAtRisk 新手流失预警
func GetNewbieAtRisk(c *gin.Context) {
	threshold, _ := strconv.Atoi(c.DefaultQuery("threshold", "40"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	list, err := analytics.GetNewbieAtRiskProfiles(threshold, limit)
	if err != nil {
		response.Fail(c, "查询失败: "+err.Error())
		return
	}
	response.OK(c, list)
}

// GetWhaleChurn 大氪沉默预警
func GetWhaleChurn(c *gin.Context) {
	threshold, _ := strconv.Atoi(c.DefaultQuery("threshold", "50"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	list, err := analytics.GetWhaleChurnProfiles(threshold, limit)
	if err != nil {
		response.Fail(c, "查询失败: "+err.Error())
		return
	}
	response.OK(c, list)
}

// GetBotSuspects 长期在线异常用户（日均活跃18h+）
func GetBotSuspects(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	minHours, _ := strconv.Atoi(c.DefaultQuery("min_hours", "18"))

	var results []map[string]any
	err := db.Get().Model(&db.UserHourlyStatsDO{}).
		Select("uid, COUNT(DISTINCT date) as active_days, COUNT(*) as total_records, ROUND(COUNT(*)::numeric / GREATEST(COUNT(DISTINCT date), 1), 1) as avg_hours").
		Where("date >= ?", time.Now().AddDate(0, 0, -14).Format("2006-01-02")).
		Group("uid").
		Having("COUNT(*)::numeric / GREATEST(COUNT(DISTINCT date), 1) >= ?", minHours).
		Order("avg_hours DESC").
		Limit(limit).
		Find(&results).Error
	if err != nil {
		response.Fail(c, "查询失败: "+err.Error())
		return
	}
	response.OK(c, results)
}
