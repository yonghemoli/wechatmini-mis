package db

import "time"

// AdminDO 内部管理员账户。
type AdminDO struct {
	ID           uint       `gorm:"primaryKey" json:"id"`
	Username     string     `gorm:"size:50;uniqueIndex;not null" json:"username"`
	PasswordHash string     `gorm:"size:255;not null;column:password_hash" json:"-"`
	Name         string     `gorm:"size:50;not null" json:"name"`
	Email        string     `gorm:"size:128;not null" json:"email"`
	RoleID       *int64     `gorm:"column:role_id" json:"roleId"`
	IsSuperAdmin bool       `gorm:"not null;default:false" json:"isSuperAdmin"`
	Status       string     `gorm:"size:20;not null;default:'active'" json:"status"`
	LastLoginAt  *time.Time `gorm:"column:last_login_at" json:"lastLoginAt"`
	CreatedAt    time.Time  `json:"createdAt"`
	UpdatedAt    time.Time  `json:"updatedAt"`
}

func (AdminDO) TableName() string { return "admins" }

type OrderDO struct {
	ID            string    `gorm:"size:32;primaryKey" json:"id"`
	UserID        string    `gorm:"size:40;index;column:user_id" json:"userId"`
	Customer      string    `gorm:"size:64;not null" json:"customer"`
	Phone         string    `gorm:"size:32;not null" json:"phone"`
	Service       string    `gorm:"size:128;not null" json:"service"`
	Amount        int       `gorm:"not null;default:0" json:"amount"`
	Status        string    `gorm:"size:32;index;not null" json:"status"`
	Source        string    `gorm:"size:32;index;not null" json:"source"`
	AppointmentAt string    `gorm:"size:32;index;column:appointment_at" json:"appointmentAt"`
	Staff         string    `gorm:"size:64;not null;default:''" json:"staff"`
	InternalNote  string    `gorm:"type:text;column:internal_note" json:"internalNote"`
	CloseReason   string    `gorm:"type:text;column:close_reason" json:"closeReason"`
	MiniDeleted   bool      `gorm:"not null;default:false;column:mini_deleted" json:"miniDeleted"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

func (OrderDO) TableName() string { return "orders" }

type CustomerDO struct {
	ID          string    `gorm:"size:32;primaryKey" json:"id"`
	Avatar      string    `gorm:"size:255;not null;default:''" json:"avatar"`
	Nickname    string    `gorm:"size:64;not null" json:"nickname"`
	Phone       string    `gorm:"size:32;not null;default:''" json:"phone"`
	Signature   string    `gorm:"size:255;not null;default:''" json:"signature"`
	TotalSpent  int       `gorm:"not null;default:0;column:total_spent" json:"totalSpent"`
	LastOrderAt string    `gorm:"size:32;column:last_order_at" json:"lastOrderAt"`
	LastLoginAt string    `gorm:"size:32;column:last_login_at" json:"lastLoginAt"`
	Status      string    `gorm:"size:32;index;not null;default:'active'" json:"status"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

func (CustomerDO) TableName() string { return "users" }

type ServiceTypeDO struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"size:64;not null" json:"name"`
	Description string    `gorm:"type:text" json:"description"`
	SortOrder   int       `gorm:"not null;default:0;column:sort_order" json:"sortOrder"`
	Status      string    `gorm:"size:20;not null;default:'active'" json:"status"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

func (ServiceTypeDO) TableName() string { return "service_types" }

type ServiceDO struct {
	ID                 uint      `gorm:"primaryKey" json:"id"`
	TypeID             uint      `gorm:"index;not null;column:type_id" json:"typeId"`
	TypeName           string    `gorm:"-" json:"typeName"`
	Name               string    `gorm:"size:128;not null" json:"name"`
	Image              string    `gorm:"size:255;not null;default:''" json:"image"`
	Price              int       `gorm:"not null;default:0" json:"price"`
	Unit               string    `gorm:"size:32;not null;default:'小时'" json:"unit"`
	Title              string    `gorm:"size:128;not null;default:''" json:"title"`
	Scene              string    `gorm:"size:255;not null;default:''" json:"scene"`
	Summary            string    `gorm:"type:text" json:"summary"`
	PriceText          string    `gorm:"size:64;not null;default:'';column:price_text" json:"priceText"`
	DurationText       string    `gorm:"size:64;not null;default:'';column:duration_text" json:"durationText"`
	RequirementLabel   string    `gorm:"size:64;not null;default:'';column:requirement_label" json:"requirementLabel"`
	RequirementOptions string    `gorm:"type:text;not null;default:'';column:requirement_options" json:"requirementOptions"`
	SuitableFor        string    `gorm:"type:text;not null;default:'';column:suitable_for" json:"suitableFor"`
	Scope              string    `gorm:"type:text;not null;default:''" json:"scope"`
	Process            string    `gorm:"type:text;not null;default:''" json:"process"`
	Notes              string    `gorm:"type:text;not null;default:''" json:"notes"`
	Description        string    `gorm:"type:text" json:"description"`
	Visible            bool      `gorm:"not null;default:true;index" json:"visible"`
	SortOrder          int       `gorm:"not null;default:0;column:sort_order" json:"sortOrder"`
	CreatedAt          time.Time `json:"createdAt"`
	UpdatedAt          time.Time `json:"updatedAt"`
}

func (ServiceDO) TableName() string { return "services" }

type ShopDO struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	Name          string    `gorm:"size:128;not null" json:"name"`
	ContactName   string    `gorm:"size:64;column:contact_name" json:"contactName"`
	Phone         string    `gorm:"size:32" json:"phone"`
	Address       string    `gorm:"size:255" json:"address"`
	BusinessHours string    `gorm:"size:128;column:business_hours" json:"businessHours"`
	Status        string    `gorm:"size:20;not null;default:'open'" json:"status"`
	Remark        string    `gorm:"type:text" json:"remark"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

func (ShopDO) TableName() string { return "shops" }

type AddressDO struct {
	ID          string    `gorm:"size:32;primaryKey" json:"id"`
	UserID      string    `gorm:"size:40;index;column:user_id" json:"userId"`
	ContactName string    `gorm:"size:64;not null;column:contact_name" json:"contactName"`
	Phone       string    `gorm:"size:32;not null" json:"phone"`
	District    string    `gorm:"size:128;not null" json:"district"`
	Detail      string    `gorm:"size:255;not null" json:"detail"`
	Tag         string    `gorm:"size:32;not null;default:''" json:"tag"`
	IsDefault   bool      `gorm:"not null;default:false;column:is_default" json:"isDefault"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

func (AddressDO) TableName() string { return "addresses" }

type ServiceTargetDO struct {
	ID        string    `gorm:"size:32;primaryKey" json:"id"`
	UserID    string    `gorm:"size:40;index;column:user_id" json:"userId"`
	Name      string    `gorm:"size:64;not null" json:"name"`
	Category  string    `gorm:"size:32;index;not null" json:"category"`
	Relation  string    `gorm:"size:32;not null" json:"relation"`
	Age       string    `gorm:"size:32;not null" json:"age"`
	Note      string    `gorm:"size:255;not null;default:''" json:"note"`
	IsDefault bool      `gorm:"not null;default:false;column:is_default" json:"isDefault"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (ServiceTargetDO) TableName() string { return "service_targets" }

type DishDO struct {
	ID          string    `gorm:"size:64;primaryKey" json:"id"`
	Name        string    `gorm:"size:128;not null;uniqueIndex" json:"name"`
	Scene       string    `gorm:"size:255;not null" json:"scene"`
	Tag         string    `gorm:"size:32;not null;default:''" json:"tag"`
	Price       int       `gorm:"not null;default:0" json:"price"`
	Ingredients string    `gorm:"type:text;not null;default:''" json:"ingredients"`
	VideoTitle  string    `gorm:"size:255;not null;default:'';column:video_title" json:"videoTitle"`
	VideoURL    string    `gorm:"size:255;not null;default:'';column:video_url" json:"videoUrl"`
	Comments    string    `gorm:"type:text;not null;default:''" json:"comments"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

func (DishDO) TableName() string { return "dishes" }

type MealPackageDO struct {
	ID          string    `gorm:"size:64;primaryKey" json:"id"`
	UserID      string    `gorm:"size:40;index;column:user_id" json:"userId"`
	PackageType string    `gorm:"size:20;index;not null;default:'official';column:package_type" json:"packageType"`
	Name        string    `gorm:"size:128;not null" json:"name"`
	Scene       string    `gorm:"size:255;not null" json:"scene"`
	Price       int       `gorm:"not null;default:0" json:"price"`
	Dishes      string    `gorm:"type:text;not null;default:''" json:"dishes"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

func (MealPackageDO) TableName() string { return "meal_packages" }

type AppConfigDO struct {
	Key       string    `gorm:"size:128;primaryKey" json:"key"`
	Value     string    `gorm:"type:longtext;not null" json:"value"`
	Note      string    `gorm:"size:255;not null;default:''" json:"note"`
	UpdatedAt time.Time `json:"updatedAt"`
	CreatedAt time.Time `json:"createdAt"`
}

func (AppConfigDO) TableName() string { return "app_configs" }

type FAQDO struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Question  string    `gorm:"size:255;not null" json:"question"`
	Answer    string    `gorm:"type:text;not null" json:"answer"`
	Category  string    `gorm:"size:64;index" json:"category"`
	SortOrder int       `gorm:"not null;default:0;column:sort_order" json:"sortOrder"`
	Visible   bool      `gorm:"not null;default:true;index" json:"visible"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (FAQDO) TableName() string { return "faqs" }

type ChatSessionDO struct {
	ID          string    `gorm:"size:40;primaryKey" json:"id"`
	UserID      string    `gorm:"size:40;column:user_id" json:"userId"`
	UserName    string    `gorm:"size:64;column:user_name" json:"userName"`
	UserAvatar  string    `gorm:"size:255;column:user_avatar" json:"userAvatar"`
	Status      string    `gorm:"size:20;not null;default:'open'" json:"status"`
	LastMessage string    `gorm:"size:255;column:last_message" json:"lastMessage"`
	UnreadCount int       `gorm:"not null;default:0;column:unread_count" json:"unreadCount"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

func (ChatSessionDO) TableName() string { return "chat_sessions" }

type ChatMessageDO struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	SessionID string    `gorm:"size:40;index;not null;column:session_id" json:"sessionId"`
	Sender    string    `gorm:"size:20;not null" json:"sender"`
	MsgType   string    `gorm:"size:20;not null;default:'text';column:msg_type" json:"msgType"`
	Content   string    `gorm:"type:text;not null" json:"content"`
	IsRead    bool      `gorm:"not null;default:false;column:is_read" json:"isRead"`
	CreatedAt time.Time `json:"createdAt"`
}

func (ChatMessageDO) TableName() string { return "chat_messages" }
