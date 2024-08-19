
# House Service

House Service — это микросервис на языке Go для управления домами, квартирами и пользователями. Проект построен с использованием архитектурных принципов Clean Architecture и контейнизирован с помощью Docker.

## Оглавление
- [Описание](#описание)
- [Проделанная работа](#проделанная-работа)
- [Структура проекта](#структура-проекта)
- [Технологии](#технологии)
- [Установка и запуск](#установка-и-запуск)
  - [Настройка файла `.env`](#настройка-файла-env)
  - [Запуск без Docker](#запуск-без-docker)
  - [Запуск с Docker](#запуск-с-docker)
  - [Выполнение миграций вручную](#выполнение-миграций-вручную)
- [API Эндпоинты](#api-эндпоинты)
- [Автор](#автор)

## Описание

House Service — это микросервис, предоставляющий возможности для управления домами, квартирами и пользователями. В системе поддерживаются разные роли пользователей (модераторы и клиенты), а также функционал подписки на дома.

## Проделанная работа

В рамках проекта была реализована следующая функциональность:
- Создание и управление домами и квартирами.
- Регистрация и аутентификация пользователей.
- Ролевое управление доступом (модераторы и клиенты).
- Подписка на обновления по домам.
- Имплементация чистой архитектуры (Clean Architecture).
- Контейнеризация приложения с использованием Docker.

## Технологии

- Go
- Docker
- PostgreSQL
- Chi (HTTP роутер для Go)
- Clean Architecture
- Pgxpool (для работы с PostgreSQL)

## Установка и запуск

### Настройка файла `.env`

Перед запуском проекта создайте файл `.env` в корневой директории и заполните его следующими переменными окружения:

```dotenv
POSTGRES_USER=your_postgres_user
POSTGRES_PASSWORD=your_postgres_password
POSTGRES_DB=your_database_name
SECRET_KEY=your_secret_key
```

### Запуск без Docker

1. Убедитесь, что у вас установлены Go и PostgreSQL.
2. Настройте подключение к базе данных в файле `config/config.yaml`.
3. Выполните миграции базы данных вручную. Например:

```bash
psql -U your_postgres_user -d your_database_name -f infrastructure/db/migrations/0001_init_schema.up.sql
```

4. Соберите и запустите приложение:

```bash
go build -o house-service cmd/main.go
./house-service
```

### Запуск с Docker

1. Убедитесь, что у вас установлен Docker.
2. Соберите Docker-образ:

```bash
docker build -t house-service .
```

3. Запустите контейнер:

```bash
docker run -d -p 8080:8081 --name house-service --env-file .env house-service
```

4. Приложение будет доступно по адресу [http://localhost:8080](http://localhost:8080).

### Выполнение миграций вручную

Если вы не используете Docker Compose для управления миграциями, выполните их вручную с помощью `psql` или другого клиента PostgreSQL:

```bash
psql -U your_postgres_user -d your_database_name -f infrastructure/db/migrations/0001_init_schema.up.sql
```

## API Эндпоинты

Примеры ключевых API эндпоинтов:

- `POST /register` — Регистрация нового пользователя.
- `POST /login` — Аутентификация пользователя.
- `POST /house/create` — Создание нового дома (доступно только модератору).
- `GET /house/{id}` — Получение информации о доме и списке квартир.
- `POST /flat/create` — Создание новой квартиры.
