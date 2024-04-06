###########
# BUILDER #
###########

# Скачивание образа
FROM golang:alpine as builder

# Установка рабочей директории
WORKDIR /usr/src/finapp

# Установка переменных окружения
ENV CGO_ENABLED 0
ENV GOOS linux

# Установка зависимостей
RUN go install github.com/go-task/task/v3/cmd/task@latest

# Копирование проекта
COPY . ./

# Установка приложения
RUN task install

#########
# FINAL #
#########

# Скачивание образа
FROM alpine:3

# Создание директории для пользователя app
RUN mkdir -p /home/app

# Создание пользователя app
RUN addgroup -S app && adduser -S app -G app

# Создание каталогов
ENV HOME=/home/app
ENV APP_HOME=/home/app/web
RUN mkdir $APP_HOME
WORKDIR $APP_HOME

# Копирование entrypoint.sh
COPY docker/entrypoint.sh $APP_HOME
RUN sed -i 's/\r$//g' $APP_HOME/entrypoint.sh && chmod +x $APP_HOME/entrypoint.sh

# Установка приложения и передача владения файлами пользовалелю app
COPY --from=builder /usr/local/bin/finapp $APP_HOME
RUN chmod +x $APP_HOME/finapp && chown -R app:app $APP_HOME

# Смена пользователя
USER app

# Запуск
ENTRYPOINT ["/home/app/web/entrypoint.sh"]
