package helper_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"goal-tracker/api/internal/helper"
	"goal-tracker/api/internal/models"
)

func TestAchievementsGrapher(t *testing.T) {
	totalAchievementsPerGame := map[int]int{
		1: 10, // no achievements achieved
		2: 20, // 10 achievements achieved
		3: 30, // 20 achievements achieved
	}

	grapher := helper.NewAchievementsGrapher(totalAchievementsPerGame)

	dateNow := time.Now().UTC()
	for i := 0; i < 10; i++ {
		grapher.AddPoint(dateNow.AddDate(0, 0, i), 2)
	}

	for i := 0; i < 20; i++ {
		grapher.AddPoint(dateNow.AddDate(0, 0, -1*i), 3)
	}

	dateSlice, valueSlice := grapher.ToSlices()

	assert.Equal(t, 29, len(dateSlice))
	assert.Equal(t, 29, len(valueSlice))

	assert.Equal(t, "58.33", valueSlice[28])
}

func TestGrapherCumulative(t *testing.T) {
	grapher := helper.NewGrapher[int](helper.Cumulative)

	dateNow := time.Now().UTC()
	for i := 0; i < 10; i++ {
		grapher.AddPoint(dateNow.AddDate(0, 0, i), 1)
	}

	for i := 1; i < 10; i++ {
		grapher.AddPoint(dateNow.AddDate(0, 0, -1*i), 1)
	}

	dateSlice, valueSlice := grapher.ToStringSlices()

	assert.Equal(t, 19, len(dateSlice))
	assert.Equal(t, 19, len(valueSlice))

	for i := 0; i < 19; i++ {
		assert.Equal(
			t,
			time.Now().UTC().AddDate(0, 0, i-9).Format(models.ProgressDateFormat),
			dateSlice[i],
		)
		assert.Equal(t, fmt.Sprint(i+1), valueSlice[i])
	}
}

func TestGrapherNormal(t *testing.T) {
	grapher := helper.NewGrapher[int](helper.Normal)

	dateNow := time.Now().UTC()
	for i := 0; i < 10; i++ {
		grapher.AddPoint(dateNow.AddDate(0, 0, i), i)
	}

	dateSlice, valueSlice := grapher.ToStringSlices()

	assert.Equal(t, 10, len(dateSlice))
	assert.Equal(t, 10, len(valueSlice))

	for i := 0; i < 10; i++ {
		assert.Equal(
			t,
			time.Now().UTC().AddDate(0, 0, i).Format(models.ProgressDateFormat),
			dateSlice[i],
		)
		assert.Equal(t, fmt.Sprint(i), valueSlice[i])
	}
}
