package main

import (
	"yonghemolimis/src/dao/db"
	"yonghemolimis/src/logger"
	"yonghemolimis/src/middlewares"
	"yonghemolimis/src/pkgs/oss"
	"yonghemolimis/src/route"
	"yonghemolimis/src/settings"

	"github.com/gin-gonic/gin"
)

func main() {
	// 1. 初始化配置
	settings.Init()

	// 2. 初始化日志
	logger.Init()
	logger.Info("系统启动中...")

	// 3. 初始化数据库
	if err := db.Init(); err != nil {
		logger.Errorf("数据库初始化失败: %v", err)
		panic("数据库初始化失败: " + err.Error())
	}
	if err := oss.Init(&oss.Config{
		Endpoint: settings.Conf.OSS.Endpoint, AccessKey: settings.Conf.OSS.AccessKey, SecretKey: settings.Conf.OSS.SecretKey,
		UseSSL: settings.Conf.OSS.UseSSL, Bucket: settings.Conf.OSS.Bucket, PublicBaseURL: settings.Conf.OSS.PublicBaseURL,
	}); err != nil {
		logger.Warnf("OSS 初始化失败，上传功能不可用: %v", err)
	}

	// 4. 设置 Gin
	if settings.Conf.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()
	r.Use(middlewares.CORS())

	// 5. 注册 API 路由
	route.SetupRoutes(r)

	// 7. 启动服务器
	addr := ":" + settings.Conf.Server.Port
	logger.Infof("服务启动在 %s", addr)
	if err := r.Run(addr); err != nil {
		logger.Errorf("服务启动失败: %v", err)
	}
}
