package mis

import (
	"encoding/csv"
	"net/http"
	"strconv"
	"time"

	"yonghemolimis/src/apps/api/response"
	"yonghemolimis/src/dao/db"

	"github.com/gin-gonic/gin"
)

func DashboardSummary(c *gin.Context) {
	data, err := db.DashboardSummary()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.OK(c, data)
}

func DashboardExceptions(c *gin.Context) {
	rows, err := db.DashboardExceptions()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.OK(c, gin.H{"list": rows})
}

func ListOrders(c *gin.Context) {
	q := orderQuery(c)
	rows, total, err := db.ListOrders(q)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.OKPage(c, rows, total, q.Page, q.Size)
}

func GetOrder(c *gin.Context) {
	row, err := db.GetOrderByID(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusNotFound, "订单不存在")
		return
	}
	response.OK(c, gin.H{"item": row})
}

func ConfirmOrder(c *gin.Context) {
	row, err := db.UpdateOrderStatus(c.Param("id"), db.OrderStatusCompleted, "")
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(c, gin.H{"item": row})
}

func RefundOrder(c *gin.Context) {
	var req struct {
		Reason string `json:"reason"`
	}
	_ = c.ShouldBindJSON(&req)
	row, err := db.UpdateOrderStatus(c.Param("id"), db.OrderStatusRefunded, req.Reason)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(c, gin.H{"item": row})
}

func UpdateOrderNote(c *gin.Context) {
	var req struct {
		InternalNote string `json:"internalNote"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, "参数错误")
		return
	}
	row, err := db.UpdateOrderNote(c.Param("id"), req.InternalNote)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(c, gin.H{"item": row})
}

func ExportOrders(c *gin.Context) {
	rows, err := db.ListOrdersForExport(orderQuery(c))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	writeCSV(c, "家政订单列表.csv", []string{"订单号", "客户", "手机", "服务", "金额", "状态", "来源", "预约时间", "创建时间", "服务人员", "内部备注"}, func(w *csv.Writer) {
		for _, row := range rows {
			_ = w.Write([]string{
				row.ID,
				row.Customer,
				row.Phone,
				row.Service,
				strconv.Itoa(row.Amount),
				row.Status,
				row.Source,
				row.AppointmentAt,
				formatTime(row.CreatedAt),
				row.Staff,
				row.InternalNote,
			})
		}
	})
}

func ListUsers(c *gin.Context) {
	rows, err := db.ListCustomers()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.OK(c, gin.H{"list": rows})
}

func BanUser(c *gin.Context) {
	updateUserStatus(c, db.CustomerStatusBanned)
}

func UnbanUser(c *gin.Context) {
	updateUserStatus(c, db.CustomerStatusActive)
}

func updateUserStatus(c *gin.Context, status string) {
	row, err := db.UpdateCustomerStatus(c.Param("id"), status)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(c, gin.H{"item": row})
}

func ExportUsers(c *gin.Context) {
	rows, err := db.ListCustomers()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	writeCSV(c, "家政用户列表.csv", []string{"用户ID", "昵称", "注册时间", "累计消费", "最后下单时间", "状态"}, func(w *csv.Writer) {
		for _, row := range rows {
			_ = w.Write([]string{
				row.ID,
				row.Nickname,
				formatTime(row.CreatedAt),
				strconv.Itoa(row.TotalSpent),
				row.LastOrderAt,
				row.Status,
			})
		}
	})
}

func orderQuery(c *gin.Context) db.OrderListQuery {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))
	return db.OrderListQuery{
		Status:  c.Query("status"),
		Keyword: c.Query("keyword"),
		Start:   c.Query("start"),
		End:     c.Query("end"),
		Page:    page,
		Size:    size,
	}
}

func writeCSV(c *gin.Context, filename string, headers []string, writeRows func(*csv.Writer)) {
	c.Header("Content-Type", "text/csv; charset=utf-8")
	c.Header("Content-Disposition", `attachment; filename="`+filename+`"`)
	c.Status(http.StatusOK)
	_, _ = c.Writer.Write([]byte{0xEF, 0xBB, 0xBF})
	writer := csv.NewWriter(c.Writer)
	_ = writer.Write(headers)
	writeRows(writer)
	writer.Flush()
}

func formatTime(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format("2006-01-02 15:04")
}
