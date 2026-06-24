package game

import "gorm.io/gorm"

// Admin 映射历史管理员表，仅用于兼容旧查询。
type Admin struct {
	ID           uint           `gorm:"column:id" json:"id"`
	Username     string         `gorm:"column:username" json:"username"`
	RoleID       *int64         `gorm:"column:role_id" json:"roleId"`
	IsSuperAdmin bool           `gorm:"column:is_super_admin" json:"isSuperAdmin"`
	Status       string         `gorm:"column:status" json:"status"`
	DeletedAt    gorm.DeletedAt `gorm:"column:deleted_at" json:"-"`
}

func (Admin) TableName() string { return "admins" }
