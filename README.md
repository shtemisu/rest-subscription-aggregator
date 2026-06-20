# Subscription Aggregator

REST-сервис для агрегации данных об онлайн-подписках пользователей.
Предоставляет CRUDL-операции над записями о подписках и подсчёт суммарной
стоимости подписок за выбранный период с фильтрацией.

## Стек

- **Go** + `net/http` (паттерн-роутинг 1.22+)
- **PostgreSQL** через `pgx/v5` (пул соединений)
- **golang-migrate** — миграции схемы
- **uber/fx** — внедрение зависимостей и управление жизненным циклом
- **slog** — структурированное логирование
- **swaggo** — Swagger/OpenAPI документация
- **Docker Compose** — запуск

## Архитектура

Слоистая, зависимости направлены к домену через интерфейсы:

```
controller  → HTTP, парсинг/валидация, коды ответов   (internal/controller)
usecase     → бизнес-логика                            (internal/usecase)
repository  → SQL через pgx                            (internal/repository)
domain      → модели и интерфейсы (порты)              (internal/domain)
```

## Запуск

### Через Docker Compose (рекомендуется)

```bash
cp .env.example .env          # при необходимости поправь значения
docker compose up -d --build
```

Поднимется PostgreSQL и приложение. Миграции накатываются автоматически при старте.

- API:     http://localhost:8080
- Swagger: http://localhost:8080/swagger/index.html

Остановить:

```bash
docker compose down       # остановить
docker compose down -v    # + удалить данные БД (volume)
```
## Конфигурация

Переменные окружения (переопределяют значения; в Docker задаются через `.env`):

| Переменная     | Описание                | Пример                   |
|----------------|-------------------------|--------------------------|
| `DB_HOST`      | хост БД                 | `postgres-sub-aggregator`|
| `DB_PORT`      | порт БД                 | `5432`                   |
| `DB_USER`      | пользователь            | `postgres`               |
| `DB_PASSWORD`  | пароль                  | `1234`                   |
| `DB_NAME`      | имя базы                | `subscriptions`          |

> В Docker Compose `DB_HOST` должен совпадать с именем сервиса Postgres
> (`postgres-sub-aggregator`), т.к. контейнеры резолвятся по имени сервиса.

## API

Базовый URL: `http://localhost:8080`. Даты — в формате `MM-YYYY` (например `07-2025`).

| Метод    | Путь                       | Описание                          |
|----------|----------------------------|-----------------------------------|
| `POST`   | `/subscriptions`           | создать подписку                  |
| `GET`    | `/subscriptions/{id}`      | получить по ID                    |
| `PUT`    | `/subscriptions/{id}`      | обновить                          |
| `DELETE` | `/subscriptions/{id}`      | удалить                           |
| `GET`    | `/subscriptions`           | список (фильтры + пагинация)      |
| `GET`    | `/subscriptions/summary`   | сумма за период                   |
| `GET`    | `/swagger/`                | Swagger UI                        |

### Поля записи

| Поле           | Тип            | Обяз. | Описание                          |
|----------------|----------------|-------|-----------------------------------|
| `service_name` | string         | да    | название сервиса                  |
| `price`        | int (рубли)    | да    | стоимость месяца, ≥ 0             |
| `user_id`      | string (UUID)  | да    | идентификатор пользователя        |
| `start_date`   | string `MM-YYYY` | да  | начало подписки                   |
| `end_date`     | string `MM-YYYY` | нет | окончание подписки (опционально)  |

### Параметры запроса для `GET /subscriptions` и `/subscriptions/summary`

| Параметр       | Описание                                   |
|----------------|--------------------------------------------|
| `user_id`      | фильтр по пользователю (UUID)              |
| `service_name` | фильтр по названию сервиса                 |
| `from`         | начало периода `MM-YYYY` (для summary)     |
| `to`           | конец периода `MM-YYYY` (для summary)      |
| `limit`        | размер страницы (по умолчанию 50, макс 100; только для списка) |
| `offset`       | смещение (только для списка)               |

### Примеры

Создать подписку:

```bash
curl -X POST http://localhost:8080/subscriptions \
  -H 'Content-Type: application/json' \
  -d '{
    "service_name": "Yandex Plus",
    "price": 400,
    "user_id": "60601fee-2bf1-4721-ae6f-7636e79a0cba",
    "start_date": "07-2025"
  }'
# 201 -> {"id":"c36082f7-c756-429e-9875-decf0ae12f96"}
```

Список подписок пользователя:

```bash
curl "http://localhost:8080/subscriptions?user_id=60601fee-2bf1-4721-ae6f-7636e79a0cba"
```

Сумма за период:

```bash
curl "http://localhost:8080/subscriptions/summary?user_id=60601fee-2bf1-4721-ae6f-7636e79a0cba&from=01-2025&to=12-2025"
# 200 -> {"total":400}
```

## Миграции

Лежат в `migrations/`, накатываются автоматически при старте приложения.
Применить вручную:

```bash
migrate -database 'postgres://postgres:1234@localhost:5432/subscriptions?sslmode=disable' \
  -path migrations up
```

## Swagger

Документация генерируется из аннотаций в хендлерах. После их изменения
перегенерировать:

```bash
swag init -g cmd/main.go -o docs --parseDependency --parseInternal
```

Папка `docs/` коммитится в репозиторий — она нужна для сборки
(подключается blank-импортом в роутере).

## Структура проекта

```
cmd/                  точка входа (fx)
config/               загрузка конфигурации из env
internal/
  controller/         HTTP-хендлеры, роутер, middleware, DTO
  usecase/            бизнес-логика
  repository/         доступ к PostgreSQL
  domain/             модели и интерфейсы
  di/                 сборка графа зависимостей (fx)
pkg/
  logger/             настройка slog
  migrator/           прогон миграций
  response/           JSON-хелперы ответов
migrations/           SQL-миграции
docs/                 сгенерированная Swagger-спецификация
build/app/Dockerfile  multi-stage сборка
docker-compose.yml
```
