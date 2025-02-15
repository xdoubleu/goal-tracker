package main

import (
	"context"
	"testing"

	"github.com/XDoubleU/essentia/pkg/logging"
	"github.com/stretchr/testify/assert"

	"goal-tracker/api/internal/dtos"
	"goal-tracker/api/internal/jobs"
	"goal-tracker/api/internal/models"
)

func TestGoodreadsJob(t *testing.T) {
	err := testApp.services.Goals.ImportGoalsFromTodoist(context.Background(), userID)
	assert.Nil(t, err)

	val := int64(12)
	val1 := "tag1"
	err = testApp.services.Goals.LinkGoal(
		context.Background(),
		goalID,
		userID,
		&dtos.LinkGoalDto{
			TypeID:      models.BooksFromSpecificTag.ID,
			TargetValue: &val,
			Tag:         &val1,
		},
	)
	assert.Nil(t, err)

	err = testApp.services.Goals.LinkGoal(
		context.Background(),
		goal2ID,
		userID,
		//nolint:exhaustruct //other fields are optional
		&dtos.LinkGoalDto{
			TypeID: models.SpecificBooks.ID,
		},
	)
	assert.Nil(t, err)

	_, err = testApp.services.Goals.SaveListItem(
		context.Background(),
		1,
		userID,
		goal2ID,
		"2",
		false,
	)
	assert.Nil(t, err)

	job := jobs.NewGoodreadsJob(
		testApp.services.Auth,
		testApp.services.Goodreads,
		testApp.services.Goals,
	)
	job.ID()
	job.RunEvery()

	err = job.Run(context.Background(), logging.NewNopLogger())
	assert.Nil(t, err)
}

func TestSteamJob(t *testing.T) {
	job := jobs.NewSteamJob(
		testApp.services.Auth,
		testApp.services.Steam,
		testApp.services.Goals,
	)
	job.ID()
	job.RunEvery()

	err := job.Run(context.Background(), logging.NewNopLogger())
	assert.Nil(t, err)
}

func TestTodoistJob(t *testing.T) {
	job := jobs.NewTodoistJob(testApp.services.Auth, testApp.services.Goals)
	job.ID()
	job.RunEvery()

	err := job.Run(context.Background(), logging.NewNopLogger())
	assert.Nil(t, err)
}
