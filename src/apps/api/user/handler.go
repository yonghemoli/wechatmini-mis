package user

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
	"yonghemolimis/src/apps/api/response"
	"yonghemolimis/src/dao/db"
	gameadmin "yonghemolimis/src/dao/game"
	"yonghemolimis/src/logger"
	"yonghemolimis/src/middlewares"
	"yonghemolimis/src/pkgs/session"
	"yonghemolimis/src/settings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// ==================== SSO 认证 ====================

// SSOConfig 返回 SSO 认证服务地址（前端用于构造跳转 URL）
func SSOConfig(c *gin.Context) {
	if settings.Conf.SSO == nil || settings.Conf.SSO.AuthBaseURL == "" {
		response.Fail(c, "SSO 未配置")
		return
	}
	response.OK(c, gin.H{
		"authURL": settings.Conf.SSO.AuthBaseURL,
	})
}

// SSOLogin 用一次性授权码换取用户信息并创建会话
func SSOLogin(c *gin.Context) {
	var req struct {
		Code string `json:"code" form:"code" binding:"required"`
	}
	if err := c.ShouldBind(&req); err != nil {
		response.Fail(c, "授权码不能为空")
		return
	}

	u, err := exchangeSSOCode(req.Code)
	if err != nil {
		logger.Errorf("[SSO] 交换授权码失败: %v", err)
		response.Error(c, http.StatusUnauthorized, "SSO 认证失败："+err.Error())
		return
	}

	// 第 1 层：仅允许内部用户访问
	if u.UserType != "internal" {
		logger.Warnf("[SSO] 非内部用户尝试访问: uid=%d type=%s", u.ID, u.UserType)
		response.Error(c, http.StatusForbidden, "仅限内部用户访问")
		return
	}

	// 第 2 层：只读校验共享身份库中的管理员
	admin, err := findAdmin(u)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Warnf("[SSO] 共享身份库未找到管理员: ssoId=%d", u.ID)
			response.Error(c, http.StatusForbidden, "未开通分析后台权限，请联系管理员")
			return
		}
		logger.Errorf("[SSO] 共享身份库校验失败: %v", err)
		response.Error(c, http.StatusInternalServerError, "身份校验失败")
		return
	}

	if admin.Status != "active" {
		response.Error(c, http.StatusForbidden, "账户已被禁用，请联系管理员")
		return
	}

	// 允许超管或已分配角色的管理员进入 logap。
	if !admin.IsSuperAdmin && admin.RoleID == nil {
		response.Error(c, http.StatusForbidden, "没有权限，请联系管理员")
		return
	}

	sid := session.Create(admin.ID, u.ID, u.DisplayName, u.Email, u.Avatar, admin.IsSuperAdmin, admin.RoleID)
	c.SetCookie(middlewares.SessionCookieName, sid, 86400*7, "/", "", false, true)

	response.OK(c, gin.H{
		"user": gin.H{
			"id":           admin.ID,
			"username":     u.DisplayName,
			"email":        u.Email,
			"avatar":       u.Avatar,
			"isSuperAdmin": admin.IsSuperAdmin,
			"roleId":       admin.RoleID,
		},
	})
}

// findAdmin 根据 SSO 用户在共享身份库中查找管理员，不做任何写入。
func findAdmin(u *ssoUser) (*gameadmin.Admin, error) {
	if db.GetGame() == nil {
		return nil, fmt.Errorf("游戏数据库未配置，无法校验共享身份")
	}

	var admin gameadmin.Admin
	if err := db.GetGame().Where("sso_user_id = ?", u.ID).First(&admin).Error; err != nil {
		return nil, err
	}

	return &admin, nil
}

// SSOValidateSession 验证 SSO 会话
func SSOValidateSession(c *gin.Context) {
	sid, _ := c.Cookie(middlewares.SessionCookieName)
	if sid == "" {
		response.FailCode(c, 401, "未登录")
		return
	}

	sess := session.Get(sid)
	if sess == nil {
		response.FailCode(c, 401, "登录已过期")
		return
	}

	response.OK(c, gin.H{
		"user": gin.H{
			"id":           sess.AdminID,
			"username":     sess.Username,
			"email":        sess.Email,
			"avatar":       sess.Avatar,
			"isSuperAdmin": sess.IsSuperAdmin,
			"roleId":       sess.RoleID,
		},
	})
}

// Logout 登出（清除本地会话）
func Logout(c *gin.Context) {
	sid, _ := c.Cookie(middlewares.SessionCookieName)
	if sid != "" {
		session.Destroy(sid)
	}
	c.SetCookie(middlewares.SessionCookieName, "", -1, "/", "", false, true)
	response.OKMsg(c, "已登出")
}

// Me 获取当前用户信息
func Me(c *gin.Context) {
	username, _ := c.Get("username")
	isSuperAdmin, _ := c.Get("isSuperAdmin")
	roleID, _ := c.Get("roleID")
	response.OK(c, gin.H{
		"username":     username,
		"isSuperAdmin": isSuperAdmin,
		"roleId":       roleID,
	})
}

// ==================== 内部辅助 ====================

type ssoExchangeResponse struct {
	User struct {
		ID          uint   `json:"id"`
		Name        string `json:"name"`
		Email       string `json:"email"`
		DisplayName string `json:"displayName"`
		Avatar      string `json:"avatar"`
		UserType    string `json:"userType"`
	} `json:"user"`
	AccessToken string `json:"accessToken"`
	ExpiresIn   int    `json:"expiresIn"`
}

type ssoUser struct {
	ID          uint
	Name        string
	Email       string
	DisplayName string
	Avatar      string
	UserType    string
}

func exchangeSSOCode(code string) (*ssoUser, error) {
	if settings.Conf.SSO == nil || settings.Conf.SSO.AuthBaseURL == "" {
		return nil, fmt.Errorf("SSO 未配置")
	}

	exchangeURL := settings.Conf.SSO.AuthBaseURL + "/api/sso/exchange"
	payload := fmt.Sprintf(`{"code":"%s"}`, code)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Post(exchangeURL, "application/json", io.NopCloser(strings.NewReader(payload)))
	if err != nil {
		return nil, fmt.Errorf("请求认证服务失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("认证服务返回 %d", resp.StatusCode)
	}

	var result ssoExchangeResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	if result.User.ID == 0 {
		return nil, fmt.Errorf("无效的用户信息")
	}

	displayName := result.User.DisplayName
	if displayName == "" {
		displayName = result.User.Name
	}

	return &ssoUser{
		ID:          result.User.ID,
		Name:        result.User.Name,
		Email:       result.User.Email,
		DisplayName: displayName,
		Avatar:      result.User.Avatar,
		UserType:    result.User.UserType,
	}, nil
}
