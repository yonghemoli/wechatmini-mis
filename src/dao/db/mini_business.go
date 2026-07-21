package db

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"
)

// NormalizeJSONObjectList 将历史上被错误存成对象的 JSON 列表恢复为数组。
// requiredKey 过滤掉表单控制字段等非业务条目，避免再次污染阿姨档案。
func NormalizeJSONObjectList(value interface{}, requiredKey string) []interface{} {
	var candidates []interface{}
	switch typed := value.(type) {
	case []interface{}:
		candidates = typed
	case map[string]interface{}:
		for _, item := range typed {
			candidates = append(candidates, item)
		}
	case string:
		var decoded interface{}
		if json.Unmarshal([]byte(typed), &decoded) != nil {
			return []interface{}{}
		}
		return NormalizeJSONObjectList(decoded, requiredKey)
	default:
		return []interface{}{}
	}
	result := make([]interface{}, 0, len(candidates))
	for _, item := range candidates {
		object, ok := item.(map[string]interface{})
		if !ok {
			continue
		}
		if value, exists := object[requiredKey]; !exists || strings.TrimSpace(fmt.Sprint(value)) == "" {
			continue
		}
		result = append(result, object)
	}
	return result
}

type CaregiverListQuery struct {
	ServiceID          string
	Keyword            string
	AvailabilityStatus string
	Recommended        *bool
	Status             string
	Page               int
	PageSize           int
}

type BusinessListQuery struct {
	Status   string
	Keyword  string
	Page     int
	PageSize int
}

func ListMiniServiceCategories(enabled *bool) ([]MiniServiceCategoryDO, error) {
	var rows []MiniServiceCategoryDO
	query := Get().Model(&MiniServiceCategoryDO{})
	if enabled != nil {
		query = query.Where("enabled = ?", *enabled)
	}
	err := query.Order("sort ASC, id ASC").Find(&rows).Error
	return rows, err
}

func GetMiniServiceCategory(id string, onlyEnabled bool) (*MiniServiceCategoryDO, error) {
	var row MiniServiceCategoryDO
	query := Get().Where("id = ?", id)
	if onlyEnabled {
		query = query.Where("enabled = ?", true)
	}
	if err := query.First(&row).Error; err != nil {
		return nil, err
	}
	return &row, nil
}

func SaveMiniServiceCategory(row *MiniServiceCategoryDO) error {
	if strings.TrimSpace(row.ID) == "" {
		return errors.New("服务分类 ID 不能为空")
	}
	var existing MiniServiceCategoryDO
	if err := Get().First(&existing, "id = ?", row.ID).Error; err == nil {
		row.CreatedAt = existing.CreatedAt
	}
	return Get().Save(row).Error
}

func DeleteMiniServiceCategory(id string) error {
	var used int64
	if err := Get().Model(&CaregiverDO{}).Where("service_ids LIKE ?", "%\""+id+"\"%").Count(&used).Error; err != nil {
		return err
	}
	if used > 0 {
		return errors.New("该分类已被服务人员使用，请停用而不是删除")
	}
	return Get().Delete(&MiniServiceCategoryDO{}, "id = ?", id).Error
}

func ListCaregivers(q CaregiverListQuery) ([]CaregiverDO, int64, error) {
	query := Get().Model(&CaregiverDO{})
	if q.Status != "" {
		query = query.Where("status = ?", q.Status)
	}
	if q.Recommended != nil {
		query = query.Where("recommended = ?", *q.Recommended)
	}
	if q.ServiceID != "" {
		query = query.Where("service_ids LIKE ?", "%\""+q.ServiceID+"\"%")
	}
	if q.AvailabilityStatus != "" {
		query = query.Where("availability_status = ?", q.AvailabilityStatus)
	}
	if keyword := strings.TrimSpace(q.Keyword); keyword != "" {
		like := "%" + keyword + "%"
		query = query.Where("name LIKE ? OR origin LIKE ? OR jobs LIKE ?", like, like, like)
	}
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if q.Page < 1 {
		q.Page = 1
	}
	if q.PageSize < 1 {
		q.PageSize = 20
	}
	var rows []CaregiverDO
	err := query.Order("sort DESC, updated_at DESC, id ASC").
		Limit(q.PageSize).Offset((q.Page - 1) * q.PageSize).Find(&rows).Error
	return rows, total, err
}

func GetCaregiver(id string, onlyCompleted bool) (*CaregiverDO, error) {
	var row CaregiverDO
	query := Get().Where("id = ?", id)
	if onlyCompleted {
		query = query.Where("status = ?", "COMPLETED")
	}
	if err := query.First(&row).Error; err != nil {
		return nil, err
	}
	return &row, nil
}

func SaveCaregiver(row *CaregiverDO) error {
	if row.ID == "" {
		row.ID = fmt.Sprintf("CG%d", time.Now().UnixNano())
	} else {
		var existing CaregiverDO
		if err := Get().First(&existing, "id = ?", row.ID).Error; err == nil {
			row.CreatedAt = existing.CreatedAt
		}
	}
	return Get().Save(row).Error
}

func DeleteCaregiver(id string) error { return Get().Delete(&CaregiverDO{}, "id = ?", id).Error }

func UpdateCaregiverStatus(id, status string) error {
	return Get().Model(&CaregiverDO{}).Where("id = ?", id).Update("status", status).Error
}

func FindDemandByIdempotency(scope, key string, since time.Time) (*DemandDO, error) {
	if key == "" {
		return nil, gorm.ErrRecordNotFound
	}
	var row DemandDO
	err := Get().Where("submission_scope = ? AND idempotency_key = ? AND created_at >= ?", scope, key, since).
		Order("created_at ASC").First(&row).Error
	return &row, err
}

func CreateDemand(row *DemandDO) error { return Get().Create(row).Error }

func FindRecentDemandDuplicate(phone, serviceID, requirements string, since time.Time) (*DemandDO, error) {
	var row DemandDO
	err := Get().Where("contact_phone = ? AND service_id = ? AND requirements = ? AND created_at >= ?", phone, serviceID, requirements, since).
		Order("created_at ASC").First(&row).Error
	return &row, err
}

func FindResumeByIdempotency(scope, key string, since time.Time) (*ResumeDO, error) {
	if key == "" {
		return nil, gorm.ErrRecordNotFound
	}
	var row ResumeDO
	err := Get().Where("submission_scope = ? AND idempotency_key = ? AND created_at >= ?", scope, key, since).
		Order("created_at ASC").First(&row).Error
	return &row, err
}

func CreateResumeAndCaregiverDraft(resume *ResumeDO, draft *CaregiverDO) error {
	return Get().Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(resume).Error; err != nil {
			return err
		}
		return tx.Create(draft).Error
	})
}

func ListDemands(q BusinessListQuery) ([]DemandDO, int64, error) {
	query := Get().Model(&DemandDO{})
	if q.Status != "" {
		query = query.Where("status = ?", q.Status)
	}
	if keyword := strings.TrimSpace(q.Keyword); keyword != "" {
		like := "%" + keyword + "%"
		query = query.Where("id LIKE ? OR contact_phone LIKE ? OR service_name LIKE ? OR caregiver_name LIKE ?", like, like, like, like)
	}
	return pageDemands(query, q)
}

func pageDemands(query *gorm.DB, q BusinessListQuery) ([]DemandDO, int64, error) {
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	q = normalizeBusinessQuery(q)
	var rows []DemandDO
	err := query.Order("created_at DESC").Limit(q.PageSize).Offset((q.Page - 1) * q.PageSize).Find(&rows).Error
	return rows, total, err
}

func ListResumes(q BusinessListQuery) ([]ResumeDO, int64, error) {
	query := Get().Model(&ResumeDO{})
	if q.Status != "" {
		query = query.Where("status = ?", q.Status)
	}
	if keyword := strings.TrimSpace(q.Keyword); keyword != "" {
		like := "%" + keyword + "%"
		query = query.Where("id LIKE ? OR contact_phone LIKE ? OR service_name LIKE ?", like, like, like)
	}
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	q = normalizeBusinessQuery(q)
	var rows []ResumeDO
	err := query.Order("created_at DESC").Limit(q.PageSize).Offset((q.Page - 1) * q.PageSize).Find(&rows).Error
	return rows, total, err
}

func normalizeBusinessQuery(q BusinessListQuery) BusinessListQuery {
	if q.Page < 1 {
		q.Page = 1
	}
	if q.PageSize < 1 || q.PageSize > 100 {
		q.PageSize = 20
	}
	return q
}

func UpdateDemandStatus(id, status, note string, operatorID uint) (*DemandDO, error) {
	var row DemandDO
	err := Get().Transaction(func(tx *gorm.DB) error {
		if err := tx.First(&row, "id = ?", id).Error; err != nil {
			return err
		}
		from := row.Status
		if !validDemandTransition(from, status) {
			return fmt.Errorf("不允许从 %s 变更为 %s", from, status)
		}
		if err := tx.Model(&row).Updates(map[string]any{"status": status, "updated_at": time.Now()}).Error; err != nil {
			return err
		}
		row.Status = status
		return tx.Create(&BusinessStatusHistoryDO{EntityType: "DEMAND", EntityID: id, FromStatus: from, ToStatus: status, OperatorID: operatorID, Note: note}).Error
	})
	return &row, err
}

func UpdateResumeStatus(id, status, note string, operatorID uint) (*ResumeDO, error) {
	var row ResumeDO
	err := Get().Transaction(func(tx *gorm.DB) error {
		if err := tx.First(&row, "id = ?", id).Error; err != nil {
			return err
		}
		from := row.Status
		if !validResumeTransition(from, status) {
			return fmt.Errorf("不允许从 %s 变更为 %s", from, status)
		}
		if err := tx.Model(&row).Updates(map[string]any{"status": status, "updated_at": time.Now()}).Error; err != nil {
			return err
		}
		row.Status = status
		return tx.Create(&BusinessStatusHistoryDO{EntityType: "RESUME", EntityID: id, FromStatus: from, ToStatus: status, OperatorID: operatorID, Note: note}).Error
	})
	return &row, err
}

func validDemandTransition(from, to string) bool {
	allowed := map[string][]string{
		"PENDING_CONTACT": {"CONTACTED", "CLOSED"},
		"CONTACTED":       {"MATCHING", "CLOSED"},
		"MATCHING":        {"CLOSED"},
	}
	return containsStatus(allowed[from], to)
}

func validResumeTransition(from, to string) bool {
	allowed := map[string][]string{
		"PENDING_CONTACT": {"CONTACTED", "CLOSED"},
		"CONTACTED":       {"VERIFYING", "CLOSED"},
		"VERIFYING":       {"APPROVED", "REJECTED", "CLOSED"},
	}
	return containsStatus(allowed[from], to)
}

func containsStatus(values []string, target string) bool {
	for _, value := range values {
		if value == target {
			return true
		}
	}
	return false
}

func ListBusinessStatusHistory(entityType, entityID string) ([]BusinessStatusHistoryDO, error) {
	var rows []BusinessStatusHistoryDO
	err := Get().Where("entity_type = ? AND entity_id = ?", entityType, entityID).Order("created_at ASC, id ASC").Find(&rows).Error
	return rows, err
}

func AssignDemand(id string, adminID uint) (*DemandDO, error) {
	if err := Get().Model(&DemandDO{}).Where("id = ?", id).Updates(map[string]any{"assigned_admin_id": adminID, "updated_at": time.Now()}).Error; err != nil {
		return nil, err
	}
	var row DemandDO
	if err := Get().First(&row, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &row, nil
}

func AssignResume(id string, adminID uint) (*ResumeDO, error) {
	if err := Get().Model(&ResumeDO{}).Where("id = ?", id).Updates(map[string]any{"assigned_admin_id": adminID, "updated_at": time.Now()}).Error; err != nil {
		return nil, err
	}
	var row ResumeDO
	if err := Get().First(&row, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &row, nil
}
