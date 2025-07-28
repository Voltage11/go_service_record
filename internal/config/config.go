package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type ServerConfig struct {
	ServerHost string
	ServerPort string
}

type DbConfig struct {
	Dsn string
}

type SecretConfig struct {
	HashKey string
	CsrfTokenKey string
}

type LogConfig struct {
	Level string
}

type Config struct {
	Server ServerConfig
	Db DbConfig
	Secret SecretConfig
	Log LogConfig
}

func GetConfig() *Config {
	loadEnv()
	
	postgresDsn := fmt.Sprintf("%s:%s@%s:%s/%s?sslmode=disable",
							getRequiredEnvStr("DB_USER"),
							getRequiredEnvStr("DB_PASS"),
							getRequiredEnvStr("DB_HOST"),
							getRequiredEnvStr("DB_PORT"),
							getRequiredEnvStr("DB_NAME"),)
	
	return &Config{
		Server: ServerConfig{
			ServerHost: getEnvStr("SERVER_HOST", "localhost"),
			ServerPort: getEnvStr("SERVER_PORT", "8080"),
		},
		Db: DbConfig{
			Dsn: postgresDsn,
			},
		Secret: SecretConfig{
			HashKey: getRequiredEnvStr("HASH_KEY"),
			CsrfTokenKey: getRequiredEnvStr("CSRF_TOKEN"),
		},
		Log: LogConfig{
			Level: getEnvStr("LOG_LEVEL", "0"),
		},
	}
	
}

func loadEnv() {
	err := godotenv.Load()
	if err != nil {
		if os.IsNotExist(err) {
			panic(fmt.Sprintf("Файл .env не найден: %v", err))
		} else {
			panic(fmt.Sprintf("Ошибка загрузки .env файла: %v", err))
		}
	}
}

func getEnvStr(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func getRequiredEnvStr(key string) string {
	value := os.Getenv(key)
	if value == "" {
		panic(fmt.Sprintf("Не задан обязательный параметр %s", key))
	}
	return value
}