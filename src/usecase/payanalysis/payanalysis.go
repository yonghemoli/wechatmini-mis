package payyonghemolimis

import "yonghemolimis/src/dao/game"

type PayAnalysisData struct {
	FirstRechargeConversion []game.RechargeConversionItem `json:"first_recharge_conversion"`
	RevenueTrend            []game.DailyRevenueItem       `json:"revenue_trend"`
	PackageStats            []game.PackageStatItem        `json:"package_stats"`
	TierTrend               []game.RechargeTierTrendItem  `json:"tier_trend"`
	PassStats               *game.PassSubscriptionStats   `json:"pass_stats"`
	VipDist                 []game.VipDistItem            `json:"vip_dist"`
}

func GetPayAnalysis(days int) (*PayAnalysisData, error) {
	fc, err := game.GetFirstRechargeConversion()
	if err != nil {
		return nil, err
	}
	rt, err := game.GetRechargeRevenueTrend(days)
	if err != nil {
		return nil, err
	}
	ps, err := game.GetRechargePackageStats(days)
	if err != nil {
		return nil, err
	}
	tt, err := game.GetRechargeTierTrend(days)
	if err != nil {
		return nil, err
	}
	pass, err := game.GetPassSubscriptionStats()
	if err != nil {
		return nil, err
	}
	vip, err := game.GetVipDistribution()
	if err != nil {
		return nil, err
	}
	return &PayAnalysisData{
		FirstRechargeConversion: fc,
		RevenueTrend:            rt,
		PackageStats:            ps,
		TierTrend:               tt,
		PassStats:               pass,
		VipDist:                 vip,
	}, nil
}
