// Пакет api содержит обработчики HTTP запросов
package api

import (
	"net/http"
)

// DateFormat константа для формата даты (2006-01-02 в Go - это эталон)
const DateFormat = "20060102"

// Init регистрирует все обработчики API
// Вызывается из main.go до запуска сервера
func Init() {
	http.HandleFunc("/api/nextdate", NextDateHandler)
	// здесь будут регистрироваться остальные обработчики
}
