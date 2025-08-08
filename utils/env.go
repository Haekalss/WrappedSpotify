package utils

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func LoadEnv() {
	// Try to load .env file, but don't fail if it doesn't exist
	// This allows the app to work in production where env vars are set directly
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found, using environment variables")
	}

	// Validate required environment variables
	requiredVars := []string{
		"SPOTIFY_CLIENT_ID",
		"SPOTIFY_CLIENT_SECRET",
		"SPOTIFY_REDIRECT_URI",
		"FRONTEND_URL",
	}

	for _, envVar := range requiredVars {
		if os.Getenv(envVar) == "" {
			log.Fatalf("Required environment variable %s is not set", envVar)
		}
	}
}
