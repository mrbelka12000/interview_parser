package service

import (
	"fmt"
	"time"

	"github.com/mrbelka12000/interview_parser/internal/models"
)

func (s *Service) calculateAnalytics(interview models.AnalyzeInterviewWithQA) models.InterviewAnalytics {

	totalQuestions := len(interview.QA)
	answered := 0
	unanswered := 0
	withReason := 0
	totalAccuracy := 0.0
	answeredAccuracy := 0.0
	highConf := 0
	mediumConf := 0
	lowConf := 0

	for _, q := range interview.QA {
		totalAccuracy += q.Accuracy

		if q.FullAnswer != "" && q.Accuracy > 0.7 {
			answered++
			answeredAccuracy += q.Accuracy

			if q.Accuracy > 0.9 {
				highConf++
			} else if q.Accuracy >= 0.85 {
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

	return models.InterviewAnalytics{
		TotalQuestions:            totalQuestions,
		AnsweredQuestions:         answered,
		UnansweredQuestions:       unanswered,
		AnsweredPercentage:        answeredPercentage,
		UnansweredPercentage:      unansweredPercentage,
		AverageAccuracy:           avgAccuracy,
		AverageAnsweredAccuracy:   avgAnsweredAccuracy,
		HighConfidenceQuestions:   highConf,
		MediumConfidenceQuestions: mediumConf,
		LowConfidenceQuestions:    lowConf,
		QuestionsWithReason:       withReason,
		CreatedAt:                 time.Now(),
		UpdatedAt:                 time.Now(),
	}
}

func (s *Service) GetGlobalAnalytics(filters *models.GetInterviewsFilters) (*models.GlobalAnalytics, error) {
	analyticsList, err := s.GetAllInterviewAnalytics(filters)
	if err != nil {
		return nil, err
	}

	if len(analyticsList) == 0 {
		return &models.GlobalAnalytics{
			TotalInterviews: 0,
			LastUpdated:     time.Now(),
		}, nil
	}

	global := &models.GlobalAnalytics{
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
		global.BestInterviewScore = bestScore

		global.WorstInterviewID = worst.ID
		global.WorstInterviewScore = worstScore
	}

	return global, nil
}

func (s *Service) GetInterviewAnalytics(id int64) (models.InterviewAnalytics, error) {
	interview, err := s.GetInterview(id)
	if err != nil {
		return models.InterviewAnalytics{}, err
	}

	return s.calculateAnalytics(*interview), nil
}

func (s *Service) GetAllInterviewAnalytics(filters *models.GetInterviewsFilters) ([]models.InterviewAnalytics, error) {
	interviews, err := s.GetAllInterviews(filters)
	if err != nil {
		return nil, fmt.Errorf("failed to get all interviews: %w", err)
	}
	result := make([]models.InterviewAnalytics, 0, len(interviews))

	for _, interview := range interviews {
		result = append(result, s.calculateAnalytics(interview))
	}

	return result, nil
}
