package database

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log/slog"
)

type Database struct {
	Client *gorm.DB
}

// Open a connection to the Postgres database.
func (d *Database) Connect(
	host string,
	user string,
	password string,
	dbname string,
	port string,
	sslmode string,
	timezone string,
) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
		host,
		user,
		password,
		dbname,
		port,
		sslmode,
		timezone,
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})

	if err != nil {
		slog.Error(
			fmt.Sprintf(
				"Could not connect to Postgres: %v",
				err,
			),
		)
		panic(err)
	}

	d.Client = db
}
