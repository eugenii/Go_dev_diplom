// Пакет db предоставляет функции для работы с базой данных SQLite
package db

import (
	"database/sql"
	"log"
	"os"

	_ "modernc.org/sqlite" // драйвер SQLite
)

// Глобальная переменная для хранения подключения к БД
// Используется во всех операциях с базой данных
var DB *sql.DB

// Schema содержит SQL команды для создания таблицы и индекса
// Обратные кавычки позволяют писать многострочные строки
const schema = `
CREATE TABLE IF NOT EXISTS scheduler (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    date CHAR(8) NOT NULL DEFAULT "",
    title VARCHAR(255) NOT NULL,
    comment TEXT,
    repeat VARCHAR(128)
);
CREATE INDEX IF NOT EXISTS idx_scheduler_date ON scheduler(date);`

// Init инициализирует подключение к базе данных
// Принимает путь к файлу БД, возвращает ошибку если что-то пошло не так
func Init(dbFile string) error {
	var err error

	// Проверяем, существует ли файл базы данных
	_, err = os.Stat(dbFile)
	needCreate := os.IsNotExist(err) // true если файла нет

	// Открываем базу данных
	// Важно: драйвер указывается как "sqlite", хотя пакет называется modernc.org/sqlite
	DB, err = sql.Open("sqlite", dbFile)
	if err != nil {
		return err
	}

	// Проверяем, что подключение работает
	// DB.Ping() отправляет тестовый запрос к БД
	if err = DB.Ping(); err != nil {
		return err
	}

	log.Println("Подключение к БД успешно установлено")

	// Если файл только что создан, нужно инициализировать таблицу
	if needCreate {
		log.Println("Файл БД не найден, создаём таблицу и индекс")
		_, err = DB.Exec(schema)
		if err != nil {
			return err
		}
		log.Println("Таблица scheduler успешно создана")
	}

	return nil
}

// Close закрывает подключение к базе данных
// Важно вызывать при завершении работы сервера (defer)
func Close() error {
	if DB != nil {
		return DB.Close()
	}
	return nil
}
