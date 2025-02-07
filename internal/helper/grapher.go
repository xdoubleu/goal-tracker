package helper

import (
	"fmt"
	"slices"
	"strconv"
	"time"

	"goal-tracker/api/internal/models"
)

type GraphType int
type Numeric interface {
	int | int64 | float64
}

const (
	Normal     GraphType = iota
	Cumulative GraphType = iota
)

type Grapher[T Numeric] struct {
	graphType   GraphType
	dateStrings []string
	values      []T
}

func NewGrapher[T Numeric](graphType GraphType) *Grapher[T] {
	return &Grapher[T]{
		graphType:   graphType,
		dateStrings: []string{},
		values:      []T{},
	}
}

func (grapher *Grapher[T]) AddPoint(date time.Time, value T) {
	dateStr := date.Format(models.ProgressDateFormat)
	dateIndex := slices.Index(grapher.dateStrings, dateStr)

	if dateIndex == -1 {
		grapher.addDays(dateStr)
		dateIndex = slices.Index(grapher.dateStrings, dateStr)
	}

	grapher.updateDays(dateIndex, value)
}

func (grapher *Grapher[T]) addDays(dateStr string) {
	if len(grapher.dateStrings) == 0 {
		grapher.dateStrings = append(grapher.dateStrings, dateStr)
		grapher.values = append(
			grapher.values,
			*new(T),
		)
		return
	}

	dateDay, _ := time.Parse(models.ProgressDateFormat, dateStr)
	smallestDate, _ := time.Parse(models.ProgressDateFormat, grapher.dateStrings[0])
	largestDate, _ := time.Parse(
		models.ProgressDateFormat,
		grapher.dateStrings[len(grapher.dateStrings)-1],
	)

	i := smallestDate
	for i.After(dateDay) {
		i = i.AddDate(0, 0, -1)

		grapher.dateStrings = append(
			[]string{i.Format(models.ProgressDateFormat)},
			grapher.dateStrings...)
		grapher.values = append(
			[]T{*new(T)},
			grapher.values...)
	}

	i = largestDate
	for i.Before(dateDay) {
		i = i.AddDate(0, 0, 1)

		grapher.dateStrings = append(
			grapher.dateStrings,
			i.Format(models.ProgressDateFormat),
		)

		indexOfI := slices.Index(
			grapher.dateStrings,
			i.Format(models.ProgressDateFormat),
		)

		grapher.values = append(
			grapher.values,
			grapher.values[indexOfI-1],
		)
	}
}

func (grapher *Grapher[T]) updateDays(dateIndex int, value T) {
	for i := dateIndex; i < len(grapher.dateStrings); i++ {
		switch grapher.graphType {
		case Normal:
			grapher.values[i] = value
		case Cumulative:
			grapher.values[i] += value
		}
	}
}

func (grapher Grapher[T]) ToStringSlices() ([]string, []string) {
	strValues := []string{}

	for _, value := range grapher.values {
		strValue := ""
		switch v := any(value).(type) {
		case int:
			strValue = strconv.Itoa(v)
		case int64:
			strValue = strconv.Itoa(int(v))
		case float64:
			strValue = fmt.Sprintf("%.2f", v)
		}

		strValues = append(strValues, strValue)
	}

	return grapher.dateStrings, strValues
}
