package api

import (
	"log"
	"net/http"
)

const DateFormat = "20060102"

// Init регистрирует все обработчики API
func Init() {
	log.Println("Регистрируем обработчики API...")

	// Обработчик для /api/nextdate
	http.HandleFunc("/api/nextdate", NextDateHandler)
	log.Println("  /api/nextdate зарегистрирован")

	// Обработчик для /api/task (все методы)
	http.HandleFunc("/api/task", taskHandler)
	log.Println("  /api/task зарегистрирован")
}

// taskHandler диспетчеризует запросы к /api/task в зависимости от метода
func taskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		addTaskHandler(w, r)
	// Здесь позже добавим:
	// case http.MethodGet:
	//     getTaskHandler(w, r)
	// case http.MethodPut:
	//     updateTaskHandler(w, r)
	// case http.MethodDelete:
	//     deleteTaskHandler(w, r)
	default:
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Метод не поддерживается"))
	}
}
