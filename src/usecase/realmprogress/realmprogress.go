package realmprogress

import (
	"yonghemolimis/src/dao/game"
)

// GetRealmProgressAnalysis 境界进阶综合分析
func GetRealmProgressAnalysis(inactiveDays int) ([]game.RealmProgressItem, error) {
	return game.GetRealmProgressAnalysis(inactiveDays)
}

// GetRealmPayCorrelation 境界-付费关联分析
func GetRealmPayCorrelation() ([]game.RealmPayCorrelationItem, error) {
	return game.GetRealmPayCorrelation()
}
