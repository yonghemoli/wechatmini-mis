package segment

import (
	"encoding/json"
	"fmt"

	"yonghemolimis/src/dao/db"
	"yonghemolimis/src/logger"

	"gorm.io/gorm"
)

// SegmentRule 单条分群规则
type SegmentRule struct {
	Field    string `json:"field"`    // 字段: lifecycle_stage, pay_tier, play_style, churn_risk, social_type, stuck_flag, resource_alert, ltv_predict
	Operator string `json:"operator"` // 操作符: eq, neq, gt, gte, lt, lte, in, contains
	Value    any    `json:"value"`    // 值: 字符串 / 数字 / 数组
}

// SegmentRuleSet 分群规则集
type SegmentRuleSet struct {
	Logic string        `json:"logic"` // "AND" | "OR", 默认 AND
	Rules []SegmentRule `json:"rules"`
}

// allowedFields 分群可查询的合法字段白名单
var allowedFields = map[string]bool{
	"lifecycle_stage": true,
	"pay_tier":        true,
	"play_style":      true,
	"churn_risk":      true,
	"social_type":     true,
	"stuck_flag":      true,
	"resource_alert":  true,
	"ltv_predict":     true,
}

// ParseRules 解析 JSON 规则
func ParseRules(rulesJSON string) (*SegmentRuleSet, error) {
	var ruleSet SegmentRuleSet
	if err := json.Unmarshal([]byte(rulesJSON), &ruleSet); err != nil {
		return nil, fmt.Errorf("规则JSON解析失败: %w", err)
	}
	if ruleSet.Logic == "" {
		ruleSet.Logic = "AND"
	}
	if len(ruleSet.Rules) == 0 {
		return nil, fmt.Errorf("规则不能为空")
	}
	// 校验字段白名单
	for _, r := range ruleSet.Rules {
		if !allowedFields[r.Field] {
			return nil, fmt.Errorf("不支持的筛选字段: %s", r.Field)
		}
	}
	return &ruleSet, nil
}

// ExecuteSegment 执行分群查询，返回匹配的用户画像列表和总数
func ExecuteSegment(seg *db.SegmentDO) ([]db.UserProfileDO, int, error) {
	ruleSet, err := ParseRules(seg.RulesJSON)
	if err != nil {
		return nil, 0, err
	}

	query := db.Get().Model(&db.UserProfileDO{})
	query = applyRules(query, ruleSet)

	var total int64
	query.Count(&total)

	var profiles []db.UserProfileDO
	err = query.Order("churn_risk DESC").Limit(1000).Find(&profiles).Error
	if err != nil {
		return nil, 0, fmt.Errorf("分群查询失败: %w", err)
	}
	return profiles, int(total), nil
}

// CountSegment 仅计算分群人数（不返回明细）
func CountSegment(rulesJSON string) (int, error) {
	ruleSet, err := ParseRules(rulesJSON)
	if err != nil {
		return 0, err
	}

	query := db.Get().Model(&db.UserProfileDO{})
	query = applyRules(query, ruleSet)

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return 0, err
	}
	return int(total), nil
}

// RefreshAllSegments 刷新所有分群的人数
func RefreshAllSegments() error {
	segs, err := db.Get().Find(&[]db.SegmentDO{}).Rows()
	if err != nil {
		return err
	}
	defer segs.Close()

	var allSegs []db.SegmentDO
	for segs.Next() {
		var s db.SegmentDO
		if err := db.Get().ScanRows(segs, &s); err != nil {
			continue
		}
		allSegs = append(allSegs, s)
	}

	for _, s := range allSegs {
		count, err := CountSegment(s.RulesJSON)
		if err != nil {
			logger.Errorf("[分群] 刷新分群 %d(%s) 人数失败: %v", s.ID, s.Name, err)
			continue
		}
		db.Get().Model(&db.SegmentDO{}).Where("id = ?", s.ID).Update("user_count", count)
		logger.Infof("[分群] 分群 %d(%s) 人数: %d", s.ID, s.Name, count)
	}
	return nil
}

// applyRules 将规则集转换为 GORM WHERE 条件
func applyRules(query *gorm.DB, ruleSet *SegmentRuleSet) *gorm.DB {
	if ruleSet.Logic == "OR" {
		return applyOR(query, ruleSet.Rules)
	}
	// AND 逻辑（默认）
	for _, r := range ruleSet.Rules {
		query = applySingleRule(query, r)
	}
	return query
}

func applyOR(query *gorm.DB, rules []SegmentRule) *gorm.DB {
	if len(rules) == 0 {
		return query
	}

	// 构建 OR 条件组
	combined := db.Get().Where("1 = 0") // base false
	for _, r := range rules {
		sub := db.Get().Model(&db.UserProfileDO{})
		sub = applySingleRule(sub, r)
		combined = combined.Or(sub)
	}
	return query.Where(combined)
}

func applySingleRule(query *gorm.DB, r SegmentRule) *gorm.DB {
	col := r.Field
	val := r.Value

	switch r.Operator {
	case "eq", "=", "==":
		return query.Where(col+" = ?", val)
	case "neq", "!=", "<>":
		return query.Where(col+" != ?", val)
	case "gt", ">":
		return query.Where(col+" > ?", val)
	case "gte", ">=":
		return query.Where(col+" >= ?", val)
	case "lt", "<":
		return query.Where(col+" < ?", val)
	case "lte", "<=":
		return query.Where(col+" <= ?", val)
	case "in":
		return query.Where(col+" IN (?)", val)
	case "contains":
		return query.Where(col+" LIKE ?", fmt.Sprintf("%%%v%%", val))
	default:
		logger.Warnf("[分群] 未知操作符: %s，跳过", r.Operator)
		return query
	}
}
