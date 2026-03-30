package api

import (
	"net/http"

	"Go_dev_diplom/pkg/db"
)

// deleteTaskHandler обрабатывает DELETE запрос на /api/task?id=<id>
func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	// Получаем ID из параметров запроса
	id := r.URL.Query().Get("id")
	if id == "" {
		writeJSON(w, map[string]string{"error": "Не указан идентификатор"}, http.StatusBadRequest)
		return
	}

	// Удаляем задачу из БД
	if err := db.DeleteTask(id); err != nil {
		writeJSON(w, map[string]string{"error": err.Error()}, http.StatusInternalServerError)
		return
	}

	// Возвращаем успешный ответ
	writeJSON(w, map[string]interface{}{}, http.StatusOK)
}
