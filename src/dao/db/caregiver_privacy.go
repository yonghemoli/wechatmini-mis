package db

import (
	"encoding/json"
	"regexp"
	"strings"
)

var mainlandPhoneRegexp = regexp.MustCompile(`1[3-9]\d{9}`)

// CaregiverPersonalInfoKeys 是允许进入小程序阿姨简历的个人信息白名单。
var CaregiverPersonalInfoKeys = map[string]bool{
	"heightCm": true, "weightKg": true, "bloodType": true, "gender": true,
	"maritalStatus": true, "religion": true, "languages": true, "liveInAvailable": true,
}

func ContainsMainlandPhone(value string) bool { return mainlandPhoneRegexp.MatchString(value) }

// ContainsPhoneInPublicContent 检查将公开展示的文本内容；媒体 URL 不参与检查。
func ContainsPhoneInPublicContent(value interface{}) bool {
	switch typed := value.(type) {
	case string:
		return ContainsMainlandPhone(typed)
	case []string:
		for _, item := range typed {
			if ContainsMainlandPhone(item) {
				return true
			}
		}
	case []interface{}:
		for _, item := range typed {
			if ContainsPhoneInPublicContent(item) {
				return true
			}
		}
	case map[string]interface{}:
		for key, item := range typed {
			if isCaregiverMediaKey(key) {
				continue
			}
			if ContainsPhoneInPublicContent(item) {
				return true
			}
		}
	}
	return false
}

// SanitizeCaregiverPersonalInfo 移除私密键和含手机号的历史数据，仅保留公开资料白名单。
func SanitizeCaregiverPersonalInfo(value interface{}) map[string]interface{} {
	object := toJSONObject(value)
	result := make(map[string]interface{})
	for key := range CaregiverPersonalInfoKeys {
		item, ok := object[key]
		if !ok || ContainsPhoneInPublicContent(item) {
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

func isCaregiverMediaKey(key string) bool {
	switch strings.ToLower(key) {
	case "url", "fileurl", "imageurl", "imageurls", "avatarurl", "photourls", "medicalreportimageurls":
		return true
	default:
		return false
	}
}
