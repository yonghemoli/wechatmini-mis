package mis

import (
	"net/http"
	"strconv"

	"yonghemolimis/src/apps/api/chatws"
	"yonghemolimis/src/apps/api/response"
	"yonghemolimis/src/dao/db"

	"github.com/gin-gonic/gin"
)

func ListFAQs(c *gin.Context) {
	rows, err := db.ListFAQs(c.Query("category"))
	if err != nil {
		response.Error(c, 500, "查询常见问题失败")
		return
	}
	response.OK(c, gin.H{"list": rows})
}
func CreateFAQ(c *gin.Context) {
	var req db.FAQDO
	if c.ShouldBindJSON(&req) != nil {
		response.Error(c, 400, "参数错误")
		return
	}
	if err := db.CreateFAQ(&req); err != nil {
		response.Error(c, 400, err.Error())
		return
	}
	response.OK(c, gin.H{"item": req})
}
func UpdateFAQ(c *gin.Context) {
	id, ok := uintParam(c)
	if !ok {
		return
	}
	var req db.FAQDO
	if c.ShouldBindJSON(&req) != nil {
		response.Error(c, 400, "参数错误")
		return
	}
	row, err := db.UpdateFAQ(id, req)
	if err != nil {
		response.Error(c, 400, err.Error())
		return
	}
	response.OK(c, gin.H{"item": row})
}
func DeleteFAQ(c *gin.Context) {
	id, ok := uintParam(c)
	if !ok {
		return
	}
	if err := db.DeleteFAQ(id); err != nil {
		response.Error(c, 400, err.Error())
		return
	}
	response.OKMsg(c, "常见问题已删除")
}
func PublishFAQ(c *gin.Context)   { updateFAQVisible(c, true) }
func UnpublishFAQ(c *gin.Context) { updateFAQVisible(c, false) }
func updateFAQVisible(c *gin.Context, visible bool) {
	id, ok := uintParam(c)
	if !ok {
		return
	}
	row, err := db.UpdateFAQVisible(id, visible)
	if err != nil {
		response.Error(c, 400, err.Error())
		return
	}
	response.OK(c, gin.H{"item": row})
}

func ListChatSessions(c *gin.Context) {
	rows, err := db.ListChatSessions()
	if err != nil {
		response.Error(c, 500, "查询客服会话失败")
		return
	}
	response.OK(c, gin.H{"list": rows})
}
func ListChatMessages(c *gin.Context) {
	rows, err := db.ListChatMessages(c.Param("id"))
	if err != nil {
		response.Error(c, 500, "查询客服消息失败")
		return
	}
	response.OK(c, gin.H{"list": rows})
}
func CreateChatMessage(c *gin.Context) {
	var req struct {
		Content string `json:"content" binding:"required"`
		MsgType string `json:"msgType"`
	}
	if c.ShouldBindJSON(&req) != nil {
		response.Error(c, 400, "消息内容不能为空")
		return
	}
	row, err := db.CreateChatMessage(c.Param("id"), "admin", req.MsgType, req.Content)
	if err != nil {
		response.Error(c, 400, err.Error())
		return
	}
	chatws.BroadcastMessage(row)
	response.OK(c, gin.H{"item": row})
}
func CloseChatSession(c *gin.Context) {
	if err := db.CloseChatSession(c.Param("id")); err != nil {
		response.Error(c, 400, err.Error())
		return
	}
	response.OKMsg(c, "会话已关闭")
}
func ReadChatSession(c *gin.Context) {
	if err := db.ReadChatSession(c.Param("id")); err != nil {
		response.Error(c, 400, err.Error())
		return
	}
	response.OKMsg(c, "已标记为已读")
}

func uintParam(c *gin.Context) (uint, bool) {
	value, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || value == 0 {
		response.Error(c, http.StatusBadRequest, "ID 参数错误")
		return 0, false
	}
	return uint(value), true
}
