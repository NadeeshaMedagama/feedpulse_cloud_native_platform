# Auth Service

`auth-service` manages admin authentication.

## Endpoints

- `GET /health`
- `POST /api/auth/login`

## Input

```json
{ "email": "admin@feedpulse.local", "password": "Admin123!" }
```

## Output

```json
{ "success": true, "data": { "token": "<jwt>" }, "message": "login successful" }
```

## Environment variables

- `AUTH_SERVICE_PORT` (default `8081`)
- `JWT_SECRET`
- `ADMIN_EMAIL`
- `ADMIN_PASSWORD`

## Run locally

```bash
go run ./services/auth-service
```

