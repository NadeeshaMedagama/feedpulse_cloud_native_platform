# AI Service

`ai-service` isolates Gemini API concerns.

## Endpoints

- `GET /health`
- `POST /internal/ai/analyze`
- `POST /internal/ai/summary`

## Environment variables

- `AI_SERVICE_PORT` (default `8083`)
- `GEMINI_API_KEY`

## Run locally

```bash
go run ./services/ai-service
```

## Notes

- Returns structured JSON for analysis fields.
- If Gemini is unavailable or key is missing, caller service should handle gracefully.

