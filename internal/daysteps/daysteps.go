// Package daysteps отвечает за учёт активности в течение дня.
//
// Он собирает переданную информацию в виде строк, парсит их и выводит
// информацию о количестве шагов, дистанции и потраченных калорий.
package daysteps

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/Kuguchev/fitness-tracker/internal/spentcalories"
)

// parsePackage разбирает строку с данными о шагах и продолжительности ходьбы.
// Принимает строку в формате "количество_шагов,продолжительность" (например, "5000,30m").
// Возвращает количество шагов, продолжительность ходьбы и ошибку в случае невалидных данных.
// Ошибка возвращается, если:
// - неверный формат строки
// - количество шагов не является положительным числом
// - продолжительность не может быть распарсена или не является положительной
func parsePackage(data string) (int, time.Duration, error) {
	parts := strings.Split(data, ",")
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("invalid data format, expected 'steps,duration', got: %s", data)
	}

	stepCount, durationText := parts[0], parts[1]
	count, err := strconv.Atoi(stepCount)
	if err != nil {
		return 0, 0, fmt.Errorf("parsing steps failed: %w", err)
	}

	if count <= 0 {
		return 0, 0, fmt.Errorf("steps must be greater than zero, got: %d", count)
	}

	duration, err := time.ParseDuration(durationText)
	if err != nil {
		return 0, 0, fmt.Errorf("parsing duration failed: %w", err)
	}

	if duration <= 0 {
		return 0, 0, fmt.Errorf("walk duration must be greater than zero, got: %s", duration)
	}

	return count, duration, nil
}

// DayActionInfo формирует информационное сообщение о дневной активности на основе пройденных шагов.
// Принимает:
//   - data: строка в формате "количество_шагов,продолжительность" (например, "5000,30m")
//   - weight: вес пользователя в килограммах (должен быть > 0)
//   - height: рост пользователя в сантиметрах (должен быть > 0)
//
// Возвращает отформатированную строку с информацией о количестве шагов, пройденной дистанции
// и потраченных калориях. В случае ошибки возвращает пустую строку.
func DayActionInfo(data string, weight, height float64) string {
	if weight <= 0.0 || height <= 0.0 {
		return ""
	}

	steps, duration, err := parsePackage(data)
	if err != nil {
		log.Println(err)
		return ""
	}

	dist := float64(steps) * spentcalories.LenStep / spentcalories.MInKm
	calories, err := spentcalories.WalkingSpentCalories(steps, weight, height, duration)
	if err != nil {
		return ""
	}

	return fmt.Sprintf("Количество шагов: %d.\nДистанция составила %.2f км.\nВы сожгли %.2f ккал.\n",
		steps, dist, calories)
}
