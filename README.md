# Account Master 🔐

**Простой сервис для управления пользовательскими аккаунтами** с возможностью регистрации, аутентификации и управления профилями.

---

## 🚀 Быстрый старт

### Запуск проекта
```bash
make all
```

## 📚 Документация API

Документация в формате Swagger доступна после запуска сервиса:

[👉 **Swagger UI**](http://localhost:8080/swagger/index.html)  
`http://localhost:8080/swagger/index.html`

### Доступные эндпоинты
```yaml
GET    /api/v1/users      - Список пользователей
POST   /api/v1/users      - Создание пользователя
GET    /api/v1/users/{id} - Получить пользователя
PUT    /api/v1/users/{id} - Обновить пользователя
DELETE /api/v1/users/{id} - Удалить пользователя



### Структура проекта
.
├── config.yaml
├── go.mod
├── go.sum
├── Makefile
├── README.md
└── src
    ├── config
    │   └── config.go
    ├── docs
    │   ├── docs.go
    │   ├── swagger.json
    │   └── swagger.yaml
    ├── internal
    │   ├── app
    │   │   └── app.go
    │   ├── controllers
    │   │   ├── api.go
    │   │   ├── controllers.go
    │   │   └── middleware.go
    │   ├── hash
    │   │   └── hash.go
    │   └── model
    │       └── model.go
    ├── main.go
    └── pkg
        ├── server
        │   ├── option.go
        │   └── server.go
        └── storage
            └── mock
                ├── mock.go
                └── mock_test.go

