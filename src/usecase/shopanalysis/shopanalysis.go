package shopyonghemolimis

import "yonghemolimis/src/dao/game"

type ShopAnalysisData struct {
	Overview       []game.ShopOverview       `json:"overview"`
	PriceRangeDist []game.ShopPriceRangeDist `json:"price_range_dist"`
	CurrencyDist   []game.ShopCurrencyDist   `json:"currency_dist"`
	CategoryDist   []game.ShopCategoryDist   `json:"category_dist"`
}

func GetShopAnalysis() (*ShopAnalysisData, error) {
	overview, err := game.GetShopOverview()
	if err != nil {
		return nil, err
	}
	priceRangeDist, err := game.GetShopPriceRangeDist()
	if err != nil {
		return nil, err
	}
	currencyDist, err := game.GetShopCurrencyDist()
	if err != nil {
		return nil, err
	}
	categoryDist, err := game.GetShopCategoryDist()
	if err != nil {
		return nil, err
	}
	return &ShopAnalysisData{
		Overview:       overview,
		PriceRangeDist: priceRangeDist,
		CurrencyDist:   currencyDist,
		CategoryDist:   categoryDist,
	}, nil
}
