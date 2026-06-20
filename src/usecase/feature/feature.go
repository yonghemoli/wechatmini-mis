package feature

import (
	"time"
	"yonghemolimis/src/dao/db"
)

// ==================== 原始事件查询 ====================

// GetEventCountByDate 获取某日事件总数
func GetEventCountByDate(date string) (int64, error) {
	var count int64
	err := db.Get().Model(&db.FeatureEventDO{}).
		Where("DATE(event_time) = ?", date).
		Count(&count).Error
	return count, err
}

// ==================== 每日聚合 ====================

// AggregateDailyStats 对指定日期的原始事件进行聚合，写入 feature_daily_stats
func AggregateDailyStats(date string) error {
	d := db.Get()

	// 先查出该日所有功能的聚合数据
	type aggRow struct {
		FeatureName     string  `gorm:"column:feature_name"`
		FeatureCategory string  `gorm:"column:feature_category"`
		TotalUses       int     `gorm:"column:total_uses"`
		UniqueUsers     int     `gorm:"column:unique_users"`
		AvgResponseMs   float64 `gorm:"column:avg_response_ms"`
		SuccessRate     float64 `gorm:"column:success_rate"`
	}

	var rows []aggRow
	err := d.Model(&db.FeatureEventDO{}).
		Select(`feature_name, feature_category, 
			COUNT(*) as total_uses, 
			COUNT(DISTINCT uid) as unique_users,
			AVG(response_time_ms) as avg_response_ms,
			AVG(CASE WHEN success THEN 1.0 ELSE 0.0 END) as success_rate`).
		Where("DATE(event_time) = ?", date).
		Group("feature_name, feature_category").
		Find(&rows).Error
	if err != nil {
		return err
	}

	// Upsert 写入每日统计表
	for _, r := range rows {
		stat := db.FeatureDailyStatsDO{
			Date:            date,
			FeatureName:     r.FeatureName,
			FeatureCategory: r.FeatureCategory,
			TotalUses:       r.TotalUses,
			UniqueUsers:     r.UniqueUsers,
			AvgResponseMs:   r.AvgResponseMs,
			SuccessRate:     r.SuccessRate,
		}
		// Upsert: 按 date + feature_name 唯一键
		result := d.Where("date = ? AND feature_name = ?", date, r.FeatureName).
			Assign(stat).FirstOrCreate(&stat)
		if result.Error != nil {
			return result.Error
		}
	}
	return nil
}

// ==================== 评分计算 ====================

// CalculateFeatureScores 计算所有功能的质量评分
func CalculateFeatureScores(totalUsers int64) error {
	d := db.Get()
	now := time.Now()
	date7 := now.AddDate(0, 0, -7).Format("2006-01-02")
	date30 := now.AddDate(0, 0, -30).Format("2006-01-02")
	date14 := now.AddDate(0, 0, -14).Format("2006-01-02")
	today := now.Format("2006-01-02")

	// 获取所有出现过的功能
	type featureRow struct {
		FeatureName     string `gorm:"column:feature_name"`
		FeatureCategory string `gorm:"column:feature_category"`
	}
	var features []featureRow
	err := d.Model(&db.FeatureDailyStatsDO{}).
		Select("DISTINCT feature_name, feature_category").
		Find(&features).Error
	if err != nil {
		return err
	}

	for _, f := range features {
		var score db.FeatureScoreDO

		// 7天数据
		type sumRow struct {
			TotalUses   int     `gorm:"column:total_uses"`
			UniqueUsers int     `gorm:"column:unique_users"`
			AvgResp     float64 `gorm:"column:avg_resp"`
			SuccRate    float64 `gorm:"column:succ_rate"`
			Days        int     `gorm:"column:days"`
		}

		var s7, s30 sumRow

		d.Model(&db.FeatureDailyStatsDO{}).
			Select(`COALESCE(SUM(total_uses),0) as total_uses, 
				COALESCE(MAX(unique_users),0) as unique_users,
				COALESCE(AVG(avg_response_ms),0) as avg_resp,
				COALESCE(AVG(success_rate),1) as succ_rate,
				COUNT(*) as days`).
			Where("feature_name = ? AND date >= ? AND date <= ?", f.FeatureName, date7, today).
			Scan(&s7)

		d.Model(&db.FeatureDailyStatsDO{}).
			Select(`COALESCE(SUM(total_uses),0) as total_uses, 
				COALESCE(MAX(unique_users),0) as unique_users,
				COALESCE(AVG(avg_response_ms),0) as avg_resp,
				COALESCE(AVG(success_rate),1) as succ_rate,
				COUNT(*) as days`).
			Where("feature_name = ? AND date >= ? AND date <= ?", f.FeatureName, date30, today).
			Scan(&s30)

		// 环比增长：近7天用量 vs 前7天（7-14天前）用量
		var prevWeek sumRow
		d.Model(&db.FeatureDailyStatsDO{}).
			Select("COALESCE(SUM(total_uses),0) as total_uses").
			Where("feature_name = ? AND date >= ? AND date < ?", f.FeatureName, date14, date7).
			Scan(&prevWeek)

		growth := 0.0
		if prevWeek.TotalUses > 0 {
			growth = float64(s7.TotalUses-prevWeek.TotalUses) / float64(prevWeek.TotalUses)
		}

		// 用户渗透率
		penetration := 0.0
		if totalUsers > 0 {
			penetration = float64(s30.UniqueUsers) / float64(totalUsers)
		}

		// 日均使用次数
		avgDaily := 0.0
		if s30.Days > 0 {
			avgDaily = float64(s30.TotalUses) / float64(s30.Days)
		}

		// 综合评分: 渗透率(30%) + 日均使用(25%) + 增长率(20%) + 成功率(15%) + 性能(10%)
		// 各维度归一化到 0-100
		penScore := penetration * 100                  // 已经是 0~1
		dailyScore := min100(avgDaily / 10 * 100)      // 日均10次计满分
		growthScore := min100((growth + 1) * 50)       // -100%~+100% → 0~100
		succScore := s30.SuccRate * 100                // 已经是 0~1
		perfScore := min100((1000 - s30.AvgResp) / 10) // <1000ms越低越好

		qualityScore := penScore*0.30 + dailyScore*0.25 + growthScore*0.20 + succScore*0.15 + perfScore*0.10

		score.FeatureName = f.FeatureName
		score.FeatureCategory = f.FeatureCategory
		score.TotalUses7d = s7.TotalUses
		score.TotalUses30d = s30.TotalUses
		score.UniqueUsers7d = s7.UniqueUsers
		score.UniqueUsers30d = s30.UniqueUsers
		score.AvgDailyUses = avgDaily
		score.AvgResponseMs = s30.AvgResp
		score.SuccessRate = s30.SuccRate
		score.UsageGrowth = growth
		score.UserPenetration = penetration
		score.QualityScore = qualityScore

		d.Where("feature_name = ?", f.FeatureName).Assign(score).FirstOrCreate(&score)
	}
	return nil
}

func min100(v float64) float64 {
	if v > 100 {
		return 100
	}
	if v < 0 {
		return 0
	}
	return v
}

// ==================== 查询接口 ====================

// GetFeatureScores 获取所有功能评分，按分数降序
func GetFeatureScores(category string) ([]db.FeatureScoreDO, error) {
	d := db.Get()
	q := d.Model(&db.FeatureScoreDO{}).Order("quality_score DESC")
	if category != "" {
		q = q.Where("feature_category = ?", category)
	}
	var rows []db.FeatureScoreDO
	err := q.Find(&rows).Error
	return rows, err
}

// GetFeatureDailyTrend 获取某功能的每日趋势
func GetFeatureDailyTrend(featureName string, days int) ([]db.FeatureDailyStatsDO, error) {
	d := db.Get()
	since := time.Now().AddDate(0, 0, -days).Format("2006-01-02")
	var rows []db.FeatureDailyStatsDO
	err := d.Where("feature_name = ? AND date >= ?", featureName, since).
		Order("date ASC").Find(&rows).Error
	return rows, err
}

// GetCategoryOverview 按分类聚合统计
func GetCategoryOverview(days int) ([]CategoryStat, error) {
	d := db.Get()
	since := time.Now().AddDate(0, 0, -days).Format("2006-01-02")

	var rows []CategoryStat
	err := d.Model(&db.FeatureDailyStatsDO{}).
		Select(`feature_category, 
			SUM(total_uses) as total_uses,
			MAX(unique_users) as unique_users,
			AVG(avg_response_ms) as avg_response_ms,
			COUNT(DISTINCT feature_name) as feature_count`).
		Where("date >= ?", since).
		Group("feature_category").
		Order("total_uses DESC").
		Find(&rows).Error
	return rows, err
}

// CategoryStat 分类统计
type CategoryStat struct {
	FeatureCategory string  `gorm:"column:feature_category" json:"feature_category"`
	TotalUses       int     `gorm:"column:total_uses" json:"total_uses"`
	UniqueUsers     int     `gorm:"column:unique_users" json:"unique_users"`
	AvgResponseMs   float64 `gorm:"column:avg_response_ms" json:"avg_response_ms"`
	FeatureCount    int     `gorm:"column:feature_count" json:"feature_count"`
}

// GetTopFeatures 热门功能 TopN
func GetTopFeatures(days, limit int) ([]db.FeatureDailyStatsDO, error) {
	d := db.Get()
	since := time.Now().AddDate(0, 0, -days).Format("2006-01-02")

	type aggRow struct {
		FeatureName     string  `gorm:"column:feature_name"`
		FeatureCategory string  `gorm:"column:feature_category"`
		TotalUses       int     `gorm:"column:total_uses"`
		UniqueUsers     int     `gorm:"column:unique_users"`
		AvgResponseMs   float64 `gorm:"column:avg_response_ms"`
	}
	var rows []aggRow
	err := d.Model(&db.FeatureDailyStatsDO{}).
		Select(`feature_name, feature_category,
			SUM(total_uses) as total_uses, 
			MAX(unique_users) as unique_users,
			AVG(avg_response_ms) as avg_response_ms`).
		Where("date >= ?", since).
		Group("feature_name, feature_category").
		Order("total_uses DESC").
		Limit(limit).
		Find(&rows).Error
	if err != nil {
		return nil, err
	}

	result := make([]db.FeatureDailyStatsDO, len(rows))
	for i, r := range rows {
		result[i] = db.FeatureDailyStatsDO{
			FeatureName:     r.FeatureName,
			FeatureCategory: r.FeatureCategory,
			TotalUses:       r.TotalUses,
			UniqueUsers:     r.UniqueUsers,
			AvgResponseMs:   r.AvgResponseMs,
		}
	}
	return result, nil
}

// ==================== 场景/频道/平台维度分析 ====================

// SceneDistItem 场景分布项
type SceneDistItem struct {
	Scene       string `gorm:"column:scene" json:"scene"`
	TotalUses   int    `gorm:"column:total_uses" json:"total_uses"`
	UniqueUsers int    `gorm:"column:unique_users" json:"unique_users"`
}

// GetSceneDistribution 获取私聊/群聊场景分布
func GetSceneDistribution(days int) ([]SceneDistItem, error) {
	d := db.Get()
	since := time.Now().AddDate(0, 0, -days).Format("2006-01-02")
	var rows []SceneDistItem
	err := d.Model(&db.FeatureEventDO{}).
		Select(`COALESCE(NULLIF(scene,''),'unknown') as scene,
			COUNT(*) as total_uses,
			COUNT(DISTINCT uid) as unique_users`).
		Where("event_time >= ?", since).
		Group("scene").
		Order("total_uses DESC").
		Find(&rows).Error
	return rows, err
}

// ChannelTopItem 频道活跃排名项
type ChannelTopItem struct {
	ChannelID   string `gorm:"column:channel_id" json:"channel_id"`
	Scene       string `gorm:"column:scene" json:"scene"`
	TotalUses   int    `gorm:"column:total_uses" json:"total_uses"`
	UniqueUsers int    `gorm:"column:unique_users" json:"unique_users"`
}

// GetChannelTop 获取频道/群活跃排行
func GetChannelTop(days, limit int) ([]ChannelTopItem, error) {
	d := db.Get()
	since := time.Now().AddDate(0, 0, -days).Format("2006-01-02")
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	var rows []ChannelTopItem
	err := d.Model(&db.FeatureEventDO{}).
		Select(`channel_id, 
			MAX(scene) as scene,
			COUNT(*) as total_uses,
			COUNT(DISTINCT uid) as unique_users`).
		Where("event_time >= ? AND channel_id != ''", since).
		Group("channel_id").
		Order("total_uses DESC").
		Limit(limit).
		Find(&rows).Error
	return rows, err
}

// PlatformDistItem 平台分布项
type PlatformDistItem struct {
	Platform    string `gorm:"column:platform" json:"platform"`
	TotalUses   int    `gorm:"column:total_uses" json:"total_uses"`
	UniqueUsers int    `gorm:"column:unique_users" json:"unique_users"`
}

// GetPlatformDistribution 获取平台分布
func GetPlatformDistribution(days int) ([]PlatformDistItem, error) {
	d := db.Get()
	since := time.Now().AddDate(0, 0, -days).Format("2006-01-02")
	var rows []PlatformDistItem
	err := d.Model(&db.FeatureEventDO{}).
		Select(`COALESCE(NULLIF(platform,''),'unknown') as platform,
			COUNT(*) as total_uses,
			COUNT(DISTINCT uid) as unique_users`).
		Where("event_time >= ?", since).
		Group("platform").
		Order("total_uses DESC").
		Find(&rows).Error
	return rows, err
}

// SceneFeatureItem 各场景下的功能使用排名
type SceneFeatureItem struct {
	Scene           string `gorm:"column:scene" json:"scene"`
	FeatureName     string `gorm:"column:feature_name" json:"feature_name"`
	FeatureCategory string `gorm:"column:feature_category" json:"feature_category"`
	TotalUses       int    `gorm:"column:total_uses" json:"total_uses"`
	UniqueUsers     int    `gorm:"column:unique_users" json:"unique_users"`
}

// GetSceneFeatureTop 获取各场景下的功能使用 TopN
func GetSceneFeatureTop(scene string, days, limit int) ([]SceneFeatureItem, error) {
	d := db.Get()
	since := time.Now().AddDate(0, 0, -days).Format("2006-01-02")
	if limit <= 0 || limit > 100 {
		limit = 15
	}
	q := d.Model(&db.FeatureEventDO{}).
		Select(`COALESCE(NULLIF(scene,''),'unknown') as scene,
			feature_name, feature_category,
			COUNT(*) as total_uses,
			COUNT(DISTINCT uid) as unique_users`).
		Where("event_time >= ?", since)
	if scene != "" {
		q = q.Where("scene = ?", scene)
	}
	var rows []SceneFeatureItem
	err := q.Group("scene, feature_name, feature_category").
		Order("total_uses DESC").
		Limit(limit).
		Find(&rows).Error
	return rows, err
}

// SceneTrendItem 场景每日趋势
type SceneTrendItem struct {
	Date        string `gorm:"column:date" json:"date"`
	Scene       string `gorm:"column:scene" json:"scene"`
	TotalUses   int    `gorm:"column:total_uses" json:"total_uses"`
	UniqueUsers int    `gorm:"column:unique_users" json:"unique_users"`
}

// GetSceneTrend 获取私聊/群聊的每日使用趋势
func GetSceneTrend(days int) ([]SceneTrendItem, error) {
	d := db.Get()
	since := time.Now().AddDate(0, 0, -days).Format("2006-01-02")
	var rows []SceneTrendItem
	err := d.Model(&db.FeatureEventDO{}).
		Select(`DATE(event_time) as date,
			COALESCE(NULLIF(scene,''),'unknown') as scene,
			COUNT(*) as total_uses,
			COUNT(DISTINCT uid) as unique_users`).
		Where("event_time >= ?", since).
		Group("DATE(event_time), scene").
		Order("date ASC, scene").
		Find(&rows).Error
	return rows, err
}
