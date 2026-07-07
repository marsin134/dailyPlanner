# Daily Planner API

REST API для ежедневника с поддержкой регистрации пользователей, JWT-аутентификации, хранения пользовательских сессий и управления событиями. Проект реализован на **Go** с использованием **PostgreSQL**.

---

## Возможности

### Авторизация
- Регистрация пользователей
- Авторизация по email и паролю
- JWT Access Token
- Refresh Token
- Хранение пользовательских сессий
- Обновление Access Token
- Выход из текущей сессии
- Выход со всех устройств
- Проверка User-Agent и IP

### Пользователи
- Получение собственного профиля
- Изменение имени
- Изменение пароля
- Удаление аккаунта
- Получение пользователя по ID (Admin)
- Назначение администратора

### События
- Создание события
- Получение события по ID
- Получение всех событий пользователя за выбранную дату
- Обновление события
- Отметка события выполненным
- Удаление события

---

# Используемые технологии

- Go
- PostgreSQL
- sqlx
- JWT (golang-jwt/jwt)
- bcrypt
- UUID
- Middleware
- REST API

---

Архитектура разделена на несколько слоев.

## Handler

Обрабатывает HTTP-запросы.

Отвечает за:

- чтение JSON
- валидацию
- получение данных из middleware
- формирование ответа

---

## Service

Бизнес-логика приложения.

Отвечает за:

- регистрацию
- авторизацию
- создание JWT
- Refresh Token
- управление пользовательскими сессиями

---

## Repository

Работа с PostgreSQL.

Использует sqlx.

Содержит CRUD-операции для:

- users
- events
- user_sessions

---

## Middleware

Реализованы:

- JWT Authentication
- CORS
- Logging

После успешной проверки JWT SessionID помещается в Context запроса.

---

# База данных

## users

| Поле | Тип |
|------|-----|
| user_id | UUID |
| user_name | varchar |
| email | varchar |
| password_hash | text |
| role | User / Admin |

---

## user_sessions

Хранит активные Refresh Token пользователей.

| Поле | Тип |
|------|-----|
| session_id | UUID |
| user_id | UUID |
| refresh_token_hash | text |
| expires_at | timestamp |
| user_agent | text |
| ip_address | varchar |
| is_active | bool |

---

## events

| Поле | Тип |
|------|-----|
| event_id | UUID |
| user_id | UUID |
| title_event | varchar |
| date_event | date |
| completed | bool |
| color | red / green / blue |

---

# Авторизация

Используются два токена.

## Access Token

Содержит:

- Session ID
- время истечения

Используется для доступа к защищённым маршрутам.


---

## Refresh Token

Хранится только в базе данных в виде bcrypt-хеша.

Используется для получения нового Access Token.

---

# Пользовательские сессии

При каждом входе:

- создаётся новая сессия;
- сохраняется User-Agent;
- сохраняется IP;
- создаётся Refresh Token;
- Refresh Token хешируется bcrypt.

При повторном входе с того же устройства старая сессия деактивируется.

---

# API

## Авторизация

### POST

```Authorization: Bearer <access_token>```


---

## Refresh Token

Хранится только в базе данных в виде bcrypt-хеша.

Используется для получения нового Access Token.

---

# Пользовательские сессии

При каждом входе:

- создаётся новая сессия;
- сохраняется User-Agent;
- сохраняется IP;
- создаётся Refresh Token;
- Refresh Token хешируется bcrypt.

При повторном входе с того же устройства старая сессия деактивируется.

---

# API

## Авторизация

### POST
```/api/auth/register```

---

# Безопасность

Реализовано:

- bcrypt для хранения паролей
- bcrypt для Refresh Token
- JWT Authentication
- Middleware авторизации
- Проверка срока действия токена
- Проверка алгоритма подписи JWT
- Проверка User-Agent
- Проверка IP
- UUID в качестве идентификаторов
- Role-based доступ (User/Admin)

---

# Индексы PostgreSQL

Созданы индексы:

- users.email
- sessions.user_id
- sessions.refresh_token_hash
- sessions.expires_at
- events.user_id
- events.date_event
- events.completed

Это значительно ускоряет поиск пользователей, сессий и событий.

---

# Особенности проекта

- Чистая многослойная архитектура
- Repository Pattern
- Service Layer
- Middleware
- JWT Authentication
- Refresh Token Rotation
- bcrypt Hashing
- UUID
- PostgreSQL
- REST API
- Разделение ответственности между слоями

