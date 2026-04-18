FROM golang:1.26-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /out/ai-service ./services/ai-service

FROM alpine:3.20
WORKDIR /app
COPY --from=builder /out/ai-service /usr/local/bin/ai-service
EXPOSE 8083
CMD ["ai-service"]

