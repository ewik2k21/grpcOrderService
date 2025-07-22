# Этап сборки
FROM golang:1.24.4 AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o order-service ./cmd/app/main.go

# Финальный этап
FROM gcr.io/distroless/base:latest
WORKDIR /app
COPY --from=builder /app/order-service .
EXPOSE 50052
CMD ["./order-service"]