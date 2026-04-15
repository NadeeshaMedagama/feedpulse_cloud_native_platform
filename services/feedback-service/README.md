# Feedback Service

`feedback-service` is the core business service.

## Responsibilities

- Public feedback submission
- Admin feedback listing and filtering
- Status transitions (`New`, `In Review`, `Resolved`)
- Delete and get-by-id operations
- Weekly summary generation trigger
- Manual AI re-analysis trigger

## Endpoints

- `GET /health`
- `POST /api/feedback` (public, rate-limited)
- `GET /api/feedback` (admin)
- `GET /api/feedback/{id}` (admin)
- `PATCH /api/feedback/{id}` (admin)
- `DELETE /api/feedback/{id}` (admin)
- `POST /api/feedback/{id}/reanalyze` (admin)
- `GET /api/feedback/summary` (admin)

## Environment variables

- `FEEDBACK_SERVICE_PORT` (default `8082`)
- `MONGO_URI` (Atlas preferred)
- `MONGO_DATABASE`
- `AI_SERVICE_URL` (default `http://localhost:8083`)
- `JWT_SECRET`
- `ADMIN_EMAIL`
- `ADMIN_PASSWORD`

## Run locally

```bash
go run ./services/feedback-service
```

