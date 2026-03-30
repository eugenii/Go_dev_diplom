package api

import (
	"log"
	"net/http"
	"time"

	"Go_dev_diplom/pkg/db"
)

// doneTaskHandler обрабатывает POST запрос на /api/task/done?id=<id>
func doneTaskHandler(w http.ResponseWriter, r *http.Request) {
	// Проверяем метод
	if r.Method != http.MethodPost {
		writeJSON(w, map[string]string{"error": "Метод не поддерживается"}, http.StatusMethodNotAllowed)
		return
	}

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

	// Проверяем, есть ли правило повторения
	if task.Repeat == "" {
		// Одноразовая задача - удаляем
		if err := db.DeleteTask(id); err != nil {
			writeJSON(w, map[string]string{"error": err.Error()}, http.StatusInternalServerError)
			return
		}
		log.Printf("Задача с ID %s выполнена и удалена (одноразовая)", id)
	} else {
		// Повторяющаяся задача - вычисляем следующую дату
		now := time.Now()

		// Используем текущую дату задачи как исходную
		nextDate, err := NextDate(now, task.Date, task.Repeat)
		if err != nil {
			log.Printf("Ошибка вычисления следующей даты: %v", err)
			writeJSON(w, map[string]string{"error": err.Error()}, http.StatusBadRequest)
			return
		}

		// Обновляем дату задачи
		if err := db.UpdateTaskDate(id, nextDate); err != nil {
			writeJSON(w, map[string]string{"error": err.Error()}, http.StatusInternalServerError)
			return
		}
		log.Printf("Задача с ID %s выполнена, следующая дата: %s", id, nextDate)
	}

	// Возвращаем успешный ответ
	writeJSON(w, map[string]interface{}{}, http.StatusOK)
}
