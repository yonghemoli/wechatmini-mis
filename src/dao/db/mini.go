package db

import (
	"encoding/json"
	"fmt"
	"time"
)

func GetMiniUserProfile(userID string) (*CustomerDO, error) {
	var row CustomerDO
	if err := Get().First(&row, "id = ?", userID).Error; err != nil {
		return nil, err
	}
	return &row, nil
}

func UpsertMiniUserProfile(row *CustomerDO) error {
	if row.ID == "" {
		return fmt.Errorf("user id required")
	}
	return Get().Save(row).Error
}

func TouchMiniUserLogin(userID string) error {
	return Get().Model(&CustomerDO{}).Where("id = ?", userID).
		Updates(map[string]any{"last_login_at": time.Now().Format("2006-01-02 15:04:05"), "updated_at": time.Now()}).Error
}

func GetAppConfig(key string) (*AppConfigDO, error) {
	var row AppConfigDO
	if err := Get().First(&row, "`key` = ?", key).Error; err != nil {
		return nil, err
	}
	return &row, nil
}

func MustAppConfigJSON(key string, fallback string) string {
	row, err := GetAppConfig(key)
	if err != nil || row == nil || row.Value == "" {
		return fallback
	}
	return row.Value
}

func UpsertAppConfig(key, value, note string) error {
	row := &AppConfigDO{Key: key, Value: value, Note: note, UpdatedAt: time.Now()}
	return Get().Save(row).Error
}

func ListMiniAddresses(userID string) ([]AddressDO, error) {
	var rows []AddressDO
	query := Get().Model(&AddressDO{})
	if userID != "" {
		query = query.Where("user_id = ?", userID)
	}
	err := query.Order("is_default DESC, id ASC").Find(&rows).Error
	return rows, err
}

func GetMiniAddress(userID, id string) (*AddressDO, error) {
	var row AddressDO
	if err := Get().First(&row, "id = ? AND user_id = ?", id, userID).Error; err != nil {
		return nil, err
	}
	return &row, nil
}

func CreateMiniAddress(row *AddressDO) error {
	if row.ID == "" {
		row.ID = fmt.Sprintf("addr_%d", time.Now().UnixNano())
	}
	return Get().Create(row).Error
}

func UpdateMiniAddress(userID string, row *AddressDO) error {
	if row.ID == "" {
		return fmt.Errorf("address id required")
	}
	return Get().Model(&AddressDO{}).Where("id = ? AND user_id = ?", row.ID, userID).
		Updates(map[string]any{
			"contact_name": row.ContactName,
			"phone":        row.Phone,
			"district":     row.District,
			"detail":       row.Detail,
			"tag":          row.Tag,
			"is_default":   row.IsDefault,
			"updated_at":   time.Now(),
		}).Error
}

func SetMiniAddressDefault(userID, id string) error {
	if err := Get().Model(&AddressDO{}).Where("user_id = ?", userID).
		Updates(map[string]any{"is_default": false, "updated_at": time.Now()}).Error; err != nil {
		return err
	}
	return Get().Model(&AddressDO{}).Where("id = ? AND user_id = ?", id, userID).
		Updates(map[string]any{"is_default": true, "updated_at": time.Now()}).Error
}

func DeleteMiniAddress(userID, id string) error {
	return Get().Where("id = ? AND user_id = ?", id, userID).Delete(&AddressDO{}).Error
}

func ListMiniServiceTargets(userID, category string) ([]ServiceTargetDO, error) {
	var rows []ServiceTargetDO
	query := Get().Model(&ServiceTargetDO{})
	if userID != "" {
		query = query.Where("user_id = ?", userID)
	}
	if category != "" {
		query = query.Where("category = ?", category)
	}
	err := query.Order("is_default DESC, id ASC").Find(&rows).Error
	return rows, err
}

func GetMiniServiceTarget(userID, id string) (*ServiceTargetDO, error) {
	var row ServiceTargetDO
	if err := Get().First(&row, "id = ? AND user_id = ?", id, userID).Error; err != nil {
		return nil, err
	}
	return &row, nil
}

func CreateMiniServiceTarget(row *ServiceTargetDO) error {
	if row.ID == "" {
		row.ID = fmt.Sprintf("target_%d", time.Now().UnixNano())
	}
	return Get().Create(row).Error
}

func UpdateMiniServiceTarget(userID string, row *ServiceTargetDO) error {
	if row.ID == "" {
		return fmt.Errorf("service target id required")
	}
	return Get().Model(&ServiceTargetDO{}).Where("id = ? AND user_id = ?", row.ID, userID).
		Updates(map[string]any{
			"name":       row.Name,
			"category":   row.Category,
			"relation":   row.Relation,
			"age":        row.Age,
			"note":       row.Note,
			"is_default": row.IsDefault,
			"updated_at": time.Now(),
		}).Error
}

func SetMiniServiceTargetDefault(userID, id string) error {
	var target ServiceTargetDO
	if err := Get().First(&target, "id = ? AND user_id = ?", id, userID).Error; err != nil {
		return err
	}
	if err := Get().Model(&ServiceTargetDO{}).Where("user_id = ? AND category = ?", userID, target.Category).
		Updates(map[string]any{"is_default": false, "updated_at": time.Now()}).Error; err != nil {
		return err
	}
	return Get().Model(&ServiceTargetDO{}).Where("id = ? AND user_id = ?", id, userID).
		Updates(map[string]any{"is_default": true, "updated_at": time.Now()}).Error
}

func DeleteMiniServiceTarget(userID, id string) error {
	return Get().Where("id = ? AND user_id = ?", id, userID).Delete(&ServiceTargetDO{}).Error
}

func ListMiniDishes() ([]DishDO, error) {
	var rows []DishDO
	err := Get().Order("id ASC").Find(&rows).Error
	return rows, err
}

func GetMiniDish(nameOrID string) (*DishDO, error) {
	var row DishDO
	if err := Get().First(&row, "id = ? OR name = ?", nameOrID, nameOrID).Error; err != nil {
		return nil, err
	}
	return &row, nil
}

func ListMiniMealPackages(userID, packageType string) ([]MealPackageDO, error) {
	var rows []MealPackageDO
	query := Get().Model(&MealPackageDO{})
	if userID != "" {
		query = query.Where("user_id = ?", userID)
	}
	if packageType != "" {
		query = query.Where("package_type = ?", packageType)
	}
	err := query.Order("package_type ASC, id ASC").Find(&rows).Error
	return rows, err
}

func GetMiniMealPackage(userID, id string) (*MealPackageDO, error) {
	var row MealPackageDO
	if err := Get().First(&row, "id = ? AND user_id = ? AND package_type = ?", id, userID, "custom").Error; err != nil {
		return nil, err
	}
	return &row, nil
}

func CreateMiniMealPackage(row *MealPackageDO) error {
	if row.ID == "" {
		row.ID = fmt.Sprintf("pkg_%d", time.Now().UnixNano())
	}
	return Get().Create(row).Error
}

func UpdateMiniMealPackage(userID string, row *MealPackageDO) error {
	if row.ID == "" {
		return fmt.Errorf("meal package id required")
	}
	return Get().Model(&MealPackageDO{}).Where("id = ? AND user_id = ? AND package_type = ?", row.ID, userID, "custom").
		Updates(map[string]any{
			"name":       row.Name,
			"scene":      row.Scene,
			"price":      row.Price,
			"dishes":     row.Dishes,
			"updated_at": time.Now(),
		}).Error
}

func DeleteMiniMealPackage(userID, id string) error {
	return Get().Where("id = ? AND user_id = ?", id, userID).Delete(&MealPackageDO{}).Error
}

func ListMiniOrders(userID, status string, page, size int) ([]OrderDO, int64, error) {
	query := Get().Model(&OrderDO{}).Where("mini_deleted = ?", false)
	if userID != "" {
		query = query.Where("(user_id = ? OR user_id = '')", userID)
	}
	if status != "" && status != "all" {
		query = query.Where("status = ?", status)
	}
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 20
	}
	var rows []OrderDO
	err := query.Order("created_at DESC").Limit(size).Offset((page - 1) * size).Find(&rows).Error
	return rows, total, err
}

func GetMiniOrder(userID, id string) (*OrderDO, error) {
	var row OrderDO
	query := Get().Where("id = ? AND mini_deleted = ?", id, false)
	if userID != "" {
		query = query.Where("(user_id = ? OR user_id = '')", userID)
	}
	if err := query.First(&row).Error; err != nil {
		return nil, err
	}
	return &row, nil
}

func DeleteMiniOrder(userID, id string) error {
	query := Get().Model(&OrderDO{}).Where("id = ?", id)
	if userID != "" {
		query = query.Where("(user_id = ? OR user_id = '')", userID)
	}
	return query.Updates(map[string]any{"mini_deleted": true, "updated_at": time.Now()}).Error
}

func MarshalStringSlice(values []string) string {
	if values == nil {
		return "[]"
	}
	data, err := json.Marshal(values)
	if err != nil {
		return "[]"
	}
	return string(data)
}

func UnmarshalStringSlice(value string) []string {
	if value == "" {
		return []string{}
	}
	var out []string
	if err := json.Unmarshal([]byte(value), &out); err != nil {
		return []string{}
	}
	return out
}
