package api

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"Go_dev_diplom/pkg/db"
)

// updateTaskHandler обрабатывает PUT запрос на /api/task
func updateTaskHandler(w http.ResponseWriter, r *http.Request) {
	// Десериализуем JSON
	var task db.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		log.Printf("Ошибка десериализации JSON: %v", err)
		writeJSON(w, map[string]string{"error": "Неверный формат JSON"}, http.StatusBadRequest)
		return
	}

	// Проверяем, что ID указан
	if task.ID == "" {
		writeJSON(w, map[string]string{"error": "Не указан идентификатор задачи"}, http.StatusBadRequest)
		return
	}

	// Проверяем обязательное поле title
	if task.Title == "" {
		writeJSON(w, map[string]string{"error": "Не указан заголовок задачи"}, http.StatusBadRequest)
		return
	}

	// Получаем текущую дату
	now := time.Now()

	// Проверяем и корректируем дату
	fixedDate, err := validateAndFixDate(task.Date, task.Repeat, now)
	if err != nil {
		log.Printf("Ошибка валидации даты: %v", err)
		writeJSON(w, map[string]string{"error": err.Error()}, http.StatusBadRequest)
		return
	}
	task.Date = fixedDate

	// Обновляем задачу в БД
	if err := db.UpdateTask(&task); err != nil {
		log.Printf("Ошибка обновления задачи: %v", err)
		writeJSON(w, map[string]string{"error": err.Error()}, http.StatusInternalServerError)
		return
	}

	// Возвращаем успешный ответ (пустой JSON)
	writeJSON(w, map[string]interface{}{}, http.StatusOK)
}
