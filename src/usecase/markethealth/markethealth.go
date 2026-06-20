package markethealth

import "yonghemolimis/src/dao/game"

type MarketHealthData struct {
	TradeStats     *game.MarketTradeStatsItem     `json:"trade_stats"`
	HotItems       []game.MarketHotItemRow        `json:"hot_items"`
	DailyTrend     []game.MarketDailyTrendItem    `json:"daily_trend"`
	PriceAnomalies []game.MarketPriceAnomalyItem  `json:"price_anomalies"`
	LargeTrades    []game.MarketLargeTradeItem    `json:"large_trades"`
	PlayerProfiles []game.MarketPlayerProfileItem `json:"player_profiles"`
}

func GetMarketHealth(days int) (*MarketHealthData, error) {
	stats, err := game.GetMarketTradeStats(days)
	if err != nil {
		return nil, err
	}
	hot, err := game.GetMarketHotItems(days, 20)
	if err != nil {
		return nil, err
	}
	trend, err := game.GetMarketDailyTrend(days)
	if err != nil {
		return nil, err
	}
	anomalies, err := game.GetMarketPriceAnomalies(days)
	if err != nil {
		return nil, err
	}
	large, err := game.GetMarketLargeTrades(days, 50, 10000)
	if err != nil {
		return nil, err
	}
	profiles, err := game.GetMarketPlayerProfiles(days, 100)
	if err != nil {
		return nil, err
	}
	return &MarketHealthData{
		TradeStats:     stats,
		HotItems:       hot,
		DailyTrend:     trend,
		PriceAnomalies: anomalies,
		LargeTrades:    large,
		PlayerProfiles: profiles,
	}, nil
}
