package db

import (
	"fmt"
	"log"
	"os"

	"ccsp-futa-alumni/models"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	_ "github.com/lib/pq" // PostgreSQL driver
)

var DB *gorm.DB

func Init() error {
	dsn := os.Getenv("DATABASE_URL")
	var dialector gorm.Dialector

	if dsn == "" {
		// fallback to sqlite for local dev
		log.Println("⚠️ DATABASE_URL not provided — falling back to local sqlite dev.db")
		dialector = sqlite.Open("dev.db")
	} else {
		// Ensure SSL disabled for local PostgreSQL unless specified
		if dsn[len(dsn)-14:] != "?sslmode=disable" && dsn[len(dsn)-14:] != "?sslmode=require" {
			dsn = dsn + "?sslmode=disable"
		}
		dialector = postgres.Open(dsn)
	}

	db, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		return fmt.Errorf("❌ could not connect to database: %w", err)
	}

	// Auto migrate core models
	if err := db.AutoMigrate(
		&models.User{},
    	&models.Profile{},
    	&models.Post{},
   	 	&models.Event{},
    	&models.RSVP{},
    	&models.ChatChannel{},
    	&models.ChatMember{},
    	&models.Message{},
    	&models.SetGroup{},
    	&models.SetMember{},
    	&models.PushToken{},
	); err != nil {
		return fmt.Errorf("❌ migration error: %w", err)
	}

	DB = db
	log.Println("✅ Connected and migrated database successfully")
	return nil
}
