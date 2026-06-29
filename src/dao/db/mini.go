package db

import (
	"encoding/json"
	"fmt"
	"time"
)

func ListMiniAddresses(userID string) ([]AddressDO, error) {
	var rows []AddressDO
	query := Get().Model(&AddressDO{})
	if userID != "" {
		query = query.Where("user_id = ?", userID)
	}
	err := query.Order("is_default DESC, id ASC").Find(&rows).Error
	return rows, err
}

func CreateMiniAddress(row *AddressDO) error {
	if row.ID == "" {
		row.ID = fmt.Sprintf("addr_%d", time.Now().UnixNano())
	}
	return Get().Create(row).Error
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

func CreateMiniServiceTarget(row *ServiceTargetDO) error {
	if row.ID == "" {
		row.ID = fmt.Sprintf("target_%d", time.Now().UnixNano())
	}
	return Get().Create(row).Error
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

func CreateMiniMealPackage(row *MealPackageDO) error {
	if row.ID == "" {
		row.ID = fmt.Sprintf("pkg_%d", time.Now().UnixNano())
	}
	return Get().Create(row).Error
}

func DeleteMiniMealPackage(userID, id string) error {
	return Get().Where("id = ? AND user_id = ?", id, userID).Delete(&MealPackageDO{}).Error
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
