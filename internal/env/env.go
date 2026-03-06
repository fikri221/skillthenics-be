package env

import (
	"os"
	"strings"

	"github.com/joho/godotenv"
)

func Load() {
	_ = godotenv.Load()
}

func GetString(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}

	return fallback
}

func GetBool(key string, fallback bool) bool {
	val := os.Getenv(key)
	if val == "" {
		return fallback
	}

	return strings.ToLower(val) == "true"
}
