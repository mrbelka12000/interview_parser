package postgres

import (
	"errors"
	"fmt"

	"gorm.io/gorm"

	"github.com/mrbelka12000/interview_parser/internal/models"
)

type ApiKeyRepo struct {
}

func NewApiKeyRepo() *ApiKeyRepo {
	return &ApiKeyRepo{}
}

func (ap *ApiKeyRepo) GetOpenAIAPIKeyFromDB() (string, error) {
	var apiKey models.APIKey
	if err := GetDB().Order("created_at DESC").First(&apiKey).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", ErrNoKey
		}
		return "", fmt.Errorf("failed to get API key: %w", err)
	}

	return apiKey.APIKey, nil
}

func (ap *ApiKeyRepo) InsertOpenAIAPIKey(openAIAPIKey string) error {
	apiKey := &models.APIKey{
		APIKey: openAIAPIKey,
	}

	if err := GetDB().Create(apiKey).Error; err != nil {
		return fmt.Errorf("failed to insert API key: %w", err)
	}

	return nil
}

func (ap *ApiKeyRepo) DeleteOpenAIAPIKey() error {
	if err := GetDB().Where("1 = 1").Delete(&models.APIKey{}).Error; err != nil {
		return fmt.Errorf("failed to delete API keys: %w", err)
	}

	return nil
}
