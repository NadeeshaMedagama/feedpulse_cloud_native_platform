# Internal Components

This folder contains reusable backend building blocks shared across services.

## Subfolders

- `config` - environment loading and default fallbacks
- `handlers` - HTTP layer and request/response mapping
- `middleware` - JWT validation, rate limiting, JSON response helpers
- `models` - DTOs and domain models
- `repository` - MongoDB persistence contracts and implementations
- `server` - monolith router composition (legacy single-service mode)
- `services` - application business logic and integrations

