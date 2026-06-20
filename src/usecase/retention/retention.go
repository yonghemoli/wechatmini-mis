package retention

import (
	"yonghemolimis/src/dao/game"
)

// GetRetentionData 获取留存分析数据
func GetRetentionData(days int) ([]game.RetentionCohortRow, error) {
	return game.GetRetentionData(days)
}
