package config

import (
	"log"
	"os"
	"fmt"
	"github.com/joho/godotenv"
)

type Config struct {
	Port string
	DBURL string
	JWTSecret string
	SMTPHost string
	SMTPPort string
	SMTPUser string
	SMTPPass string
	FromEmail string
	GeneralChannelID int
}

func Load() Config {
	_ = godotenv.Load()
	cfg := Config{
		Port: os.Getenv("PORT"),
		DBURL: os.Getenv("DATABASE_URL"),
		JWTSecret: os.Getenv("JWT_SECRET"),
		SMTPHost: os.Getenv("SMTP_HOST"),
		SMTPPort: os.Getenv("SMTP_PORT"),
		SMTPUser: os.Getenv("SMTP_USER"),
		SMTPPass: os.Getenv("SMTP_PASS"),
		FromEmail: os.Getenv("FROM_EMAIL"),
		GeneralChannelID: 1,
	}
	if v := os.Getenv("GENERAL_CHANNEL_ID"); v != "" {
		var n int; _, _ = fmt.Sscanf(v, "%d", &n); if n>0 { cfg.GeneralChannelID = n }
	}
	if cfg.Port == "" { cfg.Port = "8080" }
	if cfg.DBURL == "" { log.Fatal("DATABASE_URL is required") }
	if cfg.JWTSecret == "" { log.Fatal("JWT_SECRET is required") }
	return cfg
}