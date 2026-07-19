package mis

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"yonghemolimis/src/apps/api/response"
	"yonghemolimis/src/dao/db"

	"github.com/gin-gonic/gin"
)

const (
	decorationBannersKey         = "mini.decoration.banners"
	decorationCustomerServiceKey = "mini.decoration.customer_service"
	decorationCompanyKey         = "mini.decoration.company"
)

func GetDecoration(c *gin.Context) {
	services, err := db.ListMiniServiceCategories(nil)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "读取服务项目失败")
		return
	}
	items := make([]gin.H, 0, len(services))
	for _, item := range services {
		items = append(items, categoryAdminView(item))
	}
	response.OK(c, gin.H{
		"banners":         decorationBannersValue(),
		"customerService": decorationCustomerServiceValue(),
		"company":         decorationCompanyValue(),
		"services":        items,
	})
}

func GetDecorationBanners(c *gin.Context) { response.OK(c, decorationBannersValue()) }
func SaveDecorationBanners(c *gin.Context) {
	var req struct {
		Items []db.HomeBanner `json:"items"`
	}
	if c.ShouldBindJSON(&req) != nil {
		response.Error(c, 400, "参数错误")
		return
	}
	if req.Items == nil {
		req.Items = []db.HomeBanner{}
	}
	if len(req.Items) > 20 {
		response.Error(c, 400, "宣传图最多 20 张")
		return
	}
	ids := make(map[string]struct{}, len(req.Items))
	for i := range req.Items {
		item := &req.Items[i]
		item.ID = strings.TrimSpace(item.ID)
		item.ImageURL = strings.TrimSpace(item.ImageURL)
		item.Kicker = strings.TrimSpace(item.Kicker)
		item.Title = strings.TrimSpace(item.Title)
		item.Description = strings.TrimSpace(item.Description)
		item.ActionType = strings.ToUpper(strings.TrimSpace(item.ActionType))
		item.ActionValue = strings.TrimSpace(item.ActionValue)
		if item.ID == "" {
			item.ID = fmt.Sprintf("banner_%02d", i+1)
		}
		if item.ImageURL == "" || item.Title == "" {
			response.Error(c, 400, "轮播图图片和标题不能为空")
			return
		}
		if item.ActionType == "" {
			item.ActionType = "NONE"
		}
		if item.ActionType != "NONE" && item.ActionType != "DEMAND" && item.ActionType != "CAREGIVER_LIST" {
			response.Error(c, 400, "轮播图跳转类型无效")
			return
		}
		if item.ActionType == "DEMAND" && item.ActionValue == "" {
			response.Error(c, 400, "预约跳转需要填写服务项目 ID")
			return
		}
		if _, exists := ids[item.ID]; exists {
			response.Error(c, 400, "轮播图 ID 不能重复")
			return
		}
		ids[item.ID] = struct{}{}
		if item.Sort == 0 {
			item.Sort = (i + 1) * 10
		}
	}
	saveDecorationValue(c, decorationBannersKey, gin.H{"items": req.Items}, "首页轮播配置")
}

func GetDecorationCustomerService(c *gin.Context) { response.OK(c, decorationCustomerServiceValue()) }
func SaveDecorationCustomerService(c *gin.Context) {
	var req struct {
		Name      string `json:"name"`
		Phone     string `json:"phone"`
		AvatarURL string `json:"avatarUrl"`
	}
	if c.ShouldBindJSON(&req) != nil {
		response.Error(c, 400, "参数错误")
		return
	}
	req.Name = strings.TrimSpace(req.Name)
	req.Phone = strings.TrimSpace(req.Phone)
	req.AvatarURL = strings.TrimSpace(req.AvatarURL)
	if req.Name == "" || req.Phone == "" {
		response.Error(c, 400, "客服名字和电话不能为空")
		return
	}
	saveDecorationValue(c, decorationCustomerServiceKey, req, "客服信息配置")
}

func GetDecorationCompany(c *gin.Context) {
	response.OK(c, decorationCompanyValue())
}
func SaveDecorationCompany(c *gin.Context) {
	var req struct {
		LogoURL           string                `json:"logoUrl"`
		Name              string                `json:"name"`
		Address           string                `json:"address"`
		Introduction      string                `json:"introduction"`
		ServiceGuarantees []db.ServiceGuarantee `json:"serviceGuarantees"`
		ContactPhone      string                `json:"contactPhone"`
	}
	if c.ShouldBindJSON(&req) != nil {
		response.Error(c, 400, "参数错误")
		return
	}
	if req.ServiceGuarantees == nil {
		req.ServiceGuarantees = []db.ServiceGuarantee{}
	}
	req.LogoURL = strings.TrimSpace(req.LogoURL)
	req.Name = strings.TrimSpace(req.Name)
	req.Address = strings.TrimSpace(req.Address)
	req.Introduction = strings.TrimSpace(req.Introduction)
	req.ContactPhone = strings.TrimSpace(req.ContactPhone)
	if req.Name == "" || req.ContactPhone == "" {
		response.Error(c, 400, "公司名字和联系电话不能为空")
		return
	}
	if len(req.ServiceGuarantees) > 12 {
		response.Error(c, 400, "服务保障最多 12 项")
		return
	}
	for i := range req.ServiceGuarantees {
		req.ServiceGuarantees[i].Icon = strings.TrimSpace(req.ServiceGuarantees[i].Icon)
		req.ServiceGuarantees[i].Title = strings.TrimSpace(req.ServiceGuarantees[i].Title)
		req.ServiceGuarantees[i].Sub = strings.TrimSpace(req.ServiceGuarantees[i].Sub)
		if req.ServiceGuarantees[i].Icon == "" || req.ServiceGuarantees[i].Title == "" || req.ServiceGuarantees[i].Sub == "" {
			response.Error(c, 400, "服务保障图标、标题和副标题不能为空")
			return
		}
	}
	saveDecorationValue(c, decorationCompanyKey, req, "公司信息配置")
}

func decorationCompanyValue() gin.H {
	return companyDecorationValue(db.MustAppConfigJSON(decorationCompanyKey, `{"logoUrl":"","name":"永和护理","address":"","introduction":"","serviceGuarantees":[],"contactPhone":""}`))
}

func companyDecorationValue(raw string) gin.H {
	value := decorationObjectFromJSON(raw, `{"logoUrl":"","name":"永和护理","address":"","introduction":"","serviceGuarantees":[],"contactPhone":""}`)
	return gin.H{
		"logoUrl":           decorationString(value["logoUrl"]),
		"name":              decorationStringDefault(value["name"], "永和护理"),
		"address":           decorationString(value["address"]),
		"introduction":      decorationString(value["introduction"]),
		"serviceGuarantees": db.NormalizeServiceGuarantees(value["serviceGuarantees"]),
		"contactPhone":      decorationString(value["contactPhone"]),
	}
}

func decorationBannersValue() gin.H {
	value := decorationObject(decorationBannersKey, `{"items":[]}`)
	return gin.H{"items": db.NormalizeHomeBanners(map[string]interface{}(value))}
}

func decorationCustomerServiceValue() gin.H {
	value := decorationObject(decorationCustomerServiceKey, `{"name":"","phone":"","avatarUrl":""}`)
	return gin.H{
		"name":      decorationString(value["name"]),
		"phone":     decorationString(value["phone"]),
		"avatarUrl": decorationString(value["avatarUrl"]),
	}
}

func decorationObject(key, fallback string) gin.H {
	return decorationObjectFromJSON(db.MustAppConfigJSON(key, fallback), fallback)
}

func decorationObjectFromJSON(raw, fallback string) gin.H {
	value := map[string]interface{}{}
	if json.Unmarshal([]byte(raw), &value) != nil {
		_ = json.Unmarshal([]byte(fallback), &value)
	}
	return gin.H(value)
}

func decorationString(value interface{}) string {
	text, _ := value.(string)
	return strings.TrimSpace(text)
}

func decorationStringDefault(value interface{}, fallback string) string {
	if text := decorationString(value); text != "" {
		return text
	}
	return fallback
}

func decorationStringSlice(value interface{}) []string {
	items, ok := value.([]interface{})
	if !ok {
		return []string{}
	}
	result := make([]string, 0, len(items))
	for _, item := range items {
		if text := decorationString(item); text != "" {
			result = append(result, text)
		}
	}
	return result
}
func saveDecorationValue(c *gin.Context, key string, value interface{}, note string) {
	if err := db.UpsertAppConfig(key, mustJSON(value, "{}"), note); err != nil {
		response.Error(c, 500, "保存装修配置失败")
		return
	}
	response.OK(c, value)
}
