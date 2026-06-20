package rule

// CalcChurnRisk 计算流失风险分 (0-100)
// loginFreqDrop: 登录频率下降百分比 (0-100)
// inactiveDays: 连续不活跃天数
// resourceDepleting: 是否资源消耗异常
// stuck: 是否卡关
// isPaying: 是否付费用户
func CalcChurnRisk(loginFreqDrop int, inactiveDays int, resourceDepleting bool, stuck bool, isPaying bool) int {
	score := 0

	// 登录频率下降
	if loginFreqDrop > 50 {
		score += 30
	} else if loginFreqDrop > 20 {
		score += 15
	}

	// 连续不活跃天数
	switch {
	case inactiveDays >= 7:
		score += 35
	case inactiveDays >= 3:
		score += 20
	case inactiveDays >= 1:
		score += 10
	}

	// 资源异常消耗
	if resourceDepleting {
		score += 15
	}

	// 卡关
	if stuck {
		score += 15
	}

	// 付费用户降权（不容易流失）
	if isPaying && score > 20 {
		score = score * 70 / 100
	}

	if score > 100 {
		score = 100
	}
	return score
}

// IsStuck 判断是否卡关：7天内境界未提升且战斗次数下降
func IsStuck(realmUnchangedDays int, battleCountDrop bool) bool {
	return realmUnchangedDays >= 7 && battleCountDrop
}

// IsResourceAlert 判断资源异常：灵石当日消耗 > 收入的 200%
func IsResourceAlert(stoneIncome int64, stoneExpense int64) bool {
	if stoneIncome <= 0 {
		return stoneExpense > 0
	}
	return stoneExpense > stoneIncome*2
}
