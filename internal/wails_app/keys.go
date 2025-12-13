package wails_app

import (
	"errors"
	"fmt"
	"strings"

	"github.com/mrbelka12000/interview_parser/internal/client"
	"github.com/mrbelka12000/interview_parser/internal/repo/postgres"
)

// APIKeyResult represents the result of API key operations
type APIKeyResult struct {
	Success     bool   `json:"success"`
	Message     string `json:"message,omitempty"`
	APIKey      string `json:"apiKey,omitempty"`
	LastUpdated string `json:"lastUpdated,omitempty"`
}

// GetOpenAIAPIKey retrieves the current OpenAI API key from database
func (a *App) GetOpenAIAPIKey() (*APIKeyResult, error) {
	apiKey, err := a.service.GetAPIKey()
	if err != nil {
		if errors.Is(err, postgres.ErrNoKey) {
			return &APIKeyResult{
				Success: false,
				Message: "No API key found in database",
			}, nil
		}
		return &APIKeyResult{
			Success: false,
			Message: fmt.Sprintf("Failed to retrieve API key: %s", err),
		}, nil
	}

	return &APIKeyResult{
		Success:     true,
		APIKey:      apiKey,
		LastUpdated: "Recently updated",
	}, nil
}

// SaveOpenAIAPIKey saves a new OpenAI API key to database
func (a *App) SaveOpenAIAPIKey(apiKey string) (*APIKeyResult, error) {
	if apiKey == "" {
		return &APIKeyResult{
			Message: "API key cannot be empty",
		}, nil
	}

	// Basic validation for OpenAI API key
	if !strings.HasPrefix(apiKey, "sk-") {
		return &APIKeyResult{
			Message: "Invalid API key format. OpenAI API keys should start with 'sk-'",
		}, nil
	}

	fmt.Printf("Saving API key: %s\n", apiKey)
	// Save to config for current session

	aiClient := client.New(a.cfg, apiKey)
	err := aiClient.IsValidAPIKeysProvided()
	if err != nil {
		return &APIKeyResult{
			Message: err.Error(),
		}, nil
	}

	// Save to database
	err = a.service.InsertAPIKey(apiKey)
	if err != nil {
		return &APIKeyResult{
			Success: false,
			Message: fmt.Sprintf("Failed to save API key: %s", err),
		}, nil
	}

	return &APIKeyResult{
		Success:     true,
		Message:     "API key saved successfully",
		APIKey:      apiKey,
		LastUpdated: "Just now",
	}, nil
}

// DeleteOpenAIAPIKey removes the OpenAI API key from database
func (a *App) DeleteOpenAIAPIKey() (*APIKeyResult, error) {
	err := a.service.DeleteAPIKey()
	if err != nil {
		return &APIKeyResult{
			Success: false,
			Message: fmt.Sprintf("Failed to delete API key: %s", err),
		}, nil
	}

	return &APIKeyResult{
		Success: true,
		Message: "API key deleted successfully",
	}, nil
}
