package valueobj

// PlayStyle 玩法偏好
type PlayStyle string

const (
	PlayStyleCombat   PlayStyle = "COMBAT"
	PlayStyleCraft    PlayStyle = "CRAFT"
	PlayStyleSocial   PlayStyle = "SOCIAL"
	PlayStyleEconomy  PlayStyle = "ECONOMY"
	PlayStyleExplorer PlayStyle = "EXPLORER"
	PlayStyleBalanced PlayStyle = "BALANCED"
)

func (p PlayStyle) String() string { return string(p) }

func (p PlayStyle) Label() string {
	switch p {
	case PlayStyleCombat:
		return "战斗型"
	case PlayStyleCraft:
		return "制造型"
	case PlayStyleSocial:
		return "社交型"
	case PlayStyleEconomy:
		return "经济型"
	case PlayStyleExplorer:
		return "探索型"
	case PlayStyleBalanced:
		return "均衡型"
	default:
		return "未知"
	}
}

var AllPlayStyles = []PlayStyle{
	PlayStyleCombat, PlayStyleCraft, PlayStyleSocial,
	PlayStyleEconomy, PlayStyleExplorer, PlayStyleBalanced,
}
