package server

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"FeedPulse_Cloud_Native_Platform/internal/handlers"
	"FeedPulse_Cloud_Native_Platform/internal/models"
	"FeedPulse_Cloud_Native_Platform/internal/repository"
	"FeedPulse_Cloud_Native_Platform/internal/services"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type fakeFeedbackRepo struct {
	mu           sync.Mutex
	items        map[primitive.ObjectID]*models.Feedback
	updateAICh   chan struct{}
	lastStatusID primitive.ObjectID
	lastStatus   string
}

var _ repository.FeedbackRepository = (*fakeFeedbackRepo)(nil)

func newFakeFeedbackRepo() *fakeFeedbackRepo {
	return &fakeFeedbackRepo{
		items:      map[primitive.ObjectID]*models.Feedback{},
		updateAICh: make(chan struct{}, 1),
	}
}

func (r *fakeFeedbackRepo) Create(ctx context.Context, feedback *models.Feedback) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if feedback.ID.IsZero() {
		feedback.ID = primitive.NewObjectID()
	}
	feedback.Status = models.StatusNew
	r.items[feedback.ID] = feedback
	return nil
}

func (r *fakeFeedbackRepo) GetByID(ctx context.Context, id primitive.ObjectID) (*models.Feedback, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	item, ok := r.items[id]
	if !ok {
		return nil, mongo.ErrNoDocuments
	}
	copy := *item
	return &copy, nil
}

func (r *fakeFeedbackRepo) List(ctx context.Context, filters models.FeedbackFilters) ([]models.Feedback, int64, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	out := make([]models.Feedback, 0, len(r.items))
	for _, v := range r.items {
		out = append(out, *v)
	}
	return out, int64(len(out)), nil
}

func (r *fakeFeedbackRepo) UpdateStatus(ctx context.Context, id primitive.ObjectID, status string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	item, ok := r.items[id]
	if !ok {
		return mongo.ErrNoDocuments
	}
	item.Status = status
	r.lastStatusID = id
	r.lastStatus = status
	return nil
}

func (r *fakeFeedbackRepo) Delete(ctx context.Context, id primitive.ObjectID) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.items[id]; !ok {
		return mongo.ErrNoDocuments
	}
	delete(r.items, id)
	return nil
}

func (r *fakeFeedbackRepo) UpdateAIFields(ctx context.Context, id primitive.ObjectID, analysis models.GeminiAnalysis) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	item, ok := r.items[id]
	if !ok {
		return mongo.ErrNoDocuments
	}
	item.AIProcessed = true
	item.AICategory = analysis.Category
	item.AISentiment = analysis.Sentiment
	item.AIPriority = analysis.PriorityScore
	item.AISummary = analysis.Summary
	item.AITags = analysis.Tags
	select {
	case r.updateAICh <- struct{}{}:
	default:
	}
	return nil
}

func (r *fakeFeedbackRepo) LastSevenDays(ctx context.Context) ([]models.Feedback, error) {
	return []models.Feedback{}, nil
}

func (r *fakeFeedbackRepo) EnsureIndexes(ctx context.Context) error { return nil }
func (r *fakeFeedbackRepo) CountByStatus(ctx context.Context, status string) (int64, error) {
	return 0, nil
}
func (r *fakeFeedbackRepo) AveragePriority(ctx context.Context) (float64, error) { return 0, nil }
func (r *fakeFeedbackRepo) MostCommonTag(ctx context.Context) (string, error)    { return "-", nil }
func (r *fakeFeedbackRepo) CreatedAfter(ctx context.Context, after time.Time) ([]models.Feedback, error) {
	return []models.Feedback{}, nil
}

type fakeGemini struct{}

func (g *fakeGemini) AnalyzeFeedback(ctx context.Context, title, description string) (models.GeminiAnalysis, error) {
	return models.GeminiAnalysis{
		Category:      models.CategoryFeatureRequest,
		Sentiment:     "Positive",
		PriorityScore: 8,
		Summary:       "User requests dark mode.",
		Tags:          []string{"UI", "Settings"},
	}, nil
}

func (g *fakeGemini) SummarizeThemes(ctx context.Context, text string) (string, error) {
	return "Top themes summary", nil
}

func newTestRouter(repo *fakeFeedbackRepo) http.Handler {
	feedbackService := services.NewFeedbackService(repo, &fakeGemini{})
	authService := services.NewAuthService("admin@feedpulse.local", "Admin123!", "test-secret")
	return NewRouter(handlers.NewAuthHandler(authService), handlers.NewFeedbackHandler(feedbackService), authService)
}

func TestPostFeedback_ValidSubmissionSavesAndTriggersAI(t *testing.T) {
	repo := newFakeFeedbackRepo()
	router := newTestRouter(repo)

	payload := map[string]string{
		"title":       "Need dark mode",
		"description": "Please add a dark mode option in settings for low-light environments.",
		"category":    models.CategoryFeatureRequest,
	}
	body, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/api/feedback", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	res := httptest.NewRecorder()

	router.ServeHTTP(res, req)
	if res.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d, body=%s", res.Code, res.Body.String())
	}

	select {
	case <-repo.updateAICh:
	case <-time.After(2 * time.Second):
		t.Fatal("expected AI processing to be triggered")
	}
}

func TestPostFeedback_RejectsEmptyTitle(t *testing.T) {
	repo := newFakeFeedbackRepo()
	router := newTestRouter(repo)

	payload := map[string]string{
		"title":       "",
		"description": "This description is long enough but title should fail validation.",
		"category":    models.CategoryBug,
	}
	body, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/api/feedback", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	res := httptest.NewRecorder()

	router.ServeHTTP(res, req)
	if res.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", res.Code)
	}
}

func TestPatchFeedbackStatus_WorksWithAuth(t *testing.T) {
	repo := newFakeFeedbackRepo()
	router := newTestRouter(repo)

	id := primitive.NewObjectID()
	repo.items[id] = &models.Feedback{
		ID:          id,
		Title:       "Sample",
		Description: "This description has enough characters.",
		Category:    models.CategoryOther,
		Status:      models.StatusNew,
	}

	loginPayload, _ := json.Marshal(models.LoginRequest{Email: "admin@feedpulse.local", Password: "Admin123!"})
	loginReq := httptest.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewReader(loginPayload))
	loginReq.Header.Set("Content-Type", "application/json")
	loginRes := httptest.NewRecorder()
	router.ServeHTTP(loginRes, loginReq)
	if loginRes.Code != http.StatusOK {
		t.Fatalf("expected login 200, got %d", loginRes.Code)
	}

	var loginResult models.APIResponse
	if err := json.Unmarshal(loginRes.Body.Bytes(), &loginResult); err != nil {
		t.Fatalf("failed to parse login response: %v", err)
	}
	data, ok := loginResult.Data.(map[string]interface{})
	if !ok {
		t.Fatal("login response token payload missing")
	}
	token, _ := data["token"].(string)
	if token == "" {
		t.Fatal("expected token in login response")
	}

	patchPayload, _ := json.Marshal(map[string]string{"status": models.StatusInReview})
	patchReq := httptest.NewRequest(http.MethodPatch, "/api/feedback/"+id.Hex(), bytes.NewReader(patchPayload))
	patchReq.Header.Set("Content-Type", "application/json")
	patchReq.Header.Set("Authorization", "Bearer "+token)
	patchRes := httptest.NewRecorder()
	router.ServeHTTP(patchRes, patchReq)
	if patchRes.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d, body=%s", patchRes.Code, patchRes.Body.String())
	}

	if repo.lastStatus != models.StatusInReview {
		t.Fatalf("expected status update to %q, got %q", models.StatusInReview, repo.lastStatus)
	}
}

func TestAuthMiddleware_RejectsUnauthenticatedRequests(t *testing.T) {
	repo := newFakeFeedbackRepo()
	router := newTestRouter(repo)

	req := httptest.NewRequest(http.MethodGet, "/api/feedback", nil)
	res := httptest.NewRecorder()
	router.ServeHTTP(res, req)

	if res.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", res.Code)
	}
}
