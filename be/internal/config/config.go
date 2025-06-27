package config

import (
	"os"
	"sync"

	"github.com/joho/godotenv"
)

var once sync.Once

var (
	OpenAIAPIKey  string
	DataDir       string
	MetNetBaseURL string
)

// LoadEnv loads environment variables from .env once and populates typed vars.
func LoadEnv() {
	once.Do(func() {
		_ = godotenv.Load()
		OpenAIAPIKey = os.Getenv("OPENAI_API_KEY")
		DataDir = os.Getenv("DATA_DIR")
		MetNetBaseURL = os.Getenv("METNET_BASE_URL")
	})
}
