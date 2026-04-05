package services

import (
	"testing"

	"FeedPulse_Cloud_Native_Platform/internal/models"
)

func TestValidateFeedback(t *testing.T) {
	tests := []struct {
		name    string
		req     models.FeedbackCreateRequest
		wantErr bool
	}{
		{
			name: "valid payload",
			req: models.FeedbackCreateRequest{
				Title:       "Dark mode support",
				Description: "Please add dark mode support in the admin dashboard soon.",
				Category:    models.CategoryFeatureRequest,
			},
			wantErr: false,
		},
		{
			name: "missing title",
			req: models.FeedbackCreateRequest{
				Description: "This is a valid length description but title is missing.",
			},
			wantErr: true,
		},
		{
			name: "short description",
			req: models.FeedbackCreateRequest{
				Title:       "Short desc",
				Description: "Too short",
			},
			wantErr: true,
		},
		{
			name: "invalid email",
			req: models.FeedbackCreateRequest{
				Title:          "Bug found",
				Description:    "There is a serious issue with notifications not loading.",
				SubmitterEmail: "wrong-email",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateFeedback(tt.req)
			if (err != nil) != tt.wantErr {
				t.Fatalf("validateFeedback() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
