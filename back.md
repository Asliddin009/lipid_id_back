# Backend API — Спецификация эндпоинтов

Ниже перечислены все методы, необходимые для замены моковых данных в приложении.

**Тело `GET /profile` ответа:**

```json
{
  "gender": "male",
  "age": 28,
  "height": 175,
  "weight": 72,
  "bmi": 23.5,
  "goal": "weight_loss",
  "activityLevel": "moderate",
  "language": "ru",
  "notificationsEnabled": true
}
```

---

## 3. Питание / Еда (`AddFoodBloc`, `NutritionBloc`)

Заменяет: `InMemoryAddFoodRepository`, `IAddFood`, `INutritionRepo`.

| Метод | Эндпоинт | Описание |
|-------|----------|----------|
| `GET` | `/data/food-entries` | Список всех записей о еде (с пагинацией) |
| `GET` | `/data/food-entries?date=2026-03-20` | Записи за конкретную дату |
| `POST` | `/data/food-entries` | Добавить запись о еде |
| `DELETE` | `/data/food-entries/{id}` | Удалить запись о еде |
| `GET` | `/data/food-entries/summary?date=2026-03-20` | Дневная статистика (суммы КБЖУ) |
| `GET` | `/data/food-entries/summary/weekly` | Недельная статистика |

**Тело `POST /data/food-entries`:**

```json
{
  "dishName": "Курица с овощами",
  "mealTime": "breakfast",
  "calories": 450,
  "proteins": 25.0,
  "fats": 28.0,
  "carbs": 22.0,
  "imageUrl": null,
  "createdAt": "2026-03-20T08:30:00Z"
}
```

**Ответ `GET /data/food-entries/summary`:**

```json
{
  "date": "2026-03-20",
  "totalCalories": 1150,
  "totalProteins": 75.0,
  "totalFats": 62.0,
  "totalCarbs": 52.0,
  "byMealTime": {
    "breakfast": 450,
    "lunch": 320,
    "snack": 0,
    "dinner": 380
  }
}
```

---

## 4. Распознавание еды по фото (`FoodRecognitionScreen`)

Заменяет: моковые данные `mockDishName`, `caloriesValue('450')` и т.д.

| Метод | Эндпоинт | Описание |
|-------|----------|----------|
| `POST` | `/data/food-recognition/analyze` | Отправить фото → получить КБЖУ и название блюда |

**Тело:** `multipart/form-data` с полем `image`.

**Ответ:**

```json
{
  "dishName": "Курица с овощами",
  "calories": 450,
  "proteins": 25.0,
  "fats": 28.0,
  "saturatedFats": 12.0,
  "carbs": 22.0,
  "lipidRating": "high",
  "ratingImpact": -0.3,
  "confidence": 0.87
}
```

---

## 5. Главная (`HomeBloc`, `MockHomeRepository`)

Заменяет: `MockHomeRepository` → `HomeRepository`.

| Метод | Эндпоинт | Описание |
|-------|----------|----------|
| `GET` | `/data/home/dashboard` | Главная: рейтинг + график + последние события |

**Ответ:**

```json
{
  "score": 7.5,
  "chartData": [
    { "date": "2026-03-14", "value": 6.8 },
    { "date": "2026-03-15", "value": 7.0 },
    { "date": "2026-03-16", "value": 6.5 },
    { "date": "2026-03-17", "value": 7.2 },
    { "date": "2026-03-18", "value": 7.1 },
    { "date": "2026-03-19", "value": 7.4 },
    { "date": "2026-03-20", "value": 7.5 }
  ],
  "recentEvents": [
    {
      "type": "meal",
      "title": "Приём пищи",
      "subtitle": "Обед — Салат",
      "trailing": "-0.1",
      "color": "red",
      "createdAt": "2026-03-20T12:00:00Z"
    },
    {
      "type": "lab",
      "title": "Лабораторные измерения",
      "subtitle": "LDL: 3.2 ммоль/л",
      "trailing": "+0.5",
      "color": "green",
      "createdAt": "2026-03-19T10:00:00Z"
    },
    {
      "type": "anthropometry",
      "title": "Антропометрия",
      "subtitle": "72 кг",
      "trailing": "",
      "color": "gray",
      "createdAt": "2026-03-17T09:00:00Z"
    }
  ]
}
```

---

## 6. Графики (`ChartsScreen`)

Заменяет: моковые данные в `ChartsScreen`.

| Метод | Эндпоинт | Описание |
|-------|----------|----------|
| `GET` | `/data/charts/lipid-trend?period=7d` | Динамика липидного баланса |
| `GET` | `/data/charts/lipid-trend?period=30d` | За 30 дней |
| `GET` | `/data/charts/nutrition-trend?period=7d` | Динамика КБЖУ |
| `GET` | `/data/charts/weight-trend?period=30d` | Динамика веса |

**Ответ `GET /data/charts/lipid-trend`:**

```json
{
  "period": "7d",
  "points": [
    { "date": "2026-03-14", "value": 6.8 },
    { "date": "2026-03-15", "value": 7.0 }
  ]
}
```

---

## 7. Анализы (`AnalysesScreen`)

Заменяет: моковые данные в `AnalysesScreen`.

| Метод | Эндпоинт | Описание |
|-------|----------|----------|
| `GET` | `/data/analyses` | Список всех анализов пользователя |
| `POST` | `/data/analyses` | Добавить результаты анализов |
| `GET` | `/data/analyses/{id}` | Детали конкретного анализа |
| `DELETE` | `/data/analyses/{id}` | Удалить анализ |

**Тело `POST /data/analyses`:**

```json
{
  "date": "2026-03-20",
  "totalCholesterol": 5.2,
  "ldl": 3.1,
  "hdl": 1.4,
  "triglycerides": 1.5,
  "vldl": 0.7
}
```

---

## 8. Антропометрия / Вес

Заменяет: кнопку «Добавить вес» из `AppFloatingButton`.

| Метод | Эндпоинт | Описание |
|-------|----------|----------|
| `GET` | `/data/weight` | История замеров веса |
| `POST` | `/data/weight` | Добавить замер веса |
| `DELETE` | `/data/weight/{id}` | Удалить замер |

**Тело `POST /data/weight`:**

```json
{
  "value": 72.0,
  "date": "2026-03-20"
}
```

---

## 9. Уведомления

| Метод | Эндпоинт | Описание |
|-------|----------|----------|
| `POST` | `/data/devices` | Зарегистрировать FCM/APN токен |
| `GET` | `/data/notifications` | Списоке увдомлений |
| `PATCH` | `/data/notifications/{id}/read` | Пометить как прочитанное |
