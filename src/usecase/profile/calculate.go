package profile

import (
	"fmt"
	"time"

	"yonghemolimis/src/dao/analytics"
	"yonghemolimis/src/dao/db"
	"yonghemolimis/src/dao/game"
	"yonghemolimis/src/domain/rule"
	"yonghemolimis/src/domain/valueobj"
	"yonghemolimis/src/logger"
)

// CalculateProfile 计算单个用户画像
func CalculateProfile(uid uint) error {
	// 1. 从游戏库拉取用户信息
	gu, err := game.GetUserByID(uid)
	if err != nil {
		return err
	}

	// 2. 计算注册天数
	regDays := int(time.Since(gu.CreatedAt).Hours() / 24)

	// 3. 最近 7 天 / 前 7 天活跃天数
	now := time.Now()
	recent7Start := now.AddDate(0, 0, -7)
	prev7Start := now.AddDate(0, 0, -14)

	recentActive, _ := analytics.CountSnapshotsByDateRange(uid, recent7Start, now)
	prevActive, _ := analytics.CountSnapshotsByDateRange(uid, prev7Start, recent7Start)

	// 4. 判断曾经是否流失
	wasLost := false
	existingProfile, _ := analytics.GetProfile(uid)
	if existingProfile != nil && existingProfile.LifecycleStage == string(valueobj.LifecycleLost) {
		wasLost = true
	}

	// 5. 生命周期阶段
	lifecycle := rule.CalcLifecycleStage(regDays, int(recentActive), int(prevActive), wasLost)

	// 6. 付费分层
	payTier := rule.CalcPayTier(gu.TotalRechargeAmount)

	// 7. 玩法偏好 — 基于最近 7 天快照统计
	snapshots, _ := analytics.GetSnapshotsByUID(uid, 7)
	cc := aggregateCategoryCounts(snapshots)
	playStyle := rule.CalcPlayStyle(cc)

	// 8. 流失风险
	loginDrop := calcLoginFreqDrop(int(recentActive), int(prevActive))
	inactiveDays := calcInactiveDays(snapshots)
	resourceAlert := calcResourceAlert(snapshots)
	stuck := calcStuck(snapshots)
	isPaying := gu.TotalRechargeAmount > 0
	churnRisk := rule.CalcChurnRisk(loginDrop, inactiveDays, resourceAlert, stuck, isPaying)

	// 9. 社交类型
	socialType := calcSocialType(uid)

	// 10. LTV 预测 (简单模型: 基于历史 ARPU × 预估留存月数)
	ltvPredict := calcLTV(gu.TotalRechargeAmount, regDays, int(recentActive))

	// 11. 存储画像
	now2 := time.Now()
	profile := &db.UserProfileDO{
		UID:              uid,
		LifecycleStage:   string(lifecycle),
		PlayStyle:        string(playStyle),
		PayTier:          string(payTier),
		ChurnRisk:        churnRisk,
		StuckFlag:        stuck,
		ResourceAlert:    resourceAlert,
		SocialType:       socialType,
		LtvPredict:       ltvPredict,
		LastCalculatedAt: &now2,
	}
	if err := analytics.UpsertProfile(profile); err != nil {
		return err
	}

	// 12. 写入标签
	writeTags(uid, lifecycle, payTier, playStyle, churnRisk, socialType)

	return nil
}

// CalculateAllProfiles 批量计算所有用户画像
func CalculateAllProfiles() error {
	users, err := game.GetAllUsers()
	if err != nil {
		return err
	}

	logger.Infof("开始计算 %d 个用户画像", len(users))

	successCount := 0
	failCount := 0
	for _, u := range users {
		if err := CalculateProfile(u.ID); err != nil {
			logger.Errorf("计算用户 %d 画像失败: %v", u.ID, err)
			failCount++
		} else {
			successCount++
		}
	}

	logger.Infof("画像计算完成: 成功 %d, 失败 %d", successCount, failCount)
	return nil
}

// ========== 辅助函数 ==========

func aggregateCategoryCounts(snapshots []db.DailySnapshotDO) rule.CategoryCounts {
	var cc rule.CategoryCounts
	for _, s := range snapshots {
		cc.Battle += s.BattleCount
		cc.Craft += s.CraftCount
		cc.Social += s.SocialCount
		cc.Economy += s.EconomyCount
	}
	return cc
}

func calcLoginFreqDrop(recent, prev int) int {
	if prev == 0 {
		return 0
	}
	drop := (prev - recent) * 100 / prev
	if drop < 0 {
		return 0
	}
	return drop
}

func calcInactiveDays(snapshots []db.DailySnapshotDO) int {
	if len(snapshots) == 0 {
		return 7 // 没有快照，视为 7 天不活跃
	}
	// 找最近的快照日期
	latest := snapshots[0].SnapshotDate
	t, err := time.Parse("2006-01-02", latest)
	if err != nil {
		return 0
	}
	days := int(time.Since(t).Hours() / 24)
	return days
}

func calcResourceAlert(snapshots []db.DailySnapshotDO) bool {
	if len(snapshots) == 0 {
		return false
	}
	// 取最近一天
	s := snapshots[0]
	return rule.IsResourceAlert(s.StoneIncome, s.StoneExpense)
}

func calcStuck(snapshots []db.DailySnapshotDO) bool {
	if len(snapshots) < 2 {
		return false
	}
	// 对比 7 天内境界是否变化
	firstRealm := snapshots[len(snapshots)-1].RealmID
	lastRealm := snapshots[0].RealmID
	realmUnchanged := firstRealm == lastRealm

	// 战斗次数是否下降
	battleDrop := false
	if len(snapshots) >= 2 {
		recent := snapshots[0].BattleCount
		older := snapshots[len(snapshots)-1].BattleCount
		battleDrop = recent < older
	}

	unchangedDays := len(snapshots) // 有多少天快照就算多少天
	if realmUnchanged {
		return rule.IsStuck(unchangedDays, battleDrop)
	}
	return false
}

// ========== 社交类型 ==========

// calcSocialType 基于好友 + 宗门信息推算社交类型
func calcSocialType(uid uint) string {
	friendStats, err := game.GetFriendStats(uid)
	if err != nil {
		return "独行侠"
	}

	guildInfo, _ := game.GetGuildMemberInfo(uid)
	hasGuild := guildInfo != nil

	switch {
	case friendStats.Companion > 0:
		return "双修型" // 有道侣
	case hasGuild && friendStats.Intimate >= 5:
		return "社交达人" // 有宗门 + >=5 亲密好友
	case hasGuild && friendStats.Total >= 3:
		return "宗门活跃" // 有宗门 + >=3 好友
	case hasGuild:
		return "宗门成员" // 仅有宗门
	case friendStats.Total >= 3:
		return "交友型" // 无宗门但有好友
	default:
		return "独行侠" // 无宗门无好友
	}
}

// ========== LTV 预测 ==========

// calcLTV 简单 LTV 模型
// 基于: 月均充值 × 预估留存月数
// 留存月数基于当前活跃度和注册时长推算
func calcLTV(totalRecharge int64, regDays int, recentActive int) float64 {
	if totalRecharge <= 0 {
		return 0
	}

	// 月均充值 (分 → 元)
	months := float64(regDays) / 30.0
	if months < 1 {
		months = 1
	}
	monthlyARPU := float64(totalRecharge) / 100.0 / months

	// 预估剩余留存月数: 基于最近活跃度
	// recentActive = 最近7天活跃天数 (0-7)
	var retentionMonths float64
	switch {
	case recentActive >= 6:
		retentionMonths = 12 // 高活跃，预估留存1年
	case recentActive >= 4:
		retentionMonths = 6
	case recentActive >= 2:
		retentionMonths = 3
	case recentActive >= 1:
		retentionMonths = 1
	default:
		retentionMonths = 0.5 // 已不活跃，残余价值
	}

	return monthlyARPU * retentionMonths
}

// ========== 标签写入 ==========

// writeTags 将画像计算结果写入标签表
func writeTags(uid uint, lifecycle valueobj.LifecycleStage, payTier valueobj.PayTier,
	playStyle valueobj.PlayStyle, churnRisk int, socialType string) {

	tags := []db.UserProfileTagDO{
		{UID: uid, TagGroup: "lifecycle", TagKey: "stage", TagValue: string(lifecycle)},
		{UID: uid, TagGroup: "pay", TagKey: "tier", TagValue: string(payTier)},
		{UID: uid, TagGroup: "play", TagKey: "style", TagValue: string(playStyle)},
		{UID: uid, TagGroup: "risk", TagKey: "churn_level", TagValue: riskLevel(churnRisk)},
		{UID: uid, TagGroup: "social", TagKey: "type", TagValue: socialType},
	}

	for _, t := range tags {
		tag := t // copy
		if err := analytics.UpsertTag(&tag); err != nil {
			logger.Errorf("写入标签失败 uid=%d group=%s: %v", uid, t.TagGroup, err)
		}
	}
}

func riskLevel(score int) string {
	switch {
	case score >= 70:
		return "HIGH"
	case score >= 40:
		return "MEDIUM"
	case score >= 20:
		return "LOW"
	default:
		return "NONE"
	}
}

// ========== 行为日志 ETL ==========

// SyncActionLogs 从游戏库同步行为日志到分析库
// 基于货币流水和宗门日志，按行为分类写入 action_logs
func SyncActionLogs(uid uint, days int) error {
	// 1. 同步货币流水作为行为日志
	logs, err := game.GetCurrencyLogs(uid, days)
	if err != nil {
		return fmt.Errorf("获取货币流水失败: %w", err)
	}

	for _, log := range logs {
		category := classifyAction(log.Detail)
		action := summarizeAction(log.Detail, log.CurrencyCode, log.ChangeAmount)
		al := &db.ActionLogDO{
			UID:      log.UID,
			Action:   action,
			Category: category,
			Detail:   log.Detail,
		}
		al.CreatedAt = log.CreatedAt
		_ = analytics.SaveActionLog(al)
	}

	return nil
}

// SyncAllActionLogs 全量同步所有活跃用户的行为日志
func SyncAllActionLogs(days int) error {
	activeUIDs, err := game.GetActiveUserIDs(days)
	if err != nil {
		return err
	}

	logger.Infof("[ETL] 开始同步 %d 个用户的行为日志 (最近 %d 天)", len(activeUIDs), days)

	for _, uid := range activeUIDs {
		if err := SyncActionLogs(uid, days); err != nil {
			logger.Errorf("[ETL] 用户 %d 行为日志同步失败: %v", uid, err)
		}
	}

	logger.Info("[ETL] 行为日志同步完成")
	return nil
}

// classifyAction 根据 detail 文本分类行为
func classifyAction(detail string) string {
	patterns := map[string][]string{
		"BATTLE":   {"战斗", "副本", "击杀", "boss", "妖塔", "秘境", "怪物", "挑战", "PVP", "PVE"},
		"CRAFT":    {"炼器", "炼丹", "制造", "锻造", "分解", "合成", "强化"},
		"SOCIAL":   {"宗门", "好友", "道侣", "赠送", "红包", "师徒"},
		"ECONOMY":  {"交易", "拍卖", "商店", "购买", "出售", "兑换", "集市"},
		"GROWTH":   {"修炼", "突破", "升级", "经验", "境界"},
		"DAILY":    {"签到", "日常", "任务"},
		"RECHARGE": {"充值", "月卡", "首充"},
		"CONSUME":  {"消耗", "使用"},
	}

	for category, keywords := range patterns {
		for _, kw := range keywords {
			if containsStr(detail, kw) {
				return category
			}
		}
	}
	return "OTHER"
}

// summarizeAction 生成简短的行为摘要
func summarizeAction(detail, currencyCode string, amount int64) string {
	if detail != "" && len(detail) <= 32 {
		return detail
	}
	direction := "获得"
	if amount < 0 {
		direction = "消耗"
	}
	return fmt.Sprintf("%s%s %d", direction, currencyCode, abs(amount))
}

func containsStr(s, substr string) bool {
	return len(s) >= len(substr) && findSubstr(s, substr)
}

func findSubstr(s, sub string) bool {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}

func abs(n int64) int64 {
	if n < 0 {
		return -n
	}
	return n
}
