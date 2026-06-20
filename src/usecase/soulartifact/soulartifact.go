package soulartifact

import "yonghemolimis/src/dao/game"

type SoulArtifactData struct {
	UserStats    *game.SoulArtifactUserStats     `json:"user_stats"`
	GradeDist    []game.SoulArtifactGradeDist    `json:"grade_dist"`
	CategoryDist []game.SoulArtifactCategoryDist `json:"category_dist"`
	LevelDist    []game.SoulArtifactLevelDist    `json:"level_dist"`
	Popularity   []game.SoulArtifactPopularity   `json:"popularity"`
}

func GetSoulArtifactAnalysis() (*SoulArtifactData, error) {
	userStats, err := game.GetSoulArtifactUserStats()
	if err != nil {
		return nil, err
	}
	gradeDist, err := game.GetSoulArtifactGradeDist()
	if err != nil {
		return nil, err
	}
	categoryDist, err := game.GetSoulArtifactCategoryDist()
	if err != nil {
		return nil, err
	}
	levelDist, err := game.GetSoulArtifactLevelDist()
	if err != nil {
		return nil, err
	}
	popularity, err := game.GetSoulArtifactPopularity()
	if err != nil {
		return nil, err
	}
	return &SoulArtifactData{
		UserStats:    userStats,
		GradeDist:    gradeDist,
		CategoryDist: categoryDist,
		LevelDist:    levelDist,
		Popularity:   popularity,
	}, nil
}
