# 📦 WareHouse — Система управления складом

REST API для управления складом на Go с веб-интерфейсом. Clean Architecture, JWT-авторизация, PostgreSQL.

## Стек

- **Go 1.22** + **Fiber v2** — веб-фреймворк
- **GORM** — ORM
- **PostgreSQL** — база данных
- **JWT** — авторизация
- **Docker** — запуск БД
- **Vanilla HTML/JS** — фронтенд

## Структура проекта

```
warehouse-app/
├── adapters/           # HTTP-хендлеры и GORM-репозитории (продукты, категории, поставщики)
├── auth/
│   ├── authAdapter/    # HTTP-хендлеры и GORM-репозиторий пользователей
│   └── authCore/       # Интерфейсы и сервис авторизации
├── core/               # Бизнес-логика: сервис и репозиторий продуктов
├── database/           # Модели (Product, Supplier, Category, User)
├── frontend/           # Веб-интерфейс (index.html)
├── docker-compose.yml
├── .env.example
└── main.go
```

## Запуск

### 1. Клонировать репозиторий

```bash
git clone https://github.com/ryzzz0n/ware-house.git
cd ware-house
```

### 2. Настроить окружение

```bash
cp .env.example .env
```

Заполнить `.env`:

```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=yourpassword
DB_NAME=warehouse
JWT_SECRETKEY=your-secret-key-minimum-32-chars
```

### 3. Запустить PostgreSQL

```bash
docker-compose up -d
```

### 4. Запустить сервер

```bash
go mod tidy
go run main.go
```

Сервер запустится на `http://localhost:8000`

### 5. Открыть фронтенд

Открыть `frontend/index.html` в браузере.

## API

| Метод  | Путь                      | Описание                        | Auth |
|--------|---------------------------|---------------------------------|------|
| POST   | `/register`               | Регистрация пользователя        | —    |
| POST   | `/login`                  | Вход, возвращает JWT            | —    |
| POST   | `/supplier`               | Создать поставщика              | ✓    |
| GET    | `/supplier`               | Список поставщиков              | ✓    |
| PUT    | `/supplier/:id`           | Обновить поставщика             | ✓    |
| DELETE | `/supplier/:id`           | Удалить поставщика              | ✓    |
| POST   | `/category`               | Создать категорию               | ✓    |
| GET    | `/category`               | Список категорий                | ✓    |
| DELETE | `/category/:id`           | Удалить категорию               | ✓    |
| POST   | `/product`                | Создать товар                   | ✓    |
| GET    | `/product`                | Список товаров                  | ✓    |
| GET    | `/product/:id`            | Товар по ID                     | ✓    |
| PUT    | `/product/:id`            | Обновить товар                  | ✓    |
| DELETE | `/product/:id`            | Удалить товар                   | ✓    |
| GET    | `/category/:id/product`   | Товары по категории             | ✓    |
| GET    | `/supplier/:id/product`   | Товары по поставщику            | ✓    |

Защищённые маршруты принимают токен двумя способами:

```
Authorization: Bearer <token>
```
или cookie `jwt`.

## Авторизация

```bash
# Регистрация
curl -X POST http://localhost:8000/register \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@test.com","password":"123456"}'

# Вход
curl -X POST http://localhost:8000/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@test.com","password":"123456"}'
```

В ответе придёт `token` — передавать в заголовке `Authorization: Bearer <token>`.

## Автор

Максим Булюкин — учебный проект по Go
