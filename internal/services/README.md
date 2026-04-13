# services

Business logic layer.

- `auth_service.go` - login and JWT generation/parsing
- `feedback_service.go` - validation, workflow orchestration, stats
- `gemini_service.go` - direct Gemini client (used by ai-service)
- `ai_client_service.go` - HTTP client to ai-service (used by feedback-service)

