FROM golang:1.25.2-alpine AS builder

WORKDIR /build

# Установка зависимостей для сборки
RUN apk add --no-cache git ca-certificates

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Компиляция приложения (весь пакет cmd/app)
RUN CGO_ENABLED=0 GOOS=linux go build -o app ./cmd/app
# Компиляция мигратора (весь пакет cmd/migrator)
RUN CGO_ENABLED=0 GOOS=linux go build -o migrator ./cmd/migrator

FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

COPY --from=builder /build/app .
COPY --from=builder /build/migrator .
COPY --from=builder /build/config ./config
COPY --from=builder /build/migrations ./migrations

EXPOSE 8080

# Для сервиса app команда будет переопределена в docker-compose,
# поэтому здесь можно оставить заглушку или убрать CMD.
# CMD ["./app"] # закомментировано, так как в compose будет своя команда