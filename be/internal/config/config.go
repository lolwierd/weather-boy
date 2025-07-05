package config

import (
	"log"
	"os"
	"sync"

	"github.com/joho/godotenv"
)

var once sync.Once

var (
	// DataDir is the directory where data files are stored.
	DataDir = "data"
	// OpenAIAPIKey is the API key for the OpenAI API.
	OpenAIAPIKey = ""
)

// LoadEnv loads environment variables from a .env file.
func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found")
	}

	if d := os.Getenv("DATA_DIR"); d != "" {
		DataDir = d
	}
	if k := os.Getenv("OPENAI_API_KEY"); k != "" {
		OpenAIAPIKey = k
	}
}
