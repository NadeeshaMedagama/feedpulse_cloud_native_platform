FROM golang:1.26-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /out/api-gateway ./services/api-gateway

FROM alpine:3.20
WORKDIR /app
COPY --from=builder /out/api-gateway /usr/local/bin/api-gateway
COPY web ./web
EXPOSE 8080
CMD ["api-gateway"]

