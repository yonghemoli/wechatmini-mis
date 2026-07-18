package mis

import (
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
		"banners":         jsonValue(db.MustAppConfigJSON(decorationBannersKey, `{"urls":[]}`), gin.H{"urls": []string{}}),
		"customerService": jsonValue(db.MustAppConfigJSON(decorationCustomerServiceKey, `{"name":"","phone":"","avatarUrl":""}`), gin.H{}),
		"company":         jsonValue(db.MustAppConfigJSON(decorationCompanyKey, `{"logoUrl":"","name":"永和护理","address":"","introduction":"","serviceGuarantees":[],"contactPhone":""}`), gin.H{}),
		"services":        items,
	})
}

func GetDecorationBanners(c *gin.Context) { decorationValue(c, decorationBannersKey, `{"urls":[]}`) }
func SaveDecorationBanners(c *gin.Context) {
	var req struct {
		URLs []string `json:"urls"`
	}
	if c.ShouldBindJSON(&req) != nil {
		response.Error(c, 400, "参数错误")
		return
	}
	if len(req.URLs) > 20 {
		response.Error(c, 400, "宣传图最多 20 张")
		return
	}
	for i := range req.URLs {
		req.URLs[i] = strings.TrimSpace(req.URLs[i])
		if req.URLs[i] == "" {
			response.Error(c, 400, "宣传图 URL 不能为空")
			return
		}
	}
	saveDecorationValue(c, decorationBannersKey, gin.H{"urls": req.URLs}, "宣传图配置")
}

func GetDecorationCustomerService(c *gin.Context) {
	decorationValue(c, decorationCustomerServiceKey, `{"name":"","phone":"","avatarUrl":""}`)
}
func SaveDecorationCustomerService(c *gin.Context) {
	var req struct {
		Name      string `json:"name"`
		Phone     string `json:"phone"`
		AvatarURL string `json:"avatarUrl"`
	}
	if c.ShouldBindJSON(&req) != nil || strings.TrimSpace(req.Name) == "" || strings.TrimSpace(req.Phone) == "" {
		response.Error(c, 400, "客服名字和电话不能为空")
		return
	}
	saveDecorationValue(c, decorationCustomerServiceKey, req, "客服信息配置")
}

func GetDecorationCompany(c *gin.Context) {
	decorationValue(c, decorationCompanyKey, `{"logoUrl":"","name":"永和护理","address":"","introduction":"","serviceGuarantees":[],"contactPhone":""}`)
}
func SaveDecorationCompany(c *gin.Context) {
	var req struct {
		LogoURL           string   `json:"logoUrl"`
		Name              string   `json:"name"`
		Address           string   `json:"address"`
		Introduction      string   `json:"introduction"`
		ServiceGuarantees []string `json:"serviceGuarantees"`
		ContactPhone      string   `json:"contactPhone"`
	}
	if c.ShouldBindJSON(&req) != nil || strings.TrimSpace(req.Name) == "" || strings.TrimSpace(req.ContactPhone) == "" {
		response.Error(c, 400, "公司名字和联系电话不能为空")
		return
	}
	saveDecorationValue(c, decorationCompanyKey, req, "公司信息配置")
}

func decorationValue(c *gin.Context, key, fallback string) {
	response.OK(c, jsonValue(db.MustAppConfigJSON(key, fallback), gin.H{}))
}
func saveDecorationValue(c *gin.Context, key string, value interface{}, note string) {
	if err := db.UpsertAppConfig(key, mustJSON(value, "{}"), note); err != nil {
		response.Error(c, 500, "保存装修配置失败")
		return
	}
	response.OK(c, value)
}
