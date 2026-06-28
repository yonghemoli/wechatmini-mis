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
	{
		ID:                 "storage-organizing",
		Category:           "housekeeping",
		Name:               "收纳整理",
		Scene:              "适合衣柜、厨房和杂物区重新整理",
		Summary:            "分类归位、动线整理和日常收纳建议",
		PriceText:          "预约后报价",
		DurationText:       "约 3 小时起",
		RequirementLabel:   "整理区域",
		RequirementOptions: []string{"衣柜", "厨房", "全屋局部"},
		SuitableFor:        []string{"搬家后整理", "换季衣物整理", "家庭空间优化"},
		Scope:              []string{"物品分类", "空间规划", "收纳归位"},
		Process:            []string{"需求沟通", "上门评估", "现场整理", "交付建议"},
		Notes:              []string{"不包含大件搬运", "贵重物品请提前自行保管"},
	},
	{
		ID:                 "elder-home-care",
		Category:           "eldercare",
		Name:               "居家照护",
		Scene:              "适合老人日常陪护、起居协助",
		Summary:            "陪伴照护、基础清洁和生活协助",
		PriceText:          "预约后报价",
		DurationText:       "半天起约",
		RequirementLabel:   "照护时长",
		RequirementOptions: []string{"半天", "全天", "长期咨询"},
		SuitableFor:        []string{"行动较慢老人", "术后恢复期", "子女不在身边"},
		Scope:              []string{"起居协助", "安全陪伴", "基础照护"},
		Process:            []string{"了解情况", "匹配人员", "确认时间", "上门服务"},
		Notes:              []string{"不提供医疗诊断", "特殊护理需提前说明"},
	},
	{
		ID:                 "elder-hospital-escort",
		Category:           "eldercare",
		Name:               "陪诊服务",
		Scene:              "适合老人就医挂号、取药和检查陪同",
		Summary:            "协助就诊流程、排队取号和结果领取",
		PriceText:          "预约后报价",
		DurationText:       "约半天",
		RequirementLabel:   "医院区域",
		RequirementOptions: []string{"青秀区", "兴宁区", "西乡塘区"},
		SuitableFor:        []string{"老人复诊", "行动不便", "家属无法陪同"},
		Scope:              []string{"挂号取号", "检查陪同", "取药协助"},
		Process:            []string{"确认医院", "约定见面", "陪同就诊", "反馈结果"},
		Notes:              []string{"不代替家属作医疗决定", "请提前准备身份证件和病历"},
	},
	{
		ID:                 "child-temporary-care",
		Category:           "childcare",
		Name:               "临时看护",
		Scene:              "适合家长短时外出或临时照看",
		Summary:            "陪伴游戏、餐点提醒和基础安全看护",
		PriceText:          "预约后报价",
		DurationText:       "约 2 小时起",
		RequirementLabel:   "孩子年龄",
		RequirementOptions: []string{"3-6 岁", "7-9 岁", "10-12 岁"},
		SuitableFor:        []string{"临时加班", "短时外出", "课后陪伴"},
		Scope:              []string{"安全陪伴", "餐点提醒", "简单互动"},
		Process:            []string{"确认年龄", "匹配人员", "上门看护", "服务反馈"},
		Notes:              []string{"不提供医疗护理", "请提前说明过敏和禁忌"},
	},
	{
		ID:                 "child-pickup",
		Category:           "childcare",
		Name:               "接送陪伴",
		Scene:              "适合放学接送和短时等待",
		Summary:            "按约定地点接送，并陪伴至家长交接",
		PriceText:          "预约后报价",
		DurationText:       "按次服务",
		RequirementLabel:   "接送方式",
		RequirementOptions: []string{"步行接送", "公共交通", "家长指定方式"},
		SuitableFor:        []string{"放学接送", "兴趣班接送", "家长临时不便"},
		Scope:              []string{"身份核验", "路线确认", "安全交接"},
		Process:            []string{"确认地点", "人员接单", "到点接送", "完成交接"},
		Notes:              []string{"不使用未确认车辆", "需提前提供接送信息"},
	},
}

var storeList = []Store{
	{ID: "store_1", Name: "永和护理服务中心", ContactName: "李店长", Phone: "0771-8888888", Address: "广西南宁青秀区民族大道 100 号", BusinessHours: "09:00-18:00", Status: "open", DistanceText: "距您约 3.2km · 南宁家庭服务", Latitude: 22.81673, Longitude: 108.3669},
}

var addresses = []Address{
	{ID: "addr_1", ContactName: "张三", Phone: "13800000000", District: "广西 南宁 青秀区", Detail: "3 栋 2 单元 801", Tag: "家", IsDefault: true},
}

var targets = []ServiceTarget{
	{ID: "target_1", Name: "父亲", Category: "eldercare", Relation: "父母", Age: "70-79 岁", Note: "行动较慢", IsDefault: true},
}

var dishes = []Dish{
	{ID: "tomato-egg", Name: "番茄炒蛋", Scene: "酸甜开胃，适合日常晚餐", Tag: "家常", Price: 12, Ingredients: []string{"番茄 300g", "鸡蛋 3 个"}, VideoTitle: "番茄炒蛋 8 分钟快手做法", VideoURL: "https://example.com/tomato-egg", Comments: []string{"孩子很爱吃"}},
	{ID: "pepper-pork", Name: "青椒肉丝", Scene: "下饭快手菜", Tag: "家常", Price: 16, Ingredients: []string{"青椒 250g", "猪肉 200g"}, VideoTitle: "青椒肉丝家常做法", VideoURL: "https://example.com/pepper-pork", Comments: []string{"适合工作日晚餐"}},
	{ID: "seaweed-egg-soup", Name: "紫菜蛋花汤", Scene: "清淡快手，适合搭配主菜", Tag: "汤品", Price: 10, Ingredients: []string{"紫菜 20g", "鸡蛋 2 个"}, VideoTitle: "紫菜蛋花汤家常做法", VideoURL: "https://example.com/seaweed-egg-soup", Comments: []string{"老人孩子都适合"}},
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
	r.POST("/auth/phone-code", PhoneCode)
	r.POST("/auth/phone-login", PhoneLogin)
	r.GET("/user/profile", UserProfile)
	r.GET("/home", GetHome)
	r.GET("/appointment/home", AppointmentHome)
	r.GET("/stores/nearest", NearestStore)
	r.GET("/service-categories", ListServiceCategories)
	r.GET("/service-categories/:id/services", ListServicesByCategoryID)
	r.GET("/services", ListServices)
	r.GET("/services/search", SearchServices)
	r.GET("/services/:id", GetService)
	r.GET("/service-areas", ServiceAreas)
	r.GET("/addresses", ListAddresses)
	r.POST("/addresses", CreateAddress)
	r.PUT("/addresses/:id/default", SetDefaultAddress)
	r.DELETE("/addresses/:id", DeleteAddress)
	r.GET("/service-targets", ListServiceTargets)
	r.POST("/service-targets", CreateServiceTarget)
	r.PUT("/service-targets/:id/default", SetDefaultServiceTarget)
	r.DELETE("/service-targets/:id", DeleteServiceTarget)
	r.GET("/meal/pricing", MealPricing)
	r.GET("/meal/dishes", ListDishes)
	r.GET("/meal/dishes/:nameOrId", GetDish)
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
	r.GET("/chat/session", ChatSession)
	r.GET("/chat/messages", ChatMessages)
	r.POST("/chat/messages", CreateChatMessage)
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
		PhoneCode     string `json:"phoneCode"`
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

func PhoneCode(c *gin.Context) {
	var req struct {
		Phone string `json:"phone"`
	}
	if err := c.ShouldBindJSON(&req); err != nil || strings.TrimSpace(req.Phone) == "" {
		fail(c, "phone required")
		return
	}
	ok(c, gin.H{"success": true})
}

func PhoneLogin(c *gin.Context) {
	var req struct {
		Phone string `json:"phone"`
		Code  string `json:"code"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		fail(c, "invalid request")
		return
	}
	if strings.TrimSpace(req.Phone) == "" || strings.TrimSpace(req.Code) == "" {
		fail(c, "phone and code required")
		return
	}
	currentUser.Phone = req.Phone
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

func AppointmentHome(c *gin.Context) {
	store := nearestStore()
	ok(c, gin.H{
		"store": store,
		"tabs": []gin.H{
			{"id": "all", "name": "所有"},
			{"id": "activity1", "name": "活动1"},
			{"id": "activity2", "name": "活动2"},
		},
		"groups": appointmentGroups(),
		"activities": gin.H{
			"activity1": gin.H{
				"title": "活动1",
				"desc":  "精选家庭清洁与照护服务",
				"items": serviceItems("housekeeping", "booking"),
			},
			"activity2": gin.H{
				"title": "活动2",
				"desc":  "适合老人和孩子的安心陪伴服务",
				"items": append(serviceItems("eldercare", "booking"), serviceItems("childcare", "booking")...),
			},
		},
	})
}

func NearestStore(c *gin.Context) {
	store := nearestStore()
	if store.ID == "" {
		fail(c, "store not found")
		return
	}
	ok(c, gin.H{"item": store})
}

func ListServiceCategories(c *gin.Context) {
	ok(c, gin.H{"list": serviceCategories()})
}

func ListServicesByCategoryID(c *gin.Context) {
	categoryID := strings.TrimSpace(c.Param("id"))
	items := appointmentItemsByGroupID(categoryID)
	if items == nil {
		fail(c, "service category not found")
		return
	}
	ok(c, gin.H{"list": items})
}

func SearchServices(c *gin.Context) {
	ListServices(c)
}

func ListServices(c *gin.Context) {
	category := strings.TrimSpace(c.Query("category"))
	keyword := strings.TrimSpace(c.Query("keyword"))
	list := make([]Service, 0, len(serviceList))
	for _, item := range serviceList {
		if category != "" && item.Category != category {
			continue
		}
		if keyword != "" && !serviceMatchesKeyword(item, keyword) {
			continue
		}
		list = append(list, item)
	}
	ok(c, gin.H{"list": list})
}

func GetService(c *gin.Context) {
	id := c.Param("id")
	for _, item := range serviceList {
		if item.ID == id {
			ok(c, item)
			return
		}
	}
	fail(c, "service not found")
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

func SetDefaultServiceTarget(c *gin.Context) {
	id := c.Param("id")
	category := ""
	for _, item := range targets {
		if item.ID == id {
			category = item.Category
			break
		}
	}
	if category == "" {
		fail(c, "service target not found")
		return
	}
	for i := range targets {
		if targets[i].Category == category {
			targets[i].IsDefault = targets[i].ID == id
		}
	}
	ok(c, gin.H{"id": id})
}

func MealPricing(c *gin.Context) {
	ok(c, gin.H{"dishPrice": 12, "deliveryFee": 8})
}

func ListDishes(c *gin.Context) {
	ok(c, gin.H{"list": dishes})
}

func GetDish(c *gin.Context) {
	nameOrID := c.Param("nameOrId")
	for _, item := range dishes {
		if item.ID == nameOrID || item.Name == nameOrID {
			ok(c, item)
			return
		}
	}
	fail(c, "dish not found")
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

func ChatSession(c *gin.Context) {
	user := bearerUser(c)
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
	row, err := chatws.CreateMiniMessage(req.SessionID, user.ID, user.NickName, req.MsgType, req.Content)
	if err != nil {
		fail(c, err.Error())
		return
	}
	ok(c, gin.H{"item": row})
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

func nearestStore() Store {
	for _, item := range storeList {
		if item.Status == "open" {
			return item
		}
	}
	if len(storeList) > 0 {
		return storeList[0]
	}
	return Store{}
}

func serviceCategories() []ServiceCategory {
	groups := appointmentGroups()
	categories := make([]ServiceCategory, 0, len(groups))
	for index, group := range groups {
		categories = append(categories, ServiceCategory{
			ID:   group.ID,
			Name: group.Name,
			Desc: group.Desc,
			Icon: group.Icon,
			Sort: (index + 1) * 10,
		})
	}
	return categories
}

func appointmentItemsByGroupID(id string) []AppointmentItem {
	for _, group := range appointmentGroups() {
		if group.ID == id {
			return group.Items
		}
	}
	return nil
}

func appointmentGroups() []AppointmentGroup {
	return []AppointmentGroup{
		{
			ID:    "cleaning",
			Name:  "清洁清洗",
			Desc:  "家电、玻璃、厨房等上门清洁",
			Icon:  "housekeeping",
			Items: []AppointmentItem{serviceItemByID("air-conditioner-cleaning", "清洁", "booking"), serviceItemByID("daily-cleaning", "保洁", "booking")},
		},
		{
			ID:    "organizing",
			Name:  "收纳整理",
			Desc:  "衣柜、厨房和家庭空间整理",
			Icon:  "housekeeping",
			Items: []AppointmentItem{serviceItemByID("storage-organizing", "收纳", "booking")},
		},
		{
			ID:    "eldercare",
			Name:  "养老护理",
			Desc:  "陪诊、居家照护和日常陪伴",
			Icon:  "eldercare",
			Items: []AppointmentItem{serviceItemByID("elder-home-care", "照护", "booking"), serviceItemByID("elder-hospital-escort", "陪诊", "booking")},
		},
		{
			ID:    "maternal-care",
			Name:  "母婴护理",
			Desc:  "月嫂、育婴和产后家庭照护咨询",
			Icon:  "care",
			Items: []AppointmentItem{{ID: "maternal-care-consult", Name: "母婴护理咨询", Scene: "适合产后家庭提前沟通服务需求", Summary: "由客服协助确认月嫂、育婴等服务安排", PriceText: "预约后报价", DurationText: "按需求确认", ImageText: "母婴", Category: "childcare", Action: "coming"}},
		},
		{
			ID:    "childcare",
			Name:  "育儿早教",
			Desc:  "临时看护、接送陪伴和亲子照护",
			Icon:  "childcare",
			Items: []AppointmentItem{serviceItemByID("child-temporary-care", "看护", "booking"), serviceItemByID("child-pickup", "接送", "booking")},
		},
		{
			ID:    "mealprep",
			Name:  "食谱配菜",
			Desc:  "按套餐或自选菜品配好食材送上门",
			Icon:  "mealprep",
			Items: []AppointmentItem{serviceItemByID("meal-prep", "配菜", "mealprep")},
		},
	}
}

func serviceItems(category, action string) []AppointmentItem {
	items := make([]AppointmentItem, 0)
	for _, item := range serviceList {
		if item.Category == category {
			items = append(items, serviceToAppointmentItem(item, imageTextByCategory(item.Category), action))
		}
	}
	return items
}

func serviceItemByID(id, imageText, action string) AppointmentItem {
	for _, item := range serviceList {
		if item.ID == id {
			return serviceToAppointmentItem(item, imageText, action)
		}
	}
	return AppointmentItem{ID: id, ImageText: imageText, Action: action}
}

func serviceToAppointmentItem(item Service, imageText, action string) AppointmentItem {
	return AppointmentItem{
		ID:           item.ID,
		Name:         item.Name,
		Scene:        item.Scene,
		Summary:      item.Summary,
		PriceText:    item.PriceText,
		DurationText: item.DurationText,
		ImageText:    imageText,
		Category:     item.Category,
		Action:       action,
	}
}

func imageTextByCategory(category string) string {
	switch category {
	case "eldercare":
		return "照护"
	case "childcare":
		return "育儿"
	case "mealprep":
		return "配菜"
	default:
		return "清洁"
	}
}

func serviceMatchesKeyword(item Service, keyword string) bool {
	target := strings.Join([]string{item.Name, item.Scene, item.Summary, item.Category}, " ")
	return strings.Contains(strings.ToLower(target), strings.ToLower(keyword))
}
