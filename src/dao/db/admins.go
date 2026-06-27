package db

import (
	"time"

	"yonghemolimis/src/pkgs/password"

	"gorm.io/gorm"
)

const (
	AdminStatusActive  = "active"
	AdminStatusBlocked = "blocked"
)

type AdminCreateInput struct {
	Username     string
	Password     string
	Name         string
	Email        string
	RoleID       *int64
	IsSuperAdmin bool
}

type AdminUpdateInput struct {
	Name         string
	Email        string
	RoleID       *int64
	IsSuperAdmin *bool
	Status       string
}

func ListAdmins() ([]AdminDO, error) {
	var rows []AdminDO
	err := Get().Order("is_super_admin DESC, id ASC").Find(&rows).Error
	return rows, err
}

func GetAdminByID(id uint) (*AdminDO, error) {
	var admin AdminDO
	if err := Get().First(&admin, id).Error; err != nil {
		return nil, err
	}
	return &admin, nil
}

func GetAdminByUsername(username string) (*AdminDO, error) {
	var admin AdminDO
	if err := Get().Where("username = ?", username).First(&admin).Error; err != nil {
		return nil, err
	}
	return &admin, nil
}

func CreateAdmin(input AdminCreateInput) (*AdminDO, error) {
	hash, err := password.Hash(input.Password)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	admin := &AdminDO{
		Username:     input.Username,
		PasswordHash: hash,
		Name:         input.Name,
		Email:        input.Email,
		RoleID:       input.RoleID,
		IsSuperAdmin: input.IsSuperAdmin,
		Status:       AdminStatusActive,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	if err := Get().Create(admin).Error; err != nil {
		return nil, err
	}
	return admin, nil
}

func UpdateAdmin(id uint, input AdminUpdateInput) (*AdminDO, error) {
	admin, err := GetAdminByID(id)
	if err != nil {
		return nil, err
	}

	updates := map[string]any{
		"name":           input.Name,
		"email":          input.Email,
		"role_id":        input.RoleID,
		"updated_at":     time.Now(),
		"is_super_admin": admin.IsSuperAdmin,
	}
	if input.IsSuperAdmin != nil {
		updates["is_super_admin"] = *input.IsSuperAdmin
	}
	if input.Status != "" {
		updates["status"] = input.Status
	}
	if err := Get().Model(admin).Updates(updates).Error; err != nil {
		return nil, err
	}
	return GetAdminByID(id)
}

func UpdateAdminStatus(id uint, status string) error {
	return Get().Model(&AdminDO{}).Where("id = ?", id).
		Updates(map[string]any{"status": status, "updated_at": time.Now()}).Error
}

func ResetAdminPassword(id uint, passwordText string) error {
	hash, err := password.Hash(passwordText)
	if err != nil {
		return err
	}
	return Get().Model(&AdminDO{}).Where("id = ?", id).
		Updates(map[string]any{"password_hash": hash, "updated_at": time.Now()}).Error
}

func UpdateAdminLastLogin(id uint) error {
	return Get().Model(&AdminDO{}).Where("id = ?", id).
		Updates(map[string]any{"last_login_at": time.Now(), "updated_at": time.Now()}).Error
}

func VerifyAdminCredentials(username, rawPassword string) (*AdminDO, error) {
	admin, err := GetAdminByUsername(username)
	if err != nil {
		return nil, err
	}
	if admin.Status != AdminStatusActive {
		return nil, gorm.ErrRecordNotFound
	}
	if !password.Verify(admin.PasswordHash, rawPassword) {
		return nil, gorm.ErrRecordNotFound
	}
	return admin, nil
}
