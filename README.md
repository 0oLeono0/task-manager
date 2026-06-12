# Task Manager

fullstack task manager с backend на Go,
хранением данных в PostgreSQL и frontend на Vue.

## Стек

- Backend: Go, chi, PostgreSQL
- База данных и dev-инфраструктура: Docker Compose
- Frontend: Vue 3, TypeScript, Vite, Pinia, Axios

## Архитектура

- Frontend отправляет запросы на `/api/*`.
- Vite proxy проксирует эти запросы на Go backend.
- Backend обрабатывает REST endpoints.
- Данные хранятся в PostgreSQL.

## Текущие возможности

- CRUD API для задач: список, детали, создание, обновление, удаление
- Хранение задач в PostgreSQL
- Frontend-список задач с созданием, переключением completed и удалением
- UI-состояния загрузки, ошибки и пустого списка

## Скриншот

![Главный экран task manager](docs/screenshots/main.png)

## Переменные окружения

Перед запуском базы данных и backend создай локальный `.env` из примера:

```powershell
Copy-Item .env.example .env
```

`.env.example` хранится в репозитории. `.env` локальный и должен содержать реальные
значения.

Важные переменные:

- `DATABASE_DSN`: строка подключения backend к базе данных
- `TEST_DATABASE_DSN`: строка подключения для backend integration tests
- `SERVER_ADDR`: адрес backend-сервера, по умолчанию `:8080`

Пример формата DSN:

```text
postgres://task_manager_user:task_manager_password@localhost:5432/task_manager?sslmode=disable
```

Для `TEST_DATABASE_DSN` лучше использовать отдельную тестовую базу, чтобы тесты
случайно не трогали development database.

## Запуск

Запустить PostgreSQL:

```powershell
docker compose up -d postgres
```

Применить миграции:

```powershell
migrate -path backend/migrations -database $env:DATABASE_DSN up
```

Для этой команды нужен установленный `golang-migrate` CLI. Перед запуском должен
быть задан `DATABASE_DSN`.

Запустить backend:

```powershell
cd backend
go run ./cmd/api
```

Запустить frontend:

```powershell
cd frontend
npm install
npm run dev
```

В dev-режиме Vite проксирует запросы `/api/*` на Go backend.

## Проверки

Backend tests:

```powershell
cd backend
go test ./...
```

PostgreSQL integration tests используют `TEST_DATABASE_DSN`. Если переменная не
задана, эти тесты будут пропущены.

Frontend lint:

```powershell
cd frontend
npm run lint
```

Проверка форматирования frontend:

```powershell
cd frontend
npm run format:check
```

Production build frontend:

```powershell
cd frontend
npm run build
```

## Проверка качества

Перед финальной ручной проверкой должны проходить автоматические проверки:

- `go test ./...` в `backend`
- `npm run lint` в `frontend`
- `npm run format:check` в `frontend`
- `npm run build` в `frontend`

Ручные сценарии:

- Загрузка задач.
- Создание задачи, переключение `completed`, удаление задачи.
- Empty state: `Задач пока нет`.
- Ошибка загрузки при выключенном backend.
