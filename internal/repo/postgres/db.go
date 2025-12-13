package postgres

import (
	"errors"
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/mrbelka12000/interview_parser/internal/models"
)

var (
	ErrNoKey = errors.New("no key found")
	db       *gorm.DB
)

func InitDB(pgURL string) error {
	var err error

	db, err = gorm.Open(postgres.Open(pgURL), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

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
