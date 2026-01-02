package config

import (
	"fmt"
	"os"
)

type Config struct {
	Database DatabaseConfig
	Session  SessionConfig
	Server   ServerConfig
}

type DatabaseConfig struct {
	Host     string
	Port     string
	Name     string
	User     string
	Password string
}

type SessionConfig struct {
	Secret string
}

type ServerConfig struct {
	Port string
}

var AppConfig Config

func init() {
	AppConfig = Config{
		Database: DatabaseConfig{
			Host:     getEnvOrPanic("DB_HOST"),
			Port:     getEnvOrPanic("DB_PORT"),
			Name:     getEnvOrPanic("DB_NAME"),
			User:     getEnvOrPanic("DB_USER"),
			Password: getEnvOrPanic("DB_PASSWORD"),
		},
		Session: SessionConfig{
			Secret: getEnvOrPanic("SESSION_SECRET"),
		},
		Server: ServerConfig{
			Port: getEnvOrPanic("PORT"),
		},
	}
}

func getEnvOrPanic(key string) string {
	value := os.Getenv(key)
	if value == "" {
		panic(fmt.Sprintf("%s environment variable is required", key))
	}
	return value
}