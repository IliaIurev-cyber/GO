$baseUrl = "http://localhost:8080"

Write-Host "=== Полное тестирование Bank Service API ===" -ForegroundColor Cyan
Write-Host ""

# 1. Регистрация
Write-Host "1. Регистрация пользователя..." -ForegroundColor Yellow
$registerBody = @{
    username = "testuser"
    email = "test@example.com"
    password = "TestPassword123"
} | ConvertTo-Json

try {
    $registerResult = Invoke-RestMethod -Uri "$baseUrl/register" -Method Post -Body $registerBody -ContentType "application/json"
    Write-Host "   ✅ Пользователь зарегистрирован: $($registerResult.username)" -ForegroundColor Green
    Write-Host "   ID: $($registerResult.id)" -ForegroundColor Gray
} catch {
    Write-Host "   ⚠️  Пользователь уже существует или ошибка: $($_.Exception.Message)" -ForegroundColor Yellow
}

# 2. Вход
Write-Host "`n2. Вход в систему..." -ForegroundColor Yellow
$loginBody = @{
    email = "test@example.com"
    password = "TestPassword123"
} | ConvertTo-Json

try {
    $loginResponse = Invoke-RestMethod -Uri "$baseUrl/login" -Method Post -Body $loginBody -ContentType "application/json"
    $token = $loginResponse.token
    Write-Host "   ✅ Токен получен" -ForegroundColor Green
    Write-Host "   User: $($loginResponse.user.username)" -ForegroundColor Gray
} catch {
    Write-Host "   ❌ Ошибка входа: $($_.Exception.Message)" -ForegroundColor Red
    exit
}

# 3. Получение курса ЦБ РФ
Write-Host "`n3. Получение курса ЦБ РФ..." -ForegroundColor Yellow
try {
    $rate = Invoke-RestMethod -Uri "$baseUrl/api/cbr-rate" -Method Get
    Write-Host "   ✅ Ключевая ставка: $($rate.rate)%" -ForegroundColor Green
    Write-Host "   Дата: $($rate.date)" -ForegroundColor Gray
} catch {
    Write-Host "   ⚠️  Ошибка получения курса: $($_.Exception.Message)" -ForegroundColor Yellow
}

# 4. Создание счета
Write-Host "`n4. Создание счета..." -ForegroundColor Yellow
$accountBody = @{ currency = "RUB" } | ConvertTo-Json
try {
    $account = Invoke-RestMethod -Uri "$baseUrl/api/accounts" -Method Post -Body $accountBody -ContentType "application/json" -Headers @{Authorization = "Bearer $token"}
    Write-Host "   ✅ Счет создан" -ForegroundColor Green
    Write-Host "   Номер: $($account.number)" -ForegroundColor Gray
    Write-Host "   Баланс: $($account.balance) $($account.currency)" -ForegroundColor Gray
    $accountId = $account.id
} catch {
    Write-Host "   ❌ Ошибка создания счета: $($_.Exception.Message)" -ForegroundColor Red
    exit
}

# 5. Пополнение счета
Write-Host "`n5. Пополнение счета на 1000 RUB..." -ForegroundColor Yellow
$depositBody = @{
    account_id = $accountId
    amount = 1000
} | ConvertTo-Json
try {
    $deposit = Invoke-RestMethod -Uri "$baseUrl/api/deposit" -Method Post -Body $depositBody -ContentType "application/json" -Headers @{Authorization = "Bearer $token"}
    Write-Host "   ✅ Пополнено: $($deposit.amount) RUB" -ForegroundColor Green
} catch {
    Write-Host "   ❌ Ошибка пополнения: $($_.Exception.Message)" -ForegroundColor Red
}

# 6. Получение списка счетов
Write-Host "`n6. Получение списка счетов..." -ForegroundColor Yellow
try {
    $accounts = Invoke-RestMethod -Uri "$baseUrl/api/accounts" -Method Get -Headers @{Authorization = "Bearer $token"}
    Write-Host "   ✅ Найдено счетов: $($accounts.Count)" -ForegroundColor Green
    foreach ($acc in $accounts) {
        Write-Host "   - $($acc.number): $($acc.balance) $($acc.currency)" -ForegroundColor Gray
    }
} catch {
    Write-Host "   ❌ Ошибка получения счетов: $($_.Exception.Message)" -ForegroundColor Red
}

# 7. Создание второго счета
Write-Host "`n7. Создание второго счета..." -ForegroundColor Yellow
try {
    $account2 = Invoke-RestMethod -Uri "$baseUrl/api/accounts" -Method Post -Body $accountBody -ContentType "application/json" -Headers @{Authorization = "Bearer $token"}
    Write-Host "   ✅ Второй счет создан" -ForegroundColor Green
    Write-Host "   Номер: $($account2.number)" -ForegroundColor Gray
    $account2Id = $account2.id
} catch {
    Write-Host "   ❌ Ошибка создания второго счета: $($_.Exception.Message)" -ForegroundColor Red
}

# 8. Пополнение второго счета
Write-Host "`n8. Пополнение второго счета на 500 RUB..." -ForegroundColor Yellow
$depositBody2 = @{
    account_id = $account2Id
    amount = 500
} | ConvertTo-Json
try {
    Invoke-RestMethod -Uri "$baseUrl/api/deposit" -Method Post -Body $depositBody2 -ContentType "application/json" -Headers @{Authorization = "Bearer $token"} | Out-Null
    Write-Host "   ✅ Второй счет пополнен" -ForegroundColor Green
} catch {
    Write-Host "   ❌ Ошибка пополнения второго счета: $($_.Exception.Message)" -ForegroundColor Red
}

# 9. Перевод между счетами
Write-Host "`n9. Перевод 300 RUB со счета $accountId на счет $account2Id..." -ForegroundColor Yellow
$transferBody = @{
    from_account_id = $accountId
    to_account_id = $account2Id
    amount = 300
    description = "Test transfer"
} | ConvertTo-Json
try {
    $transfer = Invoke-RestMethod -Uri "$baseUrl/api/transfer" -Method Post -Body $transferBody -ContentType "application/json" -Headers @{Authorization = "Bearer $token"}
    Write-Host "   ✅ Перевод выполнен" -ForegroundColor Green
    Write-Host "   Сумма: $($transfer.amount) RUB" -ForegroundColor Gray
    Write-Host "   Описание: $($transfer.description)" -ForegroundColor Gray
} catch {
    Write-Host "   ❌ Ошибка перевода: $($_.Exception.Message)" -ForegroundColor Red
}

# 10. Проверка итоговых балансов
Write-Host "`n10. Итоговые балансы:" -ForegroundColor Yellow
try {
    $accounts = Invoke-RestMethod -Uri "$baseUrl/api/accounts" -Method Get -Headers @{Authorization = "Bearer $token"}
    foreach ($acc in $accounts) {
        Write-Host "   Счет $($acc.id): $($acc.balance) $($acc.currency)" -ForegroundColor Gray
    }
} catch {
    Write-Host "   ❌ Ошибка получения балансов: $($_.Exception.Message)" -ForegroundColor Red
}

# 11. Тест защищенного эндпоинта
Write-Host "`n11. Доступ к защищенному эндпоинту..." -ForegroundColor Yellow
try {
    $profile = Invoke-RestMethod -Uri "$baseUrl/api/profile" -Method Get -Headers @{Authorization = "Bearer $token"}
    Write-Host "   ✅ Доступ разрешен" -ForegroundColor Green
    Write-Host "   User ID: $($profile.user_id)" -ForegroundColor Gray
    Write-Host "   Message: $($profile.message)" -ForegroundColor Gray
} catch {
    Write-Host "   ❌ Ошибка доступа: $($_.Exception.Message)" -ForegroundColor Red
}

# 12. Попытка доступа без токена (ожидается ошибка)
Write-Host "`n12. Попытка доступа без токена (ожидается ошибка 401)..." -ForegroundColor Yellow
try {
    Invoke-RestMethod -Uri "$baseUrl/api/profile" -Method Get
    Write-Host "   ❌ Неожиданный успех!" -ForegroundColor Red
} catch {
    Write-Host "   ✅ Ожидаемая ошибка: 401 Unauthorized" -ForegroundColor Green
}

Write-Host "`n=== Все тесты завершены! ===" -ForegroundColor Cyan
