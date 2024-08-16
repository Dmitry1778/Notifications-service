# Используем базовый образ Go
FROM golang:1.21-alpine as builder

# Устанавливаем рабочую директорию
WORKDIR /app

# Устанавливаем зависимости Alpine
RUN apk --no-cache add bash git

# Копируем go.mod и go.sum
COPY ["go.mod", "go.sum", "./"]
RUN go mod download

#build
COPY . ./
RUN go build -o ./bin/notify cmd/main.go

FROM alpine AS runner

COPY --from=builder /app/bin/notify /
# Устанавливаем конфигурационные файлы
COPY config.yaml local.yaml ./

# Запускаем приложение
CMD ["./notify"]