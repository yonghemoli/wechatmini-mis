package mis

import (
	"encoding/csv"
	"net/http"
	"time"

	"yonghemolimis/src/apps/api/response"
	"yonghemolimis/src/dao/db"

	"github.com/gin-gonic/gin"
)

func ListUsers(c *gin.Context) {
	rows, err := db.ListCustomers()
	if err != nil {
		response.Error(c, 500, "查询用户失败")
		return
	}
	response.OK(c, gin.H{"list": rows})
}

func BanUser(c *gin.Context)   { updateUserStatus(c, db.CustomerStatusBanned) }
func UnbanUser(c *gin.Context) { updateUserStatus(c, db.CustomerStatusActive) }
func updateUserStatus(c *gin.Context, status string) {
	row, err := db.UpdateCustomerStatus(c.Param("id"), status)
	if err != nil {
		response.Error(c, 400, err.Error())
		return
	}
	response.OK(c, gin.H{"item": row})
}

func ExportUsers(c *gin.Context) {
	rows, err := db.ListCustomers()
	if err != nil {
		response.Error(c, 500, "查询用户失败")
		return
	}
	writeCSV(c, "用户列表.csv", []string{"用户ID", "昵称", "手机号", "注册时间", "状态"}, func(w *csv.Writer) {
		for _, row := range rows {
			_ = w.Write([]string{row.ID, row.Nickname, row.Phone, formatTime(row.CreatedAt), row.Status})
		}
	})
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
