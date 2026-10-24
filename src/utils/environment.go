package utils

import (
	"os"

	"github.com/joho/godotenv"
)

// Load the environment variables
func LoadEnv(path ...string) map[string]string {
	p := ".env"
	if len(path) > 0 {
		p = path[0]
	}

	err := godotenv.Load(p)

	if err != nil {
		panic(err)
	}

	return map[string]string{
		"NOTION_API_KEY":            os.Getenv("NOTION_API_KEY"),
		"NOTION_SAVER_RATE_LIMITER": os.Getenv("NOTION_SAVER_RATE_LIMITER"),
		"RABBITMQ_HOST":             os.Getenv("RABBITMQ_HOST"),
		"RABBITMQ_USER":             os.Getenv("RABBITMQ_USER"),
		"RABBITMQ_PASSWORD":         os.Getenv("RABBITMQ_PASSWORD"),
		"RABBITMQ_PORT":             os.Getenv("RABBITMQ_PORT"),
		"POSTGRES_HOST":             os.Getenv("POSTGRES_HOST"),
		"POSTGRES_USER":             os.Getenv("POSTGRES_USER"),
		"POSTGRES_PASSWORD":         os.Getenv("POSTGRES_PASSWORD"),
		"POSTGRES_PORT":             os.Getenv("POSTGRES_PORT"),
		"POSTGRES_DB":               os.Getenv("POSTGRES_DB"),
		"POSTGRES_SSLMODE":          os.Getenv("POSTGRES_SSL_MODE"),
		"POSTGRES_TIMEZONE":         os.Getenv("POSTGRES_TIMEZONE"),
		"LOG_LEVEL":                 os.Getenv("LOG_LEVEL"),
	}
}
