package api

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"Go_dev_diplom/pkg/db"
)

// addTaskHandler обрабатывает POST запрос на /api/task
func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	// Проверяем метод
	if r.Method != http.MethodPost {
		writeJSON(w, map[string]string{"error": "Метод не поддерживается"}, http.StatusMethodNotAllowed)
		return
	}

	// Десериализуем JSON
	var task db.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		log.Printf("Ошибка десериализации JSON: %v", err)
		writeJSON(w, map[string]string{"error": "Неверный формат JSON"}, http.StatusBadRequest)
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

	// Добавляем задачу в БД
	id, err := db.AddTask(&task)
	if err != nil {
		log.Printf("Ошибка добавления задачи в БД: %v", err)
		writeJSON(w, map[string]string{"error": "Ошибка при сохранении задачи"}, http.StatusInternalServerError)
		return
	}

	// Возвращаем успешный ответ
	writeJSON(w, map[string]interface{}{
		"id": id,
	}, http.StatusOK)
}
