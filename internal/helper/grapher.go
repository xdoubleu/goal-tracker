package helper

import (
	"goal-tracker/api/internal/models"
	"math"
	"slices"
	"time"
)

type Grapher struct {
	dateStrings               []string
	achievementsPerGamePerDay []map[int]int
	totalAchievementsPerGame  map[int]int
}

func NewGrapher(totalAchievementsPerGame map[int]int) Grapher {
	return Grapher{
		totalAchievementsPerGame: totalAchievementsPerGame,
	}
}

func (grapher *Grapher) AddPoint(date time.Time, gameID int) {
	dateStr := date.Format(models.ProgressDateFormat)
	dateIndex := slices.Index(grapher.dateStrings, dateStr)

	if dateIndex == -1 {
		grapher.addDays(dateStr)
		dateIndex = len(grapher.dateStrings) - 1
	}

	grapher.updateDays(dateIndex, gameID)
}

func (grapher Grapher) GetFirstDay() time.Time {
	date, _ := time.Parse(models.ProgressDateFormat, grapher.dateStrings[0])
	return date
}

func (grapher Grapher) GetLastDay() time.Time {
	date, _ := time.Parse(models.ProgressDateFormat, grapher.dateStrings[len(grapher.dateStrings)-1])
	return date
}

func (grapher *Grapher) addDays(dateStr string) {
	if len(grapher.dateStrings) == 0 {
		grapher.dateStrings = append(grapher.dateStrings, dateStr)
		grapher.achievementsPerGamePerDay = append(grapher.achievementsPerGamePerDay, map[int]int{})
		return
	}

	dateDay, _ := time.Parse(models.ProgressDateFormat, dateStr)
	smallestDate, _ := time.Parse(models.ProgressDateFormat, grapher.dateStrings[0])
	largestDate, _ := time.Parse(models.ProgressDateFormat, grapher.dateStrings[len(grapher.dateStrings)-1])

	if dateDay.Before(smallestDate) {
		i := smallestDate
		for i.After(dateDay) {
			i = i.AddDate(0, 0, -1)

			grapher.dateStrings = append([]string{i.Format(models.ProgressDateFormat)}, grapher.dateStrings...)
			grapher.achievementsPerGamePerDay = append([]map[int]int{{}}, grapher.achievementsPerGamePerDay...)
		}
	}

	if dateDay.After(largestDate) {
		i := largestDate
		for i.Before(dateDay) {
			i = i.AddDate(0, 0, 1)

			grapher.dateStrings = append(grapher.dateStrings, i.Format(models.ProgressDateFormat))
			grapher.achievementsPerGamePerDay = append(grapher.achievementsPerGamePerDay, copyMap(grapher.achievementsPerGamePerDay[len(grapher.achievementsPerGamePerDay)-1]))
		}
	}
}

func copyMap(original map[int]int) map[int]int {
	target := map[int]int{}

	for k, v := range original {
		target[k] = v
	}

	return target
}

func (grapher *Grapher) updateDays(dateIndex int, gameID int) {
	for i := dateIndex; i < len(grapher.dateStrings); i++ {
		grapher.achievementsPerGamePerDay[i][gameID]++
	}
}

func (grapher Grapher) ToSlices() ([]string, []int64) {
	percentages := []int64{}

	droppedCount := 0
	for i, achievementsPerGame := range grapher.achievementsPerGamePerDay {
		games := 0
		totalPercentageDay := 0.0

		for gameID, achievements := range achievementsPerGame {
			games++

			totalAchievements := grapher.totalAchievementsPerGame[int(gameID)]
			totalPercentageDay += calculateCompletionRate(achievements, totalAchievements)
		}

		avgCompletionRate := calculateAvgCompletionRate(totalPercentageDay, games)
		if avgCompletionRate == 0 {
			dateStringsIndex := i - droppedCount
			grapher.dateStrings = append(grapher.dateStrings[:dateStringsIndex], grapher.dateStrings[dateStringsIndex+1:]...)
			droppedCount++
			continue
		}

		percentages = append(percentages, avgCompletionRate)
	}

	return grapher.dateStrings, percentages
}

func calculateCompletionRate(achieved int, total int) float64 {
	return float64(achieved) / float64(total)
}

func calculateAvgCompletionRate(percentageSum float64, totalGames int) int64 {
	return int64(math.Floor(percentageSum / float64(totalGames) * 100.0))
}
