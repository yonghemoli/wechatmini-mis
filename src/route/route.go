package route

import (
	activityapi "yonghemolimis/src/apps/api/activity"
	alertapi "yonghemolimis/src/apps/api/alert"
	checkinapi "yonghemolimis/src/apps/api/checkin"
	companionapi "yonghemolimis/src/apps/api/companion"
	craftanalysisapi "yonghemolimis/src/apps/api/craftanalysis"
	dashboardapi "yonghemolimis/src/apps/api/dashboard"
	dungeonapi "yonghemolimis/src/apps/api/dungeonanalysis"
	economyhealthapi "yonghemolimis/src/apps/api/economyhealth"
	equipanalysisapi "yonghemolimis/src/apps/api/equipanalysis"
	featureapi "yonghemolimis/src/apps/api/feature"
	funnelapi "yonghemolimis/src/apps/api/funnel"
	gamestatapi "yonghemolimis/src/apps/api/gamestat"
	guildhealthapi "yonghemolimis/src/apps/api/guildhealth"
	markethealthapi "yonghemolimis/src/apps/api/markethealth"
	miniapi "yonghemolimis/src/apps/api/mini"
	newcomerfunnelapi "yonghemolimis/src/apps/api/newcomerfunnel"
	payanalysisapi "yonghemolimis/src/apps/api/payanalysis"
	profileapi "yonghemolimis/src/apps/api/profile"
	realmbottleneckapi "yonghemolimis/src/apps/api/realmbottleneck"
	realmprogressapi "yonghemolimis/src/apps/api/realmprogress"
	refreshapi "yonghemolimis/src/apps/api/refresh"
	retentionapi "yonghemolimis/src/apps/api/retention"
	seclusionapi "yonghemolimis/src/apps/api/seclusion"
	segmentapi "yonghemolimis/src/apps/api/segment"
	shopanalysisapi "yonghemolimis/src/apps/api/shopanalysis"
	skillanalysisapi "yonghemolimis/src/apps/api/skillanalysis"
	socialapi "yonghemolimis/src/apps/api/socialnetwork"
	soulartifactapi "yonghemolimis/src/apps/api/soulartifact"
	towerstoryapi "yonghemolimis/src/apps/api/towerstory"
	trackapi "yonghemolimis/src/apps/api/track"
	"yonghemolimis/src/apps/api/user"
	"yonghemolimis/src/middlewares"

	"github.com/gin-gonic/gin"
)

// SetupRoutes 注册所有 API 路由
func SetupRoutes(r *gin.Engine) {
	mini := r.Group("/api/mini")
	miniapi.RegisterRoutes(mini)

	api := r.Group("/api/v1")

	// 公开接口
	api.POST("/login", user.Login)
	api.GET("/session", user.Session)
	api.POST("/logout", user.Logout)

	// 数据上报接口（游戏 SDK 调用，无需登录认证，但需要 API Key）
	api.POST("/track/features", middlewares.TrackAPIKeyRequired(), trackapi.ReceiveFeatureEvents)
	api.POST("/track/actions", middlewares.TrackAPIKeyRequired(), trackapi.ReceiveActionEvents)

	// 需要认证的接口
	auth := api.Group("", middlewares.AuthRequired())
	{
		// 用户信息
		auth.GET("/me", user.Me)

		// 仪表盘
		auth.GET("/dashboard/overview", dashboardapi.GetOverview)
		auth.GET("/dashboard/distribution", dashboardapi.GetDistribution)
		auth.GET("/dashboard/dau-trend", dashboardapi.GetDAUTrend)
		auth.GET("/dashboard/realm-distribution", dashboardapi.GetRealmDistribution)

		// 用户画像
		auth.GET("/profiles", profileapi.ListProfiles)
		auth.GET("/profiles/:uid", profileapi.GetProfile)
		auth.POST("/profiles/:uid/refresh", profileapi.RefreshProfile)
		auth.POST("/profiles/refresh-all", profileapi.RefreshAllProfiles)
		auth.GET("/profiles/:uid/snapshots", profileapi.GetSnapshots)

		// 用户分群
		auth.GET("/segments", segmentapi.List)
		auth.POST("/segments", segmentapi.Create)
		auth.DELETE("/segments/:id", segmentapi.Delete)
		auth.GET("/segments/:id/execute", segmentapi.Execute)
		auth.POST("/segments/preview", segmentapi.Preview)

		// 预警
		auth.GET("/alerts/churn", alertapi.GetHighChurnUsers)
		auth.GET("/alerts/stuck", alertapi.GetStuckUsers)
		auth.GET("/alerts/stats", alertapi.GetProfileStats)
		auth.GET("/alerts/resource", alertapi.GetResourceAlertUsers)
		auth.GET("/alerts/newbie-risk", alertapi.GetNewbieAtRisk)
		auth.GET("/alerts/whale-churn", alertapi.GetWhaleChurn)
		auth.GET("/alerts/bot-suspects", alertapi.GetBotSuspects)

		// 功能分析
		auth.GET("/features/scores", featureapi.GetScores)
		auth.GET("/features/trend", featureapi.GetTrend)
		auth.GET("/features/categories", featureapi.GetCategoryOverview)
		auth.GET("/features/top", featureapi.GetTopFeatures)
		auth.GET("/features/scene-dist", featureapi.GetSceneDistribution)
		auth.GET("/features/channel-top", featureapi.GetChannelTop)
		auth.GET("/features/platform-dist", featureapi.GetPlatformDistribution)
		auth.GET("/features/scene-feature-top", featureapi.GetSceneFeatureTop)
		auth.GET("/features/scene-trend", featureapi.GetSceneTrend)

		// 游戏数据统计
		auth.GET("/gamestat/overview", gamestatapi.GetEnhancedOverview)
		auth.GET("/gamestat/player-ranking", gamestatapi.GetPlayerRanking)
		auth.GET("/gamestat/guild-ranking", gamestatapi.GetGuildRanking)
		auth.GET("/gamestat/new-users", gamestatapi.GetNewUsersTrend)
		auth.GET("/gamestat/realm-stages", gamestatapi.GetRealmStageDistribution)
		auth.GET("/gamestat/realm-churn", gamestatapi.GetRealmChurn)
		auth.GET("/gamestat/payment", gamestatapi.GetPaymentOverview)
		auth.GET("/gamestat/revenue-trend", gamestatapi.GetRevenueTrend)
		auth.GET("/gamestat/packages", gamestatapi.GetPackageStats)

		// 留存分析
		auth.GET("/retention", retentionapi.GetRetention)

		// 转化漏斗
		auth.GET("/funnel", funnelapi.GetFunnel)

		// 境界进阶分析
		auth.GET("/realm-progress", realmprogressapi.GetRealmProgress)
		auth.GET("/realm-progress/pay-correlation", realmprogressapi.GetRealmPayCorrelation)

		// 经济健康度分析
		auth.GET("/economy-health", economyhealthapi.GetEconomyHealth)
		auth.GET("/economy-health/currency-trend", economyhealthapi.GetCurrencyTrend)

		// 副本分析
		auth.GET("/dungeon-analysis", dungeonapi.GetDungeonAnalysis)
		auth.GET("/dungeon-analysis/realm-dist", dungeonapi.GetDungeonRealmDist)

		// 宗门健康度
		auth.GET("/guild-health", guildhealthapi.GetGuildHealth)

		// 社交网络
		auth.GET("/social-network", socialapi.GetSocialNetwork)

		// 境界卡点分析
		auth.GET("/realm-bottleneck", realmbottleneckapi.GetRealmBottleneck)

		// 通天塔与主线进度
		auth.GET("/tower-story", towerstoryapi.GetTowerStory)
		auth.GET("/tower-story/stage-difficulty", towerstoryapi.GetStageDifficulty)

		// 付费行为深度分析
		auth.GET("/pay-analysis", payanalysisapi.GetPayAnalysis)

		// 装备与战力分析
		auth.GET("/equip-analysis", equipanalysisapi.GetEquipAnalysis)

		// 交易市场健康度
		auth.GET("/market-health", markethealthapi.GetMarketHealth)

		// 签到与活跃任务
		auth.GET("/checkin-yonghemolimis", checkinapi.GetCheckIn)

		// 闭关修炼分析
		auth.GET("/seclusion-yonghemolimis", seclusionapi.GetSeclusion)

		// 功法/技能分析
		auth.GET("/skill-analysis", skillanalysisapi.GetSkillAnalysis)

		// 炼丹/炼器分析
		auth.GET("/craft-analysis", craftanalysisapi.GetCraftAnalysis)

		// 新人引导漏斗
		auth.GET("/newcomer-funnel", newcomerfunnelapi.GetNewcomerFunnel)

		// 道侣系统分析
		auth.GET("/companion-yonghemolimis", companionapi.GetCompanionAnalysis)

		// 本命物系统分析
		auth.GET("/soul-artifact-yonghemolimis", soulartifactapi.GetSoulArtifactAnalysis)

		// 商店消费分析
		auth.GET("/shop-analysis", shopanalysisapi.GetShopAnalysis)

		// 活跃度分析
		auth.GET("/activity/player/:uid/hourly", activityapi.GetPlayerHourly)
		auth.GET("/activity/player/:uid/daily", activityapi.GetPlayerDaily)
		auth.GET("/activity/player/:uid/peak", activityapi.GetPlayerPeak)
		auth.GET("/activity/player/:uid/bot-check", activityapi.CheckPlayerBot)
		auth.GET("/activity/global/hourly", activityapi.GetGlobalHourly)
		auth.GET("/activity/global/peak", activityapi.GetGlobalPeak)
		auth.POST("/activity/refresh", activityapi.RefreshActivity)

		// 手动刷新
		auth.POST("/refresh/dashboard", refreshapi.RefreshDashboard)
		auth.POST("/refresh/features", refreshapi.RefreshFeatures)
		auth.POST("/refresh/all", refreshapi.RefreshAll)
	}
}
