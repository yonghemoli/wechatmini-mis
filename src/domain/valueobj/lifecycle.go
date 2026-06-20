package valueobj

// LifecycleStage 用户生命周期阶段
type LifecycleStage string

const (
	LifecycleNew       LifecycleStage = "NEW"       // 新手期 (0-3天)
	LifecycleGrowing   LifecycleStage = "GROWING"   // 成长期 (3-14天)
	LifecycleMature    LifecycleStage = "MATURE"    // 成熟期 (14天+, 活跃)
	LifecycleDeclining LifecycleStage = "DECLINING" // 衰退期 (活跃度下降)
	LifecycleLost      LifecycleStage = "LOST"      // 流失期 (7天未登录)
	LifecycleReturned  LifecycleStage = "RETURNED"  // 回流期
)

func (l LifecycleStage) String() string { return string(l) }

func (l LifecycleStage) Label() string {
	switch l {
	case LifecycleNew:
		return "新手期"
	case LifecycleGrowing:
		return "成长期"
	case LifecycleMature:
		return "成熟期"
	case LifecycleDeclining:
		return "衰退期"
	case LifecycleLost:
		return "流失期"
	case LifecycleReturned:
		return "回流期"
	default:
		return "未知"
	}
}

var AllLifecycleStages = []LifecycleStage{
	LifecycleNew, LifecycleGrowing, LifecycleMature,
	LifecycleDeclining, LifecycleLost, LifecycleReturned,
}
