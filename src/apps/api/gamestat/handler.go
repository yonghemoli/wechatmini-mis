package gamestatapi

import (
	"strconv"
	"yonghemolimis/src/apps/api/response"
	"yonghemolimis/src/usecase/gamestat"

	"github.com/gin-gonic/gin"
)

// GetEnhancedOverview 增强版概览
func GetEnhancedOverview(c *gin.Context) {
	data, err := gamestat.GetEnhancedOverview()
	if err != nil {
		response.Fail(c, "获取概览失败: "+err.Error())
		return
	}
	response.OK(c, data)
}

// GetPlayerRanking 玩家排行
func GetPlayerRanking(c *gin.Context) {
	sortBy := c.DefaultQuery("sort", "realm")
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	items, err := gamestat.GetPlayerRanking(sortBy, limit)
	if err != nil {
		response.Fail(c, "查询失败: "+err.Error())
		return
	}
	response.OK(c, items)
}

// GetGuildRanking 宗门排行
func GetGuildRanking(c *gin.Context) {
	sortBy := c.DefaultQuery("sort", "prestige")
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "30"))
	items, err := gamestat.GetGuildRanking(sortBy, limit)
	if err != nil {
		response.Fail(c, "查询失败: "+err.Error())
		return
	}
	response.OK(c, items)
}

// GetNewUsersTrend 新增玩家趋势
func GetNewUsersTrend(c *gin.Context) {
	days, _ := strconv.Atoi(c.DefaultQuery("days", "30"))
	if days <= 0 {
		days = 30
	}
	items, err := gamestat.GetNewUsersTrend(days)
	if err != nil {
		response.Fail(c, "查询失败: "+err.Error())
		return
	}
	response.OK(c, items)
}

// GetRealmStageDistribution 大境界分布
func GetRealmStageDistribution(c *gin.Context) {
	items, err := gamestat.GetRealmStageDistribution()
	if err != nil {
		response.Fail(c, "查询失败: "+err.Error())
		return
	}
	response.OK(c, items)
}

// GetRealmChurn 境界流失率
func GetRealmChurn(c *gin.Context) {
	days, _ := strconv.Atoi(c.DefaultQuery("days", "7"))
	if days <= 0 {
		days = 7
	}
	items, err := gamestat.GetRealmChurn(days)
	if err != nil {
		response.Fail(c, "查询失败: "+err.Error())
		return
	}
	response.OK(c, items)
}

// GetPaymentOverview 付费概览
func GetPaymentOverview(c *gin.Context) {
	source := c.DefaultQuery("source", "official")
	data, err := gamestat.GetPaymentOverview(source)
	if err != nil {
		response.Fail(c, "查询失败: "+err.Error())
		return
	}
	response.OK(c, data)
}

// GetRevenueTrend 营收趋势
func GetRevenueTrend(c *gin.Context) {
	days, _ := strconv.Atoi(c.DefaultQuery("days", "30"))
	if days <= 0 {
		days = 30
	}
	source := c.DefaultQuery("source", "official")
	items, err := gamestat.GetRevenueTrend(days, source)
	if err != nil {
		response.Fail(c, "查询失败: "+err.Error())
		return
	}
	response.OK(c, items)
}

// GetPackageStats 套餐销售统计
func GetPackageStats(c *gin.Context) {
	days, _ := strconv.Atoi(c.DefaultQuery("days", "0"))
	source := c.DefaultQuery("source", "official")
	items, err := gamestat.GetPackageStats(days, source)
	if err != nil {
		response.Fail(c, "查询失败: "+err.Error())
		return
	}
	response.OK(c, items)
}
