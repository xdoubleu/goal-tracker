package repositories

import (
	"context"

	"github.com/XDoubleU/essentia/pkg/database/postgres"

	"goal-tracker/api/internal/models"
)

type ProgressRepository struct {
	db postgres.DB
}

func (repo ProgressRepository) Create(
	ctx context.Context,
	goalID string,
	value int64,
) (*models.Progress, error) {
	query := `
		INSERT INTO progress (goal_id, value)
		VALUES ($1, $2)
		RETURNING id
	`

	//nolint:exhaustruct //other fields are optional
	progress := models.Progress{
		GoalID: goalID,
		Value:  value,
	}

	err := repo.db.QueryRow(
		ctx,
		query,
		goalID,
		value,
	).Scan(&progress.ID)

	if err != nil {
		return nil, postgres.PgxErrorToHTTPError(err)
	}

	return &progress, nil
}
