package config

import (
	"log"
	"os"
)

// JWTKey is the global secret used for signing and parsing JWT tokens.
var JWTKey []byte

// InitConfig reads essential environment variables and sets up global configuration.
func InitConfig() {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		log.Fatal("FATAL: JWT_SECRET is not set in the environment")
	}
	JWTKey = []byte(secret)
}
