package spentcalories

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

// Основные константы, необходимые для расчетов.
const (
	lenStep                    = 0.65 // средняя длина шага.
	mInKm                      = 1000 // количество метров в километре.
	minInH                     = 60   // количество минут в часе.
	stepLengthCoefficient      = 0.45 // коэффициент для расчета длины шага на основе роста.
	walkingCaloriesCoefficient = 0.5  // коэффициент для расчета калорий при ходьбе
)

func parseTraining(data string) (int, string, time.Duration, error) {
	var (
		stepCount        int
		activityType     string
		activityDuration time.Duration
		params           []string
		err              error
	)

	params = strings.Split(data, ",")

	if paramsLen := len(params); paramsLen != 3 {
		return 0, "", 0, fmt.Errorf("expected input of 3 values, got %d", paramsLen)
	}

	stepCount, err = strconv.Atoi(params[0])

	if err != nil {
		return 0, "", 0, fmt.Errorf("error while converting step count to int: %w", err)
	}

	if stepCount <= 0 {
		return 0, "", 0, fmt.Errorf("expected strongly positive step count")
	}

	activityType = params[1]

	activityDuration, err = time.ParseDuration(params[2])

	if err != nil {
		return 0, "", 0, fmt.Errorf("error while converting activity duration to time.Duration: %w", err)
	}

	if activityDuration <= 0 {
		return 0, "", 0, fmt.Errorf("only positive duration expected")
	}

	return stepCount, activityType, activityDuration, nil
}

func distance(steps int, height float64) float64 {
	stepLength := height * stepLengthCoefficient
	stepDistance := stepLength * float64(steps)
	totalDistanceKm := stepDistance / mInKm
	return totalDistanceKm
}

func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	if duration <= 0 {
		return 0
	}

	totalDistanceKm := distance(steps, height)

	return totalDistanceKm / duration.Hours()
}

func TrainingInfo(data string, weight, height float64) (string, error) {
	var (
		stepCount        int
		activityType     string
		activityDuration time.Duration
		err              error
		calories         float64
	)

	stepCount, activityType, activityDuration, err = parseTraining(data)

	if err != nil {
		log.Println(err)
		return "", err
	}

	meanSpeed := meanSpeed(stepCount, height, activityDuration)
	distance := distance(stepCount, height)

	switch activityType {
	case "Бег":
		calories, err = RunningSpentCalories(stepCount, weight, height, activityDuration)
	case "Ходьба":
		calories, err = WalkingSpentCalories(stepCount, weight, height, activityDuration)
	default:
		err = fmt.Errorf("неизвестный тип тренировки")
	}

	if err != nil {
		log.Println(err)
		return "", fmt.Errorf("error while calculating calories: %w", err)
	}

	result := fmt.Sprintf(
		"Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f\n",
		activityType,
		activityDuration.Hours(),
		distance,
		meanSpeed,
		calories,
	)

	return result, nil
}

func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 {
		return 0, fmt.Errorf("steps should be positive")
	}

	if weight <= 0 {
		return 0, fmt.Errorf("weight should be positive")
	}

	if height <= 0 {
		return 0, fmt.Errorf("height should be positive")
	}

	if duration <= 0 {
		return 0, fmt.Errorf("duration should be positive")
	}

	meanSpeed := meanSpeed(steps, height, duration)
	durationMinutes := duration.Minutes()

	return (weight * meanSpeed * durationMinutes) / minInH, nil
}

func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 {
		return 0, fmt.Errorf("steps should be positive")
	}

	if weight <= 0 {
		return 0, fmt.Errorf("weight should be positive")
	}

	if height <= 0 {
		return 0, fmt.Errorf("height should be positive")
	}

	if duration <= 0 {
		return 0, fmt.Errorf("duration should be positive")
	}

	meanSpeed := meanSpeed(steps, height, duration)
	durationMinutes := duration.Minutes()

	return ((weight * meanSpeed * durationMinutes) / minInH) * walkingCaloriesCoefficient, nil
}
