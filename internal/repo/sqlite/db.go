package sqlite

import (
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

var (
	ErrNoKey = errors.New("no key found")
	db       *sql.DB
)

func InitDB(dbPath string) (err error) {
	db, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		return
	}

	if err = db.Ping(); err != nil {
		return
	}

	ddl := `
CREATE TABLE IF NOT EXISTS api_keys (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    api_key TEXT NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);`

	_, err = db.Exec(ddl)
	if err != nil {
		return fmt.Errorf("create api keys table: %w", err)
	}

	ddl = `
	CREATE TABLE IF NOT EXISTS interviews (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
	);

	CREATE INDEX IF NOT EXISTS interviews_created_at ON interviews(created_at);
	`

	_, err = db.Exec(ddl)
	if err != nil {
		return fmt.Errorf("create interviews table: %w", err)
	}

	ddl = `
	CREATE TABLE IF NOT EXISTS question_answers (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		interview_id INTEGER NOT NULL,
		question TEXT NOT NULL,
		full_answer TEXT,
		accuracy REAL NOT NULL,
		reason_unanswered TEXT,
		created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (interview_id) REFERENCES interviews(id) ON DELETE CASCADE
	);

	CREATE INDEX IF NOT EXISTS idx_question_answers_interview_id ON question_answers(interview_id);
	CREATE INDEX IF NOT EXISTS idx_question_answers_accuracy ON question_answers(accuracy);
	`

	_, err = db.Exec(ddl)
	if err != nil {
		return fmt.Errorf("create question answers table: %w", err)
	}

	ddl = `
	CREATE TABLE IF NOT EXISTS calls (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		transcript TEXT NOT NULL,
		analysis TEXT,
		created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
	);

	CREATE INDEX IF NOT EXISTS idx_calls_created_at ON calls(created_at);
	`

	_, err = db.Exec(ddl)
	if err != nil {
		return fmt.Errorf("create calls table: %w", err)
	}

	return nil
}
