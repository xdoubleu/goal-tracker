package repositories

import (
	"context"

	"github.com/XDoubleU/essentia/pkg/database/postgres"

	"goal-tracker/api/internal/models"
)

type GamesRepository struct {
	db postgres.DB
}

func (repo GamesRepository) Fetch(
	ctx context.Context,
) ([]models.Game, error) {
	query := `
		SELECT app_id, name
		FROM games
	`

	rows, err := repo.db.Query(ctx, query)
	if err != nil {
		return nil, postgres.PgxErrorToHTTPError(err)
	}

	games := []models.Game{}

	for rows.Next() {
		//nolint:exhaustruct //other fields are assigned later
		game := models.Game{}

		err = rows.Scan(
			&game.AppID,
			&game.Name,
		)

		if err != nil {
			return nil, postgres.PgxErrorToHTTPError(err)
		}

		games = append(games, game)
	}

	if err = rows.Err(); err != nil {
		return nil, postgres.PgxErrorToHTTPError(err)
	}

	return games, nil
}

func (repo GamesRepository) Save(
	ctx context.Context,
	appId int64,
	name string,
) (*models.Game, error) {
	query := `
		INSERT INTO games (app_id, name)
		VALUES ($1, $2)
		RETURNING app_id
	`

	game := models.Game{
		AppID: appId,
		Name:  name,
	}

	err := repo.db.QueryRow(
		ctx,
		query,
		game.AppID,
		game.Name,
	).Scan(&game.Name)

	if err != nil {
		return nil, postgres.PgxErrorToHTTPError(err)
	}

	return &game, nil
}
