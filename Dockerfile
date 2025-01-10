# Stage 1: Build
FROM golang:1.23.4-bullseye AS builder

WORKDIR /app

# Копируем зависимости
COPY go.mod go.sum ./
COPY vendor ./vendor

# Копируем исходники

COPY . .

# Устанавливаем флаги для сборки
ENV GOFLAGS="-mod=vendor"

# Сборка приложения
RUN make build

# Stage 2: Runtime
FROM debian:buster-slim

WORKDIR /app

# Устанавливаем корневые сертификаты
RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates \
    && rm -rf /var/lib/apt/lists/*

# Копируем собранное приложение из builder
COPY --from=builder /app/build/bin/ .

COPY internal/utils/email/verifyMail.html internal/utils/email/resetEmail.html ./templates/
COPY .env ./.env

EXPOSE 3000

# Запуск приложения
CMD [ "./api" ]
