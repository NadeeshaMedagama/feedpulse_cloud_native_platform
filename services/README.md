# Services Overview

This folder contains deployable Go microservices for FeedPulse_Cloud_Native_Platform.

## Services

- `api-gateway` - public entrypoint, serves web UI, proxies API requests
- `auth-service` - handles admin login and JWT token issuance
- `feedback-service` - feedback CRUD, filters, stats, status updates, persistence
- `ai-service` - Gemini integration for analysis and trend summaries

## Communication flow

1. Browser calls `api-gateway` (`:8080`)
2. `api-gateway` routes:
   - `/api/auth/*` -> `auth-service`
   - `/api/feedback*` -> `feedback-service`
3. `feedback-service` calls `ai-service` for AI tasks
4. `feedback-service` persists data in MongoDB Atlas/local Mongo

## SOLID mapping

- Single Responsibility: handlers/services/repositories separated
- Open/Closed: services depend on interfaces where needed (`GeminiService`, repository contract)
- Liskov: interchangeable AI providers (Gemini direct or HTTP AI client)
- Interface Segregation: small focused interfaces in `internal/repository`
- Dependency Inversion: business logic depends on abstractions, not concrete DB/HTTP implementations

