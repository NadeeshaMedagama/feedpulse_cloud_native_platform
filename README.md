# FeedPulse_Cloud_Native_Platform 

FeedPulse_Cloud_Native_Platform is an AI-powered product feedback platform built with Go, MongoDB Atlas, and Gemini.
This repository now includes a **microservices-oriented architecture** with a professional frontend, dark mode UI, admin dashboard, and Docker Compose setup.

## Architecture

### Services

- `services/api-gateway` - public entrypoint, serves UI and proxies API traffic
- `services/auth-service` - admin login and JWT token issuance
- `services/feedback-service` - feedback lifecycle, filters, stats, status updates, MongoDB persistence
- `services/ai-service` - Gemini adapter for analysis and weekly summary generation

### Shared backend components

- `internal/config` - env loading and defaults (Atlas-first DB strategy)
- `internal/handlers` - HTTP controllers
- `internal/services` - business rules and integration orchestration
- `internal/repository` - MongoDB repository contracts and implementation
- `internal/middleware` - JWT + rate limit + response helper
- `internal/models` - request/response and domain data contracts

### Frontend

- `frontend/` - Next.js 14 App Router frontend
- `frontend/app/page.tsx` - public submit page
- `frontend/app/admin/page.tsx` - admin login + dashboard
- `frontend/app/api/[...path]/route.ts` - API proxy to gateway

## Feature coverage

### Must-have coverage implemented

- Public feedback submission page with form validation
- Required fields and description minimum length checks
- Save feedback to MongoDB with schema validations
- Consistent JSON response envelope: `{ success, data, error, message }`
- Gemini AI processing support with graceful failure handling
- Admin login with JWT-based route protection
- Admin list view, status updates, search, filtering, sorting, and pagination
- REST endpoints for create/list/get/update/delete/summary/login
- MongoDB indexes for `status`, `category`, `ai_priority`, `createdAt`

### Nice-to-have coverage implemented

- Description character counter on public form
- Rate limit for public submissions (5/hour/IP)
- Manual AI re-analysis endpoint/action
- 7-day AI trend summary endpoint/action
- Stats strip in admin dashboard (total/open/avg-priority/top-tag)

## API map

- `POST /api/auth/login`
- `POST /api/feedback` (public)
- `GET /api/feedback` (admin)
- `GET /api/feedback/{id}` (admin)
- `PATCH /api/feedback/{id}` (admin)
- `DELETE /api/feedback/{id}` (admin)
- `POST /api/feedback/{id}/reanalyze` (admin)
- `GET /api/feedback/summary` (admin)

## Environment setup

Copy `.env.example` to `.env` and set values:

- `JWT_SECRET`
- `ADMIN_EMAIL`
- `ADMIN_PASSWORD`
- `MONGO_URI` (MongoDB Atlas URI, primary)
- `MONGO_LOCAL_URI` (optional fallback for local mode)
- `MONGO_DATABASE`
- `GEMINI_API_KEY`
- service/gateway ports and service URLs (already included in `.env.example`)

## Run options

### Option 1: Microservices with Docker Compose (recommended)

```bash
docker compose up --build
```

Open:

- `http://localhost:3060` - public page (Next.js)
- `http://localhost:3060/admin` - admin dashboard (Next.js)
- `http://localhost:8080` - API gateway (backend entrypoint)

Notes:

- Compose defaults `MONGO_URI` to local mongo container if not set.
- For Atlas in Docker mode, set `MONGO_URI` in `.env` to your Atlas connection string.

### Option 2: Local microservices (without Docker)

```bash
go mod tidy
go test ./...
go run ./services/ai-service
go run ./services/auth-service
go run ./services/feedback-service
go run ./services/api-gateway
```

### Option 3: Legacy monolith mode

```bash
go run ./cmd/feedpulse-api
```

## Docker files included

- `docker-compose.yml`
- `services/api-gateway/Dockerfile`
- `services/auth-service/Dockerfile`
- `services/feedback-service/Dockerfile`
- `services/ai-service/Dockerfile`
- `frontend/Dockerfile`

## Kubernetes deployment

Kubernetes manifests are available in `k8s/` with Kustomize overlays:

- `k8s/overlays/atlas` (recommended: MongoDB Atlas primary)
- `k8s/overlays/local-mongo` (development fallback)

Quick deploy example:

```bash
kubectl apply -k k8s/overlays/atlas
kubectl -n feedpulse get pods
```

See full Kubernetes deployment guide in `k8s/README.md`.

## Service/component documentation

Detailed READMEs were added for:

- Each service folder under `services/*`
- Shared backend component folders under `internal/*`
- Frontend assets folder `web/`
- Legacy monolith entrypoint folder `cmd/`

## Testing

Current automated tests focus on service validation/auth logic:

```bash
go test ./...
```

## Submission checklist support

- `.gitignore` includes `node_modules`, `.env`, and build outputs
- `.env.example` provided for safe setup
- README includes architecture, setup, and API usage

## If more time were available

- Add async queue for AI processing and retry strategy
- Add full HTTP integration tests for endpoints
- Add role-based admin users collection
- Add centralized structured logging + traces
- Add CI/CD workflow and Kubernetes manifests

