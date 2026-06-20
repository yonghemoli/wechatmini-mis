package game

import (
	"sort"
	"time"
	"yonghemolimis/src/dao/db"
)

const (
	realRechargeOrderWhere      = "status = 'SUCCESS' AND COALESCE(order_no, '') NOT LIKE 'MIS%'"
	realRechargeOrderWhereAlias = "ro.status = 'SUCCESS' AND COALESCE(ro.order_no, '') NOT LIKE 'MIS%'"
)

type RechargeOrderFilter string

const (
	RechargeOrderFilterAll      RechargeOrderFilter = "all"
	RechargeOrderFilterMIS      RechargeOrderFilter = "mis"
	RechargeOrderFilterOfficial RechargeOrderFilter = "official"
)

func normalizeRechargeOrderFilter(filter string) RechargeOrderFilter {
	switch RechargeOrderFilter(filter) {
	case RechargeOrderFilterMIS:
		return RechargeOrderFilterMIS
	case RechargeOrderFilterOfficial:
		return RechargeOrderFilterOfficial
	default:
		return RechargeOrderFilterAll
	}
}

func rechargeOrderWhere(filter string) string {
	switch normalizeRechargeOrderFilter(filter) {
	case RechargeOrderFilterMIS:
		return "status = 'SUCCESS' AND COALESCE(order_no, '') LIKE 'MIS%'"
	case RechargeOrderFilterOfficial:
		return realRechargeOrderWhere
	default:
		return "status = 'SUCCESS'"
	}
}

func rechargeOrderWhereAlias(filter string) string {
	switch normalizeRechargeOrderFilter(filter) {
	case RechargeOrderFilterMIS:
		return "ro.status = 'SUCCESS' AND COALESCE(ro.order_no, '') LIKE 'MIS%'"
	case RechargeOrderFilterOfficial:
		return realRechargeOrderWhereAlias
	default:
		return "ro.status = 'SUCCESS'"
	}
}

// ========== 游戏库只读查询 ==========

// GameUserRow 游戏用户行
type GameUserRow struct {
	ID                      uint       `gorm:"column:id" json:"id"`
	Name                    string     `gorm:"column:name" json:"name"`
	Gender                  int        `gorm:"column:gender" json:"gender"`
	Status                  int        `gorm:"column:status" json:"status"` // 状态: 1正常
	RealmID                 int        `gorm:"column:realm_id" json:"realm_id"`
	VipLevel                int        `gorm:"column:vip_level" json:"vip_level"`
	Experience              int        `gorm:"column:experience" json:"experience"`
	MentalValue             int        `gorm:"column:mental_value" json:"mental_value"` // 精力
	QiValue                 int        `gorm:"column:qi_value" json:"qi_value"`         // 元力
	SpiritStone             int64      `gorm:"column:spirit_stone" json:"spirit_stone"`
	SpiritsCrystal          int64      `gorm:"column:spirits_crystal" json:"spirits_crystal"`
	SpiritsMilk             int64      `gorm:"column:spirits_milk" json:"spirits_milk"`
	SignInCoin              int64      `gorm:"column:sign_in_coin" json:"sign_in_coin"`
	TowerScore              int64      `gorm:"column:tower_score" json:"tower_score"`
	Reputation              int64      `gorm:"column:reputation" json:"reputation"`
	Tongxin                 int        `gorm:"column:tongxin" json:"tongxin"`
	ShaQi                   int        `gorm:"column:sha_qi" json:"sha_qi"`                               // 煞气值(打劫)
	Stamina                 int        `gorm:"column:stamina" json:"stamina"`                             // 体力值
	FreeAttributePoints     int        `gorm:"column:free_attribute_points" json:"free_attribute_points"` // 自由属性点
	TotalRechargeAmount     int64      `gorm:"column:total_recharge_amount" json:"total_recharge_amount"`
	FirstRechargeTime       *time.Time `gorm:"column:first_recharge_time" json:"first_recharge_time"`
	MonthCardExpireTime     *time.Time `gorm:"column:month_card_expire_time" json:"month_card_expire_time"`
	MonthCardTotalDays      int        `gorm:"column:month_card_total_days" json:"month_card_total_days"`           // 月卡累计天数
	LastCheckInTime         *time.Time `gorm:"column:last_check_in_time" json:"last_check_in_time"`                 // 最新签到时间
	TotalCheckInDays        int        `gorm:"column:total_check_in_days" json:"total_check_in_days"`               // 总签到天数
	SeclusionPassExpireTime *time.Time `gorm:"column:seclusion_pass_expire_time" json:"seclusion_pass_expire_time"` // 闭关卡到期
	SeclusionPassTotalDays  int        `gorm:"column:seclusion_pass_total_days" json:"seclusion_pass_total_days"`   // 闭关卡累计天数
	CreatedAt               time.Time  `gorm:"column:created_at" json:"created_at"`
	UpdatedAt               time.Time  `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt               *time.Time `gorm:"column:deleted_at" json:"deleted_at"`
}

// GameRealmRow 境界行
type GameRealmRow struct {
	ID         uint   `gorm:"column:id"`
	Name       string `gorm:"column:name"`
	Level      int    `gorm:"column:level"`
	StageLevel int    `gorm:"column:stage_level"`
}

// CurrencyLogRow 货币流水行
type CurrencyLogRow struct {
	UID          uint      `gorm:"column:uid"`
	CurrencyCode string    `gorm:"column:currency_code"`
	ChangeAmount int64     `gorm:"column:change_amount"`
	Detail       string    `gorm:"column:detail"`
	CreatedAt    time.Time `gorm:"column:created_at"`
}

// GetAllUsers 获取所有用户基础信息
func GetAllUsers() ([]GameUserRow, error) {
	g := db.GetGame()
	if g == nil {
		return nil, ErrGameDBNil
	}
	var rows []GameUserRow
	err := g.Table("users").
		Select("id, name, gender, status, realm_id, vip_level, experience, mental_value, qi_value, " +
			"spirit_stone, spirits_crystal, spirits_milk, sign_in_coin, tower_score, reputation, tongxin, sha_qi, stamina, free_attribute_points, " +
			"total_recharge_amount, first_recharge_time, month_card_expire_time, month_card_total_days, " +
			"last_check_in_time, total_check_in_days, seclusion_pass_expire_time, seclusion_pass_total_days, " +
			"created_at, updated_at, deleted_at").
		Find(&rows).Error
	return rows, err
}

// GetUserByID 获取单个用户
func GetUserByID(uid uint) (*GameUserRow, error) {
	g := db.GetGame()
	if g == nil {
		return nil, ErrGameDBNil
	}
	var row GameUserRow
	err := g.Table("users").Where("id = ?", uid).First(&row).Error
	if err != nil {
		return nil, err
	}
	return &row, nil
}

// GetAllRealms 获取所有境界
func GetAllRealms() ([]GameRealmRow, error) {
	g := db.GetGame()
	if g == nil {
		return nil, ErrGameDBNil
	}
	var rows []GameRealmRow
	err := g.Table("realms").Select("id, name, level, stage_level").Order("level ASC").Find(&rows).Error
	return rows, err
}

// RealmDistItem 境界分布统计项
type RealmDistItem struct {
	RealmID int    `gorm:"column:realm_id" json:"realm_id"`
	Name    string `gorm:"column:name"     json:"name"`
	Count   int    `gorm:"column:count"    json:"count"`
}

// GetRealmDistribution 获取境界分布
func GetRealmDistribution() ([]RealmDistItem, error) {
	g := db.GetGame()
	if g == nil {
		return nil, ErrGameDBNil
	}
	var items []RealmDistItem
	err := g.Table("users u").
		Select("u.realm_id, r.name, COUNT(*) as count").
		Joins("LEFT JOIN realms r ON r.id = u.realm_id").
		Group("u.realm_id, r.name").
		Order("u.realm_id ASC").
		Find(&items).Error
	return items, err
}

// GetCurrencyLogs 获取指定用户的货币流水（最近 N 天）
func GetCurrencyLogs(uid uint, days int) ([]CurrencyLogRow, error) {
	g := db.GetGame()
	if g == nil {
		return nil, ErrGameDBNil
	}
	var rows []CurrencyLogRow
	err := g.Table("user_currency_logs").
		Where("uid = ? AND created_at >= DATE_SUB(NOW(), INTERVAL ? DAY)", uid, days).
		Order("created_at DESC").
		Find(&rows).Error
	return rows, err
}

// GetTotalUserCount 总用户数
func GetTotalUserCount() (int64, error) {
	g := db.GetGame()
	if g == nil {
		return 0, ErrGameDBNil
	}
	var count int64
	err := g.Table("users").Count(&count).Error
	return count, err
}

// GetPayingUserCount 付费用户数
func GetPayingUserCount() (int64, error) {
	g := db.GetGame()
	if g == nil {
		return 0, ErrGameDBNil
	}
	var count int64
	err := g.Table("users").Where("total_recharge_amount > 0").Count(&count).Error
	return count, err
}

// TotalRechargeResult 充值总额结果
type TotalRechargeResult struct {
	Total int64 `gorm:"column:total"`
}

// GetTotalRechargeSum 充值总额
func GetTotalRechargeSum() (int64, error) {
	g := db.GetGame()
	if g == nil {
		return 0, ErrGameDBNil
	}
	var result TotalRechargeResult
	err := g.Table("recharge_orders").
		Select("COALESCE(SUM(pay_amount), 0) as total").
		Where(realRechargeOrderWhere).
		Scan(&result).Error
	return result.Total, err
}

// GetRechargeRevenueTotal 获取指定充值来源口径的总营收
func GetRechargeRevenueTotal(filter string) (int64, error) {
	g := db.GetGame()
	if g == nil {
		return 0, ErrGameDBNil
	}
	var result TotalRechargeResult
	err := g.Table("recharge_orders").
		Select("COALESCE(SUM(pay_amount), 0) as total").
		Where(rechargeOrderWhere(filter)).
		Scan(&result).Error
	return result.Total, err
}

// GetVipUserCount VIP用户数（vip_level > 0）
func GetVipUserCount() (int64, error) {
	g := db.GetGame()
	if g == nil {
		return 0, ErrGameDBNil
	}
	var count int64
	err := g.Table("users").Where("vip_level > 0").Count(&count).Error
	return count, err
}

// GetRealPayingUserCount 付费用户数（按充值订单过滤口径）
func GetRealPayingUserCount(filter ...string) (int64, error) {
	g := db.GetGame()
	if g == nil {
		return 0, ErrGameDBNil
	}
	orderFilter := string(RechargeOrderFilterOfficial)
	if len(filter) > 0 {
		orderFilter = filter[0]
	}
	var count int64
	err := g.Table("recharge_orders").
		Where(rechargeOrderWhere(orderFilter)).
		Distinct("uid").
		Count(&count).Error
	return count, err
}

// ========== 留存分析 ==========

// RetentionCohortRow 留存分析某日注册用户的后续活跃
type RetentionCohortRow struct {
	CohortDate string `json:"cohort_date"`
	CohortSize int    `json:"cohort_size"`
	Day1       int    `json:"day1"`
	Day3       int    `json:"day3"`
	Day7       int    `json:"day7"`
	Day14      int    `json:"day14"`
	Day30      int    `json:"day30"`
}

// GetRetentionData 获取留存数据（按注册日分组）
func GetRetentionData(days int) ([]RetentionCohortRow, error) {
	g := db.GetGame()
	if g == nil {
		return nil, ErrGameDBNil
	}
	if days <= 0 {
		days = 30
	}

	// 获取每日注册用户数
	type cohortCount struct {
		Date  string `gorm:"column:date"`
		Count int    `gorm:"column:cnt"`
	}
	var cohorts []cohortCount
	err := g.Table("users").
		Select("DATE(created_at) as date, COUNT(*) as cnt").
		Where("created_at >= DATE_SUB(NOW(), INTERVAL ? DAY)", days).
		Group("DATE(created_at)").
		Order("date ASC").
		Find(&cohorts).Error
	if err != nil {
		return nil, err
	}

	var results []RetentionCohortRow
	for _, c := range cohorts {
		row := RetentionCohortRow{
			CohortDate: c.Date,
			CohortSize: c.Count,
		}
		// 对每个留存天数查询
		for _, d := range []struct {
			day    int
			target *int
		}{
			{1, &row.Day1}, {3, &row.Day3}, {7, &row.Day7}, {14, &row.Day14}, {30, &row.Day30},
		} {
			type cnt struct {
				V int `gorm:"column:v"`
			}
			var r cnt
			g.Raw(`
				SELECT COUNT(DISTINCT ucl.uid) as v
				FROM user_currency_logs ucl
				INNER JOIN users u ON u.id = ucl.uid
				WHERE DATE(u.created_at) = ?
				  AND DATE(ucl.created_at) = DATE_ADD(?, INTERVAL ? DAY)
			`, c.Date, c.Date, d.day).Scan(&r)
			if g.Error != nil {
				return nil, g.Error
			}
			*d.target = r.V
		}
		results = append(results, row)
	}
	return results, nil
}

// ========== 转化漏斗 ==========

// FunnelData 漏斗数据
type FunnelData struct {
	TotalRegistered int `json:"total_registered"` // 总注册
	FirstBattle     int `json:"first_battle"`     // 首战 (有战斗日志)
	FirstBreak      int `json:"first_break"`      // 首次突破 (realm_id > 1)
	FirstPromotion  int `json:"first_promotion"`  // 首次晋级 (stage_level >= 2)
	JoinedGuild     int `json:"joined_guild"`     // 加入宗门
	FirstRecharge   int `json:"first_recharge"`   // 充值
}

// GetFunnelData 获取转化漏斗数据
func GetFunnelData(days int) (*FunnelData, error) {
	g := db.GetGame()
	if g == nil {
		return nil, ErrGameDBNil
	}

	data := &FunnelData{}
	where := ""
	var args []interface{}
	if days > 0 {
		where = "created_at >= DATE_SUB(NOW(), INTERVAL ? DAY)"
		args = append(args, days)
	}

	// 总注册
	q := g.Table("users")
	if where != "" {
		q = q.Where(where, args...)
	}
	var totalCount int64
	q.Count(&totalCount)
	data.TotalRegistered = int(totalCount)

	// 首战: 有战斗相关货币日志
	type cnt struct {
		V int `gorm:"column:v"`
	}
	var battleCnt cnt
	battleQ := `
		SELECT COUNT(DISTINCT u.id) as v FROM users u
		WHERE EXISTS (
			SELECT 1 FROM user_currency_logs ucl
			WHERE ucl.uid = u.id
			AND ucl.detail REGEXP '战斗|副本|击杀|boss|妖塔|秘境|怪物|挑战'
		)`
	if days > 0 {
		battleQ += " AND u.created_at >= DATE_SUB(NOW(), INTERVAL ? DAY)"
		if err := g.Raw(battleQ, days).Scan(&battleCnt).Error; err != nil {
			return nil, err
		}
	} else {
		if err := g.Raw(battleQ).Scan(&battleCnt).Error; err != nil {
			return nil, err
		}
	}
	data.FirstBattle = battleCnt.V

	// 首次突破: realm_id > 1
	var breakCnt cnt
	breakQ := g.Table("users").Select("COUNT(*) as v").Where("realm_id > 1")
	if days > 0 {
		breakQ = breakQ.Where(where, args...)
	}
	if err := breakQ.Scan(&breakCnt).Error; err != nil {
		return nil, err
	}
	data.FirstBreak = breakCnt.V

	// 首次晋级: 到达 stage_level >= 2 的境界
	var promoCnt cnt
	promoQ := `SELECT COUNT(DISTINCT u.id) as v FROM users u
		INNER JOIN realms r ON r.id = u.realm_id
		WHERE r.stage_level >= 2`
	if days > 0 {
		promoQ += " AND u.created_at >= DATE_SUB(NOW(), INTERVAL ? DAY)"
		if err := g.Raw(promoQ, days).Scan(&promoCnt).Error; err != nil {
			return nil, err
		}
	} else {
		if err := g.Raw(promoQ).Scan(&promoCnt).Error; err != nil {
			return nil, err
		}
	}
	data.FirstPromotion = promoCnt.V

	// 加入宗门
	var guildCnt cnt
	guildQ := `SELECT COUNT(DISTINCT gm.user_id) as v FROM guild_members gm`
	if days > 0 {
		guildQ += " INNER JOIN users u ON u.id = gm.user_id WHERE u.created_at >= DATE_SUB(NOW(), INTERVAL ? DAY)"
		if err := g.Raw(guildQ, days).Scan(&guildCnt).Error; err != nil {
			return nil, err
		}
	} else {
		if err := g.Raw(guildQ).Scan(&guildCnt).Error; err != nil {
			return nil, err
		}
	}
	data.JoinedGuild = guildCnt.V

	// 首充
	var rechargeCnt cnt
	rechargeQ := `SELECT COUNT(DISTINCT ro.uid) as v FROM recharge_orders ro WHERE ` + realRechargeOrderWhereAlias
	if days > 0 {
		rechargeQ += " AND EXISTS (SELECT 1 FROM users u WHERE u.id = ro.uid AND u.created_at >= DATE_SUB(NOW(), INTERVAL ? DAY))"
		if err := g.Raw(rechargeQ, days).Scan(&rechargeCnt).Error; err != nil {
			return nil, err
		}
	} else {
		if err := g.Raw(rechargeQ).Scan(&rechargeCnt).Error; err != nil {
			return nil, err
		}
	}
	data.FirstRecharge = rechargeCnt.V

	return data, nil
}

// GetNewUsersCount N 天内新注册用户数
func GetNewUsersCount(days int) (int64, error) {
	g := db.GetGame()
	if g == nil {
		return 0, ErrGameDBNil
	}
	var count int64
	err := g.Table("users").
		Where("created_at >= DATE_SUB(NOW(), INTERVAL ? DAY)", days).
		Count(&count).Error
	return count, err
}

// GuildDistItem 门派分布统计项
type GuildDistItem struct {
	GuildID   uint   `gorm:"column:guild_id" json:"guild_id"`
	GuildName string `gorm:"column:guild_name" json:"guild_name"`
	Count     int    `gorm:"column:count" json:"count"`
}

// GetGuildDistribution 获取门派人数分布
func GetGuildDistribution() ([]GuildDistItem, error) {
	g := db.GetGame()
	if g == nil {
		return nil, ErrGameDBNil
	}
	var items []GuildDistItem
	err := g.Table("guild_members gm").
		Select("gm.guild_id, g.name as guild_name, COUNT(*) as count").
		Joins("LEFT JOIN guilds g ON g.id = gm.guild_id").
		Group("gm.guild_id, g.name").
		Order("count DESC").
		Find(&items).Error
	return items, err
}

// ========== 快照采集 & 社交分析所需查询 ==========

// CurrencyAggRow 某用户某天某币种的收入/支出汇总
type CurrencyAggRow struct {
	CurrencyCode string `gorm:"column:currency_code"`
	Income       int64  `gorm:"column:income"`
	Expense      int64  `gorm:"column:expense"`
}

// GetDailyCurrencyAgg 获取某用户某天的货币收支汇总
func GetDailyCurrencyAgg(uid uint, date string) ([]CurrencyAggRow, error) {
	g := db.GetGame()
	if g == nil {
		return nil, ErrGameDBNil
	}
	var rows []CurrencyAggRow
	err := g.Table("user_currency_logs").
		Select("currency_code, "+
			"COALESCE(SUM(CASE WHEN change_amount > 0 THEN change_amount ELSE 0 END), 0) as income, "+
			"COALESCE(SUM(CASE WHEN change_amount < 0 THEN ABS(change_amount) ELSE 0 END), 0) as expense").
		Where("uid = ? AND DATE(created_at) = ?", uid, date).
		Group("currency_code").
		Find(&rows).Error
	return rows, err
}

// DailyActionCountRow 某用户某天各行为类别的操作次数
type DailyActionCountRow struct {
	Category string `gorm:"column:category"`
	Count    int    `gorm:"column:cnt"`
}

// GetDailyActionCounts 根据货币流水的 detail 关键词估算行为分类计数
// 战斗类: 战斗/副本/击杀/boss/妖塔/秘境
// 制造类: 炼器/炼丹/制造/锻造/分解
// 社交类: 宗门/好友/道侣/赠送/红包
// 经济类: 交易/拍卖/商店/购买/出售/兑换
func GetDailyActionCounts(uid uint, date string) (battle, craft, social, economy int, err error) {
	g := db.GetGame()
	if g == nil {
		return 0, 0, 0, 0, ErrGameDBNil
	}

	type countResult struct {
		Cnt int `gorm:"column:cnt"`
	}

	categories := []struct {
		keywords string
		target   *int
	}{
		{"detail REGEXP '战斗|副本|击杀|boss|妖塔|秘境|怪物|挑战'", &battle},
		{"detail REGEXP '炼器|炼丹|制造|锻造|分解|合成|强化'", &craft},
		{"detail REGEXP '宗门|好友|道侣|赠送|红包|师徒|帮派'", &social},
		{"detail REGEXP '交易|拍卖|商店|购买|出售|兑换|集市'", &economy},
	}

	for _, cat := range categories {
		var r countResult
		if e := g.Table("user_currency_logs").
			Select("COUNT(*) as cnt").
			Where("uid = ? AND DATE(created_at) = ? AND "+cat.keywords, uid, date).
			Scan(&r).Error; e != nil {
			return 0, 0, 0, 0, e
		}
		*cat.target = r.Cnt
	}
	return
}

// FriendCountRow 好友数统计
type FriendCountRow struct {
	Total     int `gorm:"column:total"`
	Intimate  int `gorm:"column:intimate"`  // 亲密度>50 的好友数
	Companion int `gorm:"column:companion"` // 道侣数
}

// GetFriendStats 获取用户好友统计
func GetFriendStats(uid uint) (*FriendCountRow, error) {
	g := db.GetGame()
	if g == nil {
		return nil, ErrGameDBNil
	}
	var row FriendCountRow
	if err := g.Table("friendships").
		Select("COUNT(*) as total, "+
			"SUM(CASE WHEN intimacy > 50 THEN 1 ELSE 0 END) as intimate, "+
			"SUM(CASE WHEN relation_type = 2 THEN 1 ELSE 0 END) as companion").
		Where("user_id = ? AND deleted_at IS NULL", uid).
		Scan(&row).Error; err != nil {
		return nil, err
	}
	return &row, nil
}

// GuildMemberInfo 宗门成员信息
type GuildMemberInfo struct {
	GuildID           uint   `gorm:"column:guild_id"`
	GuildName         string `gorm:"column:guild_name"`
	RoleName          string `gorm:"column:role_name"`
	ContributionTotal int64  `gorm:"column:contribution_total"`
}

// GetGuildMemberInfo 获取用户的宗门成员信息
func GetGuildMemberInfo(uid uint) (*GuildMemberInfo, error) {
	g := db.GetGame()
	if g == nil {
		return nil, ErrGameDBNil
	}
	var info GuildMemberInfo
	err := g.Table("guild_members gm").
		Select("gm.guild_id, g.name as guild_name, gr.name as role_name, gm.contribution_total").
		Joins("LEFT JOIN guilds g ON g.id = gm.guild_id").
		Joins("LEFT JOIN guild_roles gr ON gr.id = gm.role_id").
		Where("gm.user_id = ? AND gm.deleted_at IS NULL", uid).
		First(&info).Error
	if err != nil {
		return nil, err
	}
	return &info, nil
}

// GetGuildLogCount 获取用户最近 N 天的宗门日志条数（社交行为指标）
func GetGuildLogCount(uid uint, days int) (int, error) {
	g := db.GetGame()
	if g == nil {
		return 0, ErrGameDBNil
	}
	type countResult struct {
		Cnt int `gorm:"column:cnt"`
	}
	var r countResult
	if err := g.Table("guild_logs").
		Select("COUNT(*) as cnt").
		Where("user_id = ? AND created_at >= DATE_SUB(NOW(), INTERVAL ? DAY)", uid, days).
		Scan(&r).Error; err != nil {
		return 0, err
	}
	return r.Cnt, nil
}

// GetActiveUserIDs 获取最近 N 天有行为记录的活跃用户 ID 列表
// 口径: 货币流水、宗门日志、副本通关、市场交易、成功充值。
func GetActiveUserIDs(days int) ([]uint, error) {
	g := db.GetGame()
	if g == nil {
		return nil, ErrGameDBNil
	}
	if days <= 0 {
		days = 3
	}
	var ids []uint
	err := g.Raw(`
		SELECT DISTINCT uid
		FROM (
			SELECT uid
			FROM user_currency_logs
			WHERE created_at >= DATE_SUB(NOW(), INTERVAL ? DAY)

			UNION ALL

			SELECT user_id as uid
			FROM guild_logs
			WHERE created_at >= DATE_SUB(NOW(), INTERVAL ? DAY)

			UNION ALL

			SELECT user_id as uid
			FROM dungeon_records
			WHERE completed_at >= DATE_SUB(NOW(), INTERVAL ? DAY)

			UNION ALL

			SELECT buyer_user_id as uid
			FROM market_trade_records
			WHERE created_at >= DATE_SUB(NOW(), INTERVAL ? DAY)

			UNION ALL

			SELECT seller_user_id as uid
			FROM market_trade_records
			WHERE created_at >= DATE_SUB(NOW(), INTERVAL ? DAY)

			UNION ALL

			SELECT uid
			FROM recharge_orders
			WHERE `+realRechargeOrderWhere+` AND pay_time >= DATE_SUB(NOW(), INTERVAL ? DAY)
		) active_users
		ORDER BY uid ASC
	`, days, days, days, days, days, days).Scan(&ids).Error
	return ids, err
}

// ========== 排行榜 & 扩展统计 ==========

// PlayerRankItem 玩家排名项
type PlayerRankItem struct {
	ID                  uint   `gorm:"column:id" json:"id"`
	Name                string `gorm:"column:name" json:"name"`
	RealmName           string `gorm:"column:realm_name" json:"realm_name"`
	StageLevel          int    `gorm:"column:stage_level" json:"stage_level"`
	VipLevel            int    `gorm:"column:vip_level" json:"vip_level"`
	SpiritStone         int64  `gorm:"column:spirit_stone" json:"spirit_stone"`
	TotalRechargeAmount int64  `gorm:"column:total_recharge" json:"total_recharge"`
}

// GetPlayerRanking 玩家排行榜
// sortBy: realm(境界), recharge(充值), wealth(财富), vip
func GetPlayerRanking(sortBy string, limit int) ([]PlayerRankItem, error) {
	g := db.GetGame()
	if g == nil {
		return nil, ErrGameDBNil
	}
	if limit <= 0 || limit > 100 {
		limit = 50
	}

	orderMap := map[string]string{
		"realm":    "r.level DESC, u.experience DESC",
		"recharge": "u.total_recharge_amount DESC",
		"wealth":   "u.spirit_stone DESC",
		"vip":      "u.vip_level DESC, u.total_recharge_amount DESC",
	}
	order, ok := orderMap[sortBy]
	if !ok {
		order = orderMap["realm"]
	}

	var items []PlayerRankItem
	err := g.Table("users u").
		Select("u.id, u.name, r.name as realm_name, r.stage_level, u.vip_level, u.spirit_stone, u.total_recharge_amount as total_recharge").
		Joins("LEFT JOIN realms r ON r.id = u.realm_id").
		Order(order).
		Limit(limit).
		Find(&items).Error
	return items, err
}

// GuildRankItem 宗门排名项
type GuildRankItem struct {
	ID          uint   `gorm:"column:id" json:"id"`
	Name        string `gorm:"column:name" json:"name"`
	Prestige    int64  `gorm:"column:prestige" json:"prestige"`
	SpiritStone int64  `gorm:"column:spirit_stone" json:"spirit_stone"`
	MemberCount int    `gorm:"column:member_count" json:"member_count"`
	LevelID     int    `gorm:"column:level_id" json:"level_id"`
}

// GetGuildRanking 宗门排行榜
func GetGuildRanking(sortBy string, limit int) ([]GuildRankItem, error) {
	g := db.GetGame()
	if g == nil {
		return nil, ErrGameDBNil
	}
	if limit <= 0 || limit > 100 {
		limit = 30
	}

	orderMap := map[string]string{
		"prestige": "prestige DESC",
		"members":  "member_count DESC",
		"wealth":   "spirit_stone DESC",
	}
	order, ok := orderMap[sortBy]
	if !ok {
		order = orderMap["prestige"]
	}

	var items []GuildRankItem
	err := g.Table("guilds").
		Select("id, name, prestige, spirit_stone, member_count, level_id").
		Order(order).
		Limit(limit).
		Find(&items).Error
	return items, err
}

// DailyCountItem 每日计数项
type DailyCountItem struct {
	Date  string `gorm:"column:date" json:"date"`
	Count int    `gorm:"column:count" json:"count"`
}

// GetNewUsersTrend 每日新增用户数趋势
func GetNewUsersTrend(days int) ([]DailyCountItem, error) {
	g := db.GetGame()
	if g == nil {
		return nil, ErrGameDBNil
	}
	var items []DailyCountItem
	err := g.Table("users").
		Select("DATE(created_at) as date, COUNT(*) as count").
		Where("created_at >= DATE_SUB(NOW(), INTERVAL ? DAY)", days).
		Group("DATE(created_at)").
		Order("date ASC").
		Find(&items).Error
	return items, err
}

// DailyRevenueItem 每日营收项
type DailyRevenueItem struct {
	Date        string `gorm:"column:date" json:"date"`
	Revenue     int64  `gorm:"column:revenue" json:"revenue"`
	OrderCount  int    `gorm:"column:order_count" json:"order_count"`
	PayingUsers int    `gorm:"column:paying_users" json:"paying_users"`
}

// GetRechargeRevenueTrend 每日充值营收趋势
func GetRechargeRevenueTrend(days int, filter ...string) ([]DailyRevenueItem, error) {
	g := db.GetGame()
	if g == nil {
		return nil, ErrGameDBNil
	}
	orderFilter := string(RechargeOrderFilterOfficial)
	if len(filter) > 0 {
		orderFilter = filter[0]
	}
	var items []DailyRevenueItem
	err := g.Table("recharge_orders").
		Select("DATE(pay_time) as date, COALESCE(SUM(pay_amount),0) as revenue, COUNT(*) as order_count, COUNT(DISTINCT uid) as paying_users").
		Where(rechargeOrderWhere(orderFilter)+" AND pay_time >= DATE_SUB(NOW(), INTERVAL ? DAY)", days).
		Group("DATE(pay_time)").
		Order("date ASC").
		Find(&items).Error
	return items, err
}

// PackageStatItem 充值套餐统计项
type PackageStatItem struct {
	PackageID    uint   `gorm:"column:package_id" json:"package_id"`
	PackageName  string `gorm:"column:package_name" json:"package_name"`
	PackageType  string `gorm:"column:package_type" json:"package_type"`
	Price        int64  `gorm:"column:price" json:"price"`
	SoldCount    int    `gorm:"column:sold_count" json:"sold_count"`
	TotalRevenue int64  `gorm:"column:total_revenue" json:"total_revenue"`
	BuyerCount   int    `gorm:"column:buyer_count" json:"buyer_count"`
}

// GetRechargePackageStats 充值套餐销售统计
func GetRechargePackageStats(days int, filter ...string) ([]PackageStatItem, error) {
	g := db.GetGame()
	if g == nil {
		return nil, ErrGameDBNil
	}
	orderFilter := string(RechargeOrderFilterOfficial)
	if len(filter) > 0 {
		orderFilter = filter[0]
	}
	var items []PackageStatItem
	query := g.Table("recharge_orders ro").
		Select("ro.package_id, rp.name as package_name, rp.package_type, rp.price_amount as price, " +
			"COUNT(*) as sold_count, COALESCE(SUM(ro.pay_amount),0) as total_revenue, COUNT(DISTINCT ro.uid) as buyer_count").
		Joins("LEFT JOIN recharge_packages rp ON rp.id = ro.package_id").
		Where(rechargeOrderWhereAlias(orderFilter))
	if days > 0 {
		query = query.Where("ro.pay_time >= DATE_SUB(NOW(), INTERVAL ? DAY)", days)
	}
	err := query.Group("ro.package_id, rp.name, rp.package_type, rp.price_amount").
		Order("total_revenue DESC").
		Find(&items).Error
	return items, err
}

// VipDistItem VIP等级分布
type VipDistItem struct {
	VipLevel int `gorm:"column:vip_level" json:"vip_level"`
	Count    int `gorm:"column:count" json:"count"`
}

// GetVipDistribution VIP等级分布
func GetVipDistribution() ([]VipDistItem, error) {
	g := db.GetGame()
	if g == nil {
		return nil, ErrGameDBNil
	}
	var items []VipDistItem
	err := g.Table("users").
		Select("vip_level, COUNT(*) as count").
		Group("vip_level").
		Order("vip_level ASC").
		Find(&items).Error
	return items, err
}

// RealmChurnItem 境界流失统计
type RealmChurnItem struct {
	RealmID    int     `gorm:"column:realm_id" json:"realm_id"`
	RealmName  string  `gorm:"column:realm_name" json:"realm_name"`
	StageLevel int     `gorm:"column:stage_level" json:"stage_level"`
	Total      int     `gorm:"column:total" json:"total"`
	Churned    int     `gorm:"column:churned" json:"churned"`
	ChurnRate  float64 `json:"churn_rate"`
}

// GetRealmChurn 各境界流失率（到达该境界后7天内未再活跃）
func GetRealmChurn(inactiveDays int) ([]RealmChurnItem, error) {
	g := db.GetGame()
	if g == nil {
		return nil, ErrGameDBNil
	}
	if inactiveDays <= 0 {
		inactiveDays = 7
	}
	var items []RealmChurnItem
	err := g.Table("users u").
		Select("u.realm_id, r.name as realm_name, r.stage_level, "+
			"COUNT(*) as total, "+
			"SUM(CASE WHEN u.updated_at < DATE_SUB(NOW(), INTERVAL ? DAY) THEN 1 ELSE 0 END) as churned",
			inactiveDays).
		Joins("LEFT JOIN realms r ON r.id = u.realm_id").
		Group("u.realm_id, r.name, r.stage_level").
		Order("r.stage_level ASC, u.realm_id ASC").
		Find(&items).Error
	if err != nil {
		return nil, err
	}
	// compute churn rate
	for i := range items {
		if items[i].Total > 0 {
			items[i].ChurnRate = float64(items[i].Churned) / float64(items[i].Total)
		}
	}
	return items, nil
}

// RealmStageDistItem 大境界分布
type RealmStageDistItem struct {
	StageLevel int    `gorm:"column:stage_level" json:"stage_level"`
	StageName  string `json:"stage_name"`
	Count      int    `gorm:"column:count" json:"count"`
}

// stageNames 大境界名称映射
var stageNames = map[int]string{
	0: "凡人", 1: "练气", 2: "筑基", 3: "金丹", 4: "元婴",
	5: "化神", 6: "洞虚", 7: "大乘", 8: "仙", 9: "金仙", 10: "道祖",
}

// GetRealmStageDistribution 按大境界汇总分布
func GetRealmStageDistribution() ([]RealmStageDistItem, error) {
	g := db.GetGame()
	if g == nil {
		return nil, ErrGameDBNil
	}
	var items []RealmStageDistItem
	err := g.Table("users u").
		Select("r.stage_level, COUNT(*) as count").
		Joins("LEFT JOIN realms r ON r.id = u.realm_id").
		Group("r.stage_level").
		Order("r.stage_level ASC").
		Find(&items).Error
	if err != nil {
		return nil, err
	}
	for i := range items {
		if name, ok := stageNames[items[i].StageLevel]; ok {
			items[i].StageName = name
		}
	}
	return items, nil
}

// GetDAUFromGame 从游戏库活跃数据估算某天的 DAU（基于货币流水）
func GetDAUFromGame(date string) (int, error) {
	g := db.GetGame()
	if g == nil {
		return 0, ErrGameDBNil
	}
	type countRow struct {
		Cnt int `gorm:"column:cnt"`
	}
	var r countRow
	err := g.Table("user_currency_logs").
		Select("COUNT(DISTINCT uid) as cnt").
		Where("DATE(created_at) = ?", date).
		Scan(&r).Error
	return r.Cnt, err
}

// GetMAUFromGame 从游戏库估算当月 MAU
func GetMAUFromGame(days int) (int, error) {
	g := db.GetGame()
	if g == nil {
		return 0, ErrGameDBNil
	}
	if days <= 0 {
		days = 30
	}
	type countRow struct {
		Cnt int `gorm:"column:cnt"`
	}
	var r countRow
	err := g.Table("user_currency_logs").
		Select("COUNT(DISTINCT uid) as cnt").
		Where("created_at >= DATE_SUB(NOW(), INTERVAL ? DAY)", days).
		Scan(&r).Error
	return r.Cnt, err
}

// ========== 境界进阶分析 ==========

// RealmProgressItem 境界进阶统计项
type RealmProgressItem struct {
	RealmID       int     `gorm:"column:realm_id"     json:"realm_id"`
	RealmName     string  `gorm:"column:realm_name"   json:"realm_name"`
	StageLevel    int     `gorm:"column:stage_level"  json:"stage_level"`
	UserCount     int     `gorm:"column:user_count"   json:"user_count"`
	AvgDaysPlayed float64 `gorm:"column:avg_days_played"  json:"avg_days_played"`
	AvgRecharge   float64 `gorm:"column:avg_recharge"     json:"avg_recharge"`
	AvgVipLevel   float64 `gorm:"column:avg_vip_level"    json:"avg_vip_level"`
	ActiveRate    float64 `json:"active_rate"`
	InactiveCount int     `gorm:"column:inactive_count"   json:"inactive_count"`
}

// GetRealmProgressAnalysis 境界进阶综合分析
func GetRealmProgressAnalysis(inactiveDays int) ([]RealmProgressItem, error) {
	g := db.GetGame()
	if g == nil {
		return nil, ErrGameDBNil
	}
	if inactiveDays <= 0 {
		inactiveDays = 7
	}

	var items []RealmProgressItem
	err := g.Table("users u").
		Select(`u.realm_id,
			r.name as realm_name,
			r.stage_level,
			COUNT(*) as user_count,
			AVG(DATEDIFF(NOW(), u.created_at)) as avg_days_played,
			AVG(u.total_recharge_amount) / 100.0 as avg_recharge,
			AVG(u.vip_level) as avg_vip_level,
			SUM(CASE WHEN u.updated_at < DATE_SUB(NOW(), INTERVAL ? DAY) THEN 1 ELSE 0 END) as inactive_count`,
			inactiveDays).
		Joins("LEFT JOIN realms r ON r.id = u.realm_id").
		Group("u.realm_id, r.name, r.stage_level").
		Order("r.stage_level ASC, u.realm_id ASC").
		Find(&items).Error
	if err != nil {
		return nil, err
	}
	for i := range items {
		if items[i].UserCount > 0 {
			items[i].ActiveRate = float64(items[i].UserCount-items[i].InactiveCount) / float64(items[i].UserCount) * 100
		}
	}
	return items, nil
}

// RealmPayCorrelationItem 境界-付费关联项
type RealmPayCorrelationItem struct {
	StageLevel    int     `gorm:"column:stage_level"    json:"stage_level"`
	StageName     string  `json:"stage_name"`
	UserCount     int     `gorm:"column:user_count"     json:"user_count"`
	PayingCount   int     `gorm:"column:paying_count"   json:"paying_count"`
	PayRate       float64 `json:"pay_rate"`
	AvgRecharge   float64 `gorm:"column:avg_recharge"   json:"avg_recharge"`
	TotalRecharge int64   `gorm:"column:total_recharge" json:"total_recharge"`
}

// GetRealmPayCorrelation 境界-付费关联分析（按大境界）
func GetRealmPayCorrelation() ([]RealmPayCorrelationItem, error) {
	g := db.GetGame()
	if g == nil {
		return nil, ErrGameDBNil
	}

	var items []RealmPayCorrelationItem
	err := g.Table("users u").
		Select(`r.stage_level,
			COUNT(*) as user_count,
			SUM(CASE WHEN u.total_recharge_amount > 0 THEN 1 ELSE 0 END) as paying_count,
			AVG(u.total_recharge_amount) / 100.0 as avg_recharge,
			SUM(u.total_recharge_amount) as total_recharge`).
		Joins("LEFT JOIN realms r ON r.id = u.realm_id").
		Group("r.stage_level").
		Order("r.stage_level ASC").
		Find(&items).Error
	if err != nil {
		return nil, err
	}
	for i := range items {
		if name, ok := stageNames[items[i].StageLevel]; ok {
			items[i].StageName = name
		}
		if items[i].UserCount > 0 {
			items[i].PayRate = float64(items[i].PayingCount) / float64(items[i].UserCount) * 100
		}
	}
	return items, nil
}

// ========== 经济健康度分析 ==========

// CurrencyFlowItem 货币流通概况
type CurrencyFlowItem struct {
	CurrencyCode string  `gorm:"column:currency_code" json:"currency_code"`
	TotalIncome  int64   `gorm:"column:total_income"  json:"total_income"`
	TotalExpense int64   `gorm:"column:total_expense" json:"total_expense"`
	NetFlow      int64   `json:"net_flow"`
	TxCount      int     `gorm:"column:tx_count"      json:"tx_count"`
	UniqueUsers  int     `gorm:"column:unique_users"  json:"unique_users"`
	AvgIncome    float64 `json:"avg_income"`
	AvgExpense   float64 `json:"avg_expense"`
}

// GetCurrencyFlow 获取各币种流通概况（最近N天）
func GetCurrencyFlow(days int) ([]CurrencyFlowItem, error) {
	g := db.GetGame()
	if g == nil {
		return nil, ErrGameDBNil
	}
	if days <= 0 {
		days = 30
	}
	var items []CurrencyFlowItem
	err := g.Table("user_currency_logs").
		Select(`currency_code,
			SUM(CASE WHEN change_amount > 0 THEN change_amount ELSE 0 END) as total_income,
			SUM(CASE WHEN change_amount < 0 THEN ABS(change_amount) ELSE 0 END) as total_expense,
			COUNT(*) as tx_count,
			COUNT(DISTINCT uid) as unique_users`).
		Where("created_at >= DATE_SUB(NOW(), INTERVAL ? DAY)", days).
		Group("currency_code").
		Find(&items).Error
	if err != nil {
		return nil, err
	}
	for i := range items {
		items[i].NetFlow = items[i].TotalIncome - items[i].TotalExpense
		if items[i].UniqueUsers > 0 {
			items[i].AvgIncome = float64(items[i].TotalIncome) / float64(items[i].UniqueUsers)
			items[i].AvgExpense = float64(items[i].TotalExpense) / float64(items[i].UniqueUsers)
		}
	}
	return items, nil
}

// CurrencyTrendItem 货币日趋势
type CurrencyTrendItem struct {
	Date         string `gorm:"column:date"          json:"date"`
	CurrencyCode string `gorm:"column:currency_code" json:"currency_code"`
	Income       int64  `gorm:"column:income"        json:"income"`
	Expense      int64  `gorm:"column:expense"       json:"expense"`
	TxCount      int    `gorm:"column:tx_count"      json:"tx_count"`
}

// GetCurrencyTrend 获取币种日趋势
func GetCurrencyTrend(currencyCode string, days int) ([]CurrencyTrendItem, error) {
	g := db.GetGame()
	if g == nil {
		return nil, ErrGameDBNil
	}
	if days <= 0 {
		days = 30
	}
	var items []CurrencyTrendItem
	err := g.Table("user_currency_logs").
		Select(`DATE(created_at) as date,
			currency_code,
			SUM(CASE WHEN change_amount > 0 THEN change_amount ELSE 0 END) as income,
			SUM(CASE WHEN change_amount < 0 THEN ABS(change_amount) ELSE 0 END) as expense,
			COUNT(*) as tx_count`).
		Where("currency_code = ? AND created_at >= DATE_SUB(NOW(), INTERVAL ? DAY)", currencyCode, days).
		Group("DATE(created_at), currency_code").
		Order("date ASC").
		Find(&items).Error
	return items, err
}

// WealthDistItem 财富分布区间
type WealthDistItem struct {
	Bracket    string  `json:"bracket"`
	UserCount  int     `gorm:"column:user_count"  json:"user_count"`
	Percentage float64 `json:"percentage"`
}

// GetWealthDistribution 获取灵石财富分布
func GetWealthDistribution() ([]WealthDistItem, error) {
	g := db.GetGame()
	if g == nil {
		return nil, ErrGameDBNil
	}

	type rawRow struct {
		Bracket   string `gorm:"column:bracket"`
		UserCount int    `gorm:"column:user_count"`
	}
	var rows []rawRow
	err := g.Table("users").
		Select(`CASE
			WHEN spirit_stone < 1000 THEN '0-1K'
			WHEN spirit_stone < 10000 THEN '1K-10K'
			WHEN spirit_stone < 100000 THEN '10K-100K'
			WHEN spirit_stone < 1000000 THEN '100K-1M'
			ELSE '1M+'
		END as bracket,
		COUNT(*) as user_count`).
		Group(`CASE
			WHEN spirit_stone < 1000 THEN '0-1K'
			WHEN spirit_stone < 10000 THEN '1K-10K'
			WHEN spirit_stone < 100000 THEN '10K-100K'
			WHEN spirit_stone < 1000000 THEN '100K-1M'
			ELSE '1M+'
		END`).
		Find(&rows).Error
	if err != nil {
		return nil, err
	}
	total := 0
	for _, r := range rows {
		total += r.UserCount
	}
	var items []WealthDistItem
	for _, r := range rows {
		pct := 0.0
		if total > 0 {
			pct = float64(r.UserCount) / float64(total) * 100
		}
		items = append(items, WealthDistItem{Bracket: r.Bracket, UserCount: r.UserCount, Percentage: pct})
	}
	return items, nil
}

// ========== 副本分析 ==========

// DungeonStatsItem 副本统计项
type DungeonStatsItem struct {
	DungeonID   int     `gorm:"column:dungeon_id"   json:"dungeon_id"`
	DungeonName string  `gorm:"column:dungeon_name" json:"dungeon_name"`
	DungeonType int     `gorm:"column:dungeon_type" json:"dungeon_type"`
	ClearCount  int     `gorm:"column:clear_count"  json:"clear_count"`
	UniqueUsers int     `gorm:"column:unique_users" json:"unique_users"`
	AvgClears   float64 `json:"avg_clears"`
}

// GetDungeonStats 获取副本通关统计（最近N天）
func GetDungeonStats(days int) ([]DungeonStatsItem, error) {
	g := db.GetGame()
	if g == nil {
		return nil, ErrGameDBNil
	}
	if days <= 0 {
		days = 30
	}
	var items []DungeonStatsItem
	q := g.Table("dungeon_records dr").
		Select(`dr.dungeon_id,
			d.name as dungeon_name,
			d.type as dungeon_type,
			COUNT(*) as clear_count,
			COUNT(DISTINCT dr.user_id) as unique_users`).
		Joins("LEFT JOIN dungeons d ON d.id = dr.dungeon_id").
		Group("dr.dungeon_id, d.name, d.type").
		Order("clear_count DESC")

	if days > 0 {
		q = q.Where("dr.completed_at >= DATE_SUB(NOW(), INTERVAL ? DAY)", days)
	}

	err := q.Find(&items).Error
	if err != nil {
		return nil, err
	}
	for i := range items {
		if items[i].UniqueUsers > 0 {
			items[i].AvgClears = float64(items[i].ClearCount) / float64(items[i].UniqueUsers)
		}
	}
	return items, nil
}

// DungeonDailyTrend 副本日趋势
type DungeonDailyTrend struct {
	Date       string `gorm:"column:date"        json:"date"`
	ClearCount int    `gorm:"column:clear_count" json:"clear_count"`
	UserCount  int    `gorm:"column:user_count"  json:"user_count"`
}

// GetDungeonDailyTrend 获取副本整体日趋势
func GetDungeonDailyTrend(days int) ([]DungeonDailyTrend, error) {
	g := db.GetGame()
	if g == nil {
		return nil, ErrGameDBNil
	}
	if days <= 0 {
		days = 30
	}
	var items []DungeonDailyTrend
	err := g.Table("dungeon_records").
		Select(`DATE(completed_at) as date,
			COUNT(*) as clear_count,
			COUNT(DISTINCT user_id) as user_count`).
		Where("completed_at >= DATE_SUB(NOW(), INTERVAL ? DAY)", days).
		Group("DATE(completed_at)").
		Order("date ASC").
		Find(&items).Error
	return items, err
}

// DungeonRealmDistItem 副本-境界参与分布
type DungeonRealmDistItem struct {
	DungeonName string `gorm:"column:dungeon_name" json:"dungeon_name"`
	StageName   string `json:"stage_name"`
	StageLevel  int    `gorm:"column:stage_level"  json:"stage_level"`
	UserCount   int    `gorm:"column:user_count"   json:"user_count"`
	ClearCount  int    `gorm:"column:clear_count"  json:"clear_count"`
}

// GetDungeonRealmDist 副本参与者境界分布
func GetDungeonRealmDist(days int) ([]DungeonRealmDistItem, error) {
	g := db.GetGame()
	if g == nil {
		return nil, ErrGameDBNil
	}
	if days <= 0 {
		days = 30
	}
	var items []DungeonRealmDistItem
	err := g.Table("dungeon_records dr").
		Select(`d.name as dungeon_name,
			r.stage_level,
			COUNT(DISTINCT dr.user_id) as user_count,
			COUNT(*) as clear_count`).
		Joins("LEFT JOIN dungeons d ON d.id = dr.dungeon_id").
		Joins("LEFT JOIN users u ON u.id = dr.user_id").
		Joins("LEFT JOIN realms r ON r.id = u.realm_id").
		Where("dr.completed_at >= DATE_SUB(NOW(), INTERVAL ? DAY)", days).
		Group("d.name, r.stage_level").
		Order("d.name ASC, r.stage_level ASC").
		Find(&items).Error
	if err != nil {
		return nil, err
	}
	for i := range items {
		if name, ok := stageNames[items[i].StageLevel]; ok {
			items[i].StageName = name
		}
	}
	return items, nil
}

// ========== 宗门健康度分析 ==========

// GuildHealthItem 宗门健康统计项
type GuildHealthItem struct {
	GuildID          int     `gorm:"column:guild_id"          json:"guild_id"`
	GuildName        string  `gorm:"column:guild_name"        json:"guild_name"`
	Level            int     `gorm:"column:level"             json:"level"`
	MemberCount      int     `gorm:"column:member_count"      json:"member_count"`
	MemberCapacity   int     `gorm:"column:member_capacity"   json:"member_capacity"`
	Prestige         int64   `gorm:"column:prestige"          json:"prestige"`
	ActiveMembers    int     `gorm:"column:active_members"    json:"active_members"`
	ActiveRate       float64 `json:"active_rate"`
	AvgContribution  float64 `gorm:"column:avg_contribution"  json:"avg_contribution"`
	TotalContrib     int64   `gorm:"column:total_contribution" json:"total_contribution"`
	AvgDailyActivity float64 `gorm:"column:avg_daily_activity" json:"avg_daily_activity"`
}

// GetGuildHealth 获取宗门健康度排行
func GetGuildHealth(inactiveDays int) ([]GuildHealthItem, error) {
	g := db.GetGame()
	if g == nil {
		return nil, ErrGameDBNil
	}
	if inactiveDays <= 0 {
		inactiveDays = 7
	}

	var items []GuildHealthItem
	err := g.Table("guilds g").
		Select(`g.id as guild_id,
			g.name as guild_name,
			gl.level,
			g.member_count,
			gl.member_capacity,
			g.prestige,
			SUM(CASE WHEN gm.last_active_at >= DATE_SUB(NOW(), INTERVAL ? DAY) THEN 1 ELSE 0 END) as active_members,
			AVG(gm.contribution_total) as avg_contribution,
			SUM(gm.contribution_total) as total_contribution,
			AVG(gm.daily_activity_sum) as avg_daily_activity`, inactiveDays).
		Joins("LEFT JOIN guild_members gm ON gm.guild_id = g.id AND gm.deleted_at IS NULL").
		Joins("LEFT JOIN guild_level gl ON gl.id = g.level_id").
		Group("g.id, g.name, gl.level, g.member_count, gl.member_capacity, g.prestige").
		Order("g.prestige DESC").
		Find(&items).Error
	if err != nil {
		return nil, err
	}
	for i := range items {
		if items[i].MemberCount > 0 {
			items[i].ActiveRate = float64(items[i].ActiveMembers) / float64(items[i].MemberCount) * 100
		}
	}
	return items, nil
}

// GuildOverviewStats 宗门整体概览
type GuildOverviewStats struct {
	TotalGuilds    int     `gorm:"column:total_guilds"    json:"total_guilds"`
	TotalMembers   int     `gorm:"column:total_members"   json:"total_members"`
	AvgMemberCount float64 `gorm:"column:avg_member_count" json:"avg_member_count"`
	MaxMembers     int     `gorm:"column:max_members"      json:"max_members"`
	AvgPrestige    float64 `gorm:"column:avg_prestige"    json:"avg_prestige"`
	CoverageRate   float64 `json:"coverage_rate"`
}

// GetGuildOverview 宗门整体概览
func GetGuildOverview() (*GuildOverviewStats, error) {
	g := db.GetGame()
	if g == nil {
		return nil, ErrGameDBNil
	}
	var stats GuildOverviewStats
	err := g.Table("guilds").
		Select(`COUNT(*) as total_guilds,
			SUM(member_count) as total_members,
			AVG(member_count) as avg_member_count,
			MAX(member_count) as max_members,
			AVG(prestige) as avg_prestige`).
		Scan(&stats).Error
	if err != nil {
		return nil, err
	}
	// 计算覆盖率：加入宗门的用户占总用户比例
	type cntRow struct {
		Cnt int `gorm:"column:cnt"`
	}
	var total cntRow
	if err := g.Table("users").Select("COUNT(*) as cnt").Scan(&total).Error; err != nil {
		return nil, err
	}
	if total.Cnt > 0 {
		stats.CoverageRate = float64(stats.TotalMembers) / float64(total.Cnt) * 100
	}
	return &stats, nil
}

// ========== 社交网络分析 ==========

// SocialOverviewStats 社交整体概览
type SocialOverviewStats struct {
	TotalFriendships int     `gorm:"column:total_friendships" json:"total_friendships"`
	AvgFriends       float64 `gorm:"column:avg_friends"       json:"avg_friends"`
	MaxFriends       int     `gorm:"column:max_friends"       json:"max_friends"`
	CompanionPairs   int     `gorm:"column:companion_pairs"   json:"companion_pairs"`
	UsersWithFriends int     `gorm:"column:users_with_friends" json:"users_with_friends"`
	FriendCoverage   float64 `json:"friend_coverage"`
}

// GetSocialOverview 社交整体概览
func GetSocialOverview() (*SocialOverviewStats, error) {
	g := db.GetGame()
	if g == nil {
		return nil, ErrGameDBNil
	}
	var stats SocialOverviewStats
	if err := g.Table("friendships").
		Select(`FLOOR(COUNT(*) / 2) as total_friendships,
			COUNT(DISTINCT user_id) as users_with_friends,
			FLOOR(SUM(CASE WHEN relation_type = 2 THEN 1 ELSE 0 END) / 2) as companion_pairs`).
		Where("deleted_at IS NULL").
		Scan(&stats).Error; err != nil {
		return nil, err
	}

	// 人均好友数
	type avgRow struct {
		Avg float64 `gorm:"column:avg_val"`
		Max int     `gorm:"column:max_val"`
	}
	var ar avgRow
	if err := g.Raw(`SELECT AVG(cnt) as avg_val, MAX(cnt) as max_val 
		FROM (SELECT user_id, COUNT(*) as cnt FROM friendships WHERE deleted_at IS NULL GROUP BY user_id) sub`).
		Scan(&ar).Error; err != nil {
		return nil, err
	}
	stats.AvgFriends = ar.Avg
	stats.MaxFriends = ar.Max

	// 好友覆盖率
	type cntRow struct {
		Cnt int `gorm:"column:cnt"`
	}
	var total cntRow
	if err := g.Table("users").Select("COUNT(*) as cnt").Scan(&total).Error; err != nil {
		return nil, err
	}
	if total.Cnt > 0 {
		stats.FriendCoverage = float64(stats.UsersWithFriends) / float64(total.Cnt) * 100
	}
	return &stats, nil
}

// RelationTypeDistItem 关系类型分布
type RelationTypeDistItem struct {
	RelationType int    `gorm:"column:relation_type" json:"relation_type"`
	TypeName     string `json:"type_name"`
	Count        int    `gorm:"column:count"         json:"count"`
}

var relationTypeNames = map[int]string{
	1: "道友",
	2: "道侣",
	3: "师父",
	4: "徒弟",
}

// GetRelationTypeDist 关系类型分布
func GetRelationTypeDist() ([]RelationTypeDistItem, error) {
	g := db.GetGame()
	if g == nil {
		return nil, ErrGameDBNil
	}
	var items []RelationTypeDistItem
	err := g.Table("friendships").
		Select("relation_type, COUNT(*) as count").
		Where("deleted_at IS NULL").
		Group("relation_type").
		Order("relation_type ASC").
		Find(&items).Error
	if err != nil {
		return nil, err
	}
	for i := range items {
		if name, ok := relationTypeNames[items[i].RelationType]; ok {
			items[i].TypeName = name
		}
	}
	return items, nil
}

// IntimacyDistItem 亲密度分布
type IntimacyDistItem struct {
	Bracket string `json:"bracket"`
	Count   int    `gorm:"column:count"     json:"count"`
}

// GetIntimacyDist 亲密度分布
func GetIntimacyDist() ([]IntimacyDistItem, error) {
	g := db.GetGame()
	if g == nil {
		return nil, ErrGameDBNil
	}
	var items []IntimacyDistItem
	err := g.Raw(`SELECT 
		CASE
			WHEN intimacy < 10 THEN '0-10'
			WHEN intimacy < 50 THEN '10-50'
			WHEN intimacy < 100 THEN '50-100'
			WHEN intimacy < 200 THEN '100-200'
			ELSE '200+'
		END as bracket,
		COUNT(*) as count
		FROM friendships
		WHERE deleted_at IS NULL
		GROUP BY bracket
		ORDER BY FIELD(bracket, '0-10', '10-50', '50-100', '100-200', '200+')`).
		Scan(&items).Error
	return items, err
}

// ========== 签到分析 ==========

// CheckInStatsItem 签到统计项
type CheckInStatsItem struct {
	Date             string `gorm:"column:date"               json:"date"`
	CheckInCount     int    `gorm:"column:checkin_count"       json:"checkin_count"`
	MonthCardClaimed int    `gorm:"column:month_card_claimed"  json:"month_card_claimed"`
}

// GetCheckInTrend 获取签到趋势（最近N天）
func GetCheckInTrend(days int) ([]CheckInStatsItem, error) {
	g := db.GetGame()
	if g == nil {
		return nil, ErrGameDBNil
	}
	if days <= 0 {
		days = 30
	}
	var items []CheckInStatsItem
	err := g.Table("check_in_records").
		Select(`check_in_date as date,
			SUM(base_checked) as checkin_count,
			SUM(month_card_reward_claimed) as month_card_claimed`).
		Where("check_in_date >= DATE_SUB(CURDATE(), INTERVAL ? DAY)", days).
		Group("check_in_date").
		Order("check_in_date ASC").
		Find(&items).Error
	if err != nil {
		return nil, err
	}
	if len(items) > 0 {
		return items, nil
	}

	// 兜底: 部分环境未持续写入 check_in_records，回退到 users.last_check_in_time 聚合
	err = g.Table("users").
		Select(`DATE(last_check_in_time) as date,
			COUNT(*) as checkin_count,
			SUM(CASE WHEN month_card_expire_time IS NOT NULL AND month_card_expire_time >= last_check_in_time THEN 1 ELSE 0 END) as month_card_claimed`).
		Where("deleted_at IS NULL AND status = 1 AND last_check_in_time IS NOT NULL AND DATE(last_check_in_time) >= DATE_SUB(CURDATE(), INTERVAL ? DAY)", days).
		Group("DATE(last_check_in_time)").
		Order("date ASC").
		Find(&items).Error
	return items, err
}

// CheckInRateItem 签到率统计
type CheckInRateItem struct {
	TotalUsers    int     `gorm:"column:total_users"   json:"total_users"`
	CheckedIn     int     `gorm:"column:checked_in"    json:"checked_in"`
	MonthCardUse  int     `gorm:"column:month_card_use" json:"month_card_use"`
	CheckInRate   float64 `json:"checkin_rate"`
	MonthCardRate float64 `json:"month_card_rate"`
}

// GetCheckInRate 获取今日签到率
func GetCheckInRate() (*CheckInRateItem, error) {
	g := db.GetGame()
	if g == nil {
		return nil, ErrGameDBNil
	}
	var item CheckInRateItem
	// 总用户数
	if err := g.Table("users").Select("COUNT(*) as total_users").
		Where("deleted_at IS NULL AND status = 1").Scan(&item).Error; err != nil {
		return nil, err
	}
	// 今日签到数
	type cntRow struct {
		Cnt int `gorm:"column:cnt"`
	}
	var checkedIn cntRow
	if err := g.Table("check_in_records").
		Select("COUNT(*) as cnt").
		Where("check_in_date = CURDATE() AND base_checked = 1").
		Scan(&checkedIn).Error; err != nil {
		return nil, err
	}
	item.CheckedIn = checkedIn.Cnt
	var monthCard cntRow
	if err := g.Table("check_in_records").
		Select("COUNT(*) as cnt").
		Where("check_in_date = CURDATE() AND month_card_reward_claimed = 1").
		Scan(&monthCard).Error; err != nil {
		return nil, err
	}
	item.MonthCardUse = monthCard.Cnt

	if item.CheckedIn == 0 {
		if err := g.Table("users").
			Select("COUNT(*) as cnt").
			Where("deleted_at IS NULL AND status = 1 AND last_check_in_time IS NOT NULL AND DATE(last_check_in_time) = CURDATE()").
			Scan(&checkedIn).Error; err != nil {
			return nil, err
		}
		item.CheckedIn = checkedIn.Cnt
	}

	if item.MonthCardUse == 0 {
		if err := g.Table("users").
			Select("COUNT(*) as cnt").
			Where("deleted_at IS NULL AND status = 1 AND last_check_in_time IS NOT NULL AND DATE(last_check_in_time) = CURDATE() AND month_card_expire_time IS NOT NULL AND month_card_expire_time >= last_check_in_time").
			Scan(&monthCard).Error; err != nil {
			return nil, err
		}
		item.MonthCardUse = monthCard.Cnt
	}

	if item.TotalUsers > 0 {
		item.CheckInRate = float64(item.CheckedIn) / float64(item.TotalUsers) * 100
		item.MonthCardRate = float64(item.MonthCardUse) / float64(item.TotalUsers) * 100
	}
	return &item, nil
}

// ========== 闭关分析 ==========

// SeclusionStatsItem 闭关统计项
type SeclusionStatsItem struct {
	TotalSessions    int     `gorm:"column:total_sessions"   json:"total_sessions"`
	UniqueUsers      int     `gorm:"column:unique_users"     json:"unique_users"`
	PremiumSessions  int     `gorm:"column:premium_sessions" json:"premium_sessions"`
	CompletedCount   int     `gorm:"column:completed_count"  json:"completed_count"`
	InterruptedCount int     `gorm:"column:interrupted_count" json:"interrupted_count"`
	AvgDuration      float64 `gorm:"column:avg_duration"     json:"avg_duration"`
	PremiumRate      float64 `json:"premium_rate"`
	CompletionRate   float64 `json:"completion_rate"`
}

// GetSeclusionStats 获取闭关统计（最近N天）
func GetSeclusionStats(days int) (*SeclusionStatsItem, error) {
	g := db.GetGame()
	if g == nil {
		return nil, ErrGameDBNil
	}
	if days <= 0 {
		days = 30
	}
	var item SeclusionStatsItem
	err := g.Table("seclusions").
		Select(`COUNT(*) as total_sessions,
			COUNT(DISTINCT user_id) as unique_users,
			SUM(CASE WHEN is_premium = 1 THEN 1 ELSE 0 END) as premium_sessions,
			SUM(CASE WHEN status = 'completed' THEN 1 ELSE 0 END) as completed_count,
			SUM(CASE WHEN status = 'interrupted' THEN 1 ELSE 0 END) as interrupted_count,
			AVG(duration) as avg_duration`).
		Where("start_time >= DATE_SUB(NOW(), INTERVAL ? DAY)", days).
		Scan(&item).Error
	if err != nil {
		return nil, err
	}
	if item.TotalSessions > 0 {
		item.PremiumRate = float64(item.PremiumSessions) / float64(item.TotalSessions) * 100
		item.CompletionRate = float64(item.CompletedCount) / float64(item.TotalSessions) * 100
	}
	return &item, nil
}

// SeclusionDurationDistItem 闭关时长分布
type SeclusionDurationDistItem struct {
	Duration int `gorm:"column:duration" json:"duration"`
	Count    int `gorm:"column:count"    json:"count"`
}

// GetSeclusionDurationDist 闭关时长分布
func GetSeclusionDurationDist(days int) ([]SeclusionDurationDistItem, error) {
	g := db.GetGame()
	if g == nil {
		return nil, ErrGameDBNil
	}
	if days <= 0 {
		days = 30
	}
	var items []SeclusionDurationDistItem
	err := g.Table("seclusions").
		Select("duration, COUNT(*) as count").
		Where("start_time >= DATE_SUB(NOW(), INTERVAL ? DAY)", days).
		Group("duration").
		Order("duration ASC").
		Find(&items).Error
	return items, err
}

// ========== 交易市场分析 ==========

// MarketTradeStatsItem 交易市场统计项
type MarketTradeStatsItem struct {
	TotalTrades   int     `gorm:"column:total_trades"   json:"total_trades"`
	UniqueTraders int     `gorm:"column:unique_traders"  json:"unique_traders"`
	TotalVolume   int64   `gorm:"column:total_volume"    json:"total_volume"`
	TotalFees     int64   `gorm:"column:total_fees"      json:"total_fees"`
	AvgPrice      float64 `gorm:"column:avg_price"       json:"avg_price"`
}

// GetMarketTradeStats 交易市场统计（最近N天）
func GetMarketTradeStats(days int) (*MarketTradeStatsItem, error) {
	g := db.GetGame()
	if g == nil {
		return nil, ErrGameDBNil
	}
	if days <= 0 {
		days = 30
	}
	var item MarketTradeStatsItem
	err := g.Table("market_trade_records").
		Select(`COUNT(*) as total_trades,
			COUNT(DISTINCT buyer_user_id) + COUNT(DISTINCT seller_user_id) as unique_traders,
			COALESCE(SUM(total_price_amount), 0) as total_volume,
			COALESCE(SUM(fee_amount), 0) as total_fees,
			AVG(unit_price_amount) as avg_price`).
		Where("created_at >= DATE_SUB(NOW(), INTERVAL ? DAY)", days).
		Scan(&item).Error
	return &item, err
}

// MarketHotItemRow 热门商品
type MarketHotItemRow struct {
	ItemCategory string `gorm:"column:item_category" json:"item_category"`
	ItemName     string `gorm:"column:item_name"     json:"item_name"`
	TradeCount   int    `gorm:"column:trade_count"   json:"trade_count"`
	TotalVolume  int64  `gorm:"column:total_volume"  json:"total_volume"`
	BuyerCount   int    `gorm:"column:buyer_count"   json:"buyer_count"`
}

// GetMarketHotItems 热门交易物品 Top N
func GetMarketHotItems(days, limit int) ([]MarketHotItemRow, error) {
	g := db.GetGame()
	if g == nil {
		return nil, ErrGameDBNil
	}
	if days <= 0 {
		days = 30
	}
	if limit <= 0 {
		limit = 20
	}
	var items []MarketHotItemRow
	err := g.Table("market_trade_records").
		Select(`item_category, item_name,
			COUNT(*) as trade_count,
			SUM(total_price_amount) as total_volume,
			COUNT(DISTINCT buyer_user_id) as buyer_count`).
		Where("created_at >= DATE_SUB(NOW(), INTERVAL ? DAY)", days).
		Group("item_category, item_name").
		Order("trade_count DESC").
		Limit(limit).
		Find(&items).Error
	return items, err
}

// MarketDailyTrendItem 市场日趋势
type MarketDailyTrendItem struct {
	Date        string `gorm:"column:date"         json:"date"`
	TradeCount  int    `gorm:"column:trade_count"  json:"trade_count"`
	TotalVolume int64  `gorm:"column:total_volume" json:"total_volume"`
	UserCount   int    `gorm:"column:user_count"   json:"user_count"`
}

// GetMarketDailyTrend 市场日趋势
func GetMarketDailyTrend(days int) ([]MarketDailyTrendItem, error) {
	g := db.GetGame()
	if g == nil {
		return nil, ErrGameDBNil
	}
	if days <= 0 {
		days = 30
	}
	var items []MarketDailyTrendItem
	err := g.Table("market_trade_records").
		Select(`DATE(created_at) as date,
			COUNT(*) as trade_count,
			COALESCE(SUM(total_price_amount), 0) as total_volume,
			COUNT(DISTINCT buyer_user_id) as user_count`).
		Where("created_at >= DATE_SUB(NOW(), INTERVAL ? DAY)", days).
		Group("DATE(created_at)").
		Order("date ASC").
		Find(&items).Error
	return items, err
}

// ========== 主线剧情分析 ==========

// StoryProgressItem 主线进度统计
type StoryProgressItem struct {
	ChapterNumber int    `gorm:"column:chapter_number" json:"chapter_number"`
	ChapterName   string `gorm:"column:chapter_name"   json:"chapter_name"`
	Arc           string `gorm:"column:arc"             json:"arc"`
	ClearedUsers  int    `gorm:"column:cleared_users"   json:"cleared_users"`
	TotalClears   int    `gorm:"column:total_clears"    json:"total_clears"`
}

// GetStoryProgress 主线章节通过统计
func GetStoryProgress() ([]StoryProgressItem, error) {
	g := db.GetGame()
	if g == nil {
		return nil, ErrGameDBNil
	}
	var items []StoryProgressItem
	err := g.Table("stage_records sr").
		Select(`c.chapter_number, c.name as chapter_name, c.arc,
			COUNT(DISTINCT sr.user_id) as cleared_users,
			COUNT(*) as total_clears`).
		Joins("INNER JOIN stages s ON s.id = sr.stage_id").
		Joins("INNER JOIN chapters c ON c.id = s.chapter_id").
		Group("c.chapter_number, c.name, c.arc").
		Order("c.chapter_number ASC").
		Find(&items).Error
	return items, err
}

// StoryFunnelItem 主线漏斗项
type StoryFunnelItem struct {
	ChapterNumber int    `gorm:"column:chapter_number" json:"chapter_number"`
	ChapterName   string `gorm:"column:chapter_name"   json:"chapter_name"`
	ReachedUsers  int    `gorm:"column:reached_users"   json:"reached_users"`
}

// GetStoryFunnel 主线进度漏斗（到达每章的用户数）
func GetStoryFunnel() ([]StoryFunnelItem, error) {
	g := db.GetGame()
	if g == nil {
		return nil, ErrGameDBNil
	}
	var items []StoryFunnelItem
	err := g.Raw(`SELECT c.chapter_number, c.name as chapter_name,
		COUNT(DISTINCT sr.user_id) as reached_users
		FROM chapters c
		INNER JOIN stages s ON s.chapter_id = c.id
		INNER JOIN stage_records sr ON sr.stage_id = s.id
		GROUP BY c.chapter_number, c.name
		ORDER BY c.chapter_number ASC`).
		Scan(&items).Error
	return items, err
}

// StageDifficultyItem 关卡难度分析
type StageDifficultyItem struct {
	ChapterNumber int     `gorm:"column:chapter_number" json:"chapter_number"`
	StageNumber   int     `gorm:"column:stage_number"   json:"stage_number"`
	StageLabel    string  `gorm:"column:stage_label"     json:"stage_label"`
	IsElite       bool    `gorm:"column:is_elite"        json:"is_elite"`
	ClearCount    int     `gorm:"column:clear_count"     json:"clear_count"`
	UniqueUsers   int     `gorm:"column:unique_users"    json:"unique_users"`
	AvgStars      float64 `gorm:"column:avg_stars"        json:"avg_stars"`
}

// GetStageDifficulty 关卡难度分析（通关次数、人数、平均星级）
func GetStageDifficulty(chapterNumber int) ([]StageDifficultyItem, error) {
	g := db.GetGame()
	if g == nil {
		return nil, ErrGameDBNil
	}
	var items []StageDifficultyItem
	q := g.Table("stage_records sr").
		Select(`c.chapter_number, s.stage_number, s.stage_label, s.is_elite,
			COUNT(*) as clear_count,
			COUNT(DISTINCT sr.user_id) as unique_users,
			AVG(sr.stars) as avg_stars`).
		Joins("INNER JOIN stages s ON s.id = sr.stage_id").
		Joins("INNER JOIN chapters c ON c.id = s.chapter_id").
		Group("c.chapter_number, s.stage_number, s.stage_label, s.is_elite").
		Order("c.chapter_number ASC, s.stage_number ASC")
	if chapterNumber > 0 {
		q = q.Where("c.chapter_number = ?", chapterNumber)
	}
	err := q.Find(&items).Error
	return items, err
}

// ========== 妖塔分析 ==========

// TowerStatsItem 妖塔统计项
type TowerStatsItem struct {
	TotalRuns      int     `gorm:"column:total_runs"       json:"total_runs"`
	UniqueUsers    int     `gorm:"column:unique_users"     json:"unique_users"`
	CompletedRuns  int     `gorm:"column:completed_runs"   json:"completed_runs"`
	FailedRuns     int     `gorm:"column:failed_runs"      json:"failed_runs"`
	AvgMaxLayer    float64 `gorm:"column:avg_max_layer"    json:"avg_max_layer"`
	AvgScore       float64 `gorm:"column:avg_score"        json:"avg_score"`
	CompletionRate float64 `json:"completion_rate"`
}

// GetTowerStats 妖塔统计（最近N天）
func GetTowerStats(days int) (*TowerStatsItem, error) {
	g := db.GetGame()
	if g == nil {
		return nil, ErrGameDBNil
	}
	if days <= 0 {
		days = 30
	}
	var item TowerStatsItem
	err := g.Table("demon_tower_runs").
		Select(`COUNT(*) as total_runs,
			COUNT(DISTINCT user_id) as unique_users,
			SUM(CASE WHEN is_completed = 1 THEN 1 ELSE 0 END) as completed_runs,
			SUM(CASE WHEN is_completed = 2 THEN 1 ELSE 0 END) as failed_runs,
			AVG(current_layer) as avg_max_layer,
			AVG(total_score) as avg_score`).
		Where("run_date >= DATE_SUB(CURDATE(), INTERVAL ? DAY)", days).
		Scan(&item).Error
	if err != nil {
		return nil, err
	}
	if item.TotalRuns > 0 {
		item.CompletionRate = float64(item.CompletedRuns) / float64(item.TotalRuns) * 100
	}
	return &item, nil
}

// TowerLayerDistItem 妖塔层级分布
type TowerLayerDistItem struct {
	CurrentLayer int `gorm:"column:current_layer" json:"current_layer"`
	UserCount    int `gorm:"column:user_count"    json:"user_count"`
}

// GetTowerLayerDist 妖塔最高到达层级分布
func GetTowerLayerDist(days int) ([]TowerLayerDistItem, error) {
	g := db.GetGame()
	if g == nil {
		return nil, ErrGameDBNil
	}
	if days <= 0 {
		days = 30
	}
	var items []TowerLayerDistItem
	err := g.Raw(`SELECT current_layer, COUNT(DISTINCT user_id) as user_count
		FROM demon_tower_runs
		WHERE run_date >= DATE_SUB(CURDATE(), INTERVAL ? DAY)
		GROUP BY current_layer
		ORDER BY current_layer ASC`, days).
		Scan(&items).Error
	return items, err
}

// ========== 月卡/闭关卡付费分析 ==========

// PassSubscriptionStats 卡类付费统计
type PassSubscriptionStats struct {
	ActiveMonthCards      int `gorm:"column:active_month_cards"     json:"active_month_cards"`
	ActiveSeclusionPasses int `gorm:"column:active_seclusion_passes" json:"active_seclusion_passes"`
	TotalMonthCardUsers   int `gorm:"column:total_month_card_users" json:"total_month_card_users"`
	TotalSeclusionUsers   int `gorm:"column:total_seclusion_users"  json:"total_seclusion_users"`
}

// GetPassSubscriptionStats 获取月卡/闭关卡付费用户统计
func GetPassSubscriptionStats() (*PassSubscriptionStats, error) {
	g := db.GetGame()
	if g == nil {
		return nil, ErrGameDBNil
	}
	var stats PassSubscriptionStats
	err := g.Table("users").
		Select(`SUM(CASE WHEN month_card_expire_time > NOW() THEN 1 ELSE 0 END) as active_month_cards,
			SUM(CASE WHEN seclusion_pass_expire_time > NOW() THEN 1 ELSE 0 END) as active_seclusion_passes,
			SUM(CASE WHEN month_card_total_days > 0 THEN 1 ELSE 0 END) as total_month_card_users,
			SUM(CASE WHEN seclusion_pass_total_days > 0 THEN 1 ELSE 0 END) as total_seclusion_users`).
		Where("deleted_at IS NULL AND status = 1").
		Scan(&stats).Error
	return &stats, err
}

// ========== 装备 & 战力分析 ==========

// EquipGradeDistItem 装备品级分布
type EquipGradeDistItem struct {
	GradeName string `gorm:"column:grade_name" json:"grade_name"`
	Count     int    `gorm:"column:count"      json:"count"`
}

// GetEquipGradeDist 已装备物品品级分布
func GetEquipGradeDist() ([]EquipGradeDistItem, error) {
	g := db.GetGame()
	if g == nil {
		return nil, ErrGameDBNil
	}
	var items []EquipGradeDistItem
	err := g.Raw(`SELECT g.name as grade_name, COUNT(*) as count
		FROM user_equipments ue
		INNER JOIN equipments e ON e.id = ue.equipment_id
		INNER JOIN grades g ON g.id = e.grade_id
		WHERE ue.is_equipped = 1
		GROUP BY g.name, g.id
		ORDER BY g.id ASC`).
		Scan(&items).Error
	return items, err
}

// EquipSetPopItem 套装使用排行
type EquipSetPopItem struct {
	SetName   string  `gorm:"column:set_name"   json:"set_name"`
	GradeName string  `gorm:"column:grade_name" json:"grade_name"`
	UserCount int     `gorm:"column:user_count" json:"user_count"`
	PieceAvg  float64 `gorm:"column:piece_avg" json:"piece_avg"`
}

// GetEquipSetPopularity 套装使用热度排行
func GetEquipSetPopularity(limit int) ([]EquipSetPopItem, error) {
	g := db.GetGame()
	if g == nil {
		return nil, ErrGameDBNil
	}
	if limit <= 0 {
		limit = 20
	}
	var items []EquipSetPopItem
	err := g.Raw(`SELECT es.name as set_name, g.name as grade_name,
		COUNT(DISTINCT ue.user_id) as user_count,
		COUNT(*) * 1.0 / COUNT(DISTINCT ue.user_id) as piece_avg
		FROM user_equipments ue
		INNER JOIN equipments e ON e.id = ue.equipment_id
		INNER JOIN equipment_sets es ON es.id = e.set_id
		INNER JOIN grades g ON g.id = e.grade_id
		WHERE ue.is_equipped = 1 AND e.set_id IS NOT NULL
		GROUP BY es.name, g.name, g.id
		ORDER BY user_count DESC
		LIMIT ?`, limit).
		Scan(&items).Error
	return items, err
}

// EquipRealmDistItem 装备品级按境界分布
type EquipRealmDistItem struct {
	StageName string  `gorm:"column:stage_name" json:"stage_name"`
	GradeName string  `gorm:"column:grade_name" json:"grade_name"`
	AvgCount  float64 `gorm:"column:avg_count"  json:"avg_count"`
}

// GetEquipRealmDist 各境界玩家装备品级分布
func GetEquipRealmDist() ([]EquipRealmDistItem, error) {
	g := db.GetGame()
	if g == nil {
		return nil, ErrGameDBNil
	}
	var items []EquipRealmDistItem
	err := g.Raw(`SELECT
		CASE r.stage_level
			WHEN 0 THEN '凡人' WHEN 1 THEN '练气' WHEN 2 THEN '筑基' WHEN 3 THEN '金丹'
			WHEN 4 THEN '元婴' WHEN 5 THEN '化神' WHEN 6 THEN '洞虚' WHEN 7 THEN '大乘'
			WHEN 8 THEN '仙' WHEN 9 THEN '金仙' WHEN 10 THEN '道祖'
		END as stage_name,
		g.name as grade_name,
		COUNT(*) * 1.0 / COUNT(DISTINCT u.id) as avg_count
		FROM users u
		INNER JOIN realms r ON r.id = u.realm_id
		INNER JOIN user_equipments ue ON ue.user_id = u.id AND ue.is_equipped = 1
		INNER JOIN equipments e ON e.id = ue.equipment_id
		INNER JOIN grades g ON g.id = e.grade_id
		GROUP BY r.stage_level, g.name, g.id
		ORDER BY r.stage_level ASC, g.id ASC`).
		Scan(&items).Error
	return items, err
}

// RechargeConversionItem 首充转化分析
type RechargeConversionItem struct {
	DayBucket string `gorm:"column:day_bucket" json:"day_bucket"`
	UserCount int    `gorm:"column:user_count" json:"user_count"`
}

// GetFirstRechargeConversion 注册 → 首充天数分布
func GetFirstRechargeConversion() ([]RechargeConversionItem, error) {
	g := db.GetGame()
	if g == nil {
		return nil, ErrGameDBNil
	}
	var items []RechargeConversionItem
	err := g.Raw(`SELECT
		CASE
			WHEN days <= 0 THEN '当天'
			WHEN days <= 1 THEN '1天内'
			WHEN days <= 3 THEN '2-3天'
			WHEN days <= 7 THEN '4-7天'
			WHEN days <= 14 THEN '8-14天'
			WHEN days <= 30 THEN '15-30天'
			ELSE '30天以上'
		END as day_bucket,
		COUNT(*) as user_count
		FROM (
			SELECT DATEDIFF(MIN(ro.pay_time), u.created_at) as days
			FROM users u
			INNER JOIN recharge_orders ro ON ro.uid = u.id AND ` + realRechargeOrderWhereAlias + `
			GROUP BY u.id
		) t
		GROUP BY day_bucket
		ORDER BY MIN(days) ASC`).
		Scan(&items).Error
	return items, err
}

// RechargeTierTrendItem 各档位充值趋势
type RechargeTierTrendItem struct {
	Date        string `gorm:"column:date"         json:"date"`
	PackageName string `gorm:"column:package_name" json:"package_name"`
	OrderCount  int    `gorm:"column:order_count"  json:"order_count"`
	Revenue     int64  `gorm:"column:revenue"      json:"revenue"`
}

// GetRechargeTierTrend 各档位日充值趋势
func GetRechargeTierTrend(days int) ([]RechargeTierTrendItem, error) {
	g := db.GetGame()
	if g == nil {
		return nil, ErrGameDBNil
	}
	if days <= 0 {
		days = 30
	}
	var items []RechargeTierTrendItem
	err := g.Table("recharge_orders ro").
		Select(`DATE(ro.pay_time) as date,
			rp.name as package_name,
			COUNT(*) as order_count,
			COALESCE(SUM(ro.pay_amount), 0) as revenue`).
		Joins("INNER JOIN recharge_packages rp ON rp.id = ro.package_id").
		Where(realRechargeOrderWhereAlias+" AND ro.pay_time >= DATE_SUB(NOW(), INTERVAL ? DAY)", days).
		Group("DATE(ro.pay_time), rp.name").
		Order("date ASC").
		Find(&items).Error
	return items, err
}

// MarketPriceAnomalyItem 价格异常检测
type MarketPriceAnomalyItem struct {
	ItemCategory string  `gorm:"column:item_category" json:"item_category"`
	ItemName     string  `gorm:"column:item_name"     json:"item_name"`
	BasePrice    float64 `gorm:"column:base_price"    json:"base_price"`
	AvgPrice     float64 `gorm:"column:avg_price"     json:"avg_price"`
	MaxPrice     int64   `gorm:"column:max_price"     json:"max_price"`
	MinPrice     int64   `gorm:"column:min_price"     json:"min_price"`
	Deviation    float64 `json:"deviation"`
	TradeCount   int     `gorm:"column:trade_count"   json:"trade_count"`
}

// GetMarketPriceAnomalies 交易价格异常检测
func GetMarketPriceAnomalies(days int) ([]MarketPriceAnomalyItem, error) {
	g := db.GetGame()
	if g == nil {
		return nil, ErrGameDBNil
	}
	if days <= 0 {
		days = 7
	}
	var items []MarketPriceAnomalyItem
	err := g.Raw(`SELECT cur.item_category, cur.item_name,
		COALESCE(hist.base_price, cur.avg_price) as base_price,
		cur.avg_price,
		cur.max_price,
		cur.min_price,
		cur.trade_count
		FROM (
			SELECT item_category, item_name,
				AVG(unit_price_amount) as avg_price,
				MAX(unit_price_amount) as max_price,
				MIN(unit_price_amount) as min_price,
				COUNT(*) as trade_count
			FROM market_trade_records
			WHERE created_at >= DATE_SUB(NOW(), INTERVAL ? DAY)
			GROUP BY item_category, item_name
			HAVING COUNT(*) >= 3
		) cur
		LEFT JOIN (
			SELECT item_category, item_name,
				ROUND(AVG(unit_price_amount)) as base_price
			FROM market_trade_records
			WHERE created_at < DATE_SUB(NOW(), INTERVAL ? DAY)
				AND created_at >= DATE_SUB(NOW(), INTERVAL ? DAY)
			GROUP BY item_category, item_name
		) hist ON cur.item_category = hist.item_category AND cur.item_name = hist.item_name
		ORDER BY cur.trade_count DESC`, days, days, days*2).
		Scan(&items).Error
	if err != nil {
		return nil, err
	}
	for i := range items {
		base := items[i].BasePrice
		if base > 0 {
			items[i].Deviation = (items[i].AvgPrice - base) / base * 100
		} else if items[i].AvgPrice > 0 {
			spread := float64(items[i].MaxPrice-items[i].MinPrice) / items[i].AvgPrice * 100
			items[i].Deviation = spread
		}
	}
	return items, nil
}

// MarketLargeTradeItem 大额交易监控
type MarketLargeTradeItem struct {
	BuyerID      uint   `gorm:"column:buyer_id"       json:"buyer_id"`
	SellerID     uint   `gorm:"column:seller_id"      json:"seller_id"`
	ItemCategory string `gorm:"column:item_category"  json:"item_category"`
	ItemName     string `gorm:"column:item_name"      json:"item_name"`
	TotalPrice   int64  `gorm:"column:total_price"    json:"total_price"`
	Quantity     int    `gorm:"column:quantity"       json:"quantity"`
	TradeTime    string `gorm:"column:trade_time"     json:"trade_time"`
}

// GetMarketLargeTrades 大额交易监控
func GetMarketLargeTrades(days, limit int, minPrice int64) ([]MarketLargeTradeItem, error) {
	g := db.GetGame()
	if g == nil {
		return nil, ErrGameDBNil
	}
	if days <= 0 {
		days = 7
	}
	if limit <= 0 {
		limit = 50
	}
	if minPrice <= 0 {
		minPrice = 100000 // 10万灵石
	}
	var items []MarketLargeTradeItem
	err := g.Table("market_trade_records").
		Select(`buyer_user_id as buyer_id, seller_user_id as seller_id, item_category, item_name,
			total_price_amount as total_price, quantity,
			DATE_FORMAT(created_at, '%Y-%m-%d %H:%i') as trade_time`).
		Where("created_at >= DATE_SUB(NOW(), INTERVAL ? DAY) AND total_price_amount >= ?", days, minPrice).
		Order("total_price_amount DESC").
		Limit(limit).
		Find(&items).Error
	return items, err
}

// MarketPlayerPairAggregateItem 玩家交易对手聚合
type MarketPlayerPairAggregateItem struct {
	BuyerID       uint   `gorm:"column:buyer_id" json:"buyer_id"`
	SellerID      uint   `gorm:"column:seller_id" json:"seller_id"`
	TradeCount    int    `gorm:"column:trade_count" json:"trade_count"`
	TotalVolume   int64  `gorm:"column:total_volume" json:"total_volume"`
	TradeDayCount int    `gorm:"column:trade_day_count" json:"trade_day_count"`
	FirstTradeAt  string `gorm:"column:first_trade_at" json:"first_trade_at"`
	LastTradeAt   string `gorm:"column:last_trade_at" json:"last_trade_at"`
}

// MarketPlayerCounterpartyItem 玩家主要交易对手
type MarketPlayerCounterpartyItem struct {
	UserID        uint    `json:"user_id"`
	UserName      string  `json:"user_name"`
	TradeCount    int     `json:"trade_count"`
	TotalVolume   int64   `json:"total_volume"`
	TradeDayCount int     `json:"trade_day_count"`
	Concentration float64 `json:"concentration"`
}

// MarketPlayerProfileItem 交易玩家画像
type MarketPlayerProfileItem struct {
	UserID                  uint                          `json:"user_id"`
	UserName                string                        `json:"user_name"`
	RealmID                 int                           `json:"realm_id"`
	VipLevel                int                           `json:"vip_level"`
	TotalRechargeAmount     int64                         `json:"total_recharge_amount"`
	RealRechargeAmount      int64                         `json:"real_recharge_amount"`
	AccountAgeDays          int                           `json:"account_age_days"`
	BuyTradeCount           int                           `json:"buy_trade_count"`
	SellTradeCount          int                           `json:"sell_trade_count"`
	BuyVolume               int64                         `json:"buy_volume"`
	SellVolume              int64                         `json:"sell_volume"`
	TotalTradeCount         int                           `json:"total_trade_count"`
	TotalTradeVolume        int64                         `json:"total_trade_volume"`
	BuySellRatio            float64                       `json:"buy_sell_ratio"`
	SellBuyRatio            float64                       `json:"sell_buy_ratio"`
	TradeRole               string                        `json:"trade_role"`
	CounterpartyCount       int                           `json:"counterparty_count"`
	TopBuyer                *MarketPlayerCounterpartyItem `json:"top_buyer,omitempty"`
	TopSeller               *MarketPlayerCounterpartyItem `json:"top_seller,omitempty"`
	LowPriceToUserID        uint                          `json:"low_price_to_user_id"`
	LowPriceToUserName      string                        `json:"low_price_to_user_name"`
	LowPriceItemName        string                        `json:"low_price_item_name"`
	LowPriceDiscountRate    float64                       `json:"low_price_discount_rate"`
	LowPriceTradeCount      int                           `json:"low_price_trade_count"`
	SuspiciousScore         int                           `json:"suspicious_score"`
	SuspiciousLevel         string                        `json:"suspicious_level"`
	SuspiciousReasons       []string                      `json:"suspicious_reasons"`
	SuspiciousAltLikelihood float64                       `json:"suspicious_alt_likelihood"`
}

type marketSellerLowPriceLinkItem struct {
	SellerID       uint    `gorm:"column:seller_id"`
	BuyerID        uint    `gorm:"column:buyer_id"`
	ItemCategory   string  `gorm:"column:item_category"`
	ItemName       string  `gorm:"column:item_name"`
	TradeCount     int     `gorm:"column:trade_count"`
	AvgUnitPrice   float64 `gorm:"column:avg_unit_price"`
	MarketAvgPrice float64 `gorm:"column:market_avg_price"`
}

type marketPlayerMetaRow struct {
	ID                  uint      `gorm:"column:id"`
	Name                string    `gorm:"column:name"`
	RealmID             int       `gorm:"column:realm_id"`
	VipLevel            int       `gorm:"column:vip_level"`
	TotalRechargeAmount int64     `gorm:"column:total_recharge_amount"`
	CreatedAt           time.Time `gorm:"column:created_at"`
}

type marketPlayerAccumulator struct {
	profile        *MarketPlayerProfileItem
	counterparties map[uint]struct{}
	buyPartners    map[uint]*MarketPlayerCounterpartyItem
	sellPartners   map[uint]*MarketPlayerCounterpartyItem
}

type marketUserRechargeSumRow struct {
	UID               uint  `gorm:"column:uid"`
	RealRechargeTotal int64 `gorm:"column:real_recharge_total"`
}

// GetMarketPlayerProfiles 获取交易玩家画像及疑似小号分析
func GetMarketPlayerProfiles(days, limit int) ([]MarketPlayerProfileItem, error) {
	g := db.GetGame()
	if g == nil {
		return nil, ErrGameDBNil
	}
	if days <= 0 {
		days = 30
	}
	if limit <= 0 {
		limit = 100
	}

	var pairs []MarketPlayerPairAggregateItem
	err := g.Table("market_trade_records").
		Select(`buyer_user_id as buyer_id,
			seller_user_id as seller_id,
			COUNT(*) as trade_count,
			COALESCE(SUM(total_price_amount), 0) as total_volume,
			COUNT(DISTINCT DATE(created_at)) as trade_day_count,
			DATE_FORMAT(MIN(created_at), '%Y-%m-%d %H:%i') as first_trade_at,
			DATE_FORMAT(MAX(created_at), '%Y-%m-%d %H:%i') as last_trade_at`).
		Where("created_at >= DATE_SUB(NOW(), INTERVAL ? DAY)", days).
		Group("buyer_user_id, seller_user_id").
		Find(&pairs).Error
	if err != nil {
		return nil, err
	}
	if len(pairs) == 0 {
		return []MarketPlayerProfileItem{}, nil
	}

	userIDs := make([]uint, 0, len(pairs)*2)
	userIDSet := make(map[uint]struct{}, len(pairs)*2)
	for _, pair := range pairs {
		if _, ok := userIDSet[pair.BuyerID]; !ok {
			userIDSet[pair.BuyerID] = struct{}{}
			userIDs = append(userIDs, pair.BuyerID)
		}
		if _, ok := userIDSet[pair.SellerID]; !ok {
			userIDSet[pair.SellerID] = struct{}{}
			userIDs = append(userIDs, pair.SellerID)
		}
	}

	var users []marketPlayerMetaRow
	err = g.Table("users").
		Select("id, name, realm_id, vip_level, total_recharge_amount, created_at").
		Where("id IN ?", userIDs).
		Find(&users).Error
	if err != nil {
		return nil, err
	}

	metaByID := make(map[uint]marketPlayerMetaRow, len(users))
	for _, user := range users {
		metaByID[user.ID] = user
	}

	var rechargeRows []marketUserRechargeSumRow
	err = g.Table("recharge_orders ro").
		Select("ro.uid as uid, COALESCE(SUM(ro.pay_amount), 0) as real_recharge_total").
		Where(realRechargeOrderWhereAlias+" AND ro.uid IN ?", userIDs).
		Group("ro.uid").
		Scan(&rechargeRows).Error
	if err != nil {
		return nil, err
	}
	realRechargeByUID := make(map[uint]int64, len(rechargeRows))
	for _, row := range rechargeRows {
		realRechargeByUID[row.UID] = row.RealRechargeTotal
	}

	var lowPriceLinks []marketSellerLowPriceLinkItem
	err = g.Raw(`SELECT
		s.seller_user_id as seller_id,
		s.buyer_user_id as buyer_id,
		s.item_category,
		s.item_name,
		s.trade_count,
		s.avg_unit_price,
		m.market_avg_price
	FROM (
		SELECT seller_user_id, buyer_user_id, item_category, item_name,
			COUNT(*) as trade_count,
			AVG(unit_price_amount) as avg_unit_price
		FROM market_trade_records
		WHERE created_at >= DATE_SUB(NOW(), INTERVAL ? DAY)
		GROUP BY seller_user_id, buyer_user_id, item_category, item_name
		HAVING COUNT(*) >= 3
	) s
	INNER JOIN (
		SELECT item_category, item_name, AVG(unit_price_amount) as market_avg_price
		FROM market_trade_records
		WHERE created_at >= DATE_SUB(NOW(), INTERVAL ? DAY)
		GROUP BY item_category, item_name
	) m ON s.item_category = m.item_category AND s.item_name = m.item_name
	WHERE m.market_avg_price > 0 AND s.avg_unit_price <= m.market_avg_price * 0.85`, days, days).
		Scan(&lowPriceLinks).Error
	if err != nil {
		return nil, err
	}

	lowPriceBySeller := make(map[uint]marketSellerLowPriceLinkItem)
	for _, link := range lowPriceLinks {
		discount := 0.0
		if link.MarketAvgPrice > 0 {
			discount = (link.MarketAvgPrice - link.AvgUnitPrice) / link.MarketAvgPrice * 100
		}
		best, ok := lowPriceBySeller[link.SellerID]
		if !ok {
			lowPriceBySeller[link.SellerID] = link
			continue
		}
		bestDiscount := 0.0
		if best.MarketAvgPrice > 0 {
			bestDiscount = (best.MarketAvgPrice - best.AvgUnitPrice) / best.MarketAvgPrice * 100
		}
		if link.TradeCount > best.TradeCount || (link.TradeCount == best.TradeCount && discount > bestDiscount) {
			lowPriceBySeller[link.SellerID] = link
		}
	}

	accByID := make(map[uint]*marketPlayerAccumulator, len(userIDs))
	ensureAccumulator := func(uid uint) *marketPlayerAccumulator {
		if acc, ok := accByID[uid]; ok {
			return acc
		}
		meta := metaByID[uid]
		profile := &MarketPlayerProfileItem{
			UserID:              uid,
			UserName:            meta.Name,
			RealmID:             meta.RealmID,
			VipLevel:            meta.VipLevel,
			TotalRechargeAmount: meta.TotalRechargeAmount,
			RealRechargeAmount:  realRechargeByUID[uid],
		}
		if !meta.CreatedAt.IsZero() {
			profile.AccountAgeDays = int(time.Since(meta.CreatedAt).Hours() / 24)
		}
		acc := &marketPlayerAccumulator{
			profile:        profile,
			counterparties: make(map[uint]struct{}),
			buyPartners:    make(map[uint]*MarketPlayerCounterpartyItem),
			sellPartners:   make(map[uint]*MarketPlayerCounterpartyItem),
		}
		accByID[uid] = acc
		return acc
	}

	for _, pair := range pairs {
		buyerAcc := ensureAccumulator(pair.BuyerID)
		sellerAcc := ensureAccumulator(pair.SellerID)

		buyerAcc.profile.BuyTradeCount += pair.TradeCount
		buyerAcc.profile.BuyVolume += pair.TotalVolume
		buyerAcc.counterparties[pair.SellerID] = struct{}{}
		buyerPartner, ok := buyerAcc.buyPartners[pair.SellerID]
		if !ok {
			buyerPartner = &MarketPlayerCounterpartyItem{
				UserID:   pair.SellerID,
				UserName: metaByID[pair.SellerID].Name,
			}
			buyerAcc.buyPartners[pair.SellerID] = buyerPartner
		}
		buyerPartner.TradeCount += pair.TradeCount
		buyerPartner.TotalVolume += pair.TotalVolume
		buyerPartner.TradeDayCount += pair.TradeDayCount

		sellerAcc.profile.SellTradeCount += pair.TradeCount
		sellerAcc.profile.SellVolume += pair.TotalVolume
		sellerAcc.counterparties[pair.BuyerID] = struct{}{}
		sellerPartner, ok := sellerAcc.sellPartners[pair.BuyerID]
		if !ok {
			sellerPartner = &MarketPlayerCounterpartyItem{
				UserID:   pair.BuyerID,
				UserName: metaByID[pair.BuyerID].Name,
			}
			sellerAcc.sellPartners[pair.BuyerID] = sellerPartner
		}
		sellerPartner.TradeCount += pair.TradeCount
		sellerPartner.TotalVolume += pair.TotalVolume
		sellerPartner.TradeDayCount += pair.TradeDayCount
	}

	profiles := make([]MarketPlayerProfileItem, 0, len(accByID))
	for _, acc := range accByID {
		profile := acc.profile
		profile.TotalTradeCount = profile.BuyTradeCount + profile.SellTradeCount
		profile.TotalTradeVolume = profile.BuyVolume + profile.SellVolume
		profile.CounterpartyCount = len(acc.counterparties)

		if profile.SellTradeCount == 0 {
			profile.BuySellRatio = float64(profile.BuyTradeCount)
		} else {
			profile.BuySellRatio = float64(profile.BuyTradeCount) / float64(profile.SellTradeCount)
		}
		if profile.BuyTradeCount == 0 {
			profile.SellBuyRatio = float64(profile.SellTradeCount)
		} else {
			profile.SellBuyRatio = float64(profile.SellTradeCount) / float64(profile.BuyTradeCount)
		}

		switch {
		case profile.SellTradeCount >= profile.BuyTradeCount*3 && profile.SellTradeCount >= 5:
			profile.TradeRole = "mostly_seller"
		case profile.BuyTradeCount >= profile.SellTradeCount*3 && profile.BuyTradeCount >= 5:
			profile.TradeRole = "mostly_buyer"
		default:
			profile.TradeRole = "balanced"
		}

		for _, partner := range acc.sellPartners {
			partner.Concentration = 0
			if profile.SellTradeCount > 0 {
				partner.Concentration = float64(partner.TradeCount) / float64(profile.SellTradeCount) * 100
			}
			if profile.TopBuyer == nil || partner.TradeCount > profile.TopBuyer.TradeCount ||
				(partner.TradeCount == profile.TopBuyer.TradeCount && partner.TotalVolume > profile.TopBuyer.TotalVolume) {
				candidate := *partner
				profile.TopBuyer = &candidate
			}
		}
		for _, partner := range acc.buyPartners {
			partner.Concentration = 0
			if profile.BuyTradeCount > 0 {
				partner.Concentration = float64(partner.TradeCount) / float64(profile.BuyTradeCount) * 100
			}
			if profile.TopSeller == nil || partner.TradeCount > profile.TopSeller.TradeCount ||
				(partner.TradeCount == profile.TopSeller.TradeCount && partner.TotalVolume > profile.TopSeller.TotalVolume) {
				candidate := *partner
				profile.TopSeller = &candidate
			}
		}

		score := 0
		reasons := make([]string, 0, 6)

		if profile.TradeRole == "mostly_seller" && profile.SellTradeCount >= 8 {
			score += 18
			reasons = append(reasons, "卖出占比极高，接近供货号")
		}
		if profile.TradeRole == "mostly_buyer" && profile.BuyTradeCount >= 8 {
			score += 16
			reasons = append(reasons, "买入占比极高，接近资源承接号")
		}
		if profile.CounterpartyCount <= 2 && profile.TotalTradeCount >= 8 {
			score += 18
			reasons = append(reasons, "交易对手过于集中")
		}
		if profile.TopBuyer != nil && profile.TopBuyer.Concentration >= 65 && profile.SellTradeCount >= 6 {
			score += 20
			reasons = append(reasons, "长期卖给固定对象")
		}
		if profile.TopSeller != nil && profile.TopSeller.Concentration >= 65 && profile.BuyTradeCount >= 6 {
			score += 18
			reasons = append(reasons, "长期从固定对象买入")
		}
		if lowPriceLink, ok := lowPriceBySeller[profile.UserID]; ok {
			discountRate := 0.0
			if lowPriceLink.MarketAvgPrice > 0 {
				discountRate = (lowPriceLink.MarketAvgPrice - lowPriceLink.AvgUnitPrice) / lowPriceLink.MarketAvgPrice * 100
			}
			profile.LowPriceToUserID = lowPriceLink.BuyerID
			profile.LowPriceToUserName = metaByID[lowPriceLink.BuyerID].Name
			profile.LowPriceItemName = lowPriceLink.ItemName
			profile.LowPriceDiscountRate = discountRate
			profile.LowPriceTradeCount = lowPriceLink.TradeCount

			if profile.TopBuyer != nil && profile.TopBuyer.UserID == lowPriceLink.BuyerID &&
				profile.TopBuyer.Concentration >= 60 && profile.TradeRole == "mostly_seller" {
				score += 26
				reasons = append(reasons, "卖得多、价格偏低且长期卖给同一玩家")
			} else if profile.TradeRole == "mostly_seller" {
				score += 12
				reasons = append(reasons, "存在低价卖货行为")
			}
		}
		if (profile.TopBuyer != nil && profile.TopBuyer.TradeDayCount >= 3) || (profile.TopSeller != nil && profile.TopSeller.TradeDayCount >= 3) {
			score += 10
			reasons = append(reasons, "固定交易关系持续多日")
		}
		if profile.AccountAgeDays > 0 && profile.AccountAgeDays <= 14 && profile.TotalTradeCount >= 6 {
			score += 12
			reasons = append(reasons, "账号较新但交易活跃")
		}
		if profile.RealRechargeAmount == 0 && profile.TotalTradeVolume >= 50000 {
			score += 12
			reasons = append(reasons, "真实充值为0但承接较高交易额")
		}
		if profile.VipLevel == 0 && profile.TotalTradeCount >= 10 {
			score += 8
			reasons = append(reasons, "低付费层级却高频参与交易")
		}

		if score > 100 {
			score = 100
		}
		profile.SuspiciousScore = score
		profile.SuspiciousAltLikelihood = float64(score)
		profile.SuspiciousReasons = reasons
		switch {
		case score >= 70:
			profile.SuspiciousLevel = "high"
		case score >= 40:
			profile.SuspiciousLevel = "medium"
		default:
			profile.SuspiciousLevel = "low"
		}

		profiles = append(profiles, *profile)
	}

	sort.Slice(profiles, func(i, j int) bool {
		if profiles[i].SuspiciousScore == profiles[j].SuspiciousScore {
			if profiles[i].TotalTradeCount == profiles[j].TotalTradeCount {
				return profiles[i].TotalTradeVolume > profiles[j].TotalTradeVolume
			}
			return profiles[i].TotalTradeCount > profiles[j].TotalTradeCount
		}
		return profiles[i].SuspiciousScore > profiles[j].SuspiciousScore
	})
	if len(profiles) > limit {
		profiles = profiles[:limit]
	}
	return profiles, nil
}

// RealmBottleneckItem 境界卡点分析
type RealmBottleneckItem struct {
	RealmID      int     `gorm:"column:realm_id"       json:"realm_id"`
	RealmName    string  `gorm:"column:realm_name"     json:"realm_name"`
	StageLevel   int     `gorm:"column:stage_level"    json:"stage_level"`
	UserCount    int     `gorm:"column:user_count"     json:"user_count"`
	AvgStayDays  float64 `gorm:"column:avg_stay_days"  json:"avg_stay_days"`
	MaxStayDays  int     `gorm:"column:max_stay_days"  json:"max_stay_days"`
	ChurnedCount int     `gorm:"column:churned_count"  json:"churned_count"`
	PassRate     float64 `json:"pass_rate"`
}

// GetRealmBottleneck 境界卡点分析（各境界停留时间+流失率+通过率）
func GetRealmBottleneck(inactiveDays int) ([]RealmBottleneckItem, error) {
	g := db.GetGame()
	if g == nil {
		return nil, ErrGameDBNil
	}
	if inactiveDays <= 0 {
		inactiveDays = 7
	}
	var items []RealmBottleneckItem
	err := g.Raw(`SELECT u.realm_id, r.name as realm_name, r.stage_level,
		COUNT(*) as user_count,
		AVG(DATEDIFF(NOW(), u.updated_at)) as avg_stay_days,
		MAX(DATEDIFF(NOW(), u.updated_at)) as max_stay_days,
		SUM(CASE WHEN u.updated_at < DATE_SUB(NOW(), INTERVAL ? DAY) THEN 1 ELSE 0 END) as churned_count
		FROM users u
		LEFT JOIN realms r ON r.id = u.realm_id
		WHERE u.deleted_at IS NULL AND u.status = 1
		GROUP BY u.realm_id, r.name, r.stage_level
		ORDER BY r.stage_level ASC, u.realm_id ASC`, inactiveDays).
		Scan(&items).Error
	if err != nil {
		return nil, err
	}
	// 计算通过率：当前境界的人数 vs 到达过更高境界的人数
	totalReached := 0
	for i := len(items) - 1; i >= 0; i-- {
		totalReached += items[i].UserCount
		if i < len(items)-1 {
			items[i].PassRate = float64(totalReached-items[i].UserCount) / float64(totalReached) * 100
		}
	}
	return items, nil
}

// ========== 用户状态统计 ==========

// UserStatusStats 用户状态统计
type UserStatusStats struct {
	Total   int `gorm:"column:total"    json:"total"`
	Active  int `gorm:"column:active"   json:"active"`
	Banned  int `gorm:"column:banned"   json:"banned"`
	Deleted int `gorm:"column:deleted"  json:"deleted"`
}

// GetUserStatusStats 获取用户状态分布
func GetUserStatusStats() (*UserStatusStats, error) {
	g := db.GetGame()
	if g == nil {
		return nil, ErrGameDBNil
	}
	var stats UserStatusStats
	err := g.Raw(`SELECT
		COUNT(*) as total,
		SUM(CASE WHEN status = 1 AND deleted_at IS NULL THEN 1 ELSE 0 END) as active,
		SUM(CASE WHEN status != 1 AND deleted_at IS NULL THEN 1 ELSE 0 END) as banned,
		SUM(CASE WHEN deleted_at IS NOT NULL THEN 1 ELSE 0 END) as deleted
		FROM users`).
		Scan(&stats).Error
	return &stats, err
}

// ========== 功法/技能分析 ==========

// SkillOverviewItem 功法概览
type SkillOverviewItem struct {
	GradeName  string `gorm:"column:grade_name"   json:"grade_name"`
	GradeLevel int    `gorm:"column:grade_level"  json:"grade_level"`
	SkillCount int    `gorm:"column:skill_count"  json:"skill_count"`
}

// GetSkillGradeDistribution 按品级统计功法数量（通过skill_books关联grade）
func GetSkillGradeDistribution() ([]SkillOverviewItem, error) {
	g := db.GetGame()
	if g == nil {
		return nil, ErrGameDBNil
	}
	items := make([]SkillOverviewItem, 0)
	err := g.Raw(`SELECT
		COALESCE(g.name, '未知品级') as grade_name,
		COALESCE(g.level, 0) as grade_level,
		COUNT(sb.id) as skill_count
		FROM skill_books sb
		LEFT JOIN grades g ON g.id = sb.grade_id
		GROUP BY COALESCE(g.id, 0), COALESCE(g.name, '未知品级'), COALESCE(g.level, 0)
		ORDER BY grade_level ASC`).
		Scan(&items).Error
	return items, err
}

// UserSkillStats 用户功法统计
type UserSkillStats struct {
	TotalUsers     int     `gorm:"column:total_users"      json:"total_users"`
	UsersWithSkill int     `gorm:"column:users_with_skill"  json:"users_with_skill"`
	AvgSkillCount  float64 `gorm:"column:avg_skill_count"   json:"avg_skill_count"`
	MaxSkillLevel  int     `gorm:"column:max_skill_level"   json:"max_skill_level"`
}

// GetUserSkillStats 获取用户功法概况
func GetUserSkillStats() (*UserSkillStats, error) {
	g := db.GetGame()
	if g == nil {
		return nil, ErrGameDBNil
	}
	var item UserSkillStats
	if err := g.Table("users").Select("COUNT(*) as total_users").
		Where("deleted_at IS NULL AND status = 1").Scan(&item).Error; err != nil {
		return nil, err
	}

	var aggregate struct {
		UsersWithSkill int     `gorm:"column:users_with_skill"`
		AvgSkillCount  float64 `gorm:"column:avg_skill_count"`
		MaxSkillLevel  int     `gorm:"column:max_skill_level"`
	}
	// 总活跃用户数
	// 拥有功法的用户数及平均功法数（user_skill_books存储用户功法）
	err := g.Raw(`SELECT COUNT(DISTINCT user_id) as users_with_skill,
		AVG(skill_cnt) as avg_skill_count,
		MAX(max_lv) as max_skill_level
		FROM (
			SELECT user_id, COUNT(*) as skill_cnt, MAX(upgrade_count) as max_lv
			FROM user_skill_books
			GROUP BY user_id
		) t`).Scan(&aggregate).Error
	if err != nil {
		return nil, err
	}

	item.UsersWithSkill = aggregate.UsersWithSkill
	item.AvgSkillCount = aggregate.AvgSkillCount
	item.MaxSkillLevel = aggregate.MaxSkillLevel
	return &item, err
}

// ========== 炼丹/炼器分析 ==========

// CraftProfessionItem 炼丹/炼器职业等级分布
type CraftProfessionItem struct {
	ProfessionName string `gorm:"column:profession_name" json:"profession_name"`
	LevelName      string `gorm:"column:level_name"      json:"level_name"`
	LevelRank      int    `gorm:"column:level_rank"       json:"level_rank"`
	UserCount      int    `gorm:"column:user_count"       json:"user_count"`
}

// GetCraftLevelDistribution 获取炼丹/炼器职业等级分布
func GetCraftLevelDistribution() ([]CraftProfessionItem, error) {
	g := db.GetGame()
	if g == nil {
		return nil, ErrGameDBNil
	}
	var items []CraftProfessionItem
	err := g.Raw(`SELECT p.name as profession_name, pl.title as level_name, pl.level as level_rank,
		COUNT(up.user_id) as user_count
		FROM professions p
		CROSS JOIN profession_levels pl ON pl.profession_id = p.id
		LEFT JOIN user_professions up ON up.profession_id = p.id AND up.level >= pl.level
		GROUP BY p.id, p.name, pl.id, pl.title, pl.level
		ORDER BY p.id ASC, pl.level ASC`).
		Scan(&items).Error
	return items, err
}

// CraftStatsItem 炼丹/炼器总体统计
type CraftStatsItem struct {
	TotalAlchemyRecipes int `gorm:"column:total_alchemy_recipes" json:"total_alchemy_recipes"`
	TotalLianqiRecipes  int `gorm:"column:total_lianqi_recipes"  json:"total_lianqi_recipes"`
	TotalCrafters       int `gorm:"column:total_crafters"   json:"total_crafters"`
}

// GetCraftStats 获取炼丹/炼器总体统计
func GetCraftStats() (*CraftStatsItem, error) {
	g := db.GetGame()
	if g == nil {
		return nil, ErrGameDBNil
	}
	var item CraftStatsItem
	type cntRow struct {
		Cnt int `gorm:"column:cnt"`
	}
	var alchemy, lianqi, crafters cntRow
	g.Raw(`SELECT COUNT(*) as cnt FROM alchemy_recipes`).Scan(&alchemy)
	g.Raw(`SELECT COUNT(*) as cnt FROM lianqi_recipes`).Scan(&lianqi)
	err := g.Raw(`SELECT COUNT(DISTINCT user_id) as cnt FROM user_professions`).Scan(&crafters).Error
	item.TotalAlchemyRecipes = alchemy.Cnt
	item.TotalLianqiRecipes = lianqi.Cnt
	item.TotalCrafters = crafters.Cnt
	return &item, err
}

// ========== 新人引导漏斗 ==========

// NewcomerFunnelItem 新人留存漏斗（Day1-7 每日签到人数）
type NewcomerFunnelItem struct {
	DayIndex  int `gorm:"column:day_index"  json:"day_index"`
	UserCount int `gorm:"column:user_count" json:"user_count"`
}

// GetNewcomerFunnel 获取新人7日签到漏斗
func GetNewcomerFunnel() ([]NewcomerFunnelItem, error) {
	g := db.GetGame()
	if g == nil {
		return nil, ErrGameDBNil
	}
	var items []NewcomerFunnelItem
	err := g.Raw(`SELECT day_index, COUNT(DISTINCT uid) as user_count
		FROM user_newcomer_login_claims
		GROUP BY day_index
		ORDER BY day_index ASC`).
		Scan(&items).Error
	return items, err
}

// ========== 道侣系统分析 ==========

// CompanionGiftUsage 道侣赠礼月度趋势
type CompanionGiftUsage struct {
	Month      string  `gorm:"column:month"      json:"month"`
	SendCount  int     `gorm:"column:send_count" json:"send_count"`
	UserCount  int     `gorm:"column:user_count" json:"user_count"`
	AvgPerUser float64 `gorm:"column:avg_per_user" json:"avg_per_user"`
}

// GetCompanionGiftUsage 获取道侣赠礼月度趋势
func GetCompanionGiftUsage() ([]CompanionGiftUsage, error) {
	g := db.GetGame()
	if g == nil {
		return nil, ErrGameDBNil
	}
	var items []CompanionGiftUsage
	err := g.Raw(`SELECT
		DATE_FORMAT(cdl.log_date, '%%Y-%%m') as month,
		SUM(cdl.completed_count) as send_count,
		COUNT(DISTINCT cdl.user_id) as user_count,
		SUM(cdl.completed_count) / NULLIF(COUNT(DISTINCT cdl.user_id), 0) as avg_per_user
		FROM companion_daily_logs cdl
		WHERE cdl.task_type = 'gift'
		GROUP BY DATE_FORMAT(cdl.log_date, '%%Y-%%m')
		ORDER BY month ASC`).
		Scan(&items).Error
	return items, err
}

// CompanionStatsItem 道侣总体数据
type CompanionStatsItem struct {
	TotalCouples      int `gorm:"column:total_couples"       json:"total_couples"`
	NormalOathCount   int `gorm:"column:normal_oath_count"   json:"normal_oath_count"`
	AdvancedOathCount int `gorm:"column:advanced_oath_count" json:"advanced_oath_count"`
	TotalGiftsSent    int `gorm:"column:total_gifts_sent"    json:"total_gifts_sent"`
}

// GetCompanionStats 获取道侣系统总体统计
func GetCompanionStats() (*CompanionStatsItem, error) {
	g := db.GetGame()
	if g == nil {
		return nil, ErrGameDBNil
	}
	type cntRow struct {
		Cnt int `gorm:"column:cnt"`
	}
	var couples, normalOath, advancedOath, gifts cntRow
	g.Raw(`SELECT COUNT(*) as cnt FROM friendships WHERE deleted_at IS NULL AND relation_type = 2`).Scan(&couples)
	g.Raw(`SELECT COUNT(*) as cnt FROM friendships WHERE deleted_at IS NULL AND relation_type = 2 AND oath_type = 'normal'`).Scan(&normalOath)
	g.Raw(`SELECT COUNT(*) as cnt FROM friendships WHERE deleted_at IS NULL AND relation_type = 2 AND oath_type = 'advanced'`).Scan(&advancedOath)
	err := g.Raw(`SELECT COALESCE(SUM(completed_count),0) as cnt
		FROM companion_daily_logs WHERE task_type = 'gift'`).Scan(&gifts).Error
	return &CompanionStatsItem{
		TotalCouples:      couples.Cnt,
		NormalOathCount:   normalOath.Cnt,
		AdvancedOathCount: advancedOath.Cnt,
		TotalGiftsSent:    gifts.Cnt,
	}, err
}

// CompanionIntimacyDist 灵契值分布
type CompanionIntimacyDist struct {
	Bracket   string `gorm:"column:bracket"    json:"bracket"`
	PairCount int    `gorm:"column:pair_count" json:"pair_count"`
}

// GetCompanionIntimacyDist 获取灵契值分布
func GetCompanionIntimacyDist() ([]CompanionIntimacyDist, error) {
	g := db.GetGame()
	if g == nil {
		return nil, ErrGameDBNil
	}
	var items []CompanionIntimacyDist
	err := g.Raw(`SELECT
		CASE
			WHEN intimacy < 100  THEN '0-99'
			WHEN intimacy < 500  THEN '100-499'
			WHEN intimacy < 1000 THEN '500-999'
			WHEN intimacy < 2000 THEN '1000-1999'
			WHEN intimacy < 5000 THEN '2000-4999'
			ELSE '5000+'
		END as bracket,
		COUNT(*) as pair_count
		FROM friendships
		WHERE deleted_at IS NULL AND relation_type = 2
		GROUP BY bracket
		ORDER BY MIN(intimacy) ASC`).
		Scan(&items).Error
	return items, err
}

// ========== 本命物系统分析 ==========

// SoulArtifactCategoryDist 本命物分类分布
type SoulArtifactCategoryDist struct {
	Category  string `gorm:"column:category"   json:"category"`
	UserCount int    `gorm:"column:user_count" json:"user_count"`
}

// GetSoulArtifactCategoryDist 获取用户本命物分类分布
func GetSoulArtifactCategoryDist() ([]SoulArtifactCategoryDist, error) {
	g := db.GetGame()
	if g == nil {
		return nil, ErrGameDBNil
	}
	var items []SoulArtifactCategoryDist
	err := g.Raw(`SELECT sa.category,
		COUNT(DISTINCT usa.user_id) as user_count
		FROM user_soul_artifacts usa
		LEFT JOIN soul_artifacts sa ON sa.id = usa.artifact_id
		WHERE usa.deleted_at IS NULL
		GROUP BY sa.category
		ORDER BY user_count DESC`).
		Scan(&items).Error
	return items, err
}

// SoulArtifactLevelDist 本命物等级分布
type SoulArtifactLevelDist struct {
	Bracket   string `gorm:"column:bracket"    json:"bracket"`
	UserCount int    `gorm:"column:user_count" json:"user_count"`
}

// GetSoulArtifactLevelDist 获取本命物等级分布
func GetSoulArtifactLevelDist() ([]SoulArtifactLevelDist, error) {
	g := db.GetGame()
	if g == nil {
		return nil, ErrGameDBNil
	}
	var items []SoulArtifactLevelDist
	err := g.Raw(`SELECT
		CASE
			WHEN level <= 5  THEN '1-5'
			WHEN level <= 10 THEN '6-10'
			WHEN level <= 20 THEN '11-20'
			WHEN level <= 30 THEN '21-30'
			ELSE '31+'
		END as bracket,
		COUNT(DISTINCT user_id) as user_count
		FROM user_soul_artifacts
		WHERE deleted_at IS NULL
		GROUP BY bracket
		ORDER BY MIN(level) ASC`).
		Scan(&items).Error
	return items, err
}

// SoulArtifactUserStats 本命物用户总体统计
type SoulArtifactUserStats struct {
	TotalOwners   int     `gorm:"column:total_owners"     json:"total_owners"`
	AvgGradeLevel float64 `gorm:"column:avg_grade_level"  json:"avg_grade_level"`
	AvgLevel      float64 `gorm:"column:avg_level"        json:"avg_level"`
	MaxLevel      int     `gorm:"column:max_level"        json:"max_level"`
}

// GetSoulArtifactUserStats 获取本命物用户统计
func GetSoulArtifactUserStats() (*SoulArtifactUserStats, error) {
	g := db.GetGame()
	if g == nil {
		return nil, ErrGameDBNil
	}
	var item SoulArtifactUserStats
	err := g.Raw(`SELECT
		COUNT(DISTINCT user_id) as total_owners,
		AVG(grade_id) as avg_grade_level,
		AVG(level) as avg_level,
		MAX(level) as max_level
		FROM user_soul_artifacts
		WHERE deleted_at IS NULL`).
		Scan(&item).Error
	return &item, err
}

// SoulArtifactGradeDist 用户本命物品级分布
type SoulArtifactGradeDist struct {
	GradeName string `gorm:"column:grade_name" json:"grade_name"`
	GradeID   int    `gorm:"column:grade_id"   json:"grade_id"`
	UserCount int    `gorm:"column:user_count" json:"user_count"`
}

// GetSoulArtifactGradeDist 获取用户本命物品级分布
func GetSoulArtifactGradeDist() ([]SoulArtifactGradeDist, error) {
	g := db.GetGame()
	if g == nil {
		return nil, ErrGameDBNil
	}
	var items []SoulArtifactGradeDist
	err := g.Raw(`SELECT g.name as grade_name, usa.grade_id,
		COUNT(DISTINCT usa.user_id) as user_count
		FROM user_soul_artifacts usa
		LEFT JOIN grades g ON usa.grade_id = g.id
		WHERE usa.deleted_at IS NULL
		GROUP BY usa.grade_id, g.name
		ORDER BY usa.grade_id ASC`).
		Scan(&items).Error
	return items, err
}

// SoulArtifactPopularity 本命物受欢迎度
type SoulArtifactPopularity struct {
	ArtifactID   uint   `gorm:"column:artifact_id"   json:"artifact_id"`
	ArtifactName string `gorm:"column:artifact_name"  json:"artifact_name"`
	Category     string `gorm:"column:category"       json:"category"`
	OwnerCount   int    `gorm:"column:owner_count"    json:"owner_count"`
}

// GetSoulArtifactPopularity 获取本命物受欢迎度排行
func GetSoulArtifactPopularity() ([]SoulArtifactPopularity, error) {
	g := db.GetGame()
	if g == nil {
		return nil, ErrGameDBNil
	}
	var items []SoulArtifactPopularity
	err := g.Raw(`SELECT sa.id as artifact_id, sa.name as artifact_name, sa.category,
		COUNT(usa.user_id) as owner_count
		FROM soul_artifacts sa
		LEFT JOIN user_soul_artifacts usa ON usa.artifact_id = sa.id AND usa.deleted_at IS NULL
		GROUP BY sa.id, sa.name, sa.category
		ORDER BY owner_count DESC`).
		Scan(&items).Error
	return items, err
}

// ========== 商店消费分析 ==========

// ShopPriceRangeDist 商店价格区间分布
type ShopPriceRangeDist struct {
	ShopCode     string `gorm:"column:shop_code"     json:"shop_code"`
	PriceRange   string `gorm:"column:price_range"   json:"price_range"`
	ProductCount int    `gorm:"column:product_count" json:"product_count"`
}

// GetShopPriceRangeDist 获取各商店价格区间分布
func GetShopPriceRangeDist() ([]ShopPriceRangeDist, error) {
	g := db.GetGame()
	if g == nil {
		return nil, ErrGameDBNil
	}
	var items []ShopPriceRangeDist
	err := g.Raw(`SELECT shop_code,
		CASE
			WHEN price_amount <= 100    THEN '0-100'
			WHEN price_amount <= 500    THEN '101-500'
			WHEN price_amount <= 2000   THEN '501-2000'
			WHEN price_amount <= 10000  THEN '2001-10000'
			ELSE '10000+'
		END as price_range,
		COUNT(*) as product_count
		FROM wanbaolou_products
		WHERE status = 1
		GROUP BY shop_code, price_range
		ORDER BY shop_code ASC, MIN(price_amount) ASC`).
		Scan(&items).Error
	return items, err
}

// ShopOverview 商店总览
type ShopOverview struct {
	ShopCode      string `gorm:"column:shop_code"      json:"shop_code"`
	ProductCount  int    `gorm:"column:product_count"   json:"product_count"`
	CategoryCount int    `gorm:"column:category_count"  json:"category_count"`
}

// GetShopOverview 获取各商店概览
func GetShopOverview() ([]ShopOverview, error) {
	g := db.GetGame()
	if g == nil {
		return nil, ErrGameDBNil
	}
	var items []ShopOverview
	err := g.Raw(`SELECT shop_code,
		COUNT(*) as product_count,
		COUNT(DISTINCT category) as category_count
		FROM wanbaolou_products
		WHERE status = 1
		GROUP BY shop_code
		ORDER BY shop_code ASC`).
		Scan(&items).Error
	return items, err
}

// ShopCurrencyDist 商店货币分布
type ShopCurrencyDist struct {
	ShopCode      string  `gorm:"column:shop_code"       json:"shop_code"`
	PriceCurrency string  `gorm:"column:price_currency"  json:"price_currency"`
	ProductCount  int     `gorm:"column:product_count"   json:"product_count"`
	AvgPrice      float64 `gorm:"column:avg_price"       json:"avg_price"`
}

// GetShopCurrencyDist 获取各商店货币使用分布
func GetShopCurrencyDist() ([]ShopCurrencyDist, error) {
	g := db.GetGame()
	if g == nil {
		return nil, ErrGameDBNil
	}
	var items []ShopCurrencyDist
	err := g.Raw(`SELECT shop_code, price_currency_code as price_currency,
		COUNT(*) as product_count,
		AVG(price_amount) as avg_price
		FROM wanbaolou_products
		WHERE status = 1
		GROUP BY shop_code, price_currency_code
		ORDER BY shop_code ASC`).
		Scan(&items).Error
	return items, err
}

// ShopCategoryDist 商品分类分布
type ShopCategoryDist struct {
	ShopCode     string `gorm:"column:shop_code"     json:"shop_code"`
	Category     string `gorm:"column:category"      json:"category"`
	ProductCount int    `gorm:"column:product_count" json:"product_count"`
}

// GetShopCategoryDist 获取商品分类分布
func GetShopCategoryDist() ([]ShopCategoryDist, error) {
	g := db.GetGame()
	if g == nil {
		return nil, ErrGameDBNil
	}
	var items []ShopCategoryDist
	err := g.Raw(`SELECT shop_code, category,
		COUNT(*) as product_count
		FROM wanbaolou_products
		WHERE status = 1
		GROUP BY shop_code, category
		ORDER BY shop_code ASC, product_count DESC`).
		Scan(&items).Error
	return items, err
}
