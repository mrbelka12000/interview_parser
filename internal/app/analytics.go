package app

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/mrbelka12000/interview_parser/internal/client"
	"github.com/mrbelka12000/interview_parser/internal/config"
)

// InterviewAnalytics represents analytics data for a single interview
type InterviewAnalytics struct {
	ID                     int       `json:"id"`
	InterviewPath          string    `json:"interviewPath"`
	AnalysisPath           string    `json:"analysisPath"`
	TotalQuestions         int       `json:"totalQuestions"`
	AnsweredQuestions      int       `json:"answeredQuestions"`
	UnansweredQuestions    int       `json:"unansweredQuestions"`
	AnsweredPercentage     float64   `json:"answeredPercentage"`
	UnansweredPercentage   float64   `json:"unansweredPercentage"`
	AverageAccuracy        float64   `json:"averageAccuracy"`
	AverageAnsweredAccuracy float64  `json:"averageAnsweredAccuracy"`
	HighConfidenceQuestions int      `json:"highConfidenceQuestions"` // accuracy > 0.8
	MediumConfidenceQuestions int    `json:"mediumConfidenceQuestions"` // accuracy 0.5-0.8
	LowConfidenceQuestions  int     `json:"lowConfidenceQuestions"` // accuracy < 0.5
	QuestionsWithReason    int       `json:"questionsWithReason"`
	CreatedAt              time.Time `json:"createdAt"`
	UpdatedAt              time.Time `json:"updatedAt"`
}

// GlobalAnalytics represents aggregated statistics across all interviews
type GlobalAnalytics struct {
	TotalInterviews        int     `json:"totalInterviews"`
	TotalQuestions         int     `json:"totalQuestions"`
	TotalAnswered          int     `json:"totalAnswered"`
	TotalUnanswered        int     `json:"totalUnanswered"`
	GlobalAnsweredPercent  float64 `json:"globalAnsweredPercent"`
	GlobalAverageAccuracy  float64 `json:"globalAverageAccuracy"`
	GlobalAnsweredAccuracy float64 `json:"globalAnsweredAccuracy"`
	BestInterviewID        int     `json:"bestInterviewID"`
	BestInterviewPath      string  `json:"bestInterviewPath"`
	BestInterviewScore     float64 `json:"bestInterviewScore"`
	WorstInterviewID       int     `json:"worstInterviewID"`
	WorstInterviewPath     string  `json:"worstInterviewPath"`
	WorstInterviewScore    float64 `json:"worstInterviewScore"`
	LastUpdated            time.Time `json:"lastUpdated"`
}

// AnalyticsFilters represents filters for querying analytics
type AnalyticsFilters struct {
	DateFrom     *time.Time `json:"dateFrom,omitempty"`
	DateTo       *time.Time `json:"dateTo,omitempty"`
	MinAccuracy  *float64   `json:"minAccuracy,omitempty"`
	MaxAccuracy  *float64   `json:"maxAccuracy,omitempty"`
}

// CalculateAnalytics computes analytics from an AnalyzeResponse
func (a *App) CalculateAnalytics(analyzeResp client.AnalyzeResponse, interviewPath, analysisPath string) (*InterviewAnalytics, error) {
	if len(analyzeResp.Questions) == 0 {
		return nil, fmt.Errorf("no questions in analysis response")
	}

	totalQuestions := len(analyzeResp.Questions)
	answered := 0
	unanswered := 0
	withReason := 0
	totalAccuracy := 0.0
	answeredAccuracy := 0.0
	highConf := 0
	mediumConf := 0
	lowConf := 0

	for _, q := range analyzeResp.Questions {
		totalAccuracy += q.Accuracy
		
		if q.FullAnswer != "" && q.Accuracy > 0 {
			answered++
			answeredAccuracy += q.Accuracy
			
			if q.Accuracy > 0.8 {
				highConf++
			} else if q.Accuracy >= 0.5 {
				mediumConf++
			} else {
				lowConf++
			}
		} else {
			unanswered++
		}
		
		if q.ReasonUnanswered != "" {
			withReason++
		}
	}

	answeredPercentage := float64(answered) / float64(totalQuestions) * 100
	unansweredPercentage := float64(unanswered) / float64(totalQuestions) * 100
	avgAccuracy := totalAccuracy / float64(totalQuestions)
	var avgAnsweredAccuracy float64
	if answered > 0 {
		avgAnsweredAccuracy = answeredAccuracy / float64(answered)
	}

	return &InterviewAnalytics{
		InterviewPath:            interviewPath,
		AnalysisPath:             analysisPath,
		TotalQuestions:           totalQuestions,
		AnsweredQuestions:        answered,
		UnansweredQuestions:      unanswered,
		AnsweredPercentage:       answeredPercentage,
		UnansweredPercentage:     unansweredPercentage,
		AverageAccuracy:          avgAccuracy,
		AverageAnsweredAccuracy:  avgAnsweredAccuracy,
		HighConfidenceQuestions:  highConf,
		MediumConfidenceQuestions: mediumConf,
		LowConfidenceQuestions:   lowConf,
		QuestionsWithReason:      withReason,
		CreatedAt:                time.Now(),
		UpdatedAt:                time.Now(),
	}, nil
}

// SaveAnalytics saves interview analytics to the database
func (a *App) SaveAnalytics(analytics *InterviewAnalytics) error {
	db, err := connectToAnalyticsDB(a.cfg)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	defer db.Close()

	query := `
	INSERT INTO interview_analytics (
		interview_path, analysis_path, total_questions, answered_questions, unanswered_questions,
		answered_percentage, unanswered_percentage, average_accuracy, average_answered_accuracy,
		high_confidence_questions, medium_confidence_questions, low_confidence_questions,
		questions_with_reason, created_at, updated_at
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err = db.Exec(query,
		analytics.InterviewPath, analytics.AnalysisPath, analytics.TotalQuestions, analytics.AnsweredQuestions, analytics.UnansweredQuestions,
		analytics.AnsweredPercentage, analytics.UnansweredPercentage, analytics.AverageAccuracy, analytics.AverageAnsweredAccuracy,
		analytics.HighConfidenceQuestions, analytics.MediumConfidenceQuestions, analytics.LowConfidenceQuestions,
		analytics.QuestionsWithReason, analytics.CreatedAt, analytics.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to save analytics: %w", err)
	}

	return nil
}

// GetInterviewAnalytics retrieves analytics for a specific interview
func (a *App) GetInterviewAnalytics(interviewPath string) (*InterviewAnalytics, error) {
	db, err := connectToAnalyticsDB(a.cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	defer db.Close()

	query := `
	SELECT id, interview_path, analysis_path, total_questions, answered_questions, unanswered_questions,
		   answered_percentage, unanswered_percentage, average_accuracy, average_answered_accuracy,
		   high_confidence_questions, medium_confidence_questions, low_confidence_questions,
		   questions_with_reason, created_at, updated_at
	FROM interview_analytics 
	WHERE interview_path = ?
	`

	var analytics InterviewAnalytics
	err = db.QueryRow(query, interviewPath).Scan(
		&analytics.ID, &analytics.InterviewPath, &analytics.AnalysisPath, &analytics.TotalQuestions, &analytics.AnsweredQuestions, &analytics.UnansweredQuestions,
		&analytics.AnsweredPercentage, &analytics.UnansweredPercentage, &analytics.AverageAccuracy, &analytics.AverageAnsweredAccuracy,
		&analytics.HighConfidenceQuestions, &analytics.MediumConfidenceQuestions, &analytics.LowConfidenceQuestions,
		&analytics.QuestionsWithReason, &analytics.CreatedAt, &analytics.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no analytics found for interview: %s", interviewPath)
		}
		return nil, fmt.Errorf("failed to retrieve analytics: %w", err)
	}

	return &analytics, nil
}

// GetAllInterviewAnalytics retrieves all interview analytics with optional filters
func (a *App) GetAllInterviewAnalytics(filters *AnalyticsFilters) ([]InterviewAnalytics, error) {
	db, err := connectToAnalyticsDB(a.cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	defer db.Close()

	query := `
	SELECT id, interview_path, analysis_path, total_questions, answered_questions, unanswered_questions,
		   answered_percentage, unanswered_percentage, average_accuracy, average_answered_accuracy,
		   high_confidence_questions, medium_confidence_questions, low_confidence_questions,
		   questions_with_reason, created_at, updated_at
	FROM interview_analytics
	WHERE 1=1
	`

	args := []interface{}{}
	argIndex := 1

	if filters != nil {
		if filters.DateFrom != nil {
			query += fmt.Sprintf(" AND created_at >= $%d", argIndex)
			args = append(args, filters.DateFrom)
			argIndex++
		}
		if filters.DateTo != nil {
			query += fmt.Sprintf(" AND created_at <= $%d", argIndex)
			args = append(args, filters.DateTo)
			argIndex++
		}
		if filters.MinAccuracy != nil {
			query += fmt.Sprintf(" AND average_accuracy >= $%d", argIndex)
			args = append(args, *filters.MinAccuracy)
			argIndex++
		}
		if filters.MaxAccuracy != nil {
			query += fmt.Sprintf(" AND average_accuracy <= $%d", argIndex)
			args = append(args, *filters.MaxAccuracy)
			argIndex++
		}
	}

	query += " ORDER BY created_at DESC"

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query analytics: %w", err)
	}
	defer rows.Close()

	var analyticsList []InterviewAnalytics
	for rows.Next() {
		var analytics InterviewAnalytics
		err := rows.Scan(
			&analytics.ID, &analytics.InterviewPath, &analytics.AnalysisPath, &analytics.TotalQuestions, &analytics.AnsweredQuestions, &analytics.UnansweredQuestions,
			&analytics.AnsweredPercentage, &analytics.UnansweredPercentage, &analytics.AverageAccuracy, &analytics.AverageAnsweredAccuracy,
			&analytics.HighConfidenceQuestions, &analytics.MediumConfidenceQuestions, &analytics.LowConfidenceQuestions,
			&analytics.QuestionsWithReason, &analytics.CreatedAt, &analytics.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan analytics row: %w", err)
		}
		analyticsList = append(analyticsList, analytics)
	}

	return analyticsList, nil
}

// GetGlobalAnalytics calculates aggregated statistics across all interviews
func (a *App) GetGlobalAnalytics(filters *AnalyticsFilters) (*GlobalAnalytics, error) {
	analyticsList, err := a.GetAllInterviewAnalytics(filters)
	if err != nil {
		return nil, fmt.Errorf("failed to get interview analytics: %w", err)
	}

	if len(analyticsList) == 0 {
		return &GlobalAnalytics{
			TotalInterviews: 0,
			LastUpdated:     time.Now(),
		}, nil
	}

	global := &GlobalAnalytics{
		TotalInterviews: len(analyticsList),
		LastUpdated:     time.Now(),
	}

	var totalAccuracy, totalAnsweredAccuracy float64
	var bestScore, worstScore float64
	bestIdx, worstIdx := 0, 0

	for i, analytics := range analyticsList {
		global.TotalQuestions += analytics.TotalQuestions
		global.TotalAnswered += analytics.AnsweredQuestions
		global.TotalUnanswered += analytics.UnansweredQuestions
		totalAccuracy += analytics.AverageAccuracy
		totalAnsweredAccuracy += analytics.AverageAnsweredAccuracy

		// Calculate a composite score (answered percentage + average accuracy)
		compositeScore := analytics.AnsweredPercentage + analytics.AverageAccuracy*100
		
		if i == 0 || compositeScore > bestScore {
			bestScore = compositeScore
			bestIdx = i
		}
		if i == 0 || compositeScore < worstScore {
			worstScore = compositeScore
			worstIdx = i
		}
	}

	global.GlobalAverageAccuracy = totalAccuracy / float64(len(analyticsList))
	global.GlobalAnsweredAccuracy = totalAnsweredAccuracy / float64(len(analyticsList))
	global.GlobalAnsweredPercent = float64(global.TotalAnswered) / float64(global.TotalQuestions) * 100

	if len(analyticsList) > 0 {
		best := analyticsList[bestIdx]
		worst := analyticsList[worstIdx]
		
		global.BestInterviewID = best.ID
		global.BestInterviewPath = best.InterviewPath
		global.BestInterviewScore = bestScore
		
		global.WorstInterviewID = worst.ID
		global.WorstInterviewPath = worst.InterviewPath
		global.WorstInterviewScore = worstScore
	}

	return global, nil
}

// connectToAnalyticsDB connects to the database and ensures analytics tables exist
func connectToAnalyticsDB(cfg *config.Config) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", cfg.DBPath)
	if err != nil {
		return nil, fmt.Errorf("open db: %w", err)
	}

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("ping db: %w", err)
	}

	// Create analytics tables if they don't exist
	ddl := `
	CREATE TABLE IF NOT EXISTS interview_analytics (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		interview_path TEXT NOT NULL,
		analysis_path TEXT NOT NULL,
		total_questions INTEGER NOT NULL,
		answered_questions INTEGER NOT NULL,
		unanswered_questions INTEGER NOT NULL,
		answered_percentage REAL NOT NULL,
		unanswered_percentage REAL NOT NULL,
		average_accuracy REAL NOT NULL,
		average_answered_accuracy REAL,
		high_confidence_questions INTEGER NOT NULL DEFAULT 0,
		medium_confidence_questions INTEGER NOT NULL DEFAULT 0,
		low_confidence_questions INTEGER NOT NULL DEFAULT 0,
		questions_with_reason INTEGER NOT NULL DEFAULT 0,
		created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		UNIQUE(interview_path)
	);

	CREATE INDEX IF NOT EXISTS idx_interview_analytics_created_at ON interview_analytics(created_at);
	CREATE INDEX IF NOT EXISTS idx_interview_analytics_accuracy ON interview_analytics(average_accuracy);
	`

	_, err = db.Exec(ddl)
	if err != nil {
		return nil, fmt.Errorf("create analytics tables: %w", err)
	}

	return db, nil
}
