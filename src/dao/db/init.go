package db

import (
	"errors"
	"log"
	"yonghemolimis/src/settings"

	"github.com/xiuxianjs/xiuxian-game-db-sdk/dbkit"
	"gorm.io/gorm"
)

var (
	analyticsDB *gorm.DB // 业务库 — 读写
	gameDB      *gorm.DB // 游戏库 — 只读
)

// Get 返回业务数据库连接
func Get() *gorm.DB {
	return analyticsDB
}

// GetGame 返回游戏数据库连接（只读）
func GetGame() *gorm.DB {
	return gameDB
}

// Init 初始化双数据源
func Init() error {
	dbc := settings.Conf.DB
	if dbc == nil || dbc.DSN == "" {
		return errors.New("业务数据库 DSN 未配置 (ANALYTICS_DB_DSN)")
	}

	// 初始化业务数据库
	var err error
	analyticsDB, err = dbkit.Open(dbkit.DatabaseConfig{
		Driver: dbc.Driver,
		DSN:    dbc.DSN,
	})
	if err != nil {
		return err
	}
	log.Println("[DB] 业务数据库连接成功")

	// 自动迁移
	if dbc.AutoMigrate {
		if err := AutoMigrate(); err != nil {
			log.Printf("[DB] 自动迁移失败: %v", err)
		} else {
			log.Println("[DB] 自动迁移完成")
		}
	}

	// 初始化游戏数据库（只读）
	gameConf := settings.Conf.GameDB
	if gameConf != nil && gameConf.DSN != "" {
		gameDB, err = dbkit.Open(dbkit.DatabaseConfig{
			Driver: gameConf.Driver,
			DSN:    gameConf.DSN,
		})
		if err != nil {
			return err
		}
		log.Println("[DB] 游戏数据库连接成功 (只读)")
	} else {
		log.Println("[DB] 游戏数据库未配置，部分功能将不可用")
	}

	return nil
}
