package daysteps

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/Yandex-Practicum/tracker/internal/spentcalories"
)

const (
	// Длина одного шага в метрах
	stepLength = 0.65
	// Количество метров в одном километре
	mInKm = 1000
)

func parsePackage(data string) (int, time.Duration, error) {

	var (
		stepCount int
		duration  time.Duration
		err       error
		params    []string
	)

	params = strings.Split(data, ",")

	if paramsLen := len(params); paramsLen != 2 {
		return 0, 0, fmt.Errorf("expected input of 2 values, got %d", paramsLen)
	}

	stepCount, err = strconv.Atoi(params[0])

	if err != nil {
		return 0, 0, fmt.Errorf("error while converting step count to int: %w", err)
	}

	if stepCount <= 0 {
		return 0, 0, fmt.Errorf("expected strongly positive step count")
	}

	duration, err = time.ParseDuration(params[1])

	if err != nil {
		return 0, 0, fmt.Errorf("error while converting activity duration to time.Duration: %w", err)
	}

	if duration <= 0 {
		return 0, 0, fmt.Errorf("expected strongly positive duration: %w", err)
	}

	return stepCount, duration, nil
}

func DayActionInfo(data string, weight, height float64) string {
	var (
		stepCount int
		duration  time.Duration
		err       error
		calories  float64
	)

	stepCount, duration, err = parsePackage(data)

	if err != nil {
		log.Println("error while parsing data parameter")
		return ""
	}

	distance := float64(stepCount) * stepLength
	kilometersWent := distance / mInKm

	calories, err = spentcalories.WalkingSpentCalories(stepCount, weight, height, duration)

	if err != nil {
		log.Println("error while parsing data parameter")
		return ""
	}

	result := fmt.Sprintf(
		"Количество шагов: %d.\nДистанция составила %.2f км.\nВы сожгли %.2f ккал.\n",
		stepCount,
		kilometersWent,
		calories,
	)

	return result
}
