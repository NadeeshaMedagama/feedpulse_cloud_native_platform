package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	CategoryBug            = "Bug"
	CategoryFeatureRequest = "Feature Request"
	CategoryImprovement    = "Improvement"
	CategoryOther          = "Other"
)

const (
	StatusNew      = "New"
	StatusInReview = "In Review"
	StatusResolved = "Resolved"
)

type Feedback struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Title          string             `bson:"title" json:"title"`
	Description    string             `bson:"description" json:"description"`
	Category       string             `bson:"category" json:"category"`
	Status         string             `bson:"status" json:"status"`
	SubmitterName  string             `bson:"submitterName,omitempty" json:"submitterName,omitempty"`
	SubmitterEmail string             `bson:"submitterEmail,omitempty" json:"submitterEmail,omitempty"`
	AICategory     string             `bson:"ai_category,omitempty" json:"ai_category,omitempty"`
	AISentiment    string             `bson:"ai_sentiment,omitempty" json:"ai_sentiment,omitempty"`
	AIPriority     int                `bson:"ai_priority,omitempty" json:"ai_priority,omitempty"`
	AISummary      string             `bson:"ai_summary,omitempty" json:"ai_summary,omitempty"`
	AITags         []string           `bson:"ai_tags,omitempty" json:"ai_tags,omitempty"`
	AIProcessed    bool               `bson:"ai_processed" json:"ai_processed"`
	CreatedAt      time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt      time.Time          `bson:"updatedAt" json:"updatedAt"`
}

type FeedbackCreateRequest struct {
	Title          string `json:"title"`
	Description    string `json:"description"`
	Category       string `json:"category"`
	SubmitterName  string `json:"name"`
	SubmitterEmail string `json:"email"`
}

type FeedbackUpdateStatusRequest struct {
	Status string `json:"status"`
}

type FeedbackFilters struct {
	Category string
	Status   string
	Search   string
	SortBy   string
	Order    string
	Page     int64
	Limit    int64
}

type GeminiAnalysis struct {
	Category      string   `json:"category"`
	Sentiment     string   `json:"sentiment"`
	PriorityScore int      `json:"priority_score"`
	Summary       string   `json:"summary"`
	Tags          []string `json:"tags"`
}
