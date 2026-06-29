package db

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"
)

const (
	StatusActive   = "active"
	StatusDisabled = "disabled"
	StatusOpen     = "open"
	StatusClosed   = "closed"
)

func ListServiceTypes(keyword string) ([]ServiceTypeDO, error) {
	var rows []ServiceTypeDO
	query := Get().Model(&ServiceTypeDO{})
	if keyword != "" {
		query = query.Where("name LIKE ?", "%"+keyword+"%")
	}
	err := query.Order("sort_order ASC, id DESC").Find(&rows).Error
	return rows, err
}

func CreateServiceType(row *ServiceTypeDO) error {
	if row.Status == "" {
		row.Status = StatusActive
	}
	return Get().Create(row).Error
}

func UpdateServiceType(id uint, input ServiceTypeDO) (*ServiceTypeDO, error) {
	updates := map[string]any{
		"name":        input.Name,
		"description": input.Description,
		"sort_order":  input.SortOrder,
		"updated_at":  time.Now(),
	}
	if input.Status != "" {
		updates["status"] = input.Status
	}
	if err := Get().Model(&ServiceTypeDO{}).Where("id = ?", id).Updates(updates).Error; err != nil {
		return nil, err
	}
	var row ServiceTypeDO
	if err := Get().First(&row, id).Error; err != nil {
		return nil, err
	}
	return &row, nil
}

func DeleteServiceType(id uint) error {
	return Get().Delete(&ServiceTypeDO{}, id).Error
}

func UpdateServiceTypeStatus(id uint, status string) (*ServiceTypeDO, error) {
	if err := Get().Model(&ServiceTypeDO{}).Where("id = ?", id).
		Updates(map[string]any{"status": status, "updated_at": time.Now()}).Error; err != nil {
		return nil, err
	}
	var row ServiceTypeDO
	if err := Get().First(&row, id).Error; err != nil {
		return nil, err
	}
	return &row, nil
}

func ListServices(typeID uint, keyword string) ([]ServiceDO, error) {
	var rows []ServiceDO
	query := Get().Model(&ServiceDO{})
	if typeID > 0 {
		query = query.Where("type_id = ?", typeID)
	}
	if keyword != "" {
		query = query.Where("name LIKE ?", "%"+keyword+"%")
	}
	err := query.Order("sort_order ASC, id DESC").Find(&rows).Error
	if err != nil {
		return nil, err
	}
	typeNames := map[uint]string{}
	var types []ServiceTypeDO
	_ = Get().Find(&types).Error
	for _, item := range types {
		typeNames[item.ID] = item.Name
	}
	for i := range rows {
		rows[i].TypeName = typeNames[rows[i].TypeID]
	}
	return rows, nil
}

func CreateService(row *ServiceDO) error {
	return Get().Create(row).Error
}

func UpdateService(id uint, input ServiceDO) (*ServiceDO, error) {
	updates := map[string]any{
		"type_id":     input.TypeID,
		"name":        input.Name,
		"image":       input.Image,
		"price":       input.Price,
		"unit":        input.Unit,
		"description": input.Description,
		"visible":     input.Visible,
		"sort_order":  input.SortOrder,
		"updated_at":  time.Now(),
	}
	if err := Get().Model(&ServiceDO{}).Where("id = ?", id).Updates(updates).Error; err != nil {
		return nil, err
	}
	var row ServiceDO
	if err := Get().First(&row, id).Error; err != nil {
		return nil, err
	}
	return &row, nil
}

func DeleteService(id uint) error {
	return Get().Delete(&ServiceDO{}, id).Error
}

func UpdateServiceVisible(id uint, visible bool) (*ServiceDO, error) {
	if err := Get().Model(&ServiceDO{}).Where("id = ?", id).
		Updates(map[string]any{"visible": visible, "updated_at": time.Now()}).Error; err != nil {
		return nil, err
	}
	var row ServiceDO
	if err := Get().First(&row, id).Error; err != nil {
		return nil, err
	}
	return &row, nil
}

func ListMiniServices(typeID uint, keyword string) ([]ServiceDO, error) {
	var rows []ServiceDO
	query := Get().Model(&ServiceDO{}).Where("visible = ?", true)
	if typeID > 0 {
		query = query.Where("type_id = ?", typeID)
	}
	if keyword != "" {
		like := "%" + strings.TrimSpace(keyword) + "%"
		query = query.Where("name LIKE ? OR title LIKE ? OR scene LIKE ? OR summary LIKE ?", like, like, like, like)
	}
	if err := query.Order("sort_order ASC, id DESC").Find(&rows).Error; err != nil {
		return nil, err
	}
	typeNames := map[uint]string{}
	var types []ServiceTypeDO
	_ = Get().Find(&types).Error
	for _, item := range types {
		typeNames[item.ID] = item.Name
	}
	for i := range rows {
		rows[i].TypeName = typeNames[rows[i].TypeID]
	}
	return rows, nil
}

func GetMiniService(id string) (*ServiceDO, error) {
	var row ServiceDO
	if err := Get().First(&row, "id = ? AND visible = ?", id, true).Error; err != nil {
		return nil, err
	}
	return &row, nil
}

func DecodeMiniStringSlice(value string) []string {
	var rows []string
	if err := json.Unmarshal([]byte(value), &rows); err != nil {
		return []string{}
	}
	return rows
}

func ListShops() ([]ShopDO, error) {
	var rows []ShopDO
	err := Get().Order("id DESC").Find(&rows).Error
	return rows, err
}

func CreateShop(row *ShopDO) error {
	if row.Status == "" {
		row.Status = StatusOpen
	}
	return Get().Create(row).Error
}

func UpdateShop(id uint, input ShopDO) (*ShopDO, error) {
	updates := map[string]any{
		"name":           input.Name,
		"contact_name":   input.ContactName,
		"phone":          input.Phone,
		"address":        input.Address,
		"business_hours": input.BusinessHours,
		"remark":         input.Remark,
		"updated_at":     time.Now(),
	}
	if input.Status != "" {
		updates["status"] = input.Status
	}
	if err := Get().Model(&ShopDO{}).Where("id = ?", id).Updates(updates).Error; err != nil {
		return nil, err
	}
	var row ShopDO
	if err := Get().First(&row, id).Error; err != nil {
		return nil, err
	}
	return &row, nil
}

func DeleteShop(id uint) error {
	return Get().Delete(&ShopDO{}, id).Error
}

func UpdateShopStatus(id uint, status string) (*ShopDO, error) {
	if err := Get().Model(&ShopDO{}).Where("id = ?", id).
		Updates(map[string]any{"status": status, "updated_at": time.Now()}).Error; err != nil {
		return nil, err
	}
	var row ShopDO
	if err := Get().First(&row, id).Error; err != nil {
		return nil, err
	}
	return &row, nil
}

func ListFAQs(category string) ([]FAQDO, error) {
	var rows []FAQDO
	query := Get().Model(&FAQDO{})
	if category != "" {
		query = query.Where("category = ?", category)
	}
	err := query.Order("sort_order ASC, id DESC").Find(&rows).Error
	return rows, err
}

func CreateFAQ(row *FAQDO) error {
	return Get().Create(row).Error
}

func UpdateFAQ(id uint, input FAQDO) (*FAQDO, error) {
	updates := map[string]any{
		"question":   input.Question,
		"answer":     input.Answer,
		"category":   input.Category,
		"sort_order": input.SortOrder,
		"visible":    input.Visible,
		"updated_at": time.Now(),
	}
	if err := Get().Model(&FAQDO{}).Where("id = ?", id).Updates(updates).Error; err != nil {
		return nil, err
	}
	var row FAQDO
	if err := Get().First(&row, id).Error; err != nil {
		return nil, err
	}
	return &row, nil
}

func DeleteFAQ(id uint) error {
	return Get().Delete(&FAQDO{}, id).Error
}

func UpdateFAQVisible(id uint, visible bool) (*FAQDO, error) {
	if err := Get().Model(&FAQDO{}).Where("id = ?", id).
		Updates(map[string]any{"visible": visible, "updated_at": time.Now()}).Error; err != nil {
		return nil, err
	}
	var row FAQDO
	if err := Get().First(&row, id).Error; err != nil {
		return nil, err
	}
	return &row, nil
}

func ListChatSessions() ([]ChatSessionDO, error) {
	var rows []ChatSessionDO
	err := Get().Order("updated_at DESC").Find(&rows).Error
	return rows, err
}

func GetChatSession(sessionID string) (*ChatSessionDO, error) {
	var row ChatSessionDO
	if err := Get().First(&row, "id = ?", sessionID).Error; err != nil {
		return nil, err
	}
	return &row, nil
}

func ListChatMessages(sessionID string) ([]ChatMessageDO, error) {
	var rows []ChatMessageDO
	err := Get().Where("session_id = ?", sessionID).Order("created_at ASC").Find(&rows).Error
	return rows, err
}

func CreateChatMessage(sessionID, sender, msgType, content string) (*ChatMessageDO, error) {
	if msgType == "" {
		msgType = "text"
	}
	row := &ChatMessageDO{
		SessionID: sessionID,
		Sender:    sender,
		MsgType:   msgType,
		Content:   content,
		IsRead:    sender == "admin",
		CreatedAt: time.Now(),
	}
	if err := Get().Create(row).Error; err != nil {
		return nil, err
	}
	updates := map[string]any{
		"last_message": content,
		"updated_at":   time.Now(),
	}
	if sender != "admin" {
		updates["unread_count"] = gorm.Expr("unread_count + 1")
	}
	_ = Get().Model(&ChatSessionDO{}).Where("id = ?", sessionID).Updates(updates).Error
	return row, nil
}

func CloseChatSession(sessionID string) error {
	return Get().Model(&ChatSessionDO{}).Where("id = ?", sessionID).
		Updates(map[string]any{"status": StatusClosed, "updated_at": time.Now()}).Error
}

func ReadChatSession(sessionID string) error {
	if err := Get().Model(&ChatMessageDO{}).Where("session_id = ?", sessionID).
		Update("is_read", true).Error; err != nil {
		return err
	}
	return Get().Model(&ChatSessionDO{}).Where("id = ?", sessionID).
		Updates(map[string]any{"unread_count": 0, "updated_at": time.Now()}).Error
}

func EnsureChatSession(id, userID, userName string) error {
	row := ChatSessionDO{
		ID:        id,
		UserID:    userID,
		UserName:  userName,
		Status:    StatusOpen,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	return Get().FirstOrCreate(&row, ChatSessionDO{ID: id}).Error
}

func NewChatSessionID() string {
	return fmt.Sprintf("chat_%d", time.Now().UnixNano())
}
