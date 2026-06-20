package valueobj

// PayTier 付费分层
type PayTier string

const (
	PayTierFree      PayTier = "FREE"
	PayTierMinnow    PayTier = "MINNOW"
	PayTierDolphin   PayTier = "DOLPHIN"
	PayTierOrca      PayTier = "ORCA"
	PayTierWhale     PayTier = "WHALE"
	PayTierLeviathan PayTier = "LEVIATHAN"
)

func (p PayTier) String() string { return string(p) }

func (p PayTier) Label() string {
	switch p {
	case PayTierFree:
		return "免费玩家"
	case PayTierMinnow:
		return "微氪"
	case PayTierDolphin:
		return "小氪"
	case PayTierOrca:
		return "中氪"
	case PayTierWhale:
		return "大氪"
	case PayTierLeviathan:
		return "巨鲸"
	default:
		return "未知"
	}
}

var AllPayTiers = []PayTier{
	PayTierFree, PayTierMinnow, PayTierDolphin,
	PayTierOrca, PayTierWhale, PayTierLeviathan,
}
