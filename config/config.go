package config

import (
	"crypto/rand"
	"encoding/hex"
	"log"
	"os"

	"github.com/joho/godotenv"
)

// LoadAndInitConfig loads environment variables and ensures essential configs like JWT_SECRET are set.
func LoadAndInitConfig() {
	// Attempt to load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, relying on system environment variables")
	}

	// Check if JWT_SECRET is set, if not, generate and save it
	if os.Getenv("JWT_SECRET") == "" {
		log.Println("JWT_SECRET not found. Generating a new one...")
		secretBytes := make([]byte, 32) // 256 bits
		if _, err := rand.Read(secretBytes); err != nil {
			log.Fatalf("Failed to generate random bytes for JWT secret: %v", err)
		}
		secret := hex.EncodeToString(secretBytes)

		// Append the new secret to the .env file
		file, err := os.OpenFile(".env", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatalf("Failed to open .env file to write JWT_SECRET: %v", err)
		}
		defer file.Close()

		if _, err := file.WriteString("\nJWT_SECRET=" + secret + "\n"); err != nil {
			log.Fatalf("Failed to write JWT_SECRET to .env file: %v", err)
		}

		// Reload the .env file to make the new secret available
		if err := godotenv.Load(); err != nil {
			log.Fatalf("Failed to reload .env file after writing JWT_SECRET: %v", err)
		}

		log.Println("New JWT_SECRET has been generated and saved to .env file.")
	} else {
		log.Println("JWT_SECRET loaded successfully.")
	}
}
