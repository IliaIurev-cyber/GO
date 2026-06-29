#  Bank Service API

REST API для банковского сервиса на Go с аутентификацией JWT, управлением счетами, интеграцией с ЦБ РФ и SMTP.

##  Содержание

- [Технологии](#технологии)
- [Функциональность](#функциональность)
- [Архитектура](#архитектура)
- [Установка и запуск](#установка-и-запуск)
- [API Эндпоинты](#api-эндпоинты)
- [Тестирование](#тестирование)
- [Структура проекта](#структура-проекта)
- [Планы по развитию](#планы-по-развитию)

## 🚀 Технологии

| Компонент | Технология |
|-----------|------------|
| Язык | Go 1.23 |
| Маршрутизация | gorilla/mux |
| Аутентификация | JWT (golang-jwt/jwt/v5) |
| Хеширование | bcrypt, HMAC-SHA256 |
| Логирование | logrus |
| Интеграции | SOAP (ЦБ РФ), SMTP |
| Хранилище | In-memory (для разработки) |

##  Реализованная функциональность

### Пользователи
-  Регистрация с проверкой уникальности email
-  Аутентификация с выдачей JWT токена (срок действия 24 часа)
-  Хеширование паролей (bcrypt)

### Счета
-  Создание банковских счетов
-  Просмотр списка счетов
-  Пополнение баланса
-  Переводы между счетами

### Интеграции
-  Интеграция с ЦБ РФ (получение ключевой ставки через SOAP)
-  Интеграция с SMTP (отправка email уведомлений)

### Безопасность
-  JWT аутентификация
-  Middleware для проверки токенов
-  Хеширование паролей (bcrypt)
-  HMAC-SHA256 для проверки целостности

### Дополнительно
-  Логирование всех запросов
-  Модульная архитектура (чистая архитектура)
-  In-memory хранилище (не требует БД для разработки)

### Компоненты системы

#### 1. Модели (internal/models)
- `user.go` - Пользователи
- `account.go` - Банковские счета
- `card.go` - Карты
- `transaction.go` - Транзакции
- `credit.go` - Кредиты

#### 2. Репозитории (internal/repository)
- `memory_user_repo.go` - In-memory хранилище пользователей
- `memory_account_repo.go` - In-memory хранилище счетов
- `memory_transaction_repo.go` - In-memory хранилище транзакций

#### 3. Сервисы (internal/service)
- `auth_service.go` - Аутентификация и авторизация
- `account_service.go` - Управление счетами
- `cbr_service.go` - Интеграция с ЦБ РФ
- `email_service.go` - Отправка email

#### 4. Обработчики (internal/handler)
- `auth_handler.go` - Эндпоинты аутентификации
- `account_handler.go` - Эндпоинты управления счетами

#### 5. Middleware (internal/middleware)
- `auth.go` - Проверка JWT токенов
- `logger.go` - Логирование запросов

##  Установка и запуск

### Требования
- Go 1.23 или выше
- (Опционально) PostgreSQL 17 для production

### 1. Клонирование репозитория

```bash
git clone https://github.com/IliaIurev-cyber/GO.git
cd bank-service
2. Настройка окружения
Создайте файл .env в корне проекта:

env
# Сервер
SERVER_PORT=8080

# JWT
JWT_SECRET=your-super-secret-jwt-key-change-this-in-production

# HMAC
HMAC_SECRET=your-hmac-secret-key-change-this

# SMTP (опционально)
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USER=your_email@gmail.com
SMTP_PASSWORD=your_app_password

# ЦБ РФ
CENTRAL_BANK_URL=https://www.cbr.ru/DailyInfoWebServ/DailyInfo.asmx

# База данных (опционально)
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=bankdb

3. Установка зависимостей
bash
go mod download
go mod tidy

4. Запуск сервера
bash
# Запуск в режиме разработки
go run cmd/main.go

# Сборка бинарника
go build -o bank-service.exe cmd/main.go

# Запуск собранного бинарника
./bank-service.exe

Автоматические тесты
1. Полное тестирование API
Запустите скрипт для полного тестирования всех эндпоинтов:

powershell
# Windows
.\test-all.ps1

# Linux/Mac
./test-all.sh
Скрипт выполняет:

Регистрацию пользователя

Вход и получение JWT токена

Получение курса ЦБ РФ

Создание счета

Пополнение счета

Получение списка счетов

Создание второго счета

Пополнение второго счета

Перевод между счетами

Проверка итоговых балансов

Доступ к защищенному эндпоинту

Проверка ошибки 401 без токена

2. Тестирование только счетов
powershell
.\test-accounts.ps1
Ручное тестирование с curl
Регистрация
bash
curl -X POST http://localhost:8080/register \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","email":"test@example.com","password":"TestPassword123"}'
Вход
bash
curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"TestPassword123"}'
Создание счета
bash
curl -X POST http://localhost:8080/api/accounts \
  -H "Authorization: Bearer <TOKEN>" \
  -H "Content-Type: application/json" \
  -d '{"currency":"RUB"}'
Пополнение счета
bash
curl -X POST http://localhost:8080/api/deposit \
  -H "Authorization: Bearer <TOKEN>" \
  -H "Content-Type: application/json" \
  -d '{"account_id":1,"amount":1000}'
Перевод
bash
curl -X POST http://localhost:8080/api/transfer \
  -H "Authorization: Bearer <TOKEN>" \
  -H "Content-Type: application/json" \
  -d '{"from_account_id":1,"to_account_id":2,"amount":300,"description":"Test transfer"}'

Тестирование с PowerShell
powershell
# Получение токена
$token = (Invoke-RestMethod -Uri "http://localhost:8080/login" -Method Post -Body '{"email":"test@example.com","password":"TestPassword123"}' -ContentType "application/json").token

# Создание счета
Invoke-RestMethod -Uri "http://localhost:8080/api/accounts" -Method Post -Body '{"currency":"RUB"}' -ContentType "application/json" -Headers @{Authorization = "Bearer $token"}

# Пополнение счета
Invoke-RestMethod -Uri "http://localhost:8080/api/deposit" -Method Post -Body '{"account_id":1,"amount":1000}' -ContentType "application/json" -Headers @{Authorization = "Bearer $token"}
