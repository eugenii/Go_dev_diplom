package main

import (
	"log"
	"net/http"
	"os"

	"Go_dev_diplom/pkg/api" // импортируем наш пакет api
	"Go_dev_diplom/pkg/db"  // импортируем наш пакет
)

func main() {
	// --- Настройка базы данных ---

	// Определяем путь к файлу БД
	// Приоритет: переменная окружения > значение по умолчанию
	dbFile := os.Getenv("TODO_DBFILE")
	if dbFile == "" {
		dbFile = "scheduler.db" // значение по умолчанию
	}

	// Инициализируем БД
	// Важно: db.Init создаёт глобальную переменную db.DB
	if err := db.Init(dbFile); err != nil {
		log.Fatalf("Ошибка инициализации БД: %v", err)
	}

	// Гарантированно закрываем БД при завершении программы
	defer db.Close()

	log.Printf("База данных %s готова к работе", dbFile)

	// --- Инициализация API обработчиков ---
	api.Init() // <-- ВАЖНО: регистрируем все API до запуска сервера
	log.Println("API обработчики зарегистрированы")

	// --- Настройка веб-сервера ---

	port := os.Getenv("TODO_PORT")
	if port == "" {
		port = "7540"
	}

	webDir := "./web"
	fileServer := http.FileServer(http.Dir(webDir))
	http.Handle("/", fileServer)

	log.Printf("Сервер запущен на порту %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}
