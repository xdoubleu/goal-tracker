package repositories

import (
	"context"
	"time"

	"github.com/XDoubleU/essentia/pkg/database"
	"github.com/XDoubleU/essentia/pkg/database/postgres"

	"goal-tracker/api/internal/dtos"
	"goal-tracker/api/internal/models"
	"goal-tracker/api/pkg/todoist"
)

type GoalRepository struct {
	db postgres.DB
}

func (repo GoalRepository) GetAll(
	ctx context.Context,
) ([]models.Goal, error) {
	query := `
		SELECT id, name, type_id, target_value, state_id, is_linked, parent_id, due_time, "order"
		FROM goals
		ORDER BY parent_id DESC
	`

	rows, err := repo.db.Query(ctx, query)
	if err != nil {
		return nil, postgres.PgxErrorToHTTPError(err)
	}

	goals := []models.Goal{}
	for rows.Next() {
		//nolint:exhaustruct //other fields are initialized later
		goal := models.Goal{}

		err = rows.Scan(
			&goal.ID,
			&goal.Name,
			&goal.TypeID,
			&goal.TargetValue,
			&goal.StateID,
			&goal.IsLinked,
			&goal.ParentID,
			&goal.DueTime,
			&goal.Order,
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
) (*models.Goal, error) {
	query := `
		SELECT name, type_id, target_value, state_id, is_linked, parent_id, due_time, "order"
		FROM goals
		WHERE goals.id = $1
	`

	//nolint:exhaustruct //other fields are optional
	goal := models.Goal{
		ID: id,
	}
	err := repo.db.QueryRow(
		ctx,
		query,
		id).Scan(
		&goal.Name,
		&goal.TypeID,
		&goal.TargetValue,
		&goal.StateID,
		&goal.IsLinked,
		&goal.ParentID,
		&goal.DueTime,
		&goal.Order,
	)
	if err != nil {
		return nil, postgres.PgxErrorToHTTPError(err)
	}

	return &goal, nil
}

func (repo GoalRepository) GetByTypeID(
	ctx context.Context,
	id int64,
) ([]models.Goal, error) {
	query := `
		SELECT id, name, target_value, state_id, is_linked, parent_id, due_time, "order"
		FROM goals
		WHERE goals.type_id = $1
	`

	rows, err := repo.db.Query(ctx, query, id)
	if err != nil {
		return nil, postgres.PgxErrorToHTTPError(err)
	}

	goals := []models.Goal{}

	for rows.Next() {
		//nolint:exhaustruct //other fields are assigned later
		goal := models.Goal{
			TypeID: &id,
		}

		err = rows.Scan(
			&goal.ID,
			&goal.Name,
			&goal.TargetValue,
			&goal.StateID,
			&goal.IsLinked,
			&goal.ParentID,
			&goal.DueTime,
			&goal.Order,
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

func (repo GoalRepository) Create(
	ctx context.Context,
	id string,
	parentID *string,
	name string,
	isLinked bool,
	targetValue *int64,
	typeID *int64,
	stateID string,
	due *todoist.Due,
	order int,
) (*models.Goal, error) {
	query := `
		INSERT INTO goals (id, parent_id, name, is_linked, target_value, 
		type_id, state_id, due_time, "order")
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id
	`

	var dueTime *time.Time
	if due != nil {
		dueTime = &due.Date.Time
	}

	goal := models.Goal{
		ID:          id,
		ParentID:    parentID,
		Name:        name,
		IsLinked:    isLinked,
		TargetValue: targetValue,
		TypeID:      typeID,
		StateID:     stateID,
		DueTime:     dueTime,
		Order:       order,
	}

	err := repo.db.QueryRow(
		ctx,
		query,
		id,
		parentID,
		name,
		isLinked,
		targetValue,
		typeID,
		stateID,
		dueTime,
		order,
	).Scan(&goal.ID)

	if err != nil {
		return nil, postgres.PgxErrorToHTTPError(err)
	}

	return &goal, nil
}

func (repo GoalRepository) Link(
	ctx context.Context,
	goal models.Goal,
	linkGoalDto dtos.LinkGoalDto,
) error {
	query := `
		UPDATE goals
		SET is_linked = true, target_value = $2, type_id = $3
		WHERE id = $1
	`

	result, err := repo.db.Exec(
		ctx,
		query,
		goal.ID,
		linkGoalDto.TargetValue,
		linkGoalDto.TypeID,
	)

	if err != nil {
		return postgres.PgxErrorToHTTPError(err)
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return database.ErrResourceNotFound
	}

	return nil
}

func (repo GoalRepository) Unlink(
	ctx context.Context,
	goal models.Goal,
) error {
	query := `
		UPDATE goals
		SET is_linked = false, target_value = null, type_id = null
		WHERE id = $1
	`

	result, err := repo.db.Exec(
		ctx,
		query,
		goal.ID,
	)

	if err != nil {
		return postgres.PgxErrorToHTTPError(err)
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return database.ErrResourceNotFound
	}

	return nil
}

func (repo GoalRepository) Update(
	ctx context.Context,
	goal models.Goal,
	task todoist.Task,
) (*models.Goal, error) {
	query := `
		UPDATE goals
		SET parent_id = $2, name = $3, state_id = $4, due_time = $5, "order" = $6
		WHERE id = $1
	`

	goal.ParentID = task.ParentID
	goal.Name = task.Content
	goal.StateID = task.SectionID
	goal.Order = task.Order

	if task.Due != nil {
		goal.DueTime = &task.Due.Date.Time
	}
	resultLocation, err := repo.db.Exec(
		ctx,
		query,
		goal.ID,
		goal.ParentID,
		goal.Name,
		goal.StateID,
		goal.DueTime,
		goal.Order,
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
		WHERE id = $1
	`

	result, err := repo.db.Exec(ctx, query, goal.ID)
	if err != nil {
		return postgres.PgxErrorToHTTPError(err)
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return database.ErrResourceNotFound
	}

	return nil
}
