package db

import (
	"fmt"
	"log"
	"os"
	"database/sql"
	_"github.com/lib/pq" // ✅ PostgreSQL driver
	"ccsp-futa-alumni/models"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Init() error {
	dsn := os.Getenv("DATABASE_URL")
	var dialector gorm.Dialector
	if dsn == "" {
		// fallback to sqlite for local dev
		log.Println("DATABASE_URL not provided — falling back to sqlite dev.db")
		dialector = sqlite.Open("dev.db")
	} else {
		dialector = postgres.Open(dsn)
	}

	db, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		return err
	}

	// Auto migrate core models
	if err := db.AutoMigrate(
		&models.User{},
		&models.Profile{},
		&models.ChatChannel{},
		&models.ChatMember{},
		&models.Message{},
		&models.PushToken{},
	); err != nil {
		return fmt.Errorf("migrate error: %w", err)
	}

	DB = db
	return nil
}
