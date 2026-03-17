package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	// Определяем порт
	port := os.Getenv("TODO_PORT")
	if port == "" {
		port = "7540"
	}

	// Директория с веб-файлами
	webDir := "./web"

	// Создаем файловый сервер
	fileServer := http.FileServer(http.Dir(webDir))

	// Регистрируем обработчик для корневого пути
	http.Handle("/", fileServer)

	// Запускаем сервер
	fmt.Printf("Сервер запущен на порту %s\n", port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal(err)
	}
}
