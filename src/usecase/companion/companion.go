package companion

import "yonghemolimis/src/dao/game"

type CompanionData struct {
	Stats        *game.CompanionStatsItem     `json:"stats"`
	GiftUsage    []game.CompanionGiftUsage    `json:"gift_usage"`
	IntimacyDist []game.CompanionIntimacyDist `json:"intimacy_dist"`
}

func GetCompanionAnalysis() (*CompanionData, error) {
	stats, err := game.GetCompanionStats()
	if err != nil {
		return nil, err
	}
	giftUsage, err := game.GetCompanionGiftUsage()
	if err != nil {
		return nil, err
	}
	intimacyDist, err := game.GetCompanionIntimacyDist()
	if err != nil {
		return nil, err
	}
	return &CompanionData{
		Stats:        stats,
		GiftUsage:    giftUsage,
		IntimacyDist: intimacyDist,
	}, nil
}
