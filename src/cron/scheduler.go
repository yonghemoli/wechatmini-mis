package cron

import (
	"time"

	"yonghemolimis/src/dao/game"
	"yonghemolimis/src/logger"
	"yonghemolimis/src/usecase/activity"
	"yonghemolimis/src/usecase/feature"
	profileuc "yonghemolimis/src/usecase/profile"
	"yonghemolimis/src/usecase/segment"
	"yonghemolimis/src/usecase/snapshot"
)

// StartScheduler 启动定时任务
func StartScheduler() {
	go func() {
		logger.Info("定时任务调度器已启动")
		for {
			now := time.Now()
			// 计算下一个凌晨 3 点
			next := time.Date(now.Year(), now.Month(), now.Day(), 3, 0, 0, 0, now.Location())
			if now.After(next) {
				next = next.Add(24 * time.Hour)
			}
			duration := next.Sub(now)
			logger.Infof("下次任务将在 %s 后执行 (%s)", duration.Round(time.Minute), next.Format("2006-01-02 15:04"))

			timer := time.NewTimer(duration)
			<-timer.C

			runDailyTasks()
		}
	}()
}

// runDailyTasks 每日定时执行的全部任务
func runDailyTasks() {
	start := time.Now()
	logger.Info("========== 每日定时任务开始 ==========")

	// Step 1: 快照采集 — 必须先于画像计算，因为画像依赖快照数据
	logger.Info("[1/6] 开始快照采集...")
	if err := snapshot.CollectAllSnapshots(3); err != nil {
		logger.Errorf("[1/6] 快照采集失败: %v", err)
	} else {
		logger.Info("[1/6] 快照采集完成")
	}

	// Step 2: 行为日志 ETL
	logger.Info("[2/6] 开始行为日志同步...")
	if err := profileuc.SyncAllActionLogs(1); err != nil {
		logger.Errorf("[2/6] 行为日志同步失败: %v", err)
	} else {
		logger.Info("[2/6] 行为日志同步完成")
	}

	// Step 3: 画像计算
	logger.Info("[3/6] 开始画像计算...")
	if err := profileuc.CalculateAllProfiles(); err != nil {
		logger.Errorf("[3/6] 画像计算失败: %v", err)
	} else {
		logger.Info("[3/6] 画像计算完成")
	}

	// Step 4: 分群刷新
	logger.Info("[4/6] 开始分群人数刷新...")
	if err := segment.RefreshAllSegments(); err != nil {
		logger.Errorf("[4/6] 分群刷新失败: %v", err)
	} else {
		logger.Info("[4/6] 分群刷新完成")
	}

	// Step 5: 功能埋点每日聚合（聚合昨天的原始事件数据）
	yesterday := time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	logger.Infof("[5/7] 开始功能埋点聚合 (%s)...", yesterday)
	if err := feature.AggregateDailyStats(yesterday); err != nil {
		logger.Errorf("[5/7] 功能埋点聚合失败: %v", err)
	} else {
		logger.Info("[5/7] 功能埋点聚合完成")
	}

	// Step 6: 功能质量评分计算
	logger.Info("[6/7] 开始功能质量评分计算...")
	totalUsers, _ := game.GetTotalUserCount()
	if err := feature.CalculateFeatureScores(totalUsers); err != nil {
		logger.Errorf("[6/7] 功能质量评分计算失败: %v", err)
	} else {
		logger.Info("[6/7] 功能质量评分计算完成")
	}

	// Step 7: 活跃度小时聚合
	logger.Infof("[7/7] 开始活跃度小时聚合 (%s)...", yesterday)
	if err := activity.AggregateHourlyStats(yesterday); err != nil {
		logger.Errorf("[7/7] 活跃度小时聚合失败: %v", err)
	} else {
		logger.Info("[7/7] 活跃度小时聚合完成")
	}

	logger.Infof("========== 每日定时任务完成，耗时 %s ==========", time.Since(start).Round(time.Second))
}
