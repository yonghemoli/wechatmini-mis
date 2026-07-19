package db

import "strings"

// ServiceGuarantee 是小程序公司资料中的一项服务保障。
// Icon 保存图标 URL 或前端约定的图标标识。
type ServiceGuarantee struct {
	Icon  string `json:"icon"`
	Title string `json:"title"`
	Sub   string `json:"sub"`
}

// NormalizeServiceGuarantees 兼容历史的字符串数组配置，并统一输出对象数组。
func NormalizeServiceGuarantees(raw interface{}) []ServiceGuarantee {
	items, ok := raw.([]interface{})
	if !ok {
		return []ServiceGuarantee{}
	}
	result := make([]ServiceGuarantee, 0, len(items))
	for _, item := range items {
		var guarantee ServiceGuarantee
		switch value := item.(type) {
		case string:
			guarantee.Title = strings.TrimSpace(value)
		case map[string]interface{}:
			guarantee.Icon, _ = value["icon"].(string)
			guarantee.Title, _ = value["title"].(string)
			guarantee.Sub, _ = value["sub"].(string)
		case map[string]string:
			guarantee = ServiceGuarantee{Icon: value["icon"], Title: value["title"], Sub: value["sub"]}
		}
		guarantee.Icon = strings.TrimSpace(guarantee.Icon)
		guarantee.Title = strings.TrimSpace(guarantee.Title)
		guarantee.Sub = strings.TrimSpace(guarantee.Sub)
		if guarantee.Title != "" {
			// 旧版仅保存标题。补齐默认值后，管理员可直接保存其他公司资料，
			// 不会被新版三字段校验阻断。
			if guarantee.Icon == "" {
				guarantee.Icon = "shield"
			}
			if guarantee.Sub == "" {
				guarantee.Sub = "服务保障"
			}
			result = append(result, guarantee)
		}
	}
	return result
}
