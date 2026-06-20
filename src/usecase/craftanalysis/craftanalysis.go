package craftyonghemolimis

import "yonghemolimis/src/dao/game"

type CraftAnalysisData struct {
	Stats     *game.CraftStatsItem       `json:"stats"`
	LevelDist []game.CraftProfessionItem `json:"level_dist"`
}

func GetCraftAnalysis() (*CraftAnalysisData, error) {
	stats, err := game.GetCraftStats()
	if err != nil {
		return nil, err
	}
	levelDist, err := game.GetCraftLevelDistribution()
	if err != nil {
		return nil, err
	}
	return &CraftAnalysisData{
		Stats:     stats,
		LevelDist: levelDist,
	}, nil
}
