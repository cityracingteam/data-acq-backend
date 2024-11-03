package database

import (
	"fmt"
	"log"

	"github.com/cityracingteam/data-acq-backend/environment"
	"github.com/cityracingteam/data-acq-backend/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect() *gorm.DB {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s",
		environment.GetEnvOrDefault("DB_HOST"),
		environment.GetEnvOrDefault("DB_USER"),
		environment.GetEnvOrDefault("DB_PASS"),
		environment.GetEnvOrDefault("DB_NAME"),
		environment.GetEnvOrDefault("DB_PORT"),
	)
	// Attempt to connect to the db
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect to the database")
	}

	migrate(db)

	return db
}

func migrate(db *gorm.DB) {
	log.Println("[database/migrate]: migrating...")

	// "Migrate" the schema
	// This will create tables, keys, columns, etc. Everything really.
	// See https://gorm.io/docs/migration.html
	// Note that we need to pass each struct in our schema.

	var models = []interface{}{
		// Add **all** models here
		models.JwtKey{},
		models.User{},
	}

	// Iterate over each model and migrate it, exiting
	// the program with an error if one is encountered.
	for index, model := range models {
		err := db.AutoMigrate(model)
		if err != nil {
			log.Panicf("Error with migrate call %d: '%s'", index, err)
		}
	}

	log.Println("[database/migrate]: done.")
}
