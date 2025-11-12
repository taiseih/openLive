package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

// Config はアプリケーションの設定を保持します
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Firebase FirebaseConfig
	CORS     CORSConfig
}

// ServerConfig はサーバー設定を保持します
type ServerConfig struct {
	Port string
	Env  string
}

// DatabaseConfig はデータベース設定を保持します
type DatabaseConfig struct {
	URL      string
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
}

// FirebaseConfig はFirebase設定を保持します
type FirebaseConfig struct {
	ProjectID           string
	AuthEmulatorHost    string
	CredentialsFilePath string
}

// CORSConfig はCORS設定を保持します
type CORSConfig struct {
	AllowOrigins string
}

// Load は環境変数から設定を読み込みます
func Load() (*Config, error) {
	// .envファイルが存在すれば読み込む（存在しなくてもエラーにしない）
	_ = godotenv.Load()

	config := &Config{
		Server: ServerConfig{
			Port: getEnv("PORT", "8080"),
			Env:  getEnv("ENV", "development"),
		},
		Database: DatabaseConfig{
			URL:      getEnv("DATABASE_URL", ""),
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "openlive"),
			Password: getEnv("DB_PASSWORD", "openlive"),
			Name:     getEnv("DB_NAME", "openlive_db"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		Firebase: FirebaseConfig{
			ProjectID:           getEnv("FIREBASE_PROJECT_ID", "openlive-dev"),
			AuthEmulatorHost:    getEnv("FIREBASE_AUTH_EMULATOR_HOST", ""),
			CredentialsFilePath: getEnv("GOOGLE_APPLICATION_CREDENTIALS", ""),
		},
		CORS: CORSConfig{
			AllowOrigins: getEnv("CORS_ALLOW_ORIGINS", "http://localhost:3000"),
		},
	}

	// DATABASE_URLが空の場合は個別の設定から構築
	if config.Database.URL == "" {
		config.Database.URL = fmt.Sprintf(
			"postgres://%s:%s@%s:%s/%s?sslmode=%s",
			config.Database.User,
			config.Database.Password,
			config.Database.Host,
			config.Database.Port,
			config.Database.Name,
			config.Database.SSLMode,
		)
	}

	return config, nil
}

// getEnv は環境変数を取得し、存在しない場合はデフォルト値を返します
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

