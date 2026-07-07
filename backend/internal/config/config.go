package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Redis    RedisConfig
	JWT      JWTConfig
	Encrypt  EncryptConfig
	Upload   UploadConfig
	Log      LogConfig
}

type ServerConfig struct {
	Port string
	Mode string
}

type DatabaseConfig struct {
	Host            string
	Port            string
	User            string
	Password        string
	DBName          string
	SSLMode         string
	Timezone        string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
}

func (d DatabaseConfig) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s TimeZone=%s",
		d.Host, d.Port, d.User, d.Password, d.DBName, d.SSLMode, d.Timezone,
	)
}

type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
	PoolSize int
}

func (r RedisConfig) Addr() string {
	return fmt.Sprintf("%s:%s", r.Host, r.Port)
}

type JWTConfig struct {
	AccessSecret  string
	RefreshSecret string
	AccessTTL     time.Duration
	RefreshTTL    time.Duration
}

type EncryptConfig struct {
	AESKey string
	Pepper string
}

type UploadConfig struct {
	MaxSizeMB int
	Path      string
}

type LogConfig struct {
	Level  string
	Output string
}

func Load() (*Config, error) {
	v := viper.New()

	// 读取 .env 文件
	v.SetConfigFile(".env")
	v.SetConfigType("env")
	_ = v.ReadInConfig() // 文件不存在也不报错，用默认值 + 环境变量

	v.AutomaticEnv()

	// Server
	v.SetDefault("SERVER_PORT", "8080")
	v.SetDefault("SERVER_MODE", "debug")

	// Database
	v.SetDefault("DB_HOST", "127.0.0.1")
	v.SetDefault("DB_PORT", "5432")
	v.SetDefault("DB_USER", "postgres")
	v.SetDefault("DB_PASSWORD", "")
	v.SetDefault("DB_NAME", "employment_db")
	v.SetDefault("DB_SSLMODE", "disable")
	v.SetDefault("DB_TIMEZONE", "Asia/Shanghai")
	v.SetDefault("DB_MAX_OPEN_CONNS", 50)
	v.SetDefault("DB_MAX_IDLE_CONNS", 10)
	v.SetDefault("DB_CONN_MAX_LIFETIME", "30m")

	// Redis
	v.SetDefault("REDIS_HOST", "127.0.0.1")
	v.SetDefault("REDIS_PORT", "6379")
	v.SetDefault("REDIS_PASSWORD", "")
	v.SetDefault("REDIS_DB", 0)
	v.SetDefault("REDIS_POOL_SIZE", 20)

	// JWT
	v.SetDefault("JWT_ACCESS_SECRET", "dev-access-secret-change-in-production")
	v.SetDefault("JWT_REFRESH_SECRET", "dev-refresh-secret-change-in-production")
	v.SetDefault("JWT_ACCESS_TTL", "15m")
	v.SetDefault("JWT_REFRESH_TTL", "168h")

	// Encrypt
	v.SetDefault("AES_KEY", "0123456789abcdef0123456789abcdef")
	v.SetDefault("PEPPER_SECRET", "dev-pepper-change-in-production!")

	// Upload
	v.SetDefault("UPLOAD_MAX_SIZE_MB", 10)
	v.SetDefault("UPLOAD_PATH", "./data/uploads")

	// Log
	v.SetDefault("LOG_LEVEL", "info")
	v.SetDefault("LOG_OUTPUT", "stdout")

	accessTTL, err := time.ParseDuration(v.GetString("JWT_ACCESS_TTL"))
	if err != nil { accessTTL = 15 * time.Minute }
	refreshTTL, err := time.ParseDuration(v.GetString("JWT_REFRESH_TTL"))
	if err != nil { refreshTTL = 168 * time.Hour }
	connMaxLifetime, err := time.ParseDuration(v.GetString("DB_CONN_MAX_LIFETIME"))
	if err != nil { connMaxLifetime = 30 * time.Minute }

	return &Config{
		Server: ServerConfig{
			Port: v.GetString("SERVER_PORT"),
			Mode: v.GetString("SERVER_MODE"),
		},
		Database: DatabaseConfig{
			Host:            v.GetString("DB_HOST"),
			Port:            v.GetString("DB_PORT"),
			User:            v.GetString("DB_USER"),
			Password:        v.GetString("DB_PASSWORD"),
			DBName:          v.GetString("DB_NAME"),
			SSLMode:         v.GetString("DB_SSLMODE"),
			Timezone:        v.GetString("DB_TIMEZONE"),
			MaxOpenConns:    v.GetInt("DB_MAX_OPEN_CONNS"),
			MaxIdleConns:    v.GetInt("DB_MAX_IDLE_CONNS"),
			ConnMaxLifetime: connMaxLifetime,
		},
		Redis: RedisConfig{
			Host:     v.GetString("REDIS_HOST"),
			Port:     v.GetString("REDIS_PORT"),
			Password: v.GetString("REDIS_PASSWORD"),
			DB:       v.GetInt("REDIS_DB"),
			PoolSize: v.GetInt("REDIS_POOL_SIZE"),
		},
		JWT: JWTConfig{
			AccessSecret:  v.GetString("JWT_ACCESS_SECRET"),
			RefreshSecret: v.GetString("JWT_REFRESH_SECRET"),
			AccessTTL:     accessTTL,
			RefreshTTL:    refreshTTL,
		},
		Encrypt: EncryptConfig{
			AESKey: v.GetString("AES_KEY"),
			Pepper: v.GetString("PEPPER_SECRET"),
		},
		Upload: UploadConfig{
			MaxSizeMB: v.GetInt("UPLOAD_MAX_SIZE_MB"),
			Path:      v.GetString("UPLOAD_PATH"),
		},
		Log: LogConfig{
			Level:  v.GetString("LOG_LEVEL"),
			Output: v.GetString("LOG_OUTPUT"),
		},
	}, nil
}
