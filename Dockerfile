FROM golang:1.25-alpine AS builder

WORKDIR /app

# Устанавливаем зависимости
RUN apk add --no-cache git

# Устанавливаем goose
RUN go install github.com/pressly/goose/v3/cmd/goose@latest

# Копируем go mod файлы
COPY go.mod go.sum ./
RUN go mod download

# Копируем исходный код
COPY . .

# Копируем миграции
COPY migrations ./migrations

# Собираем приложение
RUN go build -o app ./cmd

# Финальный образ
FROM alpine:3.18

WORKDIR /app

# Устанавливаем зависимости
RUN apk add --no-cache postgresql-client

# Копируем бинарник из builder stage
COPY --from=builder /app/app .
COPY --from=builder /app/migrations ./migrations
COPY --from=builder /go/bin/goose /usr/local/bin/goose

# Создаем non-root пользователя
RUN adduser -D -s /bin/sh appuser && chown -R appuser:appuser /app
USER appuser

EXPOSE 8080

CMD ["./app"]