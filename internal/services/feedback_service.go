package services

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"FeedPulse_Cloud_Native_Platform/internal/models"
	"FeedPulse_Cloud_Native_Platform/internal/repository"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type FeedbackService struct {
	repo   repository.FeedbackRepository
	gemini GeminiService
}

func NewFeedbackService(repo repository.FeedbackRepository, gemini GeminiService) *FeedbackService {
	return &FeedbackService{repo: repo, gemini: gemini}
}

func (s *FeedbackService) Create(ctx context.Context, req models.FeedbackCreateRequest) (*models.Feedback, error) {
	if err := validateFeedback(req); err != nil {
		return nil, err
	}
	feedback := &models.Feedback{
		Title:          strings.TrimSpace(req.Title),
		Description:    strings.TrimSpace(req.Description),
		Category:       normalizeCategory(req.Category),
		SubmitterName:  strings.TrimSpace(req.SubmitterName),
		SubmitterEmail: strings.TrimSpace(req.SubmitterEmail),
	}

	if err := s.repo.Create(ctx, feedback); err != nil {
		return nil, err
	}

	go s.processAI(feedback.ID, feedback.Title, feedback.Description)
	return feedback, nil
}

func (s *FeedbackService) processAI(id primitive.ObjectID, title, description string) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	analysis, err := s.gemini.AnalyzeFeedback(ctx, title, description)
	if err != nil {
		return
	}
	_ = s.repo.UpdateAIFields(ctx, id, analysis)
}

func (s *FeedbackService) Reanalyze(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid feedback id")
	}
	feedback, err := s.repo.GetByID(ctx, objectID)
	if err != nil {
		return err
	}
	analysis, err := s.gemini.AnalyzeFeedback(ctx, feedback.Title, feedback.Description)
	if err != nil {
		return err
	}
	return s.repo.UpdateAIFields(ctx, objectID, analysis)
}

func (s *FeedbackService) GetByID(ctx context.Context, id string) (*models.Feedback, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("invalid feedback id")
	}
	return s.repo.GetByID(ctx, objectID)
}

func (s *FeedbackService) List(ctx context.Context, filters models.FeedbackFilters) ([]models.Feedback, int64, error) {
	return s.repo.List(ctx, filters)
}

func (s *FeedbackService) UpdateStatus(ctx context.Context, id string, status string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid feedback id")
	}
	status = strings.TrimSpace(status)
	if status != models.StatusNew && status != models.StatusInReview && status != models.StatusResolved {
		return errors.New("invalid status")
	}
	return s.repo.UpdateStatus(ctx, objectID, status)
}

func (s *FeedbackService) Delete(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid feedback id")
	}
	return s.repo.Delete(ctx, objectID)
}

func (s *FeedbackService) Stats(ctx context.Context) (map[string]interface{}, error) {
	total, err := s.repo.CountByStatus(ctx, models.StatusNew)
	if err != nil {
		return nil, err
	}
	inReview, err := s.repo.CountByStatus(ctx, models.StatusInReview)
	if err != nil {
		return nil, err
	}
	resolved, err := s.repo.CountByStatus(ctx, models.StatusResolved)
	if err != nil {
		return nil, err
	}
	avgPriority, err := s.repo.AveragePriority(ctx)
	if err != nil {
		return nil, err
	}
	commonTag, err := s.repo.MostCommonTag(ctx)
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"totalFeedback":   total + inReview + resolved,
		"openItems":       total + inReview,
		"averagePriority": avgPriority,
		"mostCommonTag":   commonTag,
	}, nil
}

func (s *FeedbackService) WeeklySummary(ctx context.Context) (string, error) {
	items, err := s.repo.LastSevenDays(ctx)
	if err != nil {
		return "", err
	}
	if len(items) == 0 {
		return "No feedback submitted in the last 7 days.", nil
	}
	var b strings.Builder
	for _, item := range items {
		b.WriteString(fmt.Sprintf("- %s: %s\n", item.Title, item.Description))
	}
	summary, err := s.gemini.SummarizeThemes(ctx, b.String())
	if err != nil {
		return "AI summary currently unavailable, but feedback is saved and accessible.", nil
	}
	return summary, nil
}

func validateFeedback(req models.FeedbackCreateRequest) error {
	title := strings.TrimSpace(req.Title)
	description := strings.TrimSpace(req.Description)
	if title == "" {
		return errors.New("title is required")
	}
	if len(title) > 120 {
		return errors.New("title cannot exceed 120 characters")
	}
	if len(description) < 20 {
		return errors.New("description must be at least 20 characters")
	}
	if req.SubmitterEmail != "" {
		re := regexp.MustCompile(`^[^\s@]+@[^\s@]+\.[^\s@]+$`)
		if !re.MatchString(strings.TrimSpace(req.SubmitterEmail)) {
			return errors.New("invalid email format")
		}
	}
	return nil
}

func normalizeCategory(category string) string {
	switch strings.TrimSpace(category) {
	case models.CategoryBug:
		return models.CategoryBug
	case models.CategoryFeatureRequest:
		return models.CategoryFeatureRequest
	case models.CategoryImprovement:
		return models.CategoryImprovement
	default:
		return models.CategoryOther
	}
}

func IsNotFound(err error) bool {
	return errors.Is(err, mongo.ErrNoDocuments)
}
