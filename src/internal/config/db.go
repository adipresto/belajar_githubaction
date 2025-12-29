package config

import "os"

type db struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
	SSL      string
}

func loadDatabaseConfig() db {
	return db{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		Database: os.Getenv("DB_NAME"),
		SSL:      os.Getenv("DB_SSL_MODE"),
	}
}
