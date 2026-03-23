package api

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

// NextDate вычисляет следующую дату выполнения задачи
// Параметры:
//
//	now - время, от которого ищется ближайшая дата
//	dateStr - исходная дата в формате 20060102
//	repeat - правило повторения
//
// Возвращает:
//
//	string - следующая дата в формате 20060102
//	error - ошибка, если правило некорректно
func NextDate(now time.Time, dateStr string, repeat string) (string, error) {
	// Проверка на пустое правило
	if repeat == "" {
		return "", errors.New("правило повторения не может быть пустым")
	}

	// Парсим исходную дату
	date, err := time.Parse(DateFormat, dateStr)
	if err != nil {
		return "", errors.New("некорректный формат даты")
	}

	// Разбиваем правило на части
	parts := strings.Split(repeat, " ")
	if len(parts) == 0 {
		return "", errors.New("пустое правило повторения")
	}

	// Определяем тип правила по первому символу
	switch parts[0] {
	case "d":
		return handleDaily(now, date, parts)
	case "y":
		return handleYearly(now, date, parts)
	case "w":
		return handleWeekly(now, date, parts)
	case "m":
		return handleMonthly(now, date, parts)
	default:
		return "", errors.New("неподдерживаемый формат правила")
	}
}

// --- Базовые правила (обязательные) ---

// handleDaily обрабатывает правило "d N" - каждые N дней
func handleDaily(now time.Time, date time.Time, parts []string) (string, error) {
	// Проверяем, что указан интервал
	if len(parts) < 2 {
		return "", errors.New("для правила d требуется указать интервал в днях")
	}

	// Парсим интервал
	days, err := strconv.Atoi(parts[1])
	if err != nil {
		return "", errors.New("интервал должен быть числом")
	}

	// Проверяем ограничения по ТЗ
	if days <= 0 {
		return "", errors.New("интервал должен быть положительным числом")
	}
	if days > 400 {
		return "", errors.New("интервал не может превышать 400 дней")
	}

	// Вычисляем следующую дату
	next := date
	for {
		next = next.AddDate(0, 0, days)
		if afterNow(next, now) {
			break
		}
		// Защита от бесконечного цикла
		if next.Year() > now.Year()+100 {
			return "", errors.New("не удалось найти следующую дату")
		}
	}

	return next.Format(DateFormat), nil
}

// handleYearly обрабатывает правило "y" - ежегодно
func handleYearly(now time.Time, date time.Time, parts []string) (string, error) {
	// Проверяем, что нет лишних параметров
	if len(parts) > 1 {
		return "", errors.New("правило y не должно иметь дополнительных параметров")
	}

	// Вычисляем следующую дату
	next := date
	for {
		next = next.AddDate(1, 0, 0)
		if afterNow(next, now) {
			break
		}
		// Защита от бесконечного цикла
		if next.Year() > now.Year()+100 {
			return "", errors.New("не удалось найти следующую дату")
		}
	}

	return next.Format(DateFormat), nil
}

// --- Правила со звёздочкой (опционально, но реализуем) ---

// handleWeekly обрабатывает правило "w дни" - по дням недели
func handleWeekly(now time.Time, date time.Time, parts []string) (string, error) {
	if len(parts) < 2 {
		return "", errors.New("для правила w требуется указать дни недели")
	}

	// Парсим дни недели
	daysStr := strings.Split(parts[1], ",")
	weekDays := make(map[int]bool)

	for _, d := range daysStr {
		day, err := strconv.Atoi(strings.TrimSpace(d))
		if err != nil {
			return "", errors.New("дни недели должны быть числами")
		}
		if day < 1 || day > 7 {
			return "", errors.New("дни недели должны быть от 1 до 7")
		}
		weekDays[day] = true
	}

	// Начинаем поиск со следующего дня
	next := date.AddDate(0, 0, 1)

	// Ищем ближайший подходящий день
	for {
		// В Go Weekday: Sunday = 0, Monday = 1, ..., Saturday = 6
		// Нам нужно: Monday = 1, ..., Sunday = 7
		wd := int(next.Weekday())
		if wd == 0 {
			wd = 7 // Sunday
		}

		if weekDays[wd] && afterNow(next, now) {
			break
		}
		next = next.AddDate(0, 0, 1)

		// Защита от бесконечного цикла
		if next.Year() > now.Year()+100 {
			return "", errors.New("не удалось найти следующую дату")
		}
	}

	return next.Format(DateFormat), nil
}

// handleMonthly обрабатывает правило "m дни [месяцы]" - по дням месяца
func handleMonthly(now time.Time, date time.Time, parts []string) (string, error) {
	if len(parts) < 2 {
		return "", errors.New("для правила m требуется указать дни месяца")
	}

	// Парсим дни месяца
	daysStr := strings.Split(parts[1], ",")
	monthDays := make(map[int]bool)

	for _, d := range daysStr {
		day, err := strconv.Atoi(strings.TrimSpace(d))
		if err != nil {
			return "", errors.New("дни месяца должны быть числами")
		}
		if day < -2 || day == 0 || day > 31 {
			return "", errors.New("дни месяца должны быть от -2 до -1 или от 1 до 31")
		}
		monthDays[day] = true
	}

	// Парсим месяцы (опционально)
	months := make(map[int]bool)
	if len(parts) >= 3 {
		monthsStr := strings.Split(parts[2], ",")
		for _, m := range monthsStr {
			month, err := strconv.Atoi(strings.TrimSpace(m))
			if err != nil {
				return "", errors.New("месяцы должны быть числами")
			}
			if month < 1 || month > 12 {
				return "", errors.New("месяцы должны быть от 1 до 12")
			}
			months[month] = true
		}
	}

	// Начинаем поиск со следующего дня
	next := date.AddDate(0, 0, 1)

	for {
		// Проверяем день месяца
		day := next.Day()
		month := int(next.Month())

		// Для отрицательных дней (с конца месяца)
		lastDay := lastDayOfMonth(next)

		dayMatches := false
		for d := range monthDays {
			if d > 0 && d == day {
				dayMatches = true
				break
			}
			if d == -1 && day == lastDay {
				dayMatches = true
				break
			}
			if d == -2 && day == lastDay-1 {
				dayMatches = true
				break
			}
		}

		// Проверяем месяц (если указаны конкретные месяцы)
		monthMatches := len(months) == 0 || months[month]

		if dayMatches && monthMatches && afterNow(next, now) {
			break
		}

		next = next.AddDate(0, 0, 1)

		// Защита от бесконечного цикла
		if next.Year() > now.Year()+100 {
			return "", errors.New("не удалось найти следующую дату")
		}
	}

	return next.Format(DateFormat), nil
}

// lastDayOfMonth возвращает последний день месяца для указанной даты
func lastDayOfMonth(t time.Time) int {
	// Переходим на первый день следующего месяца и вычитаем один день
	firstOfNextMonth := time.Date(t.Year(), t.Month()+1, 1, 0, 0, 0, 0, t.Location())
	lastOfMonth := firstOfNextMonth.AddDate(0, 0, -1)
	return lastOfMonth.Day()
}

// afterNow проверяет, что дата date больше или равна now
func afterNow(date time.Time, now time.Time) bool {
	// Приводим обе даты к одному формату (без времени)
	d := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	n := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	return d.After(n) || d.Equal(n)
}
