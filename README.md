# TODO-планировщик задач на Go

## Описание проекта

Веб-сервер для управления задачами с поддержкой повторяющихся событий. 
Позволяет создавать, редактировать, удалять задачи и отмечать их как выполненные.
Поддерживает гибкие правила повторения (ежедневно, еженедельно, ежемесячно, ежегодно).

## Выполненные задания со звёздочкой

- Поддержка всех правил повторения (`d N`, `y`, `w`, `m` с отрицательными днями и указанием месяцев)
- Поиск задач по тексту и дате (`/api/tasks?search=...`)
- Аутентификация через JWT с переменной окружения `TODO_PASSWORD`
- Docker-контейнеризация

## Технологии

- Go 1.23
- SQLite (modernc.org/sqlite)
- JWT для аутентификации (golang-jwt/jwt)
- Стандартная библиотека net/http

## Локальный запуск

### Требования
- Go 1.25 
- SQLite (опционально, для просмотра БД)

### Переменные окружения

| Переменная | Описание | Значение по умолчанию |
|------------|----------|----------------------|
| `TODO_PORT` | Порт для веб-сервера | 7540 |
| `TODO_DBFILE` | Путь к файлу БД | scheduler.db |
| `TODO_PASSWORD` | Пароль для аутентификации | (не задан, вход без пароля) |
| `TODO_SECRET` | Секрет для JWT | default-secret-key-change-me |

### Запуск

```bash
# Клонировать репозиторий
git clone <your-repo>
cd Go_dev_diplom

# Установить зависимости
go mod download

# Запустить сервер
go run main.go
```


### Запуск с аутентификацией

# Windows PowerShell
$env:TODO_PASSWORD="12345"
$env:TODO_SECRET="my-secret"
go run main.go

# Linux/Mac
export TODO_PASSWORD=12345
export TODO_SECRET=my-secret
go run main.go


### Запуск тестов

```bash
# Запуск всех тестов
go test ./tests

# Запуск конкретного теста
go test -run ^TestAddTask$ ./tests
go test -run ^TestTasks$ ./tests
go test -run ^TestDone$ ./tests
go test -run ^TestDelTask$ ./tests
```

### Настройка tests/settings.go

```go
// Для проверки всех правил повторения
FullNextDate = true

// Для проверки поиска
Search = true

// Для аутентификации (установите токен, полученный при входе)
Token = "eyJhbGciOiJIUzI1NiIs..."
```

### Docker

```bash
# Сборка образа
docker build -t todo-scheduler .

# Создать папку для данных
mkdir data

# Запустить контейнер
docker run -d \
  --name todo-scheduler \
  -p 7540:7540 \
  -v $(pwd)/data:/app/data \
  -e TODO_PASSWORD=12345 \
  -e TODO_SECRET=my-secret \
  todo-scheduler

# Проверить работу
curl http://localhost:7540/api/tasks

# Остановить контейнер и удалить
docker stop todo-scheduler
docker rm todo-scheduler

```

### API эндпойнты

|Метод	 | Эндпоинт | 	Описание |
| --- | --- | --- |
|POST	| /api/signin	| Аутентификация, получение JWT |
|POST	| /api/task	| Добавление задачи |
|GET	|/api/task?id=	|Получение задачи по ID|
|PUT	|/api/task	|Обновление задачи|
|DELETE	|/api/task?id=	|Удаление задачи|
|GET	|/api/tasks	|Список задач (с поиском)|
|POST	|/api/task/done?id=	|Отметка задачи как выполненной|
|GET	|/api/nextdate	|Вычисление следующей даты (для тестов)|
