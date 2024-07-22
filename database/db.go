package database

import (
	"poll-app/config"
	"poll-app/data"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect() (*gorm.DB, error) {
	db, err := gorm.Open(
		postgres.Open(config.Env.DB_URL),
	)
	if err != nil {
		return nil, err
	}
	err = Migrate(db)

	if err != nil {
		return nil, err
	}

	return db, nil
}

func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(&data.User{})
}

func Seed() {
	// Seed the database
}

func Disconnect() {
	// Disconnect the database
}
