# PostService

Сервис по работе с постами и комментариями к ним

# Запуск сервиса

! Создать .env в корне по примеру .env.example

## Docker
Запуск
```
docker-compose up --build
```
Отключение
```
docker-compose down
```

## Локально

1) Запуск postgreSQL образа с миграциями
```
make db.up
```

2) Запуск приложения
```
make app.up
```

(Завершение) Отключить postgreSQL образ
```
make db.down
```

(Опционально) Запуск миграций к БД
```
make migrate.up
```

(Опционально) Команды к бд
```
make db.exec
```