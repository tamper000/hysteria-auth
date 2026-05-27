# Hysteria Auth

Простой сервер аутентификации для [Hysteria 2](https://github.com/apernet/hysteria).  
Позволяет регистрировать, удалять и проверять пользователей через HTTP API.

## Возможности

- Регистрация пользователей с ID, паролем и дополнительными данными
- Удаление пользователей по ID
- Эндпоинт аутентификации для интеграции с Hysteria 2
- Хранение данных в SQLite (в памяти или в файле)
- Простота — никаких сложных зависимостей, один бинарник

## API

### Регистрация пользователя

`POST /register`

```json
{
  "id": "user123",
  "auth": "someSecurePassword",
  "optional": "telegramID:12345"
}
```

Ответы:
- `201` — пользователь создан
- `409` — пользователь с таким `auth` уже существует

### Удаление пользователя

`POST /delete`

```json
{
  "id": "user123"
}
```

Ответы:
- `200` — пользователь удалён
- `404` — пользователь не найден

### Аутентификация (для Hysteria 2)

`POST /auth`

```json
{
  "auth": "someSecurePassword"
}
```

Ответы:
- `200` — успешная аутентификация, возвращает `{"ok": true, "id": "user123"}`
- `404` — пользователь не найден

## Конфигурация

Настройка через переменные окружения (поддерживается `.env` файл):

| Переменная  | По умолчанию  | Описание                          |
|-------------|---------------|-----------------------------------|
| `PORT`      | `8888`        | Порт HTTP сервера                 |
| `DB_PATH`   | `:memory:`    | Путь к файлу SQLite (или в памяти) |

Если `DB_PATH` не указан, данные хранятся в оперативной памяти и теряются при перезапуске.  
Для постоянного хранения укажите путь к файлу, например:

```env
DB_PATH=/data/hysteria_auth.db
```

## Быстрый старт

### Запуск через Go

```bash
git clone https://github.com/tamper000/hysteria-auth.git
cd hysteria-auth

# Запуск с временем хранения (всё в памяти)
go run ./cmd/

# Или с файлом базы данных
PORT=8888 DB_PATH=./users.db go run ./cmd/
```

### Сборка бинарника

```bash
go build -o hysteria-auth ./cmd/
./hysteria-auth
```

### Примеры использования

```bash
# Регистрация
curl -X POST http://localhost:8888/register \
  -H "Content-Type: application/json" \
  -d '{"id":"user123","auth":"mypassword","optional":"tg:12345"}'

# Аутентификация (то, что использует Hysteria 2)
curl -X POST http://localhost:8888/auth \
  -H "Content-Type: application/json" \
  -d '{"auth":"mypassword"}'

# Удаление
curl -X POST http://localhost:8888/delete \
  -H "Content-Type: application/json" \
  -d '{"id":"user123"}'
```

## Интеграция с Hysteria 2

В конфигурации Hysteria 2 укажите URL вашего auth-сервера:

```yaml
acme:
  domains:
    - your.domain.com

auth:
  type: http
  http:
    url: http://your-server:8888/auth
```

Hysteria будет отправлять POST-запросы на `/auth` с полем `auth`, а сервер отвечать `{"ok": true, "id": "..."}` при успехе.

## Технологии

- **Go 1.26+**
- [Huma](https://huma.rocks/) — OpenAPI фреймворк
- [Chi](https://github.com/go-chi/chi) — HTTP роутер
- [SQLite](https://gitlab.com/cznic/sqlite) — in-process база данных (pure Go, без CGO)
