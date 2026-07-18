package db

import "time"

const (
	CustomerStatusActive = "active"
	CustomerStatusBanned = "banned"
)

func ListCustomers() ([]CustomerDO, error) {
	var rows []CustomerDO
	err := Get().Order("created_at DESC").Find(&rows).Error
	return rows, err
}

func UpdateCustomerStatus(id, status string) (*CustomerDO, error) {
	if err := Get().Model(&CustomerDO{}).Where("id = ?", id).Updates(map[string]any{"status": status, "updated_at": time.Now()}).Error; err != nil {
		return nil, err
	}
	var row CustomerDO
	if err := Get().First(&row, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &row, nil
}
