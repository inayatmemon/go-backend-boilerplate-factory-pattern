package env_service_one

import (
	env_constants "go_boilerplate_project/constants/env"
	errors_constants "go_boilerplate_project/constants/errors"
	env_models "go_boilerplate_project/models/env"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

func LoadEnv() (*env_models.Environment, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Failed to load environment variables: %v", err)
	}
	env := &env_models.Environment{}
	envMode := getEnvAsString("ENV_MODE", string(env_constants.EnvModeDevelopment))
	if envMode != string(env_constants.EnvModeDevelopment) && envMode != string(env_constants.EnvModeStaging) && envMode != string(env_constants.EnvModeProduction) {
		return nil, errors_constants.ErrInvalidEnvironmentMode
	}
	envModeEnum := env_constants.EnvMode(envMode)
	env.App = &env_models.AppEnv{
		EnvMode:    &envModeEnum,
		ApiBaseUrl: getEnvAsString("API_BASE_URL", "http://localhost:8080"),
		AppName:    getEnvAsString("APP_NAME", "Go Boilerplate Project"),
		AppVersion: getEnvAsString("APP_VERSION", "1.0.0"),
		AppPort:    getEnvAsInt("APP_PORT", 8080),
		AppHost:    getEnvAsString("APP_HOST", "localhost"),
	}
	env.Databases = &env_models.DatabasesEnv{
		MongoDB: &env_models.MongoDB{
			Host:     getEnvAsString("MONGO_DB_HOST", "localhost"),
			Port:     getEnvAsInt("MONGO_DB_PORT", 27017),
			Database: getEnvAsString("MONGO_DB_DATABASE", "go_boilerplate_project"),
		},
		Redis: &env_models.Redis{
			Host:     getEnvAsString("REDIS_HOST", "localhost"),
			Port:     getEnvAsInt("REDIS_PORT", 6379),
			Database: getEnvAsInt("REDIS_DATABASE", 0),
			Password: getEnvAsString("REDIS_PASSWORD", ""),
			Username: getEnvAsString("REDIS_USERNAME", ""),
		},
		MySQL: &env_models.MySQL{
			Host:     getEnvAsString("MYSQL_HOST", "localhost"),
			Port:     getEnvAsInt("MYSQL_PORT", 3306),
			Database: getEnvAsString("MYSQL_DATABASE", "go_boilerplate_project"),
			Username: getEnvAsString("MYSQL_USERNAME", "root"),
			Password: getEnvAsString("MYSQL_PASSWORD", "password"),
		},
	}
	env.Logger = &env_models.LoggerEnv{
		Level:      getEnvAsString("LOGGER_LEVEL", "info"),
		FilePath:   getEnvAsString("LOGGER_FILE_PATH", "logs/app.log"),
		MaxSize:    getEnvAsInt("LOGGER_MAX_SIZE", 100),
		MaxBackups: getEnvAsInt("LOGGER_MAX_BACKUPS", 10),
		MaxAge:     getEnvAsInt("LOGGER_MAX_AGE", 30),
		Compress:   getEnvAsBool("LOGGER_COMPRESS", false),
	}
	return env, nil
}

func getEnvAsInt(key string, defaultValue int) int {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	parsedValue, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}
	return parsedValue
}

func getEnvAsString(key string, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func getEnvAsBool(key string, defaultValue bool) bool {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value == "true" || value == "1" || value == "yes" || value == "y" || defaultValue == true
}
