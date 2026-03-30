package db

import (
	"database/sql"
	"errors"
	"fmt"
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

// Tasks возвращает список задач, отсортированных по дате
// limit - максимальное количество возвращаемых задач
// search - строка поиска (опционально, может быть пустой)
func Tasks(limit int, search string) ([]*Task, error) {
	var rows *sql.Rows
	var err error

	// Базовый запрос
	query := `SELECT id, date, title, comment, repeat 
	          FROM scheduler 
	          ORDER BY date LIMIT ?`

	// Если есть строка поиска
	if search != "" {
		// Проверяем, является ли search датой в формате DD.MM.YYYY
		if isDateSearch(search) {
			// Преобразуем дату из 02.01.2006 в 20060102
			dateFormatted, err := formatSearchDate(search)
			if err != nil {
				return nil, err
			}
			query = `SELECT id, date, title, comment, repeat 
			         FROM scheduler 
			         WHERE date = ? 
			         ORDER BY date LIMIT ?`
			rows, err = DB.Query(query, dateFormatted, limit)
		} else {
			// Поиск по подстроке в title и comment
			searchPattern := "%" + search + "%"
			query = `SELECT id, date, title, comment, repeat 
			         FROM scheduler 
			         WHERE title LIKE ? OR comment LIKE ? 
			         ORDER BY date LIMIT ?`
			rows, err = DB.Query(query, searchPattern, searchPattern, limit)
		}
	} else {
		// Без поиска
		rows, err = DB.Query(query, limit)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Создаем слайс для результатов
	tasks := make([]*Task, 0)

	// Проходим по всем строкам результата
	for rows.Next() {
		var task Task
		err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, &task)
	}

	// Проверяем ошибки после итерации
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}

// isDateSearch проверяет, является ли строка поиска датой в формате DD.MM.YYYY
func isDateSearch(search string) bool {
	// Простая проверка формата: 2 цифры, точка, 2 цифры, точка, 4 цифры
	if len(search) != 10 {
		return false
	}
	if search[2] != '.' || search[5] != '.' {
		return false
	}
	// Проверяем, что все символы кроме точек - цифры
	for i, ch := range search {
		if i == 2 || i == 5 {
			continue
		}
		if ch < '0' || ch > '9' {
			return false
		}
	}
	return true
}

// formatSearchDate преобразует дату из формата DD.MM.YYYY в YYYYMMDD
func formatSearchDate(search string) (string, error) {
	// Парсим дату в формате DD.MM.YYYY
	// Разбиваем по точкам
	parts := make([]string, 0)
	lastIdx := 0
	for i, ch := range search {
		if ch == '.' {
			parts = append(parts, search[lastIdx:i])
			lastIdx = i + 1
		}
	}
	parts = append(parts, search[lastIdx:])

	if len(parts) != 3 {
		return "", errors.New("неверный формат даты")
	}

	day := parts[0]
	month := parts[1]
	year := parts[2]

	// Формируем строку в формате YYYYMMDD
	return year + month + day, nil
}

// GetTask возвращает задачу по её ID
func GetTask(id string) (*Task, error) {
	var task Task

	query := `SELECT id, date, title, comment, repeat 
              FROM scheduler 
              WHERE id = ?`

	err := DB.QueryRow(query, id).Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("задача не найдена")
		}
		return nil, err
	}

	return &task, nil
}

// UpdateTask обновляет существующую задачу
func UpdateTask(task *Task) error {
	query := `UPDATE scheduler 
              SET date = ?, title = ?, comment = ?, repeat = ? 
              WHERE id = ?`

	result, err := DB.Exec(query, task.Date, task.Title, task.Comment, task.Repeat, task.ID)
	if err != nil {
		log.Printf("Ошибка при обновлении задачи: %v", err)
		return err
	}

	// Проверяем, была ли обновлена запись
	count, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if count == 0 {
		return fmt.Errorf("задача с ID %s не найдена", task.ID)
	}

	log.Printf("Задача с ID %s успешно обновлена", task.ID)
	return nil
}

// DeleteTask удаляет задачу по ID
func DeleteTask(id string) error {
	query := `DELETE FROM scheduler WHERE id = ?`

	result, err := DB.Exec(query, id)
	if err != nil {
		log.Printf("Ошибка при удалении задачи: %v", err)
		return err
	}

	// Проверяем, была ли удалена запись
	count, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if count == 0 {
		return fmt.Errorf("задача с ID %s не найдена", id)
	}

	log.Printf("Задача с ID %s успешно удалена", id)
	return nil
}

// UpdateTaskDate обновляет только дату задачи (для повторяющихся задач)
func UpdateTaskDate(id string, newDate string) error {
	query := `UPDATE scheduler SET date = ? WHERE id = ?`

	result, err := DB.Exec(query, newDate, id)
	if err != nil {
		log.Printf("Ошибка при обновлении даты задачи: %v", err)
		return err
	}

	// Проверяем, была ли обновлена запись
	count, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if count == 0 {
		return fmt.Errorf("задача с ID %s не найдена", id)
	}

	log.Printf("Дата задачи с ID %s обновлена на %s", id, newDate)
	return nil
}
