###########
# BUILDER #
###########

# Скачивание образа
FROM golang:alpine as builder

# Установка рабочей директории
WORKDIR /usr/src/app

# Установка переменных окружения
ENV CGO_ENABLED 0
ENV GOOS linux

# Установка зависимостей
RUN go install github.com/swaggo/swag/cmd/swag@latest
COPY ./go.mod ./go.sum ./
RUN go mod download && go mod verify

# Копирование проекта
COPY . ./

# Сборка
RUN swag init -g server.go \
  && go build -v -o /usr/local/bin/finapp ./server.go

#########
# FINAL #
#########

# Скачивание образа
FROM alpine:3

# Установка рабочей директории
ENV APP_HOME=/home/app
WORKDIR $APP_HOME

# Создание пользователя app
RUN addgroup -S app && adduser -S app -G app

# Копирование entrypoint.sh
COPY ./entrypoint.sh $APP_HOME
RUN sed -i 's/\r$//g' $APP_HOME/entrypoint.sh && chmod +x $APP_HOME/entrypoint.sh

# Установка приложения и передача владения файлами пользовалелю app
COPY --from=builder /usr/local/bin/finapp $APP_HOME
RUN chmod +x $APP_HOME/finapp && chown -R app:app $APP_HOME

# Смена пользователя
USER app

# Запуск
ENTRYPOINT ["/home/app/entrypoint.sh"]
