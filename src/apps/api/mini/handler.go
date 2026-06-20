package miniapi

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

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
	Total    int         `json:"total"`
	List     interface{} `json:"list"`
}

type User struct {
	ID        string `json:"id"`
	NickName  string `json:"nickName"`
	AvatarURL string `json:"avatarUrl"`
	Phone     string `json:"phone"`
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
	ID:        "u_1",
	NickName:  "张三",
	AvatarURL: "",
	Phone:     "13800000000",
}

var serviceList = []Service{
	{
		ID:                 "air-conditioner-cleaning",
		Category:           "housekeeping",
		Name:               "空调清洗",
		Scene:              "适合换季前后",
		Summary:            "清洗滤网、蒸发器和出风口",
		PriceText:          "预约后报价",
		DurationText:       "约 1 小时/台",
		RequirementLabel:   "清洗数量",
		RequirementOptions: []string{"挂机 1 台", "挂机 2 台", "柜机 1 台"},
		SuitableFor:        []string{"卧室空调", "客厅空调", "换季清洁"},
		Scope:              []string{"滤网清洗", "蒸发器清洁", "出风口擦洗"},
		Process:            []string{"预约确认", "师傅上门", "现场清洗", "验收完成"},
		Notes:              []string{"高空外机不在服务范围内", "具体地址以客服确认为准"},
	},
	{
		ID:                 "daily-cleaning",
		Category:           "housekeeping",
		Name:               "日常保洁",
		Scene:              "适合家庭日常维护",
		Summary:            "客厅、卧室、厨房、卫生间基础清洁",
		PriceText:          "66 元/小时起",
		DurationText:       "约 3 小时",
		RequirementLabel:   "服务时长",
		RequirementOptions: []string{"3 小时", "4 小时", "5 小时"},
		SuitableFor:        []string{"日常打扫", "租房保洁", "老人家庭"},
		Scope:              []string{"地面清洁", "台面擦拭", "厨卫基础清洁"},
		Process:            []string{"预约下单", "客服确认", "阿姨上门", "验收核销"},
		Notes:              []string{"不含高空作业", "不含重油污深度清洁"},
	},
	{
		ID:                 "meal-prep",
		Category:           "mealprep",
		Name:               "配菜",
		Scene:              "适合工作日晚餐",
		Summary:            "按套餐或自选菜品配好食材并送达",
		PriceText:          "按菜品计价",
		DurationText:       "按预约时间送达",
		RequirementLabel:   "套餐",
		RequirementOptions: []string{"三菜一汤套餐", "自定义配菜"},
		SuitableFor:        []string{"老人家庭", "上班族", "亲子家庭"},
		Scope:              []string{"食材采购", "清洗分装", "送货上门"},
		Process:            []string{"选择菜品", "确认金额", "预约送达", "签收"},
		Notes:              []string{"暂不提供现场烹饪", "菜价可能随市场波动"},
	},
}

var addresses = []Address{
	{ID: "addr_1", ContactName: "张三", Phone: "13800000000", District: "广西 南宁 青秀区", Detail: "3 栋 2 单元 801", Tag: "家", IsDefault: true},
}

var targets = []ServiceTarget{
	{ID: "target_1", Name: "父亲", Category: "eldercare", Relation: "父母", Age: "70-79 岁", Note: "行动较慢", IsDefault: true},
}

var dishes = []Dish{
	{Name: "番茄炒蛋", Scene: "酸甜开胃，适合日常晚餐", Tag: "家常", Price: 12, Ingredients: []string{"番茄 300g", "鸡蛋 3 个"}, VideoTitle: "番茄炒蛋 8 分钟快手做法", VideoURL: "https://example.com/tomato-egg", Comments: []string{"孩子很爱吃"}},
	{Name: "青椒肉丝", Scene: "下饭快手菜", Tag: "家常", Price: 16, Ingredients: []string{"青椒 250g", "猪肉 200g"}, VideoTitle: "青椒肉丝家常做法", VideoURL: "https://example.com/pepper-pork", Comments: []string{"适合工作日晚餐"}},
}

var officialPackages = []MealPackage{
	{ID: "pkg_1", Name: "三菜一汤套餐", Scene: "适合 2-3 人工作日晚餐", Price: 68, Dishes: []string{"番茄炒蛋", "青椒肉丝", "紫菜蛋花汤"}},
}

var customPackages = []MealPackage{
	{ID: "custom_pkg_1", Name: "我家的晚餐组合", Scene: "少油少盐", Price: 24, Dishes: []string{"番茄炒蛋", "青椒肉丝"}},
}

var orders = []Order{
	{
		ID:              "YH20260620001",
		ServiceID:       "air-conditioner-cleaning",
		ServiceName:     "空调清洗",
		Category:        "housekeeping",
		Status:          "pending",
		AppointmentTime: "2026-06-21 10:00",
		Address:         "广西 南宁 青秀区 3 栋 801",
		ContactName:     "张三",
		Phone:           "13800000000",
		Remark:          "",
		DetailFields:    map[string]string{"清洗数量": "挂机 1 台"},
		Amount:          0,
		CreatedAt:       "2026-06-20 11:20:00",
	},
}

func RegisterRoutes(r gin.IRouter) {
	r.POST("/auth/wechat-login", WechatLogin)
	r.GET("/user/profile", UserProfile)
	r.GET("/home", GetHome)
	r.GET("/services", ListServices)
	r.GET("/service-areas", ServiceAreas)
	r.GET("/addresses", ListAddresses)
	r.POST("/addresses", CreateAddress)
	r.PUT("/addresses/:id/default", SetDefaultAddress)
	r.DELETE("/addresses/:id", DeleteAddress)
	r.GET("/service-targets", ListServiceTargets)
	r.POST("/service-targets", CreateServiceTarget)
	r.DELETE("/service-targets/:id", DeleteServiceTarget)
	r.GET("/meal/dishes", ListDishes)
	r.GET("/meal/packages", ListMealPackages)
	r.GET("/meal/custom-packages", ListCustomPackages)
	r.POST("/meal/custom-packages", CreateCustomPackage)
	r.DELETE("/meal/custom-packages/:id", DeleteCustomPackage)
	r.GET("/orders", ListOrders)
	r.GET("/orders/:id", GetOrder)
	r.POST("/orders", CreateOrder)
	r.PUT("/orders/:id/status", UpdateOrderStatus)
	r.POST("/orders/:id/cancel", CancelOrder)
	r.POST("/payments/wechat/prepay", WechatPrepay)
	r.POST("/payments/wechat/notify", WechatNotify)
}

func ok(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, R{Code: 0, Message: "ok", Data: data})
}

func fail(c *gin.Context, msg string) {
	c.JSON(http.StatusOK, R{Code: -1, Message: msg})
}

func bearerUser(c *gin.Context) User {
	_ = strings.TrimPrefix(c.GetHeader("Authorization"), "Bearer ")
	return currentUser
}

func WechatLogin(c *gin.Context) {
	var req struct {
		Code          string `json:"code"`
		EncryptedData string `json:"encryptedData"`
		IV            string `json:"iv"`
		NickName      string `json:"nickName"`
		AvatarURL     string `json:"avatarUrl"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		fail(c, "invalid request")
		return
	}
	if req.NickName != "" {
		currentUser.NickName = req.NickName
	}
	if req.AvatarURL != "" {
		currentUser.AvatarURL = req.AvatarURL
	}
	ok(c, gin.H{"token": "dev-mini-token", "user": currentUser})
}

func UserProfile(c *gin.Context) {
	ok(c, bearerUser(c))
}

func GetHome(c *gin.Context) {
	ok(c, gin.H{
		"appName":  "永和护理",
		"banners":  []string{"/assets/banners/home-care.png", "/assets/banners/meal-prep.png"},
		"services": serviceList,
		"notice":   "服务预约后，客服将在 10 分钟内确认上门时间。",
		"serverAt": time.Now().Format("2006-01-02 15:04:05"),
	})
}

func ListServices(c *gin.Context) {
	ok(c, gin.H{"list": serviceList})
}

func ServiceAreas(c *gin.Context) {
	ok(c, gin.H{
		"city":      "南宁",
		"districts": []string{"青秀区", "兴宁区", "西乡塘区", "良庆区", "江南区"},
		"notes":     []string{"具体地址以客服确认为准"},
	})
}

func ListAddresses(c *gin.Context) {
	ok(c, gin.H{"list": addresses})
}

func CreateAddress(c *gin.Context) {
	var req Address
	if err := c.ShouldBindJSON(&req); err != nil {
		fail(c, "invalid request")
		return
	}
	req.ID = fmt.Sprintf("addr_%d", len(addresses)+1)
	if req.IsDefault {
		clearDefaultAddress()
	}
	addresses = append(addresses, req)
	ok(c, req)
}

func SetDefaultAddress(c *gin.Context) {
	id := c.Param("id")
	found := false
	for i := range addresses {
		addresses[i].IsDefault = addresses[i].ID == id
		if addresses[i].ID == id {
			found = true
		}
	}
	if !found {
		fail(c, "address not found")
		return
	}
	ok(c, gin.H{"id": id})
}

func DeleteAddress(c *gin.Context) {
	id := c.Param("id")
	addresses = filterAddresses(id)
	ok(c, gin.H{"id": id})
}

func ListServiceTargets(c *gin.Context) {
	category := c.Query("category")
	list := make([]ServiceTarget, 0, len(targets))
	for _, item := range targets {
		if category == "" || item.Category == category {
			list = append(list, item)
		}
	}
	ok(c, gin.H{"list": list})
}

func CreateServiceTarget(c *gin.Context) {
	var req ServiceTarget
	if err := c.ShouldBindJSON(&req); err != nil {
		fail(c, "invalid request")
		return
	}
	req.ID = fmt.Sprintf("target_%d", len(targets)+1)
	if req.IsDefault {
		for i := range targets {
			if targets[i].Category == req.Category {
				targets[i].IsDefault = false
			}
		}
	}
	targets = append(targets, req)
	ok(c, req)
}

func DeleteServiceTarget(c *gin.Context) {
	id := c.Param("id")
	next := targets[:0]
	for _, item := range targets {
		if item.ID != id {
			next = append(next, item)
		}
	}
	targets = next
	ok(c, gin.H{"id": id})
}

func ListDishes(c *gin.Context) {
	ok(c, gin.H{"list": dishes})
}

func ListMealPackages(c *gin.Context) {
	ok(c, gin.H{"list": officialPackages})
}

func ListCustomPackages(c *gin.Context) {
	ok(c, gin.H{"list": customPackages})
}

func CreateCustomPackage(c *gin.Context) {
	var req MealPackage
	if err := c.ShouldBindJSON(&req); err != nil {
		fail(c, "invalid request")
		return
	}
	req.ID = fmt.Sprintf("custom_pkg_%d", len(customPackages)+1)
	customPackages = append(customPackages, req)
	ok(c, req)
}

func DeleteCustomPackage(c *gin.Context) {
	id := c.Param("id")
	next := customPackages[:0]
	for _, item := range customPackages {
		if item.ID != id {
			next = append(next, item)
		}
	}
	customPackages = next
	ok(c, gin.H{"id": id})
}

func ListOrders(c *gin.Context) {
	status := c.DefaultQuery("status", "all")
	page := parsePositiveInt(c.DefaultQuery("page", "1"), 1)
	pageSize := parsePositiveInt(c.DefaultQuery("pageSize", "20"), 20)
	list := make([]Order, 0, len(orders))
	for _, item := range orders {
		if status == "all" || item.Status == status {
			list = append(list, item)
		}
	}
	start := (page - 1) * pageSize
	end := start + pageSize
	if start > len(list) {
		start = len(list)
	}
	if end > len(list) {
		end = len(list)
	}
	ok(c, PageData{Page: page, PageSize: pageSize, Total: len(list), List: list[start:end]})
}

func GetOrder(c *gin.Context) {
	for _, item := range orders {
		if item.ID == c.Param("id") {
			ok(c, item)
			return
		}
	}
	fail(c, "order not found")
}

func CreateOrder(c *gin.Context) {
	var req Order
	if err := c.ShouldBindJSON(&req); err != nil {
		fail(c, "invalid request")
		return
	}
	req.ID = fmt.Sprintf("YH%s%03d", time.Now().Format("20060102150405"), len(orders)+1)
	req.Status = "pending"
	req.CreatedAt = time.Now().Format("2006-01-02 15:04:05")
	orders = append([]Order{req}, orders...)
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
	for i := range orders {
		if orders[i].ID == c.Param("id") {
			orders[i].Status = req.Status
			ok(c, orders[i])
			return
		}
	}
	fail(c, "order not found")
}

func CancelOrder(c *gin.Context) {
	for i := range orders {
		if orders[i].ID == c.Param("id") {
			orders[i].Status = "cancelled"
			ok(c, orders[i])
			return
		}
	}
	fail(c, "order not found")
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
	ok(c, gin.H{
		"timeStamp": strconv.FormatInt(time.Now().Unix(), 10),
		"nonceStr":  "dev-mini-nonce",
		"package":   "prepay_id=dev_prepay_id",
		"signType":  "RSA",
		"paySign":   "dev_pay_sign",
	})
}

func WechatNotify(c *gin.Context) {
	ok(c, gin.H{"received": true})
}

func clearDefaultAddress() {
	for i := range addresses {
		addresses[i].IsDefault = false
	}
}

func filterAddresses(id string) []Address {
	next := addresses[:0]
	for _, item := range addresses {
		if item.ID != id {
			next = append(next, item)
		}
	}
	return next
}

func parsePositiveInt(value string, fallback int) int {
	parsed, err := strconv.Atoi(value)
	if err != nil || parsed <= 0 {
		return fallback
	}
	return parsed
}
