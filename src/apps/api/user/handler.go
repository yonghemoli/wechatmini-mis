package user

import (
	"net/http"
	"strconv"

	"yonghemolimis/src/apps/api/response"
	"yonghemolimis/src/dao/db"
	"yonghemolimis/src/middlewares"
	"yonghemolimis/src/pkgs/session"

	"github.com/gin-gonic/gin"
)

func Login(c *gin.Context) {
	var req struct {
		Username string `json:"username" form:"username" binding:"required"`
		Password string `json:"password" form:"password" binding:"required"`
	}
	if err := c.ShouldBind(&req); err != nil {
		response.Fail(c, "账号和密码不能为空")
		return
	}

	admin, err := db.VerifyAdminCredentials(req.Username, req.Password)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, "账号或密码错误")
		return
	}
	_ = db.UpdateAdminLastLogin(admin.ID)

	sid := session.Create(admin.ID, admin.Username, admin.Email, "", admin.IsSuperAdmin, admin.RoleID)
	c.SetCookie(middlewares.SessionCookieName, sid, 86400*7, "/", "", false, true)

	response.OK(c, gin.H{"user": adminPayload(admin)})
}

func Session(c *gin.Context) {
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

	admin, err := db.GetAdminByID(sess.AdminID)
	if err != nil || admin.Status != db.AdminStatusActive {
		session.Destroy(sid)
		c.SetCookie(middlewares.SessionCookieName, "", -1, "/", "", false, true)
		response.FailCode(c, 401, "账号不可用，请重新登录")
		return
	}
	response.OK(c, gin.H{"user": adminPayload(admin)})
}

func Logout(c *gin.Context) {
	sid, _ := c.Cookie(middlewares.SessionCookieName)
	if sid != "" {
		session.Destroy(sid)
	}
	c.SetCookie(middlewares.SessionCookieName, "", -1, "/", "", false, true)
	response.OKMsg(c, "已登出")
}

func Me(c *gin.Context) {
	userID, ok := c.Get("userID")
	adminID, okID := userID.(uint)
	if !ok || !okID {
		response.Error(c, http.StatusUnauthorized, "未登录或会话已过期")
		return
	}
	admin, err := db.GetAdminByID(adminID)
	if err != nil || admin.Status != db.AdminStatusActive {
		response.Error(c, http.StatusUnauthorized, "账号不可用，请重新登录")
		return
	}
	response.OK(c, gin.H{"user": adminPayload(admin)})
}

func ListAccounts(c *gin.Context) {
	if !requireSuperAdmin(c) {
		return
	}
	rows, err := db.ListAdmins()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.OK(c, gin.H{"list": rows})
}

func CreateAccount(c *gin.Context) {
	if !requireSuperAdmin(c) {
		return
	}
	var req struct {
		Username     string `json:"username" binding:"required"`
		Password     string `json:"password" binding:"required"`
		Name         string `json:"name" binding:"required"`
		Email        string `json:"email" binding:"required"`
		RoleID       *int64 `json:"roleId"`
		IsSuperAdmin bool   `json:"isSuperAdmin"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, "参数错误")
		return
	}
	admin, err := db.CreateAdmin(db.AdminCreateInput{
		Username:     req.Username,
		Password:     req.Password,
		Name:         req.Name,
		Email:        req.Email,
		RoleID:       req.RoleID,
		IsSuperAdmin: req.IsSuperAdmin,
	})
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(c, gin.H{"user": adminPayload(admin)})
}

func UpdateAccount(c *gin.Context) {
	if !requireSuperAdmin(c) {
		return
	}
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		response.Fail(c, "账户ID错误")
		return
	}
	var req struct {
		Name         string `json:"name"`
		Email        string `json:"email"`
		RoleID       *int64 `json:"roleId"`
		IsSuperAdmin *bool  `json:"isSuperAdmin"`
		Status       string `json:"status"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, "参数错误")
		return
	}
	admin, err := db.UpdateAdmin(uint(id), db.AdminUpdateInput{
		Name:         req.Name,
		Email:        req.Email,
		RoleID:       req.RoleID,
		IsSuperAdmin: req.IsSuperAdmin,
		Status:       req.Status,
	})
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(c, gin.H{"user": adminPayload(admin)})
}

func ResetAccountPassword(c *gin.Context) {
	if !requireSuperAdmin(c) {
		return
	}
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		response.Fail(c, "账户ID错误")
		return
	}
	var req struct {
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, "参数错误")
		return
	}
	if err := db.ResetAdminPassword(uint(id), req.Password); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	response.OKMsg(c, "密码已重置")
}

func DisableAccount(c *gin.Context) {
	updateAccountStatus(c, db.AdminStatusBlocked)
}

func EnableAccount(c *gin.Context) {
	updateAccountStatus(c, db.AdminStatusActive)
}

func updateAccountStatus(c *gin.Context, status string) {
	if !requireSuperAdmin(c) {
		return
	}
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		response.Fail(c, "账户ID错误")
		return
	}
	if err := db.UpdateAdminStatus(uint(id), status); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	response.OKMsg(c, "状态已更新")
}

func requireSuperAdmin(c *gin.Context) bool {
	isSuperAdmin, _ := c.Get("isSuperAdmin")
	if ok, _ := isSuperAdmin.(bool); ok {
		return true
	}
	response.Error(c, http.StatusForbidden, "仅超级管理员可操作")
	return false
}

func adminPayload(admin *db.AdminDO) gin.H {
	return gin.H{
		"id":           admin.ID,
		"username":     admin.Username,
		"name":         admin.Name,
		"email":        admin.Email,
		"avatar":       "",
		"isSuperAdmin": admin.IsSuperAdmin,
		"roleId":       admin.RoleID,
		"status":       admin.Status,
		"lastLoginAt":  admin.LastLoginAt,
	}
}
