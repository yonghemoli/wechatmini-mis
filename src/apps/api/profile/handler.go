package profileapi

import (
	"strconv"
	"yonghemolimis/src/apps/api/response"
	"yonghemolimis/src/dao/analytics"
	"yonghemolimis/src/dao/game"
	profileuc "yonghemolimis/src/usecase/profile"

	"github.com/gin-gonic/gin"
)

// GetProfile 获取用户画像（含游戏数据和标签）
func GetProfile(c *gin.Context) {
	uid, err := strconv.ParseUint(c.Param("uid"), 10, 64)
	if err != nil {
		response.Fail(c, "无效的用户 ID")
		return
	}
	p, err := analytics.GetProfile(uint(uid))
	if err != nil {
		response.Fail(c, "用户画像不存在")
		return
	}

	// 获取游戏用户信息
	gameUser, _ := game.GetUserByID(uint(uid))

	// 获取用户标签
	tags, _ := analytics.GetTagsByUID(uint(uid))

	response.OK(c, gin.H{
		"profile":   p,
		"game_user": gameUser,
		"tags":      tags,
	})
}

// ListProfiles 画像列表
func ListProfiles(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	filter := analytics.ProfileFilter{
		LifecycleStage: c.Query("lifecycle_stage"),
		PayTier:        c.Query("pay_tier"),
		PlayStyle:      c.Query("play_style"),
		StuckOnly:      c.Query("stuck_only") == "true",
		Page:           page,
		PageSize:       pageSize,
	}

	if mr := c.Query("min_churn_risk"); mr != "" {
		filter.MinChurnRisk, _ = strconv.Atoi(mr)
	}

	list, total, err := analytics.GetProfilesByFilter(filter)
	if err != nil {
		response.Fail(c, "查询失败: "+err.Error())
		return
	}
	response.OKPage(c, list, total, page, pageSize)
}

// RefreshProfile 刷新单个用户画像
func RefreshProfile(c *gin.Context) {
	uid, err := strconv.ParseUint(c.Param("uid"), 10, 64)
	if err != nil {
		response.Fail(c, "无效的用户 ID")
		return
	}
	if err := profileuc.CalculateProfile(uint(uid)); err != nil {
		response.Fail(c, "计算失败: "+err.Error())
		return
	}
	response.OKMsg(c, "画像已刷新")
}

// RefreshAllProfiles 刷新所有用户画像
func RefreshAllProfiles(c *gin.Context) {
	go func() {
		if err := profileuc.CalculateAllProfiles(); err != nil {
			// 已在内部记录日志
			_ = err
		}
	}()
	response.OKMsg(c, "全量画像计算已在后台启动")
}

// GetSnapshots 获取用户快照
func GetSnapshots(c *gin.Context) {
	uid, err := strconv.ParseUint(c.Param("uid"), 10, 64)
	if err != nil {
		response.Fail(c, "无效的用户 ID")
		return
	}
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "30"))
	snaps, err := analytics.GetSnapshotsByUID(uint(uid), limit)
	if err != nil {
		response.Fail(c, "查询快照失败: "+err.Error())
		return
	}
	response.OK(c, snaps)
}
