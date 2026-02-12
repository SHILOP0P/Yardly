# Yardly

Yardly — это pet-проект для аренды/обмена вещами с бронированиями, избранным и админ-модерацией.

## Структура проекта

- `backend/` - Go API (HTTP), авторизация, бизнес-логика, работа с Postgres.
- `frontend/` - Next.js UI, который отображает и проверяет сценарии бэкенда.
- `backend/db/init/` - SQL-миграции/инициализация БД.
- `docker-compose.yml` - локальный Postgres.
- `commands.ps1` - вспомогательные команды для повседневной разработки.

## Технологии

- Backend: Go 1.25, `net/http`, `pgx`, JWT, HttpOnly refresh tokens.
- Frontend: Next.js 16, React 19, TypeScript, React Query, Zustand.
- База данных: PostgreSQL 16 (Docker).

## Быстрый старт (Windows PowerShell)

1. Запустить все для разработки:

```powershell
.\commands.ps1 dev
```

Команда поднимает БД и открывает 2 терминала (backend + frontend).

2. Или запускать по отдельности:

```powershell
.\commands.ps1 db-up
.\commands.ps1 backend
.\commands.ps1 frontend
```

3. Открыть:

- Frontend: `http://localhost:3000`
- Проверка backend: `http://localhost:8080/health`

## Вспомогательные команды

```powershell
.\commands.ps1 help
.\commands.ps1 db-up
.\commands.ps1 db-down
.\commands.ps1 db-shell
.\commands.ps1 backend
.\commands.ps1 frontend
.\commands.ps1 dev
.\commands.ps1 migration:new add_some_feature
.\commands.ps1 migration:apply 023_add_some_feature.sql
```

## Переменные окружения

- Backend: `backend/.env`
- Frontend: `frontend/.env.local`

В репозитории уже есть рабочие локальные значения для разработки.

## Документация по частям

- Backend: `backend/README.md`
- Frontend: `frontend/README.md`
