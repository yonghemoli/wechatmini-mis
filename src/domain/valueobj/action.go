package valueobj

// ActionCategory 行为大类
type ActionCategory string

const (
	ActionGrowth   ActionCategory = "GROWTH"
	ActionBattle   ActionCategory = "BATTLE"
	ActionCraft    ActionCategory = "CRAFT"
	ActionEconomy  ActionCategory = "ECONOMY"
	ActionSocial   ActionCategory = "SOCIAL"
	ActionDaily    ActionCategory = "DAILY"
	ActionConsume  ActionCategory = "CONSUME"
	ActionRecharge ActionCategory = "RECHARGE"
)

func (a ActionCategory) String() string { return string(a) }

// RiskLevel 风险等级
type RiskLevel string

const (
	RiskNone   RiskLevel = "NONE"
	RiskLow    RiskLevel = "LOW"
	RiskMedium RiskLevel = "MEDIUM"
	RiskHigh   RiskLevel = "HIGH"
)

func (r RiskLevel) String() string { return string(r) }
