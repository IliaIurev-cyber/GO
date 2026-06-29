$baseUrl = "http://localhost:8080"

Write-Host "=== Тестирование операций со счетами ===" -ForegroundColor Cyan

# 1. Регистрация
Write-Host "`n1. Регистрация..." -ForegroundColor Yellow
$registerBody = @{
    username = "testuser"
    email = "test@example.com"
    password = "TestPassword123"
} | ConvertTo-Json
Invoke-RestMethod -Uri "$baseUrl/register" -Method Post -Body $registerBody -ContentType "application/json" | Out-Null
Write-Host "   ✅ Пользователь зарегистрирован" -ForegroundColor Green

# 2. Вход
Write-Host "`n2. Вход..." -ForegroundColor Yellow
$loginBody = @{
    email = "test@example.com"
    password = "TestPassword123"
} | ConvertTo-Json
$loginResponse = Invoke-RestMethod -Uri "$baseUrl/login" -Method Post -Body $loginBody -ContentType "application/json"
$token = $loginResponse.token
Write-Host "   ✅ Токен получен" -ForegroundColor Green

# 3. Создание счета
Write-Host "`n3. Создание счета..." -ForegroundColor Yellow
$accountBody = @{ currency = "RUB" } | ConvertTo-Json
$account = Invoke-RestMethod -Uri "$baseUrl/api/accounts" -Method Post -Body $accountBody -ContentType "application/json" -Headers @{Authorization = "Bearer $token"}
Write-Host "   ✅ Счет создан: $($account.number)" -ForegroundColor Green
Write-Host "   ID: $($account.id)" -ForegroundColor Gray
Write-Host "   Баланс: $($account.balance) $($account.currency)" -ForegroundColor Gray

$accountId = $account.id

# 4. Пополнение счета
Write-Host "`n4. Пополнение счета на 1000 RUB..." -ForegroundColor Yellow
$depositBody = @{
    account_id = $accountId
    amount = 1000
} | ConvertTo-Json
$deposit = Invoke-RestMethod -Uri "$baseUrl/api/deposit" -Method Post -Body $depositBody -ContentType "application/json" -Headers @{Authorization = "Bearer $token"}
Write-Host "   ✅ Пополнение выполнено" -ForegroundColor Green
Write-Host "   Сумма: $($deposit.amount) RUB" -ForegroundColor Gray

# 5. Получение счетов
Write-Host "`n5. Получение списка счетов..." -ForegroundColor Yellow
$accounts = Invoke-RestMethod -Uri "$baseUrl/api/accounts" -Method Get -Headers @{Authorization = "Bearer $token"}
Write-Host "   ✅ Найдено счетов: $($accounts.Count)" -ForegroundColor Green
foreach ($acc in $accounts) {
    Write-Host "   - Счет: $($acc.number), Баланс: $($acc.balance) $($acc.currency)" -ForegroundColor Gray
}

# 6. Создание второго счета для перевода
Write-Host "`n6. Создание второго счета..." -ForegroundColor Yellow
$account2 = Invoke-RestMethod -Uri "$baseUrl/api/accounts" -Method Post -Body $accountBody -ContentType "application/json" -Headers @{Authorization = "Bearer $token"}
$account2Id = $account2.id
Write-Host "   ✅ Второй счет создан: $($account2.number)" -ForegroundColor Green

# 7. Пополнение второго счета
Write-Host "`n7. Пополнение второго счета на 500 RUB..." -ForegroundColor Yellow
$depositBody2 = @{
    account_id = $account2Id
    amount = 500
} | ConvertTo-Json
Invoke-RestMethod -Uri "$baseUrl/api/deposit" -Method Post -Body $depositBody2 -ContentType "application/json" -Headers @{Authorization = "Bearer $token"} | Out-Null
Write-Host "   ✅ Второй счет пополнен" -ForegroundColor Green

# 8. Перевод между счетами
Write-Host "`n8. Перевод 300 RUB со счета $accountId на счет $account2Id..." -ForegroundColor Yellow
$transferBody = @{
    from_account_id = $accountId
    to_account_id = $account2Id
    amount = 300
    description = "Test transfer"
} | ConvertTo-Json
$transfer = Invoke-RestMethod -Uri "$baseUrl/api/transfer" -Method Post -Body $transferBody -ContentType "application/json" -Headers @{Authorization = "Bearer $token"}
Write-Host "   ✅ Перевод выполнен" -ForegroundColor Green
Write-Host "   Сумма: $($transfer.amount) RUB" -ForegroundColor Gray
Write-Host "   Описание: $($transfer.description)" -ForegroundColor Gray

# 9. Проверка балансов
Write-Host "`n9. Проверка балансов..." -ForegroundColor Yellow
$accounts = Invoke-RestMethod -Uri "$baseUrl/api/accounts" -Method Get -Headers @{Authorization = "Bearer $token"}
foreach ($acc in $accounts) {
    Write-Host "   Счет: $($acc.number), Баланс: $($acc.balance) $($acc.currency)" -ForegroundColor Gray
}

Write-Host "`n=== Тестирование завершено ===" -ForegroundColor Cyan
