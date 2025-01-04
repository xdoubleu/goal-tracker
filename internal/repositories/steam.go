package repositories

import (
	"context"
	"time"

	"github.com/XDoubleU/essentia/pkg/database/postgres"

	"goal-tracker/api/internal/models"
)

type SteamRepository struct {
	db postgres.DB
}

func (repo SteamRepository) GetAllGames(
	ctx context.Context,
	userID string,
) ([]models.Game, error) {
	query := `
		SELECT id, name, is_delisted
		FROM steam_games
		WHERE user_id = $1
	`

	rows, err := repo.db.Query(ctx, query, userID)
	if err != nil {
		return nil, postgres.PgxErrorToHTTPError(err)
	}

	games := []models.Game{}
	for rows.Next() {
		var game models.Game

		err = rows.Scan(
			&game.ID,
			&game.Name,
			&game.IsDelisted,
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

func (repo SteamRepository) UpsertGame(
	ctx context.Context,
	id int,
	userID string,
	name string,
) error {
	query := `
		INSERT INTO steam_games (id, user_id, name)
		VALUES ($1, $2, $3)
		ON CONFLICT (id, user_id)
		DO UPDATE SET name = $3
	`

	_, err := repo.db.Exec(
		ctx,
		query,
		id,
		userID,
		name,
	)

	if err != nil {
		return postgres.PgxErrorToHTTPError(err)
	}

	return nil
}

func (repo SteamRepository) MarkGameAsDelisted(
	ctx context.Context,
	game *models.Game,
	userID string,
) error {
	query := `
		UPDATE steam_games
		SET is_delisted = true
		WHERE id = $1 AND user_id = $2
	`

	_, err := repo.db.Exec(
		ctx,
		query,
		game.ID,
		userID,
	)

	if err != nil {
		return postgres.PgxErrorToHTTPError(err)
	}

	return nil
}

func (repo SteamRepository) GetAchievementsForGame(
	ctx context.Context,
	gameID int,
	userID string,
) ([]models.Achievement, error) {
	query := `
		SELECT name, achieved, unlock_time
		FROM steam_achievements
		WHERE game_id = $1 AND user_id = $2
	`

	rows, err := repo.db.Query(ctx, query, gameID, userID)
	if err != nil {
		return nil, postgres.PgxErrorToHTTPError(err)
	}

	achievements := []models.Achievement{}
	for rows.Next() {
		//nolint:exhaustruct //other fields are defined later
		achievement := models.Achievement{
			GameID: gameID,
		}

		err = rows.Scan(
			&achievement.Name,
			&achievement.Achieved,
			&achievement.UnlockTime,
		)

		if err != nil {
			return nil, postgres.PgxErrorToHTTPError(err)
		}

		achievements = append(achievements, achievement)
	}

	if err = rows.Err(); err != nil {
		return nil, postgres.PgxErrorToHTTPError(err)
	}

	return achievements, nil
}

func (repo SteamRepository) UpsertAchievement(
	ctx context.Context,
	name string,
	userID string,
	gameID int,
	achieved bool,
	unlockTime *time.Time,
) error {
	query := `
		INSERT INTO steam_achievements (name, user_id, game_id, achieved, unlock_time)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (name, user_id, game_id)
		DO UPDATE SET achieved = $4, unlock_time = $5
	`

	_, err := repo.db.Exec(
		ctx,
		query,
		name,
		userID,
		gameID,
		achieved,
		unlockTime,
	)

	if err != nil {
		return postgres.PgxErrorToHTTPError(err)
	}

	return nil
}

func (repo SteamRepository) UpsertAchievementSchema(
	ctx context.Context,
	name string,
	userID string,
	gameID int,
) error {
	query := `
		INSERT INTO steam_achievements (name, user_id, game_id)
		VALUES ($1, $2, $3)
		ON CONFLICT (name, user_id, game_id)
		DO NOTHING
	`

	_, err := repo.db.Exec(
		ctx,
		query,
		name,
		userID,
		gameID,
	)

	if err != nil {
		return postgres.PgxErrorToHTTPError(err)
	}

	return nil
}
