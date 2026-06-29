package db

import (
	"strings"
	"time"

	"gorm.io/gorm"
)

const (
	OrderStatusPendingService = "pending_service"
	OrderStatusPendingConfirm = "pending_confirm"
	OrderStatusCancelled      = "cancelled"
	OrderStatusCompleted      = "completed"
	OrderStatusException      = "exception"
	OrderStatusRefunded       = "refunded"

	CustomerStatusActive = "active"
	CustomerStatusBanned = "banned"
)

type OrderListQuery struct {
	Status  string
	Keyword string
	Start   string
	End     string
	Page    int
	Size    int
}

func ListOrders(q OrderListQuery) ([]OrderDO, int64, error) {
	query := orderQuery(q)
	var total int64
	if err := query.Model(&OrderDO{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if q.Page <= 0 {
		q.Page = 1
	}
	if q.Size <= 0 {
		q.Size = 20
	}
	var rows []OrderDO
	err := query.Order("created_at DESC").
		Limit(q.Size).
		Offset((q.Page - 1) * q.Size).
		Find(&rows).Error
	return rows, total, err
}

func ListOrdersForExport(q OrderListQuery) ([]OrderDO, error) {
	var rows []OrderDO
	err := orderQuery(q).Order("created_at DESC").Find(&rows).Error
	return rows, err
}

func orderQuery(q OrderListQuery) *gorm.DB {
	query := Get().Model(&OrderDO{})
	if q.Status != "" && q.Status != "all" {
		query = query.Where("status = ?", q.Status)
	}
	if q.Keyword != "" {
		like := "%" + strings.TrimSpace(q.Keyword) + "%"
		query = query.Where(
			"id LIKE ? OR customer LIKE ? OR phone LIKE ? OR service LIKE ? OR staff LIKE ?",
			like, like, like, like, like,
		)
	}
	if q.Start != "" {
		query = query.Where("created_at >= ?", q.Start)
	}
	if q.End != "" {
		query = query.Where("created_at <= ?", q.End)
	}
	return query
}

func GetOrderByID(id string) (*OrderDO, error) {
	var row OrderDO
	if err := Get().First(&row, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &row, nil
}

func UpdateOrderStatus(id, status, reason string) (*OrderDO, error) {
	updates := map[string]any{
		"status":     status,
		"updated_at": time.Now(),
	}
	if reason != "" {
		updates["close_reason"] = reason
	}
	if err := Get().Model(&OrderDO{}).Where("id = ?", id).Updates(updates).Error; err != nil {
		return nil, err
	}
	return GetOrderByID(id)
}

func UpdateOrderNote(id, note string) (*OrderDO, error) {
	if err := Get().Model(&OrderDO{}).Where("id = ?", id).
		Updates(map[string]any{"internal_note": note, "updated_at": time.Now()}).Error; err != nil {
		return nil, err
	}
	return GetOrderByID(id)
}

func ListCustomers() ([]CustomerDO, error) {
	var rows []CustomerDO
	err := Get().Order("created_at DESC").Find(&rows).Error
	return rows, err
}

func UpdateCustomerStatus(id, status string) (*CustomerDO, error) {
	if err := Get().Model(&CustomerDO{}).Where("id = ?", id).
		Updates(map[string]any{"status": status, "updated_at": time.Now()}).Error; err != nil {
		return nil, err
	}
	var row CustomerDO
	if err := Get().First(&row, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &row, nil
}

func DashboardSummary() (map[string]any, error) {
	var todayOrders int64
	var todayRevenue int64
	var pending int64
	var users int64
	today := time.Now().Format("2006-01-02")
	err := Get().Model(&OrderDO{}).Where("DATE(created_at) = ?", today).Count(&todayOrders).Error
	if err != nil {
		return nil, err
	}
	err = Get().Model(&OrderDO{}).Where("DATE(created_at) = ? AND status <> ?", today, OrderStatusRefunded).
		Select("COALESCE(SUM(amount), 0)").Scan(&todayRevenue).Error
	if err != nil {
		return nil, err
	}
	err = Get().Model(&OrderDO{}).Where("status IN ?", []string{
		OrderStatusPendingService,
		OrderStatusPendingConfirm,
		OrderStatusException,
	}).Count(&pending).Error
	if err != nil {
		return nil, err
	}
	err = Get().Model(&CustomerDO{}).Count(&users).Error
	if err != nil {
		return nil, err
	}
	return map[string]any{
		"todayOrders":  todayOrders,
		"todayRevenue": todayRevenue,
		"pendingCount": pending,
		"userCount":    users,
	}, nil
}

func DashboardExceptions() ([]OrderDO, error) {
	var rows []OrderDO
	err := Get().Where("status IN ?", []string{
		OrderStatusPendingService,
		OrderStatusPendingConfirm,
		OrderStatusException,
	}).Order("created_at DESC").Limit(20).Find(&rows).Error
	return rows, err
}
