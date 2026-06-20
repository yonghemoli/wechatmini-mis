package economyhealth

import "yonghemolimis/src/dao/game"

var CurrencyLabels = map[string]string{
	"SPIRIT_STONE":   "灵石",
	"SPIRIT_CRYSTAL": "灵晶",
	"SPIRIT_MILK":    "道晶",
	"REPUTATION":     "声望",
	"CHECKIN_COIN":   "签到币",
	"TOWER_SCORE":    "妖塔积分",
	"TONGXIN":        "同心值",
}

type EconomyHealthData struct {
	CurrencyFlow []game.CurrencyFlowItem `json:"currency_flow"`
	WealthDist   []game.WealthDistItem   `json:"wealth_dist"`
}

func GetEconomyHealth(days int) (*EconomyHealthData, error) {
	flow, err := game.GetCurrencyFlow(days)
	if err != nil {
		return nil, err
	}
	wealth, err := game.GetWealthDistribution()
	if err != nil {
		return nil, err
	}
	return &EconomyHealthData{CurrencyFlow: flow, WealthDist: wealth}, nil
}

func GetCurrencyTrend(currency string, days int) ([]game.CurrencyTrendItem, error) {
	return game.GetCurrencyTrend(currency, days)
}
