package config

import (
	"fmt"
	"os"
)

type Config struct {
	AppEnv     string
	HTTPAddr   string
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string
	JWTSecret  string
	JWTIssuer  string
	JWTTTL     int // in seconds

	MongoUri    string
	MongoDbName string

	MinioEndpoint  string
	MinioAccessKey string
	MinioSecretKey string
	MinioBucket    string
	MinioUseSSL    bool
}

func Load() *Config {
	cfg := &Config{
		AppEnv:     getEnv("APP_ENV"),
		HTTPAddr:   getEnv("HTTP_ADDR"),
		DBHost:     getEnv("DB_HOST"),
		DBPort:     getEnv("DB_PORT"),
		DBUser:     getEnv("DB_USER"),
		DBPassword: getEnv("DB_PASSWORD"),
		DBName:     getEnv("DB_NAME"),
		DBSSLMode:  getEnv("DB_SSLMODE"),
		JWTSecret:  getEnv("JWT_SECRET"),
		JWTIssuer:  getEnv("JWT_ISSUER"),
		JWTTTL:     atoiMust(getEnv("JWT_TTL")),

		MongoUri:    getEnv("MONGO_URI"),
		MongoDbName: getEnv("MONGO_DB_NAME"),

		MinioEndpoint:  getEnv("MINIO_ENDPOINT"),
		MinioAccessKey: getEnv("MINIO_ACCESS_KEY"),
		MinioSecretKey: getEnv("MINIO_SECRET_KEY"),
		MinioBucket:    getEnv("MINIO_BUCKET"),
		MinioUseSSL:    getEnv("MINIO_USE_SSL") == "true",
	}

	cfg.validate()
	return cfg
}

func (c *Config) DatabaseURL() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		c.DBUser,
		c.DBPassword,
		c.DBHost,
		c.DBPort,
		c.DBName,
		c.DBSSLMode,
	)
}

func (c *Config) validate() {
	if c.JWTTTL <= 0 {
		panic("JWT_TTL must be greater than zero")
	}
}

func getEnv(key string) string {
	value, ok := os.LookupEnv(key)
	if !ok || value == "" {
		panic(fmt.Sprintf("environment variable %s not set", key))
	}
	return value
}

func atoiMust(s string) int {
	var i int
	_, err := fmt.Sscanf(s, "%d", &i)
	if err != nil {
		panic(fmt.Sprintf("invalid integer value: %s", s))
	}
	return i
}
