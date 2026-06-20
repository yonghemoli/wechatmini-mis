package funnel

import (
	"yonghemolimis/src/dao/game"
)

// GetFunnelData 获取转化漏斗数据
func GetFunnelData(days int) (*game.FunnelData, error) {
	return game.GetFunnelData(days)
}
