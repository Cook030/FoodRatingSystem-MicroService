package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	ServerPort       string
	DBHost           string
	DBPort           string
	DBUser           string
	DBPassword       string
	DBName           string
	DBSSLMode        string
	RedisHost        string
	RedisPort        string
	RedisPassword    string
	RedisDB          int
	JWTSecret        string
	JWTExpirationHrs int
}

var AppConfig *Config

func LoadConfig() {
	_ = godotenv.Load()

	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "5432")
	dbUser := getEnv("DB_USER", "postgres")
	dbPassword := getEnv("DB_PASSWORD", "root")
	dbName := getEnv("DB_NAME", "postgres")
	dbSSLMode := getEnv("DB_SSLMODE", "disable")

	redisHost := getEnv("REDIS_HOST", "localhost")
	redisPort := getEnv("REDIS_PORT", "6379")
	redisPassword := getEnv("REDIS_PASSWORD", "")
	redisDB := 0

	jwtSecret := getEnv("JWT_SECRET", "your-secret-key-change-this-in-production")
	jwtExp := 24

	AppConfig = &Config{
		ServerPort:       getEnv("API_GATEWAY_PORT", "8080"),
		DBHost:           dbHost,
		DBPort:           dbPort,
		DBUser:           dbUser,
		DBPassword:       dbPassword,
		DBName:           dbName,
		DBSSLMode:        dbSSLMode,
		RedisHost:        redisHost,
		RedisPort:        redisPort,
		RedisPassword:    redisPassword,
		RedisDB:          redisDB,
		JWTSecret:        jwtSecret,
		JWTExpirationHrs: jwtExp,
	}

	log.Println("配置加载完成")
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
