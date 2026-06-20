package realmbottleneck

import "yonghemolimis/src/dao/game"

type RealmBottleneckData struct {
	Bottleneck []game.RealmBottleneckItem `json:"bottleneck"`
	Churn      []game.RealmChurnItem      `json:"churn"`
}

func GetRealmBottleneckAnalysis(inactiveDays int) (*RealmBottleneckData, error) {
	bn, err := game.GetRealmBottleneck(inactiveDays)
	if err != nil {
		return nil, err
	}
	ch, err := game.GetRealmChurn(inactiveDays)
	if err != nil {
		return nil, err
	}
	return &RealmBottleneckData{
		Bottleneck: bn,
		Churn:      ch,
	}, nil
}
