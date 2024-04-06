# Бэкенд ЛФП

[wiki](https://wiki.zlp.ooo/ru/lfp)

## Запуск проекта

> Development

- Docker compose

```bash
docker compose up
```

- Go run

```bash
cd app
swag init -g server.go
go run ./server.go app:serve
```

- Go Task

```bash
task app:serve
```
