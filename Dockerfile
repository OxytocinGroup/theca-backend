FROM golang:1.23.4-bullseye AS builder

WORKDIR /app

COPY go.mod go.sum ./
COPY vendor ./vendor
COPY ./internal/utils/email ./internal/utils/email
COPY . .

ENV GOFLAGS="-mod=vendor"

RUN make build

FROM debian:buster-slim

WORKDIR /app

COPY --from=builder /build/bin .
COPY .env ./.env

EXPOSE 3000

CMD [ "./api" ]
