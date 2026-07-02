# Subscription Service

Subscription Service — REST-сервис на Go для управления онлайн-подписками пользователей и подсчёта суммарной стоимости подписок за выбранный период.

## Соответствие требованиям

| Требование | Статус | Реализация |
|---|---|---|
| CRUDL для подписок | Выполнено | Реализованы ручки создания, получения, обновления, удаления и списка подписок |
| Подсчёт суммарной стоимости | Выполнено | Реализована ручка расчёта суммы по `user_id`, `service_name` и периоду |
| PostgreSQL | Выполнено | Используется PostgreSQL, схема создаётся через миграцию |
| Миграции | Выполнено | В проекте есть `migrations/001_init.sql` |
| Логи | Выполнено | Используется стандартный `slog`, логи пишутся в stdout |
| Конфигурация | Выполнено | Конфигурация вынесена в `.env` |
| Swagger | Выполнено | Swagger UI доступен по `/swagger/index.html` |
| Docker Compose | Выполнено | Сервис и PostgreSQL запускаются через `docker compose` |
| Postman collection | Выполнено | В корне проекта есть коллекция `Subscription Service API.postman_collection.json` |

## Стек

- Go
- Gin
- PostgreSQL
- sqlx
- lib/pq
- slog
- Swagger
- Docker
- Docker Compose

## Структура проекта

```text
.
├── Dockerfile
├── Subscription Service API.postman_collection.json
├── cmd
│   └── app
│       └── main.go
├── docker-compose.yml
├── docs
│   ├── docs.go
│   ├── swagger.json
│   └── swagger.yaml
├── internal
│   ├── config
│   ├── db
│   ├── handler
│   ├── model
│   ├── repository
│   └── service
└── migrations
    └── 001_init.sql
```

## Быстрый старт

Скачать проект:

```bash
git clone https://github.com/f7rzen/subscription-service.git
cd subscription-service
```

Создать файл `.env`:

```bash
cp .env.example .env
```

Пример содержимого `.env`:

```env
APP_PORT=8080

DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=subscriptions_db
DB_SSLMODE=disable
```

Запустить сервисы:

```bash
docker compose up --build -d
```

После запуска API будет доступно по адресу:

```text
http://localhost:8080
```

Проверить работу сервиса:

```bash
curl http://localhost:8080/health
```

Ожидаемый ответ:

```json
{
  "status": "ok"
}
```

Остановить сервисы:

```bash
docker compose down
```

Остановить сервисы и удалить данные PostgreSQL:

```bash
docker compose down -v
```

## Swagger

Swagger UI доступен по адресу:

```text
http://localhost:8080/swagger/index.html
```
## Изменение порта

По умолчанию приложение запускается на порту `8080`.

Чтобы изменить порт, нужно поменять значение `PORT` в `.env`:

```env
PORT=9090
```

После этого перезапустить сервис:

```bash
docker compose down
docker compose up --build -d
```

API будет доступно по адресу:

```text
http://localhost:9090
```

## API

### Проверка состояния сервиса

```http
GET /health
```

Пример curl:

```bash
curl http://localhost:8080/health
```

Пример ответа:

```json
{
  "status": "ok"
}
```

---

### Создание подписки

```http
POST /api/v1/subscriptions
```

Пример curl:

```bash
curl -X POST http://localhost:8080/api/v1/subscriptions \
  -H "Content-Type: application/json" \
  -d '{
    "service_name": "Yandex Plus",
    "price": 400,
    "user_id": "60601fee-2bf1-4721-ae6f-7636e79a0cba",
    "start_date": "07-2025"
  }'
```

Пример ответа:

```json
{
  "id": 1,
  "service_name": "Yandex Plus",
  "price": 400,
  "user_id": "60601fee-2bf1-4721-ae6f-7636e79a0cba",
  "start_date": "07-2025",
  "end_date": null,
  "created_at": "2026-07-02T17:53:49.361724Z",
  "updated_at": "2026-07-02T17:53:49.361724Z"
}
```

---

### Получение списка подписок

```http
GET /api/v1/subscriptions
```

Пример curl:

```bash
curl http://localhost:8080/api/v1/subscriptions
```

Пример ответа:

```json
[
  {
    "id": 1,
    "service_name": "Yandex Plus",
    "price": 400,
    "user_id": "60601fee-2bf1-4721-ae6f-7636e79a0cba",
    "start_date": "07-2025",
    "end_date": null,
    "created_at": "2026-07-02T17:53:49.361724Z",
    "updated_at": "2026-07-02T17:53:49.361724Z"
  }
]
```

---

### Получение подписки по ID

```http
GET /api/v1/subscriptions/{id}
```

Пример curl:

```bash
curl http://localhost:8080/api/v1/subscriptions/1
```

Пример ответа:

```json
{
  "id": 1,
  "service_name": "Yandex Plus",
  "price": 400,
  "user_id": "60601fee-2bf1-4721-ae6f-7636e79a0cba",
  "start_date": "07-2025",
  "end_date": null,
  "created_at": "2026-07-02T17:53:49.361724Z",
  "updated_at": "2026-07-02T17:53:49.361724Z"
}
```

---

### Обновление подписки

```http
PUT /api/v1/subscriptions/{id}
```

Пример curl:

```bash
curl -X PUT http://localhost:8080/api/v1/subscriptions/1 \
  -H "Content-Type: application/json" \
  -d '{
    "service_name": "Yandex Plus",
    "price": 500,
    "user_id": "60601fee-2bf1-4721-ae6f-7636e79a0cba",
    "start_date": "07-2025",
    "end_date": "12-2025"
  }'
```

Пример ответа:

```json
{
  "id": 1,
  "service_name": "Yandex Plus",
  "price": 500,
  "user_id": "60601fee-2bf1-4721-ae6f-7636e79a0cba",
  "start_date": "07-2025",
  "end_date": "12-2025",
  "created_at": "2026-07-02T17:53:49.361724Z",
  "updated_at": "2026-07-02T17:57:46.044312Z"
}
```

---

### Удаление подписки

```http
DELETE /api/v1/subscriptions/{id}
```

Пример curl:

```bash
curl -X DELETE http://localhost:8080/api/v1/subscriptions/1
```

Пример ответа:

```json
{
  "message": "subscription deleted"
}
```

---

### Подсчёт суммарной стоимости

```http
GET /api/v1/subscriptions-summary
```

Параметры:

| Параметр | Описание |
|---|---|
| `user_id` | ID пользователя в формате UUID |
| `service_name` | Название сервиса |
| `from` | Начало периода в формате `MM-YYYY` |
| `to` | Конец периода в формате `MM-YYYY` |

Пример curl:

```bash
curl "http://localhost:8080/api/v1/subscriptions-summary?user_id=60601fee-2bf1-4721-ae6f-7636e79a0cba&service_name=Yandex%20Plus&from=07-2025&to=12-2025"
```

Пример ответа:

```json
{
  "total_price": 3000
}
```

## Postman collection

В корне проекта находится файл:

```text
Subscription Service API.postman_collection.json
```

Чтобы использовать коллекцию:

1. Открыть Postman.
2. Нажать `Import`.
3. Выбрать файл `Subscription Service API.postman_collection.json`.
4. Запустить нужный запрос.