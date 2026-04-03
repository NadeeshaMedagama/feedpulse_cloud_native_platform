FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /out/feedback-service ./services/feedback-service

FROM alpine:3.20
WORKDIR /app
COPY --from=builder /out/feedback-service /usr/local/bin/feedback-service
EXPOSE 8082
CMD ["feedback-service"]

