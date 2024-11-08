package repositories

import (
	"context"
	"goal-tracker/api/internal/dtos"
	"goal-tracker/api/internal/models"
	"time"

	"github.com/XDoubleU/essentia/pkg/database"
	"github.com/XDoubleU/essentia/pkg/database/postgres"
	"github.com/jackc/pgx/v5"
)

type GoalRepository struct {
	db postgres.DB
}

func (repo GoalRepository) GetPage(ctx context.Context, userID string, pageSize int, getAfterID *string) ([]*models.Goal, error) {
	query := `
		SELECT id, name, description, date, value, source_id, type_id, score, state_id
		FROM goals
		WHERE user_id = $1
	`

	var rows pgx.Rows
	var err error
	if getAfterID != nil {
		query += " AND id > $2 ORDER BY score DESC LIMIT $3"
		rows, err = repo.db.Query(ctx, query, userID, getAfterID, pageSize)
	} else {
		query += " ORDER BY score DESC LIMIT $2"
		rows, err = repo.db.Query(ctx, query, userID, pageSize)
	}

	if err != nil {
		return nil, postgres.PgxErrorToHTTPError(err)
	}

	goals := []*models.Goal{}
	for rows.Next() {
		goal := models.Goal{
			UserID: userID,
		}

		err = rows.Scan(
			&goal.ID,
			&goal.Name,
			&goal.Description,
			&goal.Date,
			&goal.Value,
			&goal.SourceID,
			&goal.TypeID,
			&goal.Score,
			&goal.StateID,
		)
		if err != nil {
			return nil, postgres.PgxErrorToHTTPError(err)
		}

		goals = append(goals, &goal)
	}

	if err = rows.Err(); err != nil {
		return nil, postgres.PgxErrorToHTTPError(err)
	}

	return goals, nil
}

func (repo GoalRepository) GetByID(ctx context.Context, id string, userID string) (*models.Goal, error) {
	query := `
		SELECT name, description, date, value, source_id, type_id, score, state_id
		FROM goals
		WHERE goals.id = $1 AND user_id = $2
	`

	//nolint:exhaustruct //other fields are optional
	goal := models.Goal{
		ID:     id,
		UserID: userID,
	}
	err := repo.db.QueryRow(
		ctx,
		query,
		id, userID).Scan(
		&goal.Name,
		&goal.Description,
		&goal.Date,
		&goal.Value,
		&goal.SourceID,
		&goal.TypeID,
		&goal.Score,
		&goal.StateID,
	)
	if err != nil {
		return nil, postgres.PgxErrorToHTTPError(err)
	}

	return &goal, nil
}

func (repo GoalRepository) Create(
	ctx context.Context,
	userID string,
	name string,
	description *string,
	date *time.Time,
	value *int64,
	sourceID *int64,
	typeID *int64,
	score int64,
	stateID int64,
) (*models.Goal, error) {
	query := `
		INSERT INTO goals (user_id, name, description, date, value, source_id, type_id, score, state_id)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`

	//nolint:exhaustruct //other fields are optional
	goal := models.Goal{
		UserID:      userID,
		Name:        name,
		Description: description,
		Date:        date,
		Value:       value,
		SourceID:    sourceID,
		TypeID:      typeID,
		Score:       score,
		StateID:     stateID,
	}

	err := repo.db.QueryRow(
		ctx,
		query,
		userID,
		name,
		description,
		score,
		stateID,
	).Scan(&goal.ID)

	if err != nil {
		return nil, postgres.PgxErrorToHTTPError(err)
	}

	return &goal, nil
}

func (repo GoalRepository) Update(
	ctx context.Context,
	goal models.Goal,
	updateGoalDto *dtos.UpdateGoalDto,
) (*models.Goal, error) {
	query := `
		UPDATE goals
		SET name = $3, description = $4, date = $5, value = $6, source_id = $7, type_id = $8, score = $9, state_id = $10
		WHERE id = $1 AND user_id = $2
	`

	if updateGoalDto.Name != nil {
		goal.Name = *updateGoalDto.Name
	}

	if updateGoalDto.Description != nil {
		goal.Description = updateGoalDto.Description
	}

	if updateGoalDto.Date != nil {
		goal.Date = updateGoalDto.Date
	}

	if updateGoalDto.Value != nil {
		goal.Value = updateGoalDto.Value
	}

	if updateGoalDto.SourceID != nil {
		goal.SourceID = updateGoalDto.SourceID
	}

	if updateGoalDto.TypeID != nil {
		goal.TypeID = updateGoalDto.TypeID
	}

	if updateGoalDto.Score != nil {
		goal.Score = *updateGoalDto.Score
	}

	if updateGoalDto.StateID != nil {
		goal.StateID = *updateGoalDto.StateID
	}

	resultLocation, err := repo.db.Exec(
		ctx,
		query,
		goal.ID,
		goal.UserID,
		goal.Name,
		goal.Description,
		goal.Date,
		goal.Value,
		goal.SourceID,
		goal.TypeID,
		goal.Score,
		goal.StateID,
	)

	if err != nil {
		return nil, postgres.PgxErrorToHTTPError(err)
	}

	rowsAffected := resultLocation.RowsAffected()
	if rowsAffected == 0 {
		return nil, database.ErrResourceNotFound
	}

	return &goal, nil
}

func (repo GoalRepository) Delete(
	ctx context.Context,
	goal *models.Goal,
) error {
	query := `
		DELETE FROM goals
		WHERE id = $1 AND user_id = $2
	`

	result, err := repo.db.Exec(ctx, query, goal.ID, goal.UserID)
	if err != nil {
		return postgres.PgxErrorToHTTPError(err)
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return database.ErrResourceNotFound
	}

	return nil
}
