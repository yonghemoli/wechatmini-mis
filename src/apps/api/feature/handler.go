package featureapi

import (
	"strconv"
	"yonghemolimis/src/apps/api/response"
	"yonghemolimis/src/usecase/feature"

	"github.com/gin-gonic/gin"
)

// GetScores 获取功能质量评分排行
// GET /api/v1/features/scores?category=
func GetScores(c *gin.Context) {
	category := c.Query("category")
	rows, err := feature.GetFeatureScores(category)
	if err != nil {
		response.Fail(c, "查询失败: "+err.Error())
		return
	}
	response.OK(c, rows)
}

// GetTrend 获取某功能的每日趋势
// GET /api/v1/features/trend?name=签到&days=30
func GetTrend(c *gin.Context) {
	name := c.Query("name")
	if name == "" {
		response.Fail(c, "缺少 name 参数")
		return
	}
	days, _ := strconv.Atoi(c.DefaultQuery("days", "30"))
	if days <= 0 {
		days = 30
	}
	rows, err := feature.GetFeatureDailyTrend(name, days)
	if err != nil {
		response.Fail(c, "查询失败: "+err.Error())
		return
	}
	response.OK(c, rows)
}

// GetCategoryOverview 按分类聚合概览
// GET /api/v1/features/categories?days=30
func GetCategoryOverview(c *gin.Context) {
	days, _ := strconv.Atoi(c.DefaultQuery("days", "30"))
	if days <= 0 {
		days = 30
	}
	rows, err := feature.GetCategoryOverview(days)
	if err != nil {
		response.Fail(c, "查询失败: "+err.Error())
		return
	}
	response.OK(c, rows)
}

// GetTopFeatures 热门功能 Top N
// GET /api/v1/features/top?days=7&limit=20
func GetTopFeatures(c *gin.Context) {
	days, _ := strconv.Atoi(c.DefaultQuery("days", "7"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	if days <= 0 {
		days = 7
	}
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	rows, err := feature.GetTopFeatures(days, limit)
	if err != nil {
		response.Fail(c, "查询失败: "+err.Error())
		return
	}
	response.OK(c, rows)
}

// GetSceneDistribution 场景分布（私聊/群聊）
// GET /api/v1/features/scene-dist?days=30
func GetSceneDistribution(c *gin.Context) {
	days, _ := strconv.Atoi(c.DefaultQuery("days", "30"))
	if days <= 0 {
		days = 30
	}
	rows, err := feature.GetSceneDistribution(days)
	if err != nil {
		response.Fail(c, "查询失败: "+err.Error())
		return
	}
	response.OK(c, rows)
}

// GetChannelTop 频道/群活跃排行
// GET /api/v1/features/channel-top?days=7&limit=20
func GetChannelTop(c *gin.Context) {
	days, _ := strconv.Atoi(c.DefaultQuery("days", "7"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	if days <= 0 {
		days = 7
	}
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	rows, err := feature.GetChannelTop(days, limit)
	if err != nil {
		response.Fail(c, "查询失败: "+err.Error())
		return
	}
	response.OK(c, rows)
}

// GetPlatformDistribution 平台分布
// GET /api/v1/features/platform-dist?days=30
func GetPlatformDistribution(c *gin.Context) {
	days, _ := strconv.Atoi(c.DefaultQuery("days", "30"))
	if days <= 0 {
		days = 30
	}
	rows, err := feature.GetPlatformDistribution(days)
	if err != nil {
		response.Fail(c, "查询失败: "+err.Error())
		return
	}
	response.OK(c, rows)
}

// GetSceneFeatureTop 各场景下的功能使用排名
// GET /api/v1/features/scene-feature-top?scene=private&days=7&limit=15
func GetSceneFeatureTop(c *gin.Context) {
	scene := c.Query("scene")
	days, _ := strconv.Atoi(c.DefaultQuery("days", "7"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "15"))
	if days <= 0 {
		days = 7
	}
	rows, err := feature.GetSceneFeatureTop(scene, days, limit)
	if err != nil {
		response.Fail(c, "查询失败: "+err.Error())
		return
	}
	response.OK(c, rows)
}

// GetSceneTrend 场景每日趋势
// GET /api/v1/features/scene-trend?days=30
func GetSceneTrend(c *gin.Context) {
	days, _ := strconv.Atoi(c.DefaultQuery("days", "30"))
	if days <= 0 {
		days = 30
	}
	rows, err := feature.GetSceneTrend(days)
	if err != nil {
		response.Fail(c, "查询失败: "+err.Error())
		return
	}
	response.OK(c, rows)
}
