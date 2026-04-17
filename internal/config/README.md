# config

Loads environment variables with practical defaults.

Primary logic:
- Prefer `MONGO_URI` (Atlas)
- Fallback to `MONGO_LOCAL_URI`
- Final fallback to `mongodb://localhost:27017`

