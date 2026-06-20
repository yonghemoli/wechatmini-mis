package entity

import (
	"time"
	"yonghemolimis/src/domain/valueobj"
)

// UserProfile 用户画像实体
type UserProfile struct {
	UID              uint
	LifecycleStage   valueobj.LifecycleStage
	PlayStyle        valueobj.PlayStyle
	PayTier          valueobj.PayTier
	LtvPredict       float64
	SocialType       string
	ChurnRisk        int // 0-100
	StuckFlag        bool
	ResourceAlert    bool
	LastCalculatedAt *time.Time
}

// UserTag 用户标签
type UserTag struct {
	UID      uint
	TagGroup string
	TagKey   string
	TagValue string
	Score    *float64
	ExpireAt *time.Time
}

// DailySnapshot 每日快照
type DailySnapshot struct {
	UID            uint
	SnapshotDate   string // "2006-01-02"
	RealmID        int
	VipLevel       int
	SpiritStone    int64
	SpiritCrystal  int64
	SpiritMilk     int64
	LoginCount     int
	BattleCount    int
	CraftCount     int
	SocialCount    int
	EconomyCount   int
	StoneIncome    int64
	StoneExpense   int64
	CrystalIncome  int64
	CrystalExpense int64
}

// ActionLog 行为日志
type ActionLog struct {
	UID       uint
	Action    string
	Category  string
	Detail    string
	CreatedAt time.Time
}

// Segment 用户分群
type Segment struct {
	ID          uint
	Name        string
	Description string
	RulesJSON   string
	UserCount   int
	CreatedBy   string
	CreatedAt   time.Time
}

// GameUser 游戏库用户（只读映射）
type GameUser struct {
	ID                  uint
	Name                string
	Gender              int
	RealmID             int
	VipLevel            int
	Experience          int
	SpiritStone         int64
	SpiritCrystal       int64
	SpiritMilk          int64
	TotalRechargeAmount int64
	FirstRechargeTime   *time.Time
	MonthCardExpireTime *time.Time
	SignInCoin          int64
	TowerScore          int64
	Reputation          int64
	Tongxin             int64
	CreatedAt           time.Time
	UpdatedAt           time.Time
}

// GameRealm 境界（只读映射）
type GameRealm struct {
	ID         uint
	Name       string
	Level      int
	StageLevel int
}

// CurrencyLog 货币变更流水（只读映射）
type CurrencyLog struct {
	UID          uint
	CurrencyCode string
	ChangeAmount int64
	Detail       string
	CreatedAt    time.Time
}
