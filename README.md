# 🚬 smoker

Telegram-бот, який раз на певний час пише у групу, **чи зараз хороший час вийти
на перекур** — і завжди пояснює, *чому саме зараз*, через набір смішних умов,
що парсяться з відкритих API для конкретного місця.

Приклад повідомлення:

```
🚬 ЧАС НА ПЕРЕКУР!
📍 Київ, Ukraine · Схоже на правду

Чому саме зараз:
• 🛰️ МКС просто зараз за 1240 км над вами — космонавти дивляться, вийдіть гідно
• ☀️ UV-індекс 6 — встигнеш засмагнути за одну цигарку, безкоштовний солярій
• 🕐 Час 13:37 — елітний leet-час, справжні профі курять саме зараз
• 🃏 Наукове обґрунтування дня: «...» — після такого треба вийти провітритись

⚖️ Підсумковий індекс перекуру: +5
```

## Архітектура (мікросервіси, DDD)

| Сервіс | Bounded context | Роль |
|---|---|---|
| `gateway` | Публікація | Telegram-бот, команди, реєстрація чатів (Postgres), планувальник |
| `orchestrator` | Core domain (`SmokeAdvisor`) | Агрегує умови провайдерів, рахує вердикт, складає текст |
| `cosmos` | Природа та космос | UV, вітер, тиск, пилок, Kp-індекс, МКС, фаза місяця |
| `chronos` | Час і числа | % доби/року, прості/Фібоначчі-хвилини, свята, «цього дня» |
| `chaos` | Випадковість заради сміху | dad jokes, Chuck Norris, котофакти, кубик |

```
Telegram ──▶ gateway ──gRPC(Recommend)──▶ orchestrator ──gRPC(Evaluate)──▶ cosmos
                │ (scheduler, Postgres)          │ (fan-out, scoring)    ├▶ chronos
                                                                         └▶ chaos
```

- Зв'язок між сервісами — **gRPC** (контракти у `proto/`, згенеровані стаби у `pkg/proto/`).
- Спільна доменна мова — `pkg/smoke` (shared kernel).
- Кожен `internal/<context>` поділений на `domain` (чисті правила + тести), `app`
  (use-cases), `infra` (ACL до зовнішніх API / gRPC).
- **Graceful degradation:** якщо якесь API чи навіть цілий сервіс лежить — умова
  просто пропускається, вердикт рахується з решти.

## Запуск

```bash
cp .env.example .env
# впишіть TELEGRAM_TOKEN від @BotFather у .env
docker compose up --build
```

Далі у Telegram:
1. Додайте бота у групу (і вимкніть privacy mode у @BotFather → `/setprivacy` →
   Disable, щоб бот бачив команди в групі).
2. `/setlocation Львів` — задати місто (або `/setlocation 49.84,24.03`).
3. `/setschedule 60` — нагадувати щогодини (також `30m`, `2h`).
4. `/smoke` — перевірити прямо зараз.
5. `/stop` — припинити нагадування.

## Розробка

```bash
# усі тести (домен + парсери, без мережі)
go test ./...

# перегенерувати gRPC-стаби після зміни proto/
buf generate

# зібрати конкретний сервіс
go run ./cmd/cosmos
```

Потрібні інструменти для кодогену: `buf`, `protoc-gen-go`, `protoc-gen-go-grpc`
(встановлюються через `go install`).

## Конфігурація (env)

| Змінна | Сервіс | За замовчуванням |
|---|---|---|
| `TELEGRAM_TOKEN` | gateway | — (обовʼязково) |
| `DATABASE_URL` | gateway | `postgres://smoker:smoker@postgres:5432/smoker` |
| `ORCHESTRATOR_ADDR` | gateway | `orchestrator:9100` |
| `DEFAULT_PLACE/LAT/LON/TZ` | gateway | Київ |
| `DEFAULT_INTERVAL_MIN` | gateway | `90` |
| `WORK_START_HOUR` / `WORK_END_HOUR` | gateway | `9` / `19` |
| `COSMOS_ADDR/CHRONOS_ADDR/CHAOS_ADDR` | orchestrator | `*:910x` |
| `COUNTRY_CODE` | chronos | `UA` (свята Nager.Date) |
| `WIKI_LANG` | chronos/gateway | `uk` |
| `GRPC_ADDR` | усі gRPC-сервіси | `:910x` |

## Джерела даних (усі безкоштовні, переважно без ключів)

Open-Meteo (погода, пилок, геокодинг), NOAA SWPC (Kp-індекс), wheretheiss.at
(МКС), Nager.Date (свята), Wikipedia REST (цього дня), icanhazdadjoke,
api.chucknorris.io, catfact.ninja.
