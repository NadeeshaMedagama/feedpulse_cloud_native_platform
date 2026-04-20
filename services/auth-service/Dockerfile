FROM golang:1.26-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /out/auth-service ./services/auth-service

FROM alpine:3.23
WORKDIR /app
COPY --from=builder /out/auth-service /usr/local/bin/auth-service
EXPOSE 8081
CMD ["auth-service"]

