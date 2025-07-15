package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
)

const (
	serverAddress = "SERVER_ADDRESS"
	jwtSecret     = "JWT_SECRET"
	jwtExpiration = "JWT_EXPIRATION"
	csrfKey       = "CSRF_KEY"
	dbHost        = "DB_HOST"
	dbUser        = "DB_USER"
	dbPass        = "DB_PASS"
	dbPort        = "DB_PORT"
	dbName        = "DB_NAME"
	migrationPath = "MIGRATION_PATH"
)

type Config struct {
	ServerAddress string
	DatabaseURL   string
	JWTSecret     string
	JWTExpiration string
	CSRFKey       string
	MigrationPath string
}

func GetConfig(logger *zerolog.Logger) *Config {
	err := godotenv.Load()
	if err != nil {
		logger.Fatal().Err(err).Msg("Ошибка загрузки env файла")
	}

	getRequiredEnv := func(key string) string {
		value := os.Getenv(key)
		if value == "" {
			logger.Fatal().Msgf("Отсутствует обязательная переменная окружения: %s", key)
		}
		return value
	}

	getOptionalEnv := func(key, defaultValue string) string {
		value := os.Getenv(key)
		if value == "" {
			return defaultValue
		}
		return value
	}

	serverAddress := getOptionalEnv(serverAddress, "127.0.0.1:8030")
	jwtSecret := getRequiredEnv(jwtSecret)
	jwtExpiration := getRequiredEnv(jwtExpiration)
	csrfKey := getRequiredEnv(csrfKey)

	dbDsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		getRequiredEnv(dbHost),
		getRequiredEnv(dbPort),
		getRequiredEnv(dbUser),
		getRequiredEnv(dbPass),
		getRequiredEnv(dbName))

	return &Config{
		ServerAddress: serverAddress,
		DatabaseURL:   dbDsn,
		JWTSecret:     jwtSecret,
		JWTExpiration: jwtExpiration,
		CSRFKey:       csrfKey,
		MigrationPath: getRequiredEnv(migrationPath),
	}
}
