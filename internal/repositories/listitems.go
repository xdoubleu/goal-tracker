package repositories

import (
	"context"

	"github.com/XDoubleU/essentia/pkg/database/postgres"

	"goal-tracker/api/internal/models"
)

type ListItemRepository struct {
	db postgres.DB
}

func (repo *ListItemRepository) GetByGoalID(
	ctx context.Context,
	goalID string,
	userID string,
) ([]models.ListItem, error) {
	query := `
		SELECT id, value, completed
		FROM list_items 
		WHERE goal_id = $1 AND user_id = $2
	`

	rows, err := repo.db.Query(ctx, query, goalID, userID)
	if err != nil {
		return nil, postgres.PgxErrorToHTTPError(err)
	}

	listItems := []models.ListItem{}
	for rows.Next() {
		//nolint:exhaustruct //other fields are assigned later
		listItem := models.ListItem{
			GoalID: goalID,
		}

		err = rows.Scan(
			&listItem.ID,
			&listItem.Value,
			&listItem.Completed,
		)

		if err != nil {
			return nil, postgres.PgxErrorToHTTPError(err)
		}

		listItems = append(listItems, listItem)
	}

	if err = rows.Err(); err != nil {
		return nil, postgres.PgxErrorToHTTPError(err)
	}

	return listItems, nil
}

func (repo *ListItemRepository) Upsert(
	ctx context.Context,
	id int64,
	userID string,
	goalID string,
	value string,
	completed bool,
) (*models.ListItem, error) {
	query := `
		INSERT INTO list_items (id, user_id, goal_id, value, completed)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (id, user_id, goal_id)
		DO UPDATE SET value = $4, completed = $5
		RETURNING id
	`

	listItem := models.ListItem{
		ID:        id,
		GoalID:    goalID,
		Value:     value,
		Completed: completed,
	}

	err := repo.db.QueryRow(
		ctx,
		query,
		id,
		userID,
		goalID,
		value,
		completed,
	).Scan(&listItem.ID)

	if err != nil {
		return nil, postgres.PgxErrorToHTTPError(err)
	}

	return &listItem, nil
}
