package api

import (
	"net/http"

	"Go_dev_diplom/pkg/db"
)

// getTaskHandler обрабатывает GET запрос на /api/task?id=<id>
func getTaskHandler(w http.ResponseWriter, r *http.Request) {
	// Получаем ID из параметров запроса
	id := r.URL.Query().Get("id")
	if id == "" {
		writeJSON(w, map[string]string{"error": "Не указан идентификатор"}, http.StatusBadRequest)
		return
	}

	// Получаем задачу из БД
	task, err := db.GetTask(id)
	if err != nil {
		writeJSON(w, map[string]string{"error": err.Error()}, http.StatusNotFound)
		return
	}

	// Возвращаем задачу
	writeJSON(w, task, http.StatusOK)
}
