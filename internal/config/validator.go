package config

import (
	"fmt"
	"strings"
)

// ValidationResult holds the validation results
type ValidationResult struct {
	Valid    bool
	Errors   []string
	Warnings []string
}

// Validate checks if the configuration is valid and provides helpful error messages
func (c *Config) Validate() *ValidationResult {
	result := &ValidationResult{
		Valid:    true,
		Errors:   []string{},
		Warnings: []string{},
	}

	// Check Google OAuth
	if c.GoogleClientID == "" || c.GoogleClientID == "your-google-client-id" {
		result.Valid = false
		result.Errors = append(result.Errors, "Google OAuth not configured")
		result.Warnings = append(result.Warnings,
			"To set up Google OAuth:\n"+
			"  1. Go to https://console.cloud.google.com/\n"+
			"  2. Create a project and enable Gmail API\n"+
			"  3. Create OAuth 2.0 credentials\n"+
			"  4. Set GOOGLE_CLIENT_ID and GOOGLE_CLIENT_SECRET in .env\n"+
			"  Or run: ./scripts/setup-oauth.sh")
	}

	if c.GoogleClientSecret == "" || c.GoogleClientSecret == "your-google-client-secret" {
		if !strings.Contains(strings.Join(result.Errors, ""), "Google OAuth") {
			result.Valid = false
			result.Errors = append(result.Errors, "Google Client Secret not configured")
		}
	}

	// Check Microsoft OAuth
	if c.MicrosoftClientID == "" || c.MicrosoftClientID == "your-microsoft-client-id" {
		result.Warnings = append(result.Warnings,
			"Microsoft OAuth not configured (Microsoft sign-in won't work)\n"+
			"  To enable: Run ./scripts/setup-oauth.sh")
	}

	if c.MicrosoftClientSecret == "" || c.MicrosoftClientSecret == "your-microsoft-client-secret" {
		if !strings.Contains(strings.Join(result.Warnings, ""), "Microsoft OAuth") {
			result.Warnings = append(result.Warnings, "Microsoft Client Secret not configured")
		}
	}

	// Check OpenAI
	if c.OpenAIAPIKey == "" || c.OpenAIAPIKey == "your-openai-api-key" {
		result.Warnings = append(result.Warnings,
			"OpenAI API not configured (AI features won't work)\n"+
			"  Options:\n"+
			"    - Get OpenAI key: https://platform.openai.com/api-keys\n"+
			"    - Use local LLM: Install LM Studio or Ollama\n"+
			"  Configure: Run ./scripts/setup-oauth.sh")
	}

	// Check session secret
	if c.SessionSecret == "change-me-in-production" {
		result.Warnings = append(result.Warnings,
			"Using default session secret (insecure for production)\n"+
			"  Generate a secure secret: openssl rand -base64 32")
	}

	return result
}

// PrintValidation prints the validation results in a user-friendly format
func (r *ValidationResult) PrintValidation() {
	if r.Valid && len(r.Warnings) == 0 {
		fmt.Println("✓ Configuration validated successfully!")
		return
	}

	if !r.Valid {
		fmt.Println("\n❌ Configuration Errors:")
		fmt.Println(strings.Repeat("=", 50))
		for _, err := range r.Errors {
			fmt.Printf("\n%s\n", err)
		}
		for _, warning := range r.Warnings {
			if strings.Contains(warning, "To set up") || strings.Contains(warning, "To enable") {
				fmt.Printf("\n%s\n", warning)
			}
		}
		fmt.Println(strings.Repeat("=", 50))
		fmt.Println("\nPlease configure your credentials before starting the server.")
		fmt.Println("Run: ./scripts/setup-oauth.sh")
		fmt.Println()
	} else if len(r.Warnings) > 0 {
		fmt.Println("\n⚠️  Configuration Warnings:")
		fmt.Println(strings.Repeat("=", 50))
		for _, warning := range r.Warnings {
			fmt.Printf("\n%s\n", warning)
		}
		fmt.Println(strings.Repeat("=", 50))
		fmt.Println("\nThe server will start, but some features may not work.")
		fmt.Println()
	}
}

// MustValidate validates the config and exits if there are errors
func (c *Config) MustValidate() {
	result := c.Validate()
	result.PrintValidation()

	if !result.Valid {
		fmt.Println("Exiting due to configuration errors.")
		panic("Invalid configuration")
	}
}
