package repository

import (
	"context"
	"time"

	"FeedPulse_Cloud_Native_Platform/internal/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type FeedbackRepository interface {
	Create(ctx context.Context, feedback *models.Feedback) error
	GetByID(ctx context.Context, id primitive.ObjectID) (*models.Feedback, error)
	List(ctx context.Context, filters models.FeedbackFilters) ([]models.Feedback, int64, error)
	UpdateStatus(ctx context.Context, id primitive.ObjectID, status string) error
	Delete(ctx context.Context, id primitive.ObjectID) error
	UpdateAIFields(ctx context.Context, id primitive.ObjectID, analysis models.GeminiAnalysis) error
	LastSevenDays(ctx context.Context) ([]models.Feedback, error)
	EnsureIndexes(ctx context.Context) error
	CountByStatus(ctx context.Context, status string) (int64, error)
	AveragePriority(ctx context.Context) (float64, error)
	MostCommonTag(ctx context.Context) (string, error)
	CreatedAfter(ctx context.Context, after time.Time) ([]models.Feedback, error)
}
