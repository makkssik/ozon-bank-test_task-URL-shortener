URL Shortener
REST API сервис для создания сокращенных ссылок. Выполнено в рамках тестового задания на позицию стажера-разработчика Go.

## Архитектура проекта
Проект сделан на основе Луковой многослойной архитектуры.

Domain (internal/domain): Ядро приложения. Содержит бизнес-сущности и доменные ошибки. Не имеет внешних зависимостей. Создание бизнес-сущности происходит без привязки к ошибкам (только генерация), а строгая валидация формата оригинального URL и проверка алиасов вынесены в отдельные чистые функции.

Application (internal/application): Слой Use Cases. Бизнес-логика грамотно разделена на независимые пакеты по юзкейсам. Каждый юзкейс имеет свои собственные, компактные интерфейсы для работы с хранилищем. Слой не знает о том, какая именно база данных используется.

Infrastructure (internal/infrastructure): Слой работы с внешними системами. Реализует интерфейс репозитория двумя способами:

postgres: персистентное хранилище на базе PostgreSQL.

memory: потокобезопасное In-Memory хранилище на базе хэш-таблиц и sync.RWMutex.

Presentation (internal/presentation): Роутер построен на базе go-chi/chi/v5. Включает хендлеры, DTO для запросов/ответов, валидацию и кастомные middleware.


### Запуск в Docker (PostgreSQL)
Поднимает базу данных PostgreSQL, автоматически применяет миграции (migrate/migrate) и запускает сам сервис.

`make docker-postgres`

В docker-compose.yml порт приложения проброшен наружу как 8082. Сервис будет доступен по адресу `http://localhost:8082`.

### Локальный запуск (In-Memory)
Данные хранятся в оперативной памяти и исчезают после остановки сервиса.

`make docker-memory`

Приложение использует параметры из local.yaml и будет доступно по адресу `http://localhost:8081`.

## API Эндпоинты

| Метод | Эндпоинт | Описание | Тело запроса             | Успешный ответ |
| :--- | :--- | :--- |:-------------------------| :--- |
| `POST` | `/` | Создает короткую ссылку (алиас) для переданного URL. | `{"url": "https://..."}` | `201 Created`<br>`{"status": "OK", "alias": "..."}` |
| `GET` | `/{alias}` | Выполняет редирект на оригинальную ссылку по её алиасу. | *-*                      | `302 Found`<br>*(Перенаправление)* |


## Метрики

In-Memory
````
Summary:
  Total:        0.2953 secs
  Slowest:      0.0293 secs
  Fastest:      0.0003 secs
  Average:      0.0028 secs
  Requests/sec: 33860.6088
  
  Total data:   370000 bytes
  Size/request: 37 bytes

Response time histogram:
  0.000 [1]     |
  0.003 [7507]  |■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■
  0.006 [2145]  |■■■■■■■■■■■
  0.009 [237]   |■
  0.012 [37]    |
  0.015 [2]     |
  0.018 [28]    |
  0.021 [0]     |
  0.023 [0]     |
  0.026 [0]     |
  0.029 [43]    |


Latency distribution:
  10%% in 0.0014 secs
  25%% in 0.0018 secs
  50%% in 0.0024 secs
  75%% in 0.0032 secs
  90%% in 0.0045 secs
  95%% in 0.0056 secs
  99%% in 0.0095 secs

Details (average, fastest, slowest):
  DNS+dialup:   0.0000 secs, 0.0000 secs, 0.0052 secs
  DNS-lookup:   0.0000 secs, 0.0000 secs, 0.0030 secs
  req write:    0.0000 secs, 0.0000 secs, 0.0022 secs
  resp wait:    0.0027 secs, 0.0003 secs, 0.0236 secs
  resp read:    0.0001 secs, 0.0000 secs, 0.0050 secs

Status code distribution:
  [201] 10000 responses

````

Postgres
````
Summary:
  Total:        3.6021 secs
  Slowest:      0.1566 secs
  Fastest:      0.0003 secs
  Average:      0.0174 secs
  Requests/sec: 2776.1487
  
  Total data:   370000 bytes
  Size/request: 37 bytes

Response time histogram:
  0.000 [1]     |
  0.016 [5503]  |■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■
  0.032 [2834]  |■■■■■■■■■■■■■■■■■■■■■
  0.047 [1090]  |■■■■■■■■
  0.063 [335]   |■■
  0.078 [153]   |■
  0.094 [29]    |
  0.110 [15]    |
  0.125 [16]    |
  0.141 [19]    |
  0.157 [5]     |


Latency distribution:
  10%% in 0.0016 secs
  25%% in 0.0035 secs
  50%% in 0.0126 secs
  75%% in 0.0261 secs
  90%% in 0.0386 secs
  95%% in 0.0494 secs
  99%% in 0.0752 secs

Details (average, fastest, slowest):
  DNS+dialup:   0.0000 secs, 0.0000 secs, 0.0039 secs
  DNS-lookup:   0.0000 secs, 0.0000 secs, 0.0022 secs
  req write:    0.0000 secs, 0.0000 secs, 0.0019 secs
  resp wait:    0.0174 secs, 0.0003 secs, 0.1566 secs
  resp read:    0.0000 secs, 0.0000 secs, 0.0047 secs

Status code distribution:
  [201] 10000 responses
````