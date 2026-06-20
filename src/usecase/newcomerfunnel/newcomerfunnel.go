package newcomerfunnel

import "yonghemolimis/src/dao/game"

type NewcomerFunnelData struct {
	Funnel []game.NewcomerFunnelItem `json:"funnel"`
}

func GetNewcomerFunnelAnalysis() (*NewcomerFunnelData, error) {
	funnel, err := game.GetNewcomerFunnel()
	if err != nil {
		return nil, err
	}
	return &NewcomerFunnelData{
		Funnel: funnel,
	}, nil
}
