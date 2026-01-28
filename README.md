# PostService

# Запуск сервера
Если в .env IN_MEMORY=false, то сначала запустить postgres (ниже)
```
make app.up
```

## Запуск postgreSQL образа с миграциями
```
make db.up
```

## Запуск postgreSQL образ
```
make migrate.up
```

## Команды к бд
```
make db.exec
```

## Отключить postgreSQL образ
```
make db.down
```

