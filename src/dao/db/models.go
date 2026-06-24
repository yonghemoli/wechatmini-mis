package db

import "time"

// UserDO 管理后台用户（历史兼容表）
type UserDO struct {
	ID        uint `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time

	UserName string `gorm:"size:128;uniqueIndex"`
	PassWord string `gorm:"size:256"`
}

func (UserDO) TableName() string { return "admin_users" }

// AnalyticsAdmin 管理员表（历史兼容表）
type AnalyticsAdmin struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	Username     string    `gorm:"size:50;not null" json:"username"`
	Name         string    `gorm:"size:50;not null" json:"name"`
	IsSuperAdmin bool      `gorm:"not null;default:false" json:"isSuperAdmin"`
	Status       string    `gorm:"size:20;not null;default:'active'" json:"status"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

func (AnalyticsAdmin) TableName() string { return "analytics_admins" }

// ActionLogDO 用户行为事件日志
type ActionLogDO struct {
	ID        uint `gorm:"primaryKey"`
	CreatedAt time.Time

	UID      uint   `gorm:"index:idx_uid_action;index:idx_uid_created;not null;column:uid"`
	Action   string `gorm:"size:64;index:idx_uid_action;not null"`
	Category string `gorm:"size:32;index:idx_category_created;not null"`
	Detail   string `gorm:"type:text"`
}

func (ActionLogDO) TableName() string { return "user_action_logs" }

// DailySnapshotDO 每日用户快照
type DailySnapshotDO struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"created_at"`

	UID          uint   `gorm:"uniqueIndex:uk_uid_date;not null;column:uid" json:"uid"`
	SnapshotDate string `gorm:"size:10;uniqueIndex:uk_uid_date;index:idx_date;not null;column:snapshot_date" json:"snapshot_date"`

	RealmID       int   `gorm:"column:realm_id" json:"realm_id"`
	VipLevel      int   `gorm:"column:vip_level" json:"vip_level"`
	SpiritStone   int64 `gorm:"column:spirit_stone" json:"spirit_stone"`
	SpiritCrystal int64 `gorm:"column:spirit_crystal" json:"spirit_crystal"`
	SpiritMilk    int64 `gorm:"column:spirit_milk" json:"spirit_milk"`

	LoginCount   int `gorm:"default:0;column:login_count" json:"login_count"`
	BattleCount  int `gorm:"default:0;column:battle_count" json:"battle_count"`
	CraftCount   int `gorm:"default:0;column:craft_count" json:"craft_count"`
	SocialCount  int `gorm:"default:0;column:social_count" json:"social_count"`
	EconomyCount int `gorm:"default:0;column:economy_count" json:"economy_count"`

	StoneIncome    int64 `gorm:"default:0;column:stone_income" json:"stone_income"`
	StoneExpense   int64 `gorm:"default:0;column:stone_expense" json:"stone_expense"`
	CrystalIncome  int64 `gorm:"default:0;column:crystal_income" json:"crystal_income"`
	CrystalExpense int64 `gorm:"default:0;column:crystal_expense" json:"crystal_expense"`
}

func (DailySnapshotDO) TableName() string { return "user_daily_snapshots" }

// UserProfileDO 用户画像主表
type UserProfileDO struct {
	UID       uint      `gorm:"primaryKey;column:uid" json:"uid"`
	UpdatedAt time.Time `json:"updated_at"`

	LifecycleStage   string     `gorm:"size:16;index:idx_lifecycle;column:lifecycle_stage" json:"lifecycle_stage"`
	PlayStyle        string     `gorm:"size:16;column:play_style" json:"play_style"`
	PayTier          string     `gorm:"size:16;index:idx_pay_tier;column:pay_tier" json:"pay_tier"`
	LtvPredict       float64    `gorm:"column:ltv_predict" json:"ltv_predict"`
	SocialType       string     `gorm:"size:20;column:social_type" json:"social_type"`
	ChurnRisk        int        `gorm:"default:0;index:idx_churn_risk;column:churn_risk" json:"churn_risk"`
	StuckFlag        bool       `gorm:"default:false;column:stuck_flag" json:"stuck_flag"`
	ResourceAlert    bool       `gorm:"default:false;column:resource_alert" json:"resource_alert"`
	LastCalculatedAt *time.Time `gorm:"column:last_calculated_at" json:"last_calculated_at"`
}

func (UserProfileDO) TableName() string { return "user_profiles" }

// UserProfileTagDO 用户标签明细
type UserProfileTagDO struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UpdatedAt time.Time `json:"updated_at"`

	UID      uint       `gorm:"uniqueIndex:uk_uid_group_key;not null;column:uid" json:"uid"`
	TagGroup string     `gorm:"size:32;uniqueIndex:uk_uid_group_key;not null;column:tag_group" json:"tag_group"`
	TagKey   string     `gorm:"size:64;uniqueIndex:uk_uid_group_key;index:idx_tag_key_value;not null;column:tag_key" json:"tag_key"`
	TagValue string     `gorm:"size:128;index:idx_tag_key_value;not null;column:tag_value" json:"tag_value"`
	Score    *float64   `gorm:"column:score" json:"score"`
	ExpireAt *time.Time `gorm:"column:expire_at" json:"expire_at"`
}

func (UserProfileTagDO) TableName() string { return "user_profile_tags" }

// SegmentDO 用户分群规则
type SegmentDO struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Name        string `gorm:"size:128;not null" json:"name"`
	Description string `gorm:"type:text" json:"description"`
	RulesJSON   string `gorm:"type:text;not null;column:rules_json" json:"rules_json"`
	UserCount   int    `gorm:"default:0;column:user_count" json:"user_count"`
	CreatedBy   string `gorm:"size:128;column:created_by" json:"created_by"`
}

func (SegmentDO) TableName() string { return "segments" }

// FeatureEventDO 原始功能埋点事件
type FeatureEventDO struct {
	ID        uint `gorm:"primaryKey"`
	CreatedAt time.Time

	UID             uint      `gorm:"index:idx_fe_uid;index:idx_fe_uid_created;not null;column:uid"`
	FeatureName     string    `gorm:"size:64;index:idx_fe_feature;not null;column:feature_name"`
	FeatureCategory string    `gorm:"size:32;index:idx_fe_category;not null;column:feature_category"`
	CommandText     string    `gorm:"size:255;column:command_text"`
	ResponseTimeMs  int       `gorm:"default:0;column:response_time_ms"`
	Success         bool      `gorm:"default:true;column:success"`
	Scene           string    `gorm:"size:16;index:idx_fe_scene;column:scene"`
	ChannelID       string    `gorm:"size:64;index:idx_fe_channel;column:channel_id"`
	Platform        string    `gorm:"size:32;index:idx_fe_platform;column:platform"`
	EventTime       time.Time `gorm:"index:idx_fe_uid_created;index:idx_fe_event_time;not null;column:event_time"`
}

func (FeatureEventDO) TableName() string { return "feature_events" }

// FeatureDailyStatsDO 功能每日聚合统计
type FeatureDailyStatsDO struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"created_at"`

	Date            string  `gorm:"size:10;uniqueIndex:uk_date_feature;index:idx_fds_date;not null;column:date" json:"date"`
	FeatureName     string  `gorm:"size:64;uniqueIndex:uk_date_feature;index:idx_fds_feature;not null;column:feature_name" json:"feature_name"`
	FeatureCategory string  `gorm:"size:32;index:idx_fds_category;not null;column:feature_category" json:"feature_category"`
	TotalUses       int     `gorm:"default:0;column:total_uses" json:"total_uses"`
	UniqueUsers     int     `gorm:"default:0;column:unique_users" json:"unique_users"`
	AvgResponseMs   float64 `gorm:"default:0;column:avg_response_ms" json:"avg_response_ms"`
	SuccessRate     float64 `gorm:"default:1;column:success_rate" json:"success_rate"`
}

func (FeatureDailyStatsDO) TableName() string { return "feature_daily_stats" }

// FeatureScoreDO 功能质量评分
type FeatureScoreDO struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UpdatedAt time.Time `json:"updated_at"`

	FeatureName     string  `gorm:"size:64;uniqueIndex;not null;column:feature_name" json:"feature_name"`
	FeatureCategory string  `gorm:"size:32;index:idx_fs_category;not null;column:feature_category" json:"feature_category"`
	TotalUses7d     int     `gorm:"default:0;column:total_uses_7d" json:"total_uses_7d"`
	TotalUses30d    int     `gorm:"default:0;column:total_uses_30d" json:"total_uses_30d"`
	UniqueUsers7d   int     `gorm:"default:0;column:unique_users_7d" json:"unique_users_7d"`
	UniqueUsers30d  int     `gorm:"default:0;column:unique_users_30d" json:"unique_users_30d"`
	AvgDailyUses    float64 `gorm:"default:0;column:avg_daily_uses" json:"avg_daily_uses"`
	AvgResponseMs   float64 `gorm:"default:0;column:avg_response_ms" json:"avg_response_ms"`
	SuccessRate     float64 `gorm:"default:1;column:success_rate" json:"success_rate"`
	UsageGrowth     float64 `gorm:"default:0;column:usage_growth" json:"usage_growth"`
	UserPenetration float64 `gorm:"default:0;column:user_penetration" json:"user_penetration"`
	QualityScore    float64 `gorm:"default:0;column:quality_score" json:"quality_score"`
}

func (FeatureScoreDO) TableName() string { return "feature_scores" }

// GameActionEventDO 游戏行为事件（历史兼容表）
type GameActionEventDO struct {
	ID        uint      `gorm:"primaryKey"`
	CreatedAt time.Time `gorm:"index:idx_gae_created"`

	UID       uint      `gorm:"index:idx_gae_uid;index:idx_gae_uid_action;not null;column:uid"`
	Action    string    `gorm:"size:64;index:idx_gae_uid_action;index:idx_gae_action;not null;column:action"`
	Category  string    `gorm:"size:32;index:idx_gae_category;not null;column:category"`
	Result    string    `gorm:"size:16;not null;column:result"`
	Value     int64     `gorm:"default:0;column:value"`
	Detail    string    `gorm:"type:text;column:detail"`
	Scene     string    `gorm:"size:16;index:idx_gae_scene;column:scene"`
	ChannelID string    `gorm:"size:64;index:idx_gae_channel;column:channel_id"`
	Platform  string    `gorm:"size:32;index:idx_gae_platform;column:platform"`
	EventTime time.Time `gorm:"index:idx_gae_event_time;not null;column:event_time"`
}

func (GameActionEventDO) TableName() string { return "game_action_events" }

// UserHourlyStatsDO 用户每小时活跃统计
type UserHourlyStatsDO struct {
	ID uint `gorm:"primaryKey" json:"id"`

	UID          uint   `gorm:"uniqueIndex:uk_uid_date_hour;index:idx_uhs_uid;not null;column:uid" json:"uid"`
	Date         string `gorm:"size:10;uniqueIndex:uk_uid_date_hour;index:idx_uhs_date;not null;column:date" json:"date"`
	Hour         int    `gorm:"uniqueIndex:uk_uid_date_hour;not null;column:hour" json:"hour"`
	CommandCount int    `gorm:"default:0;column:command_count" json:"command_count"`
	FirstSeenAt  string `gorm:"size:8;column:first_seen_at" json:"first_seen_at"`
	LastSeenAt   string `gorm:"size:8;column:last_seen_at" json:"last_seen_at"`
}

func (UserHourlyStatsDO) TableName() string { return "user_hourly_stats" }

// GlobalHourlyStatsDO 全局每小时活跃统计
type GlobalHourlyStatsDO struct {
	ID uint `gorm:"primaryKey" json:"id"`

	Date          string `gorm:"size:10;uniqueIndex:uk_global_date_hour;index:idx_ghs_date;not null;column:date" json:"date"`
	Hour          int    `gorm:"uniqueIndex:uk_global_date_hour;not null;column:hour" json:"hour"`
	ActiveUsers   int    `gorm:"default:0;column:active_users" json:"active_users"`
	TotalCommands int    `gorm:"default:0;column:total_commands" json:"total_commands"`
}

func (GlobalHourlyStatsDO) TableName() string { return "global_hourly_stats" }
