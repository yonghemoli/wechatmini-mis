package db

import (
	"fmt"
	"strings"
)

// HomeBanner 是小程序首页轮播项。图片、文案和跳转行为必须整体保存。
type HomeBanner struct {
	ID          string `json:"id"`
	ImageURL    string `json:"imageUrl"`
	Kicker      string `json:"kicker"`
	Title       string `json:"title"`
	Description string `json:"description"`
	ActionType  string `json:"actionType"`
	ActionValue string `json:"actionValue"`
	Sort        int    `json:"sort"`
}

// NormalizeHomeBanners 兼容历史 {urls: []} 配置，并统一输出完整轮播项。
func NormalizeHomeBanners(raw interface{}) []HomeBanner {
	container, _ := raw.(map[string]interface{})
	items := interface{}(nil)
	if container != nil {
		items = container["items"]
		if items == nil {
			items = container["urls"]
		}
	} else {
		items = raw
	}
	values, ok := items.([]interface{})
	if !ok {
		return []HomeBanner{}
	}
	result := make([]HomeBanner, 0, len(values))
	for index, item := range values {
		banner := HomeBanner{}
		switch value := item.(type) {
		case string:
			banner.ImageURL = value
		case map[string]interface{}:
			banner.ID = decorationText(value["id"])
			banner.ImageURL = decorationText(value["imageUrl"])
			banner.Kicker = decorationText(value["kicker"])
			banner.Title = decorationText(value["title"])
			banner.Description = decorationText(value["description"])
			banner.ActionType = decorationText(value["actionType"])
			banner.ActionValue = decorationText(value["actionValue"])
			switch number := value["sort"].(type) {
			case float64:
				banner.Sort = int(number)
			case int:
				banner.Sort = number
			}
		}
		banner.ImageURL = strings.TrimSpace(banner.ImageURL)
		if banner.ImageURL == "" {
			continue
		}
		if banner.ID == "" {
			banner.ID = fmt.Sprintf("banner_%02d", index+1)
		}
		if banner.ActionType == "" {
			banner.ActionType = "NONE"
		}
		if banner.Sort == 0 {
			banner.Sort = (index + 1) * 10
		}
		result = append(result, banner)
	}
	return result
}

func decorationText(value interface{}) string {
	text, _ := value.(string)
	return strings.TrimSpace(text)
}

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
