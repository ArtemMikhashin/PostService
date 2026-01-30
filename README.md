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

(Опционально) Запуск миграций к PostgreSQL
```
make migrate.up
```

(Опционально) Команды к PostgreSQL
```
make db.exec
```

## Utit тесты

Запуск unit тестов
```
make test
```

# Примеры запросов

Создание поста
```
mutation {
  createPost(input: {
    author: "Test author"
    title: "Test title"
    content: "Test content"
    commentsAllowed: true
  }) {
    id
    title
  }
}
```
___
Получение поста с комментариями по id поста
```
query {
  post(id: 1) {
    id
    title
    comments(page: 1, pageSize: 3) {
      id
      content
      author
      replies {
        id
        content
        author
      }
    }
  }
}
```
___
Создание комментария
```
mutation {
  createComment(input: {
    author: "Test author"
    content: "Test content
    post: 1
    replyTo: 2
  }) {
    id
    content
  }
}
```
___
Подписка на комментарии к посту
```
subscription {
  commentAdded(postId: 1) {
    id
    author
    content
    createdAt
  }
}
```
