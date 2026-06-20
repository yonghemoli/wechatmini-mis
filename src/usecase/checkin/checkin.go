package checkin

import "yonghemolimis/src/dao/game"

type CheckInData struct {
	Trend []game.CheckInStatsItem `json:"trend"`
	Rate  *game.CheckInRateItem   `json:"rate"`
}

func GetCheckInAnalysis(days int) (*CheckInData, error) {
	trend, err := game.GetCheckInTrend(days)
	if err != nil {
		return nil, err
	}
	rate, err := game.GetCheckInRate()
	if err != nil {
		return nil, err
	}
	return &CheckInData{
		Trend: trend,
		Rate:  rate,
	}, nil
}
