package main

import (
	"embed"
	"io/fs"
	"net/http"
	"yonghemolimis/src/cron"
	"yonghemolimis/src/dao/db"
	"yonghemolimis/src/logger"
	"yonghemolimis/src/middlewares"
	"yonghemolimis/src/route"
	"yonghemolimis/src/settings"

	"github.com/gin-gonic/gin"
)

//go:embed all:dist
var staticFiles embed.FS

func main() {
	// 1. 初始化配置
	settings.Init()

	// 2. 初始化日志
	logger.Init()
	logger.Info("修仙游戏数据分析系统启动中...")

	// 3. 初始化数据库
	if err := db.Init(); err != nil {
		logger.Errorf("数据库初始化失败: %v", err)
		panic("数据库初始化失败: " + err.Error())
	}

	// 4. 设置 Gin
	if settings.Conf.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()
	r.Use(middlewares.CORS())

	// 5. 注册 API 路由
	route.SetupRoutes(r)

	// 6. 静态文件（前端 SPA）
	distFS, err := fs.Sub(staticFiles, "dist")
	if err != nil {
		logger.Warnf("未找到前端静态文件: %v", err)
	} else {
		r.NoRoute(func(c *gin.Context) {
			path := c.Request.URL.Path
			f, e := fs.Stat(distFS, path[1:])
			if e == nil && !f.IsDir() {
				c.FileFromFS(path, http.FS(distFS))
				return
			}
			c.FileFromFS("/", http.FS(distFS))
		})
	}

	// 7. 启动定时任务
	cron.StartScheduler()

	// 8. 启动服务器
	addr := ":" + settings.Conf.Server.Port
	logger.Infof("服务启动在 %s", addr)
	if err := r.Run(addr); err != nil {
		logger.Errorf("服务启动失败: %v", err)
	}
}
