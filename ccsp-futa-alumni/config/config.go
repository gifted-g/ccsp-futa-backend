// ccsp-futa-alumni/config/config.go
package config

import (
    "fmt"
  //  "log"
    "os"
   // "time"

    "github.com/joho/godotenv"
)

// Global variable for the JWT Secret
var JwtSecret = []byte(getEnv("JWT_SECRET", "supersecretkey"))

type Config struct {
    Port             string
    DBURL            string
    JWTSecret        string
    SMTPHost         string
    SMTPPort         string
    SMTPUser         string
    SMTPPass         string
    FromEmail        string
    GeneralChannelID int
}

func getEnv(key, fallback string) string {
    if value, ok := os.LookupEnv(key); ok {
        return value
    }
    return fallback
}

func Load() Config {
    // Load .env file, ignore error if it doesn't exist
    _ = godotenv.Load() 

    cfg := Config{
        Port:             os.Getenv("PORT"),
        DBURL:            os.Getenv("DATABASE_URL"),
        JWTSecret:        os.Getenv("JWT_SECRET"),
        SMTPHost:         os.Getenv("SMTP_HOST"),
        SMTPPort:         os.Getenv("SMTP_PORT"),
        SMTPUser:         os.Getenv("SMTP_USER"),
        SMTPPass:         os.Getenv("SMTP_PASS"),
        FromEmail:        os.Getenv("FROM_EMAIL"),
        GeneralChannelID: 1,
    }
    
    // Set global secret based on loaded config
    JwtSecret = []byte(cfg.JWTSecret) 

    if v := os.Getenv("GENERAL_CHANNEL_ID"); v != "" {
        var n int; 
        _, _ = fmt.Sscanf(v, "%d", &n); 
        if n > 0 { cfg.GeneralChannelID = n }
    }
    if cfg.Port == "" { cfg.Port = "8080" }
    // Only Fatal if running migrations or connecting to DB, not just config load
    // log.Fatal is too aggressive here, removed for flexibility.
    
    return cfg
}