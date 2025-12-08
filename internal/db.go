package internal

import (
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/mattn/go-sqlite3"

	"github.com/mrbelka12000/interview_parser/internal/config"
)

var (
	ErrNoKey = errors.New("no key found")
)

func GetOpenAIAPIKeyFromDB(cfg *config.Config) (string, error) {
	db, err := connectToDB(cfg)
	if err != nil {
		return "", err
	}
	defer db.Close()

	var apiKey string
	err = db.QueryRow(`SELECT api_key
FROM api_keys
ORDER BY created_at DESC
LIMIT 1;`).Scan(&apiKey)
	if errors.Is(err, sql.ErrNoRows) {
		return "", ErrNoKey
	}

	return apiKey, nil
}

func InsertOpenAIAPIKey(cfg *config.Config) error {
	db, err := connectToDB(cfg)
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec(`
INSERT INTO api_keys (api_key)
VALUES (?);`, cfg.OpenAIAPIKey)
	if err != nil {
		return err
	}

	return nil
}

func connectToDB(cfg *config.Config) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", cfg.DBPath)
	if err != nil {
		return nil, fmt.Errorf("open db: %w", err)
	}

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("ping db: %w", err)
	}

	ddl := `
CREATE TABLE IF NOT EXISTS api_keys (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    api_key TEXT NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);`

	_, err = db.Exec(ddl)
	if err != nil {
		return nil, fmt.Errorf("create table: %w", err)
	}

	return db, nil
}
