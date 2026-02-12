# Yardly Backend

Бэкенд-сервис предоставляет API для авторизации, пользователей, вещей, бронирований, избранного и админ-функций.

## Стек

- Go 1.25+
- `net/http` (`http.ServeMux`)
- PostgreSQL через `pgx/v5`
- JWT access token + refresh token

## Основные возможности

- Регистрация и логин.
- Обновление/ротация токенов, logout/logout-all.
- Получение профиля текущего пользователя (`/api/users/me`).
- Операции с вещами для владельца (создание, списки, детали, мои вещи).
- Управление изображениями вещей.
- Жизненный цикл бронирования (создать, одобрить, передать, вернуть, отменить).
- Проверка доступности и ближайших бронирований.
- Избранное для авторизованных пользователей.
- Админ-эндпоинты для пользователей/бронирований/вещей/событий.
- Проверка ролей (`user`, `admin`, `superadmin`) и блокировки пользователя.

## Локальный запуск

1. Поднять БД:

```powershell
cd ..
.\commands.ps1 db-up
```

2. Запустить backend:

```powershell
cd backend
go run ./cmd/api
```

или из корня репозитория:

```powershell
.\commands.ps1 backend
```

Бэкенд по умолчанию: `http://localhost:8080`

## Переменные окружения

Пример значений находится в `backend/.env`:

- `APP_PORT` (по умолчанию в репо: `8080`)
- `DATABASE_URL` (локальный Postgres на `127.0.0.1:55432`)
- `JWT_SECRET`
- `JWT_TTL_MINUTES` (формат duration, пример `60m`)
- `REFRESH_TTL` (формат duration, пример `720h`)

## Текущие API маршруты

### Базовые

- `GET /`
- `GET /health`

### Авторизация

- `POST /api/auth/register`
- `POST /api/auth/login`
- `POST /api/auth/refresh`
- `POST /api/auth/logout`
- `POST /api/auth/logout_all`

### Пользователи

- `GET /api/users/me`

### Вещи

- `POST /api/items`
- `GET /api/items`
- `GET /api/items/{id}`
- `GET /api/my/items`
- `GET /api/users/{id}/items`
- `GET /api/items/{id}/images`
- `POST /api/items/{id}/images`
- `DELETE /api/items/{id}/images/{imageId}`

### Бронирования

- `POST /api/items/{id}/bookings`
- `GET /api/items/{id}/bookings`
- `GET /api/my/bookings`
- `GET /api/my/items/bookings`
- `GET /api/my/items/booking-requests`
- `GET /api/items/{id}/bookings/upcoming`
- `GET /api/items/{id}/availability`
- `POST /api/bookings/{id}/approve`
- `POST /api/bookings/{id}/handover`
- `POST /api/bookings/{id}/return`
- `POST /api/bookings/{id}/cancel`
- `GET /api/bookings/{id}/events`

### Избранное

- `POST /api/items/{id}/favorite`
- `DELETE /api/items/{id}/favorite`
- `GET /api/my/favorites`
- `GET /api/items/{id}/favorite`

### Админка

- `GET /api/admin/users`
- `GET /api/admin/users/{id}`
- `PATCH /api/admin/users/{id}`
- `GET /api/admin/bookings`
- `GET /api/admin/bookings/{id}`
- `GET /api/admin/bookings/{id}/events`
- `GET /api/admin/items`
- `PATCH /api/admin/items/{id}`
- `POST /api/admin/items/{id}/block`
- `POST /api/admin/items/{id}/unblock`
- `POST /api/admin/items/{id}/delete`
- `GET /api/admin/events`

## База данных и миграции

SQL-файлы схемы/миграций лежат в `backend/db/init/` и монтируются в контейнер Postgres.

Полезные команды из корня репозитория:

```powershell
.\commands.ps1 migration:new add_index_for_bookings
.\commands.ps1 migration:apply 023_add_index_for_bookings.sql
```

Примечания:

- `migration:new` создает следующий по номеру SQL-файл в `backend/db/init`.
- `migration:apply` применяет один SQL-файл в запущенный контейнер `yardly-db` через `psql`.
- На чистом volume Postgres автоматически выполнит init-файлы из `backend/db/init`.

## CORS

Сейчас разрешен origin фронтенда: `http://localhost:3000`.
