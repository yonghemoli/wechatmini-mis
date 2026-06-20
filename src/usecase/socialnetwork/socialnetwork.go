package socialnetwork

import "yonghemolimis/src/dao/game"

type SocialNetworkData struct {
	Overview     *game.SocialOverviewStats   `json:"overview"`
	RelationDist []game.RelationTypeDistItem `json:"relation_dist"`
	IntimacyDist []game.IntimacyDistItem     `json:"intimacy_dist"`
}

func GetSocialNetworkData() (*SocialNetworkData, error) {
	overview, err := game.GetSocialOverview()
	if err != nil {
		return nil, err
	}
	relations, err := game.GetRelationTypeDist()
	if err != nil {
		return nil, err
	}
	intimacy, err := game.GetIntimacyDist()
	if err != nil {
		return nil, err
	}
	return &SocialNetworkData{
		Overview:     overview,
		RelationDist: relations,
		IntimacyDist: intimacy,
	}, nil
}
