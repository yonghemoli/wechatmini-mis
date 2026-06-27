package route

import (
	"yonghemolimis/src/apps/api/chatws"
	miniapi "yonghemolimis/src/apps/api/mini"
	misapi "yonghemolimis/src/apps/api/mis"
	"yonghemolimis/src/apps/api/user"
	"yonghemolimis/src/middlewares"

	"github.com/gin-gonic/gin"
)

// SetupRoutes 注册所有 API 路由
func SetupRoutes(r *gin.Engine) {
	mini := r.Group("/api/mini")
	miniapi.RegisterRoutes(mini)
	chatws.RegisterRoutes(mini, chatws.RoleMini)

	api := r.Group("/api/v1")
	chatws.RegisterRoutes(api, chatws.RoleAdmin)

	// 公开接口
	api.POST("/login", user.Login)
	api.GET("/session", user.Session)
	api.POST("/logout", user.Logout)

	// 需要认证的接口
	auth := api.Group("", middlewares.AuthRequired())
	{
		// 内部账户
		auth.GET("/me", user.Me)
		auth.GET("/admin/accounts", user.ListAccounts)
		auth.POST("/admin/accounts", user.CreateAccount)
		auth.PUT("/admin/accounts/:id", user.UpdateAccount)
		auth.POST("/admin/accounts/:id/reset-password", user.ResetAccountPassword)
		auth.POST("/admin/accounts/:id/disable", user.DisableAccount)
		auth.POST("/admin/accounts/:id/enable", user.EnableAccount)

		auth.GET("/dashboard/summary", misapi.DashboardSummary)
		auth.GET("/dashboard/exceptions", misapi.DashboardExceptions)

		auth.GET("/orders", misapi.ListOrders)
		auth.GET("/orders/export", misapi.ExportOrders)
		auth.GET("/orders/:id", misapi.GetOrder)
		auth.POST("/orders/:id/confirm", misapi.ConfirmOrder)
		auth.POST("/orders/:id/refund", misapi.RefundOrder)
		auth.PUT("/orders/:id/note", misapi.UpdateOrderNote)

		auth.GET("/users", misapi.ListUsers)
		auth.GET("/users/export", misapi.ExportUsers)
		auth.POST("/users/:id/ban", misapi.BanUser)
		auth.POST("/users/:id/unban", misapi.UnbanUser)

		auth.GET("/service-types", misapi.ListServiceTypes)
		auth.POST("/service-types", misapi.CreateServiceType)
		auth.PUT("/service-types/:id", misapi.UpdateServiceType)
		auth.DELETE("/service-types/:id", misapi.DeleteServiceType)
		auth.POST("/service-types/:id/enable", misapi.EnableServiceType)
		auth.POST("/service-types/:id/disable", misapi.DisableServiceType)

		auth.GET("/services", misapi.ListServices)
		auth.POST("/services", misapi.CreateService)
		auth.PUT("/services/:id", misapi.UpdateService)
		auth.DELETE("/services/:id", misapi.DeleteService)
		auth.GET("/services/export", misapi.ExportServices)
		auth.POST("/services/:id/publish", misapi.PublishService)
		auth.POST("/services/:id/unpublish", misapi.UnpublishService)

		auth.GET("/shops", misapi.ListShops)
		auth.POST("/shops", misapi.CreateShop)
		auth.PUT("/shops/:id", misapi.UpdateShop)
		auth.DELETE("/shops/:id", misapi.DeleteShop)
		auth.POST("/shops/:id/open", misapi.OpenShop)
		auth.POST("/shops/:id/close", misapi.CloseShop)

		auth.GET("/faqs", misapi.ListFAQs)
		auth.POST("/faqs", misapi.CreateFAQ)
		auth.PUT("/faqs/:id", misapi.UpdateFAQ)
		auth.DELETE("/faqs/:id", misapi.DeleteFAQ)
		auth.POST("/faqs/:id/publish", misapi.PublishFAQ)
		auth.POST("/faqs/:id/unpublish", misapi.UnpublishFAQ)

		auth.GET("/chat/sessions", misapi.ListChatSessions)
		auth.GET("/chat/sessions/:id/messages", misapi.ListChatMessages)
		auth.POST("/chat/sessions/:id/messages", misapi.CreateChatMessage)
		auth.POST("/chat/sessions/:id/close", misapi.CloseChatSession)
		auth.POST("/chat/sessions/:id/read", misapi.ReadChatSession)
	}
}
