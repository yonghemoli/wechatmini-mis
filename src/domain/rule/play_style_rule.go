package rule

import "yonghemolimis/src/domain/valueobj"

// CategoryCounts 各行为大类计数
type CategoryCounts struct {
	Battle  int
	Craft   int
	Social  int
	Economy int
	Explore int
}

// CalcPlayStyle 根据行为计数计算玩法偏好
func CalcPlayStyle(c CategoryCounts) valueobj.PlayStyle {
	total := c.Battle + c.Craft + c.Social + c.Economy + c.Explore
	if total == 0 {
		return valueobj.PlayStyleBalanced
	}

	type kv struct {
		style valueobj.PlayStyle
		count int
	}
	items := []kv{
		{valueobj.PlayStyleCombat, c.Battle},
		{valueobj.PlayStyleCraft, c.Craft},
		{valueobj.PlayStyleSocial, c.Social},
		{valueobj.PlayStyleEconomy, c.Economy},
		{valueobj.PlayStyleExplorer, c.Explore},
	}

	maxItem := items[0]
	for _, it := range items[1:] {
		if it.count > maxItem.count {
			maxItem = it
		}
	}

	// 如果最高占比 < 40%，视为均衡型
	if float64(maxItem.count)/float64(total) < 0.40 {
		return valueobj.PlayStyleBalanced
	}
	return maxItem.style
}
