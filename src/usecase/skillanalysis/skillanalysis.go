package skillyonghemolimis

import "yonghemolimis/src/dao/game"

type SkillAnalysisData struct {
	GradeDist []game.SkillOverviewItem `json:"grade_dist"`
	UserStats *game.UserSkillStats     `json:"user_stats"`
}

func GetSkillAnalysis() (*SkillAnalysisData, error) {
	gradeDist, err := game.GetSkillGradeDistribution()
	if err != nil {
		return nil, err
	}
	userStats, err := game.GetUserSkillStats()
	if err != nil {
		return nil, err
	}
	return &SkillAnalysisData{
		GradeDist: gradeDist,
		UserStats: userStats,
	}, nil
}
