# Этап сборки (builder)
FROM golang:1.24-alpine AS builder

# Установка зависимостей для сборки
RUN apk add --no-cache git make

# Рабочая директория
WORKDIR /app

# Копируем файлы модулей для кэширования
COPY go.mod go.sum ./

# Скачиваем зависимости
RUN go mod download

# Копируем исходный код
COPY . .

# Собираем приложение
# Флаги:
# - CGO_ENABLED=0 - отключаем CGO для полностью статического бинарника
# - ldflags="-s -w" - уменьшаем размер бинарника
# Build the application (two levels down)
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" \
    -o /app/final_task_cmd ./final_task/cmd

# Этап запуска (runtime)
FROM alpine:3.21

# Устанавливаем tzdata для работы с временными зонами
RUN apk add --no-cache tzdata ca-certificates

# Рабочая директория
WORKDIR /app

# Копируем бинарник из этапа сборки
COPY --from=builder /app/final_task_cmd /app/bin/final_task
# Копируем конфигурационные файлы
COPY --from=builder /app/final_task/cmd/docker-config.yaml /app/bin/config.yaml
COPY --from=builder /app/final_task/SQL/migrations /app/migrations

# Переменные окружения по умолчанию
ENV APP_ENV=production \
    TZ=UTC

EXPOSE 8080

WORKDIR /app/bin

# Команда запуска
CMD ["./final_task"]