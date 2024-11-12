package database

import (
	"fmt"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log/slog"
	"notion_saver/src/models"
	"os"
)

func DatabaseConfig() *gorm.DB {
	err := godotenv.Load(".env")
	host := os.Getenv("POSTGRES_HOST")
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	dbname := os.Getenv("POSTGRES_DB")
	port := os.Getenv("POSTGRES_PORT")
	sslmode := os.Getenv("POSTGRES_SSL_MODE")
	timezone := os.Getenv("POSTGRES_TIME_ZONE")

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
		host, user, password, dbname, port, sslmode, timezone,
	)

	db, err := gorm.Open(
		postgres.Open(dsn), &gorm.Config{},
	)
	if err != nil {
		slog.Error(fmt.Sprintf("Error occurred when connecting to the database: %v", err))
	}

	return db
}

func MigrateSchemas() {
	err := godotenv.Load(".env")
	host := os.Getenv("POSTGRES_HOST")
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	dbname := os.Getenv("POSTGRES_DB")
	port := os.Getenv("POSTGRES_PORT")
	sslmode := os.Getenv("POSTGRES_SSL_MODE")
	timezone := os.Getenv("POSTGRES_TIME_ZONE")

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
		host, user, password, dbname, port, sslmode, timezone,
	)

	db, err := gorm.Open(
		postgres.Open(dsn), &gorm.Config{},
	)
	if err != nil {
		slog.Error(fmt.Sprintf("Error occurred when connecting to the database: %v", err))
	}

	// Migrate schemas
	db.AutoMigrate(&models.Page{})
	db.AutoMigrate(&models.Save{})
	db.AutoMigrate(&models.BlockData{})
	db.AutoMigrate(&models.Blocks{})
	db.AutoMigrate(&models.Annotation{})
}
