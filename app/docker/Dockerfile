# Скачивание образа
FROM golang:alpine

# Установка рабочей директории
WORKDIR /usr/src/finapp

# Установка зависимостей
RUN go install github.com/go-task/task/v3/cmd/task@latest

# Копирование проекта
COPY . ./

# Установка приложения
RUN task install

# Копирование entrypoint.sh
COPY docker/entrypoint.sh /usr/local/bin
#RUN chmod +x /usr/local/bin/entrypoint.sh

# Запуск
ENTRYPOINT ["/usr/local/bin/entrypoint.sh"]
CMD ["finapp", "app:serve"]
