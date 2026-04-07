package api

import (
	"encoding/json"
	"net/http"
	"os"

	"Go_dev_diplom/pkg/auth"
)

type SigninRequest struct {
	Password string `json:"password"`
}

type SigninResponse struct {
	Token string `json:"token,omitempty"`
	Error string `json:"error,omitempty"`
}

// signinHandler обрабатывает POST запрос на /api/signin
func signinHandler(w http.ResponseWriter, r *http.Request) {
	// Проверяем метод
	if r.Method != http.MethodPost {
		writeJSON(w, SigninResponse{Error: "Метод не поддерживается"}, http.StatusMethodNotAllowed)
		return
	}

	// Декодируем JSON
	var req SigninRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, SigninResponse{Error: "Неверный формат JSON"}, http.StatusBadRequest)
		return
	}

	// Получаем ожидаемый пароль из переменной окружения
	expectedPassword := os.Getenv("TODO_PASSWORD")

	// Если пароль не задан в переменной окружения, любой пароль подходит
	if expectedPassword == "" {
		// Генерируем токен (пароль не требуется)
		token, err := auth.GenerateToken("")
		if err != nil {
			writeJSON(w, SigninResponse{Error: "Ошибка генерации токена"}, http.StatusInternalServerError)
			return
		}
		writeJSON(w, SigninResponse{Token: token}, http.StatusOK)
		return
	}

	// Проверяем пароль
	if req.Password != expectedPassword {
		writeJSON(w, SigninResponse{Error: "Неверный пароль"}, http.StatusUnauthorized)
		return
	}

	// Генерируем JWT токен
	token, err := auth.GenerateToken(req.Password)
	if err != nil {
		writeJSON(w, SigninResponse{Error: "Ошибка генерации токена"}, http.StatusInternalServerError)
		return
	}

	// Возвращаем токен
	writeJSON(w, SigninResponse{Token: token}, http.StatusOK)
}
