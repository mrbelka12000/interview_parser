package models

import "time"

// InterviewAnalytics represents analytics data for a single interview
type InterviewAnalytics struct {
	ID                        int       `json:"id"`
	TotalQuestions            int       `json:"totalQuestions"`
	AnsweredQuestions         int       `json:"answeredQuestions"`
	UnansweredQuestions       int       `json:"unansweredQuestions"`
	AnsweredPercentage        float64   `json:"answeredPercentage"`
	UnansweredPercentage      float64   `json:"unansweredPercentage"`
	AverageAccuracy           float64   `json:"averageAccuracy"`
	AverageAnsweredAccuracy   float64   `json:"averageAnsweredAccuracy"`
	HighConfidenceQuestions   int       `json:"highConfidenceQuestions"`   // accuracy > 0.8
	MediumConfidenceQuestions int       `json:"mediumConfidenceQuestions"` // accuracy 0.5-0.8
	LowConfidenceQuestions    int       `json:"lowConfidenceQuestions"`    // accuracy < 0.5
	QuestionsWithReason       int       `json:"questionsWithReason"`
	CreatedAt                 time.Time `json:"createdAt"`
	UpdatedAt                 time.Time `json:"updatedAt"`
}

// GlobalAnalytics represents aggregated statistics across all interviews
type GlobalAnalytics struct {
	TotalInterviews        int       `json:"totalInterviews"`
	TotalQuestions         int       `json:"totalQuestions"`
	TotalAnswered          int       `json:"totalAnswered"`
	TotalUnanswered        int       `json:"totalUnanswered"`
	GlobalAnsweredPercent  float64   `json:"globalAnsweredPercent"`
	GlobalAverageAccuracy  float64   `json:"globalAverageAccuracy"`
	GlobalAnsweredAccuracy float64   `json:"globalAnsweredAccuracy"`
	BestInterviewID        int       `json:"bestInterviewID"`
	BestInterviewScore     float64   `json:"bestInterviewScore"`
	WorstInterviewID       int       `json:"worstInterviewID"`
	WorstInterviewScore    float64   `json:"worstInterviewScore"`
	LastUpdated            time.Time `json:"lastUpdated"`
}
