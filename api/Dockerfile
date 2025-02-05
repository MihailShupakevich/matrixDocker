# Используем официальный образ Golang как базовый
FROM golang:1.23

# Устанавливаем рабочую директорию
WORKDIR /api

# Копируем go.mod и go.sum
COPY go.mod go.sum ./

# Загружаем зависимости
RUN go mod download

# Копируем остальные файлы
COPY . .

# Устанавливаем целевую ОС и архитектуру
ENV GOOS=linux
ENV GOARCH=amd64

# Собираем приложение
RUN go build -o main ./cmd/main.go


# Указываем, что контейнер слушает на порту 8080
EXPOSE 8080
# Указываем команду для запуска приложения
CMD ["./main"]