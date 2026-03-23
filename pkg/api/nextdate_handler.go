package api

import (
	"net/http"
	"time"
)

// NextDateHandler обрабатывает GET запросы на /api/nextdate
// Параметры запроса:
//
//	now - текущая дата (опционально, если не указана - используется текущая)
//	date - исходная дата задачи
//	repeat - правило повторения
//
// Возвращает:
//
//	в теле ответа - следующую дату в формате 20060102 или текст ошибки
func NextDateHandler(w http.ResponseWriter, r *http.Request) {
	// Разрешаем только GET запросы
	if r.Method != http.MethodGet {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	// Получаем параметры из URL
	nowParam := r.FormValue("now")
	dateParam := r.FormValue("date")
	repeatParam := r.FormValue("repeat")

	// Проверяем обязательные параметры
	if dateParam == "" {
		http.Error(w, "Не указан параметр date", http.StatusBadRequest)
		return
	}
	if repeatParam == "" {
		http.Error(w, "Не указан параметр repeat", http.StatusBadRequest)
		return
	}

	// Определяем текущую дату
	var now time.Time
	if nowParam == "" {
		now = time.Now()
	} else {
		var err error
		now, err = time.Parse(DateFormat, nowParam)
		if err != nil {
			http.Error(w, "Неверный формат параметра now", http.StatusBadRequest)
			return
		}
	}

	// Вызываем функцию вычисления следующей даты
	nextDate, err := NextDate(now, dateParam, repeatParam)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Возвращаем результат
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(nextDate))
}
