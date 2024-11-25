# Используем минимальный образ Go
FROM golang:1.23.1-alpine AS builder

# Устанавливаем необходимые зависимости
RUN apk add --no-cache git

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем и устанавливаем зависимости
COPY go.mod go.sum ./
RUN go mod download

# Копируем весь код
COPY . .

# Компилируем приложение
RUN go build -o main ./cmd

# Используем минимальный образ для продакшн
FROM alpine:latest

# Устанавливаем зависимости для запуска приложения
RUN apk add --no-cache ca-certificates

# Устанавливаем рабочую директорию
WORKDIR /root/

# Копируем бинарный файл из builder
COPY --from=builder /app/main .

# Копируем папку с миграциями
COPY --from=builder /app/migrations /migrations

# Указываем порт
EXPOSE ${APP_PORT}

# Запускаем приложение
CMD ["./main"]
