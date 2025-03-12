# 1. Используем официальный образ Go (попробуй явно указать версию)
FROM golang:1.23.4 AS builder

# 2. Устанавливаем рабочую директорию внутри контейнера
WORKDIR /app

# 3. Копируем файлы go.mod и go.sum, скачиваем зависимости
COPY go.mod go.sum ./
RUN go mod download

# 4. Копируем весь код проекта
COPY . .

# 5. Компилируем Go-приложение (если `main.go` в `cmd/`, укажи `./cmd/`)
RUN go build -o main ./cmd/app/

# 6. Создаем финальный контейнер
FROM golang:1.23.4 AS runner

WORKDIR /app

# 7. Копируем скомпилированный бинарник из builder-контейнера
COPY --from=builder /app/main .

# 8. Указываем, какой файл запускать при старте контейнера
CMD ["./main"]
