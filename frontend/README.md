# FeedPulse_Cloud_Native_Platform Frontend (Next.js)

This is the Next.js App Router frontend for `FeedPulse_Cloud_Native_Platform`.

## Pages

- `/` - Public feedback submission form
- `/admin` - Admin login and feedback dashboard
- `/api/*` - Next.js proxy routes forwarding to backend gateway

## Local run

```bash
npm install
npm run dev
```

Frontend runs on `http://localhost:3060`.

## Environment

- `BACKEND_INTERNAL_URL` (server-side proxy target)
  - Docker default: `http://api-gateway:8080`
  - Local default in code: `http://api-gateway:8080` (override for local non-docker)

## Build

```bash
npm run build
npm run start
```

## Test

```bash
npm test
```


