# API Gateway

`api-gateway` is the user-facing entrypoint.

## Responsibilities

- Serves `web/index.html` and `web/admin.html`
- Serves static assets under `/web/*`
- Reverse proxies:
  - `/api/auth/*` -> auth-service
  - `/api/feedback*` -> feedback-service

## Run locally

```bash
go run ./services/api-gateway
```

## Environment variables

- `GATEWAY_PORT` (default `8080`)
- `AUTH_SERVICE_URL` (default `http://localhost:8081`)
- `FEEDBACK_SERVICE_URL` (default `http://localhost:8082`)

