package activity

import (
	"fmt"
	"math"
	"time"
	"yonghemolimis/src/dao/db"
	"yonghemolimis/src/logger"
)

// ==================== 聚合 ====================

// AggregateHourlyStats 聚合指定日期的原始 feature_events → 用户小时统计 + 全局小时统计
func AggregateHourlyStats(date string) error {
	d := db.Get()

	// ========== 1. 用户维度：按 uid + hour 聚合 ==========
	type userHourRow struct {
		UID          uint   `gorm:"column:uid"`
		Hour         int    `gorm:"column:hour"`
		CommandCount int    `gorm:"column:command_count"`
		FirstSeen    string `gorm:"column:first_seen"`
		LastSeen     string `gorm:"column:last_seen"`
	}

	var userRows []userHourRow
	err := d.Model(&db.FeatureEventDO{}).
		Select(`uid,
			HOUR(event_time) as hour,
			COUNT(*) as command_count,
			DATE_FORMAT(MIN(event_time), '%H:%i:%s') as first_seen,
			DATE_FORMAT(MAX(event_time), '%H:%i:%s') as last_seen`).
		Where("DATE(event_time) = ?", date).
		Group("uid, HOUR(event_time)").
		Find(&userRows).Error
	if err != nil {
		return fmt.Errorf("聚合用户小时数据失败: %w", err)
	}

	for _, r := range userRows {
		stat := db.UserHourlyStatsDO{
			UID:          r.UID,
			Date:         date,
			Hour:         r.Hour,
			CommandCount: r.CommandCount,
			FirstSeenAt:  r.FirstSeen,
			LastSeenAt:   r.LastSeen,
		}
		result := d.Where("uid = ? AND date = ? AND hour = ?", r.UID, date, r.Hour).
			Assign(stat).FirstOrCreate(&stat)
		if result.Error != nil {
			logger.Errorf("[activity] upsert user hourly stat uid=%d date=%s hour=%d: %v", r.UID, date, r.Hour, result.Error)
		}
	}

	// ========== 2. 全局维度：按 hour 聚合 ==========
	type globalRow struct {
		Hour          int `gorm:"column:hour"`
		ActiveUsers   int `gorm:"column:active_users"`
		TotalCommands int `gorm:"column:total_commands"`
	}

	var globalRows []globalRow
	err = d.Model(&db.FeatureEventDO{}).
		Select(`HOUR(event_time) as hour,
			COUNT(DISTINCT uid) as active_users,
			COUNT(*) as total_commands`).
		Where("DATE(event_time) = ?", date).
		Group("HOUR(event_time)").
		Find(&globalRows).Error
	if err != nil {
		return fmt.Errorf("聚合全局小时数据失败: %w", err)
	}

	for _, r := range globalRows {
		stat := db.GlobalHourlyStatsDO{
			Date:          date,
			Hour:          r.Hour,
			ActiveUsers:   r.ActiveUsers,
			TotalCommands: r.TotalCommands,
		}
		result := d.Where("date = ? AND hour = ?", date, r.Hour).
			Assign(stat).FirstOrCreate(&stat)
		if result.Error != nil {
			logger.Errorf("[activity] upsert global hourly stat date=%s hour=%d: %v", date, r.Hour, result.Error)
		}
	}

	logger.Infof("[activity] 聚合完成 date=%s: %d 条用户小时记录, %d 条全局小时记录", date, len(userRows), len(globalRows))
	return nil
}

// AggregateMultipleDays 聚合最近 N 天的小时数据
func AggregateMultipleDays(days int) error {
	now := time.Now()
	for i := 1; i <= days; i++ {
		date := now.AddDate(0, 0, -i).Format("2006-01-02")
		if err := AggregateHourlyStats(date); err != nil {
			logger.Errorf("[activity] 聚合 %s 失败: %v", date, err)
		}
	}
	return nil
}

// ==================== 查询接口 ====================

// PlayerHourlyStat 玩家小时活跃数据
type PlayerHourlyStat struct {
	Date         string `json:"date"`
	Hour         int    `json:"hour"`
	CommandCount int    `json:"command_count"`
	FirstSeenAt  string `json:"first_seen_at"`
	LastSeenAt   string `json:"last_seen_at"`
}

// GetPlayerHourlyStats 获取某玩家的小时活跃数据
func GetPlayerHourlyStats(uid uint, days int) ([]PlayerHourlyStat, error) {
	since := time.Now().AddDate(0, 0, -days).Format("2006-01-02")
	var rows []db.UserHourlyStatsDO
	err := db.Get().Where("uid = ? AND date >= ?", uid, since).
		Order("date ASC, hour ASC").Find(&rows).Error
	if err != nil {
		return nil, err
	}

	result := make([]PlayerHourlyStat, len(rows))
	for i, r := range rows {
		result[i] = PlayerHourlyStat{
			Date:         r.Date,
			Hour:         r.Hour,
			CommandCount: r.CommandCount,
			FirstSeenAt:  r.FirstSeenAt,
			LastSeenAt:   r.LastSeenAt,
		}
	}
	return result, nil
}

// PlayerDailyStat 玩家每日汇总
type PlayerDailyStat struct {
	Date         string  `json:"date"`
	TotalCmds    int     `json:"total_cmds"`
	ActiveHours  int     `json:"active_hours"`
	OnlineMinute float64 `json:"online_minute"`
	FirstSeen    string  `json:"first_seen"`
	LastSeen     string  `json:"last_seen"`
}

// GetPlayerDailyStats 获取某玩家的每日汇总
func GetPlayerDailyStats(uid uint, days int) ([]PlayerDailyStat, error) {
	since := time.Now().AddDate(0, 0, -days).Format("2006-01-02")
	var rows []db.UserHourlyStatsDO
	err := db.Get().Where("uid = ? AND date >= ?", uid, since).
		Order("date ASC, hour ASC").Find(&rows).Error
	if err != nil {
		return nil, err
	}

	// 按天分组
	dayMap := make(map[string][]db.UserHourlyStatsDO)
	var dateOrder []string
	for _, r := range rows {
		if _, exists := dayMap[r.Date]; !exists {
			dateOrder = append(dateOrder, r.Date)
		}
		dayMap[r.Date] = append(dayMap[r.Date], r)
	}

	result := make([]PlayerDailyStat, 0, len(dayMap))
	for _, date := range dateOrder {
		hrs := dayMap[date]
		totalCmds := 0
		firstSeen := "23:59:59"
		lastSeen := "00:00:00"
		for _, h := range hrs {
			totalCmds += h.CommandCount
			if h.FirstSeenAt < firstSeen {
				firstSeen = h.FirstSeenAt
			}
			if h.LastSeenAt > lastSeen {
				lastSeen = h.LastSeenAt
			}
		}

		// 估算在线时长: 活跃小时数 × 平均每小时在线比例
		onlineMin := float64(len(hrs)) * 60.0 * 0.5 // 粗估：每个活跃小时算30分钟
		if len(hrs) == 1 && totalCmds <= 2 {
			onlineMin = 5 // 只有1次活跃且指令极少，算5分钟
		}

		result = append(result, PlayerDailyStat{
			Date:         date,
			TotalCmds:    totalCmds,
			ActiveHours:  len(hrs),
			OnlineMinute: onlineMin,
			FirstSeen:    firstSeen,
			LastSeen:     lastSeen,
		})
	}
	return result, nil
}

// GlobalHourlyStat 全局小时活跃
type GlobalHourlyStat struct {
	Date          string `json:"date"`
	Hour          int    `json:"hour"`
	ActiveUsers   int    `json:"active_users"`
	TotalCommands int    `json:"total_commands"`
}

// GetGlobalHourlyStats 获取全局小时活跃数据
func GetGlobalHourlyStats(days int) ([]GlobalHourlyStat, error) {
	since := time.Now().AddDate(0, 0, -days).Format("2006-01-02")
	var rows []db.GlobalHourlyStatsDO
	err := db.Get().Where("date >= ?", since).
		Order("date ASC, hour ASC").Find(&rows).Error
	if err != nil {
		return nil, err
	}

	result := make([]GlobalHourlyStat, len(rows))
	for i, r := range rows {
		result[i] = GlobalHourlyStat{
			Date:          r.Date,
			Hour:          r.Hour,
			ActiveUsers:   r.ActiveUsers,
			TotalCommands: r.TotalCommands,
		}
	}
	return result, nil
}

// ==================== 高峰分析 ====================

// HourAvg 每小时平均
type HourAvg struct {
	Hour        int     `json:"hour"`
	AvgUsers    float64 `json:"avg_users"`
	AvgCommands float64 `json:"avg_commands"`
}

// WeekdayAvg 每星期几平均
type WeekdayAvg struct {
	Weekday     int     `json:"weekday"` // 0=Sunday
	WeekdayName string  `json:"weekday_name"`
	AvgUsers    float64 `json:"avg_users"`
	AvgCommands float64 `json:"avg_commands"`
}

// PeakAnalysis 高峰期分析结果
type PeakAnalysis struct {
	HourlyAvg   [24]HourAvg   `json:"hourly_avg"`
	WeekdayAvg  [7]WeekdayAvg `json:"weekday_avg"`
	PeakHour    int           `json:"peak_hour"`
	PeakWeekday int           `json:"peak_weekday"` // 0=Sunday
	PeakUsers   float64       `json:"peak_users"`
	TotalDays   int           `json:"total_days"`
}

var weekdayNames = [7]string{"周日", "周一", "周二", "周三", "周四", "周五", "周六"}

// GetPeakAnalysis 分析高峰时段（按小时 + 按星期几）
func GetPeakAnalysis(days int) (*PeakAnalysis, error) {
	since := time.Now().AddDate(0, 0, -days).Format("2006-01-02")
	var rows []db.GlobalHourlyStatsDO
	err := db.Get().Where("date >= ?", since).Find(&rows).Error
	if err != nil {
		return nil, err
	}
	if len(rows) == 0 {
		return &PeakAnalysis{}, nil
	}

	// 按小时聚合
	type hourAcc struct {
		totalUsers int
		totalCmds  int
		count      int
	}
	hourMap := make(map[int]*hourAcc)
	for i := 0; i < 24; i++ {
		hourMap[i] = &hourAcc{}
	}

	// 按星期几聚合
	type wdAcc struct {
		totalUsers int
		totalCmds  int
		daySet     map[string]bool
	}
	wdMap := make(map[int]*wdAcc)
	for i := 0; i < 7; i++ {
		wdMap[i] = &wdAcc{daySet: make(map[string]bool)}
	}

	dateSet := make(map[string]bool)
	for _, r := range rows {
		dateSet[r.Date] = true

		// 小时维度
		hourMap[r.Hour].totalUsers += r.ActiveUsers
		hourMap[r.Hour].totalCmds += r.TotalCommands
		hourMap[r.Hour].count++

		// 星期维度
		t, err := time.Parse("2006-01-02", r.Date)
		if err != nil {
			continue
		}
		wd := int(t.Weekday())
		wdMap[wd].totalUsers += r.ActiveUsers
		wdMap[wd].totalCmds += r.TotalCommands
		wdMap[wd].daySet[r.Date] = true
	}

	result := &PeakAnalysis{TotalDays: len(dateSet)}
	peakUsers := 0.0
	for i := 0; i < 24; i++ {
		h := hourMap[i]
		if h.count > 0 {
			avg := float64(h.totalUsers) / float64(h.count)
			result.HourlyAvg[i] = HourAvg{
				Hour:        i,
				AvgUsers:    math.Round(avg*100) / 100,
				AvgCommands: math.Round(float64(h.totalCmds)/float64(h.count)*100) / 100,
			}
			if avg > peakUsers {
				peakUsers = avg
				result.PeakHour = i
				result.PeakUsers = math.Round(avg*100) / 100
			}
		} else {
			result.HourlyAvg[i] = HourAvg{Hour: i}
		}
	}

	peakWdUsers := 0.0
	for i := 0; i < 7; i++ {
		w := wdMap[i]
		days := len(w.daySet)
		if days == 0 {
			days = 1
		}
		avg := float64(w.totalUsers) / float64(days)
		result.WeekdayAvg[i] = WeekdayAvg{
			Weekday:     i,
			WeekdayName: weekdayNames[i],
			AvgUsers:    math.Round(avg*100) / 100,
			AvgCommands: math.Round(float64(w.totalCmds)/float64(days)*100) / 100,
		}
		if avg > peakWdUsers {
			peakWdUsers = avg
			result.PeakWeekday = i
		}
	}

	return result, nil
}

// ==================== 玩家高峰 ====================

// PlayerPeakAnalysis 玩家个人高峰分析
type PlayerPeakAnalysis struct {
	HourlyAvg [24]HourAvg `json:"hourly_avg"`
	PeakHour  int         `json:"peak_hour"`
	TotalDays int         `json:"total_days"`
	TotalCmds int         `json:"total_cmds"`
}

// GetPlayerPeakAnalysis 分析某玩家的高峰时段
func GetPlayerPeakAnalysis(uid uint, days int) (*PlayerPeakAnalysis, error) {
	since := time.Now().AddDate(0, 0, -days).Format("2006-01-02")
	var rows []db.UserHourlyStatsDO
	err := db.Get().Where("uid = ? AND date >= ?", uid, since).Find(&rows).Error
	if err != nil {
		return nil, err
	}

	type hourAcc struct {
		totalCmds int
		count     int
	}
	hourMap := make(map[int]*hourAcc)
	for i := 0; i < 24; i++ {
		hourMap[i] = &hourAcc{}
	}

	dateSet := make(map[string]bool)
	totalCmds := 0
	for _, r := range rows {
		dateSet[r.Date] = true
		hourMap[r.Hour].totalCmds += r.CommandCount
		hourMap[r.Hour].count++
		totalCmds += r.CommandCount
	}

	result := &PlayerPeakAnalysis{TotalDays: len(dateSet), TotalCmds: totalCmds}
	peakCmds := 0.0
	for i := 0; i < 24; i++ {
		h := hourMap[i]
		if h.count > 0 {
			avg := float64(h.totalCmds) / float64(h.count)
			result.HourlyAvg[i] = HourAvg{
				Hour:        i,
				AvgUsers:    1, // 单用户
				AvgCommands: math.Round(avg*100) / 100,
			}
			if avg > peakCmds {
				peakCmds = avg
				result.PeakHour = i
			}
		} else {
			result.HourlyAvg[i] = HourAvg{Hour: i}
		}
	}

	return result, nil
}

// ==================== 真人检测 ====================

// BotCheckResult 机器人检测结果
type BotCheckResult struct {
	UID         uint     `json:"uid"`
	IsLikelyBot bool     `json:"is_likely_bot"`
	BotScore    int      `json:"bot_score"` // 0-100, 越高越像机器人
	Reasons     []string `json:"reasons"`
	// 参考数据
	AvgDailyCmds   float64 `json:"avg_daily_cmds"`
	ActiveDays     int     `json:"active_days"`
	AvgActiveHours float64 `json:"avg_active_hours"`
	CmdStdDev      float64 `json:"cmd_std_dev"`
	NightActivity  float64 `json:"night_activity"` // 凌晨活跃占比 (0-5点)
	Regularity     float64 `json:"regularity"`     // 规律性评分
}

// CheckBot 检测某玩家是否为脚本/机器人
func CheckBot(uid uint, days int) (*BotCheckResult, error) {
	since := time.Now().AddDate(0, 0, -days).Format("2006-01-02")
	var rows []db.UserHourlyStatsDO
	err := db.Get().Where("uid = ? AND date >= ?", uid, since).
		Order("date ASC, hour ASC").Find(&rows).Error
	if err != nil {
		return nil, err
	}

	result := &BotCheckResult{UID: uid}
	if len(rows) == 0 {
		return result, nil
	}

	// 按天分组
	dayMap := make(map[string][]db.UserHourlyStatsDO)
	for _, r := range rows {
		dayMap[r.Date] = append(dayMap[r.Date], r)
	}

	activeDays := len(dayMap)
	result.ActiveDays = activeDays

	// 1. 每日总指令数
	var dailyCmds []float64
	var dailyActiveHours []float64
	totalCmds := 0
	nightCmds := 0
	totalAllCmds := 0
	for _, hrs := range dayMap {
		dayCnt := 0
		for _, h := range hrs {
			dayCnt += h.CommandCount
			totalAllCmds += h.CommandCount
			if h.Hour >= 0 && h.Hour <= 5 {
				nightCmds += h.CommandCount
			}
		}
		dailyCmds = append(dailyCmds, float64(dayCnt))
		dailyActiveHours = append(dailyActiveHours, float64(len(hrs)))
		totalCmds += dayCnt
	}

	if activeDays > 0 {
		result.AvgDailyCmds = math.Round(float64(totalCmds)/float64(activeDays)*100) / 100
		avgHrs := 0.0
		for _, h := range dailyActiveHours {
			avgHrs += h
		}
		result.AvgActiveHours = math.Round(avgHrs/float64(activeDays)*100) / 100
	}

	// 2. 凌晨活跃占比
	if totalAllCmds > 0 {
		result.NightActivity = math.Round(float64(nightCmds)/float64(totalAllCmds)*10000) / 100
	}

	// 3. 每小时指令数标准差（跨所有记录）— 过于一致说明是脚本
	var allHourlyCmds []float64
	for _, r := range rows {
		allHourlyCmds = append(allHourlyCmds, float64(r.CommandCount))
	}
	result.CmdStdDev = math.Round(stdDev(allHourlyCmds)*100) / 100

	// 4. 规律性检测：每天活跃的小时分布是否高度重复
	hourFreq := make(map[int]int) // 每个小时出现了多少天
	for _, hrs := range dayMap {
		for _, h := range hrs {
			hourFreq[h.Hour]++
		}
	}
	// 如果某些小时几乎每天都出现，规律性高
	highFreqHours := 0
	for _, freq := range hourFreq {
		if activeDays > 2 && float64(freq)/float64(activeDays) >= 0.8 {
			highFreqHours++
		}
	}
	if activeDays > 0 {
		result.Regularity = math.Round(float64(highFreqHours)/float64(result.AvgActiveHours)*10000) / 100
	}

	// ========== 综合评分 ==========
	score := 0

	// 规则1: 凌晨活跃占比 > 30%
	if result.NightActivity > 30 {
		score += 25
		result.Reasons = append(result.Reasons, fmt.Sprintf("凌晨(0-5时)活跃占比%.1f%%，显著偏高", result.NightActivity))
	} else if result.NightActivity > 15 {
		score += 10
		result.Reasons = append(result.Reasons, fmt.Sprintf("凌晨活跃占比%.1f%%，略偏高", result.NightActivity))
	}

	// 规则2: 日均活跃超过18小时
	if result.AvgActiveHours >= 18 {
		score += 30
		result.Reasons = append(result.Reasons, fmt.Sprintf("日均活跃%.1f小时，接近全天在线", result.AvgActiveHours))
	} else if result.AvgActiveHours >= 14 {
		score += 15
		result.Reasons = append(result.Reasons, fmt.Sprintf("日均活跃%.1f小时，在线时间过长", result.AvgActiveHours))
	}

	// 规则3: 每小时指令数标准差极低（高度均匀 = 脚本特征）
	if len(allHourlyCmds) > 5 && result.CmdStdDev < 2 && result.AvgDailyCmds > 50 {
		score += 25
		result.Reasons = append(result.Reasons, fmt.Sprintf("每小时指令数标准差仅%.2f，操作过于均匀", result.CmdStdDev))
	} else if len(allHourlyCmds) > 5 && result.CmdStdDev < 5 && result.AvgDailyCmds > 100 {
		score += 10
		result.Reasons = append(result.Reasons, fmt.Sprintf("指令频率较为规律(σ=%.2f)", result.CmdStdDev))
	}

	// 规则4: 时间分布规律性极高
	if result.Regularity > 80 && activeDays > 3 {
		score += 20
		result.Reasons = append(result.Reasons, fmt.Sprintf("活跃时段规律性%.0f%%，每天几乎相同时间段操作", result.Regularity))
	}

	if score > 100 {
		score = 100
	}
	result.BotScore = score
	result.IsLikelyBot = score >= 60

	if len(result.Reasons) == 0 {
		result.Reasons = append(result.Reasons, "行为模式正常，是真实玩家")
	}

	return result, nil
}

// stdDev 计算标准差
func stdDev(data []float64) float64 {
	if len(data) < 2 {
		return 0
	}
	mean := 0.0
	for _, v := range data {
		mean += v
	}
	mean /= float64(len(data))
	sum := 0.0
	for _, v := range data {
		sum += (v - mean) * (v - mean)
	}
	return math.Sqrt(sum / float64(len(data)))
}
