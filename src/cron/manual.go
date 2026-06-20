package cron

import (
	"fmt"
	"sync"
	"time"

	"yonghemolimis/src/dao/game"
	"yonghemolimis/src/logger"
	"yonghemolimis/src/usecase/feature"
	profileuc "yonghemolimis/src/usecase/profile"
	"yonghemolimis/src/usecase/segment"
	"yonghemolimis/src/usecase/snapshot"
)

// 防止并发触发
var (
	taskRunning sync.Mutex
	lastRunTime time.Time
)

// TaskResult 单个子任务的执行结果
type TaskResult struct {
	Name    string `json:"name"`
	Success bool   `json:"success"`
	Error   string `json:"error,omitempty"`
	Elapsed string `json:"elapsed"`
}

// RunResult 整体执行结果
type RunResult struct {
	Tasks   []TaskResult `json:"tasks"`
	Total   string       `json:"total"`
	StartAt string       `json:"start_at"`
}

// RunDashboardRefresh 刷新数据总览相关（快照+日志+画像+分群）
func RunDashboardRefresh() (*RunResult, error) {
	if !taskRunning.TryLock() {
		return nil, fmt.Errorf("已有刷新任务正在执行，请稍后再试")
	}
	defer taskRunning.Unlock()

	start := time.Now()
	result := &RunResult{StartAt: start.Format("15:04:05")}

	steps := []struct {
		name string
		fn   func() error
	}{
		{"快照采集", func() error { return snapshot.CollectAllSnapshots(3) }},
		{"行为日志同步", func() error { return profileuc.SyncAllActionLogs(1) }},
		{"画像计算", func() error { return profileuc.CalculateAllProfiles() }},
		{"分群刷新", func() error { return segment.RefreshAllSegments() }},
	}

	for _, s := range steps {
		t := time.Now()
		err := s.fn()
		tr := TaskResult{Name: s.name, Success: err == nil, Elapsed: time.Since(t).Round(time.Millisecond).String()}
		if err != nil {
			tr.Error = err.Error()
			logger.Errorf("[手动刷新] %s 失败: %v", s.name, err)
		} else {
			logger.Infof("[手动刷新] %s 完成 (%s)", s.name, tr.Elapsed)
		}
		result.Tasks = append(result.Tasks, tr)
	}

	result.Total = time.Since(start).Round(time.Millisecond).String()
	lastRunTime = time.Now()
	return result, nil
}

// RunFeatureRefresh 刷新功能分析相关（聚合+评分）
func RunFeatureRefresh() (*RunResult, error) {
	if !taskRunning.TryLock() {
		return nil, fmt.Errorf("已有刷新任务正在执行，请稍后再试")
	}
	defer taskRunning.Unlock()

	start := time.Now()
	result := &RunResult{StartAt: start.Format("15:04:05")}

	// 聚合今天的数据（不只是昨天）
	dates := []string{
		time.Now().AddDate(0, 0, -1).Format("2006-01-02"),
		time.Now().Format("2006-01-02"),
	}

	steps := []struct {
		name string
		fn   func() error
	}{
		{"功能埋点聚合", func() error {
			for _, d := range dates {
				if err := feature.AggregateDailyStats(d); err != nil {
					return fmt.Errorf("日期 %s: %w", d, err)
				}
			}
			return nil
		}},
		{"功能质量评分", func() error {
			total, _ := game.GetTotalUserCount()
			return feature.CalculateFeatureScores(total)
		}},
	}

	for _, s := range steps {
		t := time.Now()
		err := s.fn()
		tr := TaskResult{Name: s.name, Success: err == nil, Elapsed: time.Since(t).Round(time.Millisecond).String()}
		if err != nil {
			tr.Error = err.Error()
			logger.Errorf("[手动刷新] %s 失败: %v", s.name, err)
		} else {
			logger.Infof("[手动刷新] %s 完成 (%s)", s.name, tr.Elapsed)
		}
		result.Tasks = append(result.Tasks, tr)
	}

	result.Total = time.Since(start).Round(time.Millisecond).String()
	lastRunTime = time.Now()
	return result, nil
}

// RunAllTasks 运行全部定时任务（手动全量刷新）
func RunAllTasks() (*RunResult, error) {
	if !taskRunning.TryLock() {
		return nil, fmt.Errorf("已有刷新任务正在执行，请稍后再试")
	}
	defer taskRunning.Unlock()

	start := time.Now()
	result := &RunResult{StartAt: start.Format("15:04:05")}

	dates := []string{
		time.Now().AddDate(0, 0, -1).Format("2006-01-02"),
		time.Now().Format("2006-01-02"),
	}

	steps := []struct {
		name string
		fn   func() error
	}{
		{"快照采集", func() error { return snapshot.CollectAllSnapshots(3) }},
		{"行为日志同步", func() error { return profileuc.SyncAllActionLogs(1) }},
		{"画像计算", func() error { return profileuc.CalculateAllProfiles() }},
		{"分群刷新", func() error { return segment.RefreshAllSegments() }},
		{"功能埋点聚合", func() error {
			for _, d := range dates {
				if err := feature.AggregateDailyStats(d); err != nil {
					return fmt.Errorf("日期 %s: %w", d, err)
				}
			}
			return nil
		}},
		{"功能质量评分", func() error {
			total, _ := game.GetTotalUserCount()
			return feature.CalculateFeatureScores(total)
		}},
	}

	for _, s := range steps {
		t := time.Now()
		err := s.fn()
		tr := TaskResult{Name: s.name, Success: err == nil, Elapsed: time.Since(t).Round(time.Millisecond).String()}
		if err != nil {
			tr.Error = err.Error()
			logger.Errorf("[手动刷新] %s 失败: %v", s.name, err)
		} else {
			logger.Infof("[手动刷新] %s 完成 (%s)", s.name, tr.Elapsed)
		}
		result.Tasks = append(result.Tasks, tr)
	}

	result.Total = time.Since(start).Round(time.Millisecond).String()
	lastRunTime = time.Now()
	return result, nil
}
