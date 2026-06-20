package equipyonghemolimis

import "yonghemolimis/src/dao/game"

type EquipAnalysisData struct {
	GradeDist     []game.EquipGradeDistItem `json:"grade_dist"`
	SetPopularity []game.EquipSetPopItem    `json:"set_popularity"`
	RealmDist     []game.EquipRealmDistItem `json:"realm_dist"`
}

func GetEquipAnalysis() (*EquipAnalysisData, error) {
	gd, err := game.GetEquipGradeDist()
	if err != nil {
		return nil, err
	}
	sp, err := game.GetEquipSetPopularity(20)
	if err != nil {
		return nil, err
	}
	rd, err := game.GetEquipRealmDist()
	if err != nil {
		return nil, err
	}
	return &EquipAnalysisData{
		GradeDist:     gd,
		SetPopularity: sp,
		RealmDist:     rd,
	}, nil
}
