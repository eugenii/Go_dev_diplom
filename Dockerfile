# Stage 1: Собираем бинарник внутри контейнера (на базе golang)
FROM golang:1.25.3 AS builder

WORKDIR /app

# Копируем только go.mod и go.sum для кэширования зависимостей
COPY go.mod go.sum ./

# Устанавливаем зависимости
RUN go mod download

# Копируем остальной код
COPY . .

# Собираем статически-linked бинарник для Linux (без CGO — безопаснее и проще)
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o scheduler .

# Stage 2: Runtime — Ubuntu (как рекомендовано)
FROM ubuntu:22.04

# Установка минимальных системных пакетов
RUN apt-get update && \
    apt-get install -y --no-install-recommends \
        ca-certificates tzdata && \
    rm -rf /var/lib/apt/lists/*

# Создаём непривилегированного пользователя
RUN useradd --create-home --shell /bin/bash appuser
USER appuser
WORKDIR /home/appuser

# Копируем только нужное: бинарник и папку web
COPY --from=builder /app/scheduler .
COPY --from=builder /app/web/ web/

# Пробрасываем порт (как в старом рабочем Dockerfile)
EXPOSE 7540

# Запуск
CMD ["./scheduler"]