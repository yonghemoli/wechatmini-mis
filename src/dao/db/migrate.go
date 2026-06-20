package db

import "log"

// AutoMigrate 自动迁移分析库表
func AutoMigrate() error {
	if Get() == nil {
		return nil
	}
	err := Get().AutoMigrate(
		&UserDO{},
		&ActionLogDO{},
		&DailySnapshotDO{},
		&UserProfileDO{},
		&UserProfileTagDO{},
		&SegmentDO{},
		&FeatureEventDO{},
		&FeatureDailyStatsDO{},
		&FeatureScoreDO{},
		&GameActionEventDO{},
		&UserHourlyStatsDO{},
		&GlobalHourlyStatsDO{},
	)
	if err != nil {
		return err
	}

	// 没有分群时，插入默认分群
	seedDefaultSegments()

	return nil
}

// seedDefaultSegments 当 segments 表为空时，插入预设的常用分群
func seedDefaultSegments() {
	var count int64
	Get().Model(&SegmentDO{}).Count(&count)
	if count > 0 {
		return // 已有数据，不插入
	}

	defaults := []SegmentDO{
		{
			Name:        "高流失风险玩家",
			Description: "流失风险 ≥ 70% 的玩家，需重点关注和挽留",
			RulesJSON:   `{"logic":"AND","rules":[{"field":"churn_risk","operator":"gte","value":70}]}`,
			CreatedBy:   "system",
		},
		{
			Name:        "高流失大氪玩家",
			Description: "付费等级为大氪或巨鲸且流失风险 ≥ 60%，优先挽留",
			RulesJSON:   `{"logic":"AND","rules":[{"field":"pay_tier","operator":"in","value":["WHALE","LEVIATHAN"]},{"field":"churn_risk","operator":"gte","value":60}]}`,
			CreatedBy:   "system",
		},
		{
			Name:        "新手卡关玩家",
			Description: "新手期且处于卡关状态，需优化引导",
			RulesJSON:   `{"logic":"AND","rules":[{"field":"lifecycle_stage","operator":"eq","value":"NEW"},{"field":"stuck_flag","operator":"eq","value":true}]}`,
			CreatedBy:   "system",
		},
		{
			Name:        "流失玩家",
			Description: "长期未登录的流失期玩家",
			RulesJSON:   `{"logic":"AND","rules":[{"field":"lifecycle_stage","operator":"eq","value":"LOST"}]}`,
			CreatedBy:   "system",
		},
		{
			Name:        "回流玩家",
			Description: "流失后重新回归的玩家，适合投放回流礼包",
			RulesJSON:   `{"logic":"AND","rules":[{"field":"lifecycle_stage","operator":"eq","value":"RETURNED"}]}`,
			CreatedBy:   "system",
		},
		{
			Name:        "免费活跃玩家",
			Description: "免费但活跃（流失风险 < 30%），有付费转化潜力",
			RulesJSON:   `{"logic":"AND","rules":[{"field":"pay_tier","operator":"eq","value":"FREE"},{"field":"churn_risk","operator":"lt","value":30}]}`,
			CreatedBy:   "system",
		},
		{
			Name:        "资源告急玩家",
			Description: "触发资源预警的玩家，可能需要补给活动",
			RulesJSON:   `{"logic":"AND","rules":[{"field":"resource_alert","operator":"eq","value":true}]}`,
			CreatedBy:   "system",
		},
		{
			Name:        "成熟期战斗型玩家",
			Description: "成熟期且偏好战斗的核心玩家",
			RulesJSON:   `{"logic":"AND","rules":[{"field":"lifecycle_stage","operator":"eq","value":"MATURE"},{"field":"play_style","operator":"eq","value":"COMBAT"}]}`,
			CreatedBy:   "system",
		},
	}

	for i := range defaults {
		if err := Get().Create(&defaults[i]).Error; err != nil {
			log.Printf("[分群] 插入默认分群 '%s' 失败: %v", defaults[i].Name, err)
		}
	}
	log.Printf("[分群] 已插入 %d 个默认分群", len(defaults))
}
