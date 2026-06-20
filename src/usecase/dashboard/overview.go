package dashboard

import (
	"time"
	"yonghemolimis/src/dao/analytics"
	"yonghemolimis/src/dao/game"
)

// OverviewData 仪表盘概览数据
type OverviewData struct {
	TotalUsers          int64   `json:"total_users"`
	VipUsers            int64   `json:"vip_users"`
	VipRate             float64 `json:"vip_rate"`
	PayingUsers         int64   `json:"paying_users"`
	PayRate             float64 `json:"pay_rate"`
	TotalRecharge       int64   `json:"total_recharge"`
	TotalRechargeYuan   float64 `json:"total_recharge_yuan"`
	ARPU                float64 `json:"arpu"`
	ARPPU               float64 `json:"arppu"`
	NewUsers7d          int64   `json:"new_users_7d"`
	HighChurnCount      int64   `json:"high_churn_count"`
	StuckCount          int64   `json:"stuck_count"`
	ResourceAlertCnt    int64   `json:"resource_alert_count"`
	SnapshotLatestDate  string  `json:"snapshot_latest_date"`
	ProfileCalculatedAt string  `json:"profile_calculated_at"`
	FeatureLastEventAt  string  `json:"feature_last_event_at"`
}

func formatNullableTime(t *time.Time) string {
	if t == nil || t.IsZero() {
		return ""
	}
	return t.Format("2006-01-02 15:04:05")
}

// GetOverview 获取仪表盘概览
func GetOverview(source string) (*OverviewData, error) {
	data := &OverviewData{}

	totalUsers, err := game.GetTotalUserCount()
	if err != nil {
		return nil, err
	}
	data.TotalUsers = totalUsers

	// VIP用户（vip_level > 0，可能来自系统活动奖励）
	vipUsers, _ := game.GetVipUserCount()
	data.VipUsers = vipUsers
	if totalUsers > 0 {
		data.VipRate = float64(vipUsers) / float64(totalUsers) * 100
	}

	// 付费用户（有真实充值订单）
	payingUsers, _ := game.GetRealPayingUserCount(source)
	data.PayingUsers = payingUsers

	if totalUsers > 0 {
		data.PayRate = float64(payingUsers) / float64(totalUsers) * 100
	}

	totalRecharge, err := game.GetRechargeRevenueTotal(source)
	if err != nil {
		return nil, err
	}
	data.TotalRecharge = totalRecharge
	data.TotalRechargeYuan = float64(totalRecharge) / 100.0

	// ARPU & ARPPU (分→元)
	if totalUsers > 0 {
		data.ARPU = float64(totalRecharge) / 100.0 / float64(totalUsers)
	}
	if payingUsers > 0 {
		data.ARPPU = float64(totalRecharge) / 100.0 / float64(payingUsers)
	}

	newUsers, err := game.GetNewUsersCount(7)
	if err != nil {
		return nil, err
	}
	data.NewUsers7d = newUsers

	stats, err := analytics.GetProfileStats()
	if err == nil && stats != nil {
		data.HighChurnCount = stats.HighChurnCount
		data.StuckCount = stats.StuckCount
		data.ResourceAlertCnt = stats.ResourceAlertCount
	}

	if latestSnapshot, err := analytics.GetLatestSnapshotDateGlobal(); err == nil {
		data.SnapshotLatestDate = latestSnapshot
	}
	if latestProfile, err := analytics.GetLatestProfileCalculatedAt(); err == nil {
		data.ProfileCalculatedAt = formatNullableTime(latestProfile)
	}
	if latestFeatureEvent, err := analytics.GetLatestFeatureEventTime(); err == nil {
		data.FeatureLastEventAt = formatNullableTime(latestFeatureEvent)
	}

	return data, nil
}

// GetDistribution 获取画像字段分布
func GetDistribution(field string) ([]analytics.DistItem, error) {
	return analytics.GetProfileDistribution(field)
}

// GetDAUTrend 获取 DAU 趋势
func GetDAUTrend(startDate, endDate string) ([]analytics.DAUItem, error) {
	return analytics.GetDAUByDateRange(startDate, endDate)
}

// GetRealmDistribution 获取境界分布
func GetRealmDistribution() ([]game.RealmDistItem, error) {
	return game.GetRealmDistribution()
}
