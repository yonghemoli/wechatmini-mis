package db

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

const (
	StatusOpen   = "open"
	StatusClosed = "closed"
)

func ListFAQs(category string) ([]FAQDO, error) {
	var rows []FAQDO
	query := Get().Model(&FAQDO{})
	if category != "" {
		query = query.Where("category = ?", category)
	}
	err := query.Order("sort_order ASC, id DESC").Find(&rows).Error
	return rows, err
}
func ListPublicFAQs(category string) ([]FAQDO, error) {
	var rows []FAQDO
	query := Get().Model(&FAQDO{}).Where("visible = ?", true)
	if category != "" {
		query = query.Where("category = ?", category)
	}
	err := query.Order("sort_order ASC, id ASC").Find(&rows).Error
	return rows, err
}
func CreateFAQ(row *FAQDO) error { return Get().Create(row).Error }
func UpdateFAQ(id uint, input FAQDO) (*FAQDO, error) {
	updates := map[string]any{"question": input.Question, "answer": input.Answer, "category": input.Category, "sort_order": input.SortOrder, "visible": input.Visible, "updated_at": time.Now()}
	if err := Get().Model(&FAQDO{}).Where("id = ?", id).Updates(updates).Error; err != nil {
		return nil, err
	}
	var row FAQDO
	if err := Get().First(&row, id).Error; err != nil {
		return nil, err
	}
	return &row, nil
}
func DeleteFAQ(id uint) error { return Get().Delete(&FAQDO{}, id).Error }
func UpdateFAQVisible(id uint, visible bool) (*FAQDO, error) {
	if err := Get().Model(&FAQDO{}).Where("id = ?", id).Updates(map[string]any{"visible": visible, "updated_at": time.Now()}).Error; err != nil {
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
	row := &ChatMessageDO{SessionID: sessionID, Sender: sender, MsgType: msgType, Content: content, IsRead: sender == "admin", CreatedAt: time.Now()}
	if err := Get().Create(row).Error; err != nil {
		return nil, err
	}
	updates := map[string]any{"last_message": content, "updated_at": time.Now()}
	if sender != "admin" {
		updates["unread_count"] = gorm.Expr("unread_count + 1")
	}
	_ = Get().Model(&ChatSessionDO{}).Where("id = ?", sessionID).Updates(updates).Error
	return row, nil
}
func CloseChatSession(sessionID string) error {
	return Get().Model(&ChatSessionDO{}).Where("id = ?", sessionID).Updates(map[string]any{"status": StatusClosed, "updated_at": time.Now()}).Error
}
func ReadChatSession(sessionID string) error {
	if err := Get().Model(&ChatMessageDO{}).Where("session_id = ?", sessionID).Update("is_read", true).Error; err != nil {
		return err
	}
	return Get().Model(&ChatSessionDO{}).Where("id = ?", sessionID).Updates(map[string]any{"unread_count": 0, "updated_at": time.Now()}).Error
}
func EnsureChatSession(id, userID, userName string) error {
	row := ChatSessionDO{ID: id, UserID: userID, UserName: userName, Status: StatusOpen, CreatedAt: time.Now(), UpdatedAt: time.Now()}
	return Get().FirstOrCreate(&row, ChatSessionDO{ID: id}).Error
}
func NewChatSessionID() string { return fmt.Sprintf("chat_%d", time.Now().UnixNano()) }
