package seclusion

import "yonghemolimis/src/dao/game"

type SeclusionData struct {
	Stats        *game.SeclusionStatsItem         `json:"stats"`
	DurationDist []game.SeclusionDurationDistItem `json:"duration_dist"`
}

func GetSeclusionAnalysis(days int) (*SeclusionData, error) {
	stats, err := game.GetSeclusionStats(days)
	if err != nil {
		return nil, err
	}
	dist, err := game.GetSeclusionDurationDist(days)
	if err != nil {
		return nil, err
	}
	return &SeclusionData{
		Stats:        stats,
		DurationDist: dist,
	}, nil
}
