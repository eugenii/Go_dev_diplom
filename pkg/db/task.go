package db

import (
	"log"
)

// Task представляет структуру задачи для работы с БД
type Task struct {
	ID      string `json:"id" db:"id"`
	Date    string `json:"date" db:"date"`
	Title   string `json:"title" db:"title"`
	Comment string `json:"comment" db:"comment"`
	Repeat  string `json:"repeat" db:"repeat"`
}

// AddTask добавляет новую задачу в базу данных
// Возвращает ID созданной задачи и ошибку
func AddTask(task *Task) (int64, error) {
	// Подготавливаем SQL запрос
	query := `INSERT INTO scheduler (date, title, comment, repeat) 
	          VALUES (?, ?, ?, ?)`

	// Выполняем запрос
	result, err := DB.Exec(query, task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		log.Printf("Ошибка при добавлении задачи: %v", err)
		return 0, err
	}

	// Получаем ID созданной записи
	id, err := result.LastInsertId()
	if err != nil {
		log.Printf("Ошибка при получении ID: %v", err)
		return 0, err
	}

	log.Printf("Задача добавлена с ID: %d", id)
	return id, nil
}
