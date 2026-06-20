package db

import "github.com/xiuxianjs/xiuxian-game-db-sdk/analytics"

// ================== 分析库模型（从 SDK 统一导出） ==================

type (
	UserDO              = analytics.UserDO
	AnalyticsAdmin      = analytics.AnalyticsAdmin
	ActionLogDO         = analytics.ActionLogDO
	DailySnapshotDO     = analytics.DailySnapshotDO
	UserProfileDO       = analytics.UserProfileDO
	UserProfileTagDO    = analytics.UserProfileTagDO
	SegmentDO           = analytics.SegmentDO
	FeatureEventDO      = analytics.FeatureEventDO
	FeatureDailyStatsDO = analytics.FeatureDailyStatsDO
	FeatureScoreDO      = analytics.FeatureScoreDO
	GameActionEventDO   = analytics.GameActionEventDO
	UserHourlyStatsDO   = analytics.UserHourlyStatsDO
	GlobalHourlyStatsDO = analytics.GlobalHourlyStatsDO
)
