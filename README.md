# GophKeeper

GophKeeper - это клиент-серверная система для безопасного хранения и синхронизации приватных данных: паролей, текстовых заметок, бинарных файлов и данных банковских карт.

## Возможности

- 🔐 Безопасное хранение паролей и личных данных
- 📱 Синхронизация между несколькими устройствами
- 🖥️ CLI клиент для Windows, Linux и macOS
- 🔒 Шифрование данных на стороне клиента
- 📄 Поддержка различных типов данных:
  - Пары логин/пароль
  - Произвольные текстовые данные
  - Бинарные файлы
  - Данные банковских карт
  - Метаинформация для всех типов данных

## Структура проекта

```
  🔧 СЕРВИСЫ (бизнес-логика)

  internal/domain/services/
  ├── auth_service.go    # Аутентификация пользователей
  └── jwt.go            # Управление JWT токенами

  🛡️ ХЭНДЛЕРСЫ (обработчики запросов)

  internal/interfaces/grpc/handlers/
  └── auth_handler.go   # gRPC обработчики для аутентификации

  🔒 МИДЛВАРЕ (промежуточное ПО)

  internal/interfaces/grpc/middleware/
  └── middleware.go     # JWT аутентификация для gRPC

  🏗️ ИНФРАСТРУКТУРА

  internal/infrastructure/crypto/
  ├── aes.go           # AES шифрование
  ├── password.go      # Хеширование паролей
  ├── utils.go         # Криптографические утилиты
  └── errors.go        # Ошибки криптографии

  📦 СУЩНОСТИ

  internal/domain/entities/
  ├── user.go          # Пользователь
  ├── data_item.go     # Элемент данных
  ├── credit_card.go   # Кредитная карта
  └── ...              # Другие сущности
```

## Быстрый старт

### Установка зависимостей

```bash
make deps
```

### Установка инструментов разработки

```bash
make dev-deps
```

### Сборка

```bash
# Собрать все компоненты
make build

# Собрать только сервер
make build-server

# Собрать только клиент
make build-client

# Собрать для всех платформ
make build-all
```

### Запуск тестов

```bash
# Все тесты
make test

# Тесты с покрытием (генерирует coverage.html)
make test-cover

# Краткая сводка по покрытию
make test-cover-summary

# Тесты конкретных пакетов
go test ./internal/domain/entities -v -cover
go test ./internal/domain/services -v -cover
go test ./internal/infrastructure/config -v -cover
go test ./internal/infrastructure/crypto -v -cover
```

### Текущее покрытие тестами

Проект имеет **40.6%** общего покрытия кода тестами:

| Пакет | Покрытие |
|-------|----------|
| `internal/domain/entities` | 89.2% |
| `internal/infrastructure/crypto` | 87.0% |
| `internal/interfaces/grpc/middleware` | 77.8% |
| `internal/infrastructure/config` | 66.1% |
| `internal/domain/services` | 43.0% |
| `pkg/version` | 100% |

Основные пакеты бизнес-логики и критичных компонентов имеют высокое покрытие тестами.

### Проверка кода

```bash
# Линтер
make lint

# Go vet
make vet

# Форматирование
make fmt
```

## Использование

### Конфигурация

Сервер настраивается через переменные окружения. Скопируйте `.env.example` в `.env` и настройте:

```bash
cp .env.example .env
# Отредактируйте .env файл
```

**Важные переменные:**
- `JWT_SIGNING_KEY` - секретный ключ для подписи JWT токенов (обязательный!)
- `DATABASE_DSN` - строка подключения к PostgreSQL
- `SERVER_PORT` - порт для gRPC сервера (по умолчанию 8080)

**Генерация безопасного JWT ключа:**
```bash
openssl rand -base64 32
```

### Сервер

```bash
# Запуск с переменными окружения
export JWT_SIGNING_KEY="your-secret-key-here"
export DATABASE_DSN="host=localhost port=5432 user=postgres password=postgres dbname=gophkeeper sslmode=disable"
./build/server

# Или через Go (для разработки)
make run-server
```

### Клиент

```bash
# Запуск скомпилированного клиента
./build/client

# Или запуск через Go (для разработки)
make run-client

# Регистрация
./build/client register --username myuser --password mypass

# Аутентификация
./build/client login --username myuser --password mypass

# Добавление пароля
./build/client add password --name "Github" --login "myuser" --password "mypass" --url "github.com"

# Получение пароля
./build/client get password --name "Github"

# Синхронизация
./build/client sync
```

## Разработка

### Генерация protobuf

```bash
make proto
```

### Docker

```bash
# Сборка образов
make docker-build

# Запуск через Docker Compose
make docker-run
```

## Архитектура

GophKeeper построен на принципах **Clean Architecture** с четким разделением слоев:

### Компоненты:
1. **Сервер** - обрабатывает запросы клиентов, управляет пользователями и хранит зашифрованные данные
2. **Клиент** - CLI приложение для взаимодействия с сервером и локального управления данными

### Слои архитектуры:
- **Domain** - бизнес-логика, сущности, доменные сервисы
- **Application** - сценарии использования (use cases)
- **Infrastructure** - внешние зависимости (БД, криптография, конфигурация)
- **Interfaces** - обработчики запросов (gRPC, HTTP)

### Безопасность

- Все данные шифруются на стороне клиента перед отправкой на сервер
- Используется end-to-end шифрование
- Сервер не имеет доступа к незашифрованным данным пользователей
- Поддержка JWT токенов для аутентификации

## Лицензия

MIT License