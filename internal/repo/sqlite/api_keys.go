package sqlite

import (
	"database/sql"
	"errors"
)

type ApiKeyRepo struct {
}

func NewApiKeyRepo() *ApiKeyRepo {
	return &ApiKeyRepo{}
}

func (ap *ApiKeyRepo) GetOpenAIAPIKeyFromDB() (string, error) {
	var apiKey string
	err := db.QueryRow(`SELECT api_key
FROM api_keys
ORDER BY created_at DESC
LIMIT 1;`).Scan(&apiKey)
	if errors.Is(err, sql.ErrNoRows) {
		return "", ErrNoKey
	}

	return apiKey, nil
}

func (ap *ApiKeyRepo) InsertOpenAIAPIKey(openAIAPIKey string) error {
	_, err := db.Exec(`
INSERT INTO api_keys (api_key)
VALUES (?);`, openAIAPIKey)
	if err != nil {
		return err
	}

	return nil
}

func (ap *ApiKeyRepo) DeleteOpenAIAPIKey() error {
	_, err := db.Exec(`DELETE FROM api_keys;`)
	if err != nil {
		return err
	}

	return nil
}
