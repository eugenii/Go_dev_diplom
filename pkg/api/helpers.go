package api

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

// writeJSON записывает данные в ответ в формате JSON
func writeJSON(w http.ResponseWriter, data interface{}, status int) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("Ошибка при кодировании JSON: %v", err)
	}
}

// afterNow проверяет, что дата date больше или равна now
func afterNow(date, now time.Time) bool {
	d := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	n := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	return d.After(n) || d.Equal(n)
}

// validateAndFixDate проверяет и корректирует дату задачи
// Возвращает исправленную дату и ошибку, если дата некорректна
func validateAndFixDate(dateStr, repeat string, now time.Time) (string, error) {
	// Если дата пустая, используем сегодняшнюю
	if dateStr == "" {
		dateStr = now.Format(DateFormat)
	}

	// Парсим дату
	date, err := time.Parse(DateFormat, dateStr)
	if err != nil {
		return "", err
	}

	// Если правило повторения не указано
	if repeat == "" {
		// Если дата в прошлом, используем сегодняшнюю
		if afterNow(now, date) {
			return now.Format(DateFormat), nil
		}
		return dateStr, nil
	}

	// Правило повторения указано
	// Проверяем корректность правила и вычисляем следующую дату
	nextDate, err := NextDate(now, dateStr, repeat)
	if err != nil {
		return "", err
	}

	// Если дата в прошлом или сегодня, используем следующую
	if !afterNow(date, now) {
		return nextDate, nil
	}

	// Дата в будущем, оставляем как есть
	return dateStr, nil
}
