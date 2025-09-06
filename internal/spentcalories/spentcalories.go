// Package spentcalories обрабатывает переданную информацию и
// рассчитывает потраченные калории в зависимости от вида активности - "Бег" или "Ходьба".
//
// Возвращает информационное сообщение о тренировке.
package spentcalories

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

// Константы, используемые для расчетов.
const (
	LenStep                    = 0.65 // средняя длина шага в метрах.
	MInKm                      = 1000 // количество метров в километре.
	minInH                     = 60   // количество минут в часе.
	stepLengthCoefficient      = 0.45 // коэффициент для расчета длины шага на основе роста.
	walkingCaloriesCoefficient = 0.5  // коэффициент для расчета калорий при ходьбе
)

// Константы, используемые для определения типа активности.
const (
	running = "Бег"    // тип активности "Бег".
	walking = "Ходьба" // тип активности "Ходьба".
)

// parseTraining разбирает строку с данными о тренировке.
// Ожидает строку в формате "количество_шагов,тип_активности,продолжительность" (например, "5000,Бег,30m").
// Возвращает количество шагов, тип активности, продолжительность и ошибку в случае невалидных данных.
func parseTraining(data string) (int, string, time.Duration, error) {
	parts := strings.Split(data, ",")

	if len(parts) != 3 {
		return 0, "", 0, fmt.Errorf("invalid data format: %s", data)
	}

	stepCount, activity, durationText := parts[0], parts[1], parts[2]

	count, err := strconv.Atoi(stepCount)
	if err != nil {
		return 0, activity, 0, fmt.Errorf("parsing steps failed: %w", err)
	}

	if count <= 0 {
		return 0, activity, 0, fmt.Errorf("steps must be greater than zero, got: %d", count)
	}

	duration, err := time.ParseDuration(durationText)
	if err != nil {
		return 0, activity, 0, fmt.Errorf("parsing duration failed: %w", err)
	}

	if duration <= 0 {
		return 0, activity, 0, fmt.Errorf("activity duration must be greater than zero, got: %s", duration)
	}

	return count, activity, duration, nil
}

// distance рассчитывает пройденную дистанцию в километрах.
// Принимает количество шагов и рост пользователя.
// Возвращает дистанцию в километрах.
func distance(steps int, height float64) float64 {
	if steps <= 0 || height <= 0 {
		return 0.0
	}

	return stepLengthCoefficient * height * float64(steps) / float64(MInKm)
}

// meanSpeed рассчитывает среднюю скорость передвижения в км/ч.
// Принимает количество шагов, рост пользователя в сантиметрах и продолжительность активности.
// Возвращает среднюю скорость в километрах в час.
func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	if steps <= 0 || height <= 0 || duration <= 0 {
		return 0.0
	}

	return distance(steps, height) / duration.Hours()
}

// TrainingInfo формирует информационное сообщение о тренировке.
// Принимает:
//   - data: строка с данными о тренировке в формате "количество_шагов,тип_активности,продолжительность"
//   - weight: вес пользователя в килограммах
//   - height: рост пользователя в сантиметрах
//
// Возвращает отформатированную строку с информацией о тренировке или ошибку в случае невалидных данных.
// Поддерживаемые типы активности: "Бег", "Ходьба".
func TrainingInfo(data string, weight, height float64) (string, error) {
	if weight <= 0.0 {
		return "", fmt.Errorf("weight must be greater than zero, got: %f", weight)
	}

	if height <= 0.0 {
		return "", fmt.Errorf("height must be greater than zero, got: %f", height)
	}

	steps, activity, duration, err := parseTraining(data)
	if err != nil {
		log.Println(err)
		return "", err
	}

	var calories float64
	switch activity {
	case running:
		calories, err = RunningSpentCalories(steps, weight, height, duration)
	case walking:
		calories, err = WalkingSpentCalories(steps, weight, height, duration)
	default:
		return "", fmt.Errorf("неизвестный тип тренировки")
	}

	if err != nil {
		return "", err
	}

	dist, speed := distance(steps, height), meanSpeed(steps, height, duration)

	return fmt.Sprintf("Тип тренировки: %s\nДлительность: %.2f ч.\n"+
		"Дистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f\n",
		activity, duration.Hours(), dist, speed, calories), nil
}

// RunningSpentCalories рассчитывает количество потраченных калорий при беге.
// Принимает:
//   - steps: количество шагов (должно быть > 0)
//   - weight: вес пользователя в килограммах (должен быть > 0)
//   - height: рост пользователя в сантиметрах (должен быть > 0)
//   - duration: продолжительность активности (должна быть > 0)
//
// Возвращает количество потраченных калорий или ошибку в случае невалидных входных данных.
func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 {
		return 0.0, fmt.Errorf("steps must be greater than zero, got: %d", steps)
	}

	if weight <= 0.0 {
		return 0.0, fmt.Errorf("weight must be greater than zero, got: %f", weight)
	}

	if height <= 0.0 {
		return 0.0, fmt.Errorf("height must be greater than zero, got: %f", height)
	}

	if duration <= 0 {
		return 0.0, fmt.Errorf("duration must be greater than zero, got: %s", duration)
	}

	return (weight * meanSpeed(steps, height, duration) * duration.Minutes()) / minInH, nil
}

// WalkingSpentCalories рассчитывает количество сожженных калорий при ходьбе.
// Использует функцию RunningSpentCalories и применяет дополнительный коэффициент walkingCaloriesCoefficient.
// Принимает:
//   - steps: количество шагов (должно быть > 0)
//   - weight: вес пользователя в килограммах (должен быть > 0)
//   - height: рост пользователя в сантиметрах (должен быть > 0)
//   - duration: продолжительность активности (должна быть > 0)
//
// Возвращает количество потраченных калорий или ошибку в случае невалидных входных данных.
func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	calories, err := RunningSpentCalories(steps, weight, height, duration)

	if err != nil {
		return 0.0, err
	}

	return calories * walkingCaloriesCoefficient, nil
}
