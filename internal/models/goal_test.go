package models_test

import (
	"goal-tracker/api/internal/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTodoistDueStringToPeriod(t *testing.T) {
	assert.Equal(t, nil, models.TodoistDueStringToPeriod(nil))

	everyYear := "every year"
	assert.Equal(t, models.Year, models.TodoistDueStringToPeriod(&everyYear))
}
