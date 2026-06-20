package analytics

import (
	"time"
	"yonghemolimis/src/dao/db"
)

// ========== 分析库 CRUD ==========

// UpsertProfile 创建或更新用户画像
func UpsertProfile(p *db.UserProfileDO) error {
	return db.Get().Save(p).Error
}

// GetProfile 获取单个用户画像
func GetProfile(uid uint) (*db.UserProfileDO, error) {
	var p db.UserProfileDO
	err := db.Get().Where("uid = ?", uid).First(&p).Error
	if err != nil {
		return nil, err
	}
	return &p, nil
}

// ProfileFilter 画像筛选条件
type ProfileFilter struct {
	LifecycleStage string
	PayTier        string
	PlayStyle      string
	MinChurnRisk   int
	StuckOnly      bool
	Page           int
	PageSize       int
}

// GetProfilesByFilter 根据筛选条件查询画像列表
func GetProfilesByFilter(f ProfileFilter) ([]db.UserProfileDO, int64, error) {
	query := db.Get().Model(&db.UserProfileDO{})

	if f.LifecycleStage != "" {
		query = query.Where("lifecycle_stage = ?", f.LifecycleStage)
	}
	if f.PayTier != "" {
		query = query.Where("pay_tier = ?", f.PayTier)
	}
	if f.PlayStyle != "" {
		query = query.Where("play_style = ?", f.PlayStyle)
	}
	if f.MinChurnRisk > 0 {
		query = query.Where("churn_risk >= ?", f.MinChurnRisk)
	}
	if f.StuckOnly {
		query = query.Where("stuck_flag = ?", true)
	}

	var total int64
	query.Count(&total)

	page := f.Page
	if page < 1 {
		page = 1
	}
	size := f.PageSize
	if size <= 0 || size > 100 {
		size = 20
	}

	var list []db.UserProfileDO
	err := query.Offset((page - 1) * size).Limit(size).Order("churn_risk DESC").Find(&list).Error
	return list, total, err
}

// DistItem 分布统计项
type DistItem struct {
	Label string `gorm:"column:label" json:"label"`
	Count int    `gorm:"column:count" json:"count"`
}

// GetProfileDistribution 获取画像字段分布
func GetProfileDistribution(field string) ([]DistItem, error) {
	allowed := map[string]bool{
		"lifecycle_stage": true,
		"pay_tier":        true,
		"play_style":      true,
		"social_type":     true,
	}
	if !allowed[field] {
		field = "lifecycle_stage"
	}

	var items []DistItem
	err := db.Get().Model(&db.UserProfileDO{}).
		Select(field + " as label, COUNT(*) as count").
		Group(field).
		Order("count DESC").
		Find(&items).Error
	return items, err
}

// SaveSnapshot 保存每日快照
func SaveSnapshot(s *db.DailySnapshotDO) error {
	return db.Get().Save(s).Error
}

// GetSnapshotsByUID 获取指定用户的快照列表
func GetSnapshotsByUID(uid uint, limit int) ([]db.DailySnapshotDO, error) {
	if limit <= 0 {
		limit = 30
	}
	var list []db.DailySnapshotDO
	err := db.Get().Where("uid = ?", uid).Order("snapshot_date DESC").Limit(limit).Find(&list).Error
	return list, err
}

// DAUItem DAU 统计项
type DAUItem struct {
	Date  string `gorm:"column:date"  json:"date"`
	Count int    `gorm:"column:count" json:"count"`
}

// GetDAUByDateRange 按日期范围统计 DAU
func GetDAUByDateRange(startDate, endDate string) ([]DAUItem, error) {
	var items []DAUItem
	err := db.Get().Model(&db.DailySnapshotDO{}).
		Select("snapshot_date as date, COUNT(DISTINCT uid) as count").
		Where("snapshot_date BETWEEN ? AND ?", startDate, endDate).
		Group("snapshot_date").
		Order("snapshot_date ASC").
		Find(&items).Error
	return items, err
}

// UpsertTag 创建或更新用户标签
func UpsertTag(t *db.UserProfileTagDO) error {
	return db.Get().Save(t).Error
}

// GetTagsByUID 获取指定用户的所有标签
func GetTagsByUID(uid uint) ([]db.UserProfileTagDO, error) {
	var tags []db.UserProfileTagDO
	err := db.Get().Where("uid = ?", uid).Find(&tags).Error
	return tags, err
}

// CreateSegment 创建分群
func CreateSegment(s *db.SegmentDO) error {
	return db.Get().Create(s).Error
}

// GetSegments 获取所有分群
func GetSegments() ([]db.SegmentDO, error) {
	var segs []db.SegmentDO
	err := db.Get().Order("created_at DESC").Find(&segs).Error
	return segs, err
}

// DeleteSegment 删除分群
func DeleteSegment(id uint) error {
	return db.Get().Delete(&db.SegmentDO{}, id).Error
}

// GetSegmentByID 获取单个分群
func GetSegmentByID(id uint) (*db.SegmentDO, error) {
	var seg db.SegmentDO
	err := db.Get().First(&seg, id).Error
	if err != nil {
		return nil, err
	}
	return &seg, nil
}

// UpdateSegmentCount 更新分群人数
func UpdateSegmentCount(id uint, count int) error {
	return db.Get().Model(&db.SegmentDO{}).Where("id = ?", id).Update("user_count", count).Error
}

// SaveActionLog 保存行为日志
func SaveActionLog(log *db.ActionLogDO) error {
	return db.Get().Create(log).Error
}

// GetHighChurnProfiles 获取高流失风险画像（churn_risk >= threshold）
func GetHighChurnProfiles(threshold int, limit int) ([]db.UserProfileDO, error) {
	if limit <= 0 {
		limit = 50
	}
	var list []db.UserProfileDO
	err := db.Get().Where("churn_risk >= ?", threshold).
		Order("churn_risk DESC").Limit(limit).
		Find(&list).Error
	return list, err
}

// GetStuckProfiles 获取卡关用户
func GetStuckProfiles(limit int) ([]db.UserProfileDO, error) {
	if limit <= 0 {
		limit = 50
	}
	var list []db.UserProfileDO
	err := db.Get().Where("stuck_flag = ?", true).
		Order("churn_risk DESC").Limit(limit).
		Find(&list).Error
	return list, err
}

// GetResourceAlertProfiles 获取资源告急用户
func GetResourceAlertProfiles(limit int) ([]db.UserProfileDO, error) {
	if limit <= 0 {
		limit = 50
	}
	var list []db.UserProfileDO
	err := db.Get().Where("resource_alert = ?", true).
		Order("churn_risk DESC").Limit(limit).
		Find(&list).Error
	return list, err
}

// GetNewbieAtRiskProfiles 获取新手流失预警用户（新手期 + 流失风险 >= threshold）
func GetNewbieAtRiskProfiles(threshold int, limit int) ([]db.UserProfileDO, error) {
	if limit <= 0 {
		limit = 50
	}
	if threshold <= 0 {
		threshold = 40
	}
	var list []db.UserProfileDO
	err := db.Get().Where("lifecycle_stage = 'NEW' AND churn_risk >= ?", threshold).
		Order("churn_risk DESC").Limit(limit).
		Find(&list).Error
	return list, err
}

// GetWhaleChurnProfiles 获取大氪沉默预警（大氪/巨鲸 + 流失风险 >= threshold）
func GetWhaleChurnProfiles(threshold int, limit int) ([]db.UserProfileDO, error) {
	if limit <= 0 {
		limit = 50
	}
	if threshold <= 0 {
		threshold = 50
	}
	var list []db.UserProfileDO
	err := db.Get().Where("pay_tier IN ('WHALE','LEVIATHAN') AND churn_risk >= ?", threshold).
		Order("churn_risk DESC").Limit(limit).
		Find(&list).Error
	return list, err
}

// ProfileStats 画像统计摘要
type ProfileStats struct {
	TotalProfiles      int64 `json:"total"`
	HighChurnCount     int64 `json:"high_churn"`
	StuckCount         int64 `json:"stuck"`
	ResourceAlertCount int64 `json:"resource_alert"`
	NewbieAtRisk       int64 `json:"newbie_at_risk"`
	WhaleChurn         int64 `json:"whale_churn"`
	PayingCount        int64 `json:"paying"`
}

// GetProfileStats 获取画像统计数据
func GetProfileStats() (*ProfileStats, error) {
	var stats ProfileStats
	d := db.Get()

	d.Model(&db.UserProfileDO{}).Count(&stats.TotalProfiles)
	d.Model(&db.UserProfileDO{}).Where("churn_risk >= 60").Count(&stats.HighChurnCount)
	d.Model(&db.UserProfileDO{}).Where("stuck_flag = ?", true).Count(&stats.StuckCount)
	d.Model(&db.UserProfileDO{}).Where("resource_alert = ?", true).Count(&stats.ResourceAlertCount)
	d.Model(&db.UserProfileDO{}).Where("lifecycle_stage = 'NEW' AND churn_risk >= 40").Count(&stats.NewbieAtRisk)
	d.Model(&db.UserProfileDO{}).Where("pay_tier IN ('WHALE','LEVIATHAN') AND churn_risk >= 50").Count(&stats.WhaleChurn)
	d.Model(&db.UserProfileDO{}).Where("pay_tier != 'FREE'").Count(&stats.PayingCount)

	return &stats, nil
}

// BatchUpsertSnapshots 批量保存快照
func BatchUpsertSnapshots(snapshots []db.DailySnapshotDO) error {
	if len(snapshots) == 0 {
		return nil
	}
	return db.Get().Save(&snapshots).Error
}

// GetLatestSnapshotDate 获取最近快照日期
func GetLatestSnapshotDate(uid uint) (string, error) {
	var snap db.DailySnapshotDO
	err := db.Get().Where("uid = ?", uid).Order("snapshot_date DESC").First(&snap).Error
	if err != nil {
		return "", err
	}
	return snap.SnapshotDate, nil
}

// GetLatestSnapshotDateGlobal 获取全局最近快照日期
func GetLatestSnapshotDateGlobal() (string, error) {
	type latestRow struct {
		SnapshotDate string `gorm:"column:snapshot_date"`
	}
	var row latestRow
	err := db.Get().Model(&db.DailySnapshotDO{}).
		Select("MAX(snapshot_date) as snapshot_date").
		Scan(&row).Error
	if err != nil {
		return "", err
	}
	return row.SnapshotDate, nil
}

// GetLatestProfileCalculatedAt 获取最近一次画像计算时间
func GetLatestProfileCalculatedAt() (*time.Time, error) {
	type latestRow struct {
		CalculatedAt *time.Time `gorm:"column:calculated_at"`
	}
	var row latestRow
	err := db.Get().Model(&db.UserProfileDO{}).
		Select("MAX(last_calculated_at) as calculated_at").
		Scan(&row).Error
	if err != nil {
		return nil, err
	}
	return row.CalculatedAt, nil
}

// GetLatestFeatureEventTime 获取最近一条功能埋点时间
func GetLatestFeatureEventTime() (*time.Time, error) {
	type latestRow struct {
		EventTime *time.Time `gorm:"column:event_time"`
	}
	var row latestRow
	err := db.Get().Model(&db.FeatureEventDO{}).
		Select("MAX(event_time) as event_time").
		Scan(&row).Error
	if err != nil {
		return nil, err
	}
	return row.EventTime, nil
}

// CountSnapshotsByDateRange 统计日期范围内的活跃天数
func CountSnapshotsByDateRange(uid uint, start, end time.Time) (int64, error) {
	var count int64
	err := db.Get().Model(&db.DailySnapshotDO{}).
		Where("uid = ? AND snapshot_date BETWEEN ? AND ?", uid, start.Format("2006-01-02"), end.Format("2006-01-02")).
		Count(&count).Error
	return count, err
}
