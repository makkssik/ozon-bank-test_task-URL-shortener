URL Shortener
REST API сервис для создания сокращенных ссылок. Выполнено в рамках тестового задания на позицию стажера-разработчика Go.

Архитектура проекта
Проект сделан на основе Луковой многослойной архитектуры.

Domain (internal/domain): Ядро приложения. Содержит бизнес-сущности и доменные ошибки. Не имеет внешних зависимостей. Строгая валидация алфавита и длины алиасов происходит на этапе создания сущности.

Application (internal/application): Слой Use Cases. Содержит интерфейсы для работы с хранилищем и саму бизнес-логику в URLService. Слой не знает о том, какая именно база данных используется.

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
Total:        0.3928 secs
Slowest:      0.0437 secs
Fastest:      0.0001 secs
Average:      0.0036 secs
Requests/sec: 25456.5095

Total data:   360000 bytes
Size/request: 36 bytes

Response time histogram:
0.000 [1]     |
0.004 [6822]  |■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■
0.009 [2082]  |■■■■■■■■■■■■
0.013 [686]   |■■■■
0.018 [256]   |■■
0.022 [97]    |■
0.026 [36]    |
0.031 [11]    |
0.035 [6]     |
0.039 [1]     |
0.044 [2]     |


Latency distribution:
10%% in 0.0003 secs
25%% in 0.0006 secs
50%% in 0.0017 secs
75%% in 0.0056 secs
90%% in 0.0092 secs
95%% in 0.0122 secs
99%% in 0.0194 secs

Details (average, fastest, slowest):
DNS+dialup:   0.0000 secs, 0.0000 secs, 0.0072 secs
DNS-lookup:   0.0000 secs, 0.0000 secs, 0.0057 secs
req write:    0.0000 secs, 0.0000 secs, 0.0021 secs
resp wait:    0.0035 secs, 0.0001 secs, 0.0437 secs
resp read:    0.0000 secs, 0.0000 secs, 0.0012 secs

Status code distribution:
[302] 10000 responses
````

Postgres
````
Summary:
  Total:        2.7574 secs
  Slowest:      0.1181 secs
  Fastest:      0.0005 secs
  Average:      0.0131 secs
  Requests/sec: 3626.6001
  
  Total data:   390000 bytes
  Size/request: 39 bytes

Response time histogram:
  0.000 [1]     |
  0.012 [6367]  |■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■
  0.024 [1620]  |■■■■■■■■■■
  0.036 [1407]  |■■■■■■■■■
  0.048 [423]   |■■■
  0.059 [118]   |■
  0.071 [37]    |
  0.083 [16]    |
  0.095 [7]     |
  0.106 [3]     |
  0.118 [1]     |


Latency distribution:
  10%% in 0.0022 secs
  25%% in 0.0036 secs
  50%% in 0.0071 secs
  75%% in 0.0216 secs
  90%% in 0.0309 secs
  95%% in 0.0372 secs
  99%% in 0.0534 secs

Details (average, fastest, slowest):
  DNS+dialup:   0.0000 secs, 0.0000 secs, 0.0062 secs
  DNS-lookup:   0.0000 secs, 0.0000 secs, 0.0025 secs
  req write:    0.0000 secs, 0.0000 secs, 0.0023 secs
  resp wait:    0.0130 secs, 0.0005 secs, 0.1180 secs
  resp read:    0.0000 secs, 0.0000 secs, 0.0021 secs

Status code distribution:
  [302] 10000 responses
````