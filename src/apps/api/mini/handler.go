package miniapi

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"yonghemolimis/src/apps/api/chatws"
	"yonghemolimis/src/dao/db"

	"github.com/gin-gonic/gin"
)

type R struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type PageData struct {
	Page     int         `json:"page"`
	PageSize int         `json:"pageSize"`
	Total    int64       `json:"total"`
	List     interface{} `json:"list"`
}

type User struct {
	ID          string `json:"id"`
	NickName    string `json:"nickName"`
	AvatarURL   string `json:"avatarUrl"`
	Signature   string `json:"signature"`
	Phone       string `json:"phone"`
	LastLoginAt string `json:"lastLoginAt"`
}

type Service struct {
	ID                 string   `json:"id"`
	Category           string   `json:"category"`
	Name               string   `json:"name"`
	Scene              string   `json:"scene"`
	Summary            string   `json:"summary"`
	PriceText          string   `json:"priceText"`
	DurationText       string   `json:"durationText"`
	RequirementLabel   string   `json:"requirementLabel"`
	RequirementOptions []string `json:"requirementOptions"`
	SuitableFor        []string `json:"suitableFor"`
	Scope              []string `json:"scope"`
	Process            []string `json:"process"`
	Notes              []string `json:"notes"`
}

type AppointmentItem struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Scene        string `json:"scene"`
	Summary      string `json:"summary"`
	PriceText    string `json:"priceText"`
	DurationText string `json:"durationText"`
	ImageText    string `json:"imageText"`
	Category     string `json:"category"`
	Action       string `json:"action"`
}

type AppointmentGroup struct {
	ID    string            `json:"id"`
	Name  string            `json:"name"`
	Desc  string            `json:"desc"`
	Icon  string            `json:"icon"`
	Items []AppointmentItem `json:"items"`
}

type Store struct {
	ID            string  `json:"id"`
	Name          string  `json:"name"`
	ContactName   string  `json:"contactName"`
	Phone         string  `json:"phone"`
	Address       string  `json:"address"`
	BusinessHours string  `json:"businessHours"`
	Status        string  `json:"status"`
	DistanceText  string  `json:"distanceText"`
	Latitude      float64 `json:"latitude"`
	Longitude     float64 `json:"longitude"`
}

type ServiceCategory struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Desc string `json:"desc"`
	Icon string `json:"icon"`
	Sort int    `json:"sort"`
}

type Address struct {
	ID          string `json:"id"`
	ContactName string `json:"contactName"`
	Phone       string `json:"phone"`
	District    string `json:"district"`
	Detail      string `json:"detail"`
	Tag         string `json:"tag"`
	IsDefault   bool   `json:"isDefault"`
}

type ServiceTarget struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Category  string `json:"category"`
	Relation  string `json:"relation"`
	Age       string `json:"age"`
	Note      string `json:"note"`
	IsDefault bool   `json:"isDefault"`
}

type Dish struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Scene       string   `json:"scene"`
	Tag         string   `json:"tag"`
	Price       int      `json:"price"`
	Ingredients []string `json:"ingredients"`
	VideoTitle  string   `json:"videoTitle"`
	VideoURL    string   `json:"videoUrl"`
	Comments    []string `json:"comments"`
}

type MealPackage struct {
	ID     string   `json:"id"`
	Name   string   `json:"name"`
	Scene  string   `json:"scene"`
	Price  int      `json:"price"`
	Dishes []string `json:"dishes"`
}

type Order struct {
	ID                string            `json:"id"`
	ServiceID         string            `json:"serviceId"`
	ServiceName       string            `json:"serviceName"`
	Category          string            `json:"category"`
	Status            string            `json:"status"`
	AppointmentTime   string            `json:"appointmentTime"`
	Address           string            `json:"address"`
	ContactName       string            `json:"contactName"`
	Phone             string            `json:"phone"`
	Remark            string            `json:"remark"`
	ServiceTargetID   string            `json:"serviceTargetId"`
	ServiceTargetName string            `json:"serviceTargetName"`
	DetailFields      map[string]string `json:"detailFields"`
	Amount            int               `json:"amount"`
	CreatedAt         string            `json:"createdAt"`
}

var currentUser = User{
	ID:        "U10081",
	NickName:  "永和护理",
	AvatarURL: "",
	Phone:     "13800000000",
}

func RegisterRoutes(r gin.IRouter) {
	r.POST("/auth/wechat-login", WechatLogin)
	r.POST("/auth/phone-code", PhoneCode)
	r.POST("/auth/phone-login", PhoneLogin)
	r.GET("/users/me", BusinessUserMe)
	r.GET("/app-config", BusinessAppConfig)
	r.GET("/agreements/:type", BusinessAgreement)
	r.GET("/about", BusinessAbout)
	r.GET("/caregivers", ListBusinessCaregivers)
	r.GET("/caregivers/:id", GetBusinessCaregiver)
	r.POST("/demands", CreateBusinessDemand)
	r.POST("/resumes", CreateBusinessResume)
	r.GET("/user/profile", UserProfile)
	r.PUT("/user/profile", UpdateUserProfile)
	r.GET("/appointment/home", AppointmentHome)
	r.GET("/stores/nearest", NearestStore)
	r.GET("/service-categories", ListServiceCategories)
	r.GET("/service-categories/:id/services", ListServicesByCategoryID)
	// 新版小程序将 /services 定义为稳定服务分类。历史服务项目查询保留在
	// /legacy-services，避免路由语义冲突。
	r.GET("/services", ListBusinessServices)
	r.GET("/legacy-services", ListServices)
	r.GET("/services/search", SearchServices)
	r.GET("/services/:id", GetService)
	r.GET("/service-areas", ServiceAreas)
	r.GET("/addresses", ListAddresses)
	r.GET("/addresses/:id", GetAddress)
	r.POST("/addresses", CreateAddress)
	r.PUT("/addresses/:id", UpdateAddress)
	r.PUT("/addresses/:id/default", SetDefaultAddress)
	r.DELETE("/addresses/:id", DeleteAddress)
	r.GET("/service-targets", ListServiceTargets)
	r.GET("/service-targets/:id", GetServiceTarget)
	r.POST("/service-targets", CreateServiceTarget)
	r.PUT("/service-targets/:id", UpdateServiceTarget)
	r.PUT("/service-targets/:id/default", SetDefaultServiceTarget)
	r.DELETE("/service-targets/:id", DeleteServiceTarget)
	r.GET("/meal/pricing", MealPricing)
	r.GET("/meal/dishes", ListDishes)
	r.GET("/meal/dishes/:nameOrId", GetDish)
	r.GET("/meal/packages", ListMealPackages)
	r.GET("/meal/custom-packages", ListCustomPackages)
	r.GET("/meal/custom-packages/:id", GetCustomPackage)
	r.POST("/meal/custom-packages", CreateCustomPackage)
	r.PUT("/meal/custom-packages/:id", UpdateCustomPackage)
	r.DELETE("/meal/custom-packages/:id", DeleteCustomPackage)
	r.GET("/orders", ListOrders)
	r.GET("/orders/:id", GetOrder)
	r.POST("/orders", CreateOrder)
	r.PUT("/orders/:id/status", UpdateOrderStatus)
	r.POST("/orders/:id/cancel", CancelOrder)
	r.DELETE("/orders/:id", DeleteOrder)
	r.POST("/payments/wechat/prepay", WechatPrepay)
	r.POST("/payments/wechat/notify", WechatNotify)
	r.GET("/chat/session", ChatSession)
	r.GET("/chat/messages", ChatMessages)
	r.POST("/chat/messages", CreateChatMessage)
}

func ok(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, R{Code: 0, Message: "ok", Data: data})
}
func fail(c *gin.Context, msg string) { c.JSON(http.StatusOK, R{Code: -1, Message: msg}) }
func bearerUser(c *gin.Context) User {
	user, _ := requestUser(c)
	return user
}

func requestUser(c *gin.Context) (User, string) {
	token := strings.TrimPrefix(c.GetHeader("Authorization"), "Bearer ")
	userID := miniUserIDFromToken(token)
	if userID == "" {
		return User{}, ""
	}
	row, err := db.GetMiniUserProfile(userID)
	if err != nil {
		return User{ID: userID, NickName: currentUser.NickName, AvatarURL: currentUser.AvatarURL, Phone: currentUser.Phone}, userID
	}
	return userFromDO(*row), row.ID
}

func requestUserID(c *gin.Context) string {
	_, userID := requestUser(c)
	return userID
}

func requireUserID(c *gin.Context) (string, bool) {
	userID := requestUserID(c)
	if userID == "" {
		fail(c, "unauthorized")
		return "", false
	}
	return userID, true
}

func WechatLogin(c *gin.Context) {
	var req struct {
		Code          string `json:"code"`
		EncryptedData string `json:"encryptedData"`
		IV            string `json:"iv"`
		PhoneCode     string `json:"phoneCode"`
		NickName      string `json:"nickName"`
		AvatarURL     string `json:"avatarUrl"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		businessFail(c, http.StatusBadRequest, 40000, "请求参数错误")
		return
	}
	if strings.TrimSpace(req.Code) == "" {
		businessFail(c, http.StatusBadRequest, 40000, "微信登录凭证不能为空")
		return
	}
	wxSession, err := wxCodeToSession(req.Code)
	if err != nil {
		businessFail(c, http.StatusBadRequest, 40002, "微信登录凭证错误或已失效")
		return
	}
	if row, err := db.GetMiniUserByOpenID(wxSession.OpenID); err == nil && strings.TrimSpace(row.Phone) != "" {
		if row.Status != db.CustomerStatusActive {
			businessFail(c, http.StatusForbidden, 40300, "账号已被禁用")
			return
		}
		row.LastLoginAt = time.Now().Format("2006-01-02 15:04:05")
		row.UpdatedAt = time.Now()
		if err := db.UpsertMiniUserProfile(row); err != nil {
			businessFail(c, http.StatusInternalServerError, 50000, "登录失败")
			return
		}
		user := userFromDO(*row)
		user.Phone = maskPhone(row.Phone)
		businessOK(c, http.StatusOK, gin.H{"token": createMiniToken(row.ID), "expiresIn": 30 * 24 * 60 * 60, "user": user, "isBoundPhone": true})
		return
	}
	if strings.TrimSpace(req.PhoneCode) == "" {
		businessFail(c, http.StatusBadRequest, 40000, "phoneCode 不能为空")
		return
	}
	phone, err := wxPhoneNumber(req.PhoneCode)
	if err != nil {
		businessFail(c, http.StatusBadRequest, 40002, "微信手机号授权错误或已失效")
		return
	}
	if strings.TrimSpace(phone) == "" {
		businessFail(c, http.StatusBadRequest, 40001, "手机号格式不正确")
		return
	}
	if !mainlandPhonePattern.MatchString(phone) {
		businessFail(c, http.StatusBadRequest, 40001, "手机号格式不正确")
		return
	}
	row, err := ensureMiniUser(wxSession.OpenID, req.NickName, req.AvatarURL, phone)
	if err != nil {
		businessFail(c, http.StatusConflict, 40900, "该手机号或微信账号已绑定其他用户")
		return
	}
	token := createMiniToken(row.ID)
	user := userFromDO(*row)
	user.Phone = maskPhone(row.Phone)
	businessOK(c, http.StatusOK, gin.H{
		"token":        token,
		"expiresIn":    30 * 24 * 60 * 60,
		"user":         user,
		"isBoundPhone": true,
	})
}

func PhoneCode(c *gin.Context) {
	phoneCode(c)
}

func PhoneLogin(c *gin.Context) {
	phoneLogin(c)
}

func UserProfile(c *gin.Context) {
	user, _ := requestUser(c)
	if user.ID == "" {
		fail(c, "unauthorized")
		return
	}
	ok(c, user)
}

func UpdateUserProfile(c *gin.Context) {
	var req User
	if err := c.ShouldBindJSON(&req); err != nil {
		fail(c, "invalid request")
		return
	}
	userID, authed := requireUserID(c)
	if !authed {
		return
	}
	row, err := db.GetMiniUserProfile(userID)
	if err != nil {
		row = &db.CustomerDO{ID: userID, Status: db.CustomerStatusActive}
	}
	if req.NickName != "" {
		row.Nickname = req.NickName
	}
	if req.AvatarURL != "" {
		row.Avatar = req.AvatarURL
	}
	row.Signature = req.Signature
	if row.Nickname == "" {
		row.Nickname = currentUser.NickName
	}
	row.UpdatedAt = time.Now()
	if err := db.UpsertMiniUserProfile(row); err != nil {
		fail(c, err.Error())
		return
	}
	ok(c, userFromDO(*row))
}

func AppointmentHome(c *gin.Context) {
	store, _ := nearestStore()
	ok(c, gin.H{
		"store":      store,
		"tabs":       []gin.H{{"id": "all", "name": "所有"}},
		"groups":     appointmentGroupsFromDB(),
		"activities": gin.H{},
	})
}

func NearestStore(c *gin.Context) {
	store, okStore := nearestStore()
	if !okStore {
		fail(c, "store not found")
		return
	}
	ok(c, gin.H{"item": store})
}

func ListServiceCategories(c *gin.Context) { ok(c, gin.H{"list": appointmentGroupsFromDB()}) }

func ListServicesByCategoryID(c *gin.Context) {
	categoryID := strings.TrimSpace(c.Param("id"))
	items := appointmentItemsByGroupID(categoryID)
	if items == nil {
		fail(c, "service category not found")
		return
	}
	ok(c, gin.H{"list": items})
}

func SearchServices(c *gin.Context) { ListServices(c) }

func ListServices(c *gin.Context) {
	category := strings.TrimSpace(c.Query("category"))
	keyword := strings.TrimSpace(c.Query("keyword"))
	typeID, _ := strconv.Atoi(category)
	rows, err := db.ListMiniServices(uint(typeID), keyword)
	if err != nil {
		fail(c, err.Error())
		return
	}
	ok(c, gin.H{"list": servicesToHome(rows)})
}

func GetService(c *gin.Context) {
	row, err := db.GetMiniService(c.Param("id"))
	if err != nil {
		fail(c, "service not found")
		return
	}
	ok(c, serviceFromDOAsService(*row))
}

func ServiceAreas(c *gin.Context) {
	ok(c, gin.H{"city": "南宁", "districts": []string{"青秀区", "兴宁区", "西乡塘区", "良庆区", "江南区"}, "notes": []string{"具体地址以客服确认为准"}})
}

func ListAddresses(c *gin.Context) {
	userID, authed := requireUserID(c)
	if !authed {
		return
	}
	rows, err := db.ListMiniAddresses(userID)
	if err != nil {
		fail(c, err.Error())
		return
	}
	ok(c, gin.H{"list": addressRows(rows)})
}

func GetAddress(c *gin.Context) {
	userID, authed := requireUserID(c)
	if !authed {
		return
	}
	row, err := db.GetMiniAddress(userID, c.Param("id"))
	if err != nil {
		fail(c, "address not found")
		return
	}
	ok(c, addressRow(*row))
}

func CreateAddress(c *gin.Context) {
	userID, authed := requireUserID(c)
	if !authed {
		return
	}
	var req Address
	if err := c.ShouldBindJSON(&req); err != nil {
		fail(c, "invalid request")
		return
	}
	row := &db.AddressDO{
		ID:          req.ID,
		UserID:      userID,
		ContactName: req.ContactName,
		Phone:       req.Phone,
		District:    req.District,
		Detail:      req.Detail,
		Tag:         req.Tag,
		IsDefault:   req.IsDefault,
	}
	if err := db.CreateMiniAddress(row); err != nil {
		fail(c, err.Error())
		return
	}
	if row.IsDefault {
		_ = db.SetMiniAddressDefault(userID, row.ID)
	}
	ok(c, req)
}

func UpdateAddress(c *gin.Context) {
	userID, authed := requireUserID(c)
	if !authed {
		return
	}
	var req Address
	if err := c.ShouldBindJSON(&req); err != nil {
		fail(c, "invalid request")
		return
	}
	req.ID = c.Param("id")
	row := &db.AddressDO{
		ID:          req.ID,
		UserID:      userID,
		ContactName: req.ContactName,
		Phone:       req.Phone,
		District:    req.District,
		Detail:      req.Detail,
		Tag:         req.Tag,
		IsDefault:   req.IsDefault,
	}
	if err := db.UpdateMiniAddress(userID, row); err != nil {
		fail(c, err.Error())
		return
	}
	if row.IsDefault {
		_ = db.SetMiniAddressDefault(userID, row.ID)
	}
	ok(c, req)
}

func SetDefaultAddress(c *gin.Context) {
	userID, authed := requireUserID(c)
	if !authed {
		return
	}
	if err := db.SetMiniAddressDefault(userID, c.Param("id")); err != nil {
		fail(c, err.Error())
		return
	}
	ok(c, gin.H{"id": c.Param("id")})
}

func DeleteAddress(c *gin.Context) {
	userID, authed := requireUserID(c)
	if !authed {
		return
	}
	if err := db.DeleteMiniAddress(userID, c.Param("id")); err != nil {
		fail(c, err.Error())
		return
	}
	ok(c, gin.H{"id": c.Param("id")})
}

func ListServiceTargets(c *gin.Context) {
	userID, authed := requireUserID(c)
	if !authed {
		return
	}
	category := c.Query("category")
	rows, err := db.ListMiniServiceTargets(userID, category)
	if err != nil {
		fail(c, err.Error())
		return
	}
	ok(c, gin.H{"list": serviceTargetRows(rows)})
}

func GetServiceTarget(c *gin.Context) {
	userID, authed := requireUserID(c)
	if !authed {
		return
	}
	row, err := db.GetMiniServiceTarget(userID, c.Param("id"))
	if err != nil {
		fail(c, "service target not found")
		return
	}
	ok(c, serviceTargetRow(*row))
}

func CreateServiceTarget(c *gin.Context) {
	userID, authed := requireUserID(c)
	if !authed {
		return
	}
	var req ServiceTarget
	if err := c.ShouldBindJSON(&req); err != nil {
		fail(c, "invalid request")
		return
	}
	row := &db.ServiceTargetDO{
		ID:        req.ID,
		UserID:    userID,
		Name:      req.Name,
		Category:  req.Category,
		Relation:  req.Relation,
		Age:       req.Age,
		Note:      req.Note,
		IsDefault: req.IsDefault,
	}
	if err := db.CreateMiniServiceTarget(row); err != nil {
		fail(c, err.Error())
		return
	}
	if row.IsDefault {
		_ = db.SetMiniServiceTargetDefault(userID, row.ID)
	}
	ok(c, req)
}

func UpdateServiceTarget(c *gin.Context) {
	userID, authed := requireUserID(c)
	if !authed {
		return
	}
	var req ServiceTarget
	if err := c.ShouldBindJSON(&req); err != nil {
		fail(c, "invalid request")
		return
	}
	req.ID = c.Param("id")
	row := &db.ServiceTargetDO{
		ID:        req.ID,
		UserID:    userID,
		Name:      req.Name,
		Category:  req.Category,
		Relation:  req.Relation,
		Age:       req.Age,
		Note:      req.Note,
		IsDefault: req.IsDefault,
	}
	if err := db.UpdateMiniServiceTarget(userID, row); err != nil {
		fail(c, err.Error())
		return
	}
	if row.IsDefault {
		_ = db.SetMiniServiceTargetDefault(userID, row.ID)
	}
	ok(c, req)
}

func DeleteServiceTarget(c *gin.Context) {
	userID, authed := requireUserID(c)
	if !authed {
		return
	}
	if err := db.DeleteMiniServiceTarget(userID, c.Param("id")); err != nil {
		fail(c, err.Error())
		return
	}
	ok(c, gin.H{"id": c.Param("id")})
}

func SetDefaultServiceTarget(c *gin.Context) {
	userID, authed := requireUserID(c)
	if !authed {
		return
	}
	if err := db.SetMiniServiceTargetDefault(userID, c.Param("id")); err != nil {
		fail(c, err.Error())
		return
	}
	ok(c, gin.H{"id": c.Param("id")})
}

func MealPricing(c *gin.Context) { ok(c, gin.H{"dishPrice": 12, "deliveryFee": 8}) }

func ListDishes(c *gin.Context) {
	rows, err := db.ListMiniDishes()
	if err != nil {
		fail(c, err.Error())
		return
	}
	ok(c, gin.H{"list": dishRows(rows)})
}

func GetDish(c *gin.Context) {
	row, err := db.GetMiniDish(c.Param("nameOrId"))
	if err != nil {
		fail(c, "dish not found")
		return
	}
	ok(c, dishRow(*row))
}

func ListMealPackages(c *gin.Context) {
	rows, err := db.ListMiniMealPackages("", "official")
	if err != nil {
		fail(c, err.Error())
		return
	}
	ok(c, gin.H{"list": mealPackageRows(rows)})
}

func ListCustomPackages(c *gin.Context) {
	userID, authed := requireUserID(c)
	if !authed {
		return
	}
	rows, err := db.ListMiniMealPackages(userID, "custom")
	if err != nil {
		fail(c, err.Error())
		return
	}
	ok(c, gin.H{"list": mealPackageRows(rows)})
}

func GetCustomPackage(c *gin.Context) {
	userID, authed := requireUserID(c)
	if !authed {
		return
	}
	row, err := db.GetMiniMealPackage(userID, c.Param("id"))
	if err != nil {
		fail(c, "meal package not found")
		return
	}
	ok(c, mealPackageRow(*row))
}

func CreateCustomPackage(c *gin.Context) {
	userID, authed := requireUserID(c)
	if !authed {
		return
	}
	var req MealPackage
	if err := c.ShouldBindJSON(&req); err != nil {
		fail(c, "invalid request")
		return
	}
	row := &db.MealPackageDO{
		ID:          req.ID,
		UserID:      userID,
		PackageType: "custom",
		Name:        req.Name,
		Scene:       req.Scene,
		Price:       req.Price,
		Dishes:      db.MarshalStringSlice(req.Dishes),
	}
	if err := db.CreateMiniMealPackage(row); err != nil {
		fail(c, err.Error())
		return
	}
	req.ID = row.ID
	ok(c, req)
}

func UpdateCustomPackage(c *gin.Context) {
	userID, authed := requireUserID(c)
	if !authed {
		return
	}
	var req MealPackage
	if err := c.ShouldBindJSON(&req); err != nil {
		fail(c, "invalid request")
		return
	}
	req.ID = c.Param("id")
	row := &db.MealPackageDO{
		ID:          req.ID,
		UserID:      userID,
		PackageType: "custom",
		Name:        req.Name,
		Scene:       req.Scene,
		Price:       req.Price,
		Dishes:      db.MarshalStringSlice(req.Dishes),
	}
	if err := db.UpdateMiniMealPackage(userID, row); err != nil {
		fail(c, err.Error())
		return
	}
	ok(c, req)
}

func DeleteCustomPackage(c *gin.Context) {
	userID, authed := requireUserID(c)
	if !authed {
		return
	}
	if err := db.DeleteMiniMealPackage(userID, c.Param("id")); err != nil {
		fail(c, err.Error())
		return
	}
	ok(c, gin.H{"id": c.Param("id")})
}

func ListOrders(c *gin.Context) {
	userID, authed := requireUserID(c)
	if !authed {
		return
	}
	status := c.DefaultQuery("status", "all")
	page := parsePositiveInt(c.DefaultQuery("page", "1"), 1)
	pageSize := parsePositiveInt(c.DefaultQuery("pageSize", "20"), 20)
	rows, total, err := db.ListMiniOrders(userID, status, page, pageSize)
	if err != nil {
		fail(c, err.Error())
		return
	}
	ok(c, PageData{Page: page, PageSize: pageSize, Total: total, List: orderRows(rows)})
}

func GetOrder(c *gin.Context) {
	userID, authed := requireUserID(c)
	if !authed {
		return
	}
	row, err := db.GetMiniOrder(userID, c.Param("id"))
	if err != nil {
		fail(c, "order not found")
		return
	}
	ok(c, orderRow(*row))
}

func CreateOrder(c *gin.Context) {
	userID, authed := requireUserID(c)
	if !authed {
		return
	}
	var req Order
	if err := c.ShouldBindJSON(&req); err != nil {
		fail(c, "invalid request")
		return
	}
	row := &db.OrderDO{
		ID:            fmt.Sprintf("YH%s%03d", time.Now().Format("20060102150405"), time.Now().Nanosecond()%1000),
		UserID:        userID,
		Customer:      req.ContactName,
		Phone:         req.Phone,
		Service:       req.ServiceName,
		Amount:        req.Amount,
		Status:        db.OrderStatusPendingConfirm,
		Source:        "mini",
		AppointmentAt: req.AppointmentTime,
		Staff:         "",
		InternalNote:  req.Remark,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
	if err := db.Get().Create(row).Error; err != nil {
		fail(c, err.Error())
		return
	}
	ok(c, req)
}

func UpdateOrderStatus(c *gin.Context) {
	var req struct {
		Status string `json:"status"`
	}
	if err := c.ShouldBindJSON(&req); err != nil || req.Status == "" {
		fail(c, "invalid status")
		return
	}
	row, err := db.UpdateOrderStatus(c.Param("id"), req.Status, "")
	if err != nil {
		fail(c, err.Error())
		return
	}
	ok(c, orderRow(*row))
}

func DeleteOrder(c *gin.Context) {
	userID, authed := requireUserID(c)
	if !authed {
		return
	}
	if err := db.DeleteMiniOrder(userID, c.Param("id")); err != nil {
		fail(c, err.Error())
		return
	}
	ok(c, gin.H{"id": c.Param("id")})
}

func CancelOrder(c *gin.Context) {
	userID, authed := requireUserID(c)
	if !authed {
		return
	}
	if _, err := db.GetMiniOrder(userID, c.Param("id")); err != nil {
		fail(c, "order not found")
		return
	}
	row, err := db.UpdateOrderStatus(c.Param("id"), db.OrderStatusCancelled, "cancelled")
	if err != nil {
		fail(c, err.Error())
		return
	}
	ok(c, orderRow(*row))
}

func WechatPrepay(c *gin.Context) {
	var req struct {
		OrderID string `json:"orderId"`
		Amount  int    `json:"amount"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		fail(c, "invalid request")
		return
	}
	ok(c, gin.H{"timeStamp": strconv.FormatInt(time.Now().Unix(), 10), "nonceStr": "dev-mini-nonce", "package": "prepay_id=dev_prepay_id", "signType": "RSA", "paySign": "dev_pay_sign"})
}

func WechatNotify(c *gin.Context) { ok(c, gin.H{"received": true}) }

func ChatSession(c *gin.Context) {
	user := bearerUser(c)
	if user.ID == "" {
		fail(c, "unauthorized")
		return
	}
	sessionID := c.DefaultQuery("sessionId", "chat_"+user.ID)
	row, err := chatws.TouchSession(sessionID, user.ID, user.NickName)
	if err != nil {
		fail(c, err.Error())
		return
	}
	ok(c, gin.H{"item": row})
}

func ChatMessages(c *gin.Context) {
	user := bearerUser(c)
	if user.ID == "" {
		fail(c, "unauthorized")
		return
	}
	sessionID := c.DefaultQuery("sessionId", "chat_"+user.ID)
	rows, err := db.ListChatMessages(sessionID)
	if err != nil {
		fail(c, err.Error())
		return
	}
	ok(c, gin.H{"list": rows})
}

func CreateChatMessage(c *gin.Context) {
	var req struct {
		SessionID string `json:"sessionId"`
		Content   string `json:"content" binding:"required"`
		MsgType   string `json:"msgType"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		fail(c, "消息内容不能为空")
		return
	}
	user := bearerUser(c)
	if user.ID == "" {
		fail(c, "unauthorized")
		return
	}
	row, err := chatws.CreateMiniMessage(req.SessionID, user.ID, user.NickName, req.MsgType, req.Content)
	if err != nil {
		fail(c, err.Error())
		return
	}
	ok(c, gin.H{"item": row})
}

func parsePositiveInt(value string, fallback int) int {
	parsed, err := strconv.Atoi(value)
	if err != nil || parsed <= 0 {
		return fallback
	}
	return parsed
}

func nearestStore() (Store, bool) {
	rows, err := db.ListShops()
	if err != nil || len(rows) == 0 {
		return Store{}, false
	}
	for _, row := range rows {
		if row.Status == db.StatusOpen {
			return storeFromDO(row), true
		}
	}
	return storeFromDO(rows[0]), true
}

func storeFromDO(row db.ShopDO) Store {
	return Store{
		ID:            fmt.Sprintf("%d", row.ID),
		Name:          row.Name,
		ContactName:   row.ContactName,
		Phone:         row.Phone,
		Address:       row.Address,
		BusinessHours: row.BusinessHours,
		Status:        row.Status,
		DistanceText:  "距您约 3.2km · 南宁家庭服务",
	}
}

func ensureMiniUser(openID, nickName, avatarURL, phone string) (*db.CustomerDO, error) {
	if strings.TrimSpace(phone) == "" {
		return nil, fmt.Errorf("phone required")
	}
	row, err := db.GetMiniUserByPhone(phone)
	if err != nil {
		row = &db.CustomerDO{
			ID:        newMiniUserID(),
			Phone:     phone,
			Nickname:  "微信用户",
			Status:    db.CustomerStatusActive,
			CreatedAt: time.Now(),
		}
	}
	if row.OpenID != "" && openID != "" && row.OpenID != openID {
		return nil, fmt.Errorf("phone already bound to another wechat account")
	}
	if row.OpenID == "" && openID != "" {
		row.OpenID = openID
	}
	if nickName != "" {
		row.Nickname = nickName
	}
	if row.Nickname == "" {
		row.Nickname = "微信用户"
	}
	if avatarURL != "" {
		row.Avatar = avatarURL
	}
	row.Phone = phone
	row.LastLoginAt = time.Now().Format("2006-01-02 15:04:05")
	row.UpdatedAt = time.Now()
	return row, db.UpsertMiniUserProfile(row)
}

func newMiniUserID() string {
	return "U" + time.Now().Format("20060102150405") + fmt.Sprintf("%03d", time.Now().Nanosecond()%1000)
}

func userFromDO(row db.CustomerDO) User {
	return User{
		ID:          row.ID,
		NickName:    row.Nickname,
		AvatarURL:   row.Avatar,
		Signature:   row.Signature,
		Phone:       row.Phone,
		LastLoginAt: row.LastLoginAt,
	}
}

func appointmentGroupsFromDB() []AppointmentGroup {
	types, err := db.ListServiceTypes("")
	if err != nil {
		return []AppointmentGroup{}
	}
	rows, _ := db.ListMiniServices(0, "")
	grouped := map[uint][]AppointmentItem{}
	for _, row := range rows {
		grouped[row.TypeID] = append(grouped[row.TypeID], appointmentItemFromService(row))
	}
	out := make([]AppointmentGroup, 0, len(types))
	for _, typ := range types {
		out = append(out, AppointmentGroup{
			ID:    fmt.Sprintf("%d", typ.ID),
			Name:  typ.Name,
			Desc:  typ.Description,
			Icon:  "housekeeping",
			Items: grouped[typ.ID],
		})
	}
	return out
}

func appointmentItemsByGroupID(id string) []AppointmentItem {
	rows := appointmentGroupsFromDB()
	for _, group := range rows {
		if group.ID == id || group.Name == id {
			return group.Items
		}
	}
	return nil
}

func servicesToHome(rows []db.ServiceDO) []Service {
	out := make([]Service, 0, len(rows))
	for _, row := range rows {
		out = append(out, serviceFromDOAsService(row))
	}
	return out
}

func appointmentItemFromService(row db.ServiceDO) AppointmentItem {
	return AppointmentItem{
		ID:           fmt.Sprintf("%d", row.ID),
		Name:         row.Name,
		Scene:        row.Scene,
		Summary:      row.Summary,
		PriceText:    row.PriceText,
		DurationText: row.DurationText,
		ImageText:    row.TypeName,
		Category:     fmt.Sprintf("%d", row.TypeID),
		Action:       "booking",
	}
}

func serviceFromDOAsService(row db.ServiceDO) Service {
	return Service{
		ID:                 fmt.Sprintf("%d", row.ID),
		Category:           fmt.Sprintf("%d", row.TypeID),
		Name:               row.Name,
		Scene:              row.Scene,
		Summary:            row.Summary,
		PriceText:          row.PriceText,
		DurationText:       row.DurationText,
		RequirementLabel:   row.RequirementLabel,
		RequirementOptions: db.DecodeMiniStringSlice(row.RequirementOptions),
		SuitableFor:        db.DecodeMiniStringSlice(row.SuitableFor),
		Scope:              db.DecodeMiniStringSlice(row.Scope),
		Process:            db.DecodeMiniStringSlice(row.Process),
		Notes:              db.DecodeMiniStringSlice(row.Notes),
	}
}

func addressRows(rows []db.AddressDO) []Address {
	out := make([]Address, 0, len(rows))
	for _, row := range rows {
		out = append(out, addressRow(row))
	}
	return out
}

func addressRow(row db.AddressDO) Address {
	return Address{
		ID:          row.ID,
		ContactName: row.ContactName,
		Phone:       row.Phone,
		District:    row.District,
		Detail:      row.Detail,
		Tag:         row.Tag,
		IsDefault:   row.IsDefault,
	}
}

func serviceTargetRows(rows []db.ServiceTargetDO) []ServiceTarget {
	out := make([]ServiceTarget, 0, len(rows))
	for _, row := range rows {
		out = append(out, serviceTargetRow(row))
	}
	return out
}

func serviceTargetRow(row db.ServiceTargetDO) ServiceTarget {
	return ServiceTarget{
		ID:        row.ID,
		Name:      row.Name,
		Category:  row.Category,
		Relation:  row.Relation,
		Age:       row.Age,
		Note:      row.Note,
		IsDefault: row.IsDefault,
	}
}

func dishRows(rows []db.DishDO) []Dish {
	out := make([]Dish, 0, len(rows))
	for _, row := range rows {
		out = append(out, dishRow(row))
	}
	return out
}

func dishRow(row db.DishDO) Dish {
	return Dish{
		ID:          row.ID,
		Name:        row.Name,
		Scene:       row.Scene,
		Tag:         row.Tag,
		Price:       row.Price,
		Ingredients: db.DecodeMiniStringSlice(row.Ingredients),
		VideoTitle:  row.VideoTitle,
		VideoURL:    row.VideoURL,
		Comments:    db.DecodeMiniStringSlice(row.Comments),
	}
}

func mealPackageRows(rows []db.MealPackageDO) []MealPackage {
	out := make([]MealPackage, 0, len(rows))
	for _, row := range rows {
		out = append(out, mealPackageRow(row))
	}
	return out
}

func mealPackageRow(row db.MealPackageDO) MealPackage {
	return MealPackage{
		ID:     row.ID,
		Name:   row.Name,
		Scene:  row.Scene,
		Price:  row.Price,
		Dishes: db.DecodeMiniStringSlice(row.Dishes),
	}
}

func orderRows(rows []db.OrderDO) []Order {
	out := make([]Order, 0, len(rows))
	for _, row := range rows {
		out = append(out, orderRow(row))
	}
	return out
}

func orderRow(row db.OrderDO) Order {
	return Order{
		ID:              row.ID,
		ServiceName:     row.Service,
		Status:          row.Status,
		AppointmentTime: row.AppointmentAt,
		Address:         "",
		ContactName:     row.Customer,
		Phone:           row.Phone,
		Remark:          row.InternalNote,
		Amount:          row.Amount,
		CreatedAt:       row.CreatedAt.Format("2006-01-02 15:04:05"),
	}
}
