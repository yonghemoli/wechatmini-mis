package mis

import (
	"encoding/json"
	"strconv"
	"strings"

	"yonghemolimis/src/apps/api/response"
	"yonghemolimis/src/dao/db"

	"github.com/gin-gonic/gin"
)

func ListCareServiceCategories(c *gin.Context) {
	rows, err := db.ListMiniServiceCategories(nil)
	if err != nil {
		response.Error(c, 500, "查询服务分类失败")
		return
	}
	items := make([]gin.H, 0, len(rows))
	for _, row := range rows {
		items = append(items, categoryAdminView(row))
	}
	response.OK(c, gin.H{"list": items})
}

func SaveCareServiceCategory(c *gin.Context) {
	id := strings.TrimSpace(c.Param("id"))
	var req struct {
		ID          string   `json:"id"`
		Name        string   `json:"name"`
		Subtitle    string   `json:"subtitle"`
		Description string   `json:"description"`
		IconURL     string   `json:"iconUrl"`
		Tags        []string `json:"tags"`
		Enabled     *bool    `json:"enabled"`
		Sort        int      `json:"sort"`
	}
	if c.ShouldBindJSON(&req) != nil {
		response.Error(c, 400, "参数错误")
		return
	}
	if id == "" {
		id = strings.TrimSpace(req.ID)
	}
	if id == "" || strings.TrimSpace(req.Name) == "" {
		response.Error(c, 400, "ID 和名称不能为空")
		return
	}
	enabled := true
	if req.Enabled != nil {
		enabled = *req.Enabled
	}
	row := &db.MiniServiceCategoryDO{ID: id, Name: strings.TrimSpace(req.Name), Subtitle: strings.TrimSpace(req.Subtitle), Description: strings.TrimSpace(req.Description), IconURL: strings.TrimSpace(req.IconURL), Tags: db.MarshalStringSlice(req.Tags), Enabled: enabled, Sort: req.Sort}
	if err := db.SaveMiniServiceCategory(row); err != nil {
		response.Error(c, 400, err.Error())
		return
	}
	response.OK(c, gin.H{"item": categoryAdminView(*row)})
}

func DeleteCareServiceCategory(c *gin.Context) {
	if err := db.DeleteMiniServiceCategory(c.Param("id")); err != nil {
		response.Error(c, 400, err.Error())
		return
	}
	response.OKMsg(c, "服务分类已删除")
}

func categoryAdminView(row db.MiniServiceCategoryDO) gin.H {
	return gin.H{"id": row.ID, "name": row.Name, "subtitle": row.Subtitle, "description": row.Description, "iconUrl": row.IconURL, "tags": db.UnmarshalStringSlice(row.Tags), "enabled": row.Enabled, "sort": row.Sort, "createdAt": row.CreatedAt, "updatedAt": row.UpdatedAt}
}

func ListCaregivers(c *gin.Context) {
	page, size := misPage(c)
	var published *bool
	if raw, exists := c.GetQuery("published"); exists {
		value, err := strconv.ParseBool(raw)
		if err != nil {
			response.Error(c, 400, "published 参数错误")
			return
		}
		published = &value
	}
	rows, total, err := db.ListCaregivers(db.CaregiverListQuery{ServiceID: c.Query("serviceId"), Keyword: c.Query("keyword"), AvailabilityStatus: c.Query("availabilityStatus"), Published: published, Page: page, PageSize: size})
	if err != nil {
		response.Error(c, 500, "查询服务人员失败")
		return
	}
	items := make([]gin.H, 0, len(rows))
	for _, row := range rows {
		items = append(items, caregiverAdminView(row))
	}
	response.OK(c, gin.H{"list": items, "total": total, "page": page, "pageSize": size})
}

func GetCaregiver(c *gin.Context) {
	row, err := db.GetCaregiver(c.Param("id"), false)
	if err != nil {
		response.Error(c, 404, "服务人员不存在")
		return
	}
	response.OK(c, gin.H{"item": caregiverAdminView(*row)})
}

func SaveCaregiver(c *gin.Context) {
	var req struct {
		ID                     string      `json:"id"`
		AvatarURL              string      `json:"avatarUrl"`
		Name                   string      `json:"name"`
		Age                    int         `json:"age"`
		ExperienceYears        int         `json:"experienceYears"`
		Origin                 string      `json:"origin"`
		ServiceIDs             []string    `json:"serviceIds"`
		Jobs                   []string    `json:"jobs"`
		AvailabilityStatus     string      `json:"availabilityStatus"`
		Rating                 float64     `json:"rating"`
		ServiceCount           int         `json:"serviceCount"`
		Recommended            bool        `json:"recommended"`
		Introduction           string      `json:"introduction"`
		Education              string      `json:"education"`
		Ethnicity              string      `json:"ethnicity"`
		Zodiac                 string      `json:"zodiac"`
		Skills                 []string    `json:"skills"`
		Certificates           interface{} `json:"certificates"`
		IdentityVerified       bool        `json:"identityVerified"`
		PhysicalExamVerified   bool        `json:"physicalExamVerified"`
		MedicalReportImageURLs []string    `json:"medicalReportImageUrls"`
		PersonalInfo           interface{} `json:"personalInfo"`
		WorkHistory            interface{} `json:"workHistory"`
		PhotoURLs              []string    `json:"photoUrls"`
		Published              bool        `json:"published"`
		Sort                   int         `json:"sort"`
	}
	if c.ShouldBindJSON(&req) != nil {
		response.Error(c, 400, "参数错误")
		return
	}
	if id := strings.TrimSpace(c.Param("id")); id != "" {
		req.ID = id
	}
	if strings.TrimSpace(req.Name) == "" || req.Age < 18 || req.Age > 70 || req.Rating < 0 || req.Rating > 5 {
		response.Error(c, 400, "姓名、年龄或评分不符合要求")
		return
	}
	for _, serviceID := range req.ServiceIDs {
		if _, err := db.GetMiniServiceCategory(serviceID, false); err != nil {
			response.Error(c, 400, "服务分类不存在: "+serviceID)
			return
		}
	}
	row := &db.CaregiverDO{ID: req.ID, AvatarURL: req.AvatarURL, Name: strings.TrimSpace(req.Name), Age: req.Age, ExperienceYears: req.ExperienceYears, Origin: req.Origin,
		ServiceIDs: db.MarshalStringSlice(req.ServiceIDs), Jobs: db.MarshalStringSlice(req.Jobs), AvailabilityStatus: req.AvailabilityStatus, Rating: req.Rating, ServiceCount: req.ServiceCount, Recommended: req.Recommended,
		Introduction: req.Introduction, Education: req.Education, Ethnicity: req.Ethnicity, Zodiac: req.Zodiac, Skills: db.MarshalStringSlice(req.Skills), Certificates: mustJSON(req.Certificates, "[]"),
		IdentityVerified: req.IdentityVerified, PhysicalExamVerified: req.PhysicalExamVerified, MedicalReportImageURLs: db.MarshalStringSlice(req.MedicalReportImageURLs), PersonalInfo: mustJSON(req.PersonalInfo, "{}"),
		WorkHistory: mustJSON(req.WorkHistory, "[]"), PhotoURLs: db.MarshalStringSlice(req.PhotoURLs), Published: req.Published, Sort: req.Sort}
	if !row.PhysicalExamVerified {
		row.MedicalReportImageURLs = "[]"
	}
	if err := db.SaveCaregiver(row); err != nil {
		response.Error(c, 400, err.Error())
		return
	}
	response.OK(c, gin.H{"item": caregiverAdminView(*row)})
}

func DeleteCaregiver(c *gin.Context) {
	if err := db.DeleteCaregiver(c.Param("id")); err != nil {
		response.Error(c, 400, err.Error())
		return
	}
	response.OKMsg(c, "服务人员已删除")
}

func caregiverAdminView(row db.CaregiverDO) gin.H {
	return gin.H{"id": row.ID, "avatarUrl": row.AvatarURL, "name": row.Name, "age": row.Age, "experienceYears": row.ExperienceYears, "origin": row.Origin, "serviceIds": db.UnmarshalStringSlice(row.ServiceIDs), "jobs": db.UnmarshalStringSlice(row.Jobs), "availabilityStatus": row.AvailabilityStatus, "rating": row.Rating, "serviceCount": row.ServiceCount, "recommended": row.Recommended, "introduction": row.Introduction, "education": row.Education, "ethnicity": row.Ethnicity, "zodiac": row.Zodiac, "skills": db.UnmarshalStringSlice(row.Skills), "certificates": jsonValue(row.Certificates, []interface{}{}), "identityVerified": row.IdentityVerified, "physicalExamVerified": row.PhysicalExamVerified, "medicalReportImageUrls": db.UnmarshalStringSlice(row.MedicalReportImageURLs), "personalInfo": jsonValue(row.PersonalInfo, map[string]interface{}{}), "workHistory": jsonValue(row.WorkHistory, []interface{}{}), "photoUrls": db.UnmarshalStringSlice(row.PhotoURLs), "published": row.Published, "sort": row.Sort, "createdAt": row.CreatedAt, "updatedAt": row.UpdatedAt}
}

func ListDemands(c *gin.Context) {
	page, size := misPage(c)
	rows, total, err := db.ListDemands(db.BusinessListQuery{Status: c.Query("status"), Keyword: c.Query("keyword"), Page: page, PageSize: size})
	if err != nil {
		response.Error(c, 500, "查询服务需求失败")
		return
	}
	response.OK(c, gin.H{"list": rows, "total": total, "page": page, "pageSize": size})
}
func ListResumes(c *gin.Context) {
	page, size := misPage(c)
	rows, total, err := db.ListResumes(db.BusinessListQuery{Status: c.Query("status"), Keyword: c.Query("keyword"), Page: page, PageSize: size})
	if err != nil {
		response.Error(c, 500, "查询简历失败")
		return
	}
	response.OK(c, gin.H{"list": rows, "total": total, "page": page, "pageSize": size})
}

func UpdateDemandStatus(c *gin.Context) { updateBusinessStatus(c, "DEMAND") }
func UpdateResumeStatus(c *gin.Context) { updateBusinessStatus(c, "RESUME") }
func AssignDemand(c *gin.Context)       { assignBusiness(c, "DEMAND") }
func AssignResume(c *gin.Context)       { assignBusiness(c, "RESUME") }
func assignBusiness(c *gin.Context, entityType string) {
	var req struct {
		AdminID uint `json:"adminId"`
	}
	if c.ShouldBindJSON(&req) != nil || req.AdminID == 0 {
		response.Error(c, 400, "顾问账号不能为空")
		return
	}
	if entityType == "DEMAND" {
		row, err := db.AssignDemand(c.Param("id"), req.AdminID)
		if err != nil {
			response.Error(c, 400, err.Error())
			return
		}
		response.OK(c, gin.H{"item": row})
		return
	}
	row, err := db.AssignResume(c.Param("id"), req.AdminID)
	if err != nil {
		response.Error(c, 400, err.Error())
		return
	}
	response.OK(c, gin.H{"item": row})
}
func updateBusinessStatus(c *gin.Context, entityType string) {
	var req struct {
		Status string `json:"status"`
		Note   string `json:"note"`
	}
	if c.ShouldBindJSON(&req) != nil || strings.TrimSpace(req.Status) == "" {
		response.Error(c, 400, "状态不能为空")
		return
	}
	operatorID, _ := c.Get("userID")
	adminID, _ := operatorID.(uint)
	if entityType == "DEMAND" {
		row, err := db.UpdateDemandStatus(c.Param("id"), req.Status, req.Note, adminID)
		if err != nil {
			response.Error(c, 400, err.Error())
			return
		}
		response.OK(c, gin.H{"item": row})
		return
	}
	row, err := db.UpdateResumeStatus(c.Param("id"), req.Status, req.Note, adminID)
	if err != nil {
		response.Error(c, 400, err.Error())
		return
	}
	response.OK(c, gin.H{"item": row})
}

func BusinessStatusHistory(c *gin.Context) {
	entityType := strings.ToUpper(c.Param("entityType"))
	if entityType != "DEMAND" && entityType != "RESUME" {
		response.Error(c, 400, "业务类型无效")
		return
	}
	rows, err := db.ListBusinessStatusHistory(entityType, c.Param("id"))
	if err != nil {
		response.Error(c, 500, "查询状态记录失败")
		return
	}
	response.OK(c, gin.H{"list": rows})
}

func GetBusinessContent(c *gin.Context) {
	key, ok := contentKey(c.Param("key"))
	if !ok {
		response.Error(c, 404, "内容配置不存在")
		return
	}
	row, err := db.GetAppConfig(key)
	if err != nil {
		response.Error(c, 404, "内容配置不存在")
		return
	}
	response.OK(c, gin.H{"key": c.Param("key"), "value": jsonValue(row.Value, map[string]interface{}{}), "updatedAt": row.UpdatedAt})
}
func SaveBusinessContent(c *gin.Context) {
	key, ok := contentKey(c.Param("key"))
	if !ok {
		response.Error(c, 404, "内容配置不存在")
		return
	}
	var req struct {
		Value interface{} `json:"value"`
	}
	if c.ShouldBindJSON(&req) != nil {
		response.Error(c, 400, "参数错误")
		return
	}
	if err := db.UpsertAppConfig(key, mustJSON(req.Value, "{}"), "护理小程序运营内容"); err != nil {
		response.Error(c, 400, err.Error())
		return
	}
	response.OKMsg(c, "内容已保存")
}
func contentKey(key string) (string, bool) {
	values := map[string]string{"app-config": "mini.business.app", "about": "mini.business.about", "agreement-privacy": "mini.business.agreement.privacy", "agreement-service": "mini.business.agreement.service"}
	value, ok := values[key]
	return value, ok
}

func misPage(c *gin.Context) (int, int) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))
	if page < 1 {
		page = 1
	}
	if size < 1 || size > 100 {
		size = 20
	}
	return page, size
}
func mustJSON(value interface{}, fallback string) string {
	if value == nil {
		return fallback
	}
	raw, err := json.Marshal(value)
	if err != nil {
		return fallback
	}
	return string(raw)
}
func jsonValue(raw string, fallback interface{}) interface{} {
	var value interface{}
	if json.Unmarshal([]byte(raw), &value) != nil {
		return fallback
	}
	return value
}
