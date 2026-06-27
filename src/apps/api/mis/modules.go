package mis

import (
	"encoding/csv"
	"net/http"
	"strconv"

	"yonghemolimis/src/apps/api/chatws"
	"yonghemolimis/src/apps/api/response"
	"yonghemolimis/src/dao/db"

	"github.com/gin-gonic/gin"
)

func ListServiceTypes(c *gin.Context) {
	rows, err := db.ListServiceTypes(c.Query("keyword"))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.OK(c, gin.H{"list": rows})
}

func CreateServiceType(c *gin.Context) {
	var req db.ServiceTypeDO
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, "参数错误")
		return
	}
	if err := db.CreateServiceType(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(c, gin.H{"item": req})
}

func UpdateServiceType(c *gin.Context) {
	id, ok := uintParam(c)
	if !ok {
		return
	}
	var req db.ServiceTypeDO
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, "参数错误")
		return
	}
	row, err := db.UpdateServiceType(id, req)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(c, gin.H{"item": row})
}

func DeleteServiceType(c *gin.Context) {
	id, ok := uintParam(c)
	if !ok {
		return
	}
	if err := db.DeleteServiceType(id); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	response.OKMsg(c, "服务类型已删除")
}

func EnableServiceType(c *gin.Context) {
	updateServiceTypeStatus(c, db.StatusActive)
}

func DisableServiceType(c *gin.Context) {
	updateServiceTypeStatus(c, db.StatusDisabled)
}

func updateServiceTypeStatus(c *gin.Context, status string) {
	id, ok := uintParam(c)
	if !ok {
		return
	}
	row, err := db.UpdateServiceTypeStatus(id, status)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(c, gin.H{"item": row})
}

func ListServices(c *gin.Context) {
	typeID, _ := strconv.Atoi(c.Query("typeId"))
	rows, err := db.ListServices(uint(typeID), c.Query("keyword"))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.OK(c, gin.H{"list": rows})
}

func CreateService(c *gin.Context) {
	var req db.ServiceDO
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, "参数错误")
		return
	}
	if err := db.CreateService(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(c, gin.H{"item": req})
}

func UpdateService(c *gin.Context) {
	id, ok := uintParam(c)
	if !ok {
		return
	}
	var req db.ServiceDO
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, "参数错误")
		return
	}
	row, err := db.UpdateService(id, req)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(c, gin.H{"item": row})
}

func DeleteService(c *gin.Context) {
	id, ok := uintParam(c)
	if !ok {
		return
	}
	if err := db.DeleteService(id); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	response.OKMsg(c, "服务已删除")
}

func PublishService(c *gin.Context) {
	updateServiceVisible(c, true)
}

func UnpublishService(c *gin.Context) {
	updateServiceVisible(c, false)
}

func updateServiceVisible(c *gin.Context, visible bool) {
	id, ok := uintParam(c)
	if !ok {
		return
	}
	row, err := db.UpdateServiceVisible(id, visible)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(c, gin.H{"item": row})
}

func ExportServices(c *gin.Context) {
	rows, err := db.ListServices(0, "")
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	writeCSV(c, "家政服务列表.csv", []string{"编号", "服务类型", "服务名称", "价格", "单位", "上架", "排序"}, func(w *csv.Writer) {
		for _, row := range rows {
			visible := "否"
			if row.Visible {
				visible = "是"
			}
			_ = w.Write([]string{
				strconv.Itoa(int(row.ID)),
				row.TypeName,
				row.Name,
				strconv.Itoa(row.Price),
				row.Unit,
				visible,
				strconv.Itoa(row.SortOrder),
			})
		}
	})
}

func ListShops(c *gin.Context) {
	rows, err := db.ListShops()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.OK(c, gin.H{"list": rows})
}

func CreateShop(c *gin.Context) {
	var req db.ShopDO
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, "参数错误")
		return
	}
	if err := db.CreateShop(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(c, gin.H{"item": req})
}

func UpdateShop(c *gin.Context) {
	id, ok := uintParam(c)
	if !ok {
		return
	}
	var req db.ShopDO
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, "参数错误")
		return
	}
	row, err := db.UpdateShop(id, req)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(c, gin.H{"item": row})
}

func DeleteShop(c *gin.Context) {
	id, ok := uintParam(c)
	if !ok {
		return
	}
	if err := db.DeleteShop(id); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	response.OKMsg(c, "店铺已删除")
}

func OpenShop(c *gin.Context) {
	updateShopStatus(c, db.StatusOpen)
}

func CloseShop(c *gin.Context) {
	updateShopStatus(c, db.StatusClosed)
}

func updateShopStatus(c *gin.Context, status string) {
	id, ok := uintParam(c)
	if !ok {
		return
	}
	row, err := db.UpdateShopStatus(id, status)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(c, gin.H{"item": row})
}

func ListFAQs(c *gin.Context) {
	rows, err := db.ListFAQs(c.Query("category"))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.OK(c, gin.H{"list": rows})
}

func CreateFAQ(c *gin.Context) {
	var req db.FAQDO
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, "参数错误")
		return
	}
	if err := db.CreateFAQ(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
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
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, "参数错误")
		return
	}
	row, err := db.UpdateFAQ(id, req)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
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
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	response.OKMsg(c, "常见问题已删除")
}

func PublishFAQ(c *gin.Context) {
	updateFAQVisible(c, true)
}

func UnpublishFAQ(c *gin.Context) {
	updateFAQVisible(c, false)
}

func updateFAQVisible(c *gin.Context, visible bool) {
	id, ok := uintParam(c)
	if !ok {
		return
	}
	row, err := db.UpdateFAQVisible(id, visible)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(c, gin.H{"item": row})
}

func ListChatSessions(c *gin.Context) {
	rows, err := db.ListChatSessions()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.OK(c, gin.H{"list": rows})
}

func ListChatMessages(c *gin.Context) {
	rows, err := db.ListChatMessages(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.OK(c, gin.H{"list": rows})
}

func CreateChatMessage(c *gin.Context) {
	var req struct {
		Content string `json:"content" binding:"required"`
		MsgType string `json:"msgType"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, "消息内容不能为空")
		return
	}
	row, err := db.CreateChatMessage(c.Param("id"), "admin", req.MsgType, req.Content)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	chatws.BroadcastMessage(row)
	if sessionRow, err := db.GetChatSession(c.Param("id")); err == nil {
		chatws.BroadcastSession(sessionRow)
	}
	response.OK(c, gin.H{"item": row})
}

func CloseChatSession(c *gin.Context) {
	if err := db.CloseChatSession(c.Param("id")); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	if sessionRow, err := db.GetChatSession(c.Param("id")); err == nil {
		chatws.BroadcastSession(sessionRow)
	}
	response.OKMsg(c, "会话已关闭")
}

func ReadChatSession(c *gin.Context) {
	if err := db.ReadChatSession(c.Param("id")); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	if sessionRow, err := db.GetChatSession(c.Param("id")); err == nil {
		chatws.BroadcastSession(sessionRow)
	}
	response.OKMsg(c, "会话已读")
}

func uintParam(c *gin.Context) (uint, bool) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		response.Fail(c, "ID错误")
		return 0, false
	}
	return uint(id), true
}
