package activityapi

import (
	"strconv"
	"yonghemolimis/src/apps/api/response"
	"yonghemolimis/src/usecase/activity"

	"github.com/gin-gonic/gin"
)

// ==================== 玩家维度 ====================

// GetPlayerHourly 获取某玩家的小时活跃数据
// GET /api/v1/activity/player/:uid/hourly?days=7
func GetPlayerHourly(c *gin.Context) {
	uid, err := strconv.ParseUint(c.Param("uid"), 10, 64)
	if err != nil {
		response.Fail(c, "无效的 UID")
		return
	}
	days, _ := strconv.Atoi(c.DefaultQuery("days", "7"))
	if days <= 0 {
		days = 7
	}
	rows, err := activity.GetPlayerHourlyStats(uint(uid), days)
	if err != nil {
		response.Fail(c, "查询失败: "+err.Error())
		return
	}
	response.OK(c, rows)
}

// GetPlayerDaily 获取某玩家的每日汇总
// GET /api/v1/activity/player/:uid/daily?days=30
func GetPlayerDaily(c *gin.Context) {
	uid, err := strconv.ParseUint(c.Param("uid"), 10, 64)
	if err != nil {
		response.Fail(c, "无效的 UID")
		return
	}
	days, _ := strconv.Atoi(c.DefaultQuery("days", "30"))
	if days <= 0 {
		days = 30
	}
	rows, err := activity.GetPlayerDailyStats(uint(uid), days)
	if err != nil {
		response.Fail(c, "查询失败: "+err.Error())
		return
	}
	response.OK(c, rows)
}

// GetPlayerPeak 获取某玩家的高峰时段分析
// GET /api/v1/activity/player/:uid/peak?days=30
func GetPlayerPeak(c *gin.Context) {
	uid, err := strconv.ParseUint(c.Param("uid"), 10, 64)
	if err != nil {
		response.Fail(c, "无效的 UID")
		return
	}
	days, _ := strconv.Atoi(c.DefaultQuery("days", "30"))
	if days <= 0 {
		days = 30
	}
	result, err := activity.GetPlayerPeakAnalysis(uint(uid), days)
	if err != nil {
		response.Fail(c, "查询失败: "+err.Error())
		return
	}
	response.OK(c, result)
}

// CheckPlayerBot 检测某玩家是否为机器人
// GET /api/v1/activity/player/:uid/bot-check?days=14
func CheckPlayerBot(c *gin.Context) {
	uid, err := strconv.ParseUint(c.Param("uid"), 10, 64)
	if err != nil {
		response.Fail(c, "无效的 UID")
		return
	}
	days, _ := strconv.Atoi(c.DefaultQuery("days", "14"))
	if days <= 0 {
		days = 14
	}
	result, err := activity.CheckBot(uint(uid), days)
	if err != nil {
		response.Fail(c, "检测失败: "+err.Error())
		return
	}
	response.OK(c, result)
}

// ==================== 全局维度 ====================

// GetGlobalHourly 获取全局小时活跃数据
// GET /api/v1/activity/global/hourly?days=7
func GetGlobalHourly(c *gin.Context) {
	days, _ := strconv.Atoi(c.DefaultQuery("days", "7"))
	if days <= 0 {
		days = 7
	}
	rows, err := activity.GetGlobalHourlyStats(days)
	if err != nil {
		response.Fail(c, "查询失败: "+err.Error())
		return
	}
	response.OK(c, rows)
}

// GetGlobalPeak 获取全局高峰时段分析
// GET /api/v1/activity/global/peak?days=30
func GetGlobalPeak(c *gin.Context) {
	days, _ := strconv.Atoi(c.DefaultQuery("days", "30"))
	if days <= 0 {
		days = 30
	}
	result, err := activity.GetPeakAnalysis(days)
	if err != nil {
		response.Fail(c, "查询失败: "+err.Error())
		return
	}
	response.OK(c, result)
}

// RefreshActivity 手动刷新活跃度聚合数据
// POST /api/v1/activity/refresh?days=3
func RefreshActivity(c *gin.Context) {
	days, _ := strconv.Atoi(c.DefaultQuery("days", "3"))
	if days <= 0 || days > 90 {
		days = 3
	}
	go func() {
		_ = activity.AggregateMultipleDays(days)
	}()
	response.OKMsg(c, "活跃度聚合任务已启动")
}
