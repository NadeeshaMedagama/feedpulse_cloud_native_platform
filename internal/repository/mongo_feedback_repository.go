package repository

import (
	"context"
	"time"

	"FeedPulse_Cloud_Native_Platform/internal/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoFeedbackRepository struct {
	collection *mongo.Collection
}

func NewMongoFeedbackRepository(db *mongo.Database) *MongoFeedbackRepository {
	return &MongoFeedbackRepository{collection: db.Collection("feedback")}
}

func (r *MongoFeedbackRepository) Create(ctx context.Context, feedback *models.Feedback) error {
	now := time.Now().UTC()
	feedback.CreatedAt = now
	feedback.UpdatedAt = now
	feedback.Status = models.StatusNew
	feedback.AIProcessed = false

	result, err := r.collection.InsertOne(ctx, feedback)
	if err != nil {
		return err
	}
	feedback.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

func (r *MongoFeedbackRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*models.Feedback, error) {
	var feedback models.Feedback
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&feedback)
	if err != nil {
		return nil, err
	}
	return &feedback, nil
}

func (r *MongoFeedbackRepository) List(ctx context.Context, filters models.FeedbackFilters) ([]models.Feedback, int64, error) {
	query := bson.M{}
	if filters.Category != "" {
		query["category"] = filters.Category
	}
	if filters.Status != "" {
		query["status"] = filters.Status
	}
	if filters.Search != "" {
		query["$or"] = []bson.M{{"title": bson.M{"$regex": filters.Search, "$options": "i"}}, {"ai_summary": bson.M{"$regex": filters.Search, "$options": "i"}}}
	}

	if filters.Page < 1 {
		filters.Page = 1
	}
	if filters.Limit < 1 {
		filters.Limit = 10
	}

	sortField := "createdAt"
	switch filters.SortBy {
	case "date":
		sortField = "createdAt"
	case "priority":
		sortField = "ai_priority"
	case "sentiment":
		sortField = "ai_sentiment"
	}
	order := int32(-1)
	if filters.Order == "asc" {
		order = 1
	}

	findOptions := options.Find().
		SetSort(bson.D{{Key: sortField, Value: order}}).
		SetSkip((filters.Page - 1) * filters.Limit).
		SetLimit(filters.Limit)

	cursor, err := r.collection.Find(ctx, query, findOptions)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var results []models.Feedback
	if err := cursor.All(ctx, &results); err != nil {
		return nil, 0, err
	}

	total, err := r.collection.CountDocuments(ctx, query)
	if err != nil {
		return nil, 0, err
	}

	return results, total, nil
}

func (r *MongoFeedbackRepository) UpdateStatus(ctx context.Context, id primitive.ObjectID, status string) error {
	result, err := r.collection.UpdateByID(ctx, id, bson.M{"$set": bson.M{"status": status, "updatedAt": time.Now().UTC()}})
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return mongo.ErrNoDocuments
	}
	return nil
}

func (r *MongoFeedbackRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return mongo.ErrNoDocuments
	}
	return nil
}

func (r *MongoFeedbackRepository) UpdateAIFields(ctx context.Context, id primitive.ObjectID, analysis models.GeminiAnalysis) error {
	_, err := r.collection.UpdateByID(ctx, id, bson.M{
		"$set": bson.M{
			"ai_category":  analysis.Category,
			"ai_sentiment": analysis.Sentiment,
			"ai_priority":  analysis.PriorityScore,
			"ai_summary":   analysis.Summary,
			"ai_tags":      analysis.Tags,
			"ai_processed": true,
			"updatedAt":    time.Now().UTC(),
		},
	})
	return err
}

func (r *MongoFeedbackRepository) LastSevenDays(ctx context.Context) ([]models.Feedback, error) {
	return r.CreatedAfter(ctx, time.Now().UTC().AddDate(0, 0, -7))
}

func (r *MongoFeedbackRepository) CreatedAfter(ctx context.Context, after time.Time) ([]models.Feedback, error) {
	cursor, err := r.collection.Find(ctx, bson.M{"createdAt": bson.M{"$gte": after}})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []models.Feedback
	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}
	return results, nil
}

func (r *MongoFeedbackRepository) EnsureIndexes(ctx context.Context) error {
	indexes := []mongo.IndexModel{
		{Keys: bson.D{{Key: "status", Value: 1}}},
		{Keys: bson.D{{Key: "category", Value: 1}}},
		{Keys: bson.D{{Key: "ai_priority", Value: -1}}},
		{Keys: bson.D{{Key: "createdAt", Value: -1}}},
	}
	_, err := r.collection.Indexes().CreateMany(ctx, indexes)
	return err
}

func (r *MongoFeedbackRepository) CountByStatus(ctx context.Context, status string) (int64, error) {
	return r.collection.CountDocuments(ctx, bson.M{"status": status})
}

func (r *MongoFeedbackRepository) AveragePriority(ctx context.Context) (float64, error) {
	pipeline := []bson.M{{"$match": bson.M{"ai_priority": bson.M{"$gt": 0}}}, {"$group": bson.M{"_id": nil, "avg": bson.M{"$avg": "$ai_priority"}}}}
	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return 0, err
	}
	defer cursor.Close(ctx)

	var result []bson.M
	if err := cursor.All(ctx, &result); err != nil {
		return 0, err
	}
	if len(result) == 0 {
		return 0, nil
	}
	if avg, ok := result[0]["avg"].(float64); ok {
		return avg, nil
	}
	return 0, nil
}

func (r *MongoFeedbackRepository) MostCommonTag(ctx context.Context) (string, error) {
	pipeline := []bson.M{
		{"$unwind": "$ai_tags"},
		{"$group": bson.M{"_id": "$ai_tags", "count": bson.M{"$sum": 1}}},
		{"$sort": bson.M{"count": -1}},
		{"$limit": 1},
	}
	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return "", err
	}
	defer cursor.Close(ctx)

	var result []bson.M
	if err := cursor.All(ctx, &result); err != nil {
		return "", err
	}
	if len(result) == 0 {
		return "-", nil
	}
	tag, _ := result[0]["_id"].(string)
	if tag == "" {
		return "-", nil
	}
	return tag, nil
}
