package segmentapi

import (
	"strconv"
	"yonghemolimis/src/apps/api/response"
	"yonghemolimis/src/dao/analytics"
	"yonghemolimis/src/dao/db"
	"yonghemolimis/src/usecase/segment"

	"github.com/gin-gonic/gin"
)

// List 获取所有分群
func List(c *gin.Context) {
	segs, err := analytics.GetSegments()
	if err != nil {
		response.Fail(c, "查询分群失败: "+err.Error())
		return
	}
	response.OK(c, segs)
}

type createReq struct {
	Name        string `json:"name" form:"name" binding:"required"`
	Description string `json:"description" form:"description"`
	RulesJSON   string `json:"rules_json" form:"rules_json" binding:"required"`
}

// Create 创建分群
func Create(c *gin.Context) {
	var req createReq
	if err := c.ShouldBind(&req); err != nil {
		response.Fail(c, "参数错误")
		return
	}

	// 验证规则 JSON 合法性
	count, err := segment.CountSegment(req.RulesJSON)
	if err != nil {
		response.Fail(c, "规则校验失败: "+err.Error())
		return
	}

	username, _ := c.Get("username")
	seg := &db.SegmentDO{
		Name:        req.Name,
		Description: req.Description,
		RulesJSON:   req.RulesJSON,
		UserCount:   count,
		CreatedBy:   username.(string),
	}
	if err := analytics.CreateSegment(seg); err != nil {
		response.Fail(c, "创建分群失败: "+err.Error())
		return
	}
	response.OK(c, seg)
}

// Delete 删除分群
func Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Fail(c, "无效的分群 ID")
		return
	}
	if err := analytics.DeleteSegment(uint(id)); err != nil {
		response.Fail(c, "删除分群失败: "+err.Error())
		return
	}
	response.OKMsg(c, "分群已删除")
}

// Execute 执行分群查询，返回匹配用户列表
func Execute(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Fail(c, "无效的分群 ID")
		return
	}

	// 获取分群
	seg, err := analytics.GetSegmentByID(uint(id))
	if err != nil {
		response.Fail(c, "分群不存在")
		return
	}

	// 执行查询
	profiles, total, err := segment.ExecuteSegment(seg)
	if err != nil {
		response.Fail(c, "执行分群查询失败: "+err.Error())
		return
	}

	// 更新人数
	analytics.UpdateSegmentCount(seg.ID, total)

	response.OK(c, gin.H{
		"total":    total,
		"profiles": profiles,
	})
}

type previewReq struct {
	RulesJSON string `json:"rules_json" form:"rules_json" binding:"required"`
}

// Preview 预览分群规则匹配人数（不保存）
func Preview(c *gin.Context) {
	var req previewReq
	if err := c.ShouldBind(&req); err != nil {
		response.Fail(c, "参数错误")
		return
	}

	count, err := segment.CountSegment(req.RulesJSON)
	if err != nil {
		response.Fail(c, "规则解析失败: "+err.Error())
		return
	}

	response.OK(c, gin.H{"count": count})
}
