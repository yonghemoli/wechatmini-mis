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

func GetMiniUserByOpenID(openID string) (*CustomerDO, error) {
	var row CustomerDO
	if err := Get().First(&row, "openid = ?", openID).Error; err != nil {
		return nil, err
	}
	return &row, nil
}

func GetMiniUserByPhone(phone string) (*CustomerDO, error) {
	var row CustomerDO
	if err := Get().First(&row, "phone = ?", phone).Error; err != nil {
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

func GetAppConfig(key string) (*AppConfigDO, error) {
	var row AppConfigDO
	if err := Get().First(&row, "`key` = ?", key).Error; err != nil {
		return nil, err
	}
	return &row, nil
}

func MustAppConfigJSON(key, fallback string) string {
	row, err := GetAppConfig(key)
	if err != nil || row.Value == "" {
		return fallback
	}
	return row.Value
}

func UpsertAppConfig(key, value, note string) error {
	row := &AppConfigDO{Key: key, Value: value, Note: note, UpdatedAt: time.Now()}
	return Get().Save(row).Error
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
	if json.Unmarshal([]byte(value), &out) != nil {
		return []string{}
	}
	return out
}
