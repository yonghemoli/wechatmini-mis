package user

import (
	"net/http"
	"yonghemolimis/src/apps/api/response"
	"yonghemolimis/src/middlewares"
	"yonghemolimis/src/pkgs/session"
	"yonghemolimis/src/settings"

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

	admin := settings.Conf.Admin
	if admin == nil || admin.Username == "" || admin.Password == "" {
		response.Error(c, http.StatusInternalServerError, "内部管理员账号未配置")
		return
	}
	if req.Username != admin.Username || req.Password != admin.Password {
		response.Error(c, http.StatusUnauthorized, "账号或密码错误")
		return
	}

	sid := session.Create(1, admin.Username, admin.Email, "", true, nil)
	c.SetCookie(middlewares.SessionCookieName, sid, 86400*7, "/", "", false, true)

	response.OK(c, gin.H{"user": adminUserPayload()})
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

func Logout(c *gin.Context) {
	sid, _ := c.Cookie(middlewares.SessionCookieName)
	if sid != "" {
		session.Destroy(sid)
	}
	c.SetCookie(middlewares.SessionCookieName, "", -1, "/", "", false, true)
	response.OKMsg(c, "已登出")
}

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

func adminUserPayload() gin.H {
	admin := settings.Conf.Admin
	return gin.H{
		"id":           1,
		"username":     admin.Username,
		"email":        admin.Email,
		"avatar":       "",
		"isSuperAdmin": true,
		"roleId":       nil,
	}
}
