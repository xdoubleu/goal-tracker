package repositories

import (
	"context"
	"time"

	"github.com/XDoubleU/essentia/pkg/database"
	"github.com/XDoubleU/essentia/pkg/database/postgres"

	"goal-tracker/api/internal/models"
)

type GoalRepository struct {
	db postgres.DB
}

func (repo GoalRepository) GetAll(
	ctx context.Context,
	userID string,
) ([]models.Goal, error) {
	query := `
		SELECT id, name, is_linked, target_value, type_id, state, due_time
		FROM goals
		WHERE user_id = $1
	`

	rows, err := repo.db.Query(ctx, query, userID)
	if err != nil {
		return nil, postgres.PgxErrorToHTTPError(err)
	}

	goals := []models.Goal{}
	for rows.Next() {
		//nolint:exhaustruct //other fields are initialized later
		goal := models.Goal{
			UserID: userID,
		}

		err = rows.Scan(
			&goal.ID,
			&goal.Name,
			&goal.IsLinked,
			&goal.TargetValue,
			&goal.TypeID,
			&goal.State,
			&goal.DueTime,
		)
		if err != nil {
			return nil, postgres.PgxErrorToHTTPError(err)
		}

		goals = append(goals, goal)
	}

	if err = rows.Err(); err != nil {
		return nil, postgres.PgxErrorToHTTPError(err)
	}

	return goals, nil
}

func (repo GoalRepository) GetByID(
	ctx context.Context,
	id string,
	userID string,
) (*models.Goal, error) {
	query := `
		SELECT name, target_value, source_id, type_id, state, due_time
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
		&goal.TargetValue,
		&goal.TypeID,
		&goal.State,
		&goal.DueTime,
	)
	if err != nil {
		return nil, postgres.PgxErrorToHTTPError(err)
	}

	return &goal, nil
}

func (repo GoalRepository) Create(
	ctx context.Context,
	id string,
	parentID *string,
	userID string,
	name string,
	isLinked bool,
	targetValue int64,
	typeID int64,
	state string,
	dueTime *time.Time,
) (*models.Goal, error) {
	query := `
		INSERT INTO goals (id, parent_id, user_id, name, 
			is_linked, target_value, type_id, state, due_time)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id
	`

	goal := models.Goal{
		ID:          id,
		ParentID:    parentID,
		UserID:      userID,
		Name:        name,
		IsLinked:    isLinked,
		TargetValue: &targetValue,
		TypeID:      &typeID,
		State:       state,
		DueTime:     dueTime,
	}

	err := repo.db.QueryRow(
		ctx,
		query,
		id,
		parentID,
		userID,
		name,
		isLinked,
		targetValue,
		typeID,
		state,
		dueTime,
	).Scan(&goal.ID)

	if err != nil {
		return nil, postgres.PgxErrorToHTTPError(err)
	}

	return &goal, nil
}

/*
func (repo GoalRepository) Update(
	ctx context.Context,
	goal models.Goal,
	updateGoalDto *dtos.UpdateGoalDto,
) (*models.Goal, error) {
	query := `
		UPDATE goals
		SET name = $3, description = $4, date = $5, value = $6, source_id = $7, type_id = $8,
		score = $9, state_id = $10
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
*/

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
