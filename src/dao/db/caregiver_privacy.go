package db

import (
	"encoding/json"
	"strings"
)

// SanitizeCaregiverPersonalInfo 仅移除阿姨联系方式；其余人工维护的简历资料原样保留。
func SanitizeCaregiverPersonalInfo(value interface{}) map[string]interface{} {
	object := toJSONObject(value)
	result := make(map[string]interface{})
	for key, item := range object {
		if isCaregiverPhoneKey(key) {
			continue
		}
		result[key] = item
	}
	return result
}

func toJSONObject(value interface{}) map[string]interface{} {
	if object, ok := value.(map[string]interface{}); ok {
		return object
	}
	if raw, ok := value.(string); ok && strings.TrimSpace(raw) != "" {
		object := map[string]interface{}{}
		if json.Unmarshal([]byte(raw), &object) == nil {
			return object
		}
	}
	return map[string]interface{}{}
}

func isCaregiverPhoneKey(key string) bool {
	switch strings.ToLower(key) {
	case "contactphone", "contact_phone", "phone", "mobile", "mobilephone", "tel", "telephone":
		return true
	default:
		return false
	}
}
