package gamestat

import (
	"yonghemolimis/src/dao/game"
)

// EnhancedOverview 增强版概览数据
type EnhancedOverview struct {
	TotalUsers    int64   `json:"total_users"`
	VipUsers      int64   `json:"vip_users"`
	VipRate       float64 `json:"vip_rate"`
	PayingUsers   int64   `json:"paying_users"`
	PayRate       float64 `json:"pay_rate"`
	TotalRecharge int64   `json:"total_recharge"`
	ARPU          float64 `json:"arpu"`
	ARPPU         float64 `json:"arppu"`
	NewUsers7d    int64   `json:"new_users_7d"`
	NewUsers1d    int64   `json:"new_users_1d"`
	DAU           int     `json:"dau"`
	MAU           int     `json:"mau"`
	DAUMAURatio   float64 `json:"dau_mau_ratio"`
}

// GetEnhancedOverview 获取增强版概览
func GetEnhancedOverview() (*EnhancedOverview, error) {
	data := &EnhancedOverview{}
	totalUsers, err := game.GetTotalUserCount()
	if err != nil {
		return nil, err
	}
	data.TotalUsers = totalUsers

	// VIP用户
	vipUsers, _ := game.GetVipUserCount()
	data.VipUsers = vipUsers
	if totalUsers > 0 {
		data.VipRate = float64(vipUsers) / float64(totalUsers) * 100
	}

	// 付费用户（有充值订单）
	payingUsers, _ := game.GetRealPayingUserCount()
	data.PayingUsers = payingUsers
	if totalUsers > 0 {
		data.PayRate = float64(payingUsers) / float64(totalUsers) * 100
	}
	totalRecharge, err := game.GetTotalRechargeSum()
	if err != nil {
		return nil, err
	}
	data.TotalRecharge = totalRecharge
	if totalUsers > 0 {
		data.ARPU = float64(totalRecharge) / 100.0 / float64(totalUsers)
	}
	if payingUsers > 0 {
		data.ARPPU = float64(totalRecharge) / 100.0 / float64(payingUsers)
	}
	newUsers7d, _ := game.GetNewUsersCount(7)
	data.NewUsers7d = newUsers7d
	newUsers1d, _ := game.GetNewUsersCount(1)
	data.NewUsers1d = newUsers1d
	dau, _ := game.GetMAUFromGame(1)
	data.DAU = dau
	mau, _ := game.GetMAUFromGame(30)
	data.MAU = mau
	if mau > 0 {
		data.DAUMAURatio = float64(dau) / float64(mau) * 100
	}
	return data, nil
}

// GetPlayerRanking 玩家排行
func GetPlayerRanking(sortBy string, limit int) ([]game.PlayerRankItem, error) {
	return game.GetPlayerRanking(sortBy, limit)
}

// GetGuildRanking 宗门排行
func GetGuildRanking(sortBy string, limit int) ([]game.GuildRankItem, error) {
	return game.GetGuildRanking(sortBy, limit)
}

// GetNewUsersTrend 新增玩家每日趋势
func GetNewUsersTrend(days int) ([]game.DailyCountItem, error) {
	return game.GetNewUsersTrend(days)
}

// GetRealmStageDistribution 大境界分布
func GetRealmStageDistribution() ([]game.RealmStageDistItem, error) {
	return game.GetRealmStageDistribution()
}

// GetRealmChurn 境界流失率
func GetRealmChurn(inactiveDays int) ([]game.RealmChurnItem, error) {
	return game.GetRealmChurn(inactiveDays)
}

// PaymentOverview 付费概览
type PaymentOverview struct {
	TotalRevenue    int64              `json:"total_revenue"`
	TotalOrders     int                `json:"total_orders"`
	PayingUsers     int64              `json:"paying_users"`
	TotalUsers      int64              `json:"total_users"`
	PayRate         float64            `json:"pay_rate"`
	ARPU            float64            `json:"arpu"`
	ARPPU           float64            `json:"arppu"`
	Revenue7d       int64              `json:"revenue_7d"`
	Revenue30d      int64              `json:"revenue_30d"`
	VipDistribution []game.VipDistItem `json:"vip_distribution"`
}

// GetPaymentOverview 付费分析概览
func GetPaymentOverview(filter string) (*PaymentOverview, error) {
	data := &PaymentOverview{}
	totalUsers, _ := game.GetTotalUserCount()
	payingUsers, _ := game.GetRealPayingUserCount(filter)
	totalRecharge, _ := game.GetRechargeRevenueTotal(filter)
	data.TotalUsers = totalUsers
	data.PayingUsers = payingUsers
	data.TotalRevenue = totalRecharge
	if totalUsers > 0 {
		data.PayRate = float64(payingUsers) / float64(totalUsers) * 100
	}
	if totalUsers > 0 {
		data.ARPU = float64(totalRecharge) / 100.0 / float64(totalUsers)
	}
	if payingUsers > 0 {
		data.ARPPU = float64(totalRecharge) / 100.0 / float64(payingUsers)
	}
	rev7d, _ := game.GetRechargeRevenueTrend(7, filter)
	for _, r := range rev7d {
		data.Revenue7d += r.Revenue
	}
	rev30d, _ := game.GetRechargeRevenueTrend(30, filter)
	for _, r := range rev30d {
		data.Revenue30d += r.Revenue
		data.TotalOrders += r.OrderCount
	}
	vipDist, _ := game.GetVipDistribution()
	data.VipDistribution = vipDist
	return data, nil
}

// GetRevenueTrend 营收趋势
func GetRevenueTrend(days int, filter string) ([]game.DailyRevenueItem, error) {
	return game.GetRechargeRevenueTrend(days, filter)
}

// GetPackageStats 套餐销售分析
func GetPackageStats(days int, filter string) ([]game.PackageStatItem, error) {
	return game.GetRechargePackageStats(days, filter)
}
