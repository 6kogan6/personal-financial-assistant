# Контракт API (v0)

База: `/api`  
Ответы — JSON.  
Ошибки всегда в формате: `{ "error": "сообщение" }`

## Общие правила

- Дата операции: `YYYY-MM-DD` (пример: `2026-02-23`)
- Месяц для отчётов: `YYYY-MM` (пример: `2026-02`)
- Деньги: в центах (целое число) `amount_cents`

---

## 1) Health

### GET /api/health

Проверка, что сервер жив.

Ответ 200:

```json
{ "status": "ok" }
```

## 2) Auth

### POST /api/auth/register

Регистрация пользователя.
Request body:

```json
{ "email": "user@example.com", "password": "строка (мин 8)" }
```

Ответ: 201
Response body:

```json
{ "id": "uuid", "email": "user@example.com" }
```

400 - неверные данные
409 - email уже существует

### POST /api/auth/login

Вход пользователя.
Request body:

```json
{ "email": "user@example.com", "password": "строка" }
```

Ответ: 200
Response body:

```json
{ "access_token": "...", "token_type": "Bearer" }
```

400 - неверные данные
401 - неверный email/пароль

## 3) Categories

### GET /api/categories

Список категорий пользователя.

Ответ: 200
Response body:

```json
{
  "items": [
    { "id": "uuid", "name": "Еда", "type": "expense" },
    { "id": "uuid", "name": "Зарплата", "type": "income" }
  ]
}
```

### POST /api/categories

Создание категории.
Request body:

```json
{ "name": "Еда", "type": "expense" }
```

Ответ: 201
Response body:

```json
{ "id": "uuid", "name": "Еда", "type": "expense" }
```

400 — неверные данные

## 4) Transactions

### GET /api/transactions

Список транзакций.

Query params (опционально):

- from=YYYY-MM-DD
- to=YYYY-MM-DD
- category_id=uuid
- q=строка

Ответ: 200
Response body:

```json
{
  "items": [
    {
      "id": "uuid",
      "occurred_at": "2026-02-23",
      "type": "expense",
      "amount_cents": 1299,
      "category_id": "uuid",
      "merchant": "Starbucks",
      "note": "латте"
    }
  ],
  "total": 1
}
```

### POST /api/transactions

Создание транзакции.

Request body:

```json
{
  "occurred_at": "2026-02-23",
  "type": "expense",
  "amount_cents": 1299,
  "category_id": "uuid",
  "merchant": "Starbucks",
  "note": "латте"
}
```

Ответ: 201
Response body:

```json
{
  "id": "uuid",
  "occurred_at": "2026-02-23",
  "type": "expense",
  "amount_cents": 1299,
  "category_id": "uuid",
  "merchant": "Starbucks",
  "note": "латте"
}
```

400 - неверные данные
404 - категория не найдена

## 5) Reports

### GET /api/reports/summary?month=YYYY-MM

Сводка за месяц для дашборда.

Ответ: 200
Response body:

```json
{
  "month": "2026-02",
  "totals": {
    "income_cents": 250000,
    "expense_cents": 180000,
    "balance_cents": 70000
  },
  "by_category": [
    {
      "category_id": "uuid",
      "category_name": "Еда",
      "type": "expense",
      "amount_cents": 45000
    }
  ],
  "daily": [{ "date": "2026-02-01", "income_cents": 0, "expense_cents": 1299 }]
}
```

400 - неверный формат month
