package miniapi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"yonghemolimis/src/dao/db"

	"github.com/gin-gonic/gin"
)

var mainlandPhonePattern = regexp.MustCompile(`^1[3-9]\d{9}$`)

var availabilityText = map[string]string{
	"AVAILABLE_NOW": "随时上岗", "AVAILABLE_IN_3_DAYS": "三天内上岗",
	"AVAILABLE_IN_1_WEEK": "一周内上岗", "OPEN_TO_OPPORTUNITIES": "在职看机会",
	"UNAVAILABLE": "暂不可上岗",
}

var demandSources = stringSet("HOME_BANNER", "HOME_SERVICE", "SERVICE_LIST", "CAREGIVER_DETAIL", "OTHER")
var experienceRanges = stringSet("LESS_THAN_1_YEAR", "YEAR_1_TO_3", "YEAR_3_TO_5", "YEAR_5_TO_10", "MORE_THAN_10_YEARS")

type businessError struct {
	Code      int         `json:"code"`
	Message   string      `json:"message"`
	Data      interface{} `json:"data"`
	RequestID string      `json:"requestId"`
}

func businessFail(c *gin.Context, httpCode, code int, message string) {
	c.JSON(httpCode, businessError{Code: code, Message: message, Data: nil, RequestID: requestID(c)})
}

func requestID(c *gin.Context) string {
	id := strings.TrimSpace(c.GetHeader("X-Request-Id"))
	if id == "" {
		id = "req_" + randomHex(12)
	}
	c.Header("X-Request-Id", id)
	return id
}

func businessOK(c *gin.Context, status int, data interface{}) {
	requestID(c)
	c.JSON(status, R{Code: 0, Message: "ok", Data: data})
}

type serviceCategoryView struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Subtitle    string   `json:"subtitle"`
	Description string   `json:"description"`
	IconURL     string   `json:"iconUrl"`
	Tags        []string `json:"tags"`
	Enabled     bool     `json:"enabled"`
	Sort        int      `json:"sort"`
}

func ListBusinessServices(c *gin.Context) {
	enabled := true
	if raw, exists := c.GetQuery("enabled"); exists {
		parsed, err := strconv.ParseBool(raw)
		if err != nil {
			businessFail(c, http.StatusBadRequest, 40000, "enabled 参数错误")
			return
		}
		enabled = parsed
	}
	var enabledFilter *bool
	if enabled {
		enabledFilter = &enabled
	}
	rows, err := db.ListMiniServiceCategories(enabledFilter)
	if err != nil {
		businessFail(c, http.StatusInternalServerError, 50000, "查询服务分类失败")
		return
	}
	list := make([]serviceCategoryView, 0, len(rows))
	for _, row := range rows {
		list = append(list, categoryView(row))
	}
	businessOK(c, http.StatusOK, gin.H{"list": list})
}

func categoryView(row db.MiniServiceCategoryDO) serviceCategoryView {
	return serviceCategoryView{ID: row.ID, Name: row.Name, Subtitle: row.Subtitle, Description: row.Description,
		IconURL: row.IconURL, Tags: db.UnmarshalStringSlice(row.Tags), Enabled: row.Enabled, Sort: row.Sort}
}

type caregiverSummary struct {
	ID                 string          `json:"id"`
	AvatarURL          string          `json:"avatarUrl"`
	Name               string          `json:"name"`
	Age                int             `json:"age"`
	ExperienceYears    int             `json:"experienceYears"`
	ExperienceText     string          `json:"experienceText"`
	Origin             string          `json:"origin"`
	ServiceIDs         []string        `json:"serviceIds"`
	Jobs               []string        `json:"jobs"`
	AvailabilityStatus string          `json:"availabilityStatus"`
	AvailabilityText   string          `json:"availabilityText"`
	Recommended        bool            `json:"recommended"`
	DisplayFields      map[string]bool `json:"displayFields"`
}

type caregiverDetail struct {
	caregiverSummary
	Introduction           string          `json:"introduction"`
	Education              string          `json:"education"`
	Ethnicity              string          `json:"ethnicity"`
	Zodiac                 string          `json:"zodiac,omitempty"`
	BirthDate              string          `json:"birthDate,omitempty"`
	Constellation          string          `json:"constellation,omitempty"`
	Skills                 []string        `json:"skills"`
	Certificates           interface{}     `json:"certificates"`
	IdentityVerified       bool            `json:"identityVerified"`
	PhysicalExamVerified   bool            `json:"physicalExamVerified"`
	MedicalReportImageURLs []string        `json:"medicalReportImageUrls"`
	PersonalInfo           interface{}     `json:"personalInfo"`
	WorkHistory            interface{}     `json:"workHistory"`
	PhotoURLs              []string        `json:"photoUrls"`
	DisplayFields          map[string]bool `json:"displayFields"`
}

func ListBusinessCaregivers(c *gin.Context) {
	page, err1 := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, err2 := strconv.Atoi(c.DefaultQuery("pageSize", "20"))
	if err1 != nil || err2 != nil || page < 1 || pageSize < 1 || pageSize > 50 {
		businessFail(c, http.StatusBadRequest, 40000, "分页参数错误")
		return
	}
	keyword := strings.TrimSpace(c.Query("keyword"))
	if utf8.RuneCountInString(keyword) > 30 {
		businessFail(c, http.StatusBadRequest, 40000, "keyword 最长 30 个字符")
		return
	}
	status := strings.TrimSpace(c.Query("availabilityStatus"))
	if status != "" {
		if _, ok := availabilityText[status]; !ok {
			businessFail(c, http.StatusBadRequest, 40000, "上岗状态无效")
			return
		}
	}
	var recommended *bool
	if raw, exists := c.GetQuery("recommended"); exists {
		value, err := strconv.ParseBool(raw)
		if err != nil {
			businessFail(c, 400, 40000, "recommended 参数错误")
			return
		}
		recommended = &value
	}
	rows, total, err := db.ListCaregivers(db.CaregiverListQuery{ServiceID: strings.TrimSpace(c.Query("serviceId")), Keyword: keyword,
		AvailabilityStatus: status, Recommended: recommended, Status: "COMPLETED", Page: page, PageSize: pageSize})
	if err != nil {
		businessFail(c, 500, 50000, "查询服务人员失败")
		return
	}
	list := make([]caregiverSummary, 0, len(rows))
	for _, row := range rows {
		list = append(list, caregiverSummaryView(row))
	}
	businessOK(c, 200, gin.H{"list": list, "page": page, "pageSize": pageSize, "total": total, "hasMore": int64(page*pageSize) < total})
}

func caregiverSummaryView(row db.CaregiverDO) caregiverSummary {
	return caregiverSummary{ID: row.ID, AvatarURL: row.AvatarURL, Name: row.Name, Age: row.Age, ExperienceYears: row.ExperienceYears,
		ExperienceText: fmt.Sprintf("%d年经验", row.ExperienceYears), Origin: row.Origin, ServiceIDs: db.UnmarshalStringSlice(row.ServiceIDs),
		Jobs: db.UnmarshalStringSlice(row.Jobs), AvailabilityStatus: row.AvailabilityStatus, AvailabilityText: availabilityText[row.AvailabilityStatus],
		Recommended: row.Recommended, DisplayFields: displayFields(row.DisplayFields)}
}

func GetBusinessCaregiver(c *gin.Context) {
	row, err := db.GetCaregiver(c.Param("id"), true)
	if err != nil {
		businessFail(c, 404, 40400, "服务人员不存在")
		return
	}
	reports := db.UnmarshalStringSlice(row.MedicalReportImageURLs)
	if !row.PhysicalExamVerified {
		reports = []string{}
	}
	detail := caregiverDetail{caregiverSummary: caregiverSummaryView(*row), Introduction: row.Introduction, Education: row.Education,
		Ethnicity: row.Ethnicity, Zodiac: row.Zodiac, BirthDate: row.BirthDate, Constellation: row.Constellation, Skills: db.UnmarshalStringSlice(row.Skills), Certificates: db.NormalizeJSONObjectList(decodeJSONValue(row.Certificates, []interface{}{}), "name"),
		IdentityVerified: row.IdentityVerified, PhysicalExamVerified: row.PhysicalExamVerified, MedicalReportImageURLs: reports,
		PersonalInfo: decodeJSONValue(row.PersonalInfo, map[string]interface{}{}), WorkHistory: db.NormalizeJSONObjectList(decodeJSONValue(row.WorkHistory, []interface{}{}), "role"), PhotoURLs: db.UnmarshalStringSlice(row.PhotoURLs), DisplayFields: displayFields(row.DisplayFields)}
	businessOK(c, 200, detail)
}

func CreateBusinessDemand(c *gin.Context) {
	var req struct {
		ServiceID    string `json:"serviceId"`
		CaregiverID  string `json:"caregiverId"`
		ContactName  string `json:"contactName"`
		ContactPhone string `json:"contactPhone"`
		Requirements string `json:"requirements"`
		Source       string `json:"source"`
	}
	if c.ShouldBindJSON(&req) != nil {
		businessFail(c, 400, 40000, "请求参数错误")
		return
	}
	req.ServiceID = strings.TrimSpace(req.ServiceID)
	req.CaregiverID = strings.TrimSpace(req.CaregiverID)
	req.ContactName = strings.TrimSpace(req.ContactName)
	req.ContactPhone = strings.TrimSpace(req.ContactPhone)
	req.Requirements = strings.TrimSpace(req.Requirements)
	req.Source = strings.TrimSpace(req.Source)
	service, err := db.GetMiniServiceCategory(req.ServiceID, true)
	if err != nil {
		businessFail(c, 400, 40000, "服务分类不存在或已停用")
		return
	}
	if !mainlandPhonePattern.MatchString(req.ContactPhone) {
		businessFail(c, 400, 40001, "手机号格式不正确")
		return
	}
	if utf8.RuneCountInString(req.Requirements) < 1 || utf8.RuneCountInString(req.Requirements) > 500 {
		businessFail(c, 400, 40000, "需求描述长度须为 1—500 字")
		return
	}
	if !demandSources[req.Source] {
		businessFail(c, 400, 40000, "来源类型无效")
		return
	}
	caregiverName := ""
	if req.CaregiverID != "" {
		caregiver, err := db.GetCaregiver(req.CaregiverID, true)
		if err != nil {
			businessFail(c, 400, 40000, "服务人员不存在或未发布")
			return
		}
		if !contains(db.UnmarshalStringSlice(caregiver.ServiceIDs), req.ServiceID) {
			businessFail(c, 400, 40000, "该服务人员不能承接所选服务")
			return
		}
		caregiverName = caregiver.Name
	}
	user, userID := requestUser(c)
	if req.ContactName == "" {
		req.ContactName = strings.TrimSpace(user.NickName)
	}
	if utf8.RuneCountInString(req.ContactName) < 1 || utf8.RuneCountInString(req.ContactName) > 64 {
		businessFail(c, 400, 40000, "联系人姓名长度须为 1—64 字")
		return
	}
	scope := submissionScope(userID, req.ContactPhone)
	key := strings.TrimSpace(c.GetHeader("Idempotency-Key"))
	if len(key) > 128 {
		businessFail(c, 400, 40000, "Idempotency-Key 最长 128 个字符")
		return
	}
	if old, err := db.FindDemandByIdempotency(scope, key, time.Now().Add(-10*time.Minute)); err == nil {
		businessOK(c, 201, createdView(old.ID, old.Status, old.CreatedAt))
		return
	}
	if _, err := db.FindRecentDemandDuplicate(req.ContactPhone, req.ServiceID, req.Requirements, time.Now().Add(-10*time.Minute)); err == nil {
		businessFail(c, 409, 40900, "请勿重复提交相同需求")
		return
	}
	now := time.Now()
	row := &db.DemandDO{ID: businessID("D", now), UserID: userID, ServiceID: service.ID, ServiceName: service.Name, CaregiverID: req.CaregiverID,
		CaregiverName: caregiverName, ContactName: req.ContactName, ContactPhone: req.ContactPhone, Requirements: req.Requirements, Source: req.Source, Status: "PENDING_CONTACT", IdempotencyKey: key, SubmissionScope: scope, CreatedAt: now, UpdatedAt: now}
	if err := db.CreateDemand(row); err != nil {
		businessFail(c, 500, 50000, "提交服务需求失败")
		return
	}
	businessOK(c, 201, createdView(row.ID, row.Status, row.CreatedAt))
}

func CreateBusinessResume(c *gin.Context) {
	var req struct {
		IntentionServiceID string `json:"intentionServiceId"`
		WorkStatus         string `json:"workStatus"`
		ExperienceRange    string `json:"experienceRange"`
		EntryYear          int    `json:"entryYear"`
		ContactName        string `json:"contactName"`
		ContactPhone       string `json:"contactPhone"`
	}
	if c.ShouldBindJSON(&req) != nil {
		businessFail(c, 400, 40000, "请求参数错误")
		return
	}
	req.IntentionServiceID = strings.TrimSpace(req.IntentionServiceID)
	req.WorkStatus = strings.TrimSpace(req.WorkStatus)
	req.ExperienceRange = strings.TrimSpace(req.ExperienceRange)
	req.ContactName = strings.TrimSpace(req.ContactName)
	req.ContactPhone = strings.TrimSpace(req.ContactPhone)
	service, err := db.GetMiniServiceCategory(req.IntentionServiceID, true)
	if err != nil {
		businessFail(c, 400, 40000, "服务分类不存在或已停用")
		return
	}
	if _, ok := availabilityText[req.WorkStatus]; !ok {
		businessFail(c, 400, 40000, "上岗状态无效")
		return
	}
	if !experienceRanges[req.ExperienceRange] {
		businessFail(c, 400, 40000, "从业年限无效")
		return
	}
	if req.EntryYear < 2000 || req.EntryYear > time.Now().Year() || !experienceMatches(time.Now().Year()-req.EntryYear, req.ExperienceRange) {
		businessFail(c, 400, 40000, "入行年份与从业年限不一致")
		return
	}
	if !mainlandPhonePattern.MatchString(req.ContactPhone) {
		businessFail(c, 400, 40001, "手机号格式不正确")
		return
	}
	user, userID := requestUser(c)
	if req.ContactName == "" {
		req.ContactName = strings.TrimSpace(user.NickName)
	}
	if utf8.RuneCountInString(req.ContactName) < 1 || utf8.RuneCountInString(req.ContactName) > 64 {
		businessFail(c, 400, 40000, "联系人姓名长度须为 1—64 字")
		return
	}
	scope := submissionScope(userID, req.ContactPhone)
	key := strings.TrimSpace(c.GetHeader("Idempotency-Key"))
	if len(key) > 128 {
		businessFail(c, 400, 40000, "Idempotency-Key 最长 128 个字符")
		return
	}
	if old, err := db.FindResumeByIdempotency(scope, key, time.Now().Add(-10*time.Minute)); err == nil {
		businessOK(c, 201, createdView(old.ID, old.Status, old.CreatedAt))
		return
	}
	now := time.Now()
	row := &db.ResumeDO{ID: businessID("R", now), UserID: userID, IntentionServiceID: service.ID, ServiceName: service.Name,
		WorkStatus: req.WorkStatus, ExperienceRange: req.ExperienceRange, EntryYear: req.EntryYear, ContactName: req.ContactName, ContactPhone: req.ContactPhone, Status: "PENDING_CONTACT",
		IdempotencyKey: key, SubmissionScope: scope, CreatedAt: now, UpdatedAt: now}
	draft := &db.CaregiverDO{
		ID: "CG_" + row.ID, ApplicationID: row.ID, ContactPhone: req.ContactPhone, Name: req.ContactName, Age: 0,
		ExperienceYears: time.Now().Year() - req.EntryYear, ServiceIDs: db.MarshalStringSlice([]string{service.ID}),
		Jobs: db.MarshalStringSlice([]string{service.Name}), AvailabilityStatus: req.WorkStatus,
		Skills: "[]", Certificates: "[]", MedicalReportImageURLs: "[]", PersonalInfo: encodeJSONValue(gin.H{"contactName": req.ContactName, "contactPhone": req.ContactPhone}, "{}"),
		WorkHistory: "[]", PhotoURLs: "[]", DisplayFields: "{}", Status: "DRAFT", Source: "SELF_SUBMITTED",
		CreatedAt: now, UpdatedAt: now,
	}
	if err := db.CreateResumeAndCaregiverDraft(row, draft); err != nil {
		businessFail(c, 500, 50000, "提交求职申请失败")
		return
	}
	businessOK(c, 201, createdView(row.ID, row.Status, row.CreatedAt))
}

func BusinessUserMe(c *gin.Context) {
	_, userID := requestUser(c)
	if userID == "" {
		businessFail(c, 401, 40100, "未登录或 Token 缺失")
		return
	}
	row, err := db.GetMiniUserProfile(userID)
	if err != nil {
		businessFail(c, 401, 40101, "Token 失效")
		return
	}
	if row.Status != db.CustomerStatusActive {
		businessFail(c, 403, 40300, "账号已被禁用")
		return
	}
	businessOK(c, 200, gin.H{"id": row.ID, "nickName": row.Nickname, "avatarUrl": row.Avatar,
		"phone": maskPhone(row.Phone), "signature": row.Signature, "loginAt": row.LastLoginAt})
}

func BusinessAppConfig(c *gin.Context) {
	banners := decodeJSONMap(db.MustAppConfigJSON("mini.decoration.banners", `{"items":[]}`))
	customerService := decodeJSONValue(db.MustAppConfigJSON("mini.decoration.customer_service", `{"name":"","phone":"","avatarUrl":""}`), map[string]interface{}{})
	company := decodeJSONMap(db.MustAppConfigJSON("mini.decoration.company", `{"logoUrl":"","name":"永和护理","address":"","introduction":"","serviceGuarantees":[],"contactPhone":""}`))
	company["serviceGuarantees"] = db.NormalizeServiceGuarantees(company["serviceGuarantees"])
	homeBanners := db.NormalizeHomeBanners(banners)
	bannerURLs := make([]string, 0, len(homeBanners))
	for _, banner := range homeBanners {
		bannerURLs = append(bannerURLs, banner.ImageURL)
	}
	businessOK(c, 200, gin.H{"bannerUrls": bannerURLs, "homeBanners": homeBanners, "consultant": customerService, "customerService": customerService, "company": company, "trustItems": company["serviceGuarantees"]})
}
func BusinessAbout(c *gin.Context) {
	company := decodeJSONMap(db.MustAppConfigJSON("mini.decoration.company", `{"logoUrl":"","name":"永和护理","address":"","introduction":"","serviceGuarantees":[],"contactPhone":""}`))
	company["serviceGuarantees"] = db.NormalizeServiceGuarantees(company["serviceGuarantees"])
	businessOK(c, 200, company)
}
func BusinessAgreement(c *gin.Context) {
	kind := c.Param("type")
	if kind != "privacy" && kind != "service" {
		businessFail(c, 404, 40400, "协议不存在")
		return
	}
	businessConfig(c, "mini.business.agreement."+kind, fmt.Sprintf(`{"title":"%s","version":"1.0","updatedAt":"2026-07-18T00:00:00+08:00","effectiveAt":"2026-07-18T00:00:00+08:00","intro":"","sections":[]}`, map[string]string{"privacy": "隐私政策", "service": "用户服务协议"}[kind]))
}

func businessConfig(c *gin.Context, key, fallback string) {
	raw := db.MustAppConfigJSON(key, fallback)
	businessOK(c, 200, decodeJSONValue(raw, map[string]interface{}{}))
}
func decodeJSONValue(raw string, fallback interface{}) interface{} {
	var value interface{}
	if json.Unmarshal([]byte(raw), &value) != nil {
		return fallback
	}
	return value
}

func encodeJSONValue(value interface{}, fallback string) string {
	raw, err := json.Marshal(value)
	if err != nil {
		return fallback
	}
	return string(raw)
}

func displayFields(raw string) map[string]bool {
	fields := map[string]bool{}
	_ = json.Unmarshal([]byte(raw), &fields)
	return fields
}
func decodeJSONMap(raw string) map[string]interface{} {
	value := map[string]interface{}{}
	_ = json.Unmarshal([]byte(raw), &value)
	return value
}

func ListBusinessFAQs(c *gin.Context) {
	rows, err := db.ListPublicFAQs(strings.TrimSpace(c.Query("category")))
	if err != nil {
		businessFail(c, 500, 50000, "查询常见问题失败")
		return
	}
	businessOK(c, 200, gin.H{"list": rows})
}
func createdView(id, status string, createdAt time.Time) gin.H {
	return gin.H{"id": id, "status": status, "createdAt": createdAt.Format(time.RFC3339)}
}
func businessID(prefix string, now time.Time) string {
	return fmt.Sprintf("%s%s%s", prefix, now.Format("20060102150405"), randomHex(4))
}
func submissionScope(userID, phone string) string {
	if userID != "" {
		return "user:" + userID
	}
	return "phone:" + phone
}
func maskPhone(phone string) string {
	if len(phone) == 11 {
		return phone[:3] + "****" + phone[7:]
	}
	return phone
}
func stringSet(values ...string) map[string]bool {
	out := map[string]bool{}
	for _, value := range values {
		out[value] = true
	}
	return out
}
func contains(values []string, target string) bool {
	for _, value := range values {
		if value == target {
			return true
		}
	}
	return false
}
func experienceMatches(years int, value string) bool {
	switch value {
	case "LESS_THAN_1_YEAR":
		return years < 1
	case "YEAR_1_TO_3":
		return years >= 1 && years <= 3
	case "YEAR_3_TO_5":
		return years >= 3 && years <= 5
	case "YEAR_5_TO_10":
		return years >= 5 && years <= 10
	case "MORE_THAN_10_YEARS":
		return years > 10
	}
	return false
}
