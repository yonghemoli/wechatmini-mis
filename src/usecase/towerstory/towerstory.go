package towerstory

import "yonghemolimis/src/dao/game"

type TowerStoryData struct {
	TowerStats    *game.TowerStatsItem      `json:"tower_stats"`
	TowerLayers   []game.TowerLayerDistItem `json:"tower_layers"`
	StoryProgress []game.StoryProgressItem  `json:"story_progress"`
	StoryFunnel   []game.StoryFunnelItem    `json:"story_funnel"`
}

func GetTowerStoryAnalysis(days int) (*TowerStoryData, error) {
	ts, err := game.GetTowerStats(days)
	if err != nil {
		return nil, err
	}
	tl, err := game.GetTowerLayerDist(days)
	if err != nil {
		return nil, err
	}
	sp, err := game.GetStoryProgress()
	if err != nil {
		return nil, err
	}
	sf, err := game.GetStoryFunnel()
	if err != nil {
		return nil, err
	}
	return &TowerStoryData{
		TowerStats:    ts,
		TowerLayers:   tl,
		StoryProgress: sp,
		StoryFunnel:   sf,
	}, nil
}

func GetStageDifficulty(chapter int) ([]game.StageDifficultyItem, error) {
	return game.GetStageDifficulty(chapter)
}
