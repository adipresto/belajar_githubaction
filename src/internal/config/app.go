package config

import (
	"fmt"
	"os"
)

type app struct {
	Port string
	Env  string
}

func loadApplicationConfig() app {
	return app{
		Port: getAppPort(),
		Env:  getAppEnv(),
	}
}

func getAppEnv() string {
	env := os.Getenv("APP_ENV")
	if env == "" {
		return "development"
	}
	return env
}

func getAppPort() string {
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "3000"
	}
	return fmt.Sprintf("0.0.0.0:%s", port)
}
