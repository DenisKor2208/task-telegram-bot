FROM golang:1.25.2-alpine AS builder

WORKDIR /app

# Установка зависимостей для сборки
RUN apk add --no-cache git ca-certificates

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Компиляция приложения
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/app/main.go

FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root/

COPY --from=builder /app/main .

# Копирование файлов миграций
COPY migrations/ /migrations/

EXPOSE 8080

CMD ["./main"]
