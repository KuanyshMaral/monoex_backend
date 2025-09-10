package config

import (
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	"log"
	"os"
	"strconv"
)

// Основная структура конфига
type Config struct {
	Server ServerConfig   `mapstructure:"server" yaml:"server"`
	DB     DatabaseConfig `mapstructure:"database" yaml:"database"`
}

type ServerConfig struct {
	Host         string `mapstructure:"host" yaml:"host"`
	Port         string `mapstructure:"port" yaml:"port"`
	ReadTimeout  int    `mapstructure:"read_timeout" yaml:"read_timeout"`
	WriteTimeout int    `mapstructure:"write_timeout" yaml:"write_timeout"`
	IdleTimeout  int    `mapstructure:"idle_timeout" yaml:"idle_timeout"`
}

type DatabaseConfig struct {
	Driver          string `mapstructure:"driver" yaml:"driver"`
	DSN             string `mapstructure:"url" yaml:"url"`
	MaxOpenConns    int    `mapstructure:"max_open_conns" yaml:"max_open_conns"`
	MaxIdleConns    int    `mapstructure:"max_idle_conns" yaml:"max_idle_conns"`
	ConnMaxLifetime int    `mapstructure:"conn_max_lifetime" yaml:"conn_max_lifetime"` // в минутах
}

var AppConfig *Config

func Load() (*Config, error) {
	_ = godotenv.Load()

	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath(getEnv("CONFIG_PATH", "./config"))
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		log.Println("⚠️ Could not read config.yaml, relying on env only:", err)
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	applyEnvOverrides(&cfg)

	// 🔹 Разворачиваем переменные в DSN и Port
	cfg.DB.DSN = os.ExpandEnv(cfg.DB.DSN)
	cfg.Server.Port = os.ExpandEnv(cfg.Server.Port)

	AppConfig = &cfg
	log.Println("✅ Config loaded successfully")
	return AppConfig, nil
}

// Получить конфиг (ленивая инициализация)
func GetConfig() *Config {
	if AppConfig == nil {
		if _, err := Load(); err != nil {
			log.Fatalf("Failed to load config: %v", err)
		}
	}
	return AppConfig
}

// --- вспомогательные функции ---
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if v, err := strconv.Atoi(value); err == nil {
			return v
		}
	}
	return defaultValue
}

func applyEnvOverrides(cfg *Config) {
	// Server
	if cfg.Server.Host == "" {
		cfg.Server.Host = getEnv("SERVER_HOST", "0.0.0.0")
	}
	if cfg.Server.Port == "" {
		cfg.Server.Port = getEnv("SERVER_PORT", "8080")
	}
	if cfg.Server.ReadTimeout == 0 {
		cfg.Server.ReadTimeout = getEnvAsInt("SERVER_READ_TIMEOUT", 15)
	}
	if cfg.Server.WriteTimeout == 0 {
		cfg.Server.WriteTimeout = getEnvAsInt("SERVER_WRITE_TIMEOUT", 15)
	}
	if cfg.Server.IdleTimeout == 0 {
		cfg.Server.IdleTimeout = getEnvAsInt("SERVER_IDLE_TIMEOUT", 60)
	}

	// Database
	if cfg.DB.Driver == "" {
		cfg.DB.Driver = getEnv("DB_DRIVER", "postgres")
	}
	if cfg.DB.DSN == "" {
		cfg.DB.DSN = getEnv("DB_URL", "")
	}
	if cfg.DB.MaxOpenConns == 0 {
		cfg.DB.MaxOpenConns = getEnvAsInt("DB_MAX_OPEN_CONNS", 25)
	}
	if cfg.DB.MaxIdleConns == 0 {
		cfg.DB.MaxIdleConns = getEnvAsInt("DB_MAX_IDLE_CONNS", 5)
	}
	if cfg.DB.ConnMaxLifetime == 0 {
		cfg.DB.ConnMaxLifetime = getEnvAsInt("DB_CONN_MAX_LIFETIME", 5)
	}
}
