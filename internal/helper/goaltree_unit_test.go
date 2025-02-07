//nolint:exhaustruct //other fields are optional
package helper_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"goal-tracker/api/internal/helper"
	"goal-tracker/api/internal/models"
)

func TestGoalTree(t *testing.T) {
	goalTree := helper.NewGoalTree()

	// Add parents
	goalTree.TryAdd(models.Goal{
		ID:       "1",
		ParentID: nil,
		Order:    1,
	})

	goalTree.TryAdd(models.Goal{
		ID:       "2",
		ParentID: nil,
		Order:    2,
	})

	goalTree.TryAdd(models.Goal{
		ID:       "3",
		ParentID: nil,
		Order:    3,
	})

	// Add children
	id1 := "1"
	id2 := "2"
	id3 := "3"

	goalTree.TryAdd(models.Goal{
		ID:       "11",
		ParentID: &id1,
		Order:    1,
	})

	goalTree.TryAdd(models.Goal{
		ID:       "21",
		ParentID: &id2,
		Order:    1,
	})

	goalTree.TryAdd(models.Goal{
		ID:       "31",
		ParentID: &id3,
		Order:    1,
	})

	// Add subchildren
	id11 := "11"
	id21 := "21"
	id31 := "31"

	goalTree.TryAdd(models.Goal{
		ID:       "111",
		ParentID: &id11,
		Order:    1,
	})

	goalTree.TryAdd(models.Goal{
		ID:       "211",
		ParentID: &id21,
		Order:    1,
	})

	goalTree.TryAdd(models.Goal{
		ID:       "311",
		ParentID: &id31,
		Order:    1,
	})

	// Parent 404
	randomID := "404"
	goalTree.TryAdd(models.Goal{
		ID:       "4041",
		ParentID: &randomID,
	})

	slice := goalTree.ToSlice()

	assert.Equal(t, 3, len(slice))
	assert.Equal(t, "1", slice[0].ID)
	assert.Equal(t, "2", slice[1].ID)
	assert.Equal(t, "3", slice[2].ID)

	assert.Equal(t, "11", slice[0].SubGoals[0].ID)
	assert.Equal(t, "21", slice[1].SubGoals[0].ID)
	assert.Equal(t, "31", slice[2].SubGoals[0].ID)

	assert.Equal(t, "111", slice[0].SubGoals[0].SubGoals[0].ID)
	assert.Equal(t, "211", slice[1].SubGoals[0].SubGoals[0].ID)
	assert.Equal(t, "311", slice[2].SubGoals[0].SubGoals[0].ID)
}
