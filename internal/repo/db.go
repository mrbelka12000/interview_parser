package repo

import (
	"errors"
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/mrbelka12000/interview_parser/internal/config"
	"github.com/mrbelka12000/interview_parser/internal/models"
)

var (
	ErrNoKey = errors.New("no key found")
	db       *gorm.DB
)

func InitDB(cfg *config.Config) error {
	var err error
	
	if cfg.DatabaseURL != "" {
		// Use PostgreSQL
		db, err = gorm.Open(postgres.Open(cfg.DatabaseURL), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		})
	} else {
		// Fallback to SQLite
		db, err = gorm.Open(sqlite.Open(cfg.DBPath), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		})
	}
	
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	// Auto migrate the schema
	err = db.AutoMigrate(
		&models.APIKey{},
		&models.AnalyzeInterview{},
		&models.QuestionAnswer{},
		&models.Call{},
	)
	if err != nil {
		return fmt.Errorf("failed to migrate database: %w", err)
	}

	return nil
}

func GetDB() *gorm.DB {
	return db
}
