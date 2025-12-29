package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DB         db
	App        app
	Aws        aws
	JwtSecret  string
	Minio      minio
	LOCAL_FILE string
}

func LoadConfig() (*Config, error) {
	_ = godotenv.Load()

	return &Config{
		App:        loadApplicationConfig(),
		Aws:        loadAwsConfig(),
		DB:         loadDatabaseConfig(),
		JwtSecret:  getEnv("JWT_SECRET", "defaultsecret"), // 🔑 baca dari env sekali
		Minio:      loadMinioConfig(),
		LOCAL_FILE: getEnv("FILE_LOCAL", "true"),
	}, nil
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
