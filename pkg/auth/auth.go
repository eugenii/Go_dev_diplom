package auth

import (
	"errors"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Claims структура для JWT токена
type Claims struct {
	PasswordHash string `json:"password_hash"`
	jwt.RegisteredClaims
}

// GenerateToken создает JWT токен на основе пароля
func GenerateToken(password string) (string, error) {
	// Создаем claims с хэшем пароля
	claims := Claims{
		PasswordHash: hashPassword(password),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(8 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	// Создаем токен
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Подписываем токен секретным ключом
	secret := os.Getenv("TODO_SECRET")
	if secret == "" {
		secret = "default-secret-key-change-me"
	}
	return token.SignedString([]byte(secret))
}

// ValidateToken проверяет валидность JWT токена
func ValidateToken(tokenString string) (bool, error) {
	secret := os.Getenv("TODO_SECRET")
	if secret == "" {
		secret = "default-secret-key-change-me"
	}

	// Парсим токен
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if err != nil {
		return false, err
	}

	// Проверяем валидность токена
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		// Получаем текущий пароль из переменной окружения
		currentPassword := os.Getenv("TODO_PASSWORD")
		if currentPassword == "" {
			return true, nil // Если пароль не задан, любой токен валиден
		}
		// Проверяем соответствие хэша пароля
		return claims.PasswordHash == hashPassword(currentPassword), nil
	}

	return false, errors.New("невалидный токен")
}

// hashPassword создает простой хэш пароля (для учебных целей)
func hashPassword(password string) string {
	// Простой хэш для учебного проекта
	// В реальном проекте используйте bcrypt или другой надежный алгоритм
	return password // В учебных целях используем сам пароль
	// Для реального проекта раскомментируйте:
	// hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	// return string(hash)
}

// AuthMiddleware создает middleware для проверки аутентификации
func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Проверяем, задан ли пароль в переменной окружения
		requiredPassword := os.Getenv("TODO_PASSWORD")
		if requiredPassword == "" {
			// Пароль не задан - аутентификация не требуется
			next(w, r)
			return
		}

		// Получаем токен из куки
		cookie, err := r.Cookie("token")
		if err != nil {
			http.Error(w, "Authentication required", http.StatusUnauthorized)
			return
		}

		// Валидируем токен
		valid, err := ValidateToken(cookie.Value)
		if err != nil || !valid {
			http.Error(w, "Authentication required", http.StatusUnauthorized)
			return
		}

		next(w, r)
	}
}
