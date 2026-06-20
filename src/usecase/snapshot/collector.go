package snapshot

import (
	"fmt"
	"time"

	"yonghemolimis/src/dao/analytics"
	"yonghemolimis/src/dao/db"
	"yonghemolimis/src/dao/game"
	"yonghemolimis/src/logger"
)

// CollectAllSnapshots 采集所有活跃用户的当日快照
// 策略: 只采集最近 N 天有货币流水(即有游戏行为)的用户，避免全量扫描
func CollectAllSnapshots(activeDays int) error {
	if activeDays <= 0 {
		activeDays = 3
	}
	today := time.Now().Format("2006-01-02")
	logger.Infof("[快照采集] 开始采集 %s 的用户快照 (活跃窗口: %d天)", today, activeDays)

	// 1. 获取活跃用户 ID
	activeUIDs, err := game.GetActiveUserIDs(activeDays)
	if err != nil {
		return fmt.Errorf("获取活跃用户列表失败: %w", err)
	}
	if len(activeUIDs) == 0 {
		logger.Info("[快照采集] 无活跃用户，跳过")
		return nil
	}
	logger.Infof("[快照采集] 找到 %d 个活跃用户", len(activeUIDs))

	// 2. 批量采集
	successCount := 0
	failCount := 0
	batch := make([]db.DailySnapshotDO, 0, 100)

	for _, uid := range activeUIDs {
		snap, err := collectOne(uid, today)
		if err != nil {
			logger.Errorf("[快照采集] 用户 %d 采集失败: %v", uid, err)
			failCount++
			continue
		}
		batch = append(batch, *snap)

		// 每 100 条批量写入
		if len(batch) >= 100 {
			if err := analytics.BatchUpsertSnapshots(batch); err != nil {
				logger.Errorf("[快照采集] 批量写入失败: %v", err)
				failCount += len(batch)
			} else {
				successCount += len(batch)
			}
			batch = batch[:0]
		}
	}

	// 剩余写入
	if len(batch) > 0 {
		if err := analytics.BatchUpsertSnapshots(batch); err != nil {
			logger.Errorf("[快照采集] 批量写入失败: %v", err)
			failCount += len(batch)
		} else {
			successCount += len(batch)
		}
	}

	logger.Infof("[快照采集] 完成: 成功 %d, 失败 %d", successCount, failCount)
	return nil
}

// collectOne 采集单个用户的当日快照
func collectOne(uid uint, date string) (*db.DailySnapshotDO, error) {
	// 1. 获取用户当前状态
	u, err := game.GetUserByID(uid)
	if err != nil {
		return nil, fmt.Errorf("获取用户信息失败: %w", err)
	}

	// 2. 获取当日货币收支汇总
	currAggs, err := game.GetDailyCurrencyAgg(uid, date)
	if err != nil {
		return nil, fmt.Errorf("获取货币汇总失败: %w", err)
	}
	var stoneIncome, stoneExpense, crystalIncome, crystalExpense int64
	for _, agg := range currAggs {
		switch agg.CurrencyCode {
		case "spirit_stone":
			stoneIncome = agg.Income
			stoneExpense = agg.Expense
		case "spirits_crystal":
			crystalIncome = agg.Income
			crystalExpense = agg.Expense
		}
	}

	// 3. 获取当日行为分类计数
	battle, craft, social, economy, err := game.GetDailyActionCounts(uid, date)
	if err != nil {
		// 非致命错误，继续
		logger.Warnf("[快照采集] 用户 %d 行为计数查询失败: %v", uid, err)
	}

	// 4. 补充社交计数：宗门日志
	guildLogCnt, _ := game.GetGuildLogCount(uid, 1)
	social += guildLogCnt

	snap := &db.DailySnapshotDO{
		UID:            uid,
		SnapshotDate:   date,
		RealmID:        u.RealmID,
		VipLevel:       u.VipLevel,
		SpiritStone:    u.SpiritStone,
		SpiritCrystal:  u.SpiritsCrystal,
		SpiritMilk:     u.SpiritsMilk,
		LoginCount:     1, // 有货币流水 = 当日活跃
		BattleCount:    battle,
		CraftCount:     craft,
		SocialCount:    social,
		EconomyCount:   economy,
		StoneIncome:    stoneIncome,
		StoneExpense:   stoneExpense,
		CrystalIncome:  crystalIncome,
		CrystalExpense: crystalExpense,
	}
	return snap, nil
}
