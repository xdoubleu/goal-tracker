package repositories

import (
	"context"

	"github.com/XDoubleU/essentia/pkg/database"
	"github.com/XDoubleU/essentia/pkg/database/postgres"

	"goal-tracker/api/internal/models"
)

type StateRepository struct {
	db postgres.DB
}

func (repo StateRepository) GetAll(
	ctx context.Context,
) ([]models.State, error) {
	query := `
		SELECT id, name, "order"
		FROM states
		ORDER BY "order"
	`

	rows, err := repo.db.Query(ctx, query)
	if err != nil {
		return nil, postgres.PgxErrorToHTTPError(err)
	}

	states := []models.State{}
	for rows.Next() {
		//nolint:exhaustruct //other fields are initialized later
		state := models.State{}

		err = rows.Scan(
			&state.ID,
			&state.Name,
			&state.Order,
		)
		if err != nil {
			return nil, postgres.PgxErrorToHTTPError(err)
		}

		states = append(states, state)
	}

	if err = rows.Err(); err != nil {
		return nil, postgres.PgxErrorToHTTPError(err)
	}

	return states, nil
}

func (repo StateRepository) Upsert(
	ctx context.Context,
	id string,
	name string,
	order int,
) (*models.State, error) {
	query := `
		INSERT INTO states (id, name, "order")
		VALUES ($1, $2, $3)
		ON CONFLICT (id)
		DO UPDATE SET name = $2, "order" = $3
		RETURNING id
	`

	state := models.State{
		ID:    id,
		Name:  name,
		Order: order,
	}

	err := repo.db.QueryRow(
		ctx,
		query,
		id,
		name,
		order,
	).Scan(&state.ID)

	if err != nil {
		return nil, postgres.PgxErrorToHTTPError(err)
	}

	return &state, nil
}

func (repo StateRepository) Delete(
	ctx context.Context,
	state *models.State,
) error {
	query := `
		DELETE FROM states
		WHERE id = $1
	`

	result, err := repo.db.Exec(ctx, query, state.ID)
	if err != nil {
		return postgres.PgxErrorToHTTPError(err)
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return database.ErrResourceNotFound
	}

	return nil
}
