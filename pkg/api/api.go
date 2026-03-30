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

	// Обработчик для /api/tasks (GET)
	http.HandleFunc("/api/tasks", tasksHandler)
	log.Println("  /api/tasks зарегистрирован")

	// Обработчик для /api/task/done (POST)
	http.HandleFunc("/api/task/done", doneTaskHandler)
	log.Println("  /api/task/done зарегистрирован")
}

// taskHandler диспетчеризует запросы к /api/task в зависимости от метода
func taskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		addTaskHandler(w, r)
	case http.MethodGet:
		getTaskHandler(w, r)
	case http.MethodPut:
		updateTaskHandler(w, r)
	case http.MethodDelete:
		deleteTaskHandler(w, r)
		// Будет реализовано в шаге 7
		// w.Header().Set("Content-Type", "text/plain")
		// w.WriteHeader(http.StatusMethodNotAllowed)
		// w.Write([]byte("Метод DELETE будет реализован позже"))
	default:
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Метод не поддерживается"))
	}
}
