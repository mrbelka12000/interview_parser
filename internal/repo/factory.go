package repo

import (
	"github.com/mrbelka12000/interview_parser/internal/config"
	"github.com/mrbelka12000/interview_parser/internal/repo/postgres"
	"github.com/mrbelka12000/interview_parser/internal/repo/sqlite"
)

// NewRepositories creates repository instances based on database configuration
func NewRepositories(cfg *config.Config) (ApiKeyRepository, InterviewRepository, CallRepository) {
	switch {
	case cfg.DBConfig.PGURL != "":
		if err := postgres.InitDB(cfg.DBConfig.PGURL); err != nil {
			panic(err)
		}
		return newPostgresRepositories()
	default:
		if err := sqlite.InitDB(cfg.DBConfig.Path); err != nil {
			panic(err)
		}
		return newSQLiteRepositories()
	}
}

// newPostgresRepositories creates PostgreSQL repository instances
func newPostgresRepositories() (ApiKeyRepository, InterviewRepository, CallRepository) {
	apiKeyRepo := postgres.NewApiKeyRepo()
	interviewRepo := postgres.NewInterviewRepo()
	callRepo := postgres.NewCallRepo()

	return apiKeyRepo, interviewRepo, callRepo
}

// newSQLiteRepositories creates SQLite repository instances
func newSQLiteRepositories() (ApiKeyRepository, InterviewRepository, CallRepository) {
	apiKeyRepo := sqlite.NewApiKeyRepo()
	interviewRepo := sqlite.NewInterviewRepo()
	callRepo := sqlite.NewCallRepo()

	return apiKeyRepo, interviewRepo, callRepo
}
