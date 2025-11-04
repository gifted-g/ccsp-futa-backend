package config

import "os"

var JwtSecret = []byte(getEnv("JWT_SECRET", "supersecretkey"))

func getEnv(key, fallback string) string {
    if value, ok := os.LookupEnv(key); ok {
        return value
    }
    return fallback
}
