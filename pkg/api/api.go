package api

import (
	"Go_dev_diplom/pkg/auth"
	"log"
	"net/http"
)

const DateFormat = "20060102"

// Init регистрирует все обработчики API
func Init() {
	log.Println("Регистрируем обработчики API...")

	// Обработчик для /api/signin (без аутентификации)
	http.HandleFunc("/api/signin", signinHandler)
	log.Println("  /api/signin зарегистрирован")

	// /api/nextdate - без аутентификации (используется тестами и для вычислений)
	http.HandleFunc("/api/nextdate", NextDateHandler)
	log.Println("  /api/nextdate зарегистрирован")

	// Обработчик для /api/task (с аутентификацией)
	http.HandleFunc("/api/task", auth.AuthMiddleware(taskHandler))
	log.Println("  /api/task зарегистрирован")

	// Обработчик для /api/tasks (с аутентификацией)
	http.HandleFunc("/api/tasks", auth.AuthMiddleware(tasksHandler))
	log.Println("  /api/tasks зарегистрирован")

	// Обработчик для /api/task/done (с аутентификацией)
	http.HandleFunc("/api/task/done", auth.AuthMiddleware(doneTaskHandler))
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
