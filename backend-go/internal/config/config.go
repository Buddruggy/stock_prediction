package config

import (
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

// Config 应用配置
type Config struct {
	Port     string
	LogLevel string
	Cache    CacheConfig
	API      APIConfig
}

// CacheConfig 缓存配置
type CacheConfig struct {
	Duration time.Duration
}

// APIConfig API配置
type APIConfig struct {
	Timeout time.Duration
}

// Load 加载配置
func Load() *Config {
	// 加载 .env 文件（如果存在）
	_ = godotenv.Load()

	// 设置默认值
	config := &Config{
		Port:     getEnv("PORT", "8000"),
		LogLevel: getEnv("LOG_LEVEL", "info"),
		Cache: CacheConfig{
			Duration: getDurationEnv("CACHE_DURATION", 5*time.Minute),
		},
		API: APIConfig{
			Timeout: getDurationEnv("API_TIMEOUT", 30*time.Second),
		},
	}

	return config
}

// getEnv 获取环境变量，如果不存在则返回默认值
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getDurationEnv 获取持续时间环境变量
func getDurationEnv(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}

// getIntEnv 获取整数环境变量
func getIntEnv(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}
