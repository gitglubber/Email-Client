package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port         string
	FrontendURL  string
	DatabasePath string

	GoogleClientID     string
	GoogleClientSecret string
	GoogleRedirectURL  string

	MicrosoftClientID     string
	MicrosoftClientSecret string
	MicrosoftRedirectURL  string

	OpenAIAPIKey  string
	OpenAIBaseURL string
	OpenAIModel   string

	SessionSecret string
}

func Load() (*Config, error) {
	godotenv.Load()

	return &Config{
		Port:         getEnv("PORT", "8080"),
		FrontendURL:  getEnv("FRONTEND_URL", "http://localhost:3000"),
		DatabasePath: getEnv("DATABASE_PATH", "./data/email-client.db"),

		GoogleClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		GoogleClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		GoogleRedirectURL:  getEnv("GOOGLE_REDIRECT_URL", "http://localhost:8080/auth/google/callback"),

		MicrosoftClientID:     os.Getenv("MICROSOFT_CLIENT_ID"),
		MicrosoftClientSecret: os.Getenv("MICROSOFT_CLIENT_SECRET"),
		MicrosoftRedirectURL:  getEnv("MICROSOFT_REDIRECT_URL", "http://localhost:8080/auth/microsoft/callback"),

		OpenAIAPIKey:  os.Getenv("OPENAI_API_KEY"),
		OpenAIBaseURL: getEnv("OPENAI_BASE_URL", "https://api.openai.com/v1"),
		OpenAIModel:   getEnv("OPENAI_MODEL", "gpt-4-turbo-preview"),

		SessionSecret: getEnv("SESSION_SECRET", "change-me-in-production"),
	}, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
