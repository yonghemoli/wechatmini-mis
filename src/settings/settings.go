package settings

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"gopkg.in/yaml.v3"
)

const (
	ServiceName        = "xiuxian-analytics"
	ServiceDescription = "修仙游戏数据分析系统"
)

var Conf = new(AppConfig)

type AppConfig struct {
	Name        string          `json:"name"`
	Mode        string          `json:"mode"`
	TrackAPIKey string          `json:"track_api_key"`
	Server      *ServerConfig   `json:"server"`
	Session     *SessionConfig  `json:"session"`
	Log         *LogConfig      `json:"log"`
	GameDB      *DatabaseConfig `json:"game_db"`
	DB          *DatabaseConfig `json:"database"`
	SSO         *SSOConfig      `json:"sso"`
}

type SSOConfig struct {
	AuthBaseURL string `json:"auth_base_url"` // 认证服务地址，如 https://auth.alemonjs.com
}

type ServerConfig struct {
	Host string `json:"host"`
	Port string `json:"port"`
}

type SessionConfig struct {
	ExpiresDays int64 `json:"expires_days"`
}

type LogConfig struct {
	Level    string `json:"level"`
	Filename string `json:"filename"`
}

type DatabaseConfig struct {
	Driver      string `json:"driver"`
	DSN         string `json:"dsn"`
	AutoMigrate bool   `json:"auto_migrate"`
}

// alemonYAML 对应 alemon.config.yaml 的结构
type alemonYAML struct {
	MySQL struct {
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		Database string `yaml:"database"`
	} `yaml:"mysql"`
}

func Init() error {
	log.Println("=== 修仙数据分析系统 配置初始化 ===")

	Conf.Name = getEnv("ANALYTICS_NAME", ServiceName)
	Conf.Mode = getEnv("ANALYTICS_MODE", "release")
	Conf.TrackAPIKey = getEnv("ANALYTICS_TRACK_API_KEY", "")

	Conf.Server = &ServerConfig{
		Host: getEnv("ANALYTICS_SERVER_HOST", "127.0.0.1"),
		Port: getEnv("ANALYTICS_SERVER_PORT", "8080"),
	}

	Conf.Session = &SessionConfig{
		ExpiresDays: getEnvAsInt64("ANALYTICS_SESSION_EXPIRES_DAYS", 30),
	}

	Conf.Log = &LogConfig{
		Level:    getEnv("ANALYTICS_LOG_LEVEL", "info"),
		Filename: getEnv("ANALYTICS_LOG_FILENAME", "work/logs"),
	}

	// ========== 游戏数据库 (MySQL, 只读) ==========
	// 优先从环境变量 ANALYTICS_GAME_DB_DSN 获取；
	// 如果未设置，尝试读取 alemon.config.yaml 自动构建 DSN
	gameDSN := getEnv("ANALYTICS_GAME_DB_DSN", "")
	if gameDSN == "" {
		yamlPath := getEnv("ANALYTICS_GAME_YAML", "")
		if yamlPath != "" {
			if dsn, err := buildGameDSNFromYAML(yamlPath); err != nil {
				log.Printf("[配置] 读取游戏 YAML 配置失败: %v", err)
			} else {
				gameDSN = dsn
				log.Printf("[配置] 已从 %s 加载游戏数据库配置", yamlPath)
			}
		}
	}
	Conf.GameDB = &DatabaseConfig{
		Driver: getEnv("ANALYTICS_GAME_DB_DRIVER", "mysql"),
		DSN:    gameDSN,
	}

	// ========== 业务数据库 (MySQL, 读写) ==========
	Conf.DB = &DatabaseConfig{
		Driver:      getEnv("ANALYTICS_DB_DRIVER", "mysql"),
		DSN:         getEnv("ANALYTICS_DB_DSN", ""),
		AutoMigrate: getEnvAsBool("ANALYTICS_DB_AUTO_MIGRATE", true),
	}

	log.Printf("应用名称: %s", Conf.Name)
	log.Printf("运行模式: %s", Conf.Mode)
	log.Printf("服务地址: http://%s:%s", Conf.Server.Host, Conf.Server.Port)
	log.Printf("游戏数据库驱动: %s (MySQL 只读)", Conf.GameDB.Driver)
	log.Printf("业务数据库驱动: %s (MySQL 读写)", Conf.DB.Driver)
	if Conf.TrackAPIKey != "" {
		log.Printf("Track API Key: %s****", Conf.TrackAPIKey[:4])
	} else {
		log.Println("⚠ Track API Key 未设置，上报接口无认证保护")
	}
	log.Printf("日志级别: %s", Conf.Log.Level)

	// ========== SSO 配置 ==========
	Conf.SSO = &SSOConfig{
		AuthBaseURL: getEnv("ANALYTICS_SSO_AUTH_URL", "https://auth.alemonjs.com"),
	}
	if Conf.SSO.AuthBaseURL != "" {
		log.Printf("SSO 认证地址: %s", Conf.SSO.AuthBaseURL)
	} else {
		log.Println("⚠ SSO 未配置 (ANALYTICS_SSO_AUTH_URL)，将无法通过 SSO 登录")
	}

	log.Println("=== 配置加载完成 ===")
	return nil
}

// buildGameDSNFromYAML 从 alemon.config.yaml 提取 mysql 配置，构建 MySQL DSN
func buildGameDSNFromYAML(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("读取文件失败: %w", err)
	}
	var cfg alemonYAML
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return "", fmt.Errorf("YAML 解析失败: %w", err)
	}
	m := cfg.MySQL
	if m.Host == "" || m.Database == "" {
		return "", fmt.Errorf("YAML 中 mysql 配置不完整")
	}
	// user:password@tcp(host:port)/database?charset=utf8mb4&parseTime=True&loc=Local
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		m.User, m.Password, m.Host, m.Port, m.Database)
	return dsn, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt64(key string, defaultValue int64) int64 {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}
	value, err := strconv.ParseInt(valueStr, 10, 64)
	if err != nil {
		log.Printf("环境变量 %s 解析为 int64 失败，使用默认值 %d: %v\n", key, defaultValue, err)
		return defaultValue
	}
	return value
}

func getEnvAsBool(key string, defaultValue bool) bool {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}
	value, err := strconv.ParseBool(valueStr)
	if err != nil {
		log.Printf("环境变量 %s 解析为 bool 失败，使用默认值 %v: %v\n", key, defaultValue, err)
		return defaultValue
	}
	return value
}

var processRunAt string

func GetProcessRunAT() string {
	if processRunAt == "" {
		processRunAt = time.Now().Format("2006-01-02 15:04:05")
	}
	return processRunAt
}

func LogServerInfo() {
	server := Conf.Server
	log.Println("http://" + server.Host + ":" + server.Port)
}

var Version string
var BuildTime string

func SetBaseInfo(version, buildTime string) {
	Version = version
	BuildTime = buildTime
	log.Println("Version: ", Version, "(", BuildTime, ")")
}

type BaseInfo struct {
	Version   string `json:"version"`
	BuildTime string `json:"build_time"`
}

func GetBaseInfo() BaseInfo {
	return BaseInfo{
		Version:   Version,
		BuildTime: BuildTime,
	}
}
