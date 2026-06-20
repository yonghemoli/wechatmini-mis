package rule

import "yonghemolimis/src/domain/valueobj"

// CalcPayTier 根据总充值金额（单位：分）计算付费分层
func CalcPayTier(totalRechargeAmount int64) valueobj.PayTier {
	yuan := totalRechargeAmount / 100
	switch {
	case yuan <= 0:
		return valueobj.PayTierFree
	case yuan <= 30:
		return valueobj.PayTierMinnow
	case yuan <= 328:
		return valueobj.PayTierDolphin
	case yuan <= 1280:
		return valueobj.PayTierOrca
	case yuan <= 6480:
		return valueobj.PayTierWhale
	default:
		return valueobj.PayTierLeviathan
	}
}

// CalcVipLevel 根据总充值金额推算 VIP 等级 (0-15)
func CalcVipLevel(totalRechargeAmount int64) int {
	yuan := totalRechargeAmount / 100
	thresholds := []int64{0, 6, 30, 98, 198, 328, 648, 1280, 1980, 3280, 4880, 6480, 9800, 16800, 29800, 50000}
	level := 0
	for i, t := range thresholds {
		if yuan >= t {
			level = i
		} else {
			break
		}
	}
	return level
}
