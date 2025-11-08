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
}

func Load() *Config {
	return &Config{
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
	}
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
