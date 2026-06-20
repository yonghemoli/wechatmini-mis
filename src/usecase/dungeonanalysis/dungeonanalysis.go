package dungeonyonghemolimis

import "yonghemolimis/src/dao/game"

type DungeonAnalysisData struct {
	Stats      []game.DungeonStatsItem  `json:"stats"`
	DailyTrend []game.DungeonDailyTrend `json:"daily_trend"`
}

func GetDungeonAnalysis(days int) (*DungeonAnalysisData, error) {
	stats, err := game.GetDungeonStats(days)
	if err != nil {
		return nil, err
	}
	trend, err := game.GetDungeonDailyTrend(days)
	if err != nil {
		return nil, err
	}
	return &DungeonAnalysisData{Stats: stats, DailyTrend: trend}, nil
}

func GetDungeonRealmDist(days int) ([]game.DungeonRealmDistItem, error) {
	return game.GetDungeonRealmDist(days)
}
