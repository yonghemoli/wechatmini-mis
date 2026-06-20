package rule

import "yonghemolimis/src/domain/valueobj"

// CalcLifecycleStage 计算用户生命周期阶段
// regDays: 注册天数, recentActive: 最近7天活跃天数, prevActive: 前7天活跃天数, wasLost: 是否曾经流失
func CalcLifecycleStage(regDays int, recentActive int, prevActive int, wasLost bool) valueobj.LifecycleStage {
	if wasLost && recentActive > 0 {
		return valueobj.LifecycleReturned
	}
	if recentActive == 0 {
		if prevActive == 0 {
			return valueobj.LifecycleLost
		}
		return valueobj.LifecycleDeclining
	}
	if regDays <= 7 {
		return valueobj.LifecycleNew
	}
	if recentActive >= prevActive {
		return valueobj.LifecycleMature
	}
	if recentActive < prevActive {
		return valueobj.LifecycleDeclining
	}
	return valueobj.LifecycleGrowing
}
