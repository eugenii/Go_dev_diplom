package api

import (
	"net/http"
	"strconv"

	"Go_dev_diplom/pkg/db"
)

// TasksResp структура ответа со списком задач
type TasksResp struct {
	Tasks []*db.Task `json:"tasks"`
}

// tasksHandler обрабатывает GET запрос на /api/tasks
func tasksHandler(w http.ResponseWriter, r *http.Request) {
	// Проверяем метод
	if r.Method != http.MethodGet {
		writeJSON(w, map[string]string{"error": "Метод не поддерживается"}, http.StatusMethodNotAllowed)
		return
	}

	// Получаем параметры запроса
	limit := 50 // значение по умолчанию
	limitStr := r.URL.Query().Get("limit")
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	// Получаем строку поиска
	search := r.URL.Query().Get("search")

	// Получаем задачи из БД
	tasks, err := db.Tasks(limit, search)
	if err != nil {
		writeJSON(w, map[string]string{"error": err.Error()}, http.StatusInternalServerError)
		return
	}

	// Если tasks == nil, создаем пустой слайс
	if tasks == nil {
		tasks = make([]*db.Task, 0)
	}

	// Возвращаем успешный ответ
	writeJSON(w, TasksResp{
		Tasks: tasks,
	}, http.StatusOK)
}
