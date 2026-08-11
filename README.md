# Health Tracker — Pet-проект для подготовки к Middle Backend Golang Developer

## Цель проекта

Получить практический опыт работы с технологиями, которые нужны для позиции Middle Backend Golang Developer:
**Redis, Kafka, RabbitMQ, gRPC, Kubernetes, Docker, микросервисная архитектура, работа с БД.**

Проект должен дать реальные истории для собеседования вместо абстрактного "я читал про Kafka".

---

## Описание проекта

**Health Tracker** — приложение для здоровья с подсчётом калорий и учётом активности за день.
Пользователь логирует приёмы пищи и тренировки, приложение считает дневной баланс калорий
относительно его цели (похудение / поддержание / набор массы) и уведомляет при превышении лимита.

---

## Архитектура

**Тип архитектуры: микросервисная.**
Внутри каждого сервиса — простая **layered-архитектура** (handler → service → repository),
без переусложнения через DDD.

### Схема

```
                    ┌─────────────┐
                    │  API Gateway │ (опционально, можно сразу без него)
                    └──────┬──────┘
                           │
        ┌──────────────────┼──────────────────┐
        │                  │                  │
   ┌────▼────┐       ┌─────▼─────┐      ┌─────▼──────┐
   │  User   │◄─gRPC─►│ Nutrition │      │  Activity  │
   │ Service │        │  Service  │      │  Service   │
   └────┬────┘        └─────┬─────┘      └─────┬──────┘
        │                   │                   │
        │              (Kafka topic: "food_logged", "activity_logged")
        │                   │                   │
        │                   └─────────┬─────────┘
        │                             │
        │                      ┌──────▼───────┐
        │                      │ Notification │
        │                      │   Service    │ ◄── RabbitMQ (email/push задачи)
        │                      └──────────────┘
        │
   Postgres (у каждого сервиса своя БД/схема)
   Redis (кэш дневной сводки, rate limiting)
```

### Сервисы и их роли

1. **User Service**
   - Регистрация/логин (JWT)
   - Профиль пользователя (вес, рост, возраст, цель)
   - Расчёт дневной нормы калорий (формула Миффлина-Сан Жеора)
   - Отдаёт данные другим сервисам по gRPC
   - Здесь же отрабатывается тема "HTTP и REST API" из плана подготовки

2. **Nutrition Service**
   - CRUD приёмов пищи
   - Справочник продуктов (можно захардкодить или взять Open Food Facts API)
   - Подсчёт калорий/БЖУ за день
   - Публикует событие `food_logged` в Kafka при добавлении записи

3. **Activity Service**
   - CRUD тренировок/активности
   - Расчёт потраченных калорий
   - Публикует событие `activity_logged` в Kafka

4. **Notification Service**
   - Consumer из Kafka (слушает события еды/активности)
   - Считает дневной баланс калорий
   - При превышении лимита кладёт задачу в RabbitMQ на отправку уведомления
   - Отдельный воркер разбирает очередь RabbitMQ и "отправляет" уведомление (можно заглушкой — лог/консоль)
   - Своей БД может не иметь (stateless), либо небольшая таблица `notifications_log`

### Роль каждой технологии

| Технология | Зачем нужна в проекте |
|---|---|
| **net/http + chi** | REST API каждого сервиса |
| **gRPC** | Внутреннее общение User ↔ Nutrition ↔ Activity |
| **Kafka** | События `food_logged` / `activity_logged`, несколько consumer'ов реагируют независимо |
| **RabbitMQ** | Очередь задач на отправку уведомлений (гарантия доставки конкретному воркеру) |
| **Redis** | Кэш дневной сводки (`daily_summary:{user_id}:{date}`), rate limiting на API |
| **Postgres** | Основное хранилище, отдельная БД/схема на каждый сервис |
| **Docker** | Упаковка каждого сервиса |
| **Kubernetes** | Оркестрация — деплой, масштабирование, self-healing |

### Важный архитектурный момент

У каждого сервиса своя БД → **нет foreign key между `user_id` в Nutrition/Activity и таблицей `users`**.
Это нормальное следствие микросервисной архитектуры — консистентность между сервисами
обеспечивается на уровне приложения (eventual consistency), а не на уровне БД.

---

## Структура репозитория (монорепо)

```
health-tracker/
├── services/
│   ├── user-service/
│   │   ├── cmd/
│   │   │   └── main.go
│   │   ├── internal/
│   │   │   ├── handler/       # HTTP-хендлеры (REST)
│   │   │   ├── grpc/          # gRPC сервер
│   │   │   ├── service/       # бизнес-логика
│   │   │   ├── repository/    # работа с БД
│   │   │   └── model/         # структуры данных
│   │   ├── migrations/
│   │   ├── Dockerfile
│   │   └── go.mod
│   ├── nutrition-service/
│   │   └── ... (та же структура)
│   ├── activity-service/
│   │   └── ...
│   └── notification-service/
│       └── ...
├── pkg/
│   ├── proto/                  # .proto файлы и сгенерированный gRPC-код
│   └── logger/
├── deploy/
│   ├── docker-compose.yml
│   └── k8s/
│       ├── user-service.yaml
│       ├── nutrition-service.yaml
│       └── ...
└── README.md
```

Каждый сервис — свой `go.mod` (подчёркивает независимость).

---

## Схема БД по сервисам

### User Service — БД `users_db`
```
users
├── id (uuid, PK)
├── email
├── password_hash
├── name
├── weight_kg
├── height_cm
├── age
├── gender
├── goal                    -- lose_weight / maintain / gain_weight
├── daily_calorie_target    -- рассчитывается при регистрации/обновлении профиля
├── created_at
```

### Nutrition Service — БД `nutrition_db`
```
food_items                      -- справочник продуктов
├── id (uuid, PK)
├── name
├── calories_per_100g
├── protein_per_100g
├── fat_per_100g
├── carbs_per_100g

food_logs                       -- записи "что съел"
├── id (uuid, PK)
├── user_id                     -- внешний id, без FK на чужую БД
├── food_item_id (FK)
├── amount_grams
├── meal_type                   -- breakfast/lunch/dinner/snack
├── logged_at
```

### Activity Service — БД `activity_db`
```
activity_types                  -- справочник (бег, ходьба, силовая и т.д.)
├── id (uuid, PK)
├── name
├── calories_per_minute_per_kg

activity_logs
├── id (uuid, PK)
├── user_id
├── activity_type_id (FK)
├── duration_minutes
├── calories_burned             -- рассчитывается при записи
├── logged_at
```

---

## План по этапам

- [ ] **Этап 0** — репозиторий, структура проекта (монорепо)
- [ ] **Этап 1** — User Service: БД, регистрация/логин, JWT, REST API
- [ ] **Этап 2** — gRPC между User и Nutrition (Nutrition получает профиль/цель пользователя)
- [ ] **Этап 3** — Nutrition Service: CRUD еды, подсчёт калорий
- [ ] **Этап 4** — Activity Service: CRUD активности, расчёт потраченных калорий
- [ ] **Этап 5** — Kafka: события `food_logged` / `activity_logged`, продюсеры в Nutrition/Activity
- [ ] **Этап 6** — Notification Service: consumer из Kafka, логика "превышен лимит", producer в RabbitMQ
- [ ] **Этап 7** — Redis: кэш дневной сводки
- [ ] **Этап 8** — Docker Compose — поднять всё локально одной командой
- [ ] **Этап 9** — Kubernetes — манифесты, деплой в minikube/kind локально

---

## Организация рабочего процесса

- **Этот чат** — планирование, архитектурные решения, разбор концепций, ревью кусков кода
- **Claude Code в VS Code** — непосредственное написание кода, ревью в контексте всего репозитория
- **Дисциплина**: для каждого нового куска — сначала план здесь, потом самостоятельное написание в VS Code,
  Claude Code зовётся для объяснений/ревью, а не для написания кода за тебя