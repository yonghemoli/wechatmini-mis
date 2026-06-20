package track

import (
	"net/http"
	"time"
	"yonghemolimis/src/apps/api/response"
	"yonghemolimis/src/dao/db"

	"github.com/gin-gonic/gin"
)

// FeatureEventPayload 单条功能埋点事件
type FeatureEventPayload struct {
	UID             uint   `json:"uid"`
	FeatureName     string `json:"feature_name"`
	FeatureCategory string `json:"feature_category"`
	CommandText     string `json:"command_text"`
	ResponseTimeMs  int    `json:"response_time_ms"`
	Success         bool   `json:"success"`
	Scene           string `json:"scene"`      // private / group
	ChannelID       string `json:"channel_id"` // 频道/群 ID
	Platform        string `json:"platform"`   // 平台标识
	Timestamp       string `json:"timestamp"`  // ISO8601
}

// BatchRequest 批量上报请求
type BatchRequest struct {
	Events []FeatureEventPayload `json:"events"`
}

// ReceiveFeatureEvents 接收游戏 SDK 上报的功能埋点事件
// POST /api/v1/track/features
func ReceiveFeatureEvents(c *gin.Context) {
	var req BatchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.R{Code: -1, Message: "参数错误"})
		return
	}

	if len(req.Events) == 0 {
		response.OKMsg(c, "无事件")
		return
	}

	// 转换为 DO 并批量写入
	dos := make([]db.FeatureEventDO, 0, len(req.Events))
	for _, e := range req.Events {
		eventTime := time.Now()
		if e.Timestamp != "" {
			if t, err := time.Parse(time.RFC3339, e.Timestamp); err == nil {
				eventTime = t
			}
		}
		dos = append(dos, db.FeatureEventDO{
			UID:             e.UID,
			FeatureName:     e.FeatureName,
			FeatureCategory: e.FeatureCategory,
			CommandText:     e.CommandText,
			ResponseTimeMs:  e.ResponseTimeMs,
			Success:         e.Success,
			Scene:           e.Scene,
			ChannelID:       e.ChannelID,
			Platform:        e.Platform,
			EventTime:       eventTime,
		})
	}

	// 批量插入（分批，每批 500 条）
	batchSize := 500
	d := db.Get()
	for i := 0; i < len(dos); i += batchSize {
		end := i + batchSize
		if end > len(dos) {
			end = len(dos)
		}
		if err := d.CreateInBatches(dos[i:end], batchSize).Error; err != nil {
			response.Fail(c, "写入失败: "+err.Error())
			return
		}
	}

	response.OK(c, gin.H{"received": len(dos)})
}

// ==================== 行为事件上报 ====================

// ActionEventPayload 单条行为事件
type ActionEventPayload struct {
	UID       uint   `json:"uid"`
	Action    string `json:"action"`     // 行为标识: breakthrough, dungeon, spar, etc.
	Category  string `json:"category"`   // 系统分类: 修炼, 战斗, 社交, 经济, 成长
	Result    string `json:"result"`     // 结果: success / fail / draw
	Value     int64  `json:"value"`      // 关键数值
	Detail    string `json:"detail"`     // JSON 扩展详情
	Scene     string `json:"scene"`      // private / group
	ChannelID string `json:"channel_id"` // 频道/群 ID
	Platform  string `json:"platform"`   // 平台标识
	Timestamp string `json:"timestamp"`  // ISO8601
}

// ActionBatchRequest 批量行为事件请求
type ActionBatchRequest struct {
	Events []ActionEventPayload `json:"events"`
}

// ReceiveActionEvents 接收游戏 SDK 上报的行为事件
// POST /api/v1/track/actions
func ReceiveActionEvents(c *gin.Context) {
	var req ActionBatchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.R{Code: -1, Message: "参数错误"})
		return
	}

	if len(req.Events) == 0 {
		response.OKMsg(c, "无事件")
		return
	}

	dos := make([]db.GameActionEventDO, 0, len(req.Events))
	for _, e := range req.Events {
		eventTime := time.Now()
		if e.Timestamp != "" {
			if t, err := time.Parse(time.RFC3339, e.Timestamp); err == nil {
				eventTime = t
			}
		}
		dos = append(dos, db.GameActionEventDO{
			UID:       e.UID,
			Action:    e.Action,
			Category:  e.Category,
			Result:    e.Result,
			Value:     e.Value,
			Detail:    e.Detail,
			Scene:     e.Scene,
			ChannelID: e.ChannelID,
			Platform:  e.Platform,
			EventTime: eventTime,
		})
	}

	batchSize := 500
	d := db.Get()
	for i := 0; i < len(dos); i += batchSize {
		end := i + batchSize
		if end > len(dos) {
			end = len(dos)
		}
		if err := d.CreateInBatches(dos[i:end], batchSize).Error; err != nil {
			response.Fail(c, "写入失败: "+err.Error())
			return
		}
	}

	response.OK(c, gin.H{"received": len(dos)})
}
