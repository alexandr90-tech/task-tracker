# 📝 Task Tracker

Простой менеджер задач на Go. Хранит задачи в JSON-файле и предоставляет HTTP API для управления ими.

## 📁 Структура проекта

task-tracker/
├── cmd/
│ └── main.go # Точка входа: HTTP-сервер и маршруты
├── internal/
│ └── task/
│ └── task.go # Логика работы с задачами и JSON-файлом
└── tasks.json # (Создаётся автоматически) Хранилище задач


## 🚀 Запуск

go run ./cmd

📬 HTTP API
Получить список задач
GET /tasks

Создать задачу
POST /tasks
Content-Type: application/json

{
  "title": "Купить молоко"
}

Обновить статус задачи
PUT /tasks?id=1
Content-Type: application/json

{
  "done": true
}

Удалить задачу
DELETE /tasks?id=1

🛠 Используемые технологии
Go

net/http

JSON

Файловая система

📦 Установка Go
https://go.dev/doc/install

👨‍💻 Автор
alexandr90-tech