package main

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"goal-tracker/api/internal/models"
)

func TestTodoistDueStringToPeriod(t *testing.T) {
	var nilPtr *int
	assert.Equal(t, nilPtr, models.TodoistDueStringToPeriod(""))
	assert.Equal(t, models.Year, *models.TodoistDueStringToPeriod("every year"))
	assert.Equal(t, models.Quarter, *models.TodoistDueStringToPeriod("every 3 months"))
}
