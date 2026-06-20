package guildhealth

import "yonghemolimis/src/dao/game"

type GuildHealthData struct {
	Overview *game.GuildOverviewStats `json:"overview"`
	Guilds   []game.GuildHealthItem   `json:"guilds"`
}

func GetGuildHealthData(inactiveDays int) (*GuildHealthData, error) {
	overview, err := game.GetGuildOverview()
	if err != nil {
		return nil, err
	}
	guilds, err := game.GetGuildHealth(inactiveDays)
	if err != nil {
		return nil, err
	}
	return &GuildHealthData{Overview: overview, Guilds: guilds}, nil
}
