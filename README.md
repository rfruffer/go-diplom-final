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

# Тесты с покрытием
make test-cover

# Юнит-тесты (по слоям)
go test ./tests/unit/...

# Тесты доменной логики
go test ./tests/unit/domain/...

# Тесты инфраструктуры
go test ./tests/unit/infrastructure/...

# Интеграционные тесты
go test ./tests/integration/...

# End-to-end тесты
go test ./tests/e2e/...
```

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

### Сервер

```bash
./build/server --config config/server.yaml
```

### Клиент

```bash
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